package detector

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cronjob"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var _ service.Server = (*server)(nil)

var Name = "detector"

type server struct {
	cronJob        *cronjob.CronJob
	databaseClient database.Client
}

func (s *server) Name() string {
	return Name
}

func (s *server) Spec() string {
	return "*/5 * * * * *"
}

func (s *server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		if err := s.updateNodeActivity(ctx); err != nil {
			zap.L().Error("detect node activity error", zap.Error(err))
			return
		}
	})
	if err != nil {
		return fmt.Errorf("add detector cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	return nil
}

func (s *server) updateNodeActivity(ctx context.Context) error {
	timeout := time.Now().Add(-5 * time.Minute)

	if err := s.databaseClient.UpdateNodesStatusOffline(ctx, timeout.Unix()); err != nil {
		zap.L().Error("update node activity error", zap.Error(err), zap.String("timeout", timeout.String()))

		return fmt.Errorf("update node activity: %w", err)
	}

	return nil
}

func New(databaseClient database.Client, redis *redis.Client) (service.Server, error) {
	instance := server{
		databaseClient: databaseClient,
		cronJob:        cronjob.New(redis, Name, 10*time.Second),
	}

	return &instance, nil
}
