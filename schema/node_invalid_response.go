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
	FaultyNode        common.Address          `json:"faultyNode"`
	FaultyResponse    json.RawMessage         `json:"faultyResponse"`
	CreatedAt         int64                   `json:"createdAt"`
}

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=NodeInvalidResponseType --linecomment --output node_invalid_response_type_string.go --json --yaml --sql
type NodeInvalidResponseType int64

const (
	NodeInvalidResponseTypeData  NodeInvalidResponseType = iota // data
	NodeInvalidResponseTypeError                                // error
)
