package enforcer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/node/schema/worker/decentralized"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

const (
	// a valid response gives 1 point
	validPointUnit = 1
	// an invalid response gives 1 point (in a bad way)
	invalidPointUnit = 1
)

func (e *SimpleEnforcer) getNodeStatsMap(ctx context.Context, responses []*model.DataResponse) (map[common.Address]*schema.Stat, error) {
	stats, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		Addresses: lo.Map(responses, func(response *model.DataResponse, _ int) common.Address {
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

func (e *SimpleEnforcer) batchUpdateScoreMaintainer(ctx context.Context, nodeStats []*schema.Stat) {
	statsPool := pool.New().WithContext(ctx).WithMaxGoroutines(lo.Ternary(len(nodeStats) < 20*runtime.NumCPU(), len(nodeStats), 20*runtime.NumCPU()))

	for i := range nodeStats {
		i := i

		statsPool.Go(func(_ context.Context) error {
			calculateReliabilityScore(nodeStats[i])

			e.updateScoreMaintainer(ctx, nodeStats[i])

			return nil
		})
	}

	if err := statsPool.Wait(); err != nil {
		zap.L().Error("failed to update score maintainer", zap.Error(err))
	}
}

func (e *SimpleEnforcer) updateScoreMaintainer(ctx context.Context, nodeStat *schema.Stat) {
	nodeCache := &model.NodeEndpointCache{
		Address:      nodeStat.Address.String(),
		Score:        nodeStat.Score,
		Endpoint:     nodeStat.Endpoint,
		InvalidCount: nodeStat.EpochInvalidRequest,
	}

	if err := e.fullNodeScoreMaintainer.addOrUpdateScore(ctx, model.FullNodeCacheKey, nodeCache); err != nil {
		zap.L().Error("failed to update full node score", zap.Error(err), zap.String("address", nodeStat.Address.String()))
	}

	if err := e.rssNodeScoreMaintainer.addOrUpdateScore(ctx, model.RssNodeCacheKey, nodeCache); err != nil {
		zap.L().Error("failed to update rss node score", zap.Error(err), zap.String("address", nodeStat.Address.String()))
	}
}

// verifyPartialActivities filter Activity based on the platform to perform a partial verification.
func (e *SimpleEnforcer) verifyPartialActivities(ctx context.Context, epochID uint64, validResponse *model.DataResponse, activities []*model.Activity, workingNodes []common.Address) {
	// platformMap is used to store the platform that has been verified
	platformMap := make(map[string]struct{}, model.RequiredVerificationCount)
	// statMap is used to store the stats that have been verified
	statMap := make(map[string]struct{})

	nodeInvalidResponse := &schema.NodeInvalidResponse{
		EpochID:       epochID,
		VerifierNodes: []common.Address{validResponse.Address},
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
	pid, err := decentralized.PlatformString(activity.Platform)
	if err != nil {
		return nil, err
	}

	workers := model.PlatformToWorkersMap[pid.String()]

	indexers, err := e.databaseClient.FindNodeWorkers(ctx, &schema.WorkerQuery{
		Networks: []string{activity.Network},
		Names:    workers,
	})

	if err != nil {
		return nil, err
	}

	nodeAddresses := excludeWorkingNodes(indexers, workingNodes)

	if len(nodeAddresses) == 0 {
		return nil, nil
	}

	stats, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		Addresses:    nodeAddresses,
		ValidRequest: lo.ToPtr(model.DemotionCountBeforeSlashing),
		PointsOrder:  lo.ToPtr("DESC"),
	})

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// excludeWorkingNodes excludes the working Nodes from the indexers.
func excludeWorkingNodes(indexers []*schema.Worker, workingNodes []common.Address) []common.Address {
	nodeAddresses := lo.Map(indexers, func(indexer *schema.Worker, _ int) common.Address {
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

				nodeInvalidResponse.Type = lo.Ternary(err != nil, schema.NodeInvalidResponseTypeError, schema.NodeInvalidResponseTypeInconsistent)
				nodeInvalidResponse.Response = generateInvalidResponse(err, activityFetched)
			} else {
				stat.TotalRequest++
				stat.EpochRequest += validPointUnit
			}

			// If the request is invalid, save the invalid response to the database.
			if stat.EpochInvalidRequest > 0 {
				nodeInvalidResponse.Node = stat.Address
				nodeInvalidResponse.Request = stat.Endpoint + "/decentralized/tx/" + activity.ID

				validData, _ := json.Marshal(activity)
				nodeInvalidResponse.VerifierResponse = validData

				if err = e.databaseClient.SaveNodeInvalidResponses(ctx, []*schema.NodeInvalidResponse{nodeInvalidResponse}); err != nil {
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

func generateInvalidResponse(err error, activity *model.ActivityResponse) json.RawMessage {
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

// updateStatsWithResults updates the stats based on the responses.
func updateStatsWithResults(statsMap map[common.Address]*schema.Stat, responses []*model.DataResponse) {
	for _, response := range responses {
		if stat, exists := statsMap[response.Address]; exists {
			stat.TotalRequest += int64(response.ValidPoint)
			stat.EpochRequest += int64(response.ValidPoint)
			stat.EpochInvalidRequest += int64(response.InvalidPoint)
		}
	}
}

// sortResponseByValidity sorts the responses based on the validity.
func sortResponseByValidity(responses []*model.DataResponse) {
	sort.SliceStable(responses, func(i, j int) bool {
		return (responses[i].Err == nil && responses[j].Err != nil) ||
			(responses[i].Err == nil && responses[j].Err == nil && responses[i].Valid && !responses[j].Valid)
	})
}

// updatePointsBasedOnIdentity updates both  based on responses identity.
func updatePointsBasedOnIdentity(responses []*model.DataResponse) {
	errResponseCount := countAndMarkErrorResponse(responses)

	if len(responses) == model.RequiredQualifiedNodeCount-1 {
		handleTwoResponses(responses)
	} else if len(responses) == model.RequiredQualifiedNodeCount-2 {
		handleSingleResponse(responses)
	} else {
		handleFullResponses(responses, errResponseCount)
	}
}

// isResponseIdentical returns true if two byte slices (responses) are identical.
func isResponseIdentical(src, des []byte) bool {
	srcActivity := &model.ActivityResponse{}
	desActivity := &model.ActivityResponse{}

	// check if the data is activity response
	if isDataValid(src, srcActivity) && isDataValid(des, desActivity) {
		if srcActivity.Data == nil && desActivity.Data == nil {
			return true
		} else if srcActivity.Data != nil && desActivity.Data != nil {
			if _, exist := model.MutablePlatformMap[srcActivity.Data.Platform]; !exist {
				return isActivityIdentical(srcActivity.Data, desActivity.Data)
			}

			return true
		}
	}

	srcActivities := &model.ActivitiesResponse{}
	desActivities := &model.ActivitiesResponse{}
	// check if the data is activities response
	if isDataValid(src, srcActivities) && isDataValid(des, desActivities) {
		if srcActivities.Data == nil && desActivities.Data == nil {
			return true
		} else if srcActivities.Data != nil && desActivities.Data != nil {
			// exclude the mutable platforms
			srcActivity, desActivity := excludeMutableActivity(srcActivities.Data), excludeMutableActivity(desActivities.Data)

			for i := range srcActivity {
				if !isActivityIdentical(srcActivity[i], desActivity[i]) {
					return false
				}
			}

			return true
		}
	}

	return false
}

// excludeMutableActivity excludes the mutable platforms from the activities.
func excludeMutableActivity(activities []*model.Activity) []*model.Activity {
	var newActivities []*model.Activity

	for i := range activities {
		if _, exist := model.MutablePlatformMap[activities[i].Platform]; !exist {
			newActivities = append(newActivities, activities[i])
		}
	}

	return newActivities
}

// isActivityIdentical returns true if two Activity are identical.
func isActivityIdentical(src, des *model.Activity) bool {
	if src.ID != des.ID ||
		src.Network != des.Network ||
		src.Index != des.Index ||
		src.From != des.From ||
		src.To != des.To ||
		src.Owner != des.Owner ||
		src.Tag != des.Tag ||
		src.Type != des.Type ||
		src.Platform != des.Platform ||
		len(src.Actions) != len(des.Actions) {
		return false
	}

	// check if the inner actions are identical
	if len(src.Actions) > 0 {
		for i := range des.Actions {
			if src.Actions[i].From != des.Actions[i].From ||
				src.Actions[i].To != des.Actions[i].To ||
				src.Actions[i].Tag != des.Actions[i].Tag ||
				src.Actions[i].Type != des.Actions[i].Type {
				return false
			}
		}
	}

	return true
}

// isDataValid returns true if the data is valid.
func isDataValid(data []byte, target any) bool {
	if err := json.Unmarshal(data, target); err == nil {
		return true
	}

	return false
}

// countAndMarkErrorResponse marks the error results and returns the count.
func countAndMarkErrorResponse(responses []*model.DataResponse) (errResponseCount int) {
	for i := range responses {
		if responses[i].Err != nil {
			responses[i].InvalidPoint = invalidPointUnit
			errResponseCount++
		}
	}

	return errResponseCount
}

// handleTwoResponses handles the case when there are two responses.
func handleTwoResponses(responses []*model.DataResponse) {
	if responses[0].Err == nil && responses[1].Err == nil {
		updateRequestBasedOnComparison(responses)
	} else {
		markErrorResponse(responses...)
	}
}

// handleSingleResponse handles the case when there is only one response.
func handleSingleResponse(responses []*model.DataResponse) {
	if responses[0].Err == nil {
		responses[0].ValidPoint = validPointUnit
	} else {
		responses[0].InvalidPoint = invalidPointUnit
	}
}

// handleFullResponses handles the case when there are more than two results.
func handleFullResponses(responses []*model.DataResponse, errResponseCount int) {
	if errResponseCount < len(responses) {
		compareAndAssignPoints(responses, errResponseCount)
	}
}

// updateRequestBasedOnComparison updates the requests based on the comparison of the data.
func updateRequestBasedOnComparison(responses []*model.DataResponse) {
	if isResponseIdentical(responses[0].Data, responses[1].Data) {
		responses[0].ValidPoint = 2 * validPointUnit
		responses[1].ValidPoint = validPointUnit
	} else {
		responses[0].ValidPoint = validPointUnit
	}
}

// markErrorResponse marks the error responses.
func markErrorResponse(responses ...*model.DataResponse) {
	for i, result := range responses {
		if result.Err != nil {
			responses[i].InvalidPoint = invalidPointUnit
		} else {
			responses[i].ValidPoint = validPointUnit
		}
	}
}

// compareAndAssignPoints compares the data for identity and assigns corresponding points.
func compareAndAssignPoints(responses []*model.DataResponse, errResponseCount int) {
	d0, d1, d2 := responses[0].Data, responses[1].Data, responses[2].Data
	diff01, diff02, diff12 := isResponseIdentical(d0, d1), isResponseIdentical(d0, d2), isResponseIdentical(d1, d2)

	switch errResponseCount {
	// responses contain 2 errors
	case len(responses) - 1:
		responses[0].ValidPoint = validPointUnit
	// responses contain 1 error
	case len(responses) - 2:
		if diff01 {
			responses[0].ValidPoint = 2 * validPointUnit
			responses[1].ValidPoint = validPointUnit
		} else {
			responses[0].ValidPoint = validPointUnit
		}
	// responses contain no errors
	case len(responses) - 3:
		if diff01 && diff02 && diff12 {
			responses[0].ValidPoint = 2 * validPointUnit
			responses[1].ValidPoint = validPointUnit
			responses[2].ValidPoint = validPointUnit
		} else if diff01 && !diff02 {
			responses[0].ValidPoint = 2 * validPointUnit
			responses[1].ValidPoint = validPointUnit
			responses[2].InvalidPoint = invalidPointUnit
		} else if diff02 && !diff01 {
			responses[0].ValidPoint = 2 * validPointUnit
			responses[1].InvalidPoint = invalidPointUnit
			responses[2].ValidPoint = validPointUnit
		} else if diff12 {
			// if the second response is non-null
			if responses[1].Valid {
				responses[0].InvalidPoint = invalidPointUnit
				responses[1].ValidPoint = validPointUnit
				responses[2].ValidPoint = validPointUnit
			} else {
				// the last two responses must include null data
				responses[0].ValidPoint = validPointUnit
			}
		} else if !diff01 {
			responses[0].ValidPoint = validPointUnit
		}
	}
}
