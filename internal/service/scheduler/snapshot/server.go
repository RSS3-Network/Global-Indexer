package snapshot

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	nodecount "github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/node_count"
	nodemintokenstostake "github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/node_min_tokens_to_stake"
	operationPool "github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/operator_profit"
	stakercount "github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/staker_count"
	stakerprofit "github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/staker_profit"
	"github.com/sourcegraph/conc/pool"
)

var Name = "snapshot"

var _ service.Server = (*server)(nil)

type server struct {
	snapshots []service.Server
}

func (s *server) Name() string {
	return Name
}

func (s *server) Run(ctx context.Context) error {
	errorPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	for _, snapshot := range s.snapshots {
		snapshot := snapshot

		errorPool.Go(func(ctx context.Context) error {
			return snapshot.Run(ctx)
		})
	}

	if err := errorPool.Wait(); err != nil {
		return err
	}

	return nil
}

func New(databaseClient database.Client, redis *redis.Client, ethereumClient *ethclient.Client) (service.Server, error) {
	chainID, err := ethereumClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	contractAddresses := l2.ContractMap[chainID.Uint64()]
	if contractAddresses == nil {
		return nil, fmt.Errorf("contract address not found for chain id: %d", chainID.Uint64())
	}

	stakingContract, err := l2.NewStaking(contractAddresses.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	return &server{
		snapshots: []service.Server{
			nodecount.New(databaseClient, redis),
			stakercount.New(databaseClient, redis),
			nodemintokenstostake.New(databaseClient, redis, stakingContract),
			stakerprofit.New(databaseClient, redis, stakingContract),
			operationPool.New(databaseClient, redis, stakingContract),
		},
	}, nil
}
