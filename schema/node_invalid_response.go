package schema

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

// NodeInvalidResponse records an alleged invalid response of a Node
// A group of Nodes are selected as verifiers to verify the response returned by a Node
// The response (with all responses from all verifiers) is saved in the database pending challenge by the penalized Node
type NodeInvalidResponse struct {
	ID               uint64                  `json:"id"`
	EpochID          uint64                  `json:"epoch_id"`
	Type             NodeInvalidResponseType `json:"type"`
	Request          string                  `json:"request"`
	VerifierNodes    []common.Address        `json:"verifier_nodes"`
	VerifierResponse json.RawMessage         `json:"verifier_response"`
	Node             common.Address          `json:"node"`
	Response         json.RawMessage         `json:"response"`
	CreatedAt        int64                   `json:"created_at"`
}

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=NodeInvalidResponseType --linecomment --output node_invalid_response_type_string.go --json --yaml --sql
type NodeInvalidResponseType int64

const (
	// NodeInvalidResponseTypeInconsistent when the Node's response differs from the majority of verifiers
	NodeInvalidResponseTypeInconsistent NodeInvalidResponseType = iota // inconsistent
	// NodeInvalidResponseTypeError when the Node returns an error
	NodeInvalidResponseTypeError // error
)
