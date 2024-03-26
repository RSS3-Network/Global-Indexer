package indexer

import (
	"context"
	"errors"
	"fmt"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/client/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config/flag"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/indexer/l1"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/indexer/l2"
	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/viper"
)

const Name = "indexer"

type Server struct {
	config                   *config.RSS3Chain
	databaseClient           database.Client
	cacheClient              cache.Client
	ethereumMultiChainClient *ethereum.MultiChainClient
}

func (s *Server) Name() string {
	return Name
}

func (s *Server) Run(ctx context.Context) error {
	errorPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	// Run L1 indexer.
	errorPool.Go(func(ctx context.Context) error {
		ethereumClient, err := s.ethereumMultiChainClient.Get(viper.GetUint64(flag.KeyChainIDL1))
		if err != nil {
			return fmt.Errorf("get ethereum client: %w", err)
		}

		l1Config := l1.Config{
			BlockThreads: s.config.BlockThreadsL1,
		}

		serverL1, err := l1.NewServer(ctx, s.databaseClient, ethereumClient, l1Config)
		if err != nil {
			return err
		}

		return serverL1.Run(ctx)
	})

	// Run L2 indexer.
	errorPool.Go(func(ctx context.Context) error {
		ethereumClient, err := s.ethereumMultiChainClient.Get(viper.GetUint64(flag.KeyChainIDL2))
		if err != nil {
			return fmt.Errorf("get ethereum client: %w", err)
		}

		l2Config := l2.Config{
			BlockThreads: s.config.BlockThreadsL2,
		}

		serverL2, err := l2.NewServer(ctx, s.databaseClient, s.cacheClient, ethereumClient, l2Config)
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

func NewServer(databaseClient database.Client, redisClient *redis.Client, ethereumMultiChainClient *ethereum.MultiChainClient, configFile *config.File) (service.Server, error) {
	instance := Server{
		config:                   configFile.RSS3Chain,
		databaseClient:           databaseClient,
		cacheClient:              cache.New(redisClient),
		ethereumMultiChainClient: ethereumMultiChainClient,
	}

	return &instance, nil
}
