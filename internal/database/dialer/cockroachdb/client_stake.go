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

func (c *client) FindStakeStaker(ctx context.Context, user, node common.Address) (*schema.StakeStaker, error) {
	var value table.StakeStaker

	if err := c.database.
		WithContext(ctx).
		Where(`"user" = ? AND "node" = ?`, user.String(), node.String()).
		First(&value).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		value = table.StakeStaker{
			User: user.String(),
			Node: node.String(),
		}
	}

	return value.Export()
}

func (c *client) SaveStakeStaker(ctx context.Context, stakeStaker *schema.StakeStaker) error {
	var value table.StakeStaker
	if err := value.Import(*stakeStaker); err != nil {
		return err
	}

	clauses := []clause.Expression{
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user"},
				{Name: "node"},
			},
			UpdateAll: true,
		},
	}

	return c.database.WithContext(ctx).Clauses(clauses...).Create(&value).Error
}

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
