package model

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	apisixHTTPAPI "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/apisix/httpapi"
	"gorm.io/gorm"
	"log"
)

type Key struct {
	table.GatewayKey

	databaseClient   *gorm.DB
	apiSixAPIService *apisixHTTPAPI.HTTPAPIService
}

func KeyCreate(ctx context.Context, accountAddress common.Address, keyName string, databaseClient *gorm.DB, apiSixAPIService *apisixHTTPAPI.HTTPAPIService) (*Key, error) {
	keyUUID := uuid.New()
	k := table.GatewayKey{
		Key:            keyUUID,
		Name:           keyName,
		AccountAddress: accountAddress,
	}
	err := databaseClient.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// DB
		err := tx.
			Create(&k).
			Error
		if err != nil {
			return err
		}
		// APISix
		err = apiSixAPIService.NewConsumer(k.ID, keyUUID.String(), accountAddress.Hex())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Key{k, databaseClient, apiSixAPIService}, nil
}

func KeyGetByID(ctx context.Context, KeyID uint, activeOnly bool, databaseClient *gorm.DB, apiSixAPIService *apisixHTTPAPI.HTTPAPIService) (*Key, bool, error) {
	queryBase := databaseClient.WithContext(ctx).Model(&table.GatewayKey{})
	if activeOnly {
		queryBase = queryBase.Unscoped()
	}
	var k table.GatewayKey
	err := queryBase.Where("id = ?", KeyID).First(&k).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	} else {
		return &Key{k, databaseClient, apiSixAPIService}, true, nil
	}
}

func (k *Key) ConsumeRu(ctx context.Context, ru int64) error {

	// Add API calls count
	err := k.databaseClient.WithContext(ctx).
		Model(&table.GatewayKey{}).
		Unscoped().
		Where("id = ?", k.ID).
		Updates(map[string]interface{}{
			"api_calls_current": gorm.Expr("api_calls_current + ?", 1),
			"ru_used_current":   gorm.Expr("ru_used_current + ?", ru),
		}).
		Error
	if err != nil {
		// Failed to consumer RU
		log.Printf("Faield to increase API call count with error: %v", err)
		return err
	}

	return nil
}

func (k *Key) GetAccount(_ context.Context) (*Account, error) {
	return &Account{k.Account, k.databaseClient, k.apiSixAPIService}, nil
}

func (k *Key) Delete(ctx context.Context) error {

	err := k.apiSixAPIService.DeleteConsumer(k.ID)
	if err != nil {
		return err
	}

	err = k.databaseClient.WithContext(ctx).
		Delete(&k).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (k *Key) UpdateInfo(ctx context.Context, name string) error {
	k.Name = name
	err := k.databaseClient.WithContext(ctx).
		Model(&table.GatewayKey{}).
		Where("id = ?", k.ID).
		Update("name", name).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (k *Key) Rotate(ctx context.Context) error {

	// Replace old consumer
	oldConsumer, err := k.apiSixAPIService.CheckConsumer(k.ID)
	if err != nil {
		return err
	}

	k.Key = uuid.New()

	err = k.databaseClient.WithContext(ctx).
		Model(&table.GatewayKey{}).
		Where("id = ?", k.ID).
		Update("key", k.Key).
		Error
	if err != nil {
		return err
	}

	// Update consumer
	err = k.apiSixAPIService.NewConsumer(k.ID, k.Key.String(), oldConsumer.Value.GroupID)
	if err != nil {
		return err
	}

	return nil
}
