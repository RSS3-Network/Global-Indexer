package handlers

import (
	"github.com/naturalselectionlabs/api-gateway/gen/entschema"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/api-gateway/app/model"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/utils"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
)

func (*App) DeleteKey(ctx echo.Context, key int) error {
	k, err := getKey(ctx, key)
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusNotFound)
	}
	if err := k.Delete(ctx.Request().Context()); err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusOK)
}

func (*App) GenerateKey(ctx echo.Context) error {
	user := ctx.Get("user").(*model.Account)

	var req oapi.KeyInfoBody
	if err := ctx.Bind(&req); err != nil || req.Name == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	k, err := model.KeyCreate(ctx.Request().Context(), user.ID, user.Address, *req.Name)
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, createKeyResponse(&k))
}

func (*App) GetKey(ctx echo.Context, keyId int) error {
	k, err := getKey(ctx, keyId)
	if err != nil {
		if entschema.IsNotFound(err) {
			return utils.SendJSONError(ctx, http.StatusNotFound)
		} else {
			log.Print(err)
			return utils.SendJSONError(ctx, http.StatusInternalServerError)
		}
	}

	return ctx.JSON(http.StatusOK, createKeyResponse(k))
}

func (*App) GetKeys(ctx echo.Context) error {
	rctx := ctx.Request().Context()

	user := ctx.Get("user").(*model.Account)

	keys, err := user.ListKeys(rctx)
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

func (*App) UpdateKeyInfo(ctx echo.Context, keyId int) error {
	var req oapi.KeyInfoBody
	if err := ctx.Bind(&req); err != nil || req.Name == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	k, err := getKey(ctx, keyId)
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusUnauthorized)
	}

	if err = k.UpdateInfo(ctx.Request().Context(), *req.Name); err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, createKeyResponse(k))
}

func (*App) RotateKey(ctx echo.Context, keyId int) error {
	k, err := getKey(ctx, keyId)
	if err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	if err = k.Rotate(ctx.Request().Context()); err != nil {
		log.Print(err)
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, createKeyResponse(k))
}

func createKeyResponse(k *model.Key) oapi.Key { // Assuming KeyType is the type of k
	return oapi.Key{
		Id:              &k.ID,
		Key:             to.String_StringPtr(k.Key.Key.String()),
		Name:            to.String_StringPtr(k.Name),
		ApiCallsTotal:   to.Int64_Int64Ptr(k.APICallsTotal),
		ApiCallsCurrent: to.Int64_Int64Ptr(k.APICallsCurrent),
		RuUsedTotal:     to.Int64_Int64Ptr(k.RuUsedTotal),
		RuUsedCurrent:   to.Int64_Int64Ptr(k.RuUsedCurrent),
	}
}
