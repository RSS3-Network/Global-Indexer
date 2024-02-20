package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/model"
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

	err = k.Delete(ctx.Request().Context())
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusOK)
}

func (app *App) GenerateKey(ctx echo.Context) error {
	user := ctx.Get("user").(*model.Account)

	var req oapi.KeyInfoBody
	if err := ctx.Bind(&req); err != nil || req.Name == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	k, err := model.KeyCreate(ctx.Request().Context(), user.Address, *req.Name, app.databaseClient, app.apiSixAPIService)
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, createKeyResponse(k))
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
	user := ctx.Get("user").(*model.Account)

	keys, err := user.ListKeys(ctx.Request().Context())
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	resp := oapi.Keys{}
	for _, k := range keys {
		resp = append(resp, createKeyResponse(k))
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

	err = k.UpdateInfo(ctx.Request().Context(), *req.Name)
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

	err = k.Rotate(ctx.Request().Context())
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, createKeyResponse(k))
}

func createKeyResponse(k *model.Key) oapi.Key { // Assuming KeyType is the type of k
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
