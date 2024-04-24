package schema

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

// NodeInvalidResponse records an alleged invalid response of a Node
// A group of Nodes are selected as validators to verify the response returned by a Node
// The response (with all responses from all validators) is saved in the database pending challenge by the penalized Node
type NodeInvalidResponse struct {
	ID                uint64                  `json:"id"`
	EpochID           uint64                  `json:"epochID"`
	InvalidType       NodeInvalidResponseType `json:"invalidType"`
	Request           string                  `json:"request"`
	ValidatorNodes    []common.Address        `json:"validatorNodes"`
	ValidatorResponse json.RawMessage         `json:"validatorResponse"`
	Node              common.Address          `json:"ode"`
	InvalidResponse   json.RawMessage         `json:"invalidResponse"`
	CreatedAt         int64                   `json:"createdAt"`
}

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=NodeInvalidResponseType --linecomment --output node_invalid_response_type_string.go --json --yaml --sql
type NodeInvalidResponseType int64

const (
	// NodeInvalidResponseTypeInconsistent when the Node's response differs from the majority of validators
	NodeInvalidResponseTypeInconsistent NodeInvalidResponseType = iota // inconsistent
	// NodeInvalidResponseTypeError when the Node returns an error
	NodeInvalidResponseTypeError // error
)
