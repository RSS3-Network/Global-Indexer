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
	User *string `query:"user"`
	Node *string `query:"node"`
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

	var transactions []*schema.StakeTransaction

	switch {
	case request.User != nil:
		transactions, err = databaseTransaction.FindStakeTransactionsByUser(c.Request().Context(), common.HexToAddress(*request.User))
	case request.Node != nil:
		transactions, err = databaseTransaction.FindStakeTransactionsByNode(c.Request().Context(), common.HexToAddress(*request.Node))
	default:
		transactions, err = databaseTransaction.FindStakeTransactions(c.Request().Context())
	}

	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find stake transactions", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	ids := lo.Map(transactions, func(transaction *schema.StakeTransaction, _ int) common.Hash {
		return transaction.ID
	})

	events, err := databaseTransaction.FindStakeEventsByIDs(c.Request().Context(), ids)
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

	transactionModels := make([]*model.StakeTransaction, 0, len(transactions))

	for _, transaction := range transactions {
		events := lo.Filter(events, func(event *schema.StakeEvent, _ int) bool {
			return event.ID == transaction.ID
		})

		transactionModels = append(transactionModels, model.NewStakeTransaction(transaction, events))
	}

	var response Response

	response.Data = transactionModels

	return c.JSON(http.StatusOK, response)
}

type GetStakeTransactionRequest struct {
	ID *string `param:"id"`
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

	transaction, err := databaseTransaction.FindStakeTransaction(c.Request().Context(), common.HexToHash(*request.ID))
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find stake transaction", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	events, err := databaseTransaction.FindStakeEventsByIDs(c.Request().Context(), []common.Hash{transaction.ID})
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

	events = lo.Filter(events, func(event *schema.StakeEvent, _ int) bool {
		return event.ID == transaction.ID
	})

	var response Response
	response.Data = model.NewStakeTransaction(transaction, events)

	return c.JSON(http.StatusOK, response)
}

type GetStakeNodeStakersRequest struct {
	Node common.Address `param:"node"`
}

func (h *Hub) GetStakeNodeStakers(c echo.Context) error {
	var request GetStakeNodeStakersRequest
	if err := c.Bind(&request); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	chips, err := h.databaseClient.FindStakeChipsByNode(c.Request().Context(), request.Node)
	if err != nil {
		zap.L().Error("find node stakers", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	response := Response{
		Data: model.NewStakeStakers(chips),
	}

	return c.JSON(http.StatusOK, response)
}

type GetStakeStakerNodesRequest struct {
	Staker common.Address `param:"staker"`
}

func (h *Hub) GetStakeStakerNodes(c echo.Context) error {
	var request GetStakeStakerNodesRequest
	if err := c.Bind(&request); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	chips, err := h.databaseClient.FindStakeChipsByOwner(c.Request().Context(), request.Staker)
	if err != nil {
		zap.L().Error("find node stakers", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	response := Response{
		Data: model.NewStakeNodes(chips),
	}

	return c.JSON(http.StatusOK, response)
}
