package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type NodeMinTokensToStakeSnapshot struct {
	ID               uint64          `gorm:"column:id"`
	Date             time.Time       `gorm:"column:date"`
	EpochID          uint64          `gorm:"column:epoch_id"`
	NodeAddress      common.Address  `gorm:"column:node_address"`
	MinTokensToStake decimal.Decimal `gorm:"column:min_tokens_to_stake"`
	CreatedAt        time.Time       `gorm:"column:created_at"`
	UpdatedAt        time.Time       `gorm:"column:updated_at"`
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
