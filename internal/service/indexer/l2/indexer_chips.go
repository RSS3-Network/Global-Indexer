package l2

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
)

func (s *server) indexChipsLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; eventHash {
	case l2.EventHashChipsTransfer:
		return s.indexChipsTransferLog(ctx, header, transaction, receipt, log, databaseTransaction)
	default: // Discard all unsupported events.
		return nil
	}
}

func (s *server) indexChipsTransferLog(ctx context.Context, _ *types.Header, _ *types.Transaction, _ *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractChips.ParseTransfer(*log)
	if err != nil {
		return fmt.Errorf("parse Transfer event: %w", err)
	}

	if err := databaseTransaction.UpdateStakeChipsOwner(ctx, event.To, event.TokenId); err != nil {
		return fmt.Errorf("update stake chips owner: %w", err)
	}

	return nil
}
