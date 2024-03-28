package hub

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model/response"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type GetStakeTransactionsRequest struct {
	Cursor  *common.Hash                 `query:"cursor"`
	Staker  *common.Address              `query:"staker"`
	Node    *common.Address              `query:"node"`
	Type    *schema.StakeTransactionType `query:"type"`
	Pending *bool                        `query:"pending"`
	Limit   int                          `query:"limit" default:"20" min:"1" max:"20"`
}

func (h *Hub) GetStakeTransactions(c echo.Context) error {
	var request GetStakeTransactionsRequest
	if err := c.Bind(&request); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return response.InternalError(c, err)
	}

	databaseTransactionOptions := sql.TxOptions{
		ReadOnly: true,
	}

	databaseTransaction, err := h.databaseClient.Begin(c.Request().Context(), &databaseTransactionOptions)
	if err != nil {
		zap.L().Error("begin database transaction", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	defer lo.Try(databaseTransaction.Rollback)

	stakeTransactionsQuery := schema.StakeTransactionsQuery{
		Cursor:  request.Cursor,
		User:    request.Staker,
		Node:    request.Node,
		Type:    request.Type,
		Pending: request.Pending,
		Limit:   request.Limit,
	}

	stakeTransactions, err := databaseTransaction.FindStakeTransactions(c.Request().Context(), stakeTransactionsQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find stake transactions", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	stakeEventsQuery := schema.StakeEventsQuery{
		IDs: lo.Map(stakeTransactions, func(transaction *schema.StakeTransaction, _ int) common.Hash {
			return transaction.ID
		}),
	}

	stakeEvents, err := databaseTransaction.FindStakeEvents(c.Request().Context(), stakeEventsQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find stake events", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	chipsIDs := lo.Flatten(lo.FilterMap(stakeTransactions, func(stakeTransaction *schema.StakeTransaction, _ int) ([]*big.Int, bool) {
		return stakeTransaction.Chips, len(stakeTransaction.Chips) != 0
	}))

	stakeChipsQuery := schema.StakeChipsQuery{
		IDs: chipsIDs,
	}

	stakeChips, err := databaseTransaction.FindStakeChips(c.Request().Context(), stakeChipsQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find stake chips", zap.Error(err), zap.Any("request", request))
	}

	if err := databaseTransaction.Commit(); err != nil {
		return fmt.Errorf("commit database transaction")
	}

	stakeTransactionModels := make([]*model.StakeTransaction, 0, len(stakeTransactions))

	for _, stakeTransaction := range stakeTransactions {
		stakeEvents := lo.Filter(stakeEvents, func(stakeEvent *schema.StakeEvent, _ int) bool {
			return stakeEvent.ID == stakeTransaction.ID
		})

		stakeTransactionModels = append(stakeTransactionModels, model.NewStakeTransaction(stakeTransaction, stakeEvents, stakeChips, baseURL(c)))
	}

	response := Response{
		Data: stakeTransactionModels,
	}

	if length := len(stakeTransactionModels); length > 0 && length == request.Limit {
		response.Cursor = stakeTransactionModels[length-1].ID.String()
	}

	return c.JSON(http.StatusOK, response)
}

type GetStakeTransactionRequest struct {
	ID   *common.Hash                 `param:"id"`
	Type *schema.StakeTransactionType `query:"type"`
}

