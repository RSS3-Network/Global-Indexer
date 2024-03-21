package hub

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model/response"
	"go.uber.org/zap"
)

func (h *Hub) GetNodeCountSnapshots(c echo.Context) error {
	nodeSnapshots, err := h.databaseClient.FindNodeCountSnapshots(c.Request().Context())
	if err != nil {
		zap.L().Error("find node snapshots", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	var response Response

	response.Data = model.NewNodeSnapshots(nodeSnapshots)

	return c.JSON(http.StatusOK, response)
}

func (h *Hub) BatchGetNodeMinTokensToStakeSnapshots(c echo.Context) error {
	var request BatchNodeMinTokensToStakeRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	nodeMinTokensToStakeSnapshots, err := h.databaseClient.FindNodeMinTokensToStakeSnapshots(c.Request().Context(), request.NodeAddresses, request.OnlyStartAndEnd, nil)
	if err != nil {
		zap.L().Error("find node min tokens to stake snapshots", zap.Error(err))

		return response.InternalError(c, fmt.Errorf("find node min tokens to stake snapshots: %w", err))
	}

	return c.JSON(http.StatusOK, Response{
		Data: model.NewNodeMinTokensToStakeSnapshots(nodeMinTokensToStakeSnapshots),
	})
}

func (h *Hub) GetStakersCountSnapshots(c echo.Context) error {
	stakeSnapshots, err := h.databaseClient.FindStakeSnapshots(c.Request().Context())
	if err != nil {
		zap.L().Error("find stake snapshots", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	var response Response

	response.Data = model.NewStakeSnapshots(stakeSnapshots)

	return c.JSON(http.StatusOK, response)
}

type BatchNodeMinTokensToStakeRequest struct {
	NodeAddresses   []*common.Address `json:"nodeAddresses" validate:"required"`
	OnlyStartAndEnd bool              `json:"onlyStartAndEnd"`
}
