package schema

import "github.com/ethereum/go-ethereum/common"

type Indexer struct {
	Address common.Address `json:"address"`
	Network string         `json:"network"`
	Worker  string         `json:"worker"`
}
