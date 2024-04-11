package nta

import "github.com/ethereum/go-ethereum/common"

type NodeChallengeRequest struct {
	Address common.Address `param:"id" validate:"required"`
	Type    string         `query:"type"`
}

type NodeChallengeResponseData string
