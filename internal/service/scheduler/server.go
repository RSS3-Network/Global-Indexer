package scheduler

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/detector"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/integrator"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/snapshot"
	"github.com/redis/go-redis/v9"
)

func New(server string, databaseClient database.Client, redis *redis.Client, ethereumClient *ethclient.Client) (service.Server, error) {
	switch server {
	case detector.Name:
		return detector.New(databaseClient, redis)
	case integrator.Name:
		return integrator.New(databaseClient, redis, ethereumClient)
	case snapshot.Name:
		return snapshot.New(databaseClient, redis)
	default:
		return nil, fmt.Errorf("unknown scheduler server: %s", server)
	}
}
