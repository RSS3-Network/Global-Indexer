package schema

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type NodeEventType string

const (
	NodeEventNodeCreated NodeEventType = "nodeCreated"
)

type NodeEvent struct {
	TransactionHash  common.Hash       `json:"transactionHash"`
	TransactionIndex uint              `json:"transactionIndex"`
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
	NodeCreatedMetadata *NodeCreatedMetadata `json:"nodeCreated"`
}

type NodeCreatedMetadata struct {
	Address            common.Address `json:"address"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	TaxRateBasisPoints uint64         `json:"taxRateBasisPoints"`
	PublicGood         bool           `json:"publicGood"`
}
