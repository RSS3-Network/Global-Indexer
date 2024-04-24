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
	FaultyNode        common.Address                 `gorm:"column:faulty_node"`
	FaultyResponse    json.RawMessage                `gorm:"column:faulty_response"`
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
	n.FaultyNode = nodeResponseFailure.FaultyNode
	n.FaultyResponse = nodeResponseFailure.FaultyResponse
}

func (n *NodeInvalidResponse) Export() *schema.NodeInvalidResponse {
	return &schema.NodeInvalidResponse{
		ID:                n.ID,
		EpochID:           n.EpochID,
		InvalidType:       n.InvalidType,
		Request:           n.Request,
		ValidatorNodes:    n.ValidatorNodes,
		ValidatorResponse: n.ValidatorResponse,
		FaultyNode:        n.FaultyNode,
		FaultyResponse:    n.FaultyResponse,
		CreatedAt:         n.CreatedAt.Unix(),
	}
}
