package schema

import "github.com/ethereum/go-ethereum/common"

type Node struct {
	Address      common.Address `json:"address"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Endpoint     string         `json:"endpoint"`
	TaxFraction  uint64         `json:"taxFraction"`
	IsPublicGood bool           `json:"isPublicGood"`
	StreamURI    string         `json:"streamURI"`
}
