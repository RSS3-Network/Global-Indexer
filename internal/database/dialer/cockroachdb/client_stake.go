package cockroachdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
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
