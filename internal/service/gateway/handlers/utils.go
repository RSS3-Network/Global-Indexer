package handlers

import (
	"context"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/constants"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/jwt"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/middlewares"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
)

func (app *App) getCtx(ctx echo.Context) (context.Context, *jwt.User) {
	return ctx.Request().Context(), middlewares.ParseUserWithToken(ctx, app.jwtClient)
}

func (app *App) getKey(ctx echo.Context, keyID int) (*table.GatewayKey, error) {
	user := ctx.Get("user").(*table.GatewayAccount)

	var k table.GatewayKey
	err := app.databaseClient.WithContext(ctx.Request().Context()).
		Model(&table.GatewayKey{}).
		Where("account_address = ? AND id = ?", user.Address, keyID).
		First(&k).
		Error
	if err != nil {
		return nil, err
	}

	return &k, nil
}

func parseDates(since *oapi.Since, until *oapi.Until) (time.Time, time.Time) {
	var startFrom, untilTo time.Time
	nowTime := time.Now()

	if since != nil {
		startFrom = time.UnixMilli(*since)
	} else {
		startFrom = nowTime.Add(-constants.DEFAULT_HISTORY_SINCE)
	}

	if until != nil {
		untilTo = time.UnixMilli(*until)
	} else {
		untilTo = nowTime
	}

	if untilTo.Before(startFrom) {
		// Swap
		startFrom, untilTo = untilTo, startFrom
	}
	return startFrom, untilTo
}

func parseLimitPage(limit *oapi.Limit, page *oapi.Page) (int, int) {
	var (
		l = constants.DEFAULT_PAGINATION_LIMIT
		p = 1
	)

	if limit != nil {
		l = int(*limit)
	}

	if page != nil && *page >= 1 {
		p = int(*page)
	}

	return l, p
}
