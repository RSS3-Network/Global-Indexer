package table

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type NodeInvalidResponse struct {
	ID                uint64                         `gorm:"column:id;type:bigint;primaryKey;"`
	EpochID           uint64                         `gorm:"column:epoch_id;type:bigint;not null;index:idx_epoch_id,sort:desc"`
	Type              schema.NodeInvalidResponseType `gorm:"column:type;type:text;not null;index:idx_type,priority:1;"`
	Request           string                         `gorm:"column:request;type:text;not null;index:idx_request,priority:1;"`
	ValidatorNodes    []common.Address               `gorm:"column:validator_nodes;type:bytea[];"`
	ValidatorResponse json.RawMessage                `gorm:"column:validator_response;type:json;"`
	Node              common.Address                 `gorm:"column:node;type:bytea;index:idx_node,priority:1;"`
	Response          json.RawMessage                `gorm:"column:response;type:json;"`
	CreatedAt         time.Time                      `gorm:"column:created_at;type:timestamp with time zone;autoCreateTime;not null;default:now();index:idx_type,priority:2,sort:desc;index:idx_request,priority:2,sort:desc;index:idx_node,priority:2,sort:desc;"`
	UpdatedAt         time.Time                      `gorm:"column:updated_at;type:timestamp with time zone;autoUpdateTime;not null;default:now();"`
}

func (*NodeInvalidResponse) TableName() string {
	return "node_invalid_response"
}

func (n *NodeInvalidResponse) Import(nodeResponseFailure *schema.NodeInvalidResponse) {
	n.EpochID = nodeResponseFailure.EpochID
	n.Type = nodeResponseFailure.Type
	n.Request = nodeResponseFailure.Request
	n.ValidatorNodes = nodeResponseFailure.ValidatorNodes
	n.ValidatorResponse = nodeResponseFailure.ValidatorResponse
	n.Node = nodeResponseFailure.Node
	n.Response = nodeResponseFailure.Response
}

func (n *NodeInvalidResponse) Export() *schema.NodeInvalidResponse {
	return &schema.NodeInvalidResponse{
		ID:                n.ID,
		EpochID:           n.EpochID,
		Type:              n.Type,
		Request:           n.Request,
		ValidatorNodes:    n.ValidatorNodes,
		ValidatorResponse: n.ValidatorResponse,
		Node:              n.Node,
		Response:          n.Response,
		CreatedAt:         n.CreatedAt.Unix(),
	}
}
