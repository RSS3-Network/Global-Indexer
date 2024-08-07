package nta

import "github.com/ethereum/go-ethereum/common"

type NodeChallengeRequest struct {
	NodeAddress common.Address `param:"node_address" validate:"required"`
	Type        string         `query:"type"`
}

type NodeChallengeResponseData string
