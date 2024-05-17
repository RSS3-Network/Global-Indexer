package schema

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BridgeEventType string

const (
	BridgeEventTypeDepositInitialized BridgeEventType = "initialized"
	BridgeEventTypeDepositFinalized   BridgeEventType = "finalized"

	BridgeEventTypeWithdrawalInitialized BridgeEventType = "initialized"
	BridgeEventTypeWithdrawalProved      BridgeEventType = "proved"
	BridgeEventTypeWithdrawalFinalized   BridgeEventType = "finalized"
)

type BridgeEventImporter interface {
	Import(bridgeEvent BridgeEvent) error
}

type BridgeEventExporter interface {
	Export() (*BridgeEvent, error)
}

type BridgeEventTransformer interface {
	BridgeEventImporter
	BridgeEventExporter
}

type BridgeEvent struct {
	ID                common.Hash     `json:"id"`
	Type              BridgeEventType `json:"type"`
	TransactionHash   common.Hash     `json:"transaction_hash"`
	TransactionIndex  uint            `json:"transaction_index"`
	TransactionStatus uint64          `json:"transaction_status"`
	ChainID           uint64          `json:"chain_id"`
	BlockHash         common.Hash     `json:"block_hash"`
	BlockNumber       *big.Int        `json:"block_number"`
	BlockTimestamp    time.Time       `json:"block_timestamp"`
}

func NewBridgeEvent(id common.Hash, eventType BridgeEventType, chainID uint64, header *types.Header, transaction *types.Transaction, receipt *types.Receipt) *BridgeEvent {
	bridgeEvent := BridgeEvent{
		ID:                id,
		Type:              eventType,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		ChainID:           chainID,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
	}

	return &bridgeEvent
}

type BridgeEventsQuery struct {
	IDs []common.Hash `query:"ids"`
}