func (h *Hub) GetStakeTransaction(c echo.Context) error {
	var request GetStakeTransactionRequest
	if err := c.Bind(&request); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	databaseTransactionOptions := sql.TxOptions{
		ReadOnly: true,
	}

	databaseTransaction, err := h.databaseClient.Begin(c.Request().Context(), &databaseTransactionOptions)
	if err != nil {
		zap.L().Error("begin database transaction", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	defer lo.Try(databaseTransaction.Rollback)

	stakeTransactionQuery := schema.StakeTransactionQuery{
		ID:   request.ID,
		Type: request.Type,
	}

	stakeTransaction, err := databaseTransaction.FindStakeTransaction(c.Request().Context(), stakeTransactionQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find stake transaction", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	stakeEventsQuery := schema.StakeEventsQuery{
		IDs: []common.Hash{stakeTransaction.ID},
	}

	stakeEvents, err := databaseTransaction.FindStakeEvents(c.Request().Context(), stakeEventsQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find stake events", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	stakeChipsQuery := schema.StakeChipsQuery{
		IDs: stakeTransaction.Chips,
	}

	stakeChips, err := databaseTransaction.FindStakeChips(c.Request().Context(), stakeChipsQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find stake chips", zap.Error(err), zap.Any("request", request))
	}

	if err := databaseTransaction.Commit(); err != nil {
		return fmt.Errorf("commit database transaction")
	}

	stakeEvents = lo.Filter(stakeEvents, func(stakeEvent *schema.StakeEvent, _ int) bool {
		return stakeEvent.ID == stakeTransaction.ID
	})

	var response Response
	response.Data = model.NewStakeTransaction(stakeTransaction, stakeEvents, stakeChips, baseURL(c))

	return c.JSON(http.StatusOK, response)
}

type GetStakeChipsRequest struct {
	Cursor *big.Int        `query:"cursor"`
	IDs    []*big.Int      `query:"id"`
	Node   *common.Address `query:"node"`
	Owner  *common.Address `query:"owner"`
	Limit  int             `query:"limit" default:"10" min:"1" max:"10"`
}

func (h *Hub) GetStakeChips(c echo.Context) error {
	var request GetStakeChipsRequest
	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return response.InternalError(c, err)
	}

	stakeChipsQuery := schema.StakeChipsQuery{
		Cursor: request.Cursor,
		IDs:    request.IDs,
		Node:   request.Node,
		Owner:  request.Owner,
		Limit:  &request.Limit,
	}

	stakeChips, err := h.databaseClient.FindStakeChips(c.Request().Context(), stakeChipsQuery)
	if err != nil {
		return fmt.Errorf("find stake chips: %w", err)
	}

	// Get current chip values
	nodeAddresses := lo.Map(stakeChips, func(stakeChip *schema.StakeChip, _ int) common.Address {
		return stakeChip.Node
	})

	node, err := h.databaseClient.FindNodes(c.Request().Context(), nodeAddresses, nil, nil, len(nodeAddresses))
	if err != nil {
		return fmt.Errorf("find nodes: %w", err)
	}

	values := lo.SliceToMap(node, func(node *schema.Node) (common.Address, decimal.Decimal) {
		return node.Address, node.MinTokensToStake
	})

	for _, chip := range stakeChips {
		chip.LatestValue = values[chip.Node]
	}

	var response Response
	response.Data = lo.Map(stakeChips, func(stakeChip *schema.StakeChip, _ int) *model.StakeChip {
		return model.NewStakeChip(stakeChip, baseURL(c))
	})

	if length := len(stakeChips); length > 0 && length == request.Limit {
		response.Cursor = stakeChips[length-1].ID.String()
	}

	return c.JSON(http.StatusOK, response)
}

type GetStakeChipRequest struct {
	ID *big.Int `param:"id"`
}

func (h *Hub) GetStakeChip(c echo.Context) error {
	var request GetStakeChipRequest
	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return response.InternalError(c, err)
	}

	stakeChipQuery := schema.StakeChipQuery{
		ID: request.ID,
	}

	stakeChip, err := h.databaseClient.FindStakeChip(c.Request().Context(), stakeChipQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNoContent)
		}

		return err
	}

	node, err := h.databaseClient.FindNode(c.Request().Context(), stakeChip.Node)
	if err != nil {
		return fmt.Errorf("find node: %w", err)
	}

	stakeChip.LatestValue = node.MinTokensToStake

	var response Response
	response.Data = model.NewStakeChip(stakeChip, baseURL(c))

	return c.JSON(http.StatusOK, response)
}

type GetStakeChipsImageRequest struct {
	ID *big.Int `param:"id"`
}

func (h *Hub) GetStakeChipImage(c echo.Context) error {
	var request GetStakeChipsImageRequest
	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return response.InternalError(c, err)
	}

	stakeChipQuery := schema.StakeChipQuery{
		ID: request.ID,
	}

	chip, err := h.databaseClient.FindStakeChip(c.Request().Context(), stakeChipQuery)
	if err != nil {
		return fmt.Errorf("find stake chip: %w", err)
	}

	var metadata l2.ChipsTokenMetadata
	if err := json.Unmarshal(chip.Metadata, &metadata); err != nil {
		return fmt.Errorf("invalid metadata: %w", err)
	}

	data, found := strings.CutPrefix(metadata.Image, "data:image/svg+xml;base64,")
	if !found {
		return fmt.Errorf("invalid image")
	}

	content, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("invalid data: %w", err)
	}

	return c.Blob(http.StatusOK, "image/svg+xml", content)
}

type GetStakeStakingsRequest struct {
	Cursor *string         `query:"cursor"`
	Staker *common.Address `query:"staker"`
	Node   *common.Address `query:"node"`
	Limit  int             `query:"limit" default:"2" min:"1" max:"10"`
}

