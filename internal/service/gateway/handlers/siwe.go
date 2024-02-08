package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	jwtext "github.com/naturalselectionlabs/api-gateway/app/jwt"
	"github.com/naturalselectionlabs/api-gateway/app/model"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/siwe"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/utils"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
)

// SIWEGetNonce implements oapi.ServerInterface
func (*App) SIWEGetNonce(ctx echo.Context) error {
	rctx, uctx := getCtx(ctx)
	uctx.User = nil // clear User

	// Get nonce
	nonce, err := siwe.GetNonce(rctx)
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	// Return
	return ctx.String(http.StatusOK, nonce)
}

// SIWEGetSession implements oapi.ServerInterface
func (*App) SIWEGetSession(ctx echo.Context) error {
	res := oapi.SIWESessionResponse{}

	_, uctx := getCtx(ctx)
	if uctx.User != nil {
		res.Address = &uctx.User.Address
		res.ChainId = &uctx.User.ChainId
	}

	// User has logged in
	return ctx.JSON(http.StatusOK, res)
}

// SIWEVerify implements oapi.ServerInterface
func (*App) SIWEVerify(ctx echo.Context) error {
	rctx, uctx := getCtx(ctx)
	uctx.User = nil // clear User

	var req oapi.SIWEVerifyBody
	if err := ctx.Bind(&req); err != nil || req.Message == nil || req.Signature == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	address, chainId, err := siwe.ValidateSIWESignature(rctx, *req.Message, *req.Signature)
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusUnauthorized)
	}

	// get or create account
	acc, err := model.AccountGetOrCreate(rctx, address)
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	// set User with expiration in 10 days
	uctx.User = &jwtext.User{
		Address: acc.Address,
		ChainId: chainId,
	}
	if err := uctx.SetUserCookie(); err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusOK)
}

// SIWELogout implements oapi.ServerInterface
func (*App) SIWELogout(ctx echo.Context) error {
	_, uctx := getCtx(ctx)

	uctx.ClearUserCookie()

	return ctx.NoContent(http.StatusOK)
}
