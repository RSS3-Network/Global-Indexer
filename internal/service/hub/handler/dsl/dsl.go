package dsl

import (
	"context"

	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/distributor"
	"github.com/rss3-network/global-indexer/internal/nameresolver"
)

type DSL struct {
	Distributor    *distributor.Distributor
	databaseClient database.Client
	cacheClient    cache.Client
	nameService    *nameresolver.NameResolver
}

func NewDSL(ctx context.Context, databaseClient database.Client, cacheClient cache.Client, nameService *nameresolver.NameResolver) *DSL {
	return &DSL{
		Distributor:    distributor.NewDistributor(ctx, databaseClient, cacheClient),
		databaseClient: databaseClient,
		cacheClient:    cacheClient,
		nameService:    nameService,
	}
}
