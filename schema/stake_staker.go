package schema

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type StakeStaker struct {
	Address             common.Address  `json:"address"`
	TotalStakedTokens   decimal.Decimal `json:"total_staked_tokens"`
	CurrentStakedNodes  uint64          `json:"current_staked_nodes"`
	CurrentOwnedChips   uint64          `json:"current_owned_chips"`
	CurrentStakedTokens decimal.Decimal `json:"current_staked_tokens"`
}
