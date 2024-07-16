package nta

import (
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type GetStakeTransactionsRequest struct {
	Cursor  *common.Hash                 `query:"cursor"`
	Staker  *common.Address              `query:"staker"`
	Node    *common.Address              `query:"node"`
	Type    *schema.StakeTransactionType `query:"type"`
	Pending *bool                        `query:"pending"`
	Limit   int                          `query:"limit" default:"20" min:"1" max:"20"`
}

type GetStakeTransactionRequest struct {
	TransactionHash *common.Hash                 `param:"transaction_hash"`
	Type            *schema.StakeTransactionType `query:"type"`
}

type GetStakeTransactionsResponseData []*StakeTransaction

type GetStakeTransactionResponseData *StakeTransaction

type StakeTransaction struct {
	ID        common.Hash                `json:"id"`
	Staker    common.Address             `json:"staker"`
	Node      common.Address             `json:"node"`
	Value     decimal.Decimal            `json:"value"`
	Chips     []*StakeChip               `json:"chips,omitempty"`
	Event     StakeTransactionEventTypes `json:"event"`
	Finalized bool                       `json:"finalized"`
}

type StakeTransactionEventTypes struct {
	Deposit  *StakeTransactionEventTypeDeposit  `json:"deposit,omitempty"`
	Withdraw *StakeTransactionEventTypeWithdraw `json:"withdraw,omitempty"`
	Stake    *StakeTransactionEventTypeStake    `json:"stake,omitempty"`
	Unstake  *StakeTransactionEventTypeUnstake  `json:"unstake,omitempty"`
}

type StakeTransactionEventTypeDeposit struct {
	Deposited *StakeTransactionEvent `json:"deposited,omitempty"`
}

type StakeTransactionEventTypeWithdraw struct {
	Requested *StakeTransactionEvent `json:"requested,omitempty"`
	Claimed   *StakeTransactionEvent `json:"claimed,omitempty"`
}

type StakeTransactionEventTypeStake struct {
	Staked *StakeTransactionEvent `json:"staked,omitempty"`
}

type StakeTransactionEventTypeUnstake struct {
	Requested *StakeTransactionEvent `json:"requested,omitempty"`
	Claimed   *StakeTransactionEvent `json:"claimed,omitempty"`
}

type StakeTransactionEvent struct {
	Block       TransactionEventBlock       `json:"block"`
	Transaction TransactionEventTransaction `json:"transaction"`
}

func NewStakeTransaction(transaction *schema.StakeTransaction, events []*schema.StakeEvent, stakeChips []*schema.StakeChip, baseURL url.URL) GetStakeTransactionResponseData {
	transactionModel := StakeTransaction{
		ID:     transaction.ID,
		Staker: transaction.User,
		Node:   transaction.Node,
		Value:  decimal.NewFromBigInt(transaction.Value, 0),
		Chips: lo.FilterMap(transaction.Chips, func(id *big.Int, _ int) (*StakeChip, bool) {
			stakeChip, found := lo.Find(stakeChips, func(stakeChip *schema.StakeChip) bool {
				return stakeChip.ID.Cmp(id) == 0
			})

			if !found {
				return nil, false
			}

			// Rewrite the owner address to restore the history.
			stakeChip.Owner = transaction.User

			return NewStakeChip(stakeChip, baseURL), true
		}),
		Finalized: transaction.Finalized,
	}

	switch transaction.Type {
	case schema.StakeTransactionTypeDeposit:
		transactionModel.Event.Deposit = new(StakeTransactionEventTypeDeposit)
	case schema.StakeTransactionTypeWithdraw:
		transactionModel.Event.Withdraw = new(StakeTransactionEventTypeWithdraw)
	case schema.StakeTransactionTypeStake:
		transactionModel.Event.Stake = new(StakeTransactionEventTypeStake)
	case schema.StakeTransactionTypeUnstake:
		transactionModel.Event.Unstake = new(StakeTransactionEventTypeUnstake)
	}

	for _, event := range events {
		if event.ID != transaction.ID {
			continue
		}

		eventModel := StakeTransactionEvent{
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
		case schema.StakeTransactionTypeDeposit:
			switch event.Type {
			case schema.StakeEventTypeDepositDeposited:
				transactionModel.Event.Deposit.Deposited = &eventModel
			}
		case schema.StakeTransactionTypeWithdraw:
			switch event.Type {
			case schema.StakeEventTypeWithdrawRequested:
				transactionModel.Event.Withdraw.Requested = &eventModel
			case schema.StakeEventTypeWithdrawClaimed:
				transactionModel.Event.Withdraw.Claimed = &eventModel
			}
		case schema.StakeTransactionTypeStake:
			switch event.Type {
			case schema.StakeEventTypeStakeStaked:
				transactionModel.Event.Stake.Staked = &eventModel
			}
		case schema.StakeTransactionTypeUnstake:
			switch event.Type {
			case schema.StakeEventTypeUnstakeRequested:
				transactionModel.Event.Unstake.Requested = &eventModel
			case schema.StakeEventTypeUnstakeClaimed:
				transactionModel.Event.Unstake.Claimed = &eventModel
			}
		}
	}

	return &transactionModel
}
