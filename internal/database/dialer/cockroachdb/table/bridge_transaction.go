package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                         = (*BridgeTransaction)(nil)
	_ schema.BridgeTransactionTransformer = (*BridgeTransaction)(nil)
)

type BridgeTransaction struct {
	ID               string          `gorm:"column:id;type:text;not null;primaryKey;"`
	Type             string          `gorm:"column:type;type:text;not null;primaryKey;"`
	Sender           string          `gorm:"column:sender;type:text;not null;index:idx_transactions_sender;index:idx_transactions_address,priority:1;"`
	Receiver         string          `gorm:"column:receiver;type:text;not null;index:idx_transactions_receiver;index:idx_transactions_address,priority:2;"`
	TokenAddressL1   *string         `gorm:"column:token_address_l1;type:text;"`
	TokenAddressL2   *string         `gorm:"column:token_address_l2;type:text"`
	TokenValue       decimal.Decimal `gorm:"column:token_value;type:decimal;not null;"`
	Data             string          `gorm:"column:data;type:text"`
	ChainID          uint64          `gorm:"column:chain_id;type:bigint;not null;"`
	BlockTimestamp   time.Time       `gorm:"column:block_timestamp;type:bigint;index:idx_transactions_order,priority:1,sort:desc;"`
	BlockNumber      uint64          `gorm:"column:block_number;type:bigint;index:idx_transactions_order,priority:2,sort:desc;"`
	TransactionIndex uint            `gorm:"column:transaction_index;type:timestamp with time zone;index:idx_transactions_order,priority:3,sort:desc"`
}

func (b *BridgeTransaction) TableName() string {
	return "bridge_transactions"
}

func (b *BridgeTransaction) Import(bridgeTransaction schema.BridgeTransaction) error {
	b.ID = bridgeTransaction.ID.String()
	b.Type = string(bridgeTransaction.Type)
	b.Sender = bridgeTransaction.Sender.String()
	b.Receiver = bridgeTransaction.Receiver.String()
	b.TokenAddressL1 = lo.ToPtr(bridgeTransaction.TokenAddressL1.String())
	b.TokenAddressL2 = lo.ToPtr(bridgeTransaction.TokenAddressL2.String())
	b.TokenValue = decimal.NewFromBigInt(bridgeTransaction.TokenValue, 0)
	b.Data = bridgeTransaction.Data
	b.ChainID = bridgeTransaction.ChainID
	b.BlockTimestamp = bridgeTransaction.BlockTimestamp
	b.BlockNumber = bridgeTransaction.BlockNumber
	b.TransactionIndex = bridgeTransaction.TransactionIndex

	return nil
}

func (b *BridgeTransaction) Export() (*schema.BridgeTransaction, error) {
	bridgeTransaction := schema.BridgeTransaction{
		ID:       common.HexToHash(b.ID),
		Type:     schema.BridgeTransactionType(b.Type),
		Sender:   common.HexToAddress(b.Sender),
		Receiver: common.HexToAddress(b.Receiver),
		TokenAddressL1: func(tokenAddress *string) *common.Address {
			if tokenAddress == nil {
				return nil
			}

			return lo.ToPtr(common.HexToAddress(*tokenAddress))
		}(b.TokenAddressL1),
		TokenAddressL2: func(tokenAddress *string) *common.Address {
			if tokenAddress == nil {
				return nil
			}

			return lo.ToPtr(common.HexToAddress(*tokenAddress))
		}(b.TokenAddressL2),
		TokenValue:       b.TokenValue.BigInt(),
		Data:             b.Data,
		ChainID:          b.ChainID,
		BlockTimestamp:   b.BlockTimestamp,
		BlockNumber:      b.BlockNumber,
		TransactionIndex: b.TransactionIndex,
	}

	return &bridgeTransaction, nil
}
