package l1

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l1"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"go.uber.org/zap"
)

func (s *server) indexBillingLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, logIndex int, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; eventHash {
	case l1.EventHashBillingTokensDeposited:
		return s.indexBillingTokensDepositedLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	case l1.EventHashBillingTokensWithdrawn:
		return s.indexBillingTokensWithdrawnLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	case l1.EventHashBillingTokensCollected:
		return s.indexBillingTokensCollectedLog(ctx, header, transaction, receipt, log, logIndex, databaseTransaction)
	default:
		return nil
	}
}

func (s *server) indexBillingTokensDepositedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, _ int, databaseTransaction database.Client) error {
	billingTokensDepositedEvent, err := s.contractBilling.ParseTokensDeposited(*log)
	if err != nil {
		return fmt.Errorf("parse TokensDeposited event: %w", err)
	}

	zap.L().Debug("indexing TokensDeposited event for Billing", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", billingTokensDepositedEvent))

	billingRecord := schema.NewBillingRecordDeposited(header, transaction, receipt, billingTokensDepositedEvent.User, billingTokensDepositedEvent.Amount)

	if err := databaseTransaction.SaveBillingRecordDeposited(ctx, billingRecord); err != nil {
		return fmt.Errorf("save billing record: %w", err)
	}

	return nil
}

func (s *server) indexBillingTokensWithdrawnLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, _ int, databaseTransaction database.Client) error {
	billingTokensWithdrawnEvent, err := s.contractBilling.ParseTokensWithdrawn(*log)
	if err != nil {
		return fmt.Errorf("parse TokensWithdrawn event: %w", err)
	}

	zap.L().Debug("indexing TokensWithdrawn event for Billing", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", billingTokensWithdrawnEvent))

	billingRecord := schema.NewBillingRecordWithdrawal(header, transaction, receipt, billingTokensWithdrawnEvent.User, billingTokensWithdrawnEvent.Amount, billingTokensWithdrawnEvent.Fee)

	if err := databaseTransaction.SaveBillingRecordWithdrawal(ctx, billingRecord); err != nil {
		return fmt.Errorf("save billing record: %w", err)
	}

	return nil
}

func (s *server) indexBillingTokensCollectedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, _ int, databaseTransaction database.Client) error {
	billingTokensCollected, err := s.contractBilling.ParseTokensCollected(*log)
	if err != nil {
		return fmt.Errorf("parse TokensCollected event: %w", err)
	}

	zap.L().Debug("indexing TokensCollected event for Billing", zap.Stringer("transaction.hash", transaction.Hash()), zap.Any("event", billingTokensCollected))

	billingRecord := schema.NewBillingRecordCollected(header, transaction, receipt, billingTokensCollected.User, billingTokensCollected.Amount)

	if err := databaseTransaction.SaveBillingRecordCollected(ctx, billingRecord); err != nil {
		return fmt.Errorf("save billing record: %w", err)
	}

	return nil
}
