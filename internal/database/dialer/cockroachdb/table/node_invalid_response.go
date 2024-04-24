package table

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type NodeInvalidResponse struct {
	ID                uint64                         `gorm:"id;primaryKey"`
	EpochID           uint64                         `gorm:"column:epoch_id"`
	InvalidType       schema.NodeInvalidResponseType `gorm:"column:invalid_type"`
	Request           string                         `gorm:"column:request"`
	ValidatorNodes    []common.Address               `gorm:"column:validator_nodes"`
	ValidatorResponse json.RawMessage                `gorm:"column:validator_response"`
	Node              common.Address                 `gorm:"column:node"`
	InvalidResponse   json.RawMessage                `gorm:"column:invalid_response"`
	CreatedAt         time.Time                      `gorm:"column:created_at"`
	UpdatedAt         time.Time                      `gorm:"column:updated_at"`
}

func (*NodeInvalidResponse) TableName() string {
	return "node_invalid_response"
}

func (n *NodeInvalidResponse) Import(nodeResponseFailure *schema.NodeInvalidResponse) {
	n.EpochID = nodeResponseFailure.EpochID
	n.InvalidType = nodeResponseFailure.InvalidType
	n.Request = nodeResponseFailure.Request
	n.ValidatorNodes = nodeResponseFailure.ValidatorNodes
	n.ValidatorResponse = nodeResponseFailure.ValidatorResponse
	n.Node = nodeResponseFailure.Node
	n.InvalidResponse = nodeResponseFailure.InvalidResponse
}

func (n *NodeInvalidResponse) Export() *schema.NodeInvalidResponse {
	return &schema.NodeInvalidResponse{
		ID:                n.ID,
		EpochID:           n.EpochID,
		InvalidType:       n.InvalidType,
		Request:           n.Request,
		ValidatorNodes:    n.ValidatorNodes,
		ValidatorResponse: n.ValidatorResponse,
		Node:              n.Node,
		InvalidResponse:   n.InvalidResponse,
		CreatedAt:         n.CreatedAt.Unix(),
	}
}
