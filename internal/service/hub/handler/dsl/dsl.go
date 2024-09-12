package dsl

import (
	"context"
	"math/big"

	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/nameresolver"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/distributor"
)

type DSL struct {
	distributor    *distributor.Distributor
	databaseClient database.Client
	cacheClient    cache.Client
	nameService    *nameresolver.NameResolver
}

func NewDSL(ctx context.Context, databaseClient database.Client, cacheClient cache.Client, nameService *nameresolver.NameResolver, stakingContract *l2.StakingV2MulticallClient, networkParamsContract *l2.NetworkParams, httpClient httputil.Client, txManager *txmgr.SimpleTxManager, settlerConfig *config.Settler, chainID *big.Int) (*DSL, error) {
	distributorService, err := distributor.NewDistributor(ctx, databaseClient, cacheClient, httpClient, stakingContract, networkParamsContract, txManager, settlerConfig, chainID)
	if err != nil {
		return nil, err
	}

	return &DSL{
		distributor:    distributorService,
		databaseClient: databaseClient,
		cacheClient:    cacheClient,
		nameService:    nameService,
	}, nil
}
