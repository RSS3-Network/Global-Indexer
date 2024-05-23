package nta

import (
	"database/sql"
	"errors"
	"fmt"
	"math/big"
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

	databaseTransactionOptions := sql.TxOptions{
		ReadOnly: true,
	}

	databaseTransaction, err := n.databaseClient.Begin(c.Request().Context(), &databaseTransactionOptions)
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

	stakeTransactionModels := make([]*nta.StakeTransaction, 0, len(stakeTransactions))

	for _, stakeTransaction := range stakeTransactions {
		stakeEvents := lo.Filter(stakeEvents, func(stakeEvent *schema.StakeEvent, _ int) bool {
			return stakeEvent.ID == stakeTransaction.ID
		})

		stakeTransactionModels = append(stakeTransactionModels, nta.NewStakeTransaction(stakeTransaction, stakeEvents, stakeChips, n.baseURL(c)))
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

	databaseTransactionOptions := sql.TxOptions{
		ReadOnly: true,
	}

	databaseTransaction, err := n.databaseClient.Begin(c.Request().Context(), &databaseTransactionOptions)
	if err != nil {
		zap.L().Error("begin database transaction", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	defer lo.Try(databaseTransaction.Rollback)

	stakeTransactionQuery := schema.StakeTransactionQuery{
		ID:   request.TransactionHash,
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

	var response nta.Response
	response.Data = nta.NewStakeTransaction(stakeTransaction, stakeEvents, stakeChips, n.baseURL(c))

	return c.JSON(http.StatusOK, response)
}
