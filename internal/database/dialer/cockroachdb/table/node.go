package table

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type Node struct {
	Address                common.Address    `gorm:"column:address;primaryKey"`
	NodeID                 uint64            `gorm:"column:id"`
	Endpoint               string            `gorm:"column:endpoint"`
	HideTaxRate            bool              `gorm:"column:hide_tax_rate"`
	IsPublicGood           bool              `gorm:"column:is_public_good"`
	Stream                 json.RawMessage   `gorm:"column:stream"`
	Config                 json.RawMessage   `gorm:"column:config;type:jsonb"`
	Status                 schema.NodeStatus `gorm:"column:status"`
	LastHeartbeatTimestamp time.Time         `gorm:"column:last_heartbeat_timestamp"`
	// TODO: rename column to Location in database once atlas is merged
	Location         json.RawMessage `gorm:"column:local;type:jsonb"`
	Avatar           json.RawMessage `gorm:"column:avatar;type:jsonb"`
	MinTokensToStake decimal.Decimal `gorm:"column:min_tokens_to_stake"`
	APY              decimal.Decimal `gorm:"column:apy"`
	Score            decimal.Decimal `gorm:"column:score"`
	CreatedAt        time.Time       `gorm:"column:created_at"`
	UpdatedAt        time.Time       `gorm:"column:updated_at"`
}

func (*Node) TableName() string {
	return "node_info"
}

func (n *Node) Import(node *schema.Node) (err error) {
	n.Address = node.Address
	n.NodeID = node.ID.Uint64()
	n.Endpoint = node.Endpoint
	n.HideTaxRate = node.HideTaxRate
	n.IsPublicGood = node.IsPublicGood
	n.Status = node.Status
	n.LastHeartbeatTimestamp = time.Unix(node.LastHeartbeatTimestamp, 0)
	n.Stream = node.Stream
	n.Config = node.Config
	n.MinTokensToStake = node.MinTokensToStake
	n.APY = node.APY
	n.Score = node.ActiveScore

	n.Location, err = json.Marshal(node.Location)
	if err != nil {
		return fmt.Errorf("marshal node local: %w", err)
	}

	n.Avatar, err = json.Marshal(node.Avatar)
	if err != nil {
		return fmt.Errorf("marshal node avatar: %w", err)
	}

	return nil
}

func (n *Node) Export() (*schema.Node, error) {
	locations := make([]*schema.NodeLocation, 0)

	if err := json.Unmarshal(n.Location, &locations); len(n.Location) > 0 && err != nil {
		return nil, fmt.Errorf("unmarshal node locations: %w", err)
	}

	var avatar *l2.ChipsTokenMetadata
	if err := json.Unmarshal(n.Avatar, &avatar); len(n.Avatar) > 0 && err != nil {
		return nil, fmt.Errorf("unmarshal node avatar: %w", err)
	}

	return &schema.Node{
		Address:                n.Address,
		ID:                     big.NewInt(int64(n.NodeID)),
		Endpoint:               n.Endpoint,
		HideTaxRate:            n.HideTaxRate,
		IsPublicGood:           n.IsPublicGood,
		Status:                 n.Status,
		LastHeartbeatTimestamp: n.LastHeartbeatTimestamp.Unix(),
		Stream:                 n.Stream,
		Config:                 n.Config,
		Location:               locations,
		Avatar:                 avatar,
		MinTokensToStake:       n.MinTokensToStake,
		APY:                    n.APY,
		ActiveScore:            n.Score,
		CreatedAt:              n.CreatedAt.Unix(),
	}, nil
}

type Nodes []Node

func (n *Nodes) Import(nodes []*schema.Node) (err error) {
	*n = make([]Node, 0, len(nodes))

	for _, node := range nodes {
		var tNode Node

		if err = tNode.Import(node); err != nil {
			return err
		}

		*n = append(*n, tNode)
	}

	return nil
}

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
