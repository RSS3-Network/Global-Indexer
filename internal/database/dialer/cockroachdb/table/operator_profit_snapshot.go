package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type OperationPoolSnapshot struct {
	ID            uint64          `gorm:"column:id"`
	Date          time.Time       `gorm:"column:date"`
	EpochID       uint64          `gorm:"column:epoch_id"`
	Operator      common.Address  `gorm:"column:operator"`
	OperationPool decimal.Decimal `gorm:"column:operation_pool"`
	CreatedAt     time.Time       `gorm:"column:created_at"`
	UpdatedAt     time.Time       `gorm:"column:updated_at"`
}

// FIXME: update table name to `node.operation_pool_snapshots`
func (s *OperationPoolSnapshot) TableName() string {
	return "node.operator_profit_snapshots"
}

func (s *OperationPoolSnapshot) Import(snapshot schema.OperationPoolSnapshot) error {
	s.Date = snapshot.Date
	s.EpochID = snapshot.EpochID
	s.Operator = snapshot.Operator
	s.OperationPool = snapshot.OperationPool
	s.CreatedAt = snapshot.CreatedAt
	s.UpdatedAt = snapshot.UpdatedAt

	return nil
}

func (s *OperationPoolSnapshot) Export() (*schema.OperationPoolSnapshot, error) {
	return &schema.OperationPoolSnapshot{
		ID:            s.ID,
		Date:          s.Date,
		EpochID:       s.EpochID,
		Operator:      s.Operator,
		OperationPool: s.OperationPool,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}, nil
}

type OperatorProfitSnapshots []OperationPoolSnapshot

func (s *OperatorProfitSnapshots) Import(snapshots []*schema.OperationPoolSnapshot) error {
	for _, snapshot := range snapshots {
		var imported OperationPoolSnapshot

		if err := imported.Import(*snapshot); err != nil {
			return err
		}

		*s = append(*s, imported)
	}

	return nil
}

func (s *OperatorProfitSnapshots) Export() ([]*schema.OperationPoolSnapshot, error) {
	snapshots := make([]*schema.OperationPoolSnapshot, 0)

	for _, snapshot := range *s {
		exported, err := snapshot.Export()
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, exported)
	}

	return snapshots, nil
}
