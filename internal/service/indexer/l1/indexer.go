package l1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l1"
	apisixHTTPAPI "github.com/naturalselectionlabs/rss3-global-indexer/internal/apisix/httpapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

var _ service.Server = (*server)(nil)

type server struct {
	databaseClient                 database.Client
	ethereumClient                 *ethclient.Client
	chainID                        *big.Int
	contractGovernanceToken        *bindings.GovernanceToken
	contractOptimismPortal         *bindings.OptimismPortal
	contractL1CrossDomainMessenger *bindings.L1CrossDomainMessenger
	contractL1StandardBridge       *bindings.L1StandardBridge
	contractBilling                *l1.Billing
	checkpoint                     *schema.Checkpoint
	blockNumberLatest              uint64
	blockThreads                   uint64
	apisixHTTPAPIService           *apisixHTTPAPI.Service // For account resume only
}

func (s *server) Run(ctx context.Context) (err error) {
	// Load checkpoint from database.
	if s.checkpoint, err = s.databaseClient.FindCheckpoint(ctx, s.chainID.Uint64()); err != nil {
		return fmt.Errorf("get checkpoint: %w", err)
	}

	// Rollback to the specified block number state.
	if err := s.databaseClient.RollbackBlock(ctx, s.checkpoint.ChainID, s.checkpoint.BlockNumber); err != nil {
		return fmt.Errorf("rollback block: %w", err)
	}

	onRetry := retry.OnRetry(func(n uint, err error) {
		zap.L().Error("run indexer", zap.Error(err), zap.Uint("attempts", n))
	})

	retryIf := retry.RetryIf(func(err error) bool {
		return !errors.Is(err, context.Canceled)
	})

	return retry.Do(func() error { return s.run(ctx) }, retry.DelayType(retry.FixedDelay), retry.Delay(time.Second), retry.Attempts(30), onRetry, retryIf)
}

func (s *server) run(ctx context.Context) (err error) {
	for {
		// Refresh the latest block number.
		if s.blockNumberLatest, err = s.ethereumClient.BlockNumber(ctx); err != nil {
			return fmt.Errorf("get latest block number: %w", err)
		}

		zap.L().Info(
			"refreshed the latest block number",
			zap.Uint64("block.number.local", s.checkpoint.BlockNumber),
			zap.Uint64("block.number.latest", s.blockNumberLatest),
		)

		// Waiting for a new block to be minted.
		if s.checkpoint.BlockNumber >= s.blockNumberLatest {
			blockConfirmationTime := time.Second // TODO Redefine it.

			zap.L().Info(
				"waiting for a new block to be minted",
				zap.Uint64("block.number.local", s.checkpoint.BlockNumber),
				zap.Uint64("block.number.latest", s.blockNumberLatest),
				zap.Duration("block.confirmationTime", blockConfirmationTime),
			)

			time.Sleep(blockConfirmationTime)

			continue
		}

		// Get blocks from RPC.
		blockResultPool := pool.NewWithResults[*types.Block]().
			WithContext(ctx).
			WithCancelOnError().
			WithFirstError()

		for offset := uint64(1); offset <= s.blockThreads; offset++ {
			blockNumber := s.checkpoint.BlockNumber + offset

			if blockNumber > s.blockNumberLatest {
				continue
			}

			blockResultPool.Go(func(ctx context.Context) (*types.Block, error) {
				// Get current block (header and transactions).
				block, err := s.ethereumClient.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
				if err != nil {
					return nil, fmt.Errorf("get block %d: %w", blockNumber, err)
				}

				return block, nil
			})
		}

		blocks, err := blockResultPool.Wait()
		if err != nil {
			return fmt.Errorf("wait block result pool: %w", err)
		}

		// Get receipts from RPC.
		var (
			receiptsMap       = make(map[uint64][]*types.Receipt)
			receiptsMapLocker sync.Mutex
		)

		receiptsPool := pool.New().
			WithContext(ctx).
			WithCancelOnError().
			WithFirstError()

		for _, block := range blocks {
			block := block

			receiptsPool.Go(func(ctx context.Context) error {
				receipts, err := s.ethereumClient.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(block.NumberU64())))
				if err != nil {
					return fmt.Errorf("get receipts for block %d: %w", block.NumberU64(), err)
				}

				receiptsMapLocker.Lock()
				defer receiptsMapLocker.Unlock()
				receiptsMap[block.NumberU64()] = receipts

				return nil
			})
		}

		if err := receiptsPool.Wait(); err != nil {
			return fmt.Errorf("wait receipts pool: %w", err)
		}

		sort.SliceStable(blocks, func(i, j int) bool {
			return blocks[i].Number().Cmp(blocks[j].Number()) < 0
		})

		for _, block := range blocks {
			if err := s.index(ctx, block, receiptsMap[block.NumberU64()]); err != nil {
				return fmt.Errorf("index block %d: %w", block.NumberU64(), err)
			}
		}
	}
}

