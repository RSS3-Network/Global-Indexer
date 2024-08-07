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
	snapshot "github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot/apy"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
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
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	epochs, err := n.databaseClient.FindEpochs(c.Request().Context(), &schema.FindEpochsQuery{
		Limit:  lo.ToPtr(request.Limit),
		Cursor: request.Cursor,
	})
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("get epochs failed", zap.Error(err))

		return errorx.InternalError(c)
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
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	epoch, err := n.databaseClient.FindEpochTransactions(c.Request().Context(), request.EpochID, request.ItemLimit, request.Cursor)
	if errors.Is(err, database.ErrorRowNotFound) || len(epoch) == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	if err != nil {
		zap.L().Error("get epoch failed", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: nta.NewEpoch(request.EpochID, epoch),
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
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	epoch, err := n.databaseClient.FindEpochTransaction(c.Request().Context(), request.TransactionHash, request.ItemLimit, request.Cursor)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("get epoch distribution failed", zap.Error(err))

		return errorx.InternalError(c)
	}

	var cursor string
	if len(epoch.RewardedNodes) > 0 && len(epoch.RewardedNodes) == request.ItemLimit {
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
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	epochs, err := n.databaseClient.FindEpochNodeRewards(c.Request().Context(), request.NodeAddress, request.Limit, request.Cursor)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("get epoch node rewards failed", zap.Error(err))

		return errorx.InternalError(c)
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

func (n *NTA) GetEpochsAPY(c echo.Context) error {
	var apy decimal.Decimal

	// Get from cache if available
	err := n.cacheClient.Get(c.Request().Context(), snapshot.CacheKeyEpochAverageAPY, &apy)
	if err == nil && !apy.IsZero() {
		return c.JSON(http.StatusOK, nta.Response{
			Data: apy,
		})
	}

	// Query the database for the epoch APY
	apy, err = n.databaseClient.FindEpochAPYSnapshotsAverage(c.Request().Context())
	if err != nil {
		zap.L().Error("get epoch apy failed", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: apy,
	})
}
