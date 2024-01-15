package cache

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/redis/rueidis"
)

var (
	globalLocker      sync.RWMutex
	globalRedisClient rueidis.Client
)

func Global() rueidis.Client {
	globalLocker.RLock()

	defer globalLocker.RUnlock()

	return globalRedisClient
}

func ReplaceGlobal(db rueidis.Client) {
	globalLocker.Lock()

	defer globalLocker.Unlock()

	globalRedisClient = db
}

func Dial(config *Config) (rueidis.Client, error) {
	clientOption := rueidis.ClientOption{
		InitAddress:  config.Endpoints,
		Username:     config.Username,
		Password:     config.Password,
		DisableCache: true,
	}

	return rueidis.NewClient(clientOption)
}

func Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	data, err := globalRedisClient.Do(ctx, globalRedisClient.B().Get().Key(key).Build()).AsBytes()

	if rueidis.IsRedisNil(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if err = json.Unmarshal(data, dest); err != nil {
		return false, err
	}

	return true, nil
}

func Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return globalRedisClient.Do(ctx, globalRedisClient.B().Set().Key(key).Value(rueidis.BinaryString(data)).Build()).Error()
}
