package l2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/common/ethereum"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/schema"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (h *handler) indexChipsLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; {
	case h.finalized && eventHash == l2.EventHashChipsTransfer:
		return h.indexChipsTransferLog(ctx, header, transaction, receipt, log, databaseTransaction)
	default: // Discard all unsupported events.
		return nil
	}
}

func (h *handler) indexChipsTransferLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexChipsTransferLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractChips.ParseTransfer(*log)
	if err != nil {
		return fmt.Errorf("parse Transfer event: %w", err)
	}

	if h.finalized {
		if err := databaseTransaction.UpdateStakeChipsOwner(ctx, event.To, event.TokenId); err != nil {
			return fmt.Errorf("update stake chips owner: %w", err)
		}
	}

	if event.To != ethereum.AddressGenesis {
		return nil
	}

	metadata, err := json.Marshal(schema.StakeEventChipsMergedMetadata{ChipID: event.TokenId})
	if err != nil {
		return fmt.Errorf("marshal chips merged metadata: %w", err)
	}

	stakeEvent := schema.StakeEvent{
		ID:                transaction.Hash(),
		Type:              schema.StakeEventTypeChipsBurned,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		LogIndex:          log.Index,
		Metadata:          metadata,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
		Finalized:         h.finalized,
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}
