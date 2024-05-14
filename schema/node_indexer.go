package schema

import "github.com/ethereum/go-ethereum/common"

type Worker struct {
	Address common.Address `json:"address"`
	Network string         `json:"network"`
	Name    string         `json:"name"`
}
