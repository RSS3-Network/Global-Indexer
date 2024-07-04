package indexer

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/client/ethereum"
	"github.com/rss3-network/global-indexer/internal/config/flag"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/internal/service/indexer/internal"
	"github.com/rss3-network/global-indexer/internal/service/indexer/internal/handler/l1"
	"github.com/rss3-network/global-indexer/internal/service/indexer/internal/handler/l2"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/viper"
)

const Name = "indexer"

type Server struct {
	databaseClient           database.Client
	cacheClient              cache.Client
	ethereumMultiChainClient *ethereum.MultiChainClient
}

func (s *Server) Name() string {
	return Name
}

func (s *Server) Run(ctx context.Context) error {
	errorPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	// Run L1 indexers.
	{
		// Run L1 finalized indexer.
		errorPool.Go(func(ctx context.Context) error {
			indexer, err := s.newL1Indexer(true)
			if err != nil {
				return fmt.Errorf("new l1 indexer: %w", err)
			}

			return indexer.Run(ctx)
		})

		// Run L1 unfinalized indexer.
		errorPool.Go(func(ctx context.Context) error {
			indexer, err := s.newL1Indexer(false)
			if err != nil {
				return fmt.Errorf("new l1 indexer: %w", err)
			}

			return indexer.Run(ctx)
		})
	}

	// Run L2 indexers.
	{
		// Run L2 finalized indexer.
		errorPool.Go(func(ctx context.Context) error {
			indexer, err := s.newL2Indexer(true)
			if err != nil {
				return fmt.Errorf("new l2 indexer: %w", err)
			}

			return indexer.Run(ctx)
		})

		// Run L2 unfinalized indexer.
		errorPool.Go(func(ctx context.Context) error {
			indexer, err := s.newL2Indexer(false)
			if err != nil {
				return fmt.Errorf("new l2 indexer: %w", err)
			}

			return indexer.Run(ctx)
		})
	}

	if err := errorPool.Wait(); err != nil {
		if !errors.Is(ctx.Err(), context.Canceled) {
			return err
		}
	}

	return nil
}

func (s *Server) newL1Indexer(finalized bool) (internal.Indexer, error) {
	chainID := viper.GetUint64(flag.KeyChainIDL1)

	ethereumClient, err := s.ethereumMultiChainClient.Get(chainID)
	if err != nil {
		return nil, fmt.Errorf("load ethereum client: %w", err)
	}

	handler, err := l1.NewHandler(chainID, ethereumClient, finalized)
	if err != nil {
		return nil, fmt.Errorf("new l1 handler: %w", err)
	}

	indexer, err := internal.NewIndexer(chainID, ethereumClient, s.databaseClient, handler, finalized)
	if err != nil {
		return nil, fmt.Errorf("new l1 indexer: %w", err)
	}

	return indexer, nil
}

func (s *Server) newL2Indexer(finalized bool) (internal.Indexer, error) {
	chainID := viper.GetUint64(flag.KeyChainIDL2)

	ethereumClient, err := s.ethereumMultiChainClient.Get(chainID)
	if err != nil {
		return nil, fmt.Errorf("load ethereum client: %w", err)
	}

	handler, err := l2.NewHandler(chainID, ethereumClient, s.cacheClient, finalized)
	if err != nil {
		return nil, fmt.Errorf("new l2 handler: %w", err)
	}

	indexer, err := internal.NewIndexer(chainID, ethereumClient, s.databaseClient, handler, finalized)
	if err != nil {
		return nil, fmt.Errorf("new l2 indexer: %w", err)
	}

	return indexer, nil
}

func NewServer(databaseClient database.Client, redisClient *redis.Client, ethereumMultiChainClient *ethereum.MultiChainClient) (service.Server, error) {
	instance := Server{
		databaseClient:           databaseClient,
		cacheClient:              cache.New(redisClient),
		ethereumMultiChainClient: ethereumMultiChainClient,
	}

	return &instance, nil
}
