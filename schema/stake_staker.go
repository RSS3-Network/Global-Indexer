package schema

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type StakeStaker struct {
	Address             common.Address  `json:"address"`
	TotalStakedNodes    uint64          `json:"total_staked_nodes"`
	TotalChips          uint64          `json:"total_chips"`
	TotalStakedTokens   decimal.Decimal `json:"total_staked_tokens"`
	CurrentStakedTokens decimal.Decimal `json:"current_staked_tokens"`
}
