package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/utils"
	"gorm.io/gorm"
)

func (app *App) GetPendingRequestWithdraw(ctx echo.Context) error {
	user := ctx.Get("user").(*model.Account)

	amount := float32(0)

	// Check if there's any pending withdraw requests
	var pendingWithdrawRequest table.GatewayPendingWithdrawRequest
	err := app.databaseClient.WithContext(ctx.Request().Context()).
		Model(&pendingWithdrawRequest).
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
	}

	return ctx.JSON(http.StatusOK, oapi.GetRequestWithdrawResponse{Amount: &amount})
}

func (app *App) SetPendingRequestWithdraw(ctx echo.Context, params oapi.SetPendingRequestWithdrawParams) error {
	user := ctx.Get("user").(*model.Account)

	// Check if there's any pending withdraw requests
	var pendingWithdrawRequest table.GatewayPendingWithdrawRequest
	err := app.databaseClient.WithContext(ctx.Request().Context()).
		Model(&pendingWithdrawRequest).
		Where("account_address = ?", user.Address).
		First(&pendingWithdrawRequest).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil // Not real error
			if params.Amount > 0 {
				// Create
				err = app.databaseClient.WithContext(ctx.Request().Context()).
					Model(&pendingWithdrawRequest).
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
				Model(&pendingWithdrawRequest).
				Where("account_address = ?", user.Address).
				Update("amount", float64(params.Amount)).
				Error
		} else {
			// Delete
			err = app.databaseClient.WithContext(ctx.Request().Context()).
				Model(&pendingWithdrawRequest).
				Where("account_address = ?", user.Address).
				Delete(&pendingWithdrawRequest).
				Error
		}
	}

	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusOK)
}
