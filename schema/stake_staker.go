package schema

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type StakeStakerImporter interface {
	Import(stakeStaker StakeStaker) error
}

type StakeStakerExporter interface {
	Export() (*StakeStaker, error)
}

type StakeStakerTransformer interface {
	StakeStakerImporter
	StakeStakerExporter
}

type StakeStaker struct {
	User  common.Address  `json:"user"`
	Node  common.Address  `json:"node"`
	Value decimal.Decimal `json:"value"`
}
