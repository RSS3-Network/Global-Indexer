package nta

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
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
	"github.com/shopspring/decimal"
)

func (n *NTA) GetStakeStakings(c echo.Context) error {
	var request nta.GetStakeStakingsRequest
	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.InternalError(c, err)
	}

	stakeStakingsQuery := schema.StakeStakingsQuery{
		Cursor: request.Cursor,
		Staker: request.Staker,
		Node:   request.Node,
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
		response.Cursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%s", stakeStakings[length-1].Staker.String(), stakeStakings[length-1].Node.String())))
	}

	return c.JSON(http.StatusOK, response)
}

func (n *NTA) GetStakeOwnerProfit(c echo.Context) error {
	var request nta.GetStakeOwnerProfitRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, err)
	}

	// Find all stake chips
	data, err := n.findChipsByOwner(c.Request().Context(), request.Owner)
	if err != nil {
		return errorx.InternalError(c, err)
	}

	// Find history profit snapshots
	changes, err := n.findStakerHistoryProfitSnapshots(c.Request().Context(), request.Owner, data)
	if err != nil {
		return errorx.InternalError(c, err)
	}

	data.OneDay, data.OneWeek, data.OneMonth = changes[0], changes[1], changes[2]

	return c.JSON(http.StatusOK, nta.Response{
		Data: data,
	})
}

func (n *NTA) findChipsByOwner(ctx context.Context, owner common.Address) (*nta.GetStakeOwnerProfitResponseData, error) {
	var (
		cursor *big.Int
		data   = &nta.GetStakeOwnerProfitResponseData{Owner: owner}
	)

	for {
		chips, err := n.databaseClient.FindStakeChips(ctx, schema.StakeChipsQuery{
			Owner:  lo.ToPtr(owner),
			Cursor: cursor,
			Limit:  lo.ToPtr(500),
		})
		if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
			return nil, fmt.Errorf("find stake chips: %w", err)
		}

		if len(chips) == 0 {
			break
		}

		nodes := make(map[common.Address]int)

		for _, chip := range chips {
			chip := chip
			nodes[chip.Node]++
		}

		for node, count := range nodes {
			minTokensToStake, err := n.stakingContract.MinTokensToStake(nil, node)
			if err != nil {
				return nil, fmt.Errorf("get min tokens from rpc: %w", err)
			}

			data.TotalChipAmount = data.TotalChipAmount.Add(decimal.NewFromInt(int64(count)))
			data.TotalChipValue = data.TotalChipValue.Add(decimal.NewFromBigInt(minTokensToStake, 0).Mul(decimal.NewFromInt(int64(count))))
		}

		cursor = chips[len(chips)-1].ID
	}

	return data, nil
}

func (n *NTA) findStakerHistoryProfitSnapshots(ctx context.Context, owner common.Address, profit *nta.GetStakeOwnerProfitResponseData) ([]*nta.GetStakeOwnerProfitChangesSinceResponseData, error) {
	if profit == nil {
		return nil, nil
	}

	now := time.Now()
	query := schema.StakerProfitSnapshotsQuery{
		OwnerAddress: lo.ToPtr(owner),
		Dates: []time.Time{
			now.Add(-24 * time.Hour),      // 1 day
			now.Add(-7 * 24 * time.Hour),  // 1 week
			now.Add(-30 * 24 * time.Hour), // 1 month
		},
	}

	snapshots, err := n.databaseClient.FindStakerProfitSnapshots(ctx, query)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return nil, fmt.Errorf("find staker profit snapshots: %w", err)
	}

	data := make([]*nta.GetStakeOwnerProfitChangesSinceResponseData, len(query.Dates))

	for _, snapshot := range snapshots {
		if snapshot.TotalChipValue.IsZero() {
			continue
		}

		var index int

		if snapshot.Date.After(query.Dates[2]) && snapshot.Date.Before(query.Dates[1]) {
			index = 2
		} else if snapshot.Date.After(query.Dates[1]) && snapshot.Date.Before(query.Dates[0]) {
			index = 1
		}

		data[index] = &nta.GetStakeOwnerProfitChangesSinceResponseData{
			Date:            snapshot.Date,
			TotalChipAmount: snapshot.TotalChipAmount,
			TotalChipValue:  snapshot.TotalChipValue,
			ProfitAndLoss:   profit.TotalChipValue.Sub(snapshot.TotalChipValue).Div(snapshot.TotalChipValue),
		}
	}

	return data, nil
}
