package dsl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"go.uber.org/zap"
)

func (d *DSL) GetRSSHub(c echo.Context) error {
	path := c.Param("*")
	query := c.Request().URL.RawQuery

	data, err := d.distributor.DistributeRSSHubData(c.Request().Context(), path, query)

	if err != nil {
		zap.L().Error("distribute rss hub data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, data)
}
