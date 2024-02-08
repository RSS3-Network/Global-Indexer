package handlers

import (
	"context"
	"errors"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/constants"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
)

func checkNextEpochRun(ctx context.Context, rc *redis.Client) bool {
	nextEpochRun, err := rc.
		Get(ctx, constants.NEXT_EPOCH_RUN_BEFORE).
		Result()
	if errors.Is(err, redis.Nil) {
		// No such key -> not ready, but should be OK
		return true
	} else if err != nil {
		log.Printf("Failed to check next epoch run with error: %v", err)
		return false
	} else {
		// No error
		nextRunTime, err := time.Parse(time.RFC3339, nextEpochRun)
		if err != nil {
			// Invalid time -> corrupted result, but should be OK
			log.Printf("Invalid time: %s", nextEpochRun)
			return true
		}

		return time.Now().Before(nextRunTime)
	}
}

// HealthCheck implements oapi.ServerInterface
func (a *App) HealthCheck(ctx echo.Context, params oapi.HealthCheckParams) error {
	if params.Type == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	switch *params.Type {
	case "liveness":
		if checkNextEpochRun(ctx.Request().Context(), a.redisClient) {
			return ctx.NoContent(http.StatusOK)
		} else {
			return ctx.NoContent(http.StatusInternalServerError)
		}
	case "readiness":
		return ctx.NoContent(http.StatusOK)
	default:
		return ctx.NoContent(http.StatusBadRequest)
	}
}
