package table

import (
	"time"

	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type EpochAPYSnapshot struct {
	EpochID   uint64          `gorm:"column:epoch_id"`
	Date      time.Time       `gorm:"column:date"`
	APY       decimal.Decimal `gorm:"column:apy"`
	CreatedAt time.Time       `gorm:"column:created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
}

func (e *EpochAPYSnapshot) TableName() string {
	return "epoch.apy_snapshots"
}

func (e *EpochAPYSnapshot) Import(epochAPYSnapshot *schema.EpochAPYSnapshot) error {
	e.EpochID = epochAPYSnapshot.EpochID
	e.Date = epochAPYSnapshot.Date
	e.APY = epochAPYSnapshot.APY

	return nil
}

func (e *EpochAPYSnapshot) Export() (*schema.EpochAPYSnapshot, error) {
	return &schema.EpochAPYSnapshot{
		EpochID:   e.EpochID,
		Date:      e.Date,
		APY:       e.APY,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}, nil
}

type EpochAPYSnapshots []EpochAPYSnapshot

func (e *EpochAPYSnapshots) Import(snapshots []*schema.EpochAPYSnapshot) error {
	for _, snapshot := range snapshots {
		var imported EpochAPYSnapshot

		if err := imported.Import(snapshot); err != nil {
			return err
		}

		*e = append(*e, imported)
	}

	return nil
}

func (e *EpochAPYSnapshots) Export() ([]*schema.EpochAPYSnapshot, error) {
	snapshots := make([]*schema.EpochAPYSnapshot, 0, len(*e))

	for _, snapshot := range *e {
		exported, err := snapshot.Export()
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, exported)
	}

	return snapshots, nil
}
