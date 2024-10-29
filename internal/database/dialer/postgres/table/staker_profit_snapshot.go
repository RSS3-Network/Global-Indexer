package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type StakerProfitSnapshot struct {
	ID              uint64          `gorm:"column:id"`
	Date            time.Time       `gorm:"column:date"`
	EpochID         uint64          `gorm:"column:epoch_id"`
	OwnerAddress    common.Address  `gorm:"column:owner_address"`
	TotalChipAmount decimal.Decimal `gorm:"column:total_chip_amounts"` // Fixme: total_chip_amounts-> total_chip_amount
	TotalChipValue  decimal.Decimal `gorm:"column:total_chip_values"`  // Fixme: total_chip_values-> total_chip_value
	CreatedAt       time.Time       `gorm:"column:created_at"`
	UpdatedAt       time.Time       `gorm:"column:updated_at"`
}

func (s *StakerProfitSnapshot) TableName() string {
	return "stake.profit_snapshots"
}

func (s *StakerProfitSnapshot) Import(snapshot schema.StakerProfitSnapshot) error {
	s.Date = snapshot.Date
	s.EpochID = snapshot.EpochID
	s.OwnerAddress = snapshot.OwnerAddress
	s.TotalChipAmount = snapshot.TotalChipAmount
	s.TotalChipValue = snapshot.TotalChipValue
	s.CreatedAt = snapshot.CreatedAt
	s.UpdatedAt = snapshot.UpdatedAt

	return nil
}

func (s *StakerProfitSnapshot) Export() (*schema.StakerProfitSnapshot, error) {
	return &schema.StakerProfitSnapshot{
		ID:              s.ID,
		Date:            s.Date,
		EpochID:         s.EpochID,
		OwnerAddress:    s.OwnerAddress,
		TotalChipAmount: s.TotalChipAmount,
		TotalChipValue:  s.TotalChipValue,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
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
