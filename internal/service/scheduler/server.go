package scheduler

import (
	"fmt"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/detector"
	"github.com/redis/go-redis/v9"
)

func New(server string, databaseClient database.Client, redis *redis.Client) (service.Server, error) {
	switch server {
	case "detector":
		return detector.New(databaseClient, redis)
	}

	return nil, fmt.Errorf("unknown scheduler server: %s", server)
}
