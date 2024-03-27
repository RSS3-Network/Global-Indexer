package score

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cronjob"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/enforcer"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var _ service.Server = (*server)(nil)

var Name = "score"

type server struct {
	cronJob        *cronjob.CronJob
	simpleEnforcer *enforcer.SimpleEnforcer
}

func (s *server) Spec() string {
	return "0 */10 * * * *"
}

func (s *server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		if err := s.simpleEnforcer.MaintainScore(ctx); err != nil {
			zap.L().Error("sort nodes error", zap.Error(err))
			return
		}
	})
	if err != nil {
		return fmt.Errorf("add detector cron job: %w", err)
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

	simpleEnforcer, err := enforcer.NewSimpleEnforcer(databaseClient, nil, cache.New(redis), stakingContract)
	if err != nil {
		return nil, fmt.Errorf("new enforcer: %w", err)
	}

	instance := server{
		simpleEnforcer: simpleEnforcer,
		cronJob:        cronjob.New(redis, Name, 10*time.Second),
	}

	return &instance, nil
}
