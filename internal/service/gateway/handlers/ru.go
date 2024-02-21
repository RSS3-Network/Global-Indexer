package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/utils"
)

func (app *App) GetRUStatus(ctx echo.Context) error {
	user := ctx.Get("user").(*model.Account)

	ruUsedTotal, ruUsedCurrent, apiCallsTotal, apiCallsCurrent, err := user.GetUsage(ctx.Request().Context())
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	resp := oapi.RUStatus{
		RuLimit:         &user.RuLimit,
		RuUsedTotal:     &ruUsedTotal,
		RuUsedCurrent:   &ruUsedCurrent,
		ApiCallsTotal:   &apiCallsTotal,
		ApiCallsCurrent: &apiCallsCurrent,
	}

	return ctx.JSON(http.StatusOK, resp)
}
