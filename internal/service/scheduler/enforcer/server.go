package enforcer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/enforcer"
	epochfresher "github.com/rss3-network/global-indexer/internal/service/scheduler/enforcer/epoch_fresher"
	federatedhandles "github.com/rss3-network/global-indexer/internal/service/scheduler/enforcer/federated_handles"
	nodestatus "github.com/rss3-network/global-indexer/internal/service/scheduler/enforcer/node_status"
	reliabilityscore "github.com/rss3-network/global-indexer/internal/service/scheduler/enforcer/reliability_score"
	"github.com/sourcegraph/conc/pool"
)

var Name = "enforcer"

var _ service.Server = (*server)(nil)

type server struct {
	enforcers []service.Server
}

func (s *server) Name() string {
	return Name
}

func (s *server) Run(ctx context.Context) error {
	errorPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	for _, e := range s.enforcers {
		e := e

		errorPool.Go(func(ctx context.Context) error {
			return e.Run(ctx)
		})
	}

	return errorPool.Wait()
}

func New(databaseClient database.Client, redis *redis.Client, ethereumClient *ethclient.Client, httpClient httputil.Client, config *config.File, txManager *txmgr.SimpleTxManager) (service.Server, error) {
	chainID, err := ethereumClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	contractAddresses := l2.ContractMap[chainID.Uint64()]
	if contractAddresses == nil {
		return nil, fmt.Errorf("contract address not found for chain id: %d", chainID.Uint64())
	}

	networkParamsContract, err := l2.NewNetworkParams(contractAddresses.AddressNetworkParamsProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new network contract: %w", err)
	}

	contractStakingEvents, err := l2.NewEvents(contractAddresses.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking events contract: %w", err)
	}

	stakingContract, err := l2.NewStakingV2MulticallClient(chainID.Uint64(), ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	simpleEnforcer, err := enforcer.NewSimpleEnforcer(context.Background(), databaseClient, cache.New(redis), stakingContract, networkParamsContract, httpClient, txManager, config.Settler, chainID, false)

	if err != nil {
		return nil, fmt.Errorf("new simple enforcer: %w", err)
	}

	settlementContract, err := l2.NewSettlement(contractAddresses.AddressSettlementProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new settlement contract: %w", err)
	}

	checkpoint, err := databaseClient.FindCheckpoint(context.Background(), chainID.Uint64())
	if err != nil {
		return nil, fmt.Errorf("get checkpoint: %w", err)
	}

	return &server{
		enforcers: []service.Server{
			nodestatus.New(redis, simpleEnforcer),
			reliabilityscore.New(redis, simpleEnforcer),
			epochfresher.New(redis, ethereumClient, checkpoint.BlockNumber, simpleEnforcer, contractStakingEvents, settlementContract, contractAddresses.AddressStakingProxy),
			federatedhandles.New(redis, databaseClient, httpClient),
		},
	}, nil
}
