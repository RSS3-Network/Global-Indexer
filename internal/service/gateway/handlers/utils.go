package handlers

import (
	"context"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/constants"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/jwt"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/middlewares"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/model"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
)

func (app *App) getCtx(ctx echo.Context) (context.Context, *jwt.User) {
	return ctx.Request().Context(), middlewares.ParseUserWithToken(ctx, app.jwtClient)
}

func (app *App) getKey(ctx echo.Context, keyID string) (*model.Key, bool, error) {
	user := ctx.Get("user").(*model.Account)

	keyIDUint64, err := strconv.ParseUint(keyID, 10, 64)
	if err != nil {
		return nil, false, err
	}

	return user.GetKey(ctx.Request().Context(), keyIDUint64)
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
