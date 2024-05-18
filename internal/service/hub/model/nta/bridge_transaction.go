package nta

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type GetBridgeTransactionsRequest struct {
	Cursor   *common.Hash                  `query:"cursor" description:"cursor for pagination"`
	Sender   *common.Address               `query:"sender" description:"sender address"`
	Receiver *common.Address               `query:"receiver" description:"receiver address"`
	Address  *common.Address               `query:"address" description:"token address"`
	Type     *schema.BridgeTransactionType `query:"type" description:"transaction type"`
	Limit    int                           `query:"limit" default:"20" min:"1" max:"20" description:"limit the number of results"`
}

type GetBridgeTransactionRequest struct {
	TransactionHash *common.Hash `param:"transaction_hash" description:"transaction hash"`
}

type GetBridgeTransactionsResponseData []*BridgeTransaction

type GetBridgeTransactionResponseData *BridgeTransaction

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

type BridgeToken struct {
	Address BridgeTokenAddress `json:"address"`
	Value   decimal.Decimal    `json:"value"`
}

type BridgeTokenAddress struct {
	L1 *common.Address `json:"l1,omitempty"`
	L2 *common.Address `json:"l2,omitempty"`
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
