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

func (c *client) FindStakeTransaction(ctx context.Context, id common.Hash) (*schema.StakeTransaction, error) {
	var row table.StakeTransaction

	if err := c.database.WithContext(ctx).Where("id = ?", id.String()).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find stake transaction: %w", err)
	}

	return row.Export()
}

func (c *client) FindStakeTransactions(ctx context.Context) ([]*schema.StakeTransaction, error) {
	var rows []table.StakeTransaction

	if err := c.database.WithContext(ctx).Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (c *client) FindStakeTransactionsByUser(ctx context.Context, address common.Address) ([]*schema.StakeTransaction, error) {
	var rows []table.StakeTransaction

	if err := c.database.WithContext(ctx).Distinct("*").Where(`"user" = ?`, address.String()).Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (c *client) FindStakeTransactionsByNode(ctx context.Context, address common.Address) ([]*schema.StakeTransaction, error) {
	var rows []table.StakeTransaction

	if err := c.database.WithContext(ctx).Distinct("*").Where(`"node" = ?`, address.String()).Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (c *client) FindStakeEventsByID(ctx context.Context, id common.Hash) ([]*schema.StakeEvent, error) {
	var rows []table.StakeEvent

	if err := c.database.WithContext(ctx).Where("id = ?", id.String()).First(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find stake event: %w", err)
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

func (c *client) FindStakeEventsByIDs(ctx context.Context, ids []common.Hash) ([]*schema.StakeEvent, error) {
	var rows []table.StakeEvent

	transactionIDs := lo.Map(ids, func(id common.Hash, _ int) string {
		return id.String()
	})

	if err := c.database.WithContext(ctx).Where("id IN ?", transactionIDs).Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find stake event: %w", err)
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

func (c *client) FindStakeChipsByOwner(ctx context.Context, owner common.Address) ([]*schema.StakeChip, error) {
	if owner == ethereum.AddressGenesis {
		return make([]*schema.StakeChip, 0), nil
	}

	var rows []table.StakeChip

	if err := c.database.WithContext(ctx).Where(`"owner" = ?`, owner.String()).Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (c *client) FindStakeChipsByNode(ctx context.Context, node common.Address) ([]*schema.StakeChip, error) {
	var rows []table.StakeChip

	if err := c.database.WithContext(ctx).Where(`"node" = ? AND "owner" != ?`, node.String(), ethereum.AddressGenesis.String()).Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
