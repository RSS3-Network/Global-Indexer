package nta

import (
	"errors"
	"net/http"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (n *NTA) GetStakeTransactions(c echo.Context) error {
	var request nta.GetStakeTransactionsRequest
	if err := c.Bind(&request); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.InternalError(c)
	}

	stakeTransactionsQuery := schema.StakeTransactionsQuery{
		Cursor:  request.Cursor,
		User:    request.Staker,
		Node:    request.Node,
		Type:    request.Type,
		Pending: request.Pending,
		Limit:   request.Limit,
	}

	// Find staking transactions
	stakeTransactions, err := n.databaseClient.FindStakeTransactions(c.Request().Context(), stakeTransactionsQuery)
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

	// Find staking events
	stakeEvents, err := n.databaseClient.FindStakeEvents(c.Request().Context(), stakeEventsQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find stake events", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	stakeTransactionModels := make([]*nta.StakeTransaction, 0, len(stakeTransactions))

	for _, stakeTransaction := range stakeTransactions {
		stakeEvents := lo.Filter(stakeEvents, func(stakeEvent *schema.StakeEvent, _ int) bool {
			return stakeEvent.ID == stakeTransaction.ID
		})

		stakeTransactionModels = append(stakeTransactionModels, nta.NewStakeTransaction(stakeTransaction, stakeEvents))
	}

	response := nta.Response{
		Data: stakeTransactionModels,
	}

	if length := len(stakeTransactionModels); length > 0 && length == request.Limit {
		response.Cursor = stakeTransactionModels[length-1].ID.String()
	}

	return c.JSON(http.StatusOK, response)
}

func (n *NTA) GetStakeTransaction(c echo.Context) error {
	var request nta.GetStakeTransactionRequest
	if err := c.Bind(&request); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	stakeTransactionQuery := schema.StakeTransactionQuery{
		ID:   request.TransactionHash,
		Type: request.Type,
	}

	stakeTransaction, err := n.databaseClient.FindStakeTransaction(c.Request().Context(), stakeTransactionQuery)
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

	stakeEvents, err := n.databaseClient.FindStakeEvents(c.Request().Context(), stakeEventsQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("find stake events", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	stakeEvents = lo.Filter(stakeEvents, func(stakeEvent *schema.StakeEvent, _ int) bool {
		return stakeEvent.ID == stakeTransaction.ID
	})

	var response nta.Response
	response.Data = nta.NewStakeTransaction(stakeTransaction, stakeEvents)

	return c.JSON(http.StatusOK, response)
}
