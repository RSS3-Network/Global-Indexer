package enforcer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Enforcer interface {
	VerifyResponses(ctx context.Context, responses []model.DataResponse) error
	VerifyPartialResponses(ctx context.Context, epochID uint64, responses []model.DataResponse)
	MaintainScore(ctx context.Context) error
	ChallengeStates(ctx context.Context) error
}

type SimpleEnforcer struct {
	cacheClient     cache.Client
	databaseClient  database.Client
	httpClient      httputil.Client
	stakingContract *l2.Staking
}

// VerifyResponses verifies the responses from the Nodes.
func (e *SimpleEnforcer) VerifyResponses(ctx context.Context, responses []*model.DataResponse) error {
	if len(responses) == 0 {
		return fmt.Errorf("no response returned from nodes")
	}

	nodeStatsMap, err := e.getNodeStatsMap(ctx, responses)
	if err != nil {
		return fmt.Errorf("failed to Find node stats: %w", err)
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
		return fmt.Errorf("save Node stats: %w", err)
	}

	return nil
}

// VerifyPartialResponses performs a partial verification of the responses from the Nodes.
func (e *SimpleEnforcer) VerifyPartialResponses(ctx context.Context, epochID uint64, responses []*model.DataResponse) {
	// Check if there are any responses
	if len(responses) == 0 {
		zap.L().Warn("no response returned from nodes")

		return
	}

	activities := &model.ActivitiesResponse{}
	// TODO: Consider selecting response that have been successfully verified as data source
	// and now select the first response as data source
	data := responses[0].Data

	// Check if the data is valid
	if !isDataValid(data, activities) {
		zap.L().Warn("failed to parse response")

		return
	}

	// Check if there are any activities in the activities responses data
	if len(activities.Data) == 0 {
		zap.L().Warn("no activities returned from nodes")

		return
	}

	workingNodes := lo.Map(responses, func(result *model.DataResponse, _ int) common.Address {
		return result.Address
	})

	e.verifyPartialActivities(ctx, epochID, responses[0], activities.Data, workingNodes)
}