func (s *server) index(ctx context.Context, block *types.Block, receipts types.Receipts) error {
	// Begin a database transaction for the block.
	databaseTransaction, err := s.databaseClient.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin database transaction: %w", err)
	}

	defer lo.Try(databaseTransaction.Rollback)

	header := block.Header()

	for _, receipt := range receipts {
		// Discard all contract creation transactions.
		if block.Transaction(receipt.TxHash).To() == nil {
			continue
		}

		// Discard all failed transactions.
		if receipt.Status != types.ReceiptStatusSuccessful {
			continue
		}

		for index, log := range receipt.Logs {
			// Discard all removed logs.
			if log.Removed {
				continue
			}

			// Discard all anonymous logs.
			if len(log.Topics) == 0 {
				continue
			}

			switch log.Address {
			case l1.AddressL1StandardBridgeProxy, l1.AddressOptimismPortalProxy:
				if err := s.indexBridgingLog(ctx, header, block.Transaction(log.TxHash), receipt, log, index, databaseTransaction); err != nil {
					return fmt.Errorf("index bridge log %s %d: %w", log.TxHash, log.Index, err)
				}
			case l1.AddressBillingProxy:
				if err := s.indexBillingLog(ctx, header, block.Transaction(log.TxHash), receipt, log, index, databaseTransaction); err != nil {
					return fmt.Errorf("index billing log %s %d: %w", log.TxHash, log.Index, err)
				}
			}
		}
	}

	// Update and save checkpoint to memory and database.
	s.checkpoint.BlockHash = block.Hash()
	s.checkpoint.BlockNumber = block.NumberU64()

	if err := databaseTransaction.SaveCheckpoint(ctx, s.checkpoint); err != nil {
		return fmt.Errorf("save checkpoint: %w", err)
	}

	if databaseTransaction.Commit() != nil {
		return fmt.Errorf("commit database transaction: %w", err)
	}

	return nil
}

func NewServer(ctx context.Context, databaseClient database.Client, apisixHTTPAPIService *apisixHTTPAPI.Service, config Config) (service.Server, error) {
	var (
		instance = server{
			databaseClient:       databaseClient,
			blockThreads:         config.BlockThreads,
			apisixHTTPAPIService: apisixHTTPAPIService,
		}
		err error
	)

	if instance.ethereumClient, err = ethclient.DialContext(ctx, config.Endpoint); err != nil {
		return nil, err
	}

	if instance.chainID, err = instance.ethereumClient.ChainID(ctx); err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	if instance.contractGovernanceToken, err = bindings.NewGovernanceToken(l1.AddressGovernanceTokenProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractOptimismPortal, err = bindings.NewOptimismPortal(l1.AddressOptimismPortalProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractL1CrossDomainMessenger, err = bindings.NewL1CrossDomainMessenger(l1.AddressL1CrossDomainMessengerProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractL1StandardBridge, err = bindings.NewL1StandardBridge(l1.AddressL1StandardBridgeProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractBilling, err = l1.NewBilling(l1.AddressBillingProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	return &instance, nil
}
