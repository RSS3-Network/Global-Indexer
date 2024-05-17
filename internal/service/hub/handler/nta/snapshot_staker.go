package nta

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (n *NTA) GetStakerCountSnapshots(c echo.Context) error {
	stakeSnapshots, err := n.databaseClient.FindStakerCountSnapshots(c.Request().Context())
	if err != nil {
		zap.L().Error("find staker_count snapshots", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	var response nta.Response

	response.Data = nta.NewStakeSnapshots(stakeSnapshots)

	return c.JSON(http.StatusOK, response)
}

func (n *NTA) GetStakerProfitSnapshots(c echo.Context) error {
	var request nta.GetStakerProfitSnapshotsRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	query := schema.StakerProfitSnapshotsQuery{
		OwnerAddress: lo.ToPtr(request.StakerAddress),
		Limit:        request.Limit,
		Cursor:       request.Cursor,
		BeforeDate:   request.BeforeDate,
		AfterDate:    request.AfterDate,
	}

	stakerProfitSnapshots, err := n.databaseClient.FindStakerProfitSnapshots(c.Request().Context(), query)
	if err != nil {
		zap.L().Error("find staker profit snapshots", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	var cursor string

	if request.Limit != nil && len(stakerProfitSnapshots) > 0 && len(stakerProfitSnapshots) == lo.FromPtr(request.Limit) {
		last, _ := lo.Last(stakerProfitSnapshots)
		cursor = fmt.Sprintf("%d", last.ID)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data:   stakerProfitSnapshots,
		Cursor: cursor,
	})
}
