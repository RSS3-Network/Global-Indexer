package dsl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/hub/model/errorx"
)

func (d *DSL) GetRSSHub(c echo.Context) error {
	path := c.Param("*")
	query := c.Request().URL.RawQuery

	data, err := d.Distributor.RouterRSSHubData(c.Request().Context(), path, query)

	if err != nil {
		return errorx.InternalError(c, err)
	}

	return c.JSONBlob(http.StatusOK, data)
}
