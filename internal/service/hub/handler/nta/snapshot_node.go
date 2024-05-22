package nta

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (n *NTA) GetNodeCountSnapshots(c echo.Context) error {
	nodeSnapshots, err := n.databaseClient.FindNodeCountSnapshots(c.Request().Context())
	if err != nil {
		zap.L().Error("find Node snapshots", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	var response nta.Response

	response.Data = nta.NewNodeCountSnapshots(nodeSnapshots)

	return c.JSON(http.StatusOK, response)
}

func (n *NTA) BatchGetNodeMinTokensToStakeSnapshots(c echo.Context) error {
	var request nta.BatchNodeMinTokensToStakeRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	nodeMinTokensToStakeSnapshots, err := n.databaseClient.FindNodeMinTokensToStakeSnapshots(c.Request().Context(), request.NodeAddresses, request.OnlyStartAndEnd, nil)
	if err != nil {
		zap.L().Error("find Node min tokens to stake snapshots", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: nta.NewNodeMinTokensToStakeSnapshots(nodeMinTokensToStakeSnapshots),
	})
}

func (n *NTA) GetNodeOperationProfitSnapshots(c echo.Context) error {
	var request nta.GetNodeOperationProfitSnapshotsRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}
	// FIXME: OperatorProfit -> NodeOperationProfit
	query := schema.OperatorProfitSnapshotsQuery{
		Operator:   lo.ToPtr(request.NodeAddress),
		Limit:      request.Limit,
		Cursor:     request.Cursor,
		BeforeDate: request.BeforeDate,
		AfterDate:  request.AfterDate,
	}

	operatorProfitSnapshots, err := n.databaseClient.FindOperatorProfitSnapshots(c.Request().Context(), query)
	if err != nil {
		zap.L().Error("find operator profit snapshots", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	var cursor string

	if request.Limit != nil && len(operatorProfitSnapshots) > 0 && len(operatorProfitSnapshots) == lo.FromPtr(request.Limit) {
		last, _ := lo.Last(operatorProfitSnapshots)
		cursor = fmt.Sprintf("%d", last.ID)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data:   operatorProfitSnapshots,
		Cursor: cursor,
	})
}

func (n *NTA) GetEpochsAPYSnapshots(c echo.Context) error {
	epochAPYSnapshots, err := n.databaseClient.FindEpochAPYSnapshots(c.Request().Context(), schema.EpochAPYSnapshotQuery{})
	if err != nil {
		zap.L().Error("find epoch APY snapshots", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: epochAPYSnapshots,
	})
}

func (n *NTA) FindNodeOperationProfitSnapshots(ctx context.Context, operator common.Address, profit *nta.GetNodeOperationProfitResponse) ([]*nta.NodeProfitChangeDetail, error) {
	if profit == nil {
		return nil, nil
	}

	now := time.Now()
	query := schema.OperatorProfitSnapshotsQuery{
		Operator: lo.ToPtr(operator),
		Dates: []time.Time{
			now.Add(-24 * time.Hour),      // 1 day
			now.Add(-7 * 24 * time.Hour),  // 1 week
			now.Add(-30 * 24 * time.Hour), // 1 month
		},
	}

	snapshots, err := n.databaseClient.FindOperatorProfitSnapshots(ctx, query)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return nil, fmt.Errorf("find operator profit snapshots: %w", err)
	}

	data := make([]*nta.NodeProfitChangeDetail, len(query.Dates))

	for _, snapshot := range snapshots {
		if snapshot.OperationPool.IsZero() {
			continue
		}

		var index int

		if snapshot.Date.After(query.Dates[2]) && snapshot.Date.Before(query.Dates[1]) {
			index = 2
		} else if snapshot.Date.After(query.Dates[1]) && snapshot.Date.Before(query.Dates[0]) {
			index = 1
		}

		data[index] = &nta.NodeProfitChangeDetail{
			Date:          snapshot.Date,
			OperationPool: snapshot.OperationPool,
			ProfitAndLoss: profit.OperationPool.Sub(snapshot.OperationPool).Div(snapshot.OperationPool),
		}
	}

	return data, nil
}
