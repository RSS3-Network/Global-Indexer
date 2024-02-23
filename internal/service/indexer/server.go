package indexer

import (
	"context"

	apisixHTTPAPI "github.com/naturalselectionlabs/rss3-global-indexer/internal/apisix/httpapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/indexer/l1"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/indexer/l2"
	"github.com/sourcegraph/conc/pool"
)

type Server struct {
	config              config.RSS3Chain
	databaseClient      database.Client
	apisixHTTPAPIClient *apisixHTTPAPI.Client // For billing - account resume only
	ruPerToken          int64                 // For billing - deposit only
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

		serverL2, err := l2.NewServer(ctx, s.databaseClient, s.apisixHTTPAPIClient, s.ruPerToken, l2Config)
		if err != nil {
			return err
		}

		return serverL2.Run(ctx)
	})

	errorChan := make(chan error)
	go func() { errorChan <- errorPool.Wait() }()

	select {
	case err := <-errorChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func New(databaseClient database.Client, apisixHTTPAPIClient *apisixHTTPAPI.Client, ruPerToken int64, config config.RSS3Chain) (*Server, error) {
	instance := Server{
		config:              config,
		databaseClient:      databaseClient,
		apisixHTTPAPIClient: apisixHTTPAPIClient,
		ruPerToken:          ruPerToken,
	}

	return &instance, nil
}
