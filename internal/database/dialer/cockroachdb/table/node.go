package table

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

type Node struct {
	Address                common.Address  `gorm:"column:address;primaryKey"`
	Endpoint               string          `gorm:"column:endpoint"`
	IsPublicGood           bool            `gorm:"column:is_public_good"`
	Stream                 json.RawMessage `gorm:"column:stream"`
	Config                 json.RawMessage `gorm:"column:config;type:jsonb"`
	Status                 schema.Status   `gorm:"column:status"`
	LastHeartbeatTimestamp time.Time       `gorm:"column:last_heartbeat_timestamp"`
	CreatedAt              time.Time       `gorm:"column:created_at"`
	UpdatedAt              time.Time       `gorm:"column:updated_at"`
}

func (*Node) TableName() string {
	return "node_info"
}

func (n *Node) Import(node *schema.Node) (err error) {
	n.Address = node.Address
	n.Endpoint = node.Endpoint
	n.IsPublicGood = node.IsPublicGood
	n.Status = node.Status
	n.LastHeartbeatTimestamp = time.Unix(node.LastHeartbeatTimestamp, 0)

	if n.Stream, err = json.Marshal(node.Stream); err != nil {
		return fmt.Errorf("failed to marshal node stream: %w", err)
	}

	if n.Config, err = json.Marshal(node.Config); err != nil {
		return fmt.Errorf("failed to marshal node config: %w", err)
	}

	return nil
}

func (n *Node) Export() (*schema.Node, error) {
	node := schema.Node{
		Address:                n.Address,
		Endpoint:               n.Endpoint,
		IsPublicGood:           n.IsPublicGood,
		Status:                 n.Status,
		LastHeartbeatTimestamp: n.LastHeartbeatTimestamp.Unix(),
	}

	if err := json.Unmarshal(n.Stream, &node.Stream); err != nil {
		return nil, fmt.Errorf("failed to unmarshal node stream: %w", err)
	}

	if err := json.Unmarshal(n.Config, &node.Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal node config: %w", err)
	}

	return &node, nil
}

type Nodes []*Node

func (n Nodes) Export() ([]*schema.Node, error) {
	nodes := make([]*schema.Node, 0)

	for _, node := range n {
		exportedNode, err := node.Export()
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, exportedNode)
	}

	return nodes, nil
}
