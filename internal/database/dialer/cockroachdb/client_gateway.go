package cockroachdb

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
)

func (c *client) ResumeGatewayAccount(ctx context.Context, address common.Address) (bool, error) {
	var account table.GatewayAccount

	err := c.database.WithContext(ctx).
		Model(&table.GatewayAccount{}).
		Where("address = ?", address).
		First(&account).
		Error

	if err != nil {
		// Failed to query account
		return false, err
	}

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
