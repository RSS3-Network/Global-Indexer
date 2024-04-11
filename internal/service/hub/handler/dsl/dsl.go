package dsl

import (
	"context"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/distributor"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/nameresolver"
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
