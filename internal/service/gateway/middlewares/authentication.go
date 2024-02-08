package middlewares

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/api-gateway/app"
	jwtext "github.com/naturalselectionlabs/api-gateway/app/jwt"
	"github.com/naturalselectionlabs/api-gateway/app/model"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/types"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/utils"
	"net/http"
	"regexp"
)

func authenticateUser(rctx context.Context, uctx *types.UserContext) (model.Account, error) {
	if uctx.User == nil {
		return model.Account{}, nil
	}
	return model.AccountGetByAddress(rctx, uctx.User.Address)
}

func ParseUserWithToken(c echo.Context) *jwtext.User {
	authToken, err := c.Request().Cookie(app.AuthTokenCookieName)
	if err != nil || authToken == nil {
		return nil
	}

	user, _ := app.JwtExt.ParseUser(authToken.Value)
	return user
}

var (
	SkipMiddlewarePaths = regexp.MustCompile("^/(users/|health)")
)

func UserAuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// this is a hack to workaround codegen and echo router group issue
		// see: https://github.com/labstack/echo/issues/1737
		// otherwise we need to manually group the routes
		// see: https://github.com/labstack/echo/issues/1737#issuecomment-906721802
		if SkipMiddlewarePaths.MatchString(c.Path()) {
			return next(c)
		}

		user := ParseUserWithToken(c)
		if user == nil {
			return utils.SendJSONError(c, http.StatusUnauthorized)
		}

		// Authenticate user
		uctx := &types.UserContext{Context: c, User: user}
		account, err := authenticateUser(c.Request().Context(), uctx)
		if err != nil || account == (model.Account{}) {
			return utils.SendJSONError(c, http.StatusUnauthorized)
		}

		// Set user in context
		c.Set("user", &account)

		// Continue with the pipeline
		return next(c)
	}
}
