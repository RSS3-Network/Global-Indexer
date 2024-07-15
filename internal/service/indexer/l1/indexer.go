package l1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rss3-network/global-indexer/contract/l1"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	checkpoint                     *schema.Checkpoint
	blockNumberLatest              uint64
	blockThreads                   uint64
}

func (s *server) Name() string {
	return "l1"
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

	retryableFunc := func() error {
		for {
			if err := s.index(ctx); err != nil {
				return err
			}
		}
	}

	onRetry := retry.OnRetry(func(n uint, err error) {
		if !errors.Is(ctx.Err(), context.Canceled) {
			zap.L().Error("run indexer", zap.Error(err), zap.Uint("attempts", n))
		}
	})

	return retry.Do(retryableFunc, retry.Context(ctx), retry.DelayType(retry.FixedDelay), retry.Delay(time.Second), retry.Attempts(30), onRetry)
}

func (s *server) refreshLatestBlockNumber(ctx context.Context) (err error) {
	ctx, span := otel.Tracer("").Start(ctx, "refreshLatestBlockNumber")
	defer span.End()

	if s.blockNumberLatest, err = s.ethereumClient.BlockNumber(ctx); err != nil {
		return fmt.Errorf("get latest block number: %w", err)
	}

	span.SetAttributes(attribute.Int64("block.number", int64(s.blockNumberLatest)))

	return nil
}

func (s *server) fetchBlocks(ctx context.Context) ([]*types.Block, error) {
	ctx, span := otel.Tracer("").Start(ctx, "fetchBlocks")
	defer span.End()

	var blockNumbers []int64

	for offset := uint64(1); offset <= s.blockThreads; offset++ {
		blockNumber := s.checkpoint.BlockNumber + offset

		if blockNumber > s.blockNumberLatest {
			continue
		}

		blockNumbers = append(blockNumbers, int64(blockNumber))
	}

	span.SetAttributes(
		attribute.Int64Slice("block.numbers", blockNumbers),
	)

	resultPool := pool.NewWithResults[*types.Block]().
		WithContext(ctx).
		WithCancelOnError().
		WithFirstError()

	for _, blockNumber := range blockNumbers {
		blockNumber := blockNumber

		resultPool.Go(func(ctx context.Context) (*types.Block, error) {
			block, err := s.ethereumClient.BlockByNumber(ctx, new(big.Int).SetInt64(blockNumber))
			if err != nil {
				return nil, fmt.Errorf("get block %d: %w", blockNumber, err)
			}

			return block, nil
		})
	}

	return resultPool.Wait()
}

func (s *server) fetchReceipts(ctx context.Context, blockNumbers []int64) ([]*types.Receipt, error) {
	ctx, span := otel.Tracer("").Start(ctx, "fetchReceipts")
	defer span.End()

	span.SetAttributes(attribute.Int64Slice("block.numbers", blockNumbers))

	resultPool := pool.NewWithResults[[]*types.Receipt]().
		WithContext(ctx).
		WithCancelOnError().
		WithFirstError()

	for _, blockNumber := range blockNumbers {
		blockNumber := blockNumber

		resultPool.Go(func(ctx context.Context) ([]*types.Receipt, error) {
			receipts, err := s.ethereumClient.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber)))
			if err != nil {
				return nil, fmt.Errorf("get receipts for block %d: %w", blockNumber, err)
			}

			return receipts, nil
		})
	}

	results, err := resultPool.Wait()
	if err != nil {
		return nil, fmt.Errorf("wait result pool: %w", err)
	}

	return lo.Flatten(results), nil
}

func (s *server) index(ctx context.Context) (err error) {
	ctx, span := otel.Tracer("").Start(ctx, "index")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("chain.id", s.chainID.Int64()),
		attribute.Int64("block.number.local", int64(s.checkpoint.BlockNumber)),
		attribute.Int64("block.number.latest", int64(s.blockNumberLatest)),
	)

	if err := s.refreshLatestBlockNumber(ctx); err != nil {
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

		timer := time.NewTimer(blockConfirmationTime)

		select {
		case <-ctx.Done():
			break
		case <-timer.C:
			return nil
		}

		return nil
	}

	blocks, err := s.fetchBlocks(ctx)
	if err != nil {
		return fmt.Errorf("fetch blocks: %w", err)
	}

	receipts, err := s.fetchReceipts(ctx, lo.Map(blocks, func(block *types.Block, _ int) int64 { return block.Number().Int64() }))
	if err != nil {
		return fmt.Errorf("fetch receipts: %w", err)
	}

	receiptsMap := lo.GroupBy(receipts, func(receipt *types.Receipt) int64 {
		return receipt.BlockNumber.Int64()
	})

	sort.SliceStable(blocks, func(i, j int) bool {
		return blocks[i].Number().Cmp(blocks[j].Number()) < 0
	})

	for _, block := range blocks {
		if err := s.indexBlock(ctx, block, receiptsMap[block.Number().Int64()]); err != nil {
			return fmt.Errorf("index block %d: %w", block.NumberU64(), err)
		}
	}

	return nil
}

func (s *server) indexBlock(ctx context.Context, block *types.Block, receipts types.Receipts) error {
	// Begin a database transaction for the block.
	databaseTransaction, err := s.databaseClient.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin database transaction: %w", err)
	}

	defer lo.Try(databaseTransaction.Rollback)

	header := block.Header()

	for _, receipt := range receipts {
		// Discard all contract creation transactions.
		if tx := block.Transaction(receipt.TxHash); tx == nil || tx.To() == nil {
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
			case l1.ContractMap[s.chainID.Uint64()].AddressL1StandardBridgeProxy, l1.ContractMap[s.chainID.Uint64()].AddressOptimismPortalProxy:
				if err := s.indexBridgingLog(ctx, header, block.Transaction(log.TxHash), receipt, log, index, databaseTransaction); err != nil {
					return fmt.Errorf("index bridge log %s %d: %w", log.TxHash, log.Index, err)
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

func NewServer(ctx context.Context, databaseClient database.Client, ethereumClient *ethclient.Client, config Config) (service.Server, error) {
	var (
		instance = server{
			databaseClient: databaseClient,
			ethereumClient: ethereumClient,
			blockThreads:   config.BlockThreads,
		}
		err error
	)

	if instance.chainID, err = instance.ethereumClient.ChainID(ctx); err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	contractAddresses := l1.ContractMap[instance.chainID.Uint64()]
	if contractAddresses == nil {
		return nil, fmt.Errorf("chain id %d is not supported", instance.chainID)
	}

	if instance.contractGovernanceToken, err = bindings.NewGovernanceToken(contractAddresses.AddressGovernanceTokenProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractOptimismPortal, err = bindings.NewOptimismPortal(contractAddresses.AddressOptimismPortalProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractL1CrossDomainMessenger, err = bindings.NewL1CrossDomainMessenger(contractAddresses.AddressL1CrossDomainMessengerProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractL1StandardBridge, err = bindings.NewL1StandardBridge(contractAddresses.AddressL1StandardBridgeProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	return &instance, nil
}
