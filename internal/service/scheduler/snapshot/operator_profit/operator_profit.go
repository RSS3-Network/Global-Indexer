package operationpool

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
	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/ethereum"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cronjob"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

var ( // FIXME: Name should be "operation_pool", update the file naem too.
	Name    = "operator_profit"
	Timeout = 3 * time.Minute
)

var _ service.Server = (*server)(nil)

type server struct {
	cronJob         *cronjob.CronJob
	databaseClient  database.Client
	redisClient     *redis.Client
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
		// Query the latest Operation Pool snapshot.
		snapshot, err := s.databaseClient.FindOperationPoolSnapshots(ctx, schema.OperationPoolSnapshotsQuery{Limit: lo.ToPtr(1)})
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("find Operation Pool snapshot", zap.Error(err))

			return
		}

		// Query the latest epoch.
		epoch, err := s.databaseClient.FindEpochs(ctx, 1, nil)
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("find epoch", zap.Error(err))

			return
		}

		// Assign Epoch Ids based on the retrieved snapshot and epoch.
		var latestSnapshotEpochID, latestEpochID uint64

		if len(snapshot) > 0 {
			latestSnapshotEpochID = snapshot[0].EpochID
		}

		if len(epoch) > 0 {
			latestEpochID = epoch[0].ID
		}

		// Only begin the snapshot process if the latest snapshot is behind the latest epoch.
		if latestSnapshotEpochID < latestEpochID {
			if err := s.beginSnapshot(ctx, latestSnapshotEpochID, latestEpochID); err != nil {
				zap.L().Error("save Operation Pool snapshot", zap.Error(err))

				return
			}
		}
	})

	if err != nil {
		return fmt.Errorf("add Operation Pool snapshot cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	return nil
}

// beginSnapshot takes new snapshots of all Nodes' operation pool up to the latest epoch.
func (s *server) beginSnapshot(ctx context.Context, currentSnapshotEpochID, latestEpochID uint64) error {
	// Query the array of Nodes.
	nodes, err := s.getNodesFromDB(ctx)
	if err != nil {
		return fmt.Errorf("get nodes: %w", err)
	}

	// Iterate until the snapshot is up to date with the latest epoch.
	// currentEpochID is the epoch being snapshotted.
	for currentEpochID := currentSnapshotEpochID + 1; currentEpochID <= latestEpochID; currentEpochID++ {
		// Fetch the epoch by currentEpochID.
		epoch, err := s.databaseClient.FindEpochTransactions(ctx, currentEpochID, 1, nil)
		if err != nil {
			zap.L().Error("find epoch transactions", zap.Error(err))

			continue
		}

		// If the epoch does not exist in the database, log an error and continue.
		// This means epochs are not indexed up to the latest epoch.
		if len(epoch) == 0 {
			zap.L().Error("an epoch does not exist in database", zap.Any("epoch ID", currentEpochID))

			continue
		}

		// If there are no nodes, continue to the next epoch.
		if len(nodes) == 0 {
			continue
		}

		var (
			errorPool = pool.New().WithContext(ctx).WithMaxGoroutines(30).WithCancelOnError().WithFirstError()
			result    = make([]*schema.OperationPoolSnapshot, len(nodes))
		)

		for i, node := range nodes {
			errorPool.Go(func(ctx context.Context) error {
				// Query the Node info from the VSL.
				nodeInfo, err := s.getNodeInfoFromVSL(ctx, epoch[0].BlockNumber, node.Address)
				if err != nil {
					return err
				}

				// should not include genesis account
				if nodeInfo.Account == ethereum.AddressGenesis {
					return nil
				}

				result[i] = &schema.OperationPoolSnapshot{
					Date:          time.Unix(epoch[0].BlockTimestamp, 0),
					EpochID:       currentEpochID,
					Operator:      nodeInfo.Account,
					OperationPool: decimal.NewFromBigInt(nodeInfo.OperationPoolTokens, 0),
				}

				return nil
			})
		}

		if err := errorPool.Wait(); err != nil {
			return fmt.Errorf("fetch operator profit: %w", err)
		}

		// Filter out nil values in the result.
		result = lo.FilterMap(result, func(snapshot *schema.OperationPoolSnapshot, _ int) (*schema.OperationPoolSnapshot, bool) {
			return snapshot, snapshot != nil
		})

		// Save snapshots into the database.
		if err := s.databaseClient.SaveOperationPoolSnapshots(ctx, result); err != nil {
			return fmt.Errorf("save Operation Pool: %w", err)
		}
	}

	return nil
}

func (s *server) getNodesFromDB(ctx context.Context) ([]*schema.Node, error) {
	nodes, err := s.databaseClient.FindNodes(ctx, schema.FindNodesQuery{})

	if err != nil {
		return nil, fmt.Errorf("find nodes from DB: %w", err)
	}

	return nodes, nil
}

func (s *server) getNodeInfoFromVSL(ctx context.Context, blockNumber *big.Int, nodeAddress common.Address) (*l2.DataTypesNode, error) {
	nodeInfo, err := s.stakingContract.GetNode(&bind.CallOpts{Context: ctx, BlockNumber: blockNumber}, nodeAddress)
	if err != nil {
		msg := "get node from VSL error"
		zap.L().Error(msg, zap.Error(err))

		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	return &nodeInfo, nil
}

func New(databaseClient database.Client, redisClient *redis.Client, stakingContract *l2.Staking) service.Server {
	return &server{
		cronJob:         cronjob.New(redisClient, Name, Timeout),
		databaseClient:  databaseClient,
		redisClient:     redisClient,
		stakingContract: stakingContract,
	}
}
