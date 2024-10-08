package stakerprofit

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/ethereum"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/cronjob"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

var (
	Name    = "staker_profit"
	Timeout = 3 * time.Minute
)

var _ service.Server = (*server)(nil)

type server struct {
	cronJob         *cronjob.CronJob
	databaseClient  database.Client
	redisClient     *redis.Client
	stakingContract *stakingv2.Staking
}

func (s *server) Name() string {
	return Name
}

func (s *server) Spec() string {
	return "0 */1 * * * *" // every minute
}

func (s *server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		// Query the latest epoch of the epoch events.
		epochEvents, err := s.databaseClient.FindEpochs(ctx, &schema.FindEpochsQuery{Limit: lo.ToPtr(1)})
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("find epochs", zap.Error(err))

			return
		}

		if len(epochEvents) == 0 {
			return
		}

		var latestEpochSnapshot uint64

		// Check the latest epoch of the staker profit snapshots.
		for epochID := epochEvents[0].ID; epochID > 0; epochID-- {
			// Query the epoch of the staker profit snapshots.
			snapshots, err := s.databaseClient.FindStakerProfitSnapshots(ctx, schema.StakerProfitSnapshotsQuery{EpochID: lo.ToPtr(epochID)})
			if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
				zap.L().Error("find staker profit snapshots", zap.Error(err))

				return
			}

			if len(snapshots) == 0 {
				continue
			}

			epochItems, err := s.databaseClient.FindEpochTransactions(ctx, epochID, 1, nil)
			if err != nil {
				zap.L().Error("find epoch transactions", zap.Error(err))

				return
			}

			stakerCount, err := s.databaseClient.FindStakerCount(ctx, schema.StakeChipsQuery{
				BlockNumber: epochItems[0].BlockNumber,
			})
			if err != nil {
				zap.L().Error("find staker count", zap.Error(err))

				return
			}

			if int64(len(snapshots)) < stakerCount {
				continue
			}

			latestEpochSnapshot = epochID

			break
		}

		// Save the staker profit snapshots.
		if latestEpochSnapshot < epochEvents[0].ID {
			if err := s.saveStakerProfitSnapshots(ctx, latestEpochSnapshot, epochEvents[0].ID); err != nil {
				zap.L().Error("save staker profit snapshots", zap.Error(err))

				return
			}
		}
	})

	if err != nil {
		return fmt.Errorf("add staker profit cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	return nil
}

func (s *server) saveStakerProfitSnapshots(ctx context.Context, latestEpochSnapshot, latestEpochEvent uint64) error {
	// Iterate the epoch id from the latest epoch snapshot to the latest epoch event.
	for epochID := latestEpochSnapshot + 1; epochID <= latestEpochEvent; epochID++ {
		if err := s.saveStakerProfitSnapshotsByEpochID(ctx, epochID); err != nil {
			return fmt.Errorf("save staker profit snapshots by epoch id: %w", err)
		}
	}

	return nil
}

func (s *server) saveStakerProfitSnapshotsByEpochID(ctx context.Context, epochID uint64) error {
	// Fetch the epoch items by the epoch id.
	epochItems, err := s.databaseClient.FindEpochTransactions(ctx, epochID, 1, nil)
	if err != nil {
		return fmt.Errorf("find epoch transactions: %w", err)
	}

	if len(epochItems) == 0 {
		return nil
	}

	var cursor *big.Int

	for {
		// Fetch the distinct stakers from the chips table.
		findStakeChips := schema.StakeChipsQuery{
			Cursor:        cursor,
			Limit:         lo.ToPtr(500),
			DistinctOwner: true,
			BlockNumber:   epochItems[0].BlockNumber,
		}

		stakers, err := s.databaseClient.FindStakeChips(ctx, findStakeChips)
		if errors.Is(err, database.ErrorRowNotFound) || len(stakers) == 0 {
			break
		}

		if err != nil {
			return fmt.Errorf("find stake chips: %w", err)
		}

		snapshots := make([]*schema.StakerProfitSnapshot, 0, len(stakers))

		// Fetch the chips by the stakers.
		for _, staker := range stakers {
			staker := staker

			if staker.Owner == ethereum.AddressGenesis {
				continue
			}

			// Query the staker profit snapshots by the owner address and the epoch id.
			exist, _ := s.databaseClient.FindStakerProfitSnapshots(ctx, schema.StakerProfitSnapshotsQuery{
				OwnerAddress: lo.ToPtr(staker.Owner),
				EpochID:      lo.ToPtr(epochID),
				Limit:        lo.ToPtr(1),
			})
			if len(exist) > 0 {
				continue
			}

			data, err := s.buildStakerProfitSnapshots(ctx, epochItems[0], staker.Owner)
			if err != nil {
				return fmt.Errorf("build staker profit snapshots: %w", err)
			}

			snapshots = append(snapshots, data)
		}

		// Save the staker profit snapshots.
		if len(snapshots) > 0 {
			if err := s.databaseClient.SaveStakerProfitSnapshots(ctx, snapshots); err != nil {
				return fmt.Errorf("save staker profit snapshots: %w", err)
			}
		}

		cursor = stakers[len(stakers)-1].ID
	}

	return nil
}

func (s *server) buildStakerProfitSnapshots(ctx context.Context, currentEpoch *schema.Epoch, staker common.Address) (*schema.StakerProfitSnapshot, error) {
	var chipsCursor *big.Int

	profit := &schema.StakerProfitSnapshot{
		EpochID:      currentEpoch.ID,
		OwnerAddress: staker,
		Date:         time.Unix(currentEpoch.BlockTimestamp, 0),
	}

	for {
		findStakeChips := schema.StakeChipsQuery{
			Owner:       lo.ToPtr(staker),
			Cursor:      chipsCursor,
			BlockNumber: currentEpoch.BlockNumber,
			Limit:       lo.ToPtr(200),
		}

		// Fetch the chips by the stakers.
		chips, err := s.databaseClient.FindStakeChips(ctx, findStakeChips)
		if errors.Is(err, database.ErrorRowNotFound) || len(chips) == 0 {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("find stake chips: %w", err)
		}

		var mutex sync.Mutex

		// Parallel processing the chips.
		errorPool := pool.New().WithContext(ctx).WithMaxGoroutines(30).WithCancelOnError().WithFirstError()

		for _, chip := range chips {
			chip := chip

			errorPool.Go(func(ctx context.Context) error {
				// Query the chip value from the staking contract.
				chipInfo, err := s.stakingContract.GetChipInfo(&bind.CallOpts{Context: ctx, BlockNumber: currentEpoch.BlockNumber}, chip.ID)
				if err != nil {
					zap.L().Error("get chip info from chain", zap.Error(err))

					return fmt.Errorf("get chip info from chain: %w", err)
				}

				chip.Value = decimal.NewFromBigInt(chipInfo.Tokens, 0)

				mutex.Lock()
				profit.TotalChipValue = profit.TotalChipValue.Add(chip.Value)
				mutex.Unlock()

				return nil
			})
		}

		if err := errorPool.Wait(); err != nil {
			return nil, fmt.Errorf("parallel processing the chips: %w", err)
		}

		profit.TotalChipAmount = profit.TotalChipAmount.Add(decimal.NewFromInt(int64(len(chips))))

		chipsCursor = chips[len(chips)-1].ID
	}

	return profit, nil
}

func New(databaseClient database.Client, redisClient *redis.Client, stakingContract *stakingv2.Staking) service.Server {
	return &server{
		cronJob:         cronjob.New(redisClient, Name, Timeout),
		databaseClient:  databaseClient,
		redisClient:     redisClient,
		stakingContract: stakingContract,
	}
}
