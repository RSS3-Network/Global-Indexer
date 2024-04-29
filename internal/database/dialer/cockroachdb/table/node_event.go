package table

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type NodeEvent struct {
	TransactionHash  string               `gorm:"column:transaction_hash;type:text;not null;primaryKey;"`
	TransactionIndex uint                 `gorm:"column:transaction_index;type:bigint;not null;primaryKey;autoIncrement:false;index:events_index_block_number,priority:2,sort:desc;"`
	NodeID           uint64               `gorm:"column:node_id;type:bigint;not null;index:events_index_node_id;"`
	AddressFrom      common.Address       `gorm:"column:address_from;type:bytea;not null;index:events_index_address,priority:1;index:events_index_address_type,priority:1;"`
	AddressTo        common.Address       `gorm:"column:address_to;type:bytea;not null;index:events_index_address,priority:2;:"`
	Type             schema.NodeEventType `gorm:"column:type;type:text;not null;index:events_index_address_type,priority:2;"`
	LogIndex         uint                 `gorm:"column:log_index;type:bigint;not null;primaryKey;autoIncrement:false;index:events_index_block_number,priority:3,sort:desc;"`
	ChainID          uint64               `gorm:"column:chain_id;type:bigint;not null;"`
	BlockHash        string               `gorm:"column:block_hash;type:text;not null;"`
	BlockNumber      uint64               `gorm:"column:block_number;type:bigint;not null;index:events_index_block_number,priority:1,sort:desc;"`
	BlockTimestamp   time.Time            `gorm:"column:block_timestamp;type:timestamp with time zone;not null;"`
	Metadata         json.RawMessage      `gorm:"column:metadata;type:jsonb;not null;"`
	CreatedAt        time.Time            `gorm:"column:created_at;type:timestamp with time zone;autoCreateTime;not null;default:now();"`
	UpdatedAt        time.Time            `gorm:"column:updated_at;type:timestamp with time zone;autoUpdateTime;not null;default:now();"`
}

func (*NodeEvent) TableName() string {
	return "node_events"
}

func (n *NodeEvent) Import(nodeEvent schema.NodeEvent) (err error) {
	n.TransactionHash = nodeEvent.TransactionHash.String()
	n.TransactionIndex = nodeEvent.TransactionIndex
	n.NodeID = nodeEvent.NodeID.Uint64()
	n.AddressFrom = nodeEvent.AddressFrom
	n.AddressTo = nodeEvent.AddressTo
	n.Type = nodeEvent.Type
	n.LogIndex = nodeEvent.LogIndex
	n.ChainID = nodeEvent.ChainID
	n.BlockHash = nodeEvent.BlockHash.String()
	n.BlockNumber = nodeEvent.BlockNumber.Uint64()
	n.BlockTimestamp = time.Unix(nodeEvent.BlockTimestamp, 0)

	n.Metadata, err = json.Marshal(nodeEvent.Metadata)
	if err != nil {
		return fmt.Errorf("marshal node event metadata: %w", err)
	}

	return nil
}

func (n *NodeEvent) Export() (*schema.NodeEvent, error) {
	nodeEvent := schema.NodeEvent{
		TransactionHash:  common.HexToHash(n.TransactionHash),
		TransactionIndex: n.TransactionIndex,
		NodeID:           big.NewInt(int64(n.NodeID)),
		AddressFrom:      n.AddressFrom,
		AddressTo:        n.AddressTo,
		Type:             n.Type,
		LogIndex:         n.LogIndex,
		ChainID:          n.ChainID,
		BlockHash:        common.HexToHash(n.BlockHash),
		BlockNumber:      big.NewInt(int64(n.BlockNumber)),
		BlockTimestamp:   n.BlockTimestamp.Unix(),
	}

	if err := json.Unmarshal(n.Metadata, &nodeEvent.Metadata); len(n.Metadata) > 0 && err != nil {
		return nil, fmt.Errorf("unmarshal node event metadata: %w", err)
	}

	return &nodeEvent, nil
}

type NodeEvents []*NodeEvent

func (n NodeEvents) Export() ([]*schema.NodeEvent, error) {
	nodeEvents := make([]*schema.NodeEvent, 0)

	for _, nodeEvent := range n {
		exported, err := nodeEvent.Export()
		if err != nil {
			return nil, fmt.Errorf("export node event: %w", err)
		}

		nodeEvents = append(nodeEvents, exported)
	}

	return nodeEvents, nil
}
