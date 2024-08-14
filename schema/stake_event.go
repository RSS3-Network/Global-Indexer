package schema

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type StakeEventType string

const (
	StakeEventTypeDepositDeposited StakeEventType = "deposited"

	StakeEventTypeWithdrawRequested StakeEventType = "withdrawRequested"
	StakeEventTypeWithdrawClaimed   StakeEventType = "withdrawClaimed"

	StakeEventTypeStakeStaked StakeEventType = "staked"

	StakeEventTypeChipsMerged = "merged"
	StakeEventTypeChipsBurned = "burned"

	StakeEventTypeUnstakeRequested StakeEventType = "unstakeRequested"
	StakeEventTypeUnstakeClaimed   StakeEventType = "unstakeClaimed"
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
	ID                common.Hash     `json:"id"`
	Type              StakeEventType  `json:"type"`
	TransactionHash   common.Hash     `json:"transaction_hash"`
	TransactionIndex  uint            `json:"transaction_index"`
	TransactionStatus uint64          `json:"transaction_status"`
	LogIndex          uint            `json:"log_index"`
	Metadata          json.RawMessage `json:"metadata"`
	BlockHash         common.Hash     `json:"block_hash"`
	BlockNumber       *big.Int        `json:"block_number"`
	BlockTimestamp    time.Time       `json:"block_timestamp"`
	Finalized         bool
}

type StakeEventQuery struct {
	ID *common.Hash `query:"id"`
}

type StakeEventsQuery struct {
	IDs []common.Hash `query:"ids"`
}

type StakeEventChipsMergedMetadata struct {
	ChipID *big.Int `json:"chip_id"`
}
