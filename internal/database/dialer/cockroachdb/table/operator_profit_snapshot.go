package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type OperatorProfitSnapshot struct {
	ID            uint64          `gorm:"column:id;type:bigint;autoIncrement;index:operator_profit_snapshots_id_idx,sort:desc;"`
	Date          time.Time       `gorm:"column:date;type:timestamp with time zone;not null;index:operator_profit_snapshots_date_idx;"`
	Operator      common.Address  `gorm:"column:operator;type:bytea;not null;primaryKey;"`
	EpochID       uint64          `gorm:"column:epoch_id;type:bigint;not null;primaryKey;index:operator_profit_snapshots_epoch_id_idx,sort:desc;"`
	OperationPool decimal.Decimal `gorm:"column:operation_pool;type:decimal;not null;index:operator_profit_snapshots_operation_pool_idx,sort:desc;"`
	CreatedAt     time.Time       `gorm:"column:created_at;type:timestamp with time zone;not null;default:now();"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now();"`
}

func (s *OperatorProfitSnapshot) TableName() string {
	return "node_operator_profit_snapshots"
}

func (s *OperatorProfitSnapshot) Import(snapshot schema.OperatorProfitSnapshot) error {
	s.Date = snapshot.Date
	s.EpochID = snapshot.EpochID
	s.Operator = snapshot.Operator
	s.OperationPool = snapshot.OperationPool
	s.CreatedAt = snapshot.CreatedAt
	s.UpdatedAt = snapshot.UpdatedAt

	return nil
}

func (s *OperatorProfitSnapshot) Export() (*schema.OperatorProfitSnapshot, error) {
	return &schema.OperatorProfitSnapshot{
		ID:            s.ID,
		Date:          s.Date,
		EpochID:       s.EpochID,
		Operator:      s.Operator,
		OperationPool: s.OperationPool,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}, nil
}

type OperatorProfitSnapshots []OperatorProfitSnapshot

func (s *OperatorProfitSnapshots) Import(snapshots []*schema.OperatorProfitSnapshot) error {
	for _, snapshot := range snapshots {
		var imported OperatorProfitSnapshot

		if err := imported.Import(*snapshot); err != nil {
			return err
		}

		*s = append(*s, imported)
	}

	return nil
}

func (s *OperatorProfitSnapshots) Export() ([]*schema.OperatorProfitSnapshot, error) {
	snapshots := make([]*schema.OperatorProfitSnapshot, 0)

	for _, snapshot := range *s {
		exported, err := snapshot.Export()
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, exported)
	}

	return snapshots, nil
}
