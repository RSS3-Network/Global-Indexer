package indexer

import (
	"context"

	"github.com/naturalselectionlabs/global-indexer/internal/config"
	"github.com/naturalselectionlabs/global-indexer/internal/database"
	"github.com/naturalselectionlabs/global-indexer/internal/service/indexer/l1"
	"github.com/naturalselectionlabs/global-indexer/internal/service/indexer/l2"
	"github.com/sourcegraph/conc/pool"
)

type Server struct {
	config         config.RSS3ChainConfig
	databaseClient database.Client
}

func (s *Server) Run(ctx context.Context) error {
	errorGroup := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	// Run L1 indexer.
	errorGroup.Go(func(ctx context.Context) error {
		l1Config := l1.Config{
			Endpoint: s.config.EndpointL1,
		}

		serverL1, err := l1.NewServer(ctx, s.databaseClient, l1Config)
		if err != nil {
			return err
		}

		return serverL1.Run(ctx)
	})

	// Run L2 indexer.
	errorGroup.Go(func(ctx context.Context) error {
		l2Config := l2.Config{
			Endpoint: s.config.EndpointL2,
		}

		serverL2, err := l2.NewServer(ctx, s.databaseClient, l2Config)
		if err != nil {
			return err
		}

		return serverL2.Run(ctx)
	})

	return errorGroup.Wait()
}

func New(databaseClient database.Client, config config.RSS3ChainConfig) (*Server, error) {
	instance := Server{
		config:         config,
		databaseClient: databaseClient,
	}

	return &instance, nil
}
