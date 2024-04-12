package l2

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (s *server) indexBridgingLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, logIndex int, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; eventHash {
	case l2.EventHashL2StandardBridgeDepositFinalized:
		return s.indexL2StandardBridgeDepositFinalizedLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	case l2.EventHashL2StandardBridgeWithdrawalInitiated:
		return s.indexL2StandardWithdrawalInitiatedLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	default: // Discard all unsupported event.
		return nil
	}
}

func (s *server) indexL2StandardBridgeDepositFinalizedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, logIndex int, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexL2StandardBridgeDepositFinalizedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	depositFinalizedEvent, err := s.contractL2StandardBridge.ParseDepositFinalized(*log)
	if err != nil {
		return fmt.Errorf("parse DepositFinalized event: %w", err)
	}

	zap.L().Debug("indexing DepositFinalize event for L2StandardBridge", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", depositFinalizedEvent))

	// Try to match bridge transaction id.
	var relayedMessageEvent *bindings.L2CrossDomainMessengerRelayedMessage

	for _, log := range receipt.Logs[logIndex:] {
		if log.Address == l2.AddressL2CrossDomainMessengerProxy {
			if len(log.Topics) > 0 && log.Topics[0] == l2.EventHashL2CrossDomainMessengerRelayedMessage {
				if relayedMessageEvent, err = s.contractL2CrossDomainMessenger.ParseRelayedMessage(*log); err != nil {
					return fmt.Errorf("parse RelayedMessage event: %w", err)
				}
			}
		}
	}

	if relayedMessageEvent == nil {
		return fmt.Errorf("no matched RelayedMessage event")
	}

	// Create the bridge event.
	bridgeEvent := schema.NewBridgeEvent(relayedMessageEvent.MsgHash, schema.BridgeEventTypeDepositFinalized, s.chainID.Uint64(), header, transaction, receipt)

	if err := databaseTransaction.SaveBridgeEvent(ctx, bridgeEvent); err != nil {
		return fmt.Errorf("save bridge transaction: %w", err)
	}

	return nil
}

func (s *server) indexL2StandardWithdrawalInitiatedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, logIndex int, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexL2StandardWithdrawalInitiatedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	withdrawalInitiatedEvent, err := s.contractL2StandardBridge.ParseWithdrawalInitiated(*log)
	if err != nil {
		return fmt.Errorf("parse WithdrawalInitiated event: %w", err)
	}

	zap.L().Debug("indexing WithdrawalInitiated event for L2StandardBridge", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", withdrawalInitiatedEvent))

	// Try to match bridge transaction id.
	var messagePassedEvent *bindings.L2ToL1MessagePasserMessagePassed

	for _, log := range receipt.Logs[logIndex:] {
		if log.Address == l2.AddressL2ToL1MessagePasser {
			if len(log.Topics) > 0 && log.Topics[0] == l2.EventHashL2ToL1MessagePasserMessagePassed {
				if messagePassedEvent, err = s.contractL2ToL1MessagePasser.ParseMessagePassed(*log); err != nil {
					return fmt.Errorf("parse MessagePassed event: %w", err)
				}
			}
		}
	}

	if messagePassedEvent == nil {
		return fmt.Errorf("no matched MessagePassed event")
	}

	// Create the bridge transaction.
	bridgeTransaction := schema.BridgeTransaction{
		ID:               messagePassedEvent.WithdrawalHash,
		Type:             schema.BridgeTransactionTypeWithdraw,
		Sender:           withdrawalInitiatedEvent.From,
		Receiver:         withdrawalInitiatedEvent.To,
		TokenAddressL1:   lo.ToPtr(withdrawalInitiatedEvent.L1Token),
		TokenAddressL2:   lo.ToPtr(withdrawalInitiatedEvent.L2Token),
		TokenValue:       withdrawalInitiatedEvent.Amount,
		Data:             hexutil.Encode(messagePassedEvent.Data),
		ChainID:          s.chainID.Uint64(),
		BlockTimestamp:   time.Unix(int64(header.Time), 0),
		BlockNumber:      header.Number.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
	}

	if err := databaseTransaction.SaveBridgeTransaction(ctx, &bridgeTransaction); err != nil {
		return fmt.Errorf("save bridge transaction: %w", err)
	}

	// Create the bridge event.
	bridgeEvent := schema.NewBridgeEvent(messagePassedEvent.WithdrawalHash, schema.BridgeEventTypeWithdrawalInitialized, s.chainID.Uint64(), header, transaction, receipt)

	if err := databaseTransaction.SaveBridgeEvent(ctx, bridgeEvent); err != nil {
		return fmt.Errorf("save bridge event: %w", err)
	}

	return nil
}
