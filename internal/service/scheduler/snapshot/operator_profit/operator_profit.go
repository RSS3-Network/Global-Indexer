package operatorprofit

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	Name    = "operator_profit"
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
		snapshot, err := s.databaseClient.FindOperatorProfitSnapshots(ctx, schema.OperatorProfitSnapshotsQuery{Limit: lo.ToPtr(1)})
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("find staker profit snapshots", zap.Error(err))

			return
		}

		// Query the latest epoch of the epoch events.
		epochEvents, err := s.databaseClient.FindEpochs(ctx, &schema.FindEpochsQuery{Limit: lo.ToPtr(1)})
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
			if err := s.saveOperatorProfitSnapshots(ctx, latestEpochSnapshot, latestEpochEvent); err != nil {
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

func (s *server) saveOperatorProfitSnapshots(ctx context.Context, latestEpochSnapshot, latestEpochEvent uint64) error {
	// Query the array of Nodes.
	nodes, err := s.databaseClient.FindNodes(ctx, schema.FindNodesQuery{})
	if err != nil {
		zap.L().Error("find Nodes", zap.Error(err))

		return fmt.Errorf("find Nodes: %w", err)
	}

	for epochID := latestEpochSnapshot + 1; epochID <= latestEpochEvent; epochID++ {
		// Fetch the epoch items by the epoch id.
		epochItems, err := s.databaseClient.FindEpochTransactions(ctx, epochID, 1, nil)
		if err != nil {
			zap.L().Error("find epoch transactions", zap.Error(err))

			continue
		}

		if len(epochItems) == 0 {
			continue
		}

		var (
			mutex     sync.Mutex
			errorPool = pool.New().WithContext(ctx).WithMaxGoroutines(30).WithCancelOnError().WithFirstError()
			data      = make([]*schema.OperatorProfitSnapshot, 0, len(nodes))
		)

		for _, node := range nodes {
			node := node

			if node.Address == ethereum.AddressGenesis {
				continue
			}

			errorPool.Go(func(ctx context.Context) error {
				// Query the Node info from the staking contract.
				nodeInfo, err := s.stakingContract.GetNode(&bind.CallOpts{Context: ctx, BlockNumber: epochItems[0].BlockNumber}, node.Address)
				if err != nil {
					zap.L().Error("get Node from rpc", zap.Error(err))

					return fmt.Errorf("get Node from rpc: %w", err)
				}

				if nodeInfo.Account == ethereum.AddressGenesis {
					return nil
				}

				mutex.Lock()
				defer mutex.Unlock()

				data = append(data, &schema.OperatorProfitSnapshot{
					Date:          time.Unix(epochItems[0].BlockTimestamp, 0),
					EpochID:       epochID,
					Operator:      nodeInfo.Account,
					OperationPool: decimal.NewFromBigInt(nodeInfo.OperationPoolTokens, 0),
				})

				return nil
			})
		}

		if err := errorPool.Wait(); err != nil {
			return fmt.Errorf("fetch operator profit: %w", err)
		}

		if err := s.databaseClient.SaveOperatorProfitSnapshots(ctx, data); err != nil {
			return fmt.Errorf("save Node min tokens to stake snapshots: %w", err)
		}
	}

	return nil
}

func New(databaseClient database.Client, redisClient *redis.Client, stakingContract *stakingv2.Staking) service.Server {
	return &server{
		cronJob:         cronjob.New(redisClient, Name, Timeout),
		databaseClient:  databaseClient,
		redisClient:     redisClient,
		stakingContract: stakingContract,
	}
}
