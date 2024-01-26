package cronjob

import (
	"context"
	"fmt"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type CronJob struct {
	crontab *cron.Cron
	mutex   *redsync.Mutex
}

var KeyPrefix = "scheduler:%s"

func (c *CronJob) AddFunc(_ context.Context, spec string, cmd func()) error {
	_, err := c.crontab.AddFunc(spec, func() {
		if err := c.mutex.Lock(); err != nil {
			zap.L().Error("lock error", zap.String("key", c.mutex.Name()), zap.Error(err))

			return
		}

		defer func() {
			if _, err := c.mutex.Unlock(); err != nil {
				zap.L().Error("release lock error", zap.String("key", c.mutex.Name()), zap.Error(err))
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

	if _, err := c.mutex.Unlock(); err != nil {
		zap.L().Error("release lock error", zap.String("key", c.mutex.Name()), zap.Error(err))
	}
}

func New(redisClient *redis.Client, name string) *CronJob {
	pool := goredis.NewPool(redisClient)
	rs := redsync.New(pool)

	return &CronJob{
		crontab: cron.New(cron.WithSeconds()),
		mutex:   rs.NewMutex(fmt.Sprintf(KeyPrefix, name)),
	}
}
