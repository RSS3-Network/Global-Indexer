package taxer

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/cronjob"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"go.uber.org/zap"
)

var _ service.Server = (*Server)(nil)

var (
	Name    = "taxer"
	Timeout = 3 * time.Minute
)

type Server struct {
	cronJob         *cronjob.CronJob
	databaseClient  database.Client
	chainID         *big.Int
	stakingContract *stakingv2.Staking
	settlerConfig   *config.Settler
	txManager       txmgr.TxManager
}

func (s *Server) Name() string {
	return Name
}

func (s *Server) Spec() string {
	return "*/10 * * * * *" // every 10 seconds
}

func (s *Server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		if err := s.checkAndSubmitAverageTaxRate(ctx); err != nil {
			zap.L().Error("submit average tax rate error", zap.Error(err))

			return
		}
	})

	if err != nil {
		return fmt.Errorf("add cron job error: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	return nil
}

func New(databaseClient database.Client, redisClient *redis.Client, ethereumClient *ethclient.Client, config *config.File, txManager *txmgr.SimpleTxManager) (*Server, error) {
	chainID, err := ethereumClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	contractAddresses := l2.ContractMap[chainID.Uint64()]
	if contractAddresses == nil {
		return nil, fmt.Errorf("contract address not found for chain id: %d", chainID.Uint64())
	}

	stakingContract, err := stakingv2.NewStaking(contractAddresses.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	server := &Server{
		cronJob:         cronjob.New(redisClient, Name, Timeout),
		databaseClient:  databaseClient,
		chainID:         chainID,
		stakingContract: stakingContract,
		settlerConfig:   config.Settler,
		txManager:       txManager,
	}

	return server, nil
}
