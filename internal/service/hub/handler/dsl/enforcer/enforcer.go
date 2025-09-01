package enforcer

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	v2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Enforcer interface {
	VerifyResponses(ctx context.Context, responses []*model.DataResponse, verify bool) error
	VerifyPartialResponses(ctx context.Context, epochID uint64, responses []*model.DataResponse)
	MaintainReliabilityScore(ctx context.Context) error
	MaintainEpochData(ctx context.Context, epoch int64) error
	ChallengeStates(ctx context.Context) error
	RetrieveQualifiedNodes(ctx context.Context, key string) ([]*model.NodeEndpointCache, error)
}

type SimpleEnforcer struct {
	cacheClient               cache.Client
	databaseClient            database.Client
	httpClient                httputil.Client
	stakingContract           *l2.StakingV2MulticallClient
	networkParamsContract     *l2.NetworkParams
	fullNodeScoreMaintainer   *ScoreMaintainer
	rssNodeScoreMaintainer    *ScoreMaintainer
	aiNodeScoreMaintainer     *ScoreMaintainer
	rsshubNodeScoreMaintainer *ScoreMaintainer
	txManager                 txmgr.TxManager
	settlerConfig             *config.Settler
	chainID                   *big.Int
}

// VerifyResponses verifies the responses from the Nodes.
func (e *SimpleEnforcer) VerifyResponses(ctx context.Context, responses []*model.DataResponse, verify bool) error {
	if len(responses) == 0 {
		return fmt.Errorf("no response returned from nodes")
	}

	// non-error and non-null results are always put in front of the list
	sortResponseByValidity(responses)

	if verify {
		// update requests based on data compare
		updatePointsBasedOnIdentity(responses)
	} else {
		// update requests based on data
		updatePointsBasedOnData(responses)
	}
	// update the cache request
	e.updateCacheRequest(ctx, responses)
	// update the score maintainer
	e.batchUpdateScoreMaintainer(ctx, responses)

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
	if err = e.processNodeStats(ctx, stats, false); err != nil {
		return err
	}

	zap.L().Info("maintain reliability score completed")

	return nil
}

// MaintainEpochData maintains the data for the new epoch.
// The data includes the range of data that all nodes can support and status of nodes in a new epoch.
func (e *SimpleEnforcer) MaintainEpochData(ctx context.Context, epoch int64) error {
	stats, err := e.getAllNodeStats(ctx, &schema.StatQuery{})
	if err != nil {
		return err
	}

	// Separate DSL and RSSHub nodes
	dslNodeStats := lo.Filter(stats, func(stat *schema.Stat, _ int) bool {
		return !stat.IsRsshubNode
	})
	rsshubNodeStats := lo.Filter(stats, func(stat *schema.Stat, _ int) bool {
		return stat.IsRsshubNode
	})

	// Process DSL nodes
	dslNodeAddressList, dslNodeStatusList, err := e.maintainNodeWorker(ctx, epoch, dslNodeStats)
	if err != nil {
		return err
	}

	// Process RSSHub nodes
	rsshubNodeAddressList, rsshubNodeStatusList, err := e.maintainRSSHubNodes(ctx, epoch, rsshubNodeStats)
	if err != nil {
		return err
	}

	// Merge all node statuses that need to be updated
	allNodeAddressList := append(rsshubNodeAddressList, dslNodeAddressList...)
	allNodeStatusList := append(rsshubNodeStatusList, dslNodeStatusList...)

	// Batch update node statuses
	if len(allNodeAddressList) > 0 {
		if err = e.updateNodeStatusAndSubmitDemotionToVSL(ctx, allNodeAddressList, allNodeStatusList, nil, nil, nil); err != nil {
			return err
		}
	}

	// Process node statistics and update cache
	if err = e.processNodeStats(ctx, append(dslNodeStats, rsshubNodeStats...), true); err != nil {
		return err
	}

	return e.updateNodeCache(ctx, epoch)
}

