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
	Address                common.Address    `gorm:"column:address;type:bytea;primaryKey;"`
	NodeID                 uint64            `gorm:"column:id;type:bigint;not null;index:idx_id,unique;"`
	Endpoint               string            `gorm:"column:endpoint;type:text;not null;index:idx_endpoint_unique,unique;"`
	HideTaxRate            bool              `gorm:"column:hide_tax_rate;type:bool;default:false;"`
	IsPublicGood           bool              `gorm:"column:is_public_good;type:bool;not null;index:idx_is_public,priority:1;"`
	Stream                 json.RawMessage   `gorm:"column:stream;type:jsonb"`
	Config                 json.RawMessage   `gorm:"column:config;type:jsonb"`
	Status                 schema.NodeStatus `gorm:"column:status;type:text;not null;default:'offline';index:idx_status;"`
	LastHeartbeatTimestamp time.Time         `gorm:"column:last_heartbeat_timestamp;type:timestamp with time zone;index:idx_last_heartbeat_timestamp;"`
	Location               json.RawMessage   `gorm:"column:location;type:jsonb;not null;default:'[]'"`
	Avatar                 json.RawMessage   `gorm:"column:avatar;type:jsonb"`
	MinTokensToStake       decimal.Decimal   `gorm:"column:min_tokens_to_stake;type:decimal;"`
	APY                    decimal.Decimal   `gorm:"column:apy;type:decimal;default:0;"`
	Score                  decimal.Decimal   `gorm:"column:score;type:decimal;default:0;index:idx_score,sort:desc;"`
	CreatedAt              time.Time         `gorm:"column:created_at;type:timestamp with time zone;autoCreateTime;not null;default:now();index:idx_is_public,priority:2,sort:desc;index:idx_created_at,sort:desc;"`
	UpdatedAt              time.Time         `gorm:"column:updated_at;type:timestamp with time zone;autoUpdateTime;not null;default:now();"`
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
