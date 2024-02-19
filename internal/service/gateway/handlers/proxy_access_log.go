package handlers

import (
	"context"
	"fmt"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	apisixKafkaLog "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/apisix/kafkalog"
	rules "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/ru_rules"
	"gorm.io/gorm"
	"log"
	"net/http"

	"strings"
)

func (app *App) ProcessAccessLog(accessLog apisixKafkaLog.AccessLog) {

	rctx := context.Background()

	// Check billing eligibility
	if accessLog.Consumer == nil {
		return
	}

	// Find user
	keyID, err := app.apiSixAPIService.RecoverKeyIDFromConsumerUsername(*accessLog.Consumer)
	if err != nil {
		log.Printf("Failed to recover key id with error: %v", err)
		return
	}
	var key table.GatewayKey
	err = app.databaseClient.WithContext(rctx).
		Model(&table.GatewayKey{}).
		Unscoped(). // Deleted key could also be used for pending bills
		Where("id = ?", keyID).
		First(&key).
		Error
	if err != nil {
		log.Printf("Failed to get key by id with error: %v", err)
		return
	}

	// Add API calls count
	err = app.databaseClient.WithContext(rctx).
		Model(&table.GatewayKey{}).
		Unscoped().
		Where("id = ?", keyID).
		Update("api_calls_current", gorm.Expr("api_calls_current + ?", 1)).
		Error
	if err != nil {
		// Failed to consumer RU
		log.Printf("Faield to increase API call count with error: %v", err)
	}

	if accessLog.Status != http.StatusOK || key.Account.IsPaused {
		// Request failed or is in free tier, only increase API call count
		return
	}

	// Consumer RU
	pathSplits := strings.Split(accessLog.URI, "/")
	ruCalculator, ok := rules.Prefix2RUCalculator[pathSplits[1]]
	if !ok {
		// Invalid path
		log.Printf("No matching route prefix")
		return
	}
	ru := ruCalculator(fmt.Sprintf("/%s", strings.Join(pathSplits[2:], "/")))
	err = app.databaseClient.WithContext(rctx).
		Model(&table.GatewayKey{}).
		Unscoped().
		Where("id = ?", keyID).
		Update("ru_used_current", gorm.Expr("ru_used_current + ?", ru)).
		Error
	if err != nil {
		// Failed to consume RU
		log.Printf("Faield to consume RU with error: %v", err)
		return
	}

	var (
		accountTotalConsumedRU int64
	)
	err = app.databaseClient.WithContext(rctx).
		Model(&table.GatewayKey{}).
		Unscoped().
		Select("SUM(ru_used_current) AS accountTotalConsumedRU").
		Where("account_address = ?", key.Account.Address).
		Scan(&accountTotalConsumedRU).
		Error
	if err != nil {
		// Failed to get remain RU
		log.Printf("Faield to get account remain RU with error: %v", err)
		return
	}

	if key.Account.RuLimit-accountTotalConsumedRU < 0 {
		log.Printf("Insufficient remain RU, pause account")
		// Pause user account
		if !key.Account.IsPaused {
			err = app.apiSixAPIService.PauseConsumerGroup(key.Account.Address.Hex())
			if err != nil {
				log.Printf("Failed to pause account with error: %v", err)
			} else {
				err = app.databaseClient.WithContext(rctx).
					Model(&table.GatewayAccount{}).
					Unscoped().
					Where("address = ?", key.Account.Address).
					Update("is_paused", true).
					Error
				if err != nil {
					log.Printf("Failed to save paused account into db with error: %v", err)
				}
			}
		}
	}

}
