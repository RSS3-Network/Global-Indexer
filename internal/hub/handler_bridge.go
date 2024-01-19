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
	"github.com/shopspring/decimal"
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

	var transactions []*schema.BridgeTransaction

	if request.Address != nil {
		if transactions, err = databaseTransaction.FindBridgeTransactionsByAddress(c.Request().Context(), common.HexToAddress(*request.Address)); err != nil {
			zap.L().Error("find bridge transactions", zap.Error(err), zap.Any("request", request))

			return c.NoContent(http.StatusInternalServerError)
		}
	} else {
		if transactions, err = databaseTransaction.FindBridgeTransactions(c.Request().Context()); err != nil {
			zap.L().Error("find bridge transactions", zap.Error(err), zap.Any("request", request))

			return c.NoContent(http.StatusInternalServerError)
		}
	}

	ids := lo.Map(transactions, func(transaction *schema.BridgeTransaction, _ int) common.Hash {
		return transaction.ID
	})

	events, err := databaseTransaction.FindBridgeEventsByIDs(c.Request().Context(), ids)
	if err != nil {
		zap.L().Error("find bridge events", zap.Error(err), zap.Any("request", request))

		return c.NoContent(http.StatusInternalServerError)
	}

	if err := databaseTransaction.Commit(); err != nil {
		return fmt.Errorf("commit database transaction")
	}

	transactionModels := make([]model.BridgeTransaction, 0, len(transactions))

	for _, transaction := range transactions {
		transactionModel := model.BridgeTransaction{
			ID:       transaction.ID,
			Sender:   transaction.Sender,
			Receiver: transaction.Receiver,
			Token: model.BridgeToken{
				Address: model.BridgeTokenAddress{
					L1: transaction.TokenAddressL1,
					L2: transaction.TokenAddressL2,
				},
				Value: decimal.NewFromBigInt(transaction.TokenValue, 0),
			},
		}

		switch transaction.Type {
		case schema.BridgeTransactionTypeDeposit:
			transactionModel.Event.Deposit = new(model.BridgeTransactionEventTypeDeposit)
		case schema.BridgeTransactionTypeWithdraw:
			transactionModel.Event.Withdraw = new(model.BridgeTransactionEventTypeWithdraw)
		}

		for _, event := range events {
			if event.ID != transaction.ID {
				continue
			}

			eventModel := model.BridgeTransactionEvent{
				Block: model.BridgeTransactionEventBlock{
					Hash:      event.BlockHash,
					Number:    decimal.NewFromBigInt(event.BlockNumber, 0),
					Timestamp: event.BlockTimestamp,
				},
				Transaction: model.BridgeTransactionEventTransaction{
					Hash:  event.TransactionHash,
					Index: event.TransactionIndex,
				},
			}

			switch transaction.Type {
			case schema.BridgeTransactionTypeDeposit:
				switch event.Type {
				case schema.BridgeEventTypeDepositInitialized:
					transactionModel.Event.Deposit.Initialized = &eventModel
				case schema.BridgeEventTypeDepositFinalized:
					transactionModel.Event.Deposit.Finalized = &eventModel
				}
			case schema.BridgeTransactionTypeWithdraw:
				switch event.Type {
				case schema.BridgeEventTypeWithdrawalInitialized:
					transactionModel.Event.Withdraw.Initialized = &eventModel
				case schema.BridgeEventTypeWithdrawalProved:
					transactionModel.Event.Withdraw.Proved = &eventModel
				case schema.BridgeEventTypeWithdrawalFinalized:
					transactionModel.Event.Withdraw.Finalized = &eventModel
				}
			}
		}

		transactionModels = append(transactionModels, transactionModel)
	}

	var response Response

	response.Data = transactionModels

	return c.JSON(http.StatusOK, response)
}
