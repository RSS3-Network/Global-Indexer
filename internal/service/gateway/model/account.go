package model

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/apisix"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"gorm.io/gorm"
)

type Account struct {
	table.GatewayAccount

	databaseClient *gorm.DB
	apisixClient   *apisix.Client
}

func AccountCreate(ctx context.Context, address common.Address, databaseClient *gorm.DB, apisixClient *apisix.Client) (*Account, error) {
	acc := table.GatewayAccount{
		Address: address,
	}
	err := databaseClient.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// DB
		err := tx.
			Save(&acc).
			Error
		if err != nil {
			return err
		}
		// APISix
		err = apisixClient.NewConsumerGroup(ctx, address.Hex())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Account{acc, databaseClient, apisixClient}, nil
}

func AccountGetByAddress(ctx context.Context, address common.Address, databaseClient *gorm.DB, apisixClient *apisix.Client) (*Account, bool, error) {
	var acc table.GatewayAccount

	err := databaseClient.WithContext(ctx).
		Model(&table.GatewayAccount{}).
		Where("address = ?", address).
		First(&acc).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &Account{acc, databaseClient, apisixClient}, true, nil
}

func AccountGetOrCreate(ctx context.Context, address common.Address, databaseClient *gorm.DB, apisixClient *apisix.Client) (*Account, error) {
	acc, exist, err := AccountGetByAddress(ctx, address, databaseClient, apisixClient)

	if err != nil {
		return nil, err
	} else if !exist {
		return AccountCreate(ctx, address, databaseClient, apisixClient)
	}

	return acc, nil
}

func (acc *Account) ListKeys(ctx context.Context) ([]*Key, error) {
	var keys []table.GatewayKey

	err := acc.databaseClient.WithContext(ctx).
		Model(&table.GatewayKey{}).
		Where("account_address = ?", acc.Address).
		Find(&keys).
		Error

	if err != nil {
		return nil, err
	}

	wrappedKeys := make([]*Key, len(keys))
	for i, k := range keys {
		wrappedKeys[i] = &Key{k, acc.databaseClient, acc.apisixClient}
	}

	return wrappedKeys, nil
}

func (acc *Account) GetUsage(ctx context.Context) (int64, int64, int64, int64, error) {
	var status struct {
		RuUsedTotal     int64 `gorm:"column:ru_used_total"`
		RuUsedCurrent   int64 `gorm:"column:ru_used_current"`
		APICallsTotal   int64 `gorm:"column:api_calls_total"`
		APICallsCurrent int64 `gorm:"column:api_calls_current"`
	}

	err := acc.databaseClient.WithContext(ctx).
		Model(&table.GatewayKey{}).
		Unscoped().
		Select("SUM(ru_used_total) AS ru_used_total, SUM(ru_used_current) AS ru_used_current, SUM(api_calls_total) AS api_calls_total, SUM(api_calls_current) AS api_calls_current").
		Where("account_address = ?", acc.Address).
		Find(&status).
		Error

	return status.RuUsedTotal, status.RuUsedCurrent, status.APICallsTotal, status.APICallsCurrent, err
}

func (acc *Account) GetUsageByDate(ctx context.Context, since time.Time, until time.Time) (*[]table.GatewayConsumptionLog, error) {
	var logs []table.GatewayConsumptionLog

	err := acc.databaseClient.WithContext(ctx).
		Model(&table.GatewayConsumptionLog{}).
		Joins("JOIN gateway.key").
		Where("account_address = ? AND consumption_date >= ? AND consumption_date <= ?", acc.Address, since, until).
		Select("SUM(ru_used) AS ru_used, SUM(api_calls) AS api_calls, (EXTRACT(EPOCH FROM consumption_date)*1000)::BIGINT AS consumption_date").
		Group("consumption_date").
		Order("consumption_date DESC").
		Find(&logs).
		Error

	if err != nil {
		return nil, err
	}

	return &logs, nil
}

func (acc *Account) GetBalance(ctx context.Context) (int64, error) {
	_, ruUsed, _, _, err := acc.GetUsage(ctx)
	if err != nil {
		return 0, err
	}

	return acc.RuLimit - ruUsed, nil
}

func (acc *Account) GetKey(ctx context.Context, keyID uint64) (*Key, bool, error) {
	var k table.GatewayKey

	err := acc.databaseClient.WithContext(ctx).
		Model(&table.GatewayKey{}).
		Where("account_address = ? AND id = ?", acc.Address, keyID).
		First(&k).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &Key{k, acc.databaseClient, acc.apisixClient}, true, nil
}
