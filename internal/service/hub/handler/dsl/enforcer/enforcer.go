package enforcer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/node/schema/worker"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Enforcer interface {
	VerifyResponses(ctx context.Context, responses []*model.DataResponse) error
	VerifyPartialResponses(ctx context.Context, epochID uint64, responses []*model.DataResponse)
	MaintainReliabilityScore(ctx context.Context) error
	ChallengeStates(ctx context.Context) error
	RetrieveQualifiedNodes(ctx context.Context, key string) ([]*model.NodeEndpointCache, error)
}

type SimpleEnforcer struct {
	cacheClient             cache.Client
	databaseClient          database.Client
	httpClient              httputil.Client
	stakingContract         *l2.Staking
	fullNodeScoreMaintainer *ScoreMaintainer
	rssNodeScoreMaintainer  *ScoreMaintainer
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

	// update the score maintainer
	e.updateScoreMaintainer(ctx, nodeStatsMap)

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

func updateStatsWithResults(statsMap map[common.Address]*schema.Stat, responses []*model.DataResponse) {
	for _, response := range responses {
		if stat, exists := statsMap[response.Address]; exists {
			stat.TotalRequest += int64(response.ValidPoint)
			stat.EpochRequest += int64(response.ValidPoint)
			stat.EpochInvalidRequest += int64(response.InvalidPoint)
		}
	}
}

func (e *SimpleEnforcer) updateScoreMaintainer(ctx context.Context, nodeStatsMap map[common.Address]*schema.Stat) {
	for _, stat := range nodeStatsMap {
		_ = CalculateReliabilityScore(stat)

		nodeCache := &model.NodeEndpointCache{
			Address:      stat.Address.String(),
			Score:        stat.Score,
			Endpoint:     stat.Endpoint,
			InvalidCount: stat.EpochInvalidRequest,
		}

		if err := e.fullNodeScoreMaintainer.addOrUpdateScore(ctx, model.FullNodeCacheKey, nodeCache); err != nil {
			zap.L().Error("failed to update full node score", zap.Error(err))
		}

		if err := e.rssNodeScoreMaintainer.addOrUpdateScore(ctx, model.RssNodeCacheKey, nodeCache); err != nil {
			zap.L().Error("failed to update rss node score", zap.Error(err))
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
	pid, err := worker.PlatformString(activity.Platform)
	if err != nil {
		return nil, err
	}

	workers := model.PlatformToWorkersMap[pid]
	indexers, err := e.databaseClient.FindNodeWorkers(ctx, nil, []string{activity.Network}, workers)

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

// MaintainReliabilityScore maintains the Reliability Score σ for all Nodes.
// σ is used to determine the probability of a Node receiving a request on DSL.
func (e *SimpleEnforcer) MaintainReliabilityScore(ctx context.Context) error {
	// Retrieve the most recently indexed epoch.
	currentEpoch, err := e.getCurrentEpoch(ctx)
	if err != nil {
		return err
	}

	var lastStatEpoch int64

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

		// A nil cursor indicates that the stats represent the initial batch of data.
		if query.Cursor == nil {
			lastStatEpoch = stats[0].Epoch
		}

		if err = e.processNodeStats(ctx, stats, currentEpoch); err != nil {
			return err
		}

		query.Cursor = lo.ToPtr(stats[len(stats)-1].Address.String())
	}

	// If the epoch of the current stat differs from that of the first stat,
	// it indicates an epoch change, necessitating a notification to the score queue.
	if currentEpoch != lastStatEpoch {
		if err = e.updateNodeCache(ctx, currentEpoch); err != nil {
			return err
		}

		zap.L().Info("update node cache for the latest epoch", zap.Int64("epoch", currentEpoch))
	}

	zap.L().Info("maintain reliability score completed")

	return nil
}

func (e *SimpleEnforcer) ChallengeStates(_ context.Context) error {
	return nil
}

// RetrieveQualifiedNodes retrieves the qualified Nodes from the sorted set.
func (e *SimpleEnforcer) RetrieveQualifiedNodes(ctx context.Context, key string) ([]*model.NodeEndpointCache, error) {
	var (
		nodesCache []*model.NodeEndpointCache
		err        error
	)

	switch key {
	case model.RssNodeCacheKey:
		nodesCache, err = e.rssNodeScoreMaintainer.retrieveQualifiedNodes(ctx, key, model.RequiredQualifiedNodeCount)
	case model.FullNodeCacheKey:
		nodesCache, err = e.fullNodeScoreMaintainer.retrieveQualifiedNodes(ctx, key, model.RequiredQualifiedNodeCount)
	default:
		return nil, fmt.Errorf("unknown cache key: %s", key)
	}

	// TODO: If there are no qualified nodes, how should the request be handled
	if len(nodesCache) == 0 {
		return nil, fmt.Errorf("no qualified nodes in the current epoch")
	}

	return nodesCache, err
}

// subscribeNodeCacheUpdate subscribes to updates of the 'epoch' key.
// Upon updating the 'epoch' key, the Node cache is refreshed.
// This cache holds the initial reliability scores of the nodes for the new epoch.
func subscribeNodeCacheUpdate(ctx context.Context, cacheClient cache.Client, databaseClient database.Client, fullNodeScoreMaintainer, rssNodeScoreMaintainer *ScoreMaintainer) {
	go func() {
		//Subscribe to changes to 'epoch'
		pubsub := cacheClient.PSubscribe(ctx, fmt.Sprintf("__keyspace@*__:%s", model.SubscribeNodeCacheKey))
		defer pubsub.Close()

		// Wait for confirmation that subscription is created before proceeding.
		if _, err := pubsub.Receive(ctx); err != nil {
			zap.L().Error("subscribe node cache failed:", zap.Error(err))

			os.Exit(1)
		}

		// Go channel to receive messages from Redis
		ch := pubsub.Channel()

		zap.L().Info("start listening to 'epoch'...")

		// A message is received whenever the 'epoch' key is updated, indicating the start of a new epoch.
		for msg := range ch {
			zap.L().Info("received message from channel", zap.String("channel", msg.Channel), zap.String("payload", msg.Payload))

			if msg.Payload == "set" {
				updateQualifiedNodesMap(ctx, model.FullNodeCacheKey, databaseClient, fullNodeScoreMaintainer)
				updateQualifiedNodesMap(ctx, model.RssNodeCacheKey, databaseClient, rssNodeScoreMaintainer)
			}
		}
	}()
}

// updateQualifiedNodesMap retrieves the node endpoint caches from the database and updates the score maintainer's map of qualified nodes.
func updateQualifiedNodesMap(ctx context.Context, key string, databaseClient database.Client, scoreMaintainer *ScoreMaintainer) {
	nodes, err := retrieveNodeEndpointCaches(ctx, key, databaseClient)
	if err != nil {
		zap.L().Error("get nodes from db", zap.Error(err))
	}

	scoreMaintainer.updateQualifiedNodesMap(nodes)
}

func (e *SimpleEnforcer) initScoreMaintainers(ctx context.Context) error {
	var err error
	if e.fullNodeScoreMaintainer, err = e.initScoreMaintainer(ctx, model.FullNodeCacheKey); err != nil {
		return err
	}

	if e.rssNodeScoreMaintainer, err = e.initScoreMaintainer(ctx, model.RssNodeCacheKey); err != nil {
		return err
	}

	return nil
}

func (e *SimpleEnforcer) initScoreMaintainer(ctx context.Context, nodeType string) (*ScoreMaintainer, error) {
	nodes, err := retrieveNodeEndpointCaches(ctx, nodeType, e.databaseClient)

	if err != nil {
		return nil, err
	}

	return newScoreMaintainer(ctx, nodeType, nodes, e.cacheClient)
}

func NewSimpleEnforcer(ctx context.Context, databaseClient database.Client, cacheClient cache.Client, stakingContract *l2.Staking, httpClient httputil.Client, initScoreMaintainer bool) (*SimpleEnforcer, error) {
	enforcer := &SimpleEnforcer{
		databaseClient:  databaseClient,
		cacheClient:     cacheClient,
		stakingContract: stakingContract,
		httpClient:      httpClient,
	}

	if initScoreMaintainer {
		if err := enforcer.initScoreMaintainers(ctx); err != nil {
			return nil, err
		}

		subscribeNodeCacheUpdate(ctx, cacheClient, databaseClient, enforcer.fullNodeScoreMaintainer, enforcer.rssNodeScoreMaintainer)
	}

	return enforcer, nil
}
