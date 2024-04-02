package cockroachdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) FindBridgeTransaction(ctx context.Context, query schema.BridgeTransactionQuery) (*schema.BridgeTransaction, error) {
	var row *table.BridgeTransaction

	databaseClient := c.database.WithContext(ctx)

	if query.ID != nil {
		databaseClient = databaseClient.Where(`"id" = ?`, query.ID.String())
	}

	if query.Sender != nil {
		databaseClient = databaseClient.Where(`"sender" = ?`, query.Sender.String())
	}

	if query.Receiver != nil {
		databaseClient = databaseClient.Where(`"receiver" = ?`, query.Receiver.String())
	}

	if query.Address != nil {
		databaseClient = databaseClient.Where(`"sender" = ? or "receiver" = ?`, query.Address.String(), query.Address.String())
	}

	if query.Type != nil {
		databaseClient = databaseClient.Where(`"type" = ?`, *query.Type)
	}

	if err := databaseClient.Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, err
	}

	result, err := row.Export()
	if err != nil {
		return nil, fmt.Errorf("export row: %w", err)
	}

	return result, nil
}

func (c *client) FindBridgeTransactions(ctx context.Context, query schema.BridgeTransactionsQuery) ([]*schema.BridgeTransaction, error) {
	var rows []table.BridgeTransaction

	databaseClient := c.database.WithContext(ctx)

	const limit = 100

	if query.Cursor != nil {
		var cursor table.BridgeTransaction
		if err := databaseClient.Where(`"id" = ?`, query.Cursor.String()).First(&cursor).Error; err != nil {
			return nil, fmt.Errorf("query cursor: %w", err)
		}

		// TODO Need a better cursor implementation.
		databaseClient = databaseClient.Where(
			`
("block_timestamp" < ?) OR
("block_timestamp" = ? AND "chain_id" = ? AND "block_number" < ?) OR
("block_timestamp" = ? AND "chain_id" = ? AND "block_number" = ? AND "transaction_index" < ?)
`,
			cursor.BlockTimestamp,
			cursor.BlockTimestamp, cursor.ChainID, cursor.BlockNumber,
			cursor.BlockTimestamp, cursor.ChainID, cursor.BlockNumber, cursor.TransactionIndex,
		)
	}

	if query.ID != nil {
		databaseClient = databaseClient.Where(`"id" = ?`, query.ID.String())
	}

	if query.Sender != nil {
		databaseClient = databaseClient.Where(`"sender" = ?`, query.Sender.String())
	}

	if query.Receiver != nil {
		databaseClient = databaseClient.Where(`"receiver" = ?`, query.Receiver.String())
	}

	if query.Address != nil {
		databaseClient = databaseClient.Where(`"sender" = ? or "receiver" = ?`, query.Address.String(), query.Address.String())
	}

	if query.Type != nil {
		databaseClient = databaseClient.Where(`"type" = ?`, *query.Type)
	}

	if err := databaseClient.Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`).Limit(limit).Find(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, err
	}

	results := make([]*schema.BridgeTransaction, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export row: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) FindBridgeEvents(ctx context.Context, query schema.BridgeEventsQuery) ([]*schema.BridgeEvent, error) {
	var rows []*table.BridgeEvent

	databaseClient := c.database.WithContext(ctx)

	if len(query.IDs) > 0 {
		databaseClient = databaseClient.Where(`"id" IN ?`, lo.Map(query.IDs, func(id common.Hash, _ int) string {
			return id.String()
		}))
	}

	if err := databaseClient.Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`).Find(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find bridge event: %w", err)
	}

	results := make([]*schema.BridgeEvent, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export row: %w", err)
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

	clauses := []clause.Expression{
		clause.OnConflict{
			UpdateAll: true,
		},
	}

	return c.database.WithContext(ctx).Clauses(clauses...).Create(&value).Error
}

func (c *client) SaveBridgeEvent(ctx context.Context, bridgeEvent *schema.BridgeEvent) error {
	var value table.BridgeEvent
	if err := value.Import(*bridgeEvent); err != nil {
		return fmt.Errorf("import bridge event: %w", err)
	}

	clauses := []clause.Expression{
		clause.OnConflict{
			UpdateAll: true,
		},
	}

	return c.database.WithContext(ctx).Clauses(clauses...).Create(&value).Error
}
