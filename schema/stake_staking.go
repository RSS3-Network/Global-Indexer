package schema

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type StakeStakingExporter interface {
	Export() (*StakeStaking, error)
}

type StakeStakingTransformer interface {
	StakeStakingExporter
}

type StakeStaking struct {
	Staker common.Address    `json:"staker"`
	Node   common.Address    `json:"node"`
	Value  decimal.Decimal   `json:"value"`
	Chips  StakeStakingChips `json:"chips"`
}

type StakeStakingChips struct {
	Total    uint64       `json:"total"`
	Showcase []*StakeChip `json:"showcase"`
}

type StakeStakingsQuery struct {
	Cursor *string
	Node   *common.Address
	Staker *common.Address
	Limit  int
}
