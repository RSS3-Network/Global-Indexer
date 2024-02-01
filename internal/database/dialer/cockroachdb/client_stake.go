package cockroachdb

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func (c *client) FindStakeTransaction(ctx context.Context, query schema.StakeTransactionQuery) (*schema.StakeTransaction, error) {
	var row table.StakeTransaction

	databaseClient := c.database.WithContext(ctx)

	if query.ID != nil {
		databaseClient = databaseClient.Where(`"id" = ?`, query.ID.String())
	}

	if query.User != nil {
		databaseClient = databaseClient.Where(`"user" = ?`, query.User.String())
	}

	if query.Node != nil {
		databaseClient = databaseClient.Where(`"node" = ?`, query.Node.String())
	}

	if query.Address != nil {
		databaseClient = databaseClient.Where(`"user" = ? OR "node" = ?`, query.Address.String())
	}

	if query.Type != nil {
		databaseClient = databaseClient.Where(`"type" = ?`, query.Type)
	}

	if err := databaseClient.Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find stake transaction: %w", err)
	}

	return row.Export()
}

func (c *client) FindStakeTransactions(ctx context.Context, query schema.StakeTransactionsQuery) ([]*schema.StakeTransaction, error) {
	databaseClient := c.database.WithContext(ctx)

	if query.IDs != nil {
		databaseClient = databaseClient.Where(`"id" = ?`, lo.Map(query.IDs, func(id common.Hash, _ int) string {
			return id.String()
		}))
	}

	if query.User != nil {
		databaseClient = databaseClient.Where(`"user" = ?`, query.User.String())
	}

	if query.Node != nil {
		databaseClient = databaseClient.Where(`"node" = ?`, query.Node.String())
	}

	if query.Address != nil {
		databaseClient = databaseClient.Where(`"user" = ? OR "node" = ?`, query.Address.String())
	}

	if query.Type != nil {
		databaseClient = databaseClient.Where(`"type" = ?`, query.Type)
	}

	var rows []table.StakeTransaction

	if err := databaseClient.Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`).Find(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find stake transactions: %w", err)
	}

	results := make([]*schema.StakeTransaction, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export stake transaction: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) FindStakeEvents(ctx context.Context, query schema.StakeEventsQuery) ([]*schema.StakeEvent, error) {
	databaseClient := c.database.WithContext(ctx)

	if len(query.IDs) > 0 {
		databaseClient.Where(`"id" IN ?`, lo.Map(query.IDs, func(id common.Hash, _ int) string {
			return id.String()
		}))
	}

	var rows []table.StakeEvent
	if err := c.database.WithContext(ctx).Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`).Find(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find stake events: %w", err)
	}

	results := make([]*schema.StakeEvent, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export stake event: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) FindStakeChips(ctx context.Context, query schema.StakeChipsQuery) ([]*schema.StakeChip, error) {
	databaseClient := c.database.WithContext(ctx)

	if query.ID != nil {
		databaseClient = databaseClient.Where(`"id" = ?`, query.ID.String())
	}

	if query.Owner != nil {
		databaseClient = databaseClient.Where(`"owner" = ?`, query.Owner.String())
	}

	if query.Node != nil {
		databaseClient = databaseClient.Where(`"node" = ?`, query.Node.String())
	}

	databaseClient = databaseClient.Where(`"owner" != ?`, ethereum.AddressGenesis.String())

	var rows []table.StakeChip
	if err := databaseClient.Order(`"id" DESC`).Find(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find stake chip: %w", err)
	}

	results := make([]*schema.StakeChip, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export stake chip: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) SaveStakeTransaction(ctx context.Context, stakeTransaction *schema.StakeTransaction) error {
	var value table.StakeTransaction
	if err := value.Import(*stakeTransaction); err != nil {
		return fmt.Errorf("import stake transaction: %w", err)
	}

	return c.database.WithContext(ctx).Create(&value).Error
}

func (c *client) SaveStakeEvent(ctx context.Context, stakeEvent *schema.StakeEvent) error {
	var value table.StakeEvent
	if err := value.Import(*stakeEvent); err != nil {
		return fmt.Errorf("import stake event: %w", err)
	}

	return c.database.WithContext(ctx).Create(&value).Error
}

func (c *client) SaveStakeChips(ctx context.Context, stakeChips ...*schema.StakeChip) error {
	values := make([]*table.StakeChip, 0, len(stakeChips))

	for _, stakeChip := range stakeChips {
		var value table.StakeChip

		if err := value.Import(*stakeChip); err != nil {
			return fmt.Errorf("import stake chip: %w", err)
		}

		values = append(values, &value)
	}

	return c.database.WithContext(ctx).Create(&values).Error
}

func (c *client) UpdateStakeChipsOwner(ctx context.Context, owner common.Address, stakeChipIDs ...*big.Int) error {
	ids := lo.Map(stakeChipIDs, func(stakeChipID *big.Int, _ int) decimal.Decimal {
		return decimal.NewFromBigInt(stakeChipID, 0)
	})

	return c.database.WithContext(ctx).Model((*table.StakeChip)(nil)).Where(`"id" IN ?`, ids).UpdateColumn("owner", owner.String()).Error
}
