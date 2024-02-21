package hub

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"go.uber.org/zap"
)

func (h *Hub) GetNodeSnapshots(c echo.Context) error {
	nodeSnapshots, err := h.databaseClient.FindNodeSnapshots(c.Request().Context())
	if err != nil {
		zap.L().Error("find node snapshots", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	var response Response

	response.Data = model.NewNodeSnapshots(nodeSnapshots)

	return c.JSON(http.StatusOK, response)
}

func (h *Hub) GetStakeSnapshots(c echo.Context) error {
	stakeSnapshots, err := h.databaseClient.FindStakeSnapshots(c.Request().Context())
	if err != nil {
		zap.L().Error("find stake snapshots", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	var response Response

	response.Data = model.NewStakeSnapshots(stakeSnapshots)

	return c.JSON(http.StatusOK, response)
}
