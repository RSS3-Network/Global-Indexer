package hub

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model/response"
)

func (h *Hub) GetEpochsHandler(c echo.Context) error {
	var request GetEpochsRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bad request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return response.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	epochs, err := h.databaseClient.FindEpochs(c.Request().Context(), request.Limit, request.Cursor)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return response.InternalError(c, fmt.Errorf("get failed: %w", err))
	}

	data := model.NewEpochs(epochs)

	var cursor string
	if len(data) > 0 && len(data) == request.Limit {
		cursor = fmt.Sprintf("%d", data[len(data)-1].ID)
	}

	return c.JSON(http.StatusOK, Response{
		Data:   data,
		Cursor: cursor,
	})
}

func (h *Hub) GetEpochHandler(c echo.Context) error {
	var request GetEpochRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bad request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return response.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	epoch, err := h.databaseClient.FindEpochTransactions(c.Request().Context(), request.ID, request.ItemsLimit, request.Cursor)
	if errors.Is(err, database.ErrorRowNotFound) || len(epoch) == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	if err != nil {
		return response.InternalError(c, fmt.Errorf("get failed: %w", err))
	}

	return c.JSON(http.StatusOK, Response{
		Data: model.NewEpoch(request.ID, epoch),
	})
}

func (h *Hub) GetEpochDistributionHandler(c echo.Context) error {
	var request GetEpochDistributionRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bad request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return response.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	epoch, err := h.databaseClient.FindEpochTransaction(c.Request().Context(), request.TransactionHash, request.ItemsLimit, request.Cursor)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return response.InternalError(c, fmt.Errorf("get failed: %w", err))
	}

	var cursor string
	if len(epoch.RewardItems) > 0 && len(epoch.RewardItems) == request.ItemsLimit {
		cursor = fmt.Sprintf("%d", epoch.RewardItems[len(epoch.RewardItems)-1].Index)
	}

	return c.JSON(http.StatusOK, Response{
		Data:   model.NewEpochTransaction(epoch),
		Cursor: cursor,
	})
}

func (h *Hub) GetEpochNodeRewardsHandler(c echo.Context) error {
	var request GetEpochNodeRewardsRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bad request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return response.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	epochs, err := h.databaseClient.FindEpochNodeRewards(c.Request().Context(), request.NodeAddress, request.Limit, request.Cursor)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return response.InternalError(c, fmt.Errorf("get failed: %w", err))
	}

	var cursor string
	if len(epochs) > 0 && len(epochs) == request.Limit {
		cursor = fmt.Sprintf("%d", epochs[len(epochs)-1].ID)
	}

	return c.JSON(http.StatusOK, Response{
		Data:   model.NewEpochs(epochs),
		Cursor: cursor,
	})
}

type GetEpochsRequest struct {
	Cursor *string `query:"cursor"`
	Limit  int     `query:"limit" validate:"min=1,max=50" default:"10"`
}

type GetEpochRequest struct {
	ID         uint64  `param:"id" validate:"required"`
	ItemsLimit int     `query:"itemsLimit" validate:"min=1,max=50" default:"10"`
	Cursor     *string `query:"cursor"`
}

type GetEpochDistributionRequest struct {
	TransactionHash common.Hash `param:"transaction" validate:"required"`
	ItemsLimit      int         `query:"itemsLimit" validate:"min=1,max=50" default:"10"`
	Cursor          *string     `query:"cursor"`
}

type GetEpochNodeRewardsRequest struct {
	NodeAddress common.Address `param:"node" validate:"required"`
	Limit       int            `query:"limit" validate:"min=1,max=50" default:"10"`
	Cursor      *string        `query:"cursor"`
}
