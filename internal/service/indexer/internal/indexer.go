package internal

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// Handler uses to process blocks and receipts.
type Handler interface {
	Process(ctx context.Context, block *types.Block, receipts types.Receipts, databaseTransaction database.Client) error
}

// Indexer uses to index blockchain data, it will process blocks and receipts using it Handler.
type Indexer interface {
	Run(ctx context.Context) error
}

type indexer struct {
	ethereumClient    *ethclient.Client
	databaseClient    database.Client
	handler           Handler
	chainID           uint64
	finalized         bool
	checkpoint        *schema.Checkpoint
	blockNumberLatest uint64
}

func (i *indexer) Run(ctx context.Context) (err error) {
	// Load checkpoint from database.
	if i.checkpoint, err = i.databaseClient.FindCheckpoint(ctx, i.chainID); err != nil {
		return fmt.Errorf("load checkpoint: %w", err)
	}

	retryableFunc := func() error {
		for {
			if err := i.index(ctx); err != nil {
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

func (i *indexer) index(ctx context.Context) (err error) {
	ctx, span := otel.Tracer("").Start(ctx, "index")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("chain.id", int64(i.chainID)),
		attribute.Int64("block.number.local", int64(i.checkpoint.BlockNumber)),
		attribute.Int64("block.number.latest", int64(i.blockNumberLatest)),
	)

	if err := i.refreshLatestBlockNumber(ctx); err != nil {
		return fmt.Errorf("get latest block number: %w", err)
	}

	zap.L().Info(
		"refreshed the latest block number",
		zap.Int("chain.id", int(i.chainID)),
		zap.Bool("finalized", i.finalized),
		zap.Uint64("block.number.local", i.checkpoint.BlockNumber),
		zap.Uint64("block.number.latest", i.blockNumberLatest),
	)

	// Waiting for a new block to be minted.
	if i.checkpoint.BlockNumber >= i.blockNumberLatest {
		blockConfirmationTime := time.Second // TODO Redefine it.

		zap.L().Info(
			"waiting for a new block to be minted",
			zap.Uint64("block.number.local", i.checkpoint.BlockNumber),
			zap.Uint64("block.number.latest", i.blockNumberLatest),
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

	blockNumber := i.checkpoint.BlockNumber + 1

	block, err := i.ethereumClient.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return fmt.Errorf("get block by number %d: %w", blockNumber, err)
	}

	receipts, err := i.ethereumClient.BlockReceipts(ctx, rpc.BlockNumberOrHashWithHash(block.Hash(), false))
	if err != nil {
		return err
	}

	// Begin a database transaction for the block.
	databaseTransaction, err := i.databaseClient.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin database transaction: %w", err)
	}

	defer lo.Try(databaseTransaction.Rollback)

	if err := i.handler.Process(ctx, block, receipts, databaseTransaction); err != nil {
		return fmt.Errorf("process block %d: %w", block.NumberU64(), err)
	}

	if i.finalized {
		// Update and save checkpoint to memory and database.
		i.checkpoint.BlockHash = block.Hash()
		i.checkpoint.BlockNumber = block.NumberU64()

		if err := databaseTransaction.SaveCheckpoint(ctx, i.checkpoint); err != nil {
			return fmt.Errorf("save checkpoint: %w", err)
		}
	}

	if databaseTransaction.Commit() != nil {
		return fmt.Errorf("commit database transaction: %w", err)
	}

	return nil
}

func (i *indexer) refreshLatestBlockNumber(ctx context.Context) (err error) {
	ctx, span := otel.Tracer("").Start(ctx, "refreshLatestBlockNumber")
	defer span.End()

	if i.finalized {
		block, err := i.ethereumClient.BlockByNumber(ctx, big.NewInt(rpc.FinalizedBlockNumber.Int64()))
		if err != nil {
			return fmt.Errorf("get finalized block number: %w", err)
		}

		i.blockNumberLatest = block.NumberU64()
	} else {
		if i.blockNumberLatest, err = i.ethereumClient.BlockNumber(ctx); err != nil {
			return fmt.Errorf("get latest block number: %w", err)
		}
	}

	span.SetAttributes(attribute.Int64("block.number", int64(i.blockNumberLatest)))

	return nil
}

func NewIndexer(chainID uint64, ethereumClient *ethclient.Client, databaseClient database.Client, handler Handler, finalized bool) (Indexer, error) {
	instance := indexer{
		ethereumClient: ethereumClient,
		databaseClient: databaseClient,
		handler:        handler,
		chainID:        chainID,
		finalized:      finalized,
	}

	return &instance, nil
}
