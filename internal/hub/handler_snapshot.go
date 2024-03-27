package hub

import (
	"fmt"
	"net/http"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model/response"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
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

func (h *Hub) GetStakerCountSnapshots(c echo.Context) error {
	stakeSnapshots, err := h.databaseClient.FindStakerCountSnapshots(c.Request().Context())
	if err != nil {
		zap.L().Error("find staker_count snapshots", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	var response Response

	response.Data = model.NewStakeSnapshots(stakeSnapshots)

	return c.JSON(http.StatusOK, response)
}

func (h *Hub) GetStakerProfitsSnapshots(c echo.Context) error {
	var request GetStakerProfitSnapshotsRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return response.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	query := schema.StakerProfitSnapshotsQuery{
		OwnerAddress: lo.ToPtr(request.OwnerAddress),
		Limit:        lo.ToPtr(request.Limit),
		Cursor:       request.Cursor,
	}

	stakerProfitSnapshots, err := h.databaseClient.FindStakerProfitSnapshots(c.Request().Context(), query)
	if err != nil {
		zap.L().Error("find staker profit snapshots", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	var cursor string

	if len(stakerProfitSnapshots) > 0 && len(stakerProfitSnapshots) == request.Limit {
		last, _ := lo.Last(stakerProfitSnapshots)
		cursor = fmt.Sprintf("%d", last.ID)
	}

	return c.JSON(http.StatusOK, Response{
		Data:   stakerProfitSnapshots,
		Cursor: cursor,
	})
}

type BatchNodeMinTokensToStakeRequest struct {
	NodeAddresses   []*common.Address `json:"nodeAddresses" validate:"required"`
	OnlyStartAndEnd bool              `json:"onlyStartAndEnd"`
}

type GetStakerProfitSnapshotsRequest struct {
	OwnerAddress common.Address `query:"ownerAddress" validate:"required"`
	Limit        int            `query:"limit" validate:"min=1,max=100" default:"100"`
	Cursor       *string        `query:"cursor"`
}
