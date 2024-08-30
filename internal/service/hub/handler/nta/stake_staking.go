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
	"github.com/shopspring/decimal"
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

	var blockTimestamp time.Time

	profit := nta.GetStakerProfitResponseData{
		Owner: owner,
		OneDay: &nta.GetStakerProfitChangesSinceResponseData{
			Date: time.Now().AddDate(0, 0, -1),
		},
		OneWeek: &nta.GetStakerProfitChangesSinceResponseData{
			Date: time.Now().AddDate(0, 0, -7),
		},
		OneMonth: &nta.GetStakerProfitChangesSinceResponseData{
			Date: time.Now().AddDate(0, -1, 0),
		},
	}

	if len(currentProfit) > 0 {
		blockTimestamp = currentProfit[0].Date

		profit.TotalChipAmount = currentProfit[0].TotalChipAmount
		profit.TotalChipValue = currentProfit[0].TotalChipValue
	}

	// Calculate profit changes from staking transactions.
	transactions, err := n.databaseClient.FindStakeTransactions(ctx, schema.StakeTransactionsQuery{
		User:           lo.ToPtr(owner),
		BlockTimestamp: lo.ToPtr(blockTimestamp),
		Order:          "block_timestamp ASC",
	})
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return nil, fmt.Errorf("find stake transactions: %w", err)
	}

	for _, transaction := range transactions {
		switch transaction.Type {
		case schema.StakeTransactionTypeStake:
			profit.TotalChipAmount = profit.TotalChipAmount.Add(decimal.NewFromInt(int64(len(transaction.Chips))))
			profit.TotalChipValue = profit.TotalChipValue.Add(decimal.NewFromBigInt(transaction.Value, 0))
		case schema.StakeTransactionTypeUnstake:
			profit.TotalChipAmount = profit.TotalChipAmount.Sub(decimal.NewFromInt(int64(len(transaction.Chips))))
			profit.TotalChipValue = profit.TotalChipValue.Sub(decimal.NewFromBigInt(transaction.Value, 0))
		case schema.StakeEventTypeChipsMerged:
			profit.TotalChipAmount = profit.TotalChipAmount.Sub(decimal.NewFromInt(int64(len(transaction.Chips) - 2))) // Exclude the merged chips.
		}
	}

	// Find history profit snapshots.
	query = schema.StakerProfitSnapshotsQuery{
		OwnerAddress: lo.ToPtr(owner),
		Dates:        []time.Time{profit.OneDay.Date, profit.OneWeek.Date, profit.OneMonth.Date},
	}

	snapshots, err := n.databaseClient.FindStakerProfitSnapshots(ctx, query)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return nil, fmt.Errorf("find staker profit snapshots: %w", err)
	}

	for _, snapshot := range snapshots {
		if snapshot.TotalChipValue.IsZero() {
			continue
		}

		switch {
		case snapshot.Date.After(profit.OneMonth.Date) && snapshot.Date.Before(profit.OneWeek.Date):
			profit.OneMonth = &nta.GetStakerProfitChangesSinceResponseData{
				Date:            snapshot.Date,
				TotalChipAmount: snapshot.TotalChipAmount,
				TotalChipValue:  snapshot.TotalChipValue,
				ProfitAndLoss:   profit.TotalChipValue.Sub(snapshot.TotalChipValue).Div(snapshot.TotalChipValue),
			}
		case snapshot.Date.After(profit.OneWeek.Date) && snapshot.Date.Before(profit.OneDay.Date):
			profit.OneWeek = &nta.GetStakerProfitChangesSinceResponseData{
				Date:            snapshot.Date,
				TotalChipAmount: snapshot.TotalChipAmount,
				TotalChipValue:  snapshot.TotalChipValue,
				ProfitAndLoss:   profit.TotalChipValue.Sub(snapshot.TotalChipValue).Div(snapshot.TotalChipValue),
			}
		default:
			profit.OneDay = &nta.GetStakerProfitChangesSinceResponseData{
				Date:            snapshot.Date,
				TotalChipAmount: snapshot.TotalChipAmount,
				TotalChipValue:  snapshot.TotalChipValue,
				ProfitAndLoss:   profit.TotalChipValue.Sub(snapshot.TotalChipValue).Div(snapshot.TotalChipValue),
			}
		}
	}

	return &profit, nil
}

func (n *NTA) GetStakingStat(c echo.Context) error {
	var request nta.GetStakingStatRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	stakingStat, err := n.databaseClient.FindStakeStaker(c.Request().Context(), request.Address)
	if err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("fetch stake staker by address %s: %w", request.Address, err))
	}

	// Calculate profit changes from staking transactions.
	transactions, err := n.databaseClient.FindStakeTransactions(c.Request().Context(), schema.StakeTransactionsQuery{
		User:      lo.ToPtr(request.Address),
		Finalized: lo.ToPtr(false),
		Order:     "block_timestamp ASC",
	})
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return errorx.ValidationFailedError(c, fmt.Errorf("find stake transactions: %w", err))
	}

	for _, transaction := range transactions {
		switch transaction.Type {
		case schema.StakeTransactionTypeStake:
			stakingStat.TotalChips = stakingStat.TotalChips + uint64(len(transaction.Chips))
			stakingStat.TotalStakedTokens = stakingStat.TotalStakedTokens.Add(decimal.NewFromBigInt(transaction.Value, 0))
			stakingStat.CurrentStakedTokens = stakingStat.CurrentStakedTokens.Add(decimal.NewFromBigInt(transaction.Value, 0))
		case schema.StakeTransactionTypeUnstake:
			stakingStat.TotalChips = stakingStat.TotalChips - uint64(len(transaction.Chips))
			stakingStat.TotalStakedTokens = stakingStat.TotalStakedTokens.Sub(decimal.NewFromBigInt(transaction.Value, 0))
			stakingStat.CurrentStakedTokens = stakingStat.CurrentStakedTokens.Sub(decimal.NewFromBigInt(transaction.Value, 0))
		case schema.StakeEventTypeChipsMerged:
			stakingStat.TotalChips = stakingStat.TotalChips - uint64(len(transaction.Chips)-2) // Exclude the merged chips.
		}
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: stakingStat,
	})
}
