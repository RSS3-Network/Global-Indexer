package model

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

type NodeEvent struct {
	Transaction TransactionEventTransaction `json:"transaction"`
	Block       TransactionEventBlock       `json:"block"`
	AddressFrom common.Address              `json:"addressFrom"`
	AddressTo   common.Address              `json:"addressTo"`
	Type        schema.NodeEventType        `json:"type"`
	LogIndex    uint                        `json:"logIndex"`
	ChainID     uint64                      `json:"chainID"`
	Metadata    schema.NodeEventMetadata    `json:"metadata"`
}

func NewNodeEvent(event *schema.NodeEvent) *NodeEvent {
	return &NodeEvent{
		Transaction: TransactionEventTransaction{
			Hash:  event.TransactionHash,
			Index: event.TransactionIndex,
		},
		Block: TransactionEventBlock{
			Hash:      event.BlockHash,
			Number:    event.BlockNumber,
			Timestamp: event.BlockTimestamp,
		},
		AddressFrom: event.AddressFrom,
		AddressTo:   event.AddressTo,
		Type:        event.Type,
		LogIndex:    event.LogIndex,
		ChainID:     event.ChainID,
		Metadata:    event.Metadata,
	}
}

func NewNodeEvents(events []*schema.NodeEvent) []*NodeEvent {
	result := make([]*NodeEvent, len(events))
	for i, event := range events {
		result[i] = NewNodeEvent(event)
	}

	return result
}
