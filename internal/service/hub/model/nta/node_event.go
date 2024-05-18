package nta

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type GetNodeEventsRequest struct {
	NodeAddress common.Address `param:"node_address" validate:"required"`
	Cursor      *string        `query:"cursor"`
	Limit       int            `query:"limit" validate:"min=1,max=100" default:"20"`
}

type GetNodeEventResponseData *NodeEvent

type GetNodeEventsResponseData []*NodeEvent

type NodeEvent struct {
	Transaction TransactionEventTransaction `json:"transaction"`
	Block       TransactionEventBlock       `json:"block"`
	AddressFrom common.Address              `json:"address_from"`
	AddressTo   common.Address              `json:"address_to"`
	NodeID      uint64                      `json:"node_id"`
	Type        schema.NodeEventType        `json:"type"`
	LogIndex    uint                        `json:"log_index"`
	ChainID     uint64                      `json:"chain_id"`
	Metadata    schema.NodeEventMetadata    `json:"metadata"`
}

func NewNodeEvent(event *schema.NodeEvent) GetNodeEventResponseData {
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
		NodeID:      event.NodeID.Uint64(),
		Type:        event.Type,
		LogIndex:    event.LogIndex,
		ChainID:     event.ChainID,
		Metadata:    event.Metadata,
	}
}

func NewNodeEvents(events []*schema.NodeEvent) GetNodeEventsResponseData {
	result := make([]*NodeEvent, len(events))
	for i, event := range events {
		result[i] = NewNodeEvent(event)
	}

	return result
}
