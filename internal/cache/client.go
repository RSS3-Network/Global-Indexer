package cache

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	globalLocker      sync.RWMutex
	globalRedisClient *redis.Client
)

func Global() *redis.Client {
	globalLocker.RLock()

	defer globalLocker.RUnlock()

	return globalRedisClient
}

func ReplaceGlobal(db *redis.Client) {
	globalLocker.Lock()

	defer globalLocker.Unlock()

	globalRedisClient = db
}

func New(config *Config) (*redis.Client, error) {
	clientOption := &redis.Options{
		Addr:     config.Endpoints,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	}

	return redis.NewClient(clientOption), nil
}

func Get(ctx context.Context, key string, dest interface{}) error {
	data, err := globalRedisClient.Get(ctx, key).Bytes()

	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, dest); err != nil {
		return err
	}

	return nil
}

func Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return globalRedisClient.Set(ctx, key, data, 0).Err()
}
