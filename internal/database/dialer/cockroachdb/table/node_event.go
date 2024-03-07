package table

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

type NodeEvent struct {
	TransactionHash  string               `gorm:"transaction_hash"`
	TransactionIndex uint                 `gorm:"transaction_index"`
	AddressFrom      common.Address       `gorm:"address_from"`
	AddressTo        common.Address       `gorm:"address_to"`
	Type             schema.NodeEventType `gorm:"type"`
	LogIndex         uint                 `gorm:"log_index"`
	ChainID          uint64               `gorm:"chain_id"`
	BlockHash        string               `gorm:"block_hash"`
	BlockNumber      uint64               `gorm:"block_number"`
	BlockTimestamp   time.Time            `gorm:"block_timestamp"`
	Metadata         json.RawMessage      `gorm:"metadata"`
}

func (*NodeEvent) TableName() string {
	return "node.events"
}

func (n *NodeEvent) Import(nodeEvent schema.NodeEvent) (err error) {
	n.TransactionHash = nodeEvent.TransactionHash.String()
	n.TransactionIndex = nodeEvent.TransactionIndex
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
