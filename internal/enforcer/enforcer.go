package enforcer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/distributor"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
)

type Enforcer interface {
	Verify(ctx context.Context, responses []distributor.DataResponse) error
	PartialVerify(ctx context.Context, responses []distributor.DataResponse) error
	MaintainScore(ctx context.Context) error
	ChallengeStates(ctx context.Context) error
}

type SimpleEnforcer struct {
	databaseClient database.Client
}

// Verify verifies the responses from the nodes.
func (e *SimpleEnforcer) Verify(ctx context.Context, responses []distributor.DataResponse) error {
	if len(responses) == 0 {
		return fmt.Errorf("no response returned from nodes")
	}

	nodeStatsMap, err := e.getNodeStatsMap(ctx, responses)
	if err != nil {
		return fmt.Errorf("failed to find node stats: %w", err)
	}

	// non-error and non-null results are always put in front of the list
	sortResponseByValidity(responses)
	// update requests based on data compare
	updatePointsBasedOnIdentity(responses)
	// update stats struct based on the above results
	updateStatsWithResults(nodeStatsMap, responses)
	// save stats to the database
	if err = e.databaseClient.SaveNodeStats(ctx, lo.MapToSlice(nodeStatsMap,
		func(_ common.Address, stat *schema.Stat) *schema.Stat {
			return stat
		})); err != nil {
		return fmt.Errorf("save node stats: %w", err)
	}

	return nil
}

func (e *SimpleEnforcer) getNodeStatsMap(ctx context.Context, responses []distributor.DataResponse) (map[common.Address]*schema.Stat, error) {
	stats, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: lo.Map(responses, func(response distributor.DataResponse, _ int) common.Address {
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

func updateStatsWithResults(statsMap map[common.Address]*schema.Stat, responses []distributor.DataResponse) {
	for _, response := range responses {
		if stat, exists := statsMap[response.Address]; exists {
			stat.TotalRequest += int64(response.ValidPoint)
			stat.EpochRequest += int64(response.ValidPoint)
			stat.EpochInvalidRequest += int64(response.InvalidPoint)
		}
	}
}

func (e *SimpleEnforcer) PartialVerify(_ context.Context, _ []distributor.DataResponse) error {
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
