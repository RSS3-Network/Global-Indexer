package table

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

// EpochItem stores information for a Node in an Epoch
// TODO: we should probably rename this to NodeRewardRecord?
type EpochItem struct {
	EpochID          uint64          `gorm:"column:epoch_id;"`
	Index            int             `gorm:"column:index;primaryKey"`
	TransactionHash  string          `gorm:"column:transaction_hash;primaryKey"`
	NodeAddress      string          `gorm:"column:node_address"`
	OperationRewards decimal.Decimal `gorm:"column:operation_rewards"`
	StakingRewards   decimal.Decimal `gorm:"column:staking_rewards"`
	// FIXME: correct the column names
	TaxCollected decimal.Decimal `gorm:"column:tax_amounts"`
	// FIXME: correct the column names
	RequestCount decimal.Decimal `gorm:"column:request_counts"`
}

func (e *EpochItem) TableName() string {
	return "epoch_item"
}

func (e *EpochItem) Import(nodeToReward *schema.RewardedNode) error {
	e.EpochID = nodeToReward.EpochID
	e.Index = nodeToReward.Index
	e.TransactionHash = nodeToReward.TransactionHash.String()
	e.NodeAddress = nodeToReward.NodeAddress.String()
	e.OperationRewards = nodeToReward.OperationRewards
	e.StakingRewards = nodeToReward.StakingRewards
	e.TaxCollected = nodeToReward.TaxCollected
	e.RequestCount = nodeToReward.RequestCount

	return nil
}

func (e *EpochItem) Export() (*schema.RewardedNode, error) {
	return &schema.RewardedNode{
		EpochID:          e.EpochID,
		Index:            e.Index,
		TransactionHash:  common.HexToHash(e.TransactionHash),
		NodeAddress:      common.HexToAddress(e.NodeAddress),
		OperationRewards: e.OperationRewards,
		StakingRewards:   e.StakingRewards,
		TaxCollected:     e.TaxCollected,
		RequestCount:     e.RequestCount,
	}, nil
}

type EpochItems []*EpochItem

func (e *EpochItems) Import(nodesToReward []*schema.RewardedNode) error {
	*e = make([]*EpochItem, 0, len(nodesToReward))

	for index, nodeToReward := range nodesToReward {
		epochItem := &EpochItem{}
		if err := epochItem.Import(nodeToReward); err != nil {
			return err
		}

		epochItem.Index = index

		*e = append(*e, epochItem)
	}

	return nil
}

func (e *EpochItems) Export() ([]*schema.RewardedNode, error) {
	items := make([]*schema.RewardedNode, 0, len(*e))

	for _, epochItem := range *e {
		epochRewardItem, err := epochItem.Export()
		if err != nil {
			return nil, err
		}

		items = append(items, epochRewardItem)
	}

	return items, nil
}
