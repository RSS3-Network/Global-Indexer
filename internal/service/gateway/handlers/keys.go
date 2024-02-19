package handlers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/utils"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func (app *App) DeleteKey(ctx echo.Context, keyID int) error {
	k, err := app.getKey(ctx, keyID)
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusNotFound)
	}

	err = app.databaseClient.WithContext(ctx.Request().Context()).
		Delete(&k).
		Error
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusOK)
}

func (app *App) GenerateKey(ctx echo.Context) error {
	user := ctx.Get("user").(*table.GatewayAccount)

	var req oapi.KeyInfoBody
	if err := ctx.Bind(&req); err != nil || req.Name == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	k := table.GatewayKey{
		Key:            uuid.New(),
		Name:           *req.Name,
		AccountAddress: user.Address,
	}
	err := app.databaseClient.WithContext(ctx.Request().Context()).
		Create(&k).
		Error
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, createKeyResponse(&k))
}

func (app *App) GetKey(ctx echo.Context, keyID int) error {
	k, err := app.getKey(ctx, keyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendJSONError(ctx, http.StatusNotFound)
		} else {
			log.Print(err)
			return utils.SendJSONError(ctx, http.StatusInternalServerError)
		}
	}

	return ctx.JSON(http.StatusOK, createKeyResponse(k))
}

func (app *App) GetKeys(ctx echo.Context) error {
	user := ctx.Get("user").(*table.GatewayAccount)

	var keys []table.GatewayKey

	err := app.databaseClient.WithContext(ctx.Request().Context()).
		Model(&table.GatewayKey{}).
		Where("account_address = ?", user.Address).
		Error

	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	resp := oapi.Keys{}
	for _, k := range keys {
		resp = append(resp, createKeyResponse(&k))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (app *App) UpdateKeyInfo(ctx echo.Context, keyID int) error {
	var req oapi.KeyInfoBody
	if err := ctx.Bind(&req); err != nil || req.Name == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	k, err := app.getKey(ctx, keyID)
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusUnauthorized)
	}

	// Update fields
	k.Name = *req.Name

	err = app.databaseClient.WithContext(ctx.Request().Context()).
		Model(&table.GatewayKey{}).
		Where("id = ?", keyID).
		Update("name", k.Name).
		Error

	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, createKeyResponse(k))
}

func (app *App) RotateKey(ctx echo.Context, keyID int) error {
	k, err := app.getKey(ctx, keyID)
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	k.Key = uuid.New()

	err = app.databaseClient.WithContext(ctx.Request().Context()).
		Model(&table.GatewayKey{}).
		Where("id = ?", keyID).
		Update("key", k.Key).
		Error

	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, createKeyResponse(k))
}

func createKeyResponse(k *table.GatewayKey) oapi.Key { // Assuming KeyType is the type of k
	return oapi.Key{
		Id:              lo.ToPtr(int(k.ID)),
		Key:             lo.ToPtr(k.Key.String()),
		Name:            &k.Name,
		ApiCallsTotal:   &k.ApiCallsTotal,
		ApiCallsCurrent: &k.ApiCallsCurrent,
		RuUsedTotal:     &k.RuUsedTotal,
		RuUsedCurrent:   &k.RuUsedCurrent,
	}
}
