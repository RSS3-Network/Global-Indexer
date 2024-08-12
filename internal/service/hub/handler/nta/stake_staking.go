package nta

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (n *NTA) GetStakeStakings(c echo.Context) error {
	var request nta.GetStakeStakingsRequest
	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	stakeStakingsQuery := schema.StakeStakingsQuery{
		Cursor: request.Cursor,
		Staker: request.StakerAddress,
		Node:   request.NodeAddress,
		Limit:  request.Limit,
	}

	stakeStakings, err := n.databaseClient.FindStakeStakings(c.Request().Context(), stakeStakingsQuery)
	if err != nil {
		return err
	}

	response := nta.Response{
		Data: nta.NewStakeStaking(stakeStakings, n.baseURL(c)),
	}

	if length := len(stakeStakings); length > 0 && length == request.Limit {
		latestStakeStaking := stakeStakings[length-1]

		response.Cursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%s-%s", latestStakeStaking.Value, latestStakeStaking.Staker, latestStakeStaking.Node)))
	}

	return c.JSON(http.StatusOK, response)
}

func (n *NTA) GetStakerProfit(c echo.Context) error {
	var request nta.GetStakerProfitRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	// Find history profit snapshots
	data, err := n.findStakerHistoryProfitSnapshots(c.Request().Context(), request.StakerAddress)
	if err != nil {
		zap.L().Error("find staker history profit snapshots", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: data,
	})
}

func (n *NTA) findStakerHistoryProfitSnapshots(ctx context.Context, owner common.Address) (*nta.GetStakerProfitResponseData, error) {
	// Find current profit snapshot.
	query := schema.StakerProfitSnapshotsQuery{
		OwnerAddress: lo.ToPtr(owner),
		Limit:        lo.ToPtr(1),
	}

	currentProfit, err := n.databaseClient.FindStakerProfitSnapshots(ctx, query)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return nil, fmt.Errorf("find staker profit snapshots: %w", err)
	}

	if len(currentProfit) == 0 {
		return nil, nil
	}

	// Find history profit snapshots.
	yesterday := currentProfit[0].Date.AddDate(0, 0, -1)
	weekAgo := currentProfit[0].Date.AddDate(0, 0, -7)
	monthAgo := currentProfit[0].Date.AddDate(0, -1, 0)

	query = schema.StakerProfitSnapshotsQuery{
		OwnerAddress: lo.ToPtr(owner),
		Dates:        []time.Time{yesterday, weekAgo, monthAgo},
	}

	snapshots, err := n.databaseClient.FindStakerProfitSnapshots(ctx, query)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return nil, fmt.Errorf("find staker profit snapshots: %w", err)
	}

	data := &nta.GetStakerProfitResponseData{
		Owner:           owner,
		TotalChipAmount: currentProfit[0].TotalChipAmount,
		TotalChipValue:  currentProfit[0].TotalChipValue,
	}

	for _, snapshot := range snapshots {
		if snapshot.TotalChipValue.IsZero() {
			continue
		}

		switch {
		case snapshot.Date.After(monthAgo) && snapshot.Date.Before(weekAgo):
			data.OneMonth = &nta.GetStakerProfitChangesSinceResponseData{
				Date:            snapshot.Date,
				TotalChipAmount: snapshot.TotalChipAmount,
				TotalChipValue:  snapshot.TotalChipValue,
				ProfitAndLoss:   data.TotalChipValue.Sub(snapshot.TotalChipValue).Div(snapshot.TotalChipValue),
			}
		case snapshot.Date.After(weekAgo) && snapshot.Date.Before(yesterday):
			data.OneWeek = &nta.GetStakerProfitChangesSinceResponseData{
				Date:            snapshot.Date,
				TotalChipAmount: snapshot.TotalChipAmount,
				TotalChipValue:  snapshot.TotalChipValue,
				ProfitAndLoss:   data.TotalChipValue.Sub(snapshot.TotalChipValue).Div(snapshot.TotalChipValue),
			}
		default:
			data.OneDay = &nta.GetStakerProfitChangesSinceResponseData{
				Date:            snapshot.Date,
				TotalChipAmount: snapshot.TotalChipAmount,
				TotalChipValue:  snapshot.TotalChipValue,
				ProfitAndLoss:   data.TotalChipValue.Sub(snapshot.TotalChipValue).Div(snapshot.TotalChipValue),
			}
		}
	}

	return data, nil
}

func (n *NTA) GetStakingStat(c echo.Context) error {
	var request nta.GetStakingStatRequest

	if err := c.Bind(&request); err != nil {
		return fmt.Errorf("bind request: %w", err)
	}

	if err := c.Validate(&request); err != nil {
		return fmt.Errorf("validate request: %w", err)
	}

	stakingStat, err := n.databaseClient.FindStakeStaker(c.Request().Context(), request.Address)
	if err != nil {
		return fmt.Errorf("fetch stake staker by address %s: %w", request.Address, err)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: stakingStat,
	})
}
