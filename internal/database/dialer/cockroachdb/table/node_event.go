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
	TransactionHash  string               `gorm:"column:transaction_hash"`
	TransactionIndex uint                 `gorm:"column:transaction_index"`
	NodeID           uint64               `gorm:"column:node_id"`
	AddressFrom      common.Address       `gorm:"column:address_from"`
	AddressTo        common.Address       `gorm:"column:address_to"`
	Type             schema.NodeEventType `gorm:"column:type"`
	LogIndex         uint                 `gorm:"column:log_index"`
	ChainID          uint64               `gorm:"column:chain_id"`
	BlockHash        string               `gorm:"column:block_hash"`
	BlockNumber      uint64               `gorm:"column:block_number"`
	BlockTimestamp   time.Time            `gorm:"column:block_timestamp"`
	Metadata         json.RawMessage      `gorm:"column:metadata"`
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
