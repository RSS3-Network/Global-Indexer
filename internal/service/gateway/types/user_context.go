package types

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/api-gateway/app"
	"github.com/naturalselectionlabs/api-gateway/app/jwt"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/constants"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/variables"
)

type UserContext struct {
	echo.Context
	User *jwtext.User
}

// SetUserCookie if *User in context is not nil, sign new token and set cookie
func (c *UserContext) SetUserCookie() error {
	var err error
	token := ""
	if c.User != nil && c.User.Address != "" {
		expires := time.Now().Add(constants.AUTH_TOKEN_DURATION)
		c.User.Expires = expires.Unix()
		token, err = app.JwtExt.SignToken(c.User)
		if err != nil {
			return err
		}

		cookie := &http.Cookie{
			Name:     app.AuthTokenCookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Expires:  expires,
			Domain:   variables.SIWEDomain,
		}

		c.SetCookie(cookie)
	}

	return nil
}

func (c *UserContext) ClearUserCookie() {

	cookie := &http.Cookie{
		Name:     app.AuthTokenCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Unix(0, 0),
		Domain:   variables.SIWEDomain,
	}

	c.SetCookie(cookie)
}
