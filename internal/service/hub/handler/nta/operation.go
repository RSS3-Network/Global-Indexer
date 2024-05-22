package nta

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

func (n *NTA) GetOperatorProfit(c echo.Context) error {
	var request nta.GetOperatorProfitRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	node, err := n.stakingContract.GetNode(&bind.CallOpts{}, request.Operator)
	if err != nil {
		return errorx.InternalError(c, fmt.Errorf("get Node from rpc: %w", err))
	}

	data := nta.GetOperatorProfitRepsonseData{
		Operator:      request.Operator,
		OperationPool: decimal.NewFromBigInt(node.OperationPoolTokens, 0),
	}

	changes, err := n.findOperatorHistoryProfitSnapshots(c.Request().Context(), request.Operator, &data)
	if err != nil {
		return errorx.InternalError(c, fmt.Errorf("find operator history profit snapshots: %w", err))
	}

	data.OneDay, data.OneWeek, data.OneMonth = changes[0], changes[1], changes[2]

	return c.JSON(http.StatusOK, nta.Response{
		Data: data,
	})
}

func (n *NTA) findOperatorHistoryProfitSnapshots(ctx context.Context, operator common.Address, profit *nta.GetOperatorProfitRepsonseData) ([]*nta.GetOperatorProfitChangesSinceResponseData, error) {
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

	data := make([]*nta.GetOperatorProfitChangesSinceResponseData, len(query.Dates))

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

		data[index] = &nta.GetOperatorProfitChangesSinceResponseData{
			Date:          snapshot.Date,
			OperationPool: snapshot.OperationPool,
			PNL:           profit.OperationPool.Sub(snapshot.OperationPool).Div(snapshot.OperationPool),
		}
	}

	return data, nil
}
