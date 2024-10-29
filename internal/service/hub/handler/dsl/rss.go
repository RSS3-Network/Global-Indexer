package dsl

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"go.uber.org/zap"
)

func (d *DSL) GetRSSHub(c echo.Context) error {
	path := c.Param("*")
	query := c.Request().URL.RawQuery

	requestCounter.WithLabelValues("GetRSSHub").Inc()

	data, err := d.distributor.DistributeRSSHubData(c.Request().Context(), path, query)

	if err != nil {
		if errors.Is(err, errorx.ErrNoNodesAvailable) {
			return errorx.BadRequestError(c, err)
		}

		zap.L().Error("distribute rss hub data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, data)
}
