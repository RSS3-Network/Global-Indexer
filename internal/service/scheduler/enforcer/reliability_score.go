package enforcer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/cronjob"
	"github.com/rss3-network/global-indexer/internal/database"
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
	return "0 */10 * * * *"
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

func New(databaseClient database.Client, redis *redis.Client, ethereumClient *ethclient.Client) (service.Server, error) {
	chainID, err := ethereumClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	contractAddresses := l2.ContractMap[chainID.Uint64()]
	if contractAddresses == nil {
		return nil, fmt.Errorf("contract address not found for chain id: %s", chainID.String())
	}

	stakingContract, err := l2.NewStaking(contractAddresses.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	instance := server{
		cronJob: cronjob.New(redis, Name, 10*time.Second),

		simpleEnforcer: enforcer.NewSimpleEnforcer(databaseClient, cache.New(redis), stakingContract, nil),
	}

	return &instance, nil
}
