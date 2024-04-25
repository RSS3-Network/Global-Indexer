package dsl

import (
	"context"

	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
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

func NewDSL(ctx context.Context, databaseClient database.Client, cacheClient cache.Client, nameService *nameresolver.NameResolver, stakingContract *l2.Staking, httpClient httputil.Client) (*DSL, error) {
	distributor, err := distributor.NewDistributor(ctx, databaseClient, cacheClient, httpClient, stakingContract)
	if err != nil {
		return nil, err
	}

	return &DSL{
		distributor:    distributor,
		databaseClient: databaseClient,
		cacheClient:    cacheClient,
		nameService:    nameService,
	}, nil
}
