package l1

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum-optimism/optimism/op-chain-ops/crossdomain"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l1"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (s *server) indexBridgingLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, logIndex int, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; eventHash {
	case l1.EventHashL1StandardBridgeERC20DepositInitiated:
		return s.indexL1StandardBridgeERC20DepositInitiatedLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	case l1.EventHashL1StandardBridgeERC20WithdrawalFinalized:
		return s.indexL1StandardBridgeERC20WithdrawalFinalizedLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	case l1.EventHashOptimismPortalWithdrawalProven:
		return s.indexOptimismPortalWithdrawalProvenLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	default:
		return nil
	}
}

func (s *server) indexL1StandardBridgeERC20DepositInitiatedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, logIndex int, databaseTransaction database.Client) error {
	erc20DepositInitiatedEvent, err := s.contractL1StandardBridge.ParseERC20DepositInitiated(*log)
	if err != nil {
		return fmt.Errorf("parse ERC20DepositInitiated event: %w", err)
	}

	zap.L().Debug("indexing ERC20DepositInitiated event for L1StandardBridge", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", erc20DepositInitiatedEvent))

	// Try to match bridge transaction id.
	var sentMessageEvent *bindings.L1CrossDomainMessengerSentMessage

	for _, log := range receipt.Logs[logIndex:] {
		if log.Address == l1.ContractMap[s.chainID.Uint64()].AddressL1CrossDomainMessengerProxy {
			if len(log.Topics) > 0 && log.Topics[0] == l1.EventHashL1CrossDomainMessengerSentMessage {
				if sentMessageEvent, err = s.contractL1CrossDomainMessenger.ParseSentMessage(*log); err != nil {
					return fmt.Errorf("parse SentMessage event: %w", err)
				}
			}
		}
	}

	if sentMessageEvent == nil {
		return fmt.Errorf("no matched SentMessage event")
	}

	messageHash, err := crossdomain.HashCrossDomainMessageV1(sentMessageEvent.MessageNonce, sentMessageEvent.Sender, sentMessageEvent.Target, big.NewInt(0) /* TODO Refactor it. */, sentMessageEvent.GasLimit, sentMessageEvent.Message)
	if err != nil {
		return fmt.Errorf("hash cross domain message: %w", err)
	}

	// Create the bridge transaction.
	bridgeTransaction := schema.BridgeTransaction{
		ID:               messageHash,
		Type:             schema.BridgeTransactionTypeDeposit,
		Sender:           erc20DepositInitiatedEvent.From,
		Receiver:         erc20DepositInitiatedEvent.To,
		TokenAddressL1:   lo.ToPtr(erc20DepositInitiatedEvent.L1Token),
		TokenAddressL2:   lo.ToPtr(erc20DepositInitiatedEvent.L2Token),
		TokenValue:       erc20DepositInitiatedEvent.Amount,
		Data:             hexutil.Encode(erc20DepositInitiatedEvent.ExtraData),
		ChainID:          s.chainID.Uint64(),
		BlockTimestamp:   time.Unix(int64(header.Time), 0),
		BlockNumber:      header.Number.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
	}

	if err := databaseTransaction.SaveBridgeTransaction(ctx, &bridgeTransaction); err != nil {
		return fmt.Errorf("save bridge transaction: %w", err)
	}

	// Create the bridge event.
	bridgeEvent := schema.NewBridgeEvent(messageHash, schema.BridgeEventTypeDepositInitialized, s.chainID.Uint64(), header, transaction, receipt)

	if err := databaseTransaction.SaveBridgeEvent(ctx, bridgeEvent); err != nil {
		return fmt.Errorf("save bridge event: %w", err)
	}

	return nil
}

func (s *server) indexL1StandardBridgeERC20WithdrawalFinalizedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, logIndex int, databaseTransaction database.Client) (err error) {
	withdrawalFinalizedEvent, err := s.contractL1StandardBridge.ParseERC20WithdrawalFinalized(*log)
	if err != nil {
		return fmt.Errorf("parse ERC20DepositInitiated event: %w", err)
	}

	zap.L().Debug("indexing ERC20DepositInitiated event for L1StandardBridge", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", withdrawalFinalizedEvent))

	// Try to match bridge transaction id.
	var optimismPortalEvent *bindings.OptimismPortalWithdrawalFinalized

	for _, log := range receipt.Logs[logIndex:] {
		if log.Address == l1.ContractMap[s.chainID.Uint64()].AddressOptimismPortalProxy {
			if len(log.Topics) > 0 && log.Topics[0] == l1.EventHashOptimismPortalWithdrawalFinalized {
				if optimismPortalEvent, err = s.contractOptimismPortal.ParseWithdrawalFinalized(*log); err != nil {
					return fmt.Errorf("parse RelayedMessage event: %w", err)
				}
			}
		}
	}

	if optimismPortalEvent == nil {
		return fmt.Errorf("no matched OptimismPortal WithdrawalFinalized event")
	}

	// Create the bridge event.
	bridgeEvent := schema.NewBridgeEvent(optimismPortalEvent.WithdrawalHash, schema.BridgeEventTypeWithdrawalFinalized, s.chainID.Uint64(), header, transaction, receipt)

	if err := databaseTransaction.SaveBridgeEvent(ctx, bridgeEvent); err != nil {
		return fmt.Errorf("save bridge event: %w", err)
	}

	return nil
}

func (s *server) indexOptimismPortalWithdrawalProvenLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, _ int, databaseTransaction database.Client) (err error) {
	withdrawalProvenEvent, err := s.contractOptimismPortal.ParseWithdrawalProven(*log)
	if err != nil {
		return fmt.Errorf("parse WithdrawalProven event: %w", err)
	}

	zap.L().Debug("indexing WithdrawalProven event for OptimismPortal", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", withdrawalProvenEvent))

	// Create the bridge event.
	bridgeEvent := schema.NewBridgeEvent(withdrawalProvenEvent.WithdrawalHash, schema.BridgeEventTypeWithdrawalProved, s.chainID.Uint64(), header, transaction, receipt)

	if err := databaseTransaction.SaveBridgeEvent(ctx, bridgeEvent); err != nil {
		return fmt.Errorf("save bridge event: %w", err)
	}

	return nil
}
