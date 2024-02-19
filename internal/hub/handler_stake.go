package hub

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type GetStakeTransactionsRequest struct {
	Cursor  *common.Hash                 `query:"cursor"`
	User    *common.Address              `query:"user"`
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
		User:    request.User,
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

type GetStakeWalletsRequest struct {
	Cursor *common.Address `query:"cursor"`
}

func (h *Hub) GetStakeWallets(c echo.Context) error {
	var request GetStakeWalletsRequest
	if err := c.Bind(&request); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	stakeChipsQuery := schema.StakeChipsQuery{
		Cursor: request.Cursor,
	}

	stakeChips, err := h.databaseClient.FindStakeChips(c.Request().Context(), stakeChipsQuery)
	if err != nil {
		zap.L().Error("find node wallets", zap.Error(err))

		return c.NoContent(http.StatusInternalServerError)
	}

	response := Response{
		Data: model.NewStakeStakers(stakeChips),
	}

	if length := len(stakeChips); length > 0 {
		response.Cursor = stakeChips[length-1].Owner.String()
	}

	return c.JSON(http.StatusOK, response)
}

type GetStakeNodeChipsRequest struct {
	Node common.Address `param:"node"`
}

func (h *Hub) GetStakeNodeChips(c echo.Context) error {
	var request GetStakeNodeChipsRequest
	if err := c.Bind(&request); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	stakeChipsQuery := schema.StakeChipsQuery{
		Node:   &request.Node,
		Direct: true,
	}

	stakeChips, err := h.databaseClient.FindStakeChips(c.Request().Context(), stakeChipsQuery)
	if err != nil {
		zap.L().Error("find node chips", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	response := Response{
		Data: model.NewStakeStakers(stakeChips),
	}

	if length := len(stakeChips); length > 0 {
		response.Cursor = stakeChips[length-1].Owner.String()
	}

	return c.JSON(http.StatusOK, response)
}

type GetStakeWalletChipsRequest struct {
	Wallet common.Address `param:"wallet"`
}

func (h *Hub) GetStakeWalletChips(c echo.Context) error {
	var request GetStakeWalletChipsRequest
	if err := c.Bind(&request); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	stakeChipsQuery := schema.StakeChipsQuery{
		Owner: &request.Wallet,
	}

	stakeChips, err := h.databaseClient.FindStakeChips(c.Request().Context(), stakeChipsQuery)
	if err != nil {
		zap.L().Error("find node chips", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	response := Response{
		Data: model.NewStakeNodes(stakeChips),
	}

	if length := len(stakeChips); length > 0 {
		response.Cursor = stakeChips[length-1].Owner.String()
	}

	return c.JSON(http.StatusOK, response)
}
