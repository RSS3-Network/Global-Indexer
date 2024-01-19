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

type GetBridgeTransactionsRequest struct {
	Address *string `query:"address"`
}

func (h *Hub) GetBridgeTransactions(c echo.Context) error {
	var request GetBridgeTransactionsRequest
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

	var transactions []*schema.BridgeTransaction

	if request.Address != nil {
		transactions, err = databaseTransaction.FindBridgeTransactionsByAddress(c.Request().Context(), common.HexToAddress(*request.Address))
	} else {
		transactions, err = databaseTransaction.FindBridgeTransactions(c.Request().Context())
	}

	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find bridge transactions", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	ids := lo.Map(transactions, func(transaction *schema.BridgeTransaction, _ int) common.Hash {
		return transaction.ID
	})

	events, err := databaseTransaction.FindBridgeEventsByIDs(c.Request().Context(), ids)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find bridge events", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	if err := databaseTransaction.Commit(); err != nil {
		return fmt.Errorf("commit database transaction")
	}

	transactionModels := make([]*model.BridgeTransaction, 0, len(transactions))

	for _, transaction := range transactions {
		events := lo.Filter(events, func(event *schema.BridgeEvent, _ int) bool {
			return event.ID == transaction.ID
		})

		transactionModels = append(transactionModels, model.NewBridgeTransaction(transaction, events))
	}

	var response Response

	response.Data = transactionModels

	return c.JSON(http.StatusOK, response)
}

type GetBridgeTransactionRequest struct {
	ID *string `param:"id"`
}

func (h *Hub) GetBridgeTransaction(c echo.Context) error {
	var request GetBridgeTransactionRequest
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

	transaction, err := databaseTransaction.FindBridgeTransaction(c.Request().Context(), common.HexToHash(*request.ID))
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find bridge transaction", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	events, err := databaseTransaction.FindBridgeEventsByIDs(c.Request().Context(), []common.Hash{transaction.ID})
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find bridge events", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	if err := databaseTransaction.Commit(); err != nil {
		return fmt.Errorf("commit database transaction")
	}

	events = lo.Filter(events, func(event *schema.BridgeEvent, _ int) bool {
		return event.ID == transaction.ID
	})

	var response Response
	response.Data = model.NewBridgeTransaction(transaction, events)

	return c.JSON(http.StatusOK, response)
}