func (h *Hub) GetStakeStakings(c echo.Context) error {
	var request GetStakeStakingsRequest
	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return response.InternalError(c, err)
	}

	stakeStakingsQuery := schema.StakeStakingsQuery{
		Cursor: request.Cursor,
		Staker: request.Staker,
		Node:   request.Node,
		Limit:  request.Limit,
	}

	stakeStakings, err := h.databaseClient.FindStakeStakings(c.Request().Context(), stakeStakingsQuery)
	if err != nil {
		return err
	}

	response := Response{
		Data: model.NewStakeStaking(stakeStakings, baseURL(c)),
	}

	if length := len(stakeStakings); length > 0 && length == request.Limit {
		response.Cursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%s", stakeStakings[length-1].Staker.String(), stakeStakings[length-1].Node.String())))
	}

	return c.JSON(http.StatusOK, response)
}

type GetStakeOwnerProfitRequest struct {
	Owner common.Address `param:"owner" validate:"required"`
}

type GetStakeOwnerProfitResponse struct {
	Owner            common.Address                           `json:"owner"`
	TotalChipAmounts decimal.Decimal                          `json:"totalChipAmounts"`
	TotalChipValues  decimal.Decimal                          `json:"totalChipValues"`
	OneDay           *GetStakeOwnerProfitChangesSinceResponse `json:"oneDay"`
	OneWeek          *GetStakeOwnerProfitChangesSinceResponse `json:"oneWeek"`
	OneMonth         *GetStakeOwnerProfitChangesSinceResponse `json:"oneMonth"`
}

type GetStakeOwnerProfitChangesSinceResponse struct {
	Date             time.Time       `json:"date"`
	TotalChipAmounts decimal.Decimal `json:"totalChipAmounts"`
	TotalChipValues  decimal.Decimal `json:"totalChipValues"`
	PNL              decimal.Decimal `json:"pnl"`
}

func (h *Hub) GetStakeOwnerProfit(c echo.Context) error {
	var request GetStakeOwnerProfitRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, err)
	}

	// Find all stake chips
	data, err := h.findChipsByOwner(c.Request().Context(), request.Owner)
	if err != nil {
		return response.InternalError(c, err)
	}

	// Find history profit snapshots
	changes, err := h.findStakerHistoryProfitSnapshots(c.Request().Context(), request.Owner, data)
	if err != nil {
		return response.InternalError(c, err)
	}

	data.OneDay = changes[0]
	data.OneWeek = changes[1]
	data.OneMonth = changes[2]

	return c.JSON(http.StatusOK, Response{
		Data: data,
	})
}

func (h *Hub) findChipsByOwner(ctx context.Context, owner common.Address) (*GetStakeOwnerProfitResponse, error) {
	var (
		cursor *big.Int
		data   = &GetStakeOwnerProfitResponse{Owner: owner}
	)

	for {
		chips, err := h.databaseClient.FindStakeChips(ctx, schema.StakeChipsQuery{
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
			minTokensToStake, err := h.stakingContract.MinTokensToStake(nil, node)
			if err != nil {
				return nil, fmt.Errorf("get min tokens from rpc: %w", err)
			}

			data.TotalChipAmounts = data.TotalChipAmounts.Add(decimal.NewFromInt(int64(count)))
			data.TotalChipValues = data.TotalChipValues.Add(decimal.NewFromBigInt(minTokensToStake, 0).Mul(decimal.NewFromInt(int64(count))))
		}

		cursor = chips[len(chips)-1].ID
	}

	return data, nil
}

func (h *Hub) findStakerHistoryProfitSnapshots(ctx context.Context, owner common.Address, profit *GetStakeOwnerProfitResponse) ([]*GetStakeOwnerProfitChangesSinceResponse, error) {
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

	snapshots, err := h.databaseClient.FindStakerProfitSnapshots(ctx, query)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return nil, fmt.Errorf("find staker profit snapshots: %w", err)
	}

	data := make([]*GetStakeOwnerProfitChangesSinceResponse, len(query.Dates))

	for _, snapshot := range snapshots {
		if snapshot.TotalChipValues.IsZero() {
			continue
		}

		var index int

		if snapshot.Date.After(query.Dates[2]) && snapshot.Date.Before(query.Dates[1]) {
			index = 2
		} else if snapshot.Date.After(query.Dates[1]) && snapshot.Date.Before(query.Dates[0]) {
			index = 1
		}

		data[index] = &GetStakeOwnerProfitChangesSinceResponse{
			Date:             snapshot.Date,
			TotalChipAmounts: snapshot.TotalChipAmounts,
			TotalChipValues:  snapshot.TotalChipValues,
			PNL:              profit.TotalChipValues.Sub(snapshot.TotalChipValues).Div(snapshot.TotalChipValues),
		}
	}

	return data, nil
}
