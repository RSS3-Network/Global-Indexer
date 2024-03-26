package model

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type BridgeToken struct {
	Address BridgeTokenAddress `json:"address"`
	Value   decimal.Decimal    `json:"value"`
}

type BridgeTokenAddress struct {
	L1 *common.Address `json:"l1,omitempty"`
	L2 *common.Address `json:"l2,omitempty"`
}
