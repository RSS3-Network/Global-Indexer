package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	IncrBy(ctx context.Context, key string, value int64) error
	PSubscribe(ctx context.Context, pattern string) *redis.PubSub
	ZAdd(ctx context.Context, key string, members ...redis.Z) error
	ZRem(ctx context.Context, key string, members ...interface{}) error
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	Exists(ctx context.Context, key string) (int64, error)
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

func (c *client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.redisClient.Set(ctx, key, data, expiration).Err()
}

func (c *client) IncrBy(ctx context.Context, key string, value int64) error {
	return c.redisClient.IncrBy(ctx, key, value).Err()
}

func (c *client) PSubscribe(ctx context.Context, pattern string) *redis.PubSub {
	return c.redisClient.PSubscribe(ctx, pattern)
}

func (c *client) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return c.redisClient.ZAdd(ctx, key, members...).Err()
}

func (c *client) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.redisClient.ZRem(ctx, key, members...).Err()
}

func (c *client) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return c.redisClient.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

func (c *client) Exists(ctx context.Context, key string) (int64, error) {
	return c.redisClient.Exists(ctx, key).Result()
}

func New(redisClient *redis.Client) Client {
	return &client{
		redisClient: redisClient,
	}
}
