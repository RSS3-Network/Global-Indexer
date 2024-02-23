package cockroachdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func (c *client) PrepareBillingCollectTokens(ctx context.Context, nowTime time.Time) (*map[common.Address]schema.BillingCollectDataPerAddress, error) {
	// Get all keys whose ru_used_current is > 0
	var activeKeys []table.GatewayKey

	err := c.database.WithContext(ctx).
		Model(&table.GatewayKey{}).
		Where("ru_used_current > ?", 0).
		Find(&activeKeys).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No active keys
			return nil, nil
		}

		// We don't know what happened, so let's just return this error
		return nil, fmt.Errorf("prepare billing collect: %w", err)
	}

	// Count into accounts
	ruConsumptions := make(map[common.Address]schema.BillingCollectDataPerAddress)

	for _, k := range activeKeys {
		// w/ database tx
		err = c.database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Create or update consumption log
			var possibleExistLog table.GatewayConsumptionLog
			err = tx.Where("consumption_date = ? AND key_id = ?", nowTime, k.ID).
				First(&possibleExistLog).
				Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Fine, let's create it
					err = tx.Create(&table.GatewayConsumptionLog{
						KeyID:           k.ID,
						ConsumptionDate: nowTime,
						RuUsed:          k.RuUsedCurrent,
						APICalls:        k.APICallsCurrent,
					}).Error
				} else {
					// Error happens, but we don't know what's this, create a new record for now.
					err = tx.Create(&table.GatewayConsumptionLog{
						KeyID:           k.ID,
						ConsumptionDate: nowTime,
						RuUsed:          k.RuUsedCurrent,
						APICalls:        k.APICallsCurrent,
					}).Error
				}
			} else {
				// Already exists - this shouldn't happen; but when it happens, it happens
				err = tx.Model(&table.GatewayConsumptionLog{}).
					Where("id = ?", possibleExistLog.ID).
					Updates(map[string]interface{}{
						"ru_used":   gorm.Expr("ru_used + ?", k.RuUsedCurrent),
						"api_calls": gorm.Expr("api_calls + ?", k.APICallsCurrent),
					}).Error
			}
			if err != nil {
				zap.L().Error("create or update consumption log", zap.Error(err), zap.Any("key", k))
				// but no need to stop here - data error can be fixed later, let's focus on billing now
			}

			// Reset current and update total usages
			err = tx.Model(&table.GatewayKey{}).
				Where("id = ?", k.ID).
				Updates(map[string]interface{}{
					"ru_used_total":     gorm.Expr("ru_used_total + ru_used_current"),
					"ru_used_current":   0,
					"api_calls_total":   gorm.Expr("api_calls_total + api_calls_current"),
					"api_calls_current": 0,
				}).
				Error

			if err != nil {
				zap.L().Error("reset usage", zap.Error(err), zap.Any("key", k))
				// this is unacceptable, it may influence billing. abort tx and rollback.
				return err
			}

			// Commit
			return nil
		})

		if err != nil {
			// Something wrong in database transaction, skip this key
			continue
		}

		// Initialize account RU counter
		if _, exist := ruConsumptions[k.Account.Address]; !exist {
			ruConsumptions[k.Account.Address] = schema.BillingCollectDataPerAddress{
				Ru:          0,
				BillingRate: k.Account.BillingRate,
			}
		}

		// Finally we can count ru used records
		ruConsumptions[k.Account.Address] = schema.BillingCollectDataPerAddress{
			Ru:          ruConsumptions[k.Account.Address].Ru + k.RuUsedCurrent,
			BillingRate: ruConsumptions[k.Account.Address].BillingRate,
		}
	}

	return &ruConsumptions, nil
}

func (c *client) PrepareBillingWithdrawTokens(ctx context.Context) (*map[common.Address]float64, error) {
	var billingWithdrawRequests []table.GatewayPendingWithdrawRequest

	err := c.database.WithContext(ctx).
		Find(&billingWithdrawRequests).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No pending request withdraw
			return nil, nil
		}

		// We don't know what happened, so let's just return this error
		return nil, fmt.Errorf("prepare billing withdraw: %w", err)
	}

	withdrawRequests := make(map[common.Address]float64)

	for _, req := range billingWithdrawRequests {
		// Append
		withdrawRequests[req.AccountAddress] = req.Amount
	}

	// Delete
	err = c.database.WithContext(ctx).
		Delete(&billingWithdrawRequests).
		Error

	if err != nil {
		// wait, why delete failed?
		zap.L().Error("delete billing withdraw", zap.Error(err))
	}

	return &withdrawRequests, nil
}

func (c *client) UpdateBillingRuLimit(ctx context.Context, succeededUsersWithRu map[common.Address]int64) error {
	for address, ruLimit := range succeededUsersWithRu {
		c.database.WithContext(ctx).
			Where("address = ?", address).
			Update("ru_limit", ruLimit)
	}

	return nil
}
