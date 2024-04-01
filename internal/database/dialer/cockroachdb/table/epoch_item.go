package table

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/shopspring/decimal"
)

type EpochItem struct {
	EpochID          uint64          `gorm:"column:epoch_id;"`
	Index            int             `gorm:"column:index;primaryKey"`
	TransactionHash  string          `gorm:"column:transaction_hash;primaryKey"`
	NodeAddress      string          `gorm:"column:node_address"`
	OperationRewards decimal.Decimal `gorm:"column:operation_rewards"`
	StakingRewards   decimal.Decimal `gorm:"column:staking_rewards"`
	TaxAmounts       decimal.Decimal `gorm:"column:tax_amounts"`
	RequestCounts    decimal.Decimal `gorm:"column:request_counts"`
}

func (e *EpochItem) TableName() string {
	return "epoch_item"
}

func (e *EpochItem) Import(epochRewardItem *schema.EpochItem) error {
	e.EpochID = epochRewardItem.EpochID
	e.Index = epochRewardItem.Index
	e.TransactionHash = epochRewardItem.TransactionHash.String()
	e.NodeAddress = epochRewardItem.NodeAddress.String()
	e.OperationRewards = epochRewardItem.OperationRewards
	e.StakingRewards = epochRewardItem.StakingRewards
	e.TaxAmounts = epochRewardItem.TaxAmounts
	e.RequestCounts = epochRewardItem.RequestCounts

	return nil
}

func (e *EpochItem) Export() (*schema.EpochItem, error) {
	return &schema.EpochItem{
		EpochID:          e.EpochID,
		Index:            e.Index,
		TransactionHash:  common.HexToHash(e.TransactionHash),
		NodeAddress:      common.HexToAddress(e.NodeAddress),
		OperationRewards: e.OperationRewards,
		StakingRewards:   e.StakingRewards,
		TaxAmounts:       e.TaxAmounts,
		RequestCounts:    e.RequestCounts,
	}, nil
}

type EpochItems []*EpochItem

func (e *EpochItems) Import(epochRewardItems []*schema.EpochItem) error {
	*e = make([]*EpochItem, 0, len(epochRewardItems))

	for index, epochRewardItem := range epochRewardItems {
		epochItem := &EpochItem{}
		if err := epochItem.Import(epochRewardItem); err != nil {
			return err
		}

		epochItem.Index = index

		*e = append(*e, epochItem)
	}

	return nil
}

func (e *EpochItems) Export() ([]*schema.EpochItem, error) {
	items := make([]*schema.EpochItem, 0, len(*e))

	for _, epochItem := range *e {
		epochRewardItem, err := epochItem.Export()
		if err != nil {
			return nil, err
		}

		items = append(items, epochRewardItem)
	}

	return items, nil
}
