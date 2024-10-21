package stakercumulativeearnings

import (
	"context"
	"encoding/json"
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
	Name    = "staker_cumulative_earnings"
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

		// Find the latest epoch of the staker cumulative earnings snapshots.
		latestEpochSnapshot, err := s.findLatestStakerCumulativeEarningsSnapshot(ctx, epochEvents[0].ID)
		if err != nil {
			zap.L().Error("find latest staker cumulative earnings snapshot", zap.Error(err))

			return
		}

		// Save the staker cumulative earnings snapshots.
		if latestEpochSnapshot < epochEvents[0].ID {
			if err := s.saveStakerCumulativeEarningsSnapshots(ctx, latestEpochSnapshot, epochEvents[0].ID); err != nil {
				zap.L().Error("save staker cumulative earnings snapshots", zap.Error(err))
			}
		}
	})

	if err != nil {
		return fmt.Errorf("add staker_cumulative_earnings cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	return nil
}

func (s *server) findLatestStakerCumulativeEarningsSnapshot(ctx context.Context, latestEpochID uint64) (uint64, error) {
	// Check the latest epoch of the staker cumulative earning snapshots.
	for epochID := latestEpochID; epochID > 0; epochID-- {
		// Query the epoch of the staker cumulative earnings snapshots.
		snapshots, err := s.databaseClient.FindStakerCumulativeEarningSnapshots(ctx, schema.StakerCumulativeEarningSnapshotsQuery{EpochID: lo.ToPtr(epochID)})
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("find staker cumulative earnings snapshots", zap.Error(err))

			return 0, err
		}

		if len(snapshots) == 0 {
			continue
		}

		epochItems, err := s.databaseClient.FindEpochTransactions(ctx, epochID, 1, nil)
		if err != nil {
			zap.L().Error("find epoch transactions", zap.Error(err))

			return 0, err
		}

		stakerCount, err := s.databaseClient.FindStakerCount(ctx, schema.StakeChipsQuery{
			BlockNumber: epochItems[0].BlockNumber,
		})
		if err != nil {
			zap.L().Error("find staker count", zap.Error(err))

			return 0, err
		}

		if int64(len(snapshots)) < stakerCount {
			continue
		}

		return epochID, nil
	}

	return 0, nil
}

func (s *server) saveStakerCumulativeEarningsSnapshots(ctx context.Context, latestEpochSnapshot, latestEpochEvent uint64) error {
	// Iterate the epoch id from the latest epoch snapshot to the latest epoch event.
	for epochID := latestEpochSnapshot + 1; epochID <= latestEpochEvent; epochID++ {
		if err := s.saveStakerCumulativeEarningsSnapshotsByEpochID(ctx, epochID); err != nil {
			return fmt.Errorf("save staker profit snapshots by epoch id: %w", err)
		}
	}

	return nil
}

// saveStakerCumulativeEarningsSnapshotsByEpochID saves the staker cumulative earnings snapshots by the epoch id.
// Cumulative earnings = last epoch cumulative earnings + current chip value changes.
// But it is worth noting that merged chips need to be handled specially.
func (s *server) saveStakerCumulativeEarningsSnapshotsByEpochID(ctx context.Context, epochID uint64) error {
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

		snapshots := make([]*schema.StakerCumulativeEarningSnapshot, 0, len(stakers))

		// Fetch the chips by the stakers.
		for _, staker := range stakers {
			staker := staker

			if staker.Owner == ethereum.AddressGenesis {
				continue
			}

			// Query the staker profit snapshots by the owner address and the epoch id.
			exist, _ := s.databaseClient.FindStakerCumulativeEarningSnapshots(ctx, schema.StakerCumulativeEarningSnapshotsQuery{
				OwnerAddress: lo.ToPtr(staker.Owner),
				EpochID:      lo.ToPtr(epochID),
				Limit:        lo.ToPtr(1),
			})
			if len(exist) > 0 {
				continue
			}

			data, err := s.buildSaveStakerCumulativeEarningsSnapshot(ctx, epochItems[0], staker.Owner)
			if err != nil {
				return fmt.Errorf("build staker profit snapshots: %w", err)
			}

			snapshots = append(snapshots, data)
		}

		// Save the staker profit snapshots.
		if len(snapshots) > 0 {
			if err := s.databaseClient.SaveStakerCumulativeEarningSnapshots(ctx, snapshots); err != nil {
				return fmt.Errorf("save staker profit snapshots: %w", err)
			}
		}

		cursor = stakers[len(stakers)-1].ID
	}

	return nil
}

