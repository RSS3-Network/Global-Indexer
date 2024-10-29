package table

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/rss3-network/global-indexer/schema"
)

type NodeInvalidResponse struct {
	ID               uint64                         `gorm:"id;primaryKey"`
	EpochID          uint64                         `gorm:"column:epoch_id"`
	Type             schema.NodeInvalidResponseType `gorm:"column:type"`
	Request          string                         `gorm:"column:request"`
	VerifierNodes    pq.ByteaArray                  `gorm:"column:verifier_nodes;type:bytea[]"`
	VerifierResponse json.RawMessage                `gorm:"column:verifier_response;type:jsonb"`
	Node             common.Address                 `gorm:"column:node"`
	Response         json.RawMessage                `gorm:"column:response;type:jsonb"`
	CreatedAt        time.Time                      `gorm:"column:created_at"`
	UpdatedAt        time.Time                      `gorm:"column:updated_at"`
}

func (*NodeInvalidResponse) TableName() string {
	return "node_invalid_response"
}

func (n *NodeInvalidResponse) Import(nodeInvalidResponse *schema.NodeInvalidResponse) {
	n.EpochID = nodeInvalidResponse.EpochID
	n.Type = nodeInvalidResponse.Type
	n.Request = nodeInvalidResponse.Request

	for _, verifierNode := range nodeInvalidResponse.VerifierNodes {
		n.VerifierNodes = append(n.VerifierNodes, verifierNode.Bytes())
	}

	n.VerifierResponse = nodeInvalidResponse.VerifierResponse
	n.Node = nodeInvalidResponse.Node
	n.Response = nodeInvalidResponse.Response
}

func (n *NodeInvalidResponse) Export() *schema.NodeInvalidResponse {
	var verifierNodes = make([]common.Address, len(n.VerifierNodes))

	for i, verifierNode := range n.VerifierNodes {
		verifierNodes[i] = common.BytesToAddress(verifierNode)
	}

	return &schema.NodeInvalidResponse{
		ID:               n.ID,
		EpochID:          n.EpochID,
		Type:             n.Type,
		Request:          n.Request,
		VerifierNodes:    verifierNodes,
		VerifierResponse: n.VerifierResponse,
		Node:             n.Node,
		Response:         n.Response,
		CreatedAt:        n.CreatedAt.Unix(),
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
