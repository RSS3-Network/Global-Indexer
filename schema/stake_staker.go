package schema

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type StakeStakerExporter interface {
	Export() (*StakeStaker, error)
}

type StakeStakerTransformer interface {
	StakeStakerExporter
}

type StakeStaker struct {
	Address common.Address   `json:"address"`
	Nodes   uint64           `json:"nodes"`
	Chips   StakeStakerChips `json:"chips"`
}

type StakeStakerChips struct {
	TotalNumber uint64          `json:"total_number"`
	TotalValue  decimal.Decimal `json:"total_value"`
}
