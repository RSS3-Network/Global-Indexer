package hub

import (
	"net/url"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Data   any    `json:"data"`
	Cursor string `json:"cursor,omitempty"`
}

func baseURL(c echo.Context) url.URL {
	return url.URL{
		Scheme: c.Scheme(),
		Host:   c.Request().Host,
	}
}
