package schema

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type NodeEventType string

const (
	NodeEventNodeCreated            NodeEventType = "nodeCreated"
	NodeEventNodeUpdated            NodeEventType = "nodeUpdated"
	NodeEventNodeUpdated2PublicGood NodeEventType = "nodeUpdated2PublicGood"
)

type NodeEvent struct {
	TransactionHash  common.Hash       `json:"transaction_hash"`
	TransactionIndex uint              `json:"transaction_index"`
	NodeID           *big.Int          `json:"node_id"`
	AddressFrom      common.Address    `json:"address_from"`
	AddressTo        common.Address    `json:"address_to"`
	Type             NodeEventType     `json:"type"`
	LogIndex         uint              `json:"log_index"`
	ChainID          uint64            `json:"chain_id"`
	BlockHash        common.Hash       `json:"block_hash"`
	BlockNumber      *big.Int          `json:"block_number"`
	BlockTimestamp   int64             `json:"block_timestamp"`
	Metadata         NodeEventMetadata `json:"metadata"`
}

type NodeEventMetadata struct {
	NodeCreatedMetadata            *NodeCreatedMetadata            `json:"node_created,omitempty"`
	NodeUpdatedMetadata            *NodeUpdatedMetadata            `json:"node_updated,omitempty"`
	NodeUpdated2PublicGoodMetadata *NodeUpdated2PublicGoodMetadata `json:"node_updated_to_public_good,omitempty"`
}

type NodeCreatedMetadata struct {
	NodeID             *big.Int       `json:"node_id"`
	Address            common.Address `json:"address"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	TaxRateBasisPoints uint64         `json:"tax_rate_basis_points"`
	PublicGood         bool           `json:"public_good"`
}

type NodeUpdatedMetadata struct {
	Address     common.Address `json:"address"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}

type NodeUpdated2PublicGoodMetadata struct {
	Address    common.Address `json:"address"`
	PublicGood bool           `json:"public_good"`
}
