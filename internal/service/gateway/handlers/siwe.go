package handlers

import (
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/constants"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/jwt"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/utils"
	"github.com/samber/lo"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
)

// SIWEGetNonce implements oapi.ServerInterface
func (app *App) SIWEGetNonce(ctx echo.Context) error {
	// Get nonce
	nonce, err := app.siweClient.GetNonce(ctx.Request().Context())
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	// Return
	return ctx.String(http.StatusOK, nonce)
}

// SIWEGetSession implements oapi.ServerInterface
func (app *App) SIWEGetSession(ctx echo.Context) error {
	res := oapi.SIWESessionResponse{}

	_, user := app.getCtx(ctx)
	if user != nil {
		res.Address = lo.ToPtr(user.Address.Hex())
		res.ChainId = &user.ChainId
	}

	// User has logged in
	return ctx.JSON(http.StatusOK, res)
}

// SIWEVerify implements oapi.ServerInterface
func (app *App) SIWEVerify(ctx echo.Context) error {
	var req oapi.SIWEVerifyBody
	if err := ctx.Bind(&req); err != nil || req.Message == nil || req.Signature == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	address, chainId, err := app.siweClient.ValidateSIWESignature(ctx.Request().Context(), *req.Message, *req.Signature)
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusUnauthorized)
	}

	// get or create account
	var acc table.GatewayAccount
	err = app.databaseClient.WithContext(ctx.Request().Context()).
		FirstOrCreate(&acc, "address = ?", *address).
		Error
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	// set User with expiration
	expires := time.Now().Add(constants.AUTH_TOKEN_DURATION)
	token, err := app.jwtClient.SignToken(&jwt.User{
		Address: acc.Address,
		ChainId: chainId,
		Expires: expires.Unix(),
	})
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	cookie := &http.Cookie{
		Name:     constants.AuthTokenCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  expires,
		Domain:   app.siweClient.Domain(),
	}

	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}

// SIWELogout implements oapi.ServerInterface
func (app *App) SIWELogout(ctx echo.Context) error {
	cookie := &http.Cookie{
		Name:     constants.AuthTokenCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Unix(0, 0),
		Domain:   app.siweClient.Domain(),
	}

	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}
