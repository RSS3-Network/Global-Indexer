package federatedhandles

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/cronjob"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"go.uber.org/zap"
)

var _ service.Server = (*server)(nil)

var Name = "federated_handles"

type server struct {
	cronJob        *cronjob.CronJob
	databaseClient database.Client
	cacheClient    cache.Client
	httpClient     httputil.Client
}

func (s *server) Name() string {
	return Name
}

func (s *server) Spec() string {
	return "@every 15m"
}

func (s *server) Run(ctx context.Context) error {
	// initial execution of maintaining federated handles
	if err := s.maintainFederatedHandles(ctx); err != nil {
		zap.L().Error("initial execution of maintaining federated handles failed", zap.Error(err))
	}

	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		if err := s.maintainFederatedHandles(ctx); err != nil {
			zap.L().Error("maintain federated handles error", zap.Error(err))
			return
		}
	})

	if err != nil {
		return fmt.Errorf("add maintain federated handles cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopChan

	return nil
}

func New(redisClient *redis.Client, databaseClient database.Client, httpClient httputil.Client) service.Server {
	return &server{
		cronJob:        cronjob.New(redisClient, Name, 10*time.Second),
		databaseClient: databaseClient,
		cacheClient:    cache.New(redisClient),
		httpClient:     httpClient,
	}
}
