package schema

import "github.com/ethereum/go-ethereum/common"

type StakeAddress struct {
	Address common.Address    `json:"address"`
	Chips   *StakeAddressChip `json:"chips"`
}

type StakeAddressChip struct {
	Total    int64        `json:"total"`
	Showcase []*StakeChip `json:"showcase"`
}

type StakeNodeUsersQuery struct {
	Cursor *common.Address
	Node   *common.Address
	Limit  int
}

type StakeUserNodesQuery struct {
	Cursor *common.Address
	Owner  *common.Address
	Limit  int
}
