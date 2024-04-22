package schema

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type NodeInvalidResponse struct {
	ID                uint64                    `json:"id"`
	EpochID           uint64                    `json:"epochID"`
	Status            NodeInvalidResponseStatus `json:"status"`
	ValidatorNode     common.Address            `json:"validatorNode"`
	ValidatorRequest  string                    `json:"validatorRequest"`
	ValidatorResponse json.RawMessage           `json:"validatorResponse"`
	VerifiedNode      common.Address            `json:"verifiedNode"`
	VerifiedRequest   string                    `json:"verifiedRequest"`
	VerifiedResponse  json.RawMessage           `json:"verifiedResponse"`
	CreatedAt         int64                     `json:"createdAt"`
}

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=NodeInvalidResponseStatus --linecomment --output node_invalid_response_status_string.go --json --yaml --sql
type NodeInvalidResponseStatus int64

const (
	// NodeInvalidResponseStatusChallengeable
	// A node is in this status, it is possible to initiate a challenge.
	// Possible reasons :
	// - Incorrect data submission by the node.
	// - Errors encountered during the processing of node requests.
	NodeInvalidResponseStatusChallengeable NodeInvalidResponseStatus = iota // challengeable
)