func (s *server) buildSaveStakerCumulativeEarningsSnapshot(ctx context.Context, currentEpoch *schema.Epoch, staker common.Address) (*schema.StakerCumulativeEarningSnapshot, error) {
	var chipsCursor *big.Int

	profit := &schema.StakerCumulativeEarningSnapshot{
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
				// Query the previous profit by the owner address and the epoch id.
				previousProfit, err := s.databaseClient.FindStakerCumulativeEarningSnapshots(ctx, schema.StakerCumulativeEarningSnapshotsQuery{
					OwnerAddress: lo.ToPtr(staker),
					EpochID:      lo.ToPtr(currentEpoch.ID - 1),
				})
				if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
					zap.L().Error("find staker cumulative earning snapshots", zap.Error(err))

					return fmt.Errorf("find staker cumulative earning snapshots: %w", err)
				}

				if len(previousProfit) > 0 {
					profit.CumulativeEarning = previousProfit[0].CumulativeEarning
				}

				// Query the chip value by the current epoch.
				chipInfo, err := s.stakingContract.GetChipInfo(&bind.CallOpts{Context: ctx, BlockNumber: currentEpoch.BlockNumber}, chip.ID)
				if err != nil {
					zap.L().Error("get chip info from chain", zap.Error(err), zap.String("chipID", chip.ID.String()), zap.Uint64("blockNumber", currentEpoch.BlockNumber.Uint64()))

					return fmt.Errorf("get chip info from chain: %w", err)
				}

				previousChipTokens, err := s.calculatePreviousChipTokens(ctx, staker, currentEpoch.BlockNumber, chip.ID)
				if err != nil {
					return fmt.Errorf("calculate previous chip tokens: %w", err)
				}

				if previousChipTokens != nil {
					mutex.Lock()
					profit.CumulativeEarning = profit.CumulativeEarning.Add(decimal.NewFromBigInt(new(big.Int).Sub(chipInfo.Tokens, previousChipTokens), 0))
					mutex.Unlock()
				}

				return nil
			})
		}

		if err := errorPool.Wait(); err != nil {
			return nil, fmt.Errorf("parallel processing the chips: %w", err)
		}

		chipsCursor = chips[len(chips)-1].ID
	}

	return profit, nil
}

func (s *server) calculatePreviousChipTokens(ctx context.Context, staker common.Address, blockNumber *big.Int, chipID *big.Int) (*big.Int, error) {
	var previousChipTokens *big.Int

	// Query the chip value by the current block number  - 1.
	previousChipInfo, err := s.stakingContract.GetChipInfo(&bind.CallOpts{Context: ctx, BlockNumber: new(big.Int).Sub(blockNumber, big.NewInt(1))}, chipID)
	if err != nil {
		zap.L().Error("get chip info from chain", zap.Error(err))

		return nil, fmt.Errorf("get chip info from chain: %w", err)
	}

	if previousChipInfo.Tokens != nil && previousChipInfo.Tokens.Cmp(big.NewInt(0)) != 0 {
		return previousChipInfo.Tokens, nil
	}

	// If the previous chip value is nil or zero, it is possible that the chip is merged.
	mergedChips, err := s.databaseClient.FindStakeTransactions(ctx, schema.StakeTransactionsQuery{
		User:        lo.ToPtr(staker),
		BlockNumber: lo.ToPtr(blockNumber.Uint64()),
	})
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		zap.L().Error("find stake transactions", zap.Error(err))

		return nil, fmt.Errorf("find stake transactions: %w", err)
	}

	if len(mergedChips) > 0 {
		for _, mergedChip := range mergedChips {
			var exist bool

			for _, mergedChipID := range mergedChip.ChipIDs {
				if mergedChipID.Cmp(chipID) == 0 {
					exist = true

					break
				}
			}

			if exist {
				events, err := s.databaseClient.FindStakeEvents(ctx, schema.StakeEventsQuery{
					IDs: []common.Hash{mergedChip.ID},
				})
				if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
					zap.L().Error("find stake events", zap.Error(err))

					return nil, fmt.Errorf("find stake events: %w", err)
				}

				for _, event := range events {
					var metadata schema.StakeEventChipsMergedMetadata

					if err := json.Unmarshal(event.Metadata, &metadata); err != nil {
						return nil, fmt.Errorf("unmarshal stake event metadata: %w", err)
					}

					if metadata.NewTokenID.Cmp(chipID) == 0 {
						for _, burnedTokenID := range metadata.BurnedTokenIDs {
							chipInfo, err := s.stakingContract.GetChipInfo(&bind.CallOpts{Context: ctx, BlockNumber: new(big.Int).Sub(blockNumber, big.NewInt(1))}, burnedTokenID)
							if err != nil {
								zap.L().Error("get chip info from chain", zap.Error(err))

								return nil, fmt.Errorf("get chip info from chain: %w", err)
							}

							previousChipTokens = new(big.Int).Add(previousChipTokens, chipInfo.Tokens)
						}

						return previousChipTokens, nil
					}
				}
			}
		}
	}

	return nil, nil
}

func New(databaseClient database.Client, redisClient *redis.Client, stakingContract *stakingv2.Staking) service.Server {
	return &server{
		cronJob:         cronjob.New(redisClient, Name, Timeout),
		databaseClient:  databaseClient,
		redisClient:     redisClient,
		stakingContract: stakingContract,
	}
}
