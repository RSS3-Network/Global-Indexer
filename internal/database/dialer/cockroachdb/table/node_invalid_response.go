package table

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/rss3-network/global-indexer/schema"
)

type NodeInvalidResponse struct {
	ID                uint64                         `gorm:"id;primaryKey"`
	EpochID           uint64                         `gorm:"column:epoch_id"`
	Type              schema.NodeInvalidResponseType `gorm:"column:type"`
	Request           string                         `gorm:"column:request"`
	ValidatorNodes    pq.ByteaArray                  `gorm:"column:validator_nodes;type:bytea[]"`
	ValidatorResponse json.RawMessage                `gorm:"column:validator_response;type:jsonb"`
	Node              common.Address                 `gorm:"column:node"`
	Response          json.RawMessage                `gorm:"column:response;type:jsonb"`
	CreatedAt         time.Time                      `gorm:"column:created_at"`
	UpdatedAt         time.Time                      `gorm:"column:updated_at"`
}

func (*NodeInvalidResponse) TableName() string {
	return "node_invalid_response"
}

func (n *NodeInvalidResponse) Import(nodeResponseFailure *schema.NodeInvalidResponse) {
	n.EpochID = nodeResponseFailure.EpochID
	n.Type = nodeResponseFailure.Type
	n.Request = nodeResponseFailure.Request

	for _, validatorNode := range nodeResponseFailure.ValidatorNodes {
		n.ValidatorNodes = append(n.ValidatorNodes, validatorNode.Bytes())
	}

	n.ValidatorResponse = nodeResponseFailure.ValidatorResponse
	n.Node = nodeResponseFailure.Node
	n.Response = nodeResponseFailure.Response
}

func (n *NodeInvalidResponse) Export() *schema.NodeInvalidResponse {
	var validatorNodes = make([]common.Address, len(n.ValidatorNodes))

	for _, validatorNode := range n.ValidatorNodes {
		validatorNodes = append(validatorNodes, common.BytesToAddress(validatorNode))
	}

	return &schema.NodeInvalidResponse{
		ID:                n.ID,
		EpochID:           n.EpochID,
		Type:              n.Type,
		Request:           n.Request,
		ValidatorNodes:    validatorNodes,
		ValidatorResponse: n.ValidatorResponse,
		Node:              n.Node,
		Response:          n.Response,
		CreatedAt:         n.CreatedAt.Unix(),
	}
}

type NodeInvalidResponses []NodeInvalidResponse

func (ns *NodeInvalidResponses) Import(nodeInvalidResponses []*schema.NodeInvalidResponse) {
	*ns = make([]NodeInvalidResponse, 0, len(nodeInvalidResponses))

	for _, nodeInvalidResponse := range nodeInvalidResponses {
		var tNodeInvalidResponse NodeInvalidResponse

		tNodeInvalidResponse.Import(nodeInvalidResponse)

		*ns = append(*ns, tNodeInvalidResponse)
	}
}
