package schema

import "github.com/ethereum/go-ethereum/common"

type StakeStaking struct {
	Staker common.Address    `json:"staker,omitempty"`
	Node   common.Address    `json:"node,omitempty"`
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
