package epoch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	gicrypto "github.com/naturalselectionlabs/rss3-global-indexer/common/crypto"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/txmgr"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

const (
	// epochInterval is the interval between epochs.
	// The epoch interval is set to be 18 hours.
	epochInterval = 18 * time.Hour
)

type Server struct {
	txManager      txmgr.TxManager
	checkpoint     uint64
	chainID        *big.Int
	mutex          *redsync.Mutex
	currentEpoch   uint64
	gasLimit       uint64
	ethereumClient *ethclient.Client
	databaseClient database.Client
}

func (s *Server) Run(ctx context.Context) error {
	errorPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	// Listen epoch event
	errorPool.Go(func(ctx context.Context) error {
		if err := s.listenEpochEvent(ctx); err != nil {
			zap.L().Error("listen epoch event", zap.Error(err))

			return err
		}

		return nil
	})

	errorChan := make(chan error)
	go func() { errorChan <- errorPool.Wait() }()

	select {
	case err := <-errorChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Server) listenEpochEvent(ctx context.Context) error {
	for {
		// Load checkpoint and latest block number.
		indexedBlock, latestBlock, err := s.loadCheckpoint(ctx)
		if err != nil {
			zap.L().Error("get checkpoint and latest block number", zap.Error(err), zap.Any("chain_id", s.chainID),
				zap.Any("checkpoint", indexedBlock), zap.Uint64("block_number_latest", latestBlock))

			return err
		}

		// If indexer doesn't work or catch up the latest block, wait for 5 seconds.
		if int(latestBlock-indexedBlock) > 5 {
			zap.L().Error("indexer encountered errors or is still catching up with the latest block", zap.Uint64("checkpoint", indexedBlock),
				zap.Uint64("last checkpoint", s.checkpoint), zap.Uint64("block_number_latest", latestBlock))

			<-time.NewTimer(5 * time.Second).C

			continue
		}

		s.checkpoint = indexedBlock

		// Find the latest epoch event from database.
		epochEvent, err := s.databaseClient.FindEpochs(ctx, 1, nil)
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("get latest epoch event from database", zap.Error(err))

			return err
		}

		// Find the latest epoch trigger from database.
		epochTrigger, err := s.databaseClient.FindLatestEpochTrigger(ctx)
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			zap.L().Error("get latest epoch trigger from database", zap.Error(err))

			return err
		}

		var lastEpochEventTime, lastEpochTriggerTime time.Time

		if len(epochEvent) > 0 {
			lastEpochEventTime = time.Unix(epochEvent[0].BlockTimestamp, 0)
			s.currentEpoch = epochEvent[0].ID
		}

		if epochTrigger != nil {
			lastEpochTriggerTime = epochTrigger.CreatedAt
		}

		now := time.Now()

		// The time elapsed since the last epoch event was included on the VSL
		timeSinceLastEpoch := now.Sub(lastEpochEventTime)
		// The time elapsed since the last epoch trigger was sent
		timeSinceLastTrigger := now.Sub(lastEpochTriggerTime)

		// Special case for genesis epoch
		// No longer applicable for the subsequent epochs
		if genesisEpoch, exist := l2.GenesisEpochMap[s.chainID.Uint64()]; exist {
			genesisEpochTime := time.Unix(genesisEpoch, 0)

			// Wait for genesis epoch
			if now.Sub(genesisEpochTime) < epochInterval-1*time.Hour {
				zap.L().Info("wait for genesis epoch start", zap.Time("genesis_epoch_time", genesisEpochTime),
					zap.Time("estimated_epoch_start_time", now.Add(epochInterval-1*time.Hour-now.Sub(genesisEpochTime))))

				<-time.NewTimer(epochInterval - 1*time.Hour - now.Sub(genesisEpochTime)).C

				zap.L().Info("genesis epoch start", zap.Time("start_time", time.Now()))
			}
		}

		// Check if epochInterval has passed since the last epoch event
		if timeSinceLastEpoch >= epochInterval {
			// Check if the epochInterval has passed since the last epoch trigger
			if timeSinceLastTrigger >= epochInterval {
				// Trigger new epoch
				if err := s.trigger(ctx, s.currentEpoch+1); err != nil {
					zap.L().Error("trigger new epoch", zap.Error(err))

					return err
				}
			} else if timeSinceLastTrigger < epochInterval {
				// Check if epochInterval has NOT passed since the last epoch trigger
				// If so, delay the trigger by 5 seconds
				zap.L().Info("wait for epoch event indexer", zap.Time("last_epoch_event_time", lastEpochEventTime),
					zap.Time("last_epoch_trigger_time", lastEpochTriggerTime))

				<-time.NewTimer(5 * time.Second).C
			}
		} else if timeSinceLastEpoch < epochInterval {
			// If epochInterval has NOT passed since the last epoch event
			// Wait for the remaining time until the next epoch event

			remainingTime := epochInterval - now.Sub(lastEpochEventTime)
			<-time.NewTimer(remainingTime).C

			if err := s.trigger(ctx, s.currentEpoch+1); err != nil {
				zap.L().Error("trigger new epoch", zap.Error(err))
				return err
			}
		}
	}
}

func (s *Server) loadCheckpoint(ctx context.Context) (uint64, uint64, error) {
	// Load checkpoint from database.
	// A checkpoint is basically the last indexed block
	indexedBlock, err := s.databaseClient.FindCheckpoint(ctx, s.chainID.Uint64())
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return 0, 0, nil
		}

		return 0, 0, fmt.Errorf("get checkpoint from database: %w", err)
	}

	// Load latest block number from RPC.
	latestBlock, err := s.ethereumClient.BlockNumber(ctx)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return 0, 0, nil
		}

		return 0, 0, fmt.Errorf("get latest block from rpc: %w", err)
	}

	return indexedBlock.BlockNumber, latestBlock, nil
}

func New(ctx context.Context, databaseClient database.Client, redisClient *redis.Client, config config.File) (*Server, error) {
	redisPool := goredis.NewPool(redisClient)
	rs := redsync.New(redisPool)

	ethereumClient, err := ethclient.Dial(config.RSS3Chain.EndpointL2)
	if err != nil {
		return nil, fmt.Errorf("dial ethereum client: %w", err)
	}

	chainID, err := ethereumClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	signerFactory, from, err := gicrypto.NewSignerFactory(config.Epoch.PrivateKey, config.Epoch.SignerEndpoint, config.Epoch.WalletAddress)

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
		chainID:        chainID,
		mutex:          rs.NewMutex("epoch", redsync.WithExpiry(5*time.Minute)),
		gasLimit:       config.Epoch.GasLimit,
		ethereumClient: ethereumClient,
		databaseClient: databaseClient,
		txManager:      txManager,
	}

	return server, nil
}
