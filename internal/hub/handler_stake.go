package hub

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
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

	var transactions []*schema.StakeTransaction

	switch {
	case request.User != nil:
		if transactions, err = databaseTransaction.FindStakeTransactionsByUser(c.Request().Context(), common.HexToAddress(*request.User)); err != nil {
			zap.L().Error("find stake transactions", zap.Error(err), zap.Any("request", request))

			return c.NoContent(http.StatusInternalServerError)
		}
	case request.Node != nil:
		if transactions, err = databaseTransaction.FindStakeTransactionsByNode(c.Request().Context(), common.HexToAddress(*request.Node)); err != nil {
			zap.L().Error("find stake transactions", zap.Error(err), zap.Any("request", request))

			return c.NoContent(http.StatusInternalServerError)
		}
	default:
		if transactions, err = databaseTransaction.FindStakeTransactions(c.Request().Context()); err != nil {
			zap.L().Error("find stake transactions", zap.Error(err), zap.Any("request", request))

			return c.NoContent(http.StatusInternalServerError)
		}
	}

	ids := lo.Map(transactions, func(transaction *schema.StakeTransaction, _ int) common.Hash {
		return transaction.ID
	})

	events, err := databaseTransaction.FindStakeEventsByIDs(c.Request().Context(), ids)
	if err != nil {
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

	transaction, err := databaseTransaction.FindStakeTransaction(c.Request().Context(), common.HexToHash(*request.ID))
	if err != nil {
		zap.L().Error("find stake transaction", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	events, err := databaseTransaction.FindStakeEventsByIDs(c.Request().Context(), []common.Hash{transaction.ID})
	if err != nil {
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
