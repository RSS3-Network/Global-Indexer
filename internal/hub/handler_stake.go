package hub

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model/response"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type GetStakeTransactionsRequest struct {
	Cursor  *common.Hash                 `query:"cursor"`
	Staker  *common.Address              `query:"staker"`
	Node    *common.Address              `query:"node"`
	Type    *schema.StakeTransactionType `query:"type"`
	Pending *bool                        `query:"pending"`
}

func (h *Hub) GetStakeTransactions(c echo.Context) error {
	var request GetStakeTransactionsRequest
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

	stakeTransactionsQuery := schema.StakeTransactionsQuery{
		Cursor:  request.Cursor,
		User:    request.Staker,
		Node:    request.Node,
		Type:    request.Type,
		Pending: request.Pending,
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

	if err := databaseTransaction.Commit(); err != nil {
		return fmt.Errorf("commit database transaction")
	}

	stakeTransactionModels := make([]*model.StakeTransaction, 0, len(stakeTransactions))

	for _, stakeTransaction := range stakeTransactions {
		stakeEvents := lo.Filter(stakeEvents, func(stakeEvent *schema.StakeEvent, _ int) bool {
			return stakeEvent.ID == stakeTransaction.ID
		})

		stakeTransactionModels = append(stakeTransactionModels, model.NewStakeTransaction(stakeTransaction, stakeEvents))
	}

	response := Response{
		Data: stakeTransactionModels,
	}

	if length := len(stakeTransactionModels); length > 0 {
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

	if err := databaseTransaction.Commit(); err != nil {
		return fmt.Errorf("commit database transaction")
	}

	stakeEvents = lo.Filter(stakeEvents, func(stakeEvent *schema.StakeEvent, _ int) bool {
		return stakeEvent.ID == stakeTransaction.ID
	})

	var response Response
	response.Data = model.NewStakeTransaction(stakeTransaction, stakeEvents)

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
		Limit:  request.Limit,
	}

	stakeChips, err := h.databaseClient.FindStakeChips(c.Request().Context(), stakeChipsQuery)
	if err != nil {
		return err
	}

	var response Response
	response.Data = lo.Map(stakeChips, func(stakeChip *schema.StakeChip, _ int) *model.StakeChip {
		return model.NewStakeChip(stakeChip, baseURL(c))
	})

	if length := len(stakeChips); length > 0 {
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

	if length := len(stakeStakings); length > 0 {
		response.Cursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%s", stakeStakings[length-1].Staker.String(), stakeStakings[length-1].Node.String())))
	}

	return c.JSON(http.StatusOK, response)
}
