package hub

import (
	"net/url"

	"github.com/labstack/echo/v4"
)

func baseURL(c echo.Context) url.URL {
	return url.URL{
		Scheme: c.Scheme(),
		Host:   c.Request().Host,
	}
}
