package cockroachdb

import (
	"context"
	"fmt"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

func (c *client) SaveBillingRecordDeposited(ctx context.Context, billingRecord *schema.BillingRecordDeposited) error {
	var value table.BillingRecordDeposited
	if err := value.Import(*billingRecord); err != nil {
		return fmt.Errorf("import billing record: %w", err)
	}

	return c.database.WithContext(ctx).Create(&value).Error
}

func (c *client) SaveBillingRecordWithdrawal(ctx context.Context, billingRecord *schema.BillingRecordWithdrawal) error {
	var value table.BillingRecordWithdrawal
	if err := value.Import(*billingRecord); err != nil {
		return fmt.Errorf("import billing record: %w", err)
	}

	return c.database.WithContext(ctx).Create(&value).Error
}

func (c *client) SaveBillingRecordCollected(ctx context.Context, billingRecord *schema.BillingRecordCollected) error {
	var value table.BillingRecordCollected
	if err := value.Import(*billingRecord); err != nil {
		return fmt.Errorf("import billing record: %w", err)
	}

	return c.database.WithContext(ctx).Create(&value).Error
}
