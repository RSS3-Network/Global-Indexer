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
	"github.com/rss3-network/global-indexer/contract/l1"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (h *handler) indexBridgingLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, logIndex int, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; eventHash {
	case l1.EventHashL1StandardBridgeERC20DepositInitiated:
		return h.indexL1StandardBridgeERC20DepositInitiatedLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	case l1.EventHashL1StandardBridgeERC20WithdrawalFinalized:
		return h.indexL1StandardBridgeERC20WithdrawalFinalizedLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	case l1.EventHashOptimismPortalWithdrawalProven:
		return h.indexOptimismPortalWithdrawalProvenLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	default:
		return nil
	}
}

func (h *handler) indexL1StandardBridgeERC20DepositInitiatedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, logIndex int, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexL1StandardBridgeERC20DepositInitiatedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	erc20DepositInitiatedEvent, err := h.contractL1StandardBridge.ParseERC20DepositInitiated(*log)
	if err != nil {
		return fmt.Errorf("parse ERC20DepositInitiated event: %w", err)
	}

	zap.L().Debug("indexing ERC20DepositInitiated event for L1StandardBridge", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", erc20DepositInitiatedEvent))

	// Try to match bridge transaction id.
	var sentMessageEvent *bindings.L1CrossDomainMessengerSentMessage

	for _, log := range receipt.Logs[logIndex:] {
		if log.Address == l1.ContractMap[h.chainID].AddressL1CrossDomainMessengerProxy {
			if len(log.Topics) > 0 && log.Topics[0] == l1.EventHashL1CrossDomainMessengerSentMessage {
				if sentMessageEvent, err = h.contractL1CrossDomainMessenger.ParseSentMessage(*log); err != nil {
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
		ChainID:          h.chainID,
		BlockTimestamp:   time.Unix(int64(header.Time), 0),
		BlockNumber:      header.Number.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
		Finalized:        h.finalized,
	}

	if err := databaseTransaction.SaveBridgeTransaction(ctx, &bridgeTransaction); err != nil {
		return fmt.Errorf("save bridge transaction: %w", err)
	}

	// Create the bridge event.
	bridgeEvent := schema.NewBridgeEvent(messageHash, schema.BridgeEventTypeDepositInitialized, h.chainID, header, transaction, receipt, h.finalized)

	if err := databaseTransaction.SaveBridgeEvent(ctx, bridgeEvent); err != nil {
		return fmt.Errorf("save bridge event: %w", err)
	}

	return nil
}

func (h *handler) indexL1StandardBridgeERC20WithdrawalFinalizedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, logIndex int, databaseTransaction database.Client) (err error) {
	ctx, span := otel.Tracer("").Start(ctx, "indexL1StandardBridgeERC20WithdrawalFinalizedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	withdrawalFinalizedEvent, err := h.contractL1StandardBridge.ParseERC20WithdrawalFinalized(*log)
	if err != nil {
		return fmt.Errorf("parse ERC20DepositInitiated event: %w", err)
	}

	zap.L().Debug("indexing ERC20DepositInitiated event for L1StandardBridge", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", withdrawalFinalizedEvent))

	// Try to match bridge transaction id.
	var optimismPortalEvent *bindings.OptimismPortalWithdrawalFinalized

	for _, log := range receipt.Logs[logIndex:] {
		if log.Address == l1.ContractMap[h.chainID].AddressOptimismPortalProxy {
			if len(log.Topics) > 0 && log.Topics[0] == l1.EventHashOptimismPortalWithdrawalFinalized {
				if optimismPortalEvent, err = h.contractOptimismPortal.ParseWithdrawalFinalized(*log); err != nil {
					return fmt.Errorf("parse RelayedMessage event: %w", err)
				}
			}
		}
	}

	if optimismPortalEvent == nil {
		return fmt.Errorf("no matched OptimismPortal WithdrawalFinalized event")
	}

	// Create the bridge event.
	bridgeEvent := schema.NewBridgeEvent(optimismPortalEvent.WithdrawalHash, schema.BridgeEventTypeWithdrawalFinalized, h.chainID, header, transaction, receipt, h.finalized)

	if err := databaseTransaction.SaveBridgeEvent(ctx, bridgeEvent); err != nil {
		return fmt.Errorf("save bridge event: %w", err)
	}

	return nil
}

func (h *handler) indexOptimismPortalWithdrawalProvenLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, _ int, databaseTransaction database.Client) (err error) {
	ctx, span := otel.Tracer("").Start(ctx, "indexOptimismPortalWithdrawalProvenLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	withdrawalProvenEvent, err := h.contractOptimismPortal.ParseWithdrawalProven(*log)
	if err != nil {
		return fmt.Errorf("parse WithdrawalProven event: %w", err)
	}

	zap.L().Debug("indexing WithdrawalProven event for OptimismPortal", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", withdrawalProvenEvent))

	// Create the bridge event.
	bridgeEvent := schema.NewBridgeEvent(withdrawalProvenEvent.WithdrawalHash, schema.BridgeEventTypeWithdrawalProved, h.chainID, header, transaction, receipt, h.finalized)

	if err := databaseTransaction.SaveBridgeEvent(ctx, bridgeEvent); err != nil {
		return fmt.Errorf("save bridge event: %w", err)
	}

	return nil
}
