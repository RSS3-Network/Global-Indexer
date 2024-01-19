package cockroachdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func (c *client) FindBridgeTransactions(ctx context.Context) ([]*schema.BridgeTransaction, error) {
	var rows []table.BridgeTransaction

	if err := c.database.WithContext(ctx).Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("fin bridge transactions: %w", err)
	}

	results := make([]*schema.BridgeTransaction, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export bridge transaction: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) FindBridgeTransactionsByAddress(ctx context.Context, address common.Address) ([]*schema.BridgeTransaction, error) {
	var rows []table.BridgeTransaction

	if err := c.database.WithContext(ctx).Distinct("*").Where("sender = ? OR receiver = ?", address.String(), address.String()).Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("fin bridge transactions: %w", err)
	}

	results := make([]*schema.BridgeTransaction, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export bridge transaction: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) FindBridgeEventsByIDs(ctx context.Context, ids []common.Hash) ([]*schema.BridgeEvent, error) {
	var rows []table.BridgeEvent

	transactionIDs := lo.Map(ids, func(id common.Hash, _ int) string {
		return id.String()
	})

	if err := c.database.WithContext(ctx).Where("id IN ?", transactionIDs).Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("fin bridge event: %w", err)
	}

	results := make([]*schema.BridgeEvent, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export bridge event: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) SaveBridgeTransaction(ctx context.Context, bridgeTransaction *schema.BridgeTransaction) error {
	var value table.BridgeTransaction
	if err := value.Import(*bridgeTransaction); err != nil {
		return fmt.Errorf("import bridge transaction: %w", err)
	}

	return c.database.WithContext(ctx).Create(&value).Error
}

func (c *client) SaveBridgeEvent(ctx context.Context, bridgeEvent *schema.BridgeEvent) error {
	var value table.BridgeEvent
	if err := value.Import(*bridgeEvent); err != nil {
		return fmt.Errorf("import bridge event: %w", err)
	}

	return c.database.WithContext(ctx).Create(&value).Error
}
