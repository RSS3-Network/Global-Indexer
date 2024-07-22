package taxer

import (
	"context"
	"fmt"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	gicrypto "github.com/rss3-network/global-indexer/common/crypto"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
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

func New(databaseClient database.Client, redisClient *redis.Client, ethereumClient *ethclient.Client, config *config.File) (*Server, error) {
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

	signerFactory, from, err := gicrypto.NewSignerFactory(config.Settler.PrivateKey, config.Settler.SignerEndpoint, config.Settler.WalletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer")
	}

	defaultTxConfig := txmgr.Config{
		ResubmissionTimeout:       20 * time.Second,
		FeeLimitMultiplier:        5,
		TxSendTimeout:             5 * time.Minute,
		TxNotInMempoolTimeout:     1 * time.Hour,
		NetworkTimeout:            5 * time.Minute,
		ReceiptQueryInterval:      500 * time.Millisecond,
		NumConfirmations:          5,
		SafeAbortNonceTooLowCount: 3,
	}

	txManager, err := txmgr.NewSimpleTxManager(defaultTxConfig, chainID, nil, ethereumClient, from, signerFactory(chainID))
	if err != nil {
		return nil, fmt.Errorf("failed to create tx manager")
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
