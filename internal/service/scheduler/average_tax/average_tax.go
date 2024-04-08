package averagetax

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	gicrypto "github.com/naturalselectionlabs/rss3-global-indexer/common/crypto"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/txmgr"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cronjob"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

var _ service.Server = (*Server)(nil)

var (
	Name    = "average_tax"
	Timeout = 3 * time.Minute
)

type Server struct {
	cronJob         *cronjob.CronJob
	databaseClient  database.Client
	chainID         *big.Int
	stakingContract *l2.Staking
	settlerConfig   *config.Settler
	txManager       txmgr.TxManager
}

func (s *Server) Spec() string {
	return "*/10 * * * * *" // every 10 seconds
}

func (s *Server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		// Query the submission record of the average tax rate
		submissions, err := s.databaseClient.FindAverageTaxSubmissions(ctx, schema.AverageTaxSubmissionQuery{
			Limit: lo.ToPtr(1),
		})
		if err != nil {
			zap.L().Error("find average tax submissions", zap.Error(err))

			return
		}

		// Query the latest of epoch events
		latestEvent, err := s.databaseClient.FindEpochs(ctx, 1, nil)
		if err != nil {
			zap.L().Error("find epochs", zap.Error(err))

			return
		}

		if len(latestEvent) == 0 {
			return
		}

		if len(submissions) > 0 && submissions[0].EpochID == latestEvent[0].ID {
			return
		}

		// Submit a new average tax rate and save record
		if err := s.submitAverageTax(ctx, latestEvent[0].ID); err != nil {
			zap.L().Error("submit average tax", zap.Error(err))

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

func New(databaseClient database.Client, redisClient *redis.Client, config *config.File) (*Server, error) {
	ethereumClient, err := ethclient.Dial(config.RSS3Chain.EndpointL2)
	if err != nil {
		return nil, fmt.Errorf("dial ethereum client: %w", err)
	}

	chainID, err := ethereumClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	contractAddresses := l2.ContractMap[chainID.Uint64()]
	if contractAddresses == nil {
		return nil, fmt.Errorf("contract address not found for chain id: %d", chainID.Uint64())
	}

	stakingContract, err := l2.NewStaking(contractAddresses.AddressStakingProxy, ethereumClient)
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
