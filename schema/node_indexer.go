package schema

import "github.com/ethereum/go-ethereum/common"

type Worker struct {
	EpochID  uint64         `json:"epoch_id"`
	Address  common.Address `json:"address"`
	Network  string         `json:"network"`
	Name     string         `json:"name"`
	IsActive bool           `json:"is_active"`
}

type WorkerQuery struct {
	NodeAddresses []common.Address
	Networks      []string
	Names         []string
	EpochID       uint64
	IsActive      *bool
}
