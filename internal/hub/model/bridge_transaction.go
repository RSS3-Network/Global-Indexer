package model

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	Block       BridgeTransactionEventBlock       `json:"block"`
	Transaction BridgeTransactionEventTransaction `json:"transaction"`
}

type BridgeTransactionEventBlock struct {
	Hash      common.Hash     `json:"hash"`
	Number    decimal.Decimal `json:"number"`
	Timestamp time.Time       `json:"timestamp"`
}

type BridgeTransactionEventTransaction struct {
	Hash  common.Hash `json:"hash"`
	Index uint        `json:"index"`
}
