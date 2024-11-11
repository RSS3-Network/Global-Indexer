package snapshot

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/contract/l2"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/apy"
	nodecount "github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/node_count"
	operatorprofit "github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/operator_profit"
	stakercount "github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/staker_count"
	stakercumulativeearnings "github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/staker_cumulative_earnings"
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

	return errorPool.Wait()
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

	stakingContract, err := stakingv2.NewStaking(contractAddresses.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	return &server{
		snapshots: []service.Server{
			nodecount.New(databaseClient, redis),
			stakercount.New(databaseClient, redis),
			stakerprofit.New(databaseClient, redis, stakingContract),
			operatorprofit.New(databaseClient, redis, stakingContract),
			apy.New(databaseClient, redis, stakingContract),
			stakercumulativeearnings.New(databaseClient, redis, stakingContract),
		},
	}, nil
}
