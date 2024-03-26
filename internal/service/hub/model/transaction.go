package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransactionEventBlock struct {
	Hash      common.Hash `json:"hash"`
	Number    *big.Int    `json:"number"`
	Timestamp int64       `json:"timestamp"`
}

type TransactionEventTransaction struct {
	Hash  common.Hash `json:"hash"`
	Index uint        `json:"index"`
}
