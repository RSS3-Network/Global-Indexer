package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/naturalselectionlabs/api-gateway/app"
	apisixHTTPAPI "github.com/naturalselectionlabs/api-gateway/app/apisix/httpapi"
	apisixKafkaLog "github.com/naturalselectionlabs/api-gateway/app/apisix/kafkalog"
	"github.com/naturalselectionlabs/api-gateway/app/model"
	"github.com/naturalselectionlabs/api-gateway/app/reverseproxy/rules"
	"github.com/naturalselectionlabs/api-gateway/gen/entschema/account"
)

func ProcessAccessLog(accessLog apisixKafkaLog.AccessLog) {

	rctx := context.Background()

	// Check billing eligibility
	if accessLog.Consumer == nil {
		return
	}

	// Find user
	keyId, err := apisixHTTPAPI.RecoverKeyIDFromConsumerUsername(*accessLog.Consumer)
	if err != nil {
		log.Printf("Failed to recover key id with error: %v", err)
		return
	}
	key, err := model.KeyGetById(rctx, keyId, false) // Deleted key could also be used for pending bills
	if err != nil {
		log.Printf("Failed to get key by id with error: %v", err)
		return
	}

	// Check RU remain
	user, err := key.GetAccount(rctx)
	if err != nil {
		// Failed to get account
		log.Printf("Faield to get account with error: %v", err)
		return
	}

	if accessLog.Status != http.StatusOK || user.IsPaused {
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
		// Failed to consumer RU
		log.Printf("Faield to consumer RU with error: %v", err)
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
		if !user.IsPaused {
			err = apisixHTTPAPI.PauseConsumerGroup(user.Address)
			if err != nil {
				log.Printf("Failed to pause account with error: %v", err)
			} else {
				err = app.EntClient.Account.Update().SetIsPaused(true).Where(
					account.ID(user.ID),
				).Exec(rctx)
				if err != nil {
					log.Printf("Failed to save paused account into db with error: %v", err)
				}
			}
		}
	}

}
