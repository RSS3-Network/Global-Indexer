package model

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type BridgeTransaction struct {
	ID       common.Hash                 `json:"id"`
	Sender   common.Address              `json:"sender"`
	Receiver common.Address              `json:"receiver"`
	Token    BridgeToken                 `json:"token"`
	Event    BridgeTransactionEventTypes `json:"event"`
}

type BridgeTransactionEventTypes struct {
	Deposit  *BridgeTransactionEventTypeDeposit  `json:"deposit,omitempty"`
	Withdraw *BridgeTransactionEventTypeWithdraw `json:"withdraw,omitempty"`
}

type BridgeTransactionEventTypeDeposit struct {
	Initialized *BridgeTransactionEvent `json:"initialized,omitempty"`
	Finalized   *BridgeTransactionEvent `json:"finalized,omitempty"`
}

type BridgeTransactionEventTypeWithdraw struct {
	Initialized *BridgeTransactionEvent `json:"initialized,omitempty"`
	Proved      *BridgeTransactionEvent `json:"proved,omitempty"`
	Finalized   *BridgeTransactionEvent `json:"finalized,omitempty"`
}

type BridgeTransactionEvent struct {
	Block       TransactionEventBlock       `json:"block"`
	Transaction TransactionEventTransaction `json:"transaction"`
}

func NewBridgeTransaction(transaction *schema.BridgeTransaction, events []*schema.BridgeEvent) *BridgeTransaction {
	transactionModel := BridgeTransaction{
		ID:       transaction.ID,
		Sender:   transaction.Sender,
		Receiver: transaction.Receiver,
		Token: BridgeToken{
			Address: BridgeTokenAddress{
				L1: transaction.TokenAddressL1,
				L2: transaction.TokenAddressL2,
			},
			Value: decimal.NewFromBigInt(transaction.TokenValue, 0),
		},
	}

	switch transaction.Type {
	case schema.BridgeTransactionTypeDeposit:
		transactionModel.Event.Deposit = new(BridgeTransactionEventTypeDeposit)
	case schema.BridgeTransactionTypeWithdraw:
		transactionModel.Event.Withdraw = new(BridgeTransactionEventTypeWithdraw)
	}

	for _, event := range events {
		if event.ID != transaction.ID {
			continue
		}

		eventModel := BridgeTransactionEvent{
			Block: TransactionEventBlock{
				Hash:      event.BlockHash,
				Number:    event.BlockNumber,
				Timestamp: event.BlockTimestamp.Unix(),
			},
			Transaction: TransactionEventTransaction{
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

	return &transactionModel
}
