package provider

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/config"
)

func ProvideRedisClient(config *config.File) (*redis.Client, error) {
	options, err := redis.ParseURL(config.Redis.URI)
	if err != nil {
		return nil, fmt.Errorf("parse redis uri: %w", err)
	}

	return redis.NewClient(options), nil
}
