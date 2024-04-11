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
	Verify(ctx context.Context, responses []model.DataResponse) error
	PartialVerify(ctx context.Context, responses []model.DataResponse) error
	MaintainScore(ctx context.Context) error
	ChallengeStates(ctx context.Context) error
}

type SimpleEnforcer struct {
	databaseClient database.Client
}

func (e *SimpleEnforcer) Verify(ctx context.Context, responses []model.DataResponse) error {
	if len(responses) == 0 {
		return fmt.Errorf("no response returned from nodes")
	}

	nodeStatsMap, err := e.getNodeStatsMap(ctx, responses)
	if err != nil {
		return fmt.Errorf("failed to find node stats: %w", err)
	}

	// non-error and non-null results are always in front of the list
	sort.SliceStable(responses, func(i, j int) bool {
		return (responses[i].Err == nil && responses[j].Err != nil) ||
			(responses[i].Err == nil && responses[j].Err == nil && responses[i].Valid && !responses[j].Valid)
	})

	// update requests based on data compare
	updateRequestsBasedOnDataCompare(responses)

	updateStatsWithResults(nodeStatsMap, responses)

	if err = e.databaseClient.SaveNodeStats(ctx, lo.MapToSlice(nodeStatsMap,
		func(_ common.Address, stat *schema.Stat) *schema.Stat {
			return stat
		})); err != nil {
		return fmt.Errorf("save node stats: %w", err)
	}

	return nil
}

func (e *SimpleEnforcer) getNodeStatsMap(ctx context.Context, responses []model.DataResponse) (map[common.Address]*schema.Stat, error) {
	stats, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: lo.Map(responses, func(response model.DataResponse, _ int) common.Address {
			return response.Address
		}),
	})

	if err != nil {
		return nil, err
	}

	return lo.SliceToMap(stats, func(stat *schema.Stat) (common.Address, *schema.Stat) {
		return stat.Address, stat
	}), nil
}

func updateStatsWithResults(statsMap map[common.Address]*schema.Stat, responses []model.DataResponse) {
	for _, response := range responses {
		if stat, exists := statsMap[response.Address]; exists {
			stat.TotalRequest += int64(response.Request)
			stat.EpochRequest += int64(response.Request)
			stat.EpochInvalidRequest += int64(response.InvalidRequest)
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
