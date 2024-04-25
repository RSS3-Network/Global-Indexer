package cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	PSubscribe(ctx context.Context, pattern string) *redis.PubSub
}

var _ Client = (*client)(nil)

type client struct {
	redisClient *redis.Client
}

func (c *client) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

func (c *client) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.redisClient.Set(ctx, key, data, 0).Err()
}

func (c *client) PSubscribe(ctx context.Context, pattern string) *redis.PubSub {
	return c.redisClient.PSubscribe(ctx, pattern)
}

func New(redisClient *redis.Client) Client {
	return &client{
		redisClient: redisClient,
	}
}
