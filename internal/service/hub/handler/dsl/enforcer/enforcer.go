package enforcer

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Enforcer interface {
	VerifyResponses(ctx context.Context, responses []*model.DataResponse) error
	VerifyPartialResponses(ctx context.Context, epochID uint64, responses []*model.DataResponse)
	MaintainReliabilityScore(ctx context.Context) error
	MaintainEpochData(ctx context.Context, epoch int64) error
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

	nodeStatsSlice := lo.MapToSlice(nodeStatsMap, func(_ common.Address, stat *schema.Stat) *schema.Stat {
		return stat
	})
	// save stats to the database
	if err = e.databaseClient.SaveNodeStats(ctx, nodeStatsSlice); err != nil {
		return fmt.Errorf("save Node stats: %w", err)
	}

	// update the score maintainer
	e.batchUpdateScoreMaintainer(ctx, nodeStatsSlice)

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

// MaintainReliabilityScore maintains the Reliability Score σ for all Nodes.
// σ is used to determine the probability of a Node receiving a request on DSL.
func (e *SimpleEnforcer) MaintainReliabilityScore(ctx context.Context) error {
	stats, err := e.getAllNodeStats(ctx, &schema.StatQuery{
		ValidRequest: lo.ToPtr(model.DemotionCountBeforeSlashing),
		Limit:        lo.ToPtr(defaultLimit),
	})
	if err != nil {
		return err
	}

	// Update the stats of the Nodes.
	if err = e.processNodeStats(ctx, stats); err != nil {
		return err
	}

	zap.L().Info("maintain reliability score completed")

	return nil
}

// MaintainEpochData maintains the data for the new epoch.
// The data includes the range of data that all nodes can support in a new epoch.
func (e *SimpleEnforcer) MaintainEpochData(ctx context.Context, epoch int64) error {
	stats, err := e.getAllNodeStats(ctx, &schema.StatQuery{
		Limit: lo.ToPtr(defaultLimit),
	})
	if err != nil {
		return err
	}

	err = e.maintainNodeWorker(ctx, epoch, stats)
	if err != nil {
		return err
	}

	if err = e.processNodeStats(ctx, stats); err != nil {
		return err
	}

	return e.updateNodeCache(ctx, epoch)
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

func NewSimpleEnforcer(ctx context.Context, databaseClient database.Client, cacheClient cache.Client, stakingContract *l2.Staking, httpClient httputil.Client, initCacheData bool) (*SimpleEnforcer, error) {
	enforcer := &SimpleEnforcer{
		databaseClient:  databaseClient,
		cacheClient:     cacheClient,
		stakingContract: stakingContract,
		httpClient:      httpClient,
	}
	if initCacheData {
		if err := enforcer.initWorkerMap(ctx); err != nil {
			return nil, err
		}

		if err := enforcer.initScoreMaintainers(ctx); err != nil {
			return nil, err
		}

		subscribeNodeCacheUpdate(ctx, cacheClient, databaseClient, enforcer.fullNodeScoreMaintainer, enforcer.rssNodeScoreMaintainer)
	}

	return enforcer, nil
}

// subscribeNodeCacheUpdate subscribes to updates of the 'epoch' key.
// Upon updating the 'epoch' key, the Node cache is refreshed.
// This cache holds the initial reliability scores and related maps of the nodes for the new epoch.
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
				var epoch int64
				if err := cacheClient.Get(ctx, model.SubscribeNodeCacheKey, &epoch); err != nil {
					zap.L().Error("get epoch from cache", zap.Error(err))

					continue
				}

				updateQualifiedNodesMap(ctx, model.FullNodeCacheKey, databaseClient, fullNodeScoreMaintainer)
				updateQualifiedNodesMap(ctx, model.RssNodeCacheKey, databaseClient, rssNodeScoreMaintainer)

				zap.L().Info("update qualified nodes map completed", zap.Int64("epoch", epoch))

				errChan := getWorkerMapFromCache(ctx, cacheClient)
				for err := range errChan {
					if err != nil {
						zap.L().Error("get worker map from cache", zap.Error(err), zap.Int64("epoch", epoch))
					}
				}

				zap.L().Info("update worker map completed", zap.Int64("epoch", epoch))
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
