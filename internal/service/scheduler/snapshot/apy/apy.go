package apy

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/cronjob"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	Name                    = "apy"
	Timeout                 = 3 * time.Minute
	CacheKeyEpochAverageAPY = "epoch_average_apy"
)

var _ service.Server = (*server)(nil)

type server struct {
	cronJob         *cronjob.CronJob
	databaseClient  database.Client
	cacheClient     cache.Client
	stakingContract *l2.Staking
}

func (s *server) Name() string {
	return Name
}

func (s *server) Spec() string {
	return "0 */1 * * * *" // every minute
}

func (s *server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		// Query the latest of the epoch apy snapshots.
		snapshots, err := s.databaseClient.FindEpochAPYSnapshots(ctx, schema.EpochAPYSnapshotQuery{Limit: lo.ToPtr(1)})
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("find epoch APY snapshots", zap.Error(err))

			return
		}

		// Query the latest epoch of the epoch events.
		epochEvents, err := s.databaseClient.FindEpochs(ctx, 1, nil)
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("find epochs", zap.Error(err))

			return
		}

		var latestSnapshotEpochID uint64

		if len(snapshots) > 0 {
			latestSnapshotEpochID = snapshots[0].EpochID
		}

		// Save the minTokensToStake snapshots.
		if latestSnapshotEpochID < epochEvents[0].ID {
			if err := s.saveAPYToSnapshots(ctx, latestSnapshotEpochID, epochEvents[0]); err != nil {
				zap.L().Error("save APY to snapshots", zap.Error(err))

				return
			}
		}
	})
	if err != nil {
		return fmt.Errorf("add apy cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	return nil
}

func (s *server) saveAPYToSnapshots(ctx context.Context, latestEpochSnapshot uint64, latestEpochEvent *schema.Epoch) error {
	for id := latestEpochSnapshot + 1; id <= latestEpochEvent.ID; id++ {
		// Query the epoch transactions by the epoch id.
		transactions, err := s.databaseClient.FindEpochTransactions(ctx, id, latestEpochEvent.TotalRewardedNodes, nil)
		if err != nil {
			return fmt.Errorf("find epoch transactions: %w", err)
		}

		if len(transactions) == 0 {
			continue
		}

		var (
			nodeAPYSnapshots = make([]*schema.NodeAPYSnapshot, 0)
			sum              = decimal.NewFromInt(0)
		)

		for _, transaction := range transactions {
			for _, item := range transaction.RewardedNodes {
				node, err := s.stakingContract.GetNode(&bind.CallOpts{BlockNumber: transaction.BlockNumber}, item.NodeAddress)
				if err != nil {
					zap.L().Error("get node from rpc", zap.Error(err), zap.String("nodeAddress", item.NodeAddress.String()), zap.Any("blockNumber", transaction.BlockNumber))

					return fmt.Errorf("get node from rpc: %w", err)
				}

				// Calculate the APY.
				// APY = (operationRewards + stakingRewards) / (stakingPoolTokens) * (1 - tax) * number of epochs in a year
				// number of epochs in a year = 365 * 24 / 18 = 486.6666666666667
				if node.StakingPoolTokens.Cmp(big.NewInt(0)) > 0 {
					tax := 1 - float64(node.TaxRateBasisPoints)/10000

					apy := item.OperationRewards.Add(item.StakingRewards).
						Div(decimal.NewFromBigInt(node.StakingPoolTokens, 0)).
						Mul(decimal.NewFromFloat(tax)).
						Mul(decimal.NewFromFloat(486.6666666666667))

					nodeAPYSnapshots = append(nodeAPYSnapshots, &schema.NodeAPYSnapshot{
						Date:        time.Unix(transaction.EndTimestamp, 0),
						EpochID:     id,
						NodeAddress: item.NodeAddress,
						APY:         apy,
					})

					sum = sum.Add(apy)
				}
			}
		}

		zap.L().Info("save APY to snapshots", zap.Uint64("epochID", id), zap.Int("nodeAPYSnapshots", len(nodeAPYSnapshots)))

		// Save the node APY snapshots.
		if err := s.databaseClient.SaveNodeAPYSnapshots(ctx, nodeAPYSnapshots); err != nil {
			return fmt.Errorf("save node APY snapshots: %w", err)
		}

		epochAPYSnapshot := schema.EpochAPYSnapshot{
			Date:    time.Unix(transactions[0].EndTimestamp, 0),
			EpochID: id,
			APY:     sum.Div(decimal.NewFromInt(int64(len(nodeAPYSnapshots)))),
		}

		// Save the epoch APY snapshot.
		if err := s.databaseClient.SaveEpochAPYSnapshot(ctx, &epochAPYSnapshot); err != nil {
			return fmt.Errorf("save epoch APY snapshot: %w", err)
		}
	}

	// Save the epoch average APY to cache.
	apy, err := s.databaseClient.FindEpochAPYSnapshotsAverage(ctx)
	if err != nil {
		return fmt.Errorf("find epoch APY snapshots average: %w", err)
	}

	if err := s.cacheClient.Set(ctx, CacheKeyEpochAverageAPY, apy); err != nil {
		return fmt.Errorf("set epoch average APY to cache: %w", err)
	}

	return nil
}

func New(databaseClient database.Client, redisClient *redis.Client, stakingContract *l2.Staking) service.Server {
	return &server{
		cronJob:         cronjob.New(redisClient, Name, Timeout),
		cacheClient:     cache.New(redisClient),
		databaseClient:  databaseClient,
		stakingContract: stakingContract,
	}
}
