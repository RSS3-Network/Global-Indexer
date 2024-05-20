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
	TransactionHash  common.Hash       `json:"transactionHash"`
	TransactionIndex uint              `json:"transactionIndex"`
	NodeID           *big.Int          `json:"nodeID"`
	AddressFrom      common.Address    `json:"addressFrom"`
	AddressTo        common.Address    `json:"addressTo"`
	Type             NodeEventType     `json:"type"`
	LogIndex         uint              `json:"logIndex"`
	ChainID          uint64            `json:"chainID"`
	BlockHash        common.Hash       `json:"blockHash"`
	BlockNumber      *big.Int          `json:"blockNumber"`
	BlockTimestamp   int64             `json:"blockTimestamp"`
	Metadata         NodeEventMetadata `json:"metadata"`
}

type NodeEventMetadata struct {
	NodeCreatedMetadata            *NodeCreatedMetadata            `json:"nodeCreated,omitempty"`
	NodeUpdatedMetadata            *NodeUpdatedMetadata            `json:"nodeUpdated,omitempty"`
	NodeUpdated2PublicGoodMetadata *NodeUpdated2PublicGoodMetadata `json:"nodeUpdated2PublicGood,omitempty"`
}

type NodeCreatedMetadata struct {
	NodeID             *big.Int       `json:"nodeID"`
	Address            common.Address `json:"address"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	TaxRateBasisPoints uint64         `json:"taxRateBasisPoints"`
	PublicGood         bool           `json:"publicGood"`
}

type NodeUpdatedMetadata struct {
	Address     common.Address `json:"address"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}

type NodeUpdated2PublicGoodMetadata struct {
	Address    common.Address `json:"address"`
	PublicGood bool           `json:"publicGood"`
}
