package table

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/global-indexer/schema"
)

type Node struct {
	Address      common.Address  `gorm:"column:address;primaryKey"`
	Endpoint     string          `gorm:"column:endpoint"`
	IsPublicGood bool            `gorm:"column:is_public_good"`
	Stream       json.RawMessage `gorm:"column:stream"`
	Config       json.RawMessage `gorm:"column:config;type:jsonb"`
}

func (*Node) TableName() string {
	return "node_info"
}

func (n *Node) Import(node *schema.Node) (err error) {
	n.Address = node.Address
	n.Endpoint = node.Endpoint
	n.IsPublicGood = node.IsPublicGood

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
		Address:      n.Address,
		Endpoint:     n.Endpoint,
		IsPublicGood: n.IsPublicGood,
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

type Stat struct {
	Address              common.Address `gorm:"column:address;primaryKey"`
	Endpoint             string         `gorm:"column:endpoint"`
	Points               float64        `gorm:"column:points"`
	IsPublicGood         bool           `gorm:"column:is_public_good"`
	IsFullNode           bool           `gorm:"column:is_full_node"`
	IsRssNode            bool           `gorm:"column:is_rss_node"`
	Staking              float64        `gorm:"column:staking"`
	TotalRequest         int64          `gorm:"column:total_request_count"`
	EpochRequest         int64          `gorm:"column:epoch_request_count"`
	EpochInvalidRequest  int64          `gorm:"column:epoch_invalid_request_count"`
	DecentralizedNetwork int            `gorm:"column:decentralized_network_count"`
	FederatedNetwork     int            `gorm:"column:federated_network_count"`
	Indexer              int            `gorm:"column:indexer_count"`
	ResetAt              time.Time      `gorm:"column:reset_at"`
}

func (*Stat) TableName() string {
	return "node_stat"
}

func (s *Stat) Import(stat *schema.Stat) (err error) {
	s.Address = stat.Address
	s.Points = stat.Points
	s.IsPublicGood = stat.IsPublicGood
	s.IsFullNode = stat.IsFullNode
	s.IsRssNode = stat.IsRssNode
	s.Staking = stat.Staking
	s.TotalRequest = stat.TotalRequest
	s.EpochRequest = stat.EpochRequest
	s.EpochInvalidRequest = stat.EpochInvalidRequest
	s.DecentralizedNetwork = stat.DecentralizedNetwork
	s.FederatedNetwork = stat.FederatedNetwork
	s.Indexer = stat.Indexer
	s.ResetAt = stat.ResetAt

	return nil
}

func (s *Stat) Export() (*schema.Stat, error) {
	stat := schema.Stat{
		Address:              s.Address,
		Points:               s.Points,
		IsPublicGood:         s.IsPublicGood,
		IsFullNode:           s.IsFullNode,
		IsRssNode:            s.IsRssNode,
		Staking:              s.Staking,
		TotalRequest:         s.TotalRequest,
		EpochRequest:         s.EpochRequest,
		EpochInvalidRequest:  s.EpochInvalidRequest,
		DecentralizedNetwork: s.DecentralizedNetwork,
		FederatedNetwork:     s.FederatedNetwork,
		Indexer:              s.Indexer,
		ResetAt:              s.ResetAt,
	}

	return &stat, nil
}

type Stats []*Stat

func (s Stats) Export() ([]*schema.Stat, error) {
	stats := make([]*schema.Stat, 0)

	for _, stat := range s {
		exportedStat, err := stat.Export()
		if err != nil {
			return nil, err
		}

		stats = append(stats, exportedStat)
	}

	return stats, nil
}
