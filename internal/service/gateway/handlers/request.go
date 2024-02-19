package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/utils"
	"gorm.io/gorm"
	"net/http"
)

func (app *App) GetPendingRequestWithdraw(ctx echo.Context) error {
	user := ctx.Get("user").(*table.GatewayAccount)

	amount := float32(0)

	// Check if there's any pending withdraw requests
	var pendingWithdrawRequest table.GatewayPendingWithdrawRequest
	err := app.databaseClient.WithContext(ctx.Request().Context()).
		Where("account_address = ?", user.Address).
		First(&pendingWithdrawRequest).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

func (app *App) SetPendingRequestWithdraw(ctx echo.Context, params oapi.SetPendingRequestWithdrawParams) error {
	user := ctx.Get("user").(*table.GatewayAccount)

	// Check if there's any pending withdraw requests
	var pendingWithdrawRequest table.GatewayPendingWithdrawRequest
	err := app.databaseClient.WithContext(ctx.Request().Context()).
		Where("account_address = ?", user.Address).
		First(&pendingWithdrawRequest).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil // Not real error
			if params.Amount > 0 {
				// Create
				err = app.databaseClient.WithContext(ctx.Request().Context()).
					Create(&table.GatewayPendingWithdrawRequest{
						AccountAddress: user.Address,
						Amount:         float64(params.Amount),
					}).
					Error
			}
		}
	} else {
		// Found record with no error
		if params.Amount > 0 {
			// Update
			err = app.databaseClient.WithContext(ctx.Request().Context()).
				Where("account_address = ?", user.Address).
				Update("amount", float64(params.Amount)).
				Error
		} else {
			// Delete
			err = app.databaseClient.WithContext(ctx.Request().Context()).
				Where("account_address = ?", user.Address).
				Delete(&pendingWithdrawRequest).
				Error
		}
	}

	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	} else {
		return ctx.NoContent(http.StatusOK)
	}

}
