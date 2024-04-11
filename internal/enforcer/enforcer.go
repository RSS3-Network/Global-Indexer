package enforcer

import (
	"context"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type Enforcer interface {
	Verify(ctx context.Context, results []model.DataResponse) error
	PartialVerify(ctx context.Context, results []model.DataResponse) error
	MaintainScore(ctx context.Context) error
	ChallengeStates(ctx context.Context) error
}

type SimpleEnforcer struct {
	databaseClient database.Client
}

func (e *SimpleEnforcer) Verify(ctx context.Context, results []model.DataResponse) error {
	if len(results) == 0 {
		return fmt.Errorf("no response returned from nodes")
	}

	nodeStatsMap, err := e.getNodeStatsMap(ctx, results)
	if err != nil {
		return fmt.Errorf("failed to find node stats: %w", err)
	}

	// non-error and non-null results are always in front of the list
	sort.SliceStable(results, func(i, j int) bool {
		return (results[i].Err == nil && results[j].Err != nil) ||
			(results[i].Err == nil && results[j].Err == nil && results[i].Valid && !results[j].Valid)
	})

	// update requests based on data compare
	updateRequestsBasedOnDataCompare(results)

	updateStatsWithResults(nodeStatsMap, results)

	if err = e.databaseClient.SaveNodeStats(ctx, lo.MapToSlice(nodeStatsMap,
		func(_ common.Address, stat *schema.Stat) *schema.Stat {
			return stat
		})); err != nil {
		return fmt.Errorf("save node stats: %w", err)
	}

	return nil
}

func (e *SimpleEnforcer) getNodeStatsMap(ctx context.Context, results []model.DataResponse) (map[common.Address]*schema.Stat, error) {
	stats, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: lo.Map(results, func(result model.DataResponse, _ int) common.Address {
			return result.Address
		}),
	})

	if err != nil {
		return nil, err
	}

	return lo.SliceToMap(stats, func(stat *schema.Stat) (common.Address, *schema.Stat) {
		return stat.Address, stat
	}), nil
}

func updateStatsWithResults(statsMap map[common.Address]*schema.Stat, results []model.DataResponse) {
	for _, result := range results {
		if stat, exists := statsMap[result.Address]; exists {
			stat.TotalRequest += int64(result.Request)
			stat.EpochRequest += int64(result.Request)
			stat.EpochInvalidRequest += int64(result.InvalidRequest)
		}
	}
}

func (e *SimpleEnforcer) PartialVerify(_ context.Context, _ []model.DataResponse) error {
	return nil
}

func (e *SimpleEnforcer) MaintainScore(_ context.Context) error {
	return nil
}

func (e *SimpleEnforcer) ChallengeStates(_ context.Context) error {
	return nil
}

func NewSimpleEnforcer(databaseClient database.Client) (*SimpleEnforcer, error) {
	return &SimpleEnforcer{
		databaseClient: databaseClient,
	}, nil
}
