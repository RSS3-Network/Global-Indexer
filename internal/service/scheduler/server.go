package scheduler

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/detector"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/sort"
	"github.com/redis/go-redis/v9"
)

func New(server string, databaseClient database.Client, redis *redis.Client, ethereumClient *ethclient.Client) (service.Server, error) {
	switch server {
	case "detector":
		return detector.New(databaseClient, redis)
	case "sort":
		return sort.New(databaseClient, redis, ethereumClient)
	}

	return nil, fmt.Errorf("unknown scheduler server: %s", server)
}
