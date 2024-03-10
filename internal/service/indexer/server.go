package indexer

import (
	"context"
	"errors"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/indexer/l1"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/indexer/l2"
	"github.com/sourcegraph/conc/pool"
)

type Server struct {
	config         config.RSS3Chain
	databaseClient database.Client
	cacheClient    cache.Client
}

func (s *Server) Run(ctx context.Context) error {
	errorPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	// Run L1 indexer.
	errorPool.Go(func(ctx context.Context) error {
		l1Config := l1.Config{
			Endpoint:     s.config.EndpointL1,
			BlockThreads: s.config.BlockThreadsL1,
		}

		serverL1, err := l1.NewServer(ctx, s.databaseClient, l1Config)
		if err != nil {
			return err
		}

		return serverL1.Run(ctx)
	})

	// Run L2 indexer.
	errorPool.Go(func(ctx context.Context) error {
		l2Config := l2.Config{
			Endpoint: s.config.EndpointL2,
		}

		serverL2, err := l2.NewServer(ctx, s.databaseClient, s.cacheClient, l2Config)
		if err != nil {
			return err
		}

		return serverL2.Run(ctx)
	})

	if err := errorPool.Wait(); err != nil {
		if !errors.Is(ctx.Err(), context.Canceled) {
			return err
		}
	}

	return nil
}

func New(databaseClient database.Client, cacheClient cache.Client, config config.RSS3Chain) (*Server, error) {
	instance := Server{
		config:         config,
		databaseClient: databaseClient,
		cacheClient:    cacheClient,
	}

	return &instance, nil
}
