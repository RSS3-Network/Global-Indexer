package l2

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

var _ service.Server = (*server)(nil)

type server struct {
	databaseClient                 database.Client
	ethereumClient                 *ethclient.Client
	chainID                        *big.Int
	contractGovernanceToken        *bindings.GovernanceToken
	contractL2CrossDomainMessenger *bindings.L2CrossDomainMessenger
	contractL2StandardBridge       *bindings.L2StandardBridge
	contractL2ToL1MessagePasser    *bindings.L2ToL1MessagePasser
	contractStaking                *l2.Staking
	contractChips                  *l2.Chips
	checkpoint                     *schema.Checkpoint
	blockNumberLatest              uint64
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

	return retry.Do(func() error { return s.run(ctx) }, retry.Delay(time.Second), retry.Attempts(30), onRetry, retryIf)
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

		blockNumberCurrent := s.checkpoint.BlockNumber + 1

		// Get current block (header and transactions).
		block, err := s.ethereumClient.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumberCurrent))
		if err != nil {
			return fmt.Errorf("get block %d: %w", blockNumberCurrent, err)
		}

		// Get all receipts of the current block.
		receipts, err := s.ethereumClient.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumberCurrent)))
		if err != nil {
			return fmt.Errorf("get receipts for block %d: %w", block.NumberU64(), err)
		}

		if err := s.index(ctx, block, receipts); err != nil {
			return fmt.Errorf("index block %d: %w", blockNumberCurrent, err)
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
			case l2.AddressL2StandardBridgeProxy:
				if err := s.indexBridgingLog(ctx, header, block.Transaction(log.TxHash), receipt, log, index, databaseTransaction); err != nil {
					return fmt.Errorf("index bridge log: %w", err)
				}
			case l2.AddressStakingProxy:
				if err := s.indexStakingLog(ctx, header, block.Transaction(log.TxHash), receipt, log, databaseTransaction); err != nil {
					return fmt.Errorf("index staking log: %w", err)
				}
			case l2.AddressChipsProxy:
				if err := s.indexChipsLog(ctx, header, block.Transaction(log.TxHash), receipt, log, databaseTransaction); err != nil {
					return fmt.Errorf("index staking log: %w", err)
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

func NewServer(ctx context.Context, databaseClient database.Client, config Config) (service.Server, error) {
	var (
		instance = server{
			databaseClient: databaseClient,
		}
		err error
	)

	if instance.ethereumClient, err = ethclient.DialContext(ctx, config.Endpoint); err != nil {
		return nil, err
	}

	if instance.chainID, err = instance.ethereumClient.ChainID(ctx); err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	if instance.contractGovernanceToken, err = bindings.NewGovernanceToken(l2.AddressGovernanceTokenProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractL2CrossDomainMessenger, err = bindings.NewL2CrossDomainMessenger(l2.AddressL2CrossDomainMessengerProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractL2StandardBridge, err = bindings.NewL2StandardBridge(l2.AddressL2StandardBridgeProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractL2ToL1MessagePasser, err = bindings.NewL2ToL1MessagePasser(l2.AddressL2ToL1MessagePasser, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractStaking, err = l2.NewStaking(l2.AddressStakingProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	if instance.contractChips, err = l2.NewChips(l2.AddressChipsProxy, instance.ethereumClient); err != nil {
		return nil, err
	}

	return &instance, nil
}
