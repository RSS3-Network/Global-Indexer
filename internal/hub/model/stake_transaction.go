package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type StakeTransaction struct {
	ID     common.Hash                `json:"id"`
	Staker common.Address             `json:"staker"`
	Node   common.Address             `json:"node"`
	Value  decimal.Decimal            `json:"value"`
	Chips  []decimal.Decimal          `json:"chips,omitempty"`
	Event  StakeTransactionEventTypes `json:"event"`
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

func NewStakeTransaction(transaction *schema.StakeTransaction, events []*schema.StakeEvent) *StakeTransaction {
	transactionModel := StakeTransaction{
		ID:     transaction.ID,
		Staker: transaction.User,
		Node:   transaction.Node,
		Value:  decimal.NewFromBigInt(transaction.Value, 0),
		Chips: lo.Map(transaction.Chips, func(item *big.Int, _ int) decimal.Decimal {
			return decimal.NewFromBigInt(item, 0)
		}),
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
