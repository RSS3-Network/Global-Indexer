package enforcer

import (
	"context"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/distributor"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Enforcer interface {
	Verify(ctx context.Context, responses []distributor.DataResponse) error
	PartialVerify(ctx context.Context, responses []distributor.DataResponse)
	MaintainScore(ctx context.Context) error
	ChallengeStates(ctx context.Context) error
}

type SimpleEnforcer struct {
	databaseClient database.Client
	httpClient     httputil.Client
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
			stat.TotalRequest++
			stat.EpochRequest += int64(response.ValidPoint)
			stat.EpochInvalidRequest += int64(response.InvalidPoint)
		}
	}
}

// PartialVerify performs a partial verification of the responses from the nodes.
func (e *SimpleEnforcer) PartialVerify(ctx context.Context, responses []distributor.DataResponse) {
	// Check if there are any responses
	if len(responses) == 0 {
		zap.L().Warn("no response returned from nodes")

		return
	}

	activities := &distributor.ActivitiesResponse{}
	// TODO: Consider selecting response that have been successfully verified as data source
	// and now select the first response as data source
	data := responses[0].Data

	// Check if the data is valid
	if !isDataValid(data, activities) {
		zap.L().Warn("failed to parse response")

		return
	}

	// Check if there are any feeds in the activities data
	if len(activities.Data) == 0 {
		zap.L().Warn("no feed returned from nodes")

		return
	}

	workingNodes := lo.Map(responses, func(result distributor.DataResponse, _ int) common.Address {
		return result.Address
	})

	e.verifyPartialFeeds(ctx, activities.Data, workingNodes)
}

// verifyPartialFeeds filter feeds based on the platform to perform partial verification.
func (e *SimpleEnforcer) verifyPartialFeeds(ctx context.Context, feeds []*distributor.Feed, workingNodes []common.Address) {
	// platformMap is used to store the platform that has been verified
	platformMap := make(map[string]struct{}, distributor.DefaultVerifyCount)
	// statMap is used to store the stats that have been verified
	statMap := make(map[string]struct{})

	for _, feed := range feeds {
		// This usually indicates that the feed belongs to the fallback worker.
		// We cannot determine whether this feed belongs to a readable worker，
		// therefore it is skipped.
		if len(feed.Platform) == 0 {
			continue
		}

		// Find stats that related to the platform
		stats, err := e.findStatsByPlatform(ctx, feed, workingNodes)

		if err != nil {
			zap.L().Error("failed to verify platform", zap.Error(err))

			continue
		}

		if len(stats) == 0 {
			zap.L().Warn("no stats match the platform")

			continue
		}

		// Verify the feed by stats
		e.verifyFeedByStats(ctx, feed, stats, statMap, platformMap)

		// If the platform count reaches the DefaultVerifyCount, exit the verification loop.
		if _, exists := platformMap[feed.Platform]; !exists {
			if len(platformMap) == distributor.DefaultVerifyCount {
				break
			}
		}
	}
}

// findStatsByPlatform finds the stats by platform.
func (e *SimpleEnforcer) findStatsByPlatform(ctx context.Context, feed *distributor.Feed, workingNodes []common.Address) ([]*schema.Stat, error) {
	pid, err := filter.PlatformString(feed.Platform)
	if err != nil {
		return nil, err
	}

	worker := distributor.PlatformToWorkerMap[pid]
	indexers, err := e.databaseClient.FindNodeIndexers(ctx, nil, []string{feed.Network}, []string{worker})

	if err != nil {
		return nil, err
	}

	nodeAddresses := excludeWorkingNodes(indexers, workingNodes)

	stats, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList:  nodeAddresses,
		ValidRequest: lo.ToPtr(distributor.DefaultSlashCount),
		PointsOrder:  lo.ToPtr("DESC"),
	})

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// excludeWorkingNodes excludes the working nodes from the indexers.
func excludeWorkingNodes(indexers []*schema.Indexer, workingNodes []common.Address) []common.Address {
	nodeAddresses := lo.Map(indexers, func(indexer *schema.Indexer, _ int) common.Address {
		return indexer.Address
	})

	// filter out the working nodes
	return lo.Filter(nodeAddresses, func(item common.Address, _ int) bool {
		return !lo.Contains(workingNodes, item)
	})
}

// verifyFeedByStats verifies the feed by stats.
func (e *SimpleEnforcer) verifyFeedByStats(ctx context.Context, feed *distributor.Feed, stats []*schema.Stat, statMap, platformMap map[string]struct{}) {
	for _, stat := range stats {
		if _, exists := statMap[stat.Address.String()]; !exists {
			statMap[stat.Address.String()] = struct{}{}

			activity, err := e.fetchActivityByTxID(ctx, stat.Endpoint, feed.ID)

			if err != nil {
				stat.EpochInvalidRequest += invalidPointUnit
			} else {
				if activity.Data == nil || !isActivityIdentical(feed, activity.Data) {
					// TODO: if false, save the record to the database
					stat.EpochInvalidRequest += invalidPointUnit
				} else {
					stat.TotalRequest++
					stat.EpochRequest += validPointUnit
				}
			}

			platformMap[feed.Platform] = struct{}{}

			if err = e.databaseClient.SaveNodeStat(ctx, stat); err != nil {
				zap.L().Warn("[verifyStat] failed to save node stat", zap.Error(err))
			}

			break
		}
	}
}

// fetchActivityByTxID fetches the activity by txID.
func (e *SimpleEnforcer) fetchActivityByTxID(ctx context.Context, endpoint, txID string) (*distributor.ActivityResponse, error) {
	fullURL := endpoint + "/decentralized/tx/" + txID

	body, err := e.httpClient.Fetch(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	activity := &distributor.ActivityResponse{}
	if isDataValid(data, activity) {
		return activity, nil
	}

	return nil, fmt.Errorf("invalid data")
}

func (e *SimpleEnforcer) MaintainScore(_ context.Context) error {
	return nil
}

func (e *SimpleEnforcer) ChallengeStates(_ context.Context) error {
	return nil
}

func NewSimpleEnforcer(databaseClient database.Client, httpClient httputil.Client) (*SimpleEnforcer, error) {
	return &SimpleEnforcer{
		databaseClient: databaseClient,
		httpClient:     httpClient,
	}, nil
}
