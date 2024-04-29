package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type NodeMinTokensToStakeSnapshot struct {
	ID               uint64          `gorm:"column:id;type:bigint;not null;autoIncrement;index:min_tokens_to_stake_snapshots_epoch_id_idx,priority:2,sort:desc;index:min_tokens_to_stake_snapshots_id_idx,sort:desc;"`
	Date             time.Time       `gorm:"column:date;type:timestamp with time zone;not null;index:min_tokens_to_stake_snapshots_date_idx;"`
	NodeAddress      common.Address  `gorm:"column:node_address;type:bytea;not null;primaryKey;"`
	EpochID          uint64          `gorm:"column:epoch_id;type:bigint;not null;primaryKey;index:min_tokens_to_stake_snapshots_epoch_id_idx,priority:1,sort:desc;"`
	MinTokensToStake decimal.Decimal `gorm:"column:min_tokens_to_stake;type:decimal;not null;"`
	CreatedAt        time.Time       `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()"`
	UpdatedAt        time.Time       `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()"`
}

func (s *NodeMinTokensToStakeSnapshot) TableName() string {
	return "node_min_tokens_to_stake_snapshots"
}

func (s *NodeMinTokensToStakeSnapshot) Import(snapshot schema.NodeMinTokensToStakeSnapshot) error {
	s.Date = snapshot.Date
	s.EpochID = snapshot.EpochID
	s.NodeAddress = snapshot.NodeAddress
	s.MinTokensToStake = snapshot.MinTokensToStake
	s.CreatedAt = snapshot.CreatedAt
	s.UpdatedAt = snapshot.UpdatedAt

	return nil
}

func (s *NodeMinTokensToStakeSnapshot) Export() (*schema.NodeMinTokensToStakeSnapshot, error) {
	return &schema.NodeMinTokensToStakeSnapshot{
		ID:               s.ID,
		Date:             s.Date,
		EpochID:          s.EpochID,
		NodeAddress:      s.NodeAddress,
		MinTokensToStake: s.MinTokensToStake,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}, nil
}

type NodeMinTokensToStakeSnapshots []NodeMinTokensToStakeSnapshot

func (s *NodeMinTokensToStakeSnapshots) Import(snapshots []*schema.NodeMinTokensToStakeSnapshot) error {
	for _, snapshot := range snapshots {
		var imported NodeMinTokensToStakeSnapshot

		if err := imported.Import(*snapshot); err != nil {
			return err
		}

		*s = append(*s, imported)
	}

	return nil
}

func (s *NodeMinTokensToStakeSnapshots) Export() ([]*schema.NodeMinTokensToStakeSnapshot, error) {
	snapshots := make([]*schema.NodeMinTokensToStakeSnapshot, 0)

	for _, snapshot := range *s {
		exported, err := snapshot.Export()
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, exported)
	}

	return snapshots, nil
}
