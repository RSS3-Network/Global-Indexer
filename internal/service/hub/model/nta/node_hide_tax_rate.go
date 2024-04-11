package nta

import "github.com/ethereum/go-ethereum/common"

type NodeHideTaxRateRequest struct {
	Address   common.Address `param:"id" validate:"required"`
	Signature string         `json:"signature" validate:"required"`
}
