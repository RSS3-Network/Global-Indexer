package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

// EpochItem stores information for a Node in an Epoch
// TODO: we should probably rename this to NodeRewardRecord?
type EpochItem struct {
	EpochID          uint64          `gorm:"column:epoch_id;type:bigint;not null;index:idx_epoch_item_epoch_id;"`
	TransactionHash  string          `gorm:"column:transaction_hash;type:text;not null;primaryKey"`
	Index            int             `gorm:"column:index;type:bigint;not null;primaryKey;autoIncrement:false;"`
	NodeAddress      string          `gorm:"column:node_address;type:bytea;not null;index:idx_epoch_item_node_address;"`
	OperationRewards decimal.Decimal `gorm:"column:operation_rewards;type:decimal;not null;"`
	StakingRewards   decimal.Decimal `gorm:"column:staking_rewards;type:decimal;not null;"`
	TaxCollected     decimal.Decimal `gorm:"column:tax_collected;type:decimal;not null;"`
	RequestCount     decimal.Decimal `gorm:"column:request_count;type:decimal;not null;default:0;"`
	CreatedAt        time.Time       `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()"`
	UpdatedAt        time.Time       `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()"`
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
