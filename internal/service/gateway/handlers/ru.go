package handlers

import (
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
)

type ruStatus struct {
	RuUsedTotal     int64
	RuUsedCurrent   int64
	ApiCallsTotal   int64
	ApiCallsCurrent int64
}

func (app *App) GetRUStatus(ctx echo.Context) error {
	user := ctx.Get("user").(*table.GatewayAccount)

	var status ruStatus

	err := app.databaseClient.WithContext(ctx.Request().Context()).
		Model(&table.GatewayKey{}).
		Unscoped().
		Select("SUM(ru_used_total) AS ruUsedTotal, SUM(ru_used_current) AS ruUsedCurrent, SUM(api_calls_total) AS apiCallsTotal, SUM(api_calls_current) AS apiCallsCurrent").
		Where("account_address = ?", user.Address).
		Scan(&status).
		Error
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	resp := oapi.RUStatus{
		RuLimit:         &user.RuLimit,
		RuUsedTotal:     &status.RuUsedTotal,
		RuUsedCurrent:   &status.RuUsedCurrent,
		ApiCallsTotal:   &status.ApiCallsTotal,
		ApiCallsCurrent: &status.ApiCallsCurrent,
	}

	return ctx.JSON(http.StatusOK, resp)
}
