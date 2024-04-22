package schema

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type NodeFailureResponse struct {
	EpochID           uint64                    `json:"epochID"`
	Status            NodeFailureResponseStatus `json:"status"`
	ValidatorNode     common.Address            `json:"validatorNode"`
	ValidatorRequest  string                    `json:"validatorRequest"`
	ValidatorResponse json.RawMessage           `json:"validatorResponse"`
	VerifiedNode      common.Address            `json:"verifiedNode"`
	VerifiedRequest   string                    `json:"verifiedRequest"`
	VerifiedResponse  json.RawMessage           `json:"verifiedResponse"`
	CreatedAt         int64                     `json:"createdAt"`
}

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=NodeFailureResponseStatus --linecomment --output node_failure_response_status_string.go --json --yaml --sql
type NodeFailureResponseStatus int64

const (
	// NodeFailureResponseStatusChallengeable
	// A node is in this status, it is possible to initiate a challenge.
	// Possible reasons :
	// - Incorrect data submission by the node.
	// - Errors encountered during the processing of node requests.
	NodeFailureResponseStatusChallengeable NodeFailureResponseStatus = iota // challengeable
)
