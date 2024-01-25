package cronjob

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type CronJob struct {
	crontab     *cron.Cron
	redisClient *redis.Client
	lock        sync.RWMutex
}

var KeyPrefix = "cronJob:%s"

func (c *CronJob) AddFunc(ctx context.Context, key, spec string, cmd func()) error {
	_, err := c.crontab.AddFunc(spec, func() {
		acquireLock, err := c.AcquireLock(ctx, key, nil)
		if err != nil {
			zap.L().Error("acquire lock error", zap.String("key", key), zap.Error(err))
			return
		}

		if !acquireLock {
			zap.L().Info("lock is not acquired", zap.String("key", key))

			return
		}

		defer func() {
			if _, err = c.ReleaseLock(ctx, key); err != nil {
				zap.L().Error("release lock error", zap.String("key", key), zap.Error(err))
			}
		}()

		cmd()
	})

	return err
}

func (c *CronJob) Start() {
	c.crontab.Start()
}

func (c *CronJob) Stop() {
	c.crontab.Stop()
}

func (c *CronJob) AcquireLock(ctx context.Context, key string, expiration *time.Duration) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Try to acquire the lock
	result, err := c.redisClient.SetNX(ctx, key, "locked", lo.Ternary(expiration == nil, 0, lo.FromPtr(expiration))).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

func (c *CronJob) ReleaseLock(ctx context.Context, key string) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Release the lock
	if err := c.redisClient.Del(ctx, key).Err(); err != nil {
		return false, err
	}

	return true, nil
}

func New(redisClient *redis.Client) *CronJob {
	return &CronJob{
		crontab:     cron.New(cron.WithSeconds()),
		redisClient: redisClient,
	}
}
