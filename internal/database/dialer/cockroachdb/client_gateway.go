package cockroachdb

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"gorm.io/gorm"
)

func (c *client) GatewayDeposit(ctx context.Context, address common.Address, ruIncrease int64) (bool, error) {
	account := table.GatewayAccount{
		Address: address,
	}

	// Get account
	err := c.database.WithContext(ctx).
		FirstOrCreate(&account).
		Error

	if err != nil {
		// Failed to get account
		return false, err
	}

	// Increase RU
	err = c.database.WithContext(ctx).
		Model(&table.GatewayAccount{}).
		Where("address = ?", address).
		Update("ru_limit", gorm.Expr("ru_limit + ?", ruIncrease)).
		Error

	if err != nil {
		// Failed to increase RU
		return false, err
	}

	// Check if account has been paused
	if !account.IsPaused {
		// Not paused
		return false, nil
	}

	// else is paused, resume account
	err = c.database.WithContext(ctx).
		Model(&table.GatewayAccount{}).
		Where("address = ?", address).
		Update("is_paused", false).
		Error

	if err != nil {
		// Failed to update account
		return true, err
	}

	// else has no error, return true as account has been resumed
	return true, nil
}
