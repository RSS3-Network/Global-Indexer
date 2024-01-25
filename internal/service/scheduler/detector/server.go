package detector

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cronJob"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"go.uber.org/zap"
)

var _ service.Server = (*server)(nil)

type server struct {
	cronJob        *cronjob.CronJob
	databaseClient database.Client
}

func (s *server) Run(ctx context.Context) error {
	key := fmt.Sprintf(cronjob.KeyPrefix, "detector")

	err := s.cronJob.AddFunc(ctx, key, "*/5 * * * * *", func() {
		if err := s.updateNodeActivity(ctx); err != nil {
			zap.L().Error("detect node activity error", zap.Error(err))
			return
		}
	})
	if err != nil {
		return fmt.Errorf("add detector cron job: %w", err)
	}

	s.cronJob.Start()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	if _, err := s.cronJob.ReleaseLock(ctx, key); err != nil {
		zap.L().Error("release lock error", zap.Error(err))

		return fmt.Errorf("release lock: %w", err)
	}

	return nil
}

func (s *server) updateNodeActivity(ctx context.Context) error {
	timeout := time.Now().Add(-5 * time.Minute)

	if err := s.databaseClient.UpdateNodesStatus(ctx, timeout.Unix()); err != nil {
		zap.L().Error("update node activity error", zap.Error(err), zap.String("timeout", timeout.String()))

		return fmt.Errorf("update node activity: %w", err)
	}

	return nil
}

func New(databaseClient database.Client, redis *redis.Client) (service.Server, error) {
	instance := server{
		databaseClient: databaseClient,
		cronJob:        cronjob.New(redis),
	}

	return &instance, nil
}