func (e *SimpleEnforcer) getNodeStatsMap(ctx context.Context, responses []*model.DataResponse) (map[common.Address]*schema.Stat, error) {
	stats, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: lo.Map(responses, func(response *model.DataResponse, _ int) common.Address {
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

func updateStatsWithResults(statsMap map[common.Address]*schema.Stat, responses []*model.DataResponse) {
	for _, response := range responses {
		if stat, exists := statsMap[response.Address]; exists {
			stat.TotalRequest++
			stat.EpochRequest += int64(response.ValidPoint)
			stat.EpochInvalidRequest += int64(response.InvalidPoint)
		}
	}
}

// verifyPartialActivities filter Activity based on the platform to perform a partial verification.
func (e *SimpleEnforcer) verifyPartialActivities(ctx context.Context, epochID uint64, validResponse *model.DataResponse, activities []*model.Activity, workingNodes []common.Address) {
	// platformMap is used to store the platform that has been verified
	platformMap := make(map[string]struct{}, model.RequiredVerificationCount)
	// statMap is used to store the stats that have been verified
	statMap := make(map[string]struct{})

	nodeInvalidResponse := &schema.NodeInvalidResponse{
		EpochID:        epochID,
		ValidatorNodes: []common.Address{validResponse.Address},
	}

	for _, activity := range activities {
		// This usually indicates that the activity belongs to the fallback worker.
		// We cannot determine whether this activity belongs to a readable worker，
		// therefore it is skipped.
		if len(activity.Platform) == 0 {
			continue
		}

		// Find stats that related to the platform
		stats, err := e.findStatsByPlatform(ctx, activity, workingNodes)

		if err != nil {
			zap.L().Error("failed to verify platform", zap.Error(err))

			continue
		}

		if len(stats) == 0 {
			zap.L().Warn("no stats match the platform")

			continue
		}

		e.verifyActivityByStats(ctx, activity, stats, statMap, platformMap, nodeInvalidResponse)

		// If the platform count reaches the RequiredVerificationCount, exit the verification loop.
		if _, exists := platformMap[activity.Platform]; !exists {
			if len(platformMap) == model.RequiredVerificationCount {
				break
			}
		}
	}
}

// findStatsByPlatform finds the required stats based on the platform.
func (e *SimpleEnforcer) findStatsByPlatform(ctx context.Context, activity *model.Activity, workingNodes []common.Address) ([]*schema.Stat, error) {
	pid, err := filter.PlatformString(activity.Platform)
	if err != nil {
		return nil, err
	}

	worker := model.PlatformToWorkerMap[pid]
	indexers, err := e.databaseClient.FindNodeIndexers(ctx, nil, []string{activity.Network}, []string{worker})

	if err != nil {
		return nil, err
	}

	nodeAddresses := excludeWorkingNodes(indexers, workingNodes)

	stats, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList:  nodeAddresses,
		ValidRequest: lo.ToPtr(model.DemotionCountBeforeSlashing),
		PointsOrder:  lo.ToPtr("DESC"),
	})

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// excludeWorkingNodes excludes the working Nodes from the indexers.
func excludeWorkingNodes(indexers []*schema.Indexer, workingNodes []common.Address) []common.Address {
	nodeAddresses := lo.Map(indexers, func(indexer *schema.Indexer, _ int) common.Address {
		return indexer.Address
	})

	// filter out the working nodes
	return lo.Filter(nodeAddresses, func(item common.Address, _ int) bool {
		return !lo.Contains(workingNodes, item)
	})
}

// verifyActivityByStats verifies the activity based on stats nodes that meet specific criteria.
func (e *SimpleEnforcer) verifyActivityByStats(ctx context.Context, activity *model.Activity, stats []*schema.Stat, statMap, platformMap map[string]struct{}, nodeInvalidResponse *schema.NodeInvalidResponse) {
	for _, stat := range stats {
		if _, exists := statMap[stat.Address.String()]; !exists {
			statMap[stat.Address.String()] = struct{}{}

			activityFetched, err := e.fetchActivityByTxID(ctx, stat.Endpoint, activity.ID)

			if err != nil || activityFetched.Data == nil || !isActivityIdentical(activity, activityFetched.Data) {
				stat.EpochInvalidRequest += invalidPointUnit

				nodeInvalidResponse.InvalidType = lo.Ternary(err != nil, schema.NodeInvalidResponseTypeError, schema.NodeInvalidResponseTypeInconsistent)
				nodeInvalidResponse.InvalidResponse = generateFaultyResponse(err, activityFetched)
			} else {
				stat.TotalRequest++
				stat.EpochRequest += validPointUnit
			}

			// If the request is invalid, save the invalid response to the database.
			if stat.EpochInvalidRequest > 0 {
				nodeInvalidResponse.Node = stat.Address
				nodeInvalidResponse.Request = stat.Endpoint + "/decentralized/tx/" + activity.ID

				validData, _ := json.Marshal(activity)
				nodeInvalidResponse.ValidatorResponse = validData

				if err = e.databaseClient.SaveNodeInvalidResponse(ctx, nodeInvalidResponse); err != nil {
					zap.L().Error("save node invalid response", zap.Error(err))
				}
			}

			platformMap[activity.Platform] = struct{}{}

			if err = e.databaseClient.SaveNodeStat(ctx, stat); err != nil {
				zap.L().Warn("[verifyStat] failed to save node stat", zap.Error(err))
			}

			break
		}
	}
}

func generateFaultyResponse(err error, activity *model.ActivityResponse) json.RawMessage {
	if err != nil {
		return json.RawMessage(err.Error())
	}

	rawData, _ := json.Marshal(activity.Data)

	return rawData
}

// fetchActivityByTxID fetches the activity by txID from a Node.
func (e *SimpleEnforcer) fetchActivityByTxID(ctx context.Context, nodeEndpoint, txID string) (*model.ActivityResponse, error) {
	fullURL := nodeEndpoint + "/decentralized/tx/" + txID

	body, err := e.httpClient.Fetch(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	activity := &model.ActivityResponse{}
	if isDataValid(data, activity) {
		return activity, nil
	}

	return nil, fmt.Errorf("invalid data")
}

// MaintainReliabilityScore maintains the Reliability Score σ for all Nodes.
// σ is used to determine the probability of a Node receiving a request on DSL.
func (e *SimpleEnforcer) MaintainReliabilityScore(ctx context.Context) error {
	// Retrieve the most recently indexed epoch.
	currentEpoch, err := e.getCurrentEpoch(ctx)
	if err != nil {
		return err
	}

	query := &schema.StatQuery{Limit: lo.ToPtr(defaultLimit)}

	// Traverse the entire node and update its score.
	for {
		stats, err := e.databaseClient.FindNodeStats(ctx, query)
		if err != nil {
			return err
		}

		// If there are no stats, exit the loop.
		if len(stats) == 0 {
			break
		}

		if err = e.processNodeStats(ctx, stats, currentEpoch); err != nil {
			return err
		}

		lastStat := stats[len(stats)-1]
		query.Cursor = lo.ToPtr(lastStat.Address.String())
	}

	return e.updateNodeCache(ctx)
}

func (e *SimpleEnforcer) ChallengeStates(_ context.Context) error {
	return nil
}

func NewSimpleEnforcer(databaseClient database.Client, cacheClient cache.Client, stakingContract *l2.Staking, httpClient httputil.Client) *SimpleEnforcer {
	return &SimpleEnforcer{
		databaseClient:  databaseClient,
		cacheClient:     cacheClient,
		stakingContract: stakingContract,
		httpClient:      httpClient,
	}
}
