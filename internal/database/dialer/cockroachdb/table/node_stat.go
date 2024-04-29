package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type Stat struct {
	Address              common.Address `gorm:"column:address;type:bytea;not null;primaryKey"`
	Endpoint             string         `gorm:"column:endpoint;type:text;not null;"`
	Points               float64        `gorm:"column:points;type:decimal;not null;index:idx_indexes_points,sort:desc;index:idx_indexes_is_full_node,priority:2,sort:desc;index:idx_indexes_is_rss_node,priority:2,sort:desc;"`
	IsPublicGood         bool           `gorm:"column:is_public_good;type:bool;not null;"`
	IsFullNode           bool           `gorm:"column:is_full_node;type:bool;not null;index:idx_indexes_is_full_node,priority:1;"`
	IsRssNode            bool           `gorm:"column:is_rss_node;type:bool;not null;index:idx_indexes_is_rss_node,priority:1;"`
	Staking              float64        `gorm:"column:staking;type:decimal;not null;"`
	Epoch                int64          `gorm:"column:epoch;type:bigint;not null;"`
	TotalRequest         int64          `gorm:"column:total_request_count;type:bigint;not null;"`
	EpochRequest         int64          `gorm:"column:epoch_request_count;type:bigint;not null;"`
	EpochInvalidRequest  int64          `gorm:"column:epoch_invalid_request_count;type:bigint;not null;index:idx_indexes_epoch_invalid_request_count,sort:asc;"`
	DecentralizedNetwork int            `gorm:"column:decentralized_network_count;type:bigint;not null;"`
	FederatedNetwork     int            `gorm:"column:federated_network_count;type:bigint;not null;"`
	Indexer              int            `gorm:"column:indexer_count;type:bigint;not null;"`
	ResetAt              time.Time      `gorm:"column:reset_at;type:timestamp with time zone;not null;"`
	CreatedAt            time.Time      `gorm:"column:created_at;type:timestamp with time zone;autoCreateTime;not null;default:now();index:idx_indexes_created_at,sort:asc;"`
	UpdatedAt            time.Time      `gorm:"column:updated_at;type:timestamp with time zone;autoUpdateTime;not null;default:now();"`
}

func (*Stat) TableName() string {
	return "node_stat"
}

func (s *Stat) Import(stat *schema.Stat) (err error) {
	s.Address = stat.Address
	s.Endpoint = stat.Endpoint
	s.Points = stat.Score
	s.IsPublicGood = stat.IsPublicGood
	s.IsFullNode = stat.IsFullNode
	s.IsRssNode = stat.IsRssNode
	s.Staking = stat.Staking
	s.Epoch = stat.Epoch
	s.TotalRequest = stat.TotalRequest
	s.EpochRequest = stat.EpochRequest
	s.EpochInvalidRequest = stat.EpochInvalidRequest
	s.DecentralizedNetwork = stat.DecentralizedNetwork
	s.FederatedNetwork = stat.FederatedNetwork
	s.Indexer = stat.Indexer
	s.ResetAt = stat.ResetAt

	return nil
}

func (s *Stat) Export() (*schema.Stat, error) {
	stat := schema.Stat{
		Address:              s.Address,
		Endpoint:             s.Endpoint,
		Score:                s.Points,
		IsPublicGood:         s.IsPublicGood,
		IsFullNode:           s.IsFullNode,
		IsRssNode:            s.IsRssNode,
		Staking:              s.Staking,
		Epoch:                s.Epoch,
		TotalRequest:         s.TotalRequest,
		EpochRequest:         s.EpochRequest,
		EpochInvalidRequest:  s.EpochInvalidRequest,
		DecentralizedNetwork: s.DecentralizedNetwork,
		FederatedNetwork:     s.FederatedNetwork,
		Indexer:              s.Indexer,
		ResetAt:              s.ResetAt,
	}

	return &stat, nil
}

type Stats []Stat

func (s *Stats) Export() ([]*schema.Stat, error) {
	stats := make([]*schema.Stat, 0)

	for _, stat := range *s {
		exportedStat, err := stat.Export()
		if err != nil {
			return nil, err
		}

		stats = append(stats, exportedStat)
	}

	return stats, nil
}

func (s *Stats) Import(stats []*schema.Stat) (err error) {
	*s = make([]Stat, 0, len(stats))

	for _, stat := range stats {
		var tStat Stat

		if err = tStat.Import(stat); err != nil {
			return err
		}

		*s = append(*s, tStat)
	}

	return nil
}
