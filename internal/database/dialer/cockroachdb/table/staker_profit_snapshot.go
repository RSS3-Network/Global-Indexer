package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type StakerProfitSnapshot struct {
	gorm.Model
	ID               uint64          `gorm:"column:id"`
	Date             time.Time       `gorm:"column:date"`
	EpochID          uint64          `gorm:"column:epoch_id"`
	OwnerAddress     common.Address  `gorm:"column:owner_address"`
	TotalChipAmounts decimal.Decimal `gorm:"column:total_chip_amounts"`
	TotalChipValues  decimal.Decimal `gorm:"column:total_chip_values"`
	CreatedAt        time.Time       `gorm:"column:created_at"`
	UpdatedAt        time.Time       `gorm:"column:updated_at"`
}

func (s *StakerProfitSnapshot) TableName() string {
	return "profit_snapshots"
}

func (s *StakerProfitSnapshot) Import(snapshot schema.StakerProfitSnapshot) error {
	s.Date = snapshot.Date
	s.EpochID = snapshot.EpochID
	s.OwnerAddress = snapshot.OwnerAddress
	s.TotalChipAmounts = snapshot.TotalChipAmounts
	s.TotalChipValues = snapshot.TotalChipValues
	s.CreatedAt = snapshot.CreatedAt
	s.UpdatedAt = snapshot.UpdatedAt

	return nil
}

func (s *StakerProfitSnapshot) Export() (*schema.StakerProfitSnapshot, error) {
	return &schema.StakerProfitSnapshot{
		ID:               s.ID,
		Date:             s.Date,
		EpochID:          s.EpochID,
		OwnerAddress:     s.OwnerAddress,
		TotalChipAmounts: s.TotalChipAmounts,
		TotalChipValues:  s.TotalChipValues,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}, nil
}

type StakerProfitSnapshots []StakerProfitSnapshot

func (s *StakerProfitSnapshots) Import(snapshots []*schema.StakerProfitSnapshot) error {
	for _, snapshot := range snapshots {
		var imported StakerProfitSnapshot

		if err := imported.Import(*snapshot); err != nil {
			return err
		}

		*s = append(*s, imported)
	}

	return nil
}

func (s *StakerProfitSnapshots) Export() ([]*schema.StakerProfitSnapshot, error) {
	snapshots := make([]*schema.StakerProfitSnapshot, 0)

	for _, snapshot := range *s {
		exported, err := snapshot.Export()
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, exported)
	}

	return snapshots, nil
}
