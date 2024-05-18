package nta

import "github.com/ethereum/go-ethereum/common"

type PostNodeHideTaxRateRequest struct {
	NodeAddress common.Address `param:"node_address" validate:"required"`
	Signature   string         `json:"signature" validate:"required"`
}
