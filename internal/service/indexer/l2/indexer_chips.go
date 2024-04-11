package l2

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *server) indexChipsLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; eventHash {
	case l2.EventHashChipsTransfer:
		return s.indexChipsTransferLog(ctx, header, transaction, receipt, log, databaseTransaction)
	default: // Discard all unsupported events.
		return nil
	}
}

func (s *server) indexChipsTransferLog(ctx context.Context, header *types.Header, transaction *types.Transaction, _ *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexChipsTransferLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := s.contractChips.ParseTransfer(*log)
	if err != nil {
		return fmt.Errorf("parse Transfer event: %w", err)
	}

	if err := databaseTransaction.UpdateStakeChipsOwner(ctx, event.To, event.TokenId); err != nil {
		return fmt.Errorf("update stake chips owner: %w", err)
	}

	return nil
}
