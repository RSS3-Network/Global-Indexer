package nodemintokenstostake

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cronjob"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	Name    = "node_min_tokens_to_stake"
	Timeout = 3 * time.Minute
)

var _ service.Server = (*server)(nil)

type server struct {
	cronJob         *cronjob.CronJob
	databaseClient  database.Client
	redisClient     *redis.Client
	stakingContract *l2.Staking
}

func (s *server) Spec() string {
	return "0 */1 * * * *" // every minute
}

func (s *server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		// Query the latest epoch of the minToknsToStake snapshots.
		snapshot, err := s.databaseClient.FindNodeMinTokensToStakeSnapshots(ctx, nil, false, lo.ToPtr(1))
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("find node min tokens to stake snapshots", zap.Error(err))

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

		// Save the minTokensToStake snapshots.
		if latestEpochSnapshot < latestEpochEvent {
			if err := s.saveMinTokensToStakeSnapshots(ctx, latestEpochSnapshot, latestEpochEvent); err != nil {
				zap.L().Error("save min tokens to stake snapshots", zap.Error(err))

				return
			}
		}
	})
	if err != nil {
		return fmt.Errorf("add node min tokens to stake cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	return nil
}

func (s *server) saveMinTokensToStakeSnapshots(ctx context.Context, latestEpochSnapshot, latestEpochEvent uint64) error {
	for id := latestEpochSnapshot + 1; id <= latestEpochEvent; id++ {
		// Query the epoch items by the epoch id.
		epochItems, err := s.databaseClient.FindEpochTransactions(ctx, id, 1, nil)
		if err != nil {
			return fmt.Errorf("find epoch transactions: %w", err)
		}

		if len(epochItems) == 0 {
			continue
		}

		// Query the nodes from node table.
		var cursor *string

		for {
			nodes, err := s.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
				Cursor: cursor,
				Limit:  lo.ToPtr(1000),
			})
			if errors.Is(err, database.ErrorRowNotFound) || len(nodes) == 0 {
				break
			}

			if err != nil {
				return fmt.Errorf("find nodes: %w", err)
			}

			snapshots := make([]*schema.NodeMinTokensToStakeSnapshot, 0, len(nodes))

			for _, node := range nodes {
				// Query the min tokens to stake from the staking contract.
				if id < latestEpochEvent {
					minTokensToStake, err := s.stakingContract.MinTokensToStake(&bind.CallOpts{Context: ctx, BlockNumber: epochItems[0].BlockNumber}, node.Address)
					if err != nil {
						zap.L().Error("get min tokens to stake", zap.Error(err), zap.String("nodeAddress", node.Address.String()), zap.Any("blockNumber", epochItems[0].BlockNumber))

						return fmt.Errorf("get min tokens to stake: %w", err)
					}

					node.MinTokensToStake = decimal.NewFromBigInt(minTokensToStake, 0)
				}

				if node.MinTokensToStake.IsZero() {
					zap.L().Info("min tokens to stake is zero", zap.String("nodeAddress", node.Address.String()), zap.Uint64("epochID", id))

					continue
				}

				snapshots = append(snapshots, &schema.NodeMinTokensToStakeSnapshot{
					Date:             time.Unix(epochItems[0].BlockTimestamp, 0),
					EpochID:          id,
					NodeAddress:      node.Address,
					MinTokensToStake: node.MinTokensToStake,
				})
			}

			// Save the node min tokens to snapshots.
			if err := s.databaseClient.SaveNodeMinTokensToStakeSnapshots(ctx, snapshots); err != nil {
				return fmt.Errorf("save node min tokens to stake snapshots: %w", err)
			}

			cursor = lo.ToPtr(nodes[len(nodes)-1].Address.String())
		}
	}

	return nil
}

func New(databaseClient database.Client, redisClient *redis.Client, stakingContract *l2.Staking) service.Server {
	return &server{
		cronJob:         cronjob.New(redisClient, Name, Timeout),
		databaseClient:  databaseClient,
		redisClient:     redisClient,
		stakingContract: stakingContract,
	}
}
