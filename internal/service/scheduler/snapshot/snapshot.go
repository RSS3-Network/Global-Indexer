package snapshot

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
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var _ service.Server = (*server)(nil)

const (
	Name = "snapshot"

	DefaultCronJobTimeout = 10 * time.Second
)

type server struct {
	cronJob        *cronjob.CronJob
	databaseClient database.Client
}

func (s *server) Spec() string {
	return "* * * * * *"
}

func (s *server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		if err := s.snapshot(ctx); err != nil {
			zap.L().Error("snapshot", zap.Error(err))
		}
	})
	if err != nil {
		return fmt.Errorf("register cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopChan

	return nil
}

func (s *server) snapshot(ctx context.Context) error {
	return s.databaseClient.WithTransaction(ctx, func(ctx context.Context, client database.Client) error {
		year, month, day := time.Now().UTC().Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		nodeSnapshot := schema.NodeSnapshot{
			Date: date,
		}

		if err := client.SaveNodeSnapshot(ctx, &nodeSnapshot); err != nil {
			return fmt.Errorf("save node snapshot: %w", err)
		}

		stakeSnapshot := schema.StakeSnapshot{
			Date: date,
		}

		if err := client.SaveStakeSnapshot(ctx, &stakeSnapshot); err != nil {
			return fmt.Errorf("save stake snapshot: %w", err)
		}

		return nil
	})
}

func New(databaseClient database.Client, redis *redis.Client) (service.Server, error) {
	instance := server{
		databaseClient: databaseClient,
		cronJob:        cronjob.New(redis, Name, DefaultCronJobTimeout),
	}

	return &instance, nil
}
