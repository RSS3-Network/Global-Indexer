package nta

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/sourcegraph/conc/pool"
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
		response.Cursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%s", stakeStakings[length-1].Staker.String(), stakeStakings[length-1].Node.String())))
	}

	return c.JSON(http.StatusOK, response)
}

// FIXME: what is stake owner?
func (n *NTA) GetStakeOwnerProfit(c echo.Context) error {
	var request nta.GetStakeOwnerProfitRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	// Find all stake chips
	data, err := n.findChipsByOwner(c.Request().Context(), request.StakerAddress)
	if err != nil {
		zap.L().Error("find chips by owner", zap.Error(err))

		return errorx.InternalError(c)
	}

	// Find history profit snapshots
	changes, err := n.findStakerHistoryProfitSnapshots(c.Request().Context(), request.StakerAddress, data)
	if err != nil {
		zap.L().Error("find staker history profit snapshots", zap.Error(err))

		return errorx.InternalError(c)
	}

	data.OneDay, data.OneWeek, data.OneMonth = changes[0], changes[1], changes[2]

	return c.JSON(http.StatusOK, nta.Response{
		Data: data,
	})
}

func (n *NTA) findChipsByOwner(ctx context.Context, owner common.Address) (*nta.GetStakeOwnerProfitResponseData, error) {
	var (
		cursor *big.Int
		mu     sync.Mutex
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

		errorPool := pool.New().WithContext(ctx).WithMaxGoroutines(30).WithCancelOnError().WithFirstError()

		for _, chip := range chips {
			chip := chip
			data.TotalChipAmount = data.TotalChipAmount.Add(decimal.NewFromInt(int64(1)))

			errorPool.Go(func(ctx context.Context) error {
				chipInfo, err := n.stakingContract.GetChipInfo(&bind.CallOpts{Context: ctx}, chip.ID)
				if err != nil {
					zap.L().Error("get chip info from rpc", zap.Error(err), zap.String("chipID", chip.ID.String()))

					return fmt.Errorf("get chip info: %w", err)
				}

				mu.Lock()
				defer mu.Unlock()

				data.TotalChipValue = data.TotalChipValue.Add(decimal.NewFromBigInt(chipInfo.Tokens, 0))

				return nil
			})
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
