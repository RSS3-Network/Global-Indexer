package nta

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
)

func (n *NTA) GetEpochs(c echo.Context) error {
	var request nta.GetEpochsRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bad request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	epochs, err := n.databaseClient.FindEpochs(c.Request().Context(), request.Limit, request.Cursor)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return errorx.InternalError(c, fmt.Errorf("get failed: %w", err))
	}

	data := nta.NewEpochs(epochs)

	var cursor string
	if len(data) > 0 && len(data) == request.Limit {
		cursor = fmt.Sprintf("%d", data[len(data)-1].ID)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data:   data,
		Cursor: cursor,
	})
}

func (n *NTA) GetEpoch(c echo.Context) error {
	var request nta.GetEpochRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bad request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	epoch, err := n.databaseClient.FindEpochTransactions(c.Request().Context(), request.ID, request.ItemsLimit, request.Cursor)
	if errors.Is(err, database.ErrorRowNotFound) || len(epoch) == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	if err != nil {
		return errorx.InternalError(c, fmt.Errorf("get failed: %w", err))
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: nta.NewEpoch(request.ID, epoch),
	})
}

func (n *NTA) GetEpochDistribution(c echo.Context) error {
	var request nta.GetEpochDistributionRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bad request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	epoch, err := n.databaseClient.FindEpochTransaction(c.Request().Context(), request.TransactionHash, request.ItemsLimit, request.Cursor)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return errorx.InternalError(c, fmt.Errorf("get failed: %w", err))
	}

	var cursor string
	if len(epoch.RewardedNodes) > 0 && len(epoch.RewardedNodes) == request.ItemsLimit {
		cursor = fmt.Sprintf("%d", epoch.RewardedNodes[len(epoch.RewardedNodes)-1].Index)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data:   nta.NewEpochTransaction(epoch),
		Cursor: cursor,
	})
}

func (n *NTA) GetEpochNodeRewards(c echo.Context) error {
	var request nta.GetEpochNodeRewardsRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bad request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	epochs, err := n.databaseClient.FindEpochNodeRewards(c.Request().Context(), request.NodeAddress, request.Limit, request.Cursor)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return errorx.InternalError(c, fmt.Errorf("get failed: %w", err))
	}

	var cursor string
	if len(epochs) > 0 && len(epochs) == request.Limit {
		cursor = fmt.Sprintf("%d", epochs[len(epochs)-1].ID)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data:   nta.NewEpochs(epochs),
		Cursor: cursor,
	})
}
