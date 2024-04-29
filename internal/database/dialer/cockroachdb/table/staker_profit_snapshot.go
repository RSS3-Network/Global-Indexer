package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type StakerProfitSnapshot struct {
	ID               uint64          `gorm:"column:id;type:bigint;not null;autoIncrement;index:profit_snapshots_epoch_id_idx,priority:2,sort:desc;index:profit_snapshots_id_idx,sort:desc;"`
	Date             time.Time       `gorm:"column:date;type:timestamp with time zone;not null;index:profit_snapshots_date_idx;"`
	OwnerAddress     common.Address  `gorm:"column:owner_address;type:bytea;not null;primaryKey;"`
	EpochID          uint64          `gorm:"column:epoch_id;type:bigint;not null;primaryKey;index:profit_snapshots_epoch_id_idx,priority:1,sort:desc;"`
	TotalChipAmounts decimal.Decimal `gorm:"column:total_chip_amounts;type:decimal;not null;index:profit_snapshots_total_chip_amounts_idx,sort:desc;"`
	TotalChipValues  decimal.Decimal `gorm:"column:total_chip_values;type:decimal;not null;index:profit_snapshots_total_chip_values_idx,sort:desc;"`
	CreatedAt        time.Time       `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()"`
	UpdatedAt        time.Time       `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()"`
}

func (s *StakerProfitSnapshot) TableName() string {
	return "stake_profit_snapshots"
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