// maintainRSSHubNodes handles the status maintenance for RSSHub nodes
func (e *SimpleEnforcer) maintainRSSHubNodes(ctx context.Context, epoch int64, rsshubNodeStats []*schema.Stat) ([]common.Address, []uint8, error) {
	if len(rsshubNodeStats) == 0 {
		return nil, nil, nil
	}

	// Get node address list
	addresses := lo.Map(rsshubNodeStats, func(stat *schema.Stat, _ int) common.Address {
		return stat.Address
	})

	// Get node information from chain
	var nodeVSLInfo []v2.Node
	nodeVSLInfo, err := e.stakingContract.GetNodes(&bind.CallOpts{}, addresses)

	if err != nil {
		return nil, nil, err
	}

	var (
		nodeAddressList []common.Address
		nodeStatusList  []uint8
	)

	// Process each RSSHub node concurrently with limited concurrency
	// Limit concurrent network requests
	const maxConcurrency = 20

	var mu sync.Mutex

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrency)

	for i, stat := range rsshubNodeStats {
		i, stat := i, stat

		g.Go(func() error {
			// Set epoch data
			stat.Epoch = epoch
			stat.EpochRequest = 0

			// Check node status concurrently
			status, err := e.getRSSHubNodeStatus(gCtx, stat.Endpoint, stat.AccessToken)
			if err != nil {
				// Continue processing even if status check fails
				zap.L().Error("get RSSHub node status", zap.Error(err), zap.String("endpoint", stat.Endpoint), zap.String("address", stat.Address.String()))
			}

			mu.Lock()
			defer mu.Unlock()

			if err != nil || !status {
				// Node is offline
				stat.EpochInvalidRequest = int64(model.DemotionCountBeforeSlashing)
				if schema.NodeStatus(nodeVSLInfo[i].Status) != schema.NodeStatusOffline {
					nodeAddressList = append(nodeAddressList, stat.Address)
					nodeStatusList = append(nodeStatusList, uint8(schema.NodeStatusOffline))
				}
			} else {
				// Node is online
				stat.EpochInvalidRequest = 0
				if schema.NodeStatus(nodeVSLInfo[i].Status) != schema.NodeStatusOnline {
					nodeAddressList = append(nodeAddressList, stat.Address)
					nodeStatusList = append(nodeStatusList, uint8(schema.NodeStatusOnline))
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	return nodeAddressList, nodeStatusList, nil
}

func (e *SimpleEnforcer) MaintainNodeStatus(ctx context.Context) error {
	if err := e.maintainNodeStatus(ctx); err != nil {
		return err
	}

	zap.L().Info("maintain node status completed")

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
	case model.RsshubNodeCacheKey:
		nodesCache, err = e.rsshubNodeScoreMaintainer.retrieveQualifiedNodes(ctx, key, model.RequiredQualifiedNodeCount)
	case model.RssNodeCacheKey:
		nodesCache, err = e.rssNodeScoreMaintainer.retrieveQualifiedNodes(ctx, key, model.RequiredQualifiedNodeCount)
	case model.FullNodeCacheKey:
		nodesCache, err = e.fullNodeScoreMaintainer.retrieveQualifiedNodes(ctx, key, model.RequiredQualifiedNodeCount)
	case model.AINodeCacheKey:
		nodesCache, err = e.aiNodeScoreMaintainer.retrieveQualifiedNodes(ctx, key, model.RequiredQualifiedNodeCount)
	default:
		return nil, fmt.Errorf("unknown cache key: %s", key)
	}

	// TODO: If there are no qualified nodes, how should the request be handled
	if len(nodesCache) == 0 {
		return nil, fmt.Errorf("no qualified nodes in the current epoch")
	}

	return nodesCache, err
}

func NewSimpleEnforcer(ctx context.Context, databaseClient database.Client, cacheClient cache.Client, stakingContract *l2.StakingV2MulticallClient, networkParamsContract *l2.NetworkParams, httpClient httputil.Client, txManager *txmgr.SimpleTxManager, settlerConfig *config.Settler, chainID *big.Int, initCacheData bool) (*SimpleEnforcer, error) {
	enforcer := &SimpleEnforcer{
		databaseClient:        databaseClient,
		cacheClient:           cacheClient,
		stakingContract:       stakingContract,
		networkParamsContract: networkParamsContract,
		httpClient:            httpClient,
		txManager:             txManager,
		settlerConfig:         settlerConfig,
		chainID:               chainID,
	}

	if initCacheData {
		if err := enforcer.initWorkerMap(ctx); err != nil {
			return nil, err
		}

		if err := enforcer.initScoreMaintainers(ctx); err != nil {
			return nil, err
		}

		subscribeNodeCacheUpdate(ctx, cacheClient, databaseClient, enforcer.fullNodeScoreMaintainer, enforcer.rssNodeScoreMaintainer, enforcer.aiNodeScoreMaintainer, enforcer.rsshubNodeScoreMaintainer)
	}

	return enforcer, nil
}

// subscribeNodeCacheUpdate subscribes to updates of the 'epoch' key.
// Upon updating the 'epoch' key, the Node cache is refreshed.
// This cache holds the initial reliability scores and related maps of the nodes for the new epoch.
func subscribeNodeCacheUpdate(ctx context.Context, cacheClient cache.Client, databaseClient database.Client, fullNodeScoreMaintainer, rssNodeScoreMaintainer, aiNodeScoreMaintainer, rsshubNodeScoreMaintainer *ScoreMaintainer) {
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
				updateQualifiedNodesMap(ctx, model.AINodeCacheKey, databaseClient, aiNodeScoreMaintainer)
				updateQualifiedNodesMap(ctx, model.RsshubNodeCacheKey, databaseClient, rsshubNodeScoreMaintainer)

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
	nodes, err := retrieveNodeStatsFromDB(ctx, key, databaseClient)
	if err != nil {
		zap.L().Error("get nodes from db", zap.Error(err))
	}

	if err = scoreMaintainer.updateQualifiedNodesMap(ctx, nodes); err != nil {
		zap.L().Error("update qualified nodes map", zap.Error(err))
	}
}

// initScoreMaintainers initializes the score maintainers for the full and rss nodes.
func (e *SimpleEnforcer) initScoreMaintainers(ctx context.Context) error {
	var err error
	if e.fullNodeScoreMaintainer, err = e.initScoreMaintainer(ctx, model.FullNodeCacheKey); err != nil {
		return err
	}

	if e.rssNodeScoreMaintainer, err = e.initScoreMaintainer(ctx, model.RssNodeCacheKey); err != nil {
		return err
	}

	if e.aiNodeScoreMaintainer, err = e.initScoreMaintainer(ctx, model.AINodeCacheKey); err != nil {
		return err
	}

	if e.rsshubNodeScoreMaintainer, err = e.initScoreMaintainer(ctx, model.RsshubNodeCacheKey); err != nil {
		return err
	}

	return nil
}

// initScoreMaintainer initializes the score maintainer for the given node type.
func (e *SimpleEnforcer) initScoreMaintainer(ctx context.Context, nodeType string) (*ScoreMaintainer, error) {
	// Retrieve the node stats from the database.
	nodeStats, err := retrieveNodeStatsFromDB(ctx, nodeType, e.databaseClient)
	if err != nil {
		return nil, err
	}

	return newScoreMaintainer(ctx, nodeType, nodeStats, e.cacheClient)
}
