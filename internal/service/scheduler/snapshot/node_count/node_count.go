package nodecount

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/cronjob"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/schema"
	"go.uber.org/zap"
)

var (
	Name    = "node_count"
	Timeout = 10 * time.Second
)

var _ service.Server = (*server)(nil)

type server struct {
	cronJob        *cronjob.CronJob
	databaseClient database.Client
	redisClient    *redis.Client
}

func (s *server) Name() string {
	return Name
}

func (s *server) Spec() string {
	return "0 0 0 * * *"
}

func (s *server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		year, month, day := time.Now().UTC().Date()
		date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		nodeSnapshot := schema.NodeSnapshot{
			Date: date,
		}

		if err := s.databaseClient.SaveNodeCountSnapshot(ctx, &nodeSnapshot); err != nil {
			zap.L().Error("save Node count snapshot error", zap.Error(err))

			return
		}
	})
	if err != nil {
		return fmt.Errorf("add node count cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	return nil
}

func New(databaseClient database.Client, redis *redis.Client) service.Server {
	return &server{
		cronJob:        cronjob.New(redis, Name, Timeout),
		databaseClient: databaseClient,
		redisClient:    redis,
	}
}
