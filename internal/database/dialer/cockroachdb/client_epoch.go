package cockroachdb

import (
	"context"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

func (c *client) SaveEpoch(ctx context.Context, epoch *schema.Epoch) error {
	// Save epoch.
	var data table.Epoch
	if err := data.Import(epoch); err != nil {
		zap.L().Error("import epoch", zap.Error(err), zap.Any("epoch", epoch))

		return err
	}

	onConflict := clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}

	if err := c.database.WithContext(ctx).Clauses(onConflict).Create(&data).Error; err != nil {
		zap.L().Error("insert epoch", zap.Error(err), zap.Any("epoch", epoch))

		return err
	}

	// Save epoch items.
	var items table.EpochItems
	if err := items.Import(epoch.RewardItems); err != nil {
		zap.L().Error("import epoch items", zap.Error(err), zap.Any("epoch", epoch))

		return err
	}

	onConflict = clause.OnConflict{
		Columns: []clause.Column{
			{Name: "epoch_id"},
			{Name: "index"},
		},
		UpdateAll: true,
	}

	if err := c.database.WithContext(ctx).Clauses(onConflict).CreateInBatches(&items, 500).Error; err != nil {
		zap.L().Error("insert epoch items", zap.Error(err), zap.Any("epoch", epoch))

		return err
	}

	return nil
}
