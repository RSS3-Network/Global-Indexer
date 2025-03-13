package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type Stat struct {
	Address              common.Address `gorm:"column:address;primaryKey"`
	Endpoint             string         `gorm:"column:endpoint"`
	AccessToken          string         `gorm:"column:access_token"`
	Points               float64        `gorm:"column:points"`
	IsPublicGood         bool           `gorm:"column:is_public_good"`
	IsFullNode           bool           `gorm:"column:is_full_node"`
	IsRssNode            bool           `gorm:"column:is_rss_node"`
	IsAINode             bool           `gorm:"column:is_ai_node"`
	Staking              float64        `gorm:"column:staking"`
	Epoch                int64          `gorm:"column:epoch"`
	TotalRequest         int64          `gorm:"column:total_request_count"`
	EpochRequest         int64          `gorm:"column:epoch_request_count"`
	EpochInvalidRequest  int64          `gorm:"column:epoch_invalid_request_count"`
	DecentralizedNetwork int            `gorm:"column:decentralized_network_count"`
	FederatedNetwork     int            `gorm:"column:federated_network_count"`
	Indexer              int            `gorm:"column:indexer_count"`
	ResetAt              time.Time      `gorm:"column:reset_at"`
	CreatedAt            time.Time      `gorm:"column:created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at"`
}

func (*Stat) TableName() string {
	return "node_stat"
}

func (s *Stat) Import(stat *schema.Stat) (err error) {
	s.Address = stat.Address
	s.Endpoint = stat.Endpoint
	s.AccessToken = stat.AccessToken
	s.Points = stat.Score
	s.IsPublicGood = stat.IsPublicGood
	s.IsFullNode = stat.IsFullNode
	s.IsRssNode = stat.IsRssNode
	s.IsAINode = stat.IsAINode
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
		AccessToken:          s.AccessToken,
		Score:                s.Points,
		IsPublicGood:         s.IsPublicGood,
		IsFullNode:           s.IsFullNode,
		IsRssNode:            s.IsRssNode,
		IsAINode:             s.IsAINode,
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
