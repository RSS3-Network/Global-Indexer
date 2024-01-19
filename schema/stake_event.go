package schema

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type StakeEventType string

const (
	StakeEventTypeDepositDeposited StakeEventType = "deposited"

	StakeEventTypeWithdrawRequested StakeEventType = "requested"
	StakeEventTypeWithdrawClaimed   StakeEventType = "claimed"

	StakeEventTypeStakeStaked StakeEventType = "staked"

	StakeEventTypeUnstakeRequested StakeEventType = "requested"
	StakeEventTypeUnstakeClaimed   StakeEventType = "claimed"
)

type StakeEventImporter interface {
	Import(stakeEvent StakeEvent) error
}

type StakeEventExporter interface {
	Export() (*StakeEvent, error)
}

type StakeEventTransformer interface {
	StakeEventImporter
	StakeEventExporter
}

type StakeEvent struct {
	ID                common.Hash    `json:"id"`
	Type              StakeEventType `json:"type"`
	TransactionHash   common.Hash    `json:"transactionHash"`
	TransactionIndex  uint           `json:"transactionIndex"`
	TransactionStatus uint64         `json:"transactionStatus"`
	BlockHash         common.Hash    `json:"blockHash"`
	BlockNumber       *big.Int       `json:"blockNumber"`
	BlockTimestamp    time.Time      `json:"blockTimestamp"`
}
