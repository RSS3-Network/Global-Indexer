package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type StakerCumulativeEarningSnapshot struct {
	ID                uint64          `gorm:"column:id"`
	Date              time.Time       `gorm:"column:date"`
	EpochID           uint64          `gorm:"column:epoch_id"`
	OwnerAddress      common.Address  `gorm:"column:owner_address"`
	CumulativeEarning decimal.Decimal `gorm:"column:cumulative_earning"`
	CreatedAt         time.Time       `gorm:"column:created_at"`
	UpdatedAt         time.Time       `gorm:"column:updated_at"`
}

func (s *StakerCumulativeEarningSnapshot) TableName() string {
	return "stake.cumulative_earning_snapshots"
}

func (s *StakerCumulativeEarningSnapshot) Import(snapshot schema.StakerCumulativeEarningSnapshot) error {
	s.Date = snapshot.Date
	s.EpochID = snapshot.EpochID
	s.OwnerAddress = snapshot.OwnerAddress
	s.CumulativeEarning = snapshot.CumulativeEarning
	s.CreatedAt = snapshot.CreatedAt
	s.UpdatedAt = snapshot.UpdatedAt

	return nil
}

func (s *StakerCumulativeEarningSnapshot) Export() (*schema.StakerCumulativeEarningSnapshot, error) {
	return &schema.StakerCumulativeEarningSnapshot{
		ID:                s.ID,
		Date:              s.Date,
		EpochID:           s.EpochID,
		OwnerAddress:      s.OwnerAddress,
		CumulativeEarning: s.CumulativeEarning,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}, nil
}

type StakerCumulativeEarningSnapshots []StakerCumulativeEarningSnapshot

func (s *StakerCumulativeEarningSnapshots) Import(snapshots []*schema.StakerCumulativeEarningSnapshot) error {
	for _, snapshot := range snapshots {
		var imported StakerCumulativeEarningSnapshot

		if err := imported.Import(*snapshot); err != nil {
			return err
		}

		*s = append(*s, imported)
	}

	return nil
}

func (s *StakerCumulativeEarningSnapshots) Export() ([]*schema.StakerCumulativeEarningSnapshot, error) {
	snapshots := make([]*schema.StakerCumulativeEarningSnapshot, 0)

	for _, snapshot := range *s {
		exported, err := snapshot.Export()
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, exported)
	}

	return snapshots, nil
}
