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
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	nodeMinTokensToStakeSnapshots, err := n.databaseClient.FindNodeMinTokensToStakeSnapshots(c.Request().Context(), request.NodeAddresses, request.OnlyStartAndEnd, nil)
	if err != nil {
		zap.L().Error("find Node min tokens to stake snapshots", zap.Error(err))

		return errorx.InternalError(c, fmt.Errorf("find Node min tokens to stake snapshots: %w", err))
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: nta.NewNodeMinTokensToStakeSnapshots(nodeMinTokensToStakeSnapshots),
	})
}

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

func (n *NTA) GetStakerProfitsSnapshots(c echo.Context) error {
	var request nta.GetStakerProfitSnapshotsRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	query := schema.StakerProfitSnapshotsQuery{
		OwnerAddress: lo.ToPtr(request.OwnerAddress),
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

func (n *NTA) GetOperatorProfitsSnapshots(c echo.Context) error {
	var request nta.GetOperatorProfitSnapshotsRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	query := schema.OperatorProfitSnapshotsQuery{
		Operator:   lo.ToPtr(request.Operator),
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
