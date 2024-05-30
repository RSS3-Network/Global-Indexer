package schema

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type EpochTrigger struct {
	TransactionHash common.Hash    `json:"transaction_hash"`
	EpochID         uint64         `json:"epoch_id"`
	Data            SettlementData `json:"data"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type SettlementData struct {
	Epoch            *big.Int         `json:"epoch"`
	NodeAddress      []common.Address `json:"node_addresses"`
	OperationRewards []*big.Int       `json:"operation_rewards"`
	RequestCount     []*big.Int       `json:"request_count"`
	IsFinal          bool             `json:"is_final"`
}
