package reliabilityscore

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/cronjob"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/enforcer"
	"go.uber.org/zap"
)

var _ service.Server = (*server)(nil)

var Name = "reliability_score"

type server struct {
	cronJob        *cronjob.CronJob
	simpleEnforcer *enforcer.SimpleEnforcer
}

func (s *server) Name() string {
	return Name
}

func (s *server) Spec() string {
	return "0 */5 * * * *"
}

func (s *server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		if err := s.simpleEnforcer.MaintainReliabilityScore(ctx); err != nil {
			zap.L().Error("maintain reliability_score error", zap.Error(err))
			return
		}
	})

	if err != nil {
		return fmt.Errorf("add maintain reliability score cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopChan

	return nil
}

func New(redis *redis.Client, simpleEnforcer *enforcer.SimpleEnforcer) service.Server {
	return &server{
		cronJob:        cronjob.New(redis, Name, 10*time.Second),
		simpleEnforcer: simpleEnforcer,
	}
}
