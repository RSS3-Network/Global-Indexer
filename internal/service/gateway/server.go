package gateway

import (
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/redis/go-redis/v9"
)

func New(databaseClient database.Client, redis *redis.Client) (service.Server, error) {

}
