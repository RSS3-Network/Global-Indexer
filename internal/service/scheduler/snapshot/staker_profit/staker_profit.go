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
		// Query the latest epoch of the staker profit snapshots.
		snapshot, err := s.databaseClient.FindStakerProfitSnapshots(ctx, schema.StakerProfitSnapshotsQuery{Limit: lo.ToPtr(1)})
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("find staker profit snapshots", zap.Error(err))

			return
		}

		// Query the latest epoch of the epoch events.
		epochEvents, err := s.databaseClient.FindEpochs(ctx, 1, nil)
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("find epochs", zap.Error(err))

			return
		}

		var latestEpochSnapshot, latestEpochEvent uint64

		if len(snapshot) > 0 {
			latestEpochSnapshot = snapshot[0].EpochID
		}

		if len(epochEvents) > 0 {
			latestEpochEvent = epochEvents[0].ID
		}

		// Save the staker profit snapshots.
		if latestEpochSnapshot < latestEpochEvent {
			if err := s.saveStakerProfitSnapshots(ctx, latestEpochSnapshot, latestEpochEvent); err != nil {
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
		// Fetch the epoch items by the epoch id.
		epochItems, err := s.databaseClient.FindEpochTransactions(ctx, epochID, 1, nil)
		if err != nil {
			return fmt.Errorf("find epoch transactions: %w", err)
		}

		if len(epochItems) == 0 {
			continue
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

				if data.TotalChipValue.IsZero() {
					continue
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
					zap.L().Error("fetch min tokens to stake", zap.Error(err), zap.String("node", chip.Node.String()), zap.Uint64("block_number", currentEpoch.BlockNumber.Uint64()))

					return fmt.Errorf("fetch the min tokens to stake: %w", err)
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
