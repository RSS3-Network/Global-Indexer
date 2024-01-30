package cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
}

var _ Client = (*client)(nil)

type client struct {
	cacheClient *redis.Client
}

func (c *client) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.cacheClient.Get(ctx, key).Bytes()

	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, dest); err != nil {
		return err
	}

	return nil
}

func (c *client) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.cacheClient.Set(ctx, key, data, 0).Err()
}

func New(redisClient *redis.Client) Client {
	return &client{
		cacheClient: redisClient,
	}
}
