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
	Address           common.Address  `json:"address"`
	TotalStakedNodes  uint64          `json:"total_staked_nodes"`
	TotalOwnedChips   uint64          `json:"total_owned_chips"`
	TotalStakedTokens decimal.Decimal `json:"total_staked_tokens"`
}
