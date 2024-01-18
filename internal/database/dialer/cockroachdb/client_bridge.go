package cockroachdb

import (
	"context"
	"fmt"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

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
