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
	Distributor    *distributor.Distributor
	databaseClient database.Client
	cacheClient    cache.Client
	nameService    *nameresolver.NameResolver
}

func NewDSL(ctx context.Context, databaseClient database.Client, cacheClient cache.Client, nameService *nameresolver.NameResolver, stakingContract *l2.Staking, httpClient httputil.Client) *DSL {
	return &DSL{
		Distributor:    distributor.NewDistributor(ctx, databaseClient, cacheClient, httpClient, stakingContract),
		databaseClient: databaseClient,
		cacheClient:    cacheClient,
		nameService:    nameService,
	}
}
