package cockroachdb

import (
	"context"
	"fmt"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"gorm.io/gorm/clause"
)

func (c *client) FindCheckpoint(ctx context.Context, chainID uint64) (*schema.Checkpoint, error) {
	var checkpoint table.Checkpoint

	if err := c.database.
		WithContext(ctx).
		FirstOrInit(&checkpoint, table.Checkpoint{ChainID: chainID}).Error; err != nil {
		return nil, err
	}

	return checkpoint.Export()
}

func (c *client) SaveCheckpoint(ctx context.Context, checkpoint *schema.Checkpoint) error {
	var value table.Checkpoint
	if err := value.Import(*checkpoint); err != nil {
		return fmt.Errorf("import checkpoint: %w", err)
	}

	clauses := []clause.Expression{
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "chain_id"}},
			UpdateAll: true,
		},
	}

	return c.database.WithContext(ctx).Clauses(clauses...).Create(&value).Error
}
