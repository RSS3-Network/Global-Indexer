package epoch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

type Server struct {
	chainID        *big.Int
	checkpoint     uint64
	timer          *time.Timer
	mutex          *redsync.Mutex
	currentEpoch   uint64
	gasLimit       uint64
	fromAddress    common.Address
	rpcClient      *rpc.Client
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

	// Listen timer
	errorPool.Go(func(ctx context.Context) error {
		if err := s.listenTimer(ctx); err != nil {
			zap.L().Error("listen timer", zap.Error(err))

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
		checkpoint, blockNumberLatest, err := s.loadCheckpoint(ctx)
		if err != nil {
			zap.L().Error("get checkpoint and latest block number", zap.Error(err), zap.Any("chain_id", s.chainID),
				zap.Any("checkpoint", checkpoint), zap.Uint64("block_number_latest", blockNumberLatest))

			return err
		}

		// If indexer doesn't work or catch up the latest block, wait for 5 seconds.
		if int(blockNumberLatest-checkpoint) > 5 {
			zap.L().Info("indexer doesn't work or catch up the latest block", zap.Uint64("checkpoint", checkpoint),
				zap.Uint64("last checkpoint", s.checkpoint), zap.Uint64("block_number_latest", blockNumberLatest))

			time.Sleep(5 * time.Second)

			continue
		}

		s.checkpoint = checkpoint

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

		if now.Sub(lastEpochEventTime) >= 18*time.Hour && now.Sub(lastEpochTriggerTime) >= 18*time.Hour {
			// Trigger new epoch
			if err := s.trigger(ctx, s.currentEpoch+1); err != nil {
				zap.L().Error("trigger new epoch", zap.Error(err))

				return err
			}
		} else if now.Sub(lastEpochEventTime) >= 18*time.Hour && now.Sub(lastEpochTriggerTime) < 18*time.Hour {
			// Wait for epoch event indexer
			time.Sleep(5 * time.Second)
		} else if now.Sub(lastEpochEventTime) < 18*time.Hour {
			// Set timer
			s.timer = time.NewTimer(18*time.Hour - now.Sub(lastEpochEventTime))
			time.Sleep(time.Minute)
		}
	}
}

func (s *Server) listenTimer(ctx context.Context) error {
	for {
		if s.timer == nil {
			continue
		}

		select {
		case <-s.timer.C:
			// Timer expired, trigger new epoch
			err := s.trigger(ctx, s.currentEpoch+1)
			if err != nil {
				zap.L().Error("trigger new epoch", zap.Error(err))

				return err
			}

			s.timer = nil
		case <-ctx.Done():
			// Context cancelled, stop the goroutine
			return nil
		}
	}
}

func (s *Server) loadCheckpoint(ctx context.Context) (uint64, uint64, error) {
	// Load checkpoint from database.
	checkpoint, err := s.databaseClient.FindCheckpoint(ctx, s.chainID.Uint64())
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return 0, 0, nil
		}

		return 0, 0, fmt.Errorf("get checkpoint from database: %w", err)
	}

	// Load latest block number from RPC.
	blockNumberLatest, err := s.ethereumClient.BlockNumber(ctx)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return 0, 0, nil
		}

		return 0, 0, fmt.Errorf("get latest block number from rpc: %w", err)
	}

	return checkpoint.BlockNumber, blockNumberLatest, nil
}

func (s *Server) ping() (string, error) {
	var v string

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	if err := s.rpcClient.CallContext(ctx, &v, "health_status"); err != nil {
		return "", err
	}

	return v, nil
}

func New(ctx context.Context, databaseClient database.Client, redisClient *redis.Client, config config.File) (*Server, error) {
	ethereumClient, err := ethclient.Dial(config.RSS3Chain.EndpointL2)
	if err != nil {
		return nil, fmt.Errorf("dial ethereum client: %w", err)
	}

	chainID, err := ethereumClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	rpcClient, err := rpc.DialOptions(ctx, config.Epoch.SignerEndpoint, rpc.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return nil, fmt.Errorf("dial rpc client: %w", err)
	}

	redisPool := goredis.NewPool(redisClient)
	rs := redsync.New(redisPool)

	server := &Server{
		chainID:        chainID,
		mutex:          rs.NewMutex("epoch", redsync.WithExpiry(5*time.Minute)),
		gasLimit:       config.Epoch.GasLimit,
		fromAddress:    common.HexToAddress(config.Epoch.WalletAddress),
		rpcClient:      rpcClient,
		ethereumClient: ethereumClient,
		databaseClient: databaseClient,
	}

	// Check signer if reachable
	status, err := server.ping()
	if err != nil {
		return nil, err
	}

	if status != "ok" {
		return nil, fmt.Errorf("signer service unreachable")
	}

	return server, nil
}
