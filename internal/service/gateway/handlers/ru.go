package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/api-gateway/app/model"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/utils"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
)

func (*App) GetRUStatus(ctx echo.Context) error {
	rctx, _ := getCtx(ctx)

	user := ctx.Get("user").(*model.Account)

	ruUsedTotal, ruUsedCurrent, apiCallsTotal, apiCallsCurrent, err := user.GetUsage(rctx)
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	resp := oapi.RUStatus{
		RuLimit:         to.Int64_Int64Ptr(user.RuLimit),
		RuUsedTotal:     to.Int64_Int64Ptr(ruUsedTotal),
		RuUsedCurrent:   to.Int64_Int64Ptr(ruUsedCurrent),
		ApiCallsTotal:   to.Int64_Int64Ptr(apiCallsTotal),
		ApiCallsCurrent: to.Int64_Int64Ptr(apiCallsCurrent),
	}

	return ctx.JSON(http.StatusOK, resp)
}
