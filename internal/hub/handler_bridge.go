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
	Cursor   *common.Hash                  `query:"cursor"`
	Sender   *common.Address               `query:"sender"`
	Receiver *common.Address               `query:"receiver"`
	Address  *common.Address               `query:"address"`
	Type     *schema.BridgeTransactionType `query:"type"`
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

	bridgeTransactionsQuery := schema.BridgeTransactionsQuery{
		Cursor:  request.Cursor,
		Address: request.Address,
		Type:    request.Type,
	}

	transactions, err := databaseTransaction.FindBridgeTransactions(c.Request().Context(), bridgeTransactionsQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find bridge transactions", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	bridgeEventsQuery := schema.BridgeEventsQuery{
		IDs: lo.Map(transactions, func(transaction *schema.BridgeTransaction, _ int) common.Hash {
			return transaction.ID
		}),
	}

	events, err := databaseTransaction.FindBridgeEvents(c.Request().Context(), bridgeEventsQuery)
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

	response := Response{
		Data: transactionModels,
	}

	if length := len(transactionModels); length > 0 {
		response.Cursor = transactionModels[length-1].ID.String()
	}

	return c.JSON(http.StatusOK, response)
}

type GetBridgeTransactionRequest struct {
	ID *common.Hash `param:"id"`
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

	bridgeTransactionQuery := schema.BridgeTransactionQuery{
		ID: request.ID,
	}

	bridgeTransaction, err := databaseTransaction.FindBridgeTransaction(c.Request().Context(), bridgeTransactionQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find bridge transaction", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	bridgeEventsQuery := schema.BridgeEventsQuery{
		IDs: []common.Hash{
			bridgeTransaction.ID,
		},
	}

	bridgeEvents, err := databaseTransaction.FindBridgeEvents(c.Request().Context(), bridgeEventsQuery)
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

	bridgeEvents = lo.Filter(bridgeEvents, func(bridgeEvent *schema.BridgeEvent, _ int) bool {
		return bridgeEvent.ID == bridgeTransaction.ID
	})

	var response Response
	response.Data = model.NewBridgeTransaction(bridgeTransaction, bridgeEvents)

	return c.JSON(http.StatusOK, response)
}
