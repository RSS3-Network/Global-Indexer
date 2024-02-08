package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/api-gateway/app"
	"github.com/naturalselectionlabs/api-gateway/app/model"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/utils"
	"github.com/naturalselectionlabs/api-gateway/gen/entschema"
	"github.com/naturalselectionlabs/api-gateway/gen/entschema/account"
	"github.com/naturalselectionlabs/api-gateway/gen/entschema/pendingwithdrawrequest"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
	"net/http"
)

func (*App) GetPendingRequestWithdraw(ctx echo.Context) error {
	user := ctx.Get("user").(*model.Account)
	rctx := ctx.Request().Context()

	amount := float32(0)
	// Check if there's any pending withdraw requests
	pendingWithdrawRequest, err := app.EntClient.PendingWithdrawRequest.Query().Where(
		pendingwithdrawrequest.HasAccountWith(
			account.ID(user.ID),
		),
	).First(rctx)
	if err != nil {
		if entschema.IsNotFound(err) {
			err = nil // Not real error
		}
	} else {
		amount = float32(pendingWithdrawRequest.Amount)
	}

	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	} else {
		return ctx.JSON(http.StatusOK, oapi.GetRequestWithdrawResponse{Amount: &amount})
	}
}

func (*App) SetPendingRequestWithdraw(ctx echo.Context, params oapi.SetPendingRequestWithdrawParams) error {
	user := ctx.Get("user").(*model.Account)
	rctx := ctx.Request().Context()

	// Check if there's any pending withdraw requests
	pendingWithdrawRequest, err := app.EntClient.PendingWithdrawRequest.Query().Where(
		pendingwithdrawrequest.HasAccountWith(
			account.ID(user.ID),
		),
	).First(rctx)
	if err != nil {
		if entschema.IsNotFound(err) {
			err = nil // Not real error
			if params.Amount > 0 {
				// Create
				err = app.EntClient.PendingWithdrawRequest.Create().
					SetAccountID(user.ID).
					SetAmount(float64(params.Amount)).
					Exec(rctx)
			}
		}
	} else {
		// Found record with no error
		if params.Amount > 0 {
			// Update
			err = app.EntClient.PendingWithdrawRequest.UpdateOneID(pendingWithdrawRequest.ID).
				SetAmount(float64(params.Amount)).
				Exec(rctx)
		} else {
			// Delete
			err = app.EntClient.PendingWithdrawRequest.DeleteOneID(pendingWithdrawRequest.ID).Exec(rctx)
		}
	}

	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	} else {
		return ctx.NoContent(http.StatusOK)
	}

}
