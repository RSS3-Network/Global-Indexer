package snapshot

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	nodecount "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/snapshot/node_count"
	nodemintokenstostake "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/snapshot/node_min_tokens_to_stake"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/snapshot/stake"
	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/pool"
)

var Name = "snapshot"

var _ service.Server = (*server)(nil)

type server struct {
	snapshots []service.Server
}

func (s *server) Run(ctx context.Context) error {
	errorPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	for _, snapshot := range s.snapshots {
		snapshot := snapshot

		errorPool.Go(func(ctx context.Context) error {
			return snapshot.Run(ctx)
		})
	}

	if err := errorPool.Wait(); err != nil {
		return err
	}

	return nil
}

func New(databaseClient database.Client, redis *redis.Client, ethereumClient *ethclient.Client) (service.Server, error) {
	nodeMinTokensToStakeSnapshot, err := nodemintokenstostake.New(databaseClient, redis, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new node min tokens to stake snapshot: %w", err)
	}

	return &server{
		snapshots: []service.Server{
			nodecount.New(databaseClient, redis),
			stake.New(databaseClient, redis),
			nodeMinTokensToStakeSnapshot,
		},
	}, nil
}
