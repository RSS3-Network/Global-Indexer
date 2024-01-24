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
	s.Endpoint = stat.Endpoint
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
		Endpoint:             s.Endpoint,
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

type Stats []Stat

func (s *Stats) Export() ([]*schema.Stat, error) {
	stats := make([]*schema.Stat, 0)

	for _, stat := range *s {
		exportedStat, err := stat.Export()
		if err != nil {
			return nil, err
		}

		stats = append(stats, exportedStat)
	}

	return stats, nil
}

func (s *Stats) Import(stats []*schema.Stat) (err error) {
	*s = make([]Stat, 0, len(stats))

	for _, stat := range stats {
		var tStat Stat

		if err = tStat.Import(stat); err != nil {
			return err
		}

		*s = append(*s, tStat)
	}

	return nil
}

type Indexer struct {
	Address common.Address `gorm:"column:address;primaryKey"`
	Network string         `gorm:"column:network;primaryKey"`
	Worker  string         `gorm:"column:worker;primaryKey"`
}

func (*Indexer) TableName() string {
	return "node_indexer"
}

func (i *Indexer) Import(indexer *schema.Indexer) (err error) {
	i.Address = indexer.Address
	i.Network = indexer.Network
	i.Worker = indexer.Worker

	return nil
}

func (i *Indexer) Export() (*schema.Indexer, error) {
	indexer := schema.Indexer{
		Address: i.Address,
		Network: i.Network,
		Worker:  i.Worker,
	}

	return &indexer, nil
}

type Indexers []Indexer

func (i *Indexers) Export() ([]*schema.Indexer, error) {
	indexers := make([]*schema.Indexer, 0)

	for _, indexer := range *i {
		exportedIndexer, err := indexer.Export()
		if err != nil {
			return nil, err
		}

		indexers = append(indexers, exportedIndexer)
	}

	return indexers, nil
}

func (i *Indexers) Import(indexers []*schema.Indexer) (err error) {
	*i = make([]Indexer, 0, len(indexers))

	for _, indexer := range indexers {
		var tIndexer Indexer

		if err = tIndexer.Import(indexer); err != nil {
			return err
		}

		*i = append(*i, tIndexer)
	}

	return nil
}
