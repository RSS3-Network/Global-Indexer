package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/accesslog"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/model"
	rules "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/ru_rules"
)

func (app *App) ProcessAccessLog(accessLog accesslog.AccessLog) {
	rctx := context.Background()

	// Check billing eligibility
	if accessLog.Consumer == nil {
		return
	}

	// Find user
	keyID, err := app.apisixClient.RecoverKeyIDFromConsumerUsername(*accessLog.Consumer)

	if err != nil {
		log.Printf("Failed to recover key id with error: %v", err)
		return
	}

	key, _, err := model.KeyGetByID(rctx, keyID, false, app.databaseClient, app.apisixClient) // Deleted key could also be used for pending bills

	if err != nil {
		log.Printf("Failed to get key by id with error: %v", err)

		return
	}

	user, err := key.GetAccount(rctx)

	if err != nil {
		// Failed to get account
		log.Printf("Faield to get account with error: %v", err)

		return
	}

	if accessLog.Status != http.StatusOK || key.Account.IsPaused {
		err = key.ConsumeRu(rctx, 0) // Request failed or is in free tier, only increase API call count
		if err != nil {
			// Failed to consumer RU
			log.Printf("Faield to increase API call count with error: %v", err)
		}

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
	err = key.ConsumeRu(rctx, ru)

	if err != nil {
		// Failed to consume RU
		log.Printf("Faield to consume RU with error: %v", err)

		return
	}

	ruRemain, err := user.GetBalance(rctx)

	if err != nil {
		// Failed to get remain RU
		log.Printf("Faield to get account remain RU with error: %v", err)

		return
	}

	if ruRemain < 0 {
		log.Printf("Insufficient remain RU, pause account")
		// Pause user account
		if !key.Account.IsPaused {
			err = app.apisixClient.PauseConsumerGroup(rctx, key.Account.Address.Hex())
			if err != nil {
				log.Printf("Failed to pause account with error: %v", err)
			} else {
				err = app.databaseClient.WithContext(rctx).
					Model(&table.GatewayAccount{}).
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
