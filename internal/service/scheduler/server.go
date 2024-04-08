package scheduler

import (
	"fmt"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	averagetax "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/average_tax"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/detector"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/integrator"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/snapshot"
	"github.com/redis/go-redis/v9"
)

func New(server string, databaseClient database.Client, redis *redis.Client, config *config.File) (service.Server, error) {
	switch server {
	case detector.Name:
		return detector.New(databaseClient, redis)
	case integrator.Name:
		return integrator.New(databaseClient, redis, config)
	case snapshot.Name:
		return snapshot.New(databaseClient, redis, config)
	case averagetax.Name:
		return averagetax.New(databaseClient, redis, config)
	default:
		return nil, fmt.Errorf("unknown scheduler server: %s", server)
	}
}
