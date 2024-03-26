package provider

import (
	"fmt"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/redis/go-redis/v9"
)

func ProvideRedisClient(config *config.File) (*redis.Client, error) {
	options, err := redis.ParseURL(config.Redis.URI)
	if err != nil {
		return nil, fmt.Errorf("parse redis uri: %w", err)
	}

	return redis.NewClient(options), nil
}
