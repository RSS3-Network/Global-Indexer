package distributor

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/enforcer"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/router"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Distributor struct {
	simpleEnforcer *enforcer.SimpleEnforcer
	simpleRouter   *router.SimpleRouter
	databaseClient database.Client
	cacheClient    cache.Client
}

// RouterRSSHubData routes RSS Hub data retrieval requests.
// It takes a context, path, and query string as input parameters.
// It returns the retrieved data or an error if any occurred.
func (d *Distributor) RouterRSSHubData(ctx context.Context, path, query string) ([]byte, error) {
	nodes, err := d.retrieveNodes(ctx, model.RssNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeMap, err := d.buildRSSHubPath(path, query, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := d.simpleRouter.DistributeRequest(ctx, nodeMap, d.processRSSHubResponses)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

// RouterActivityData routes activity data retrieval requests.
// It takes a context and an activity request as input parameters.
// It returns the retrieved data or an error if any occurred.
func (d *Distributor) RouterActivityData(ctx context.Context, request dsl.ActivityRequest) ([]byte, error) {
	nodes, err := d.retrieveNodes(ctx, model.FullNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeMap, err := d.buildActivityPathByID(request, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := d.simpleRouter.DistributeRequest(ctx, nodeMap, d.processActivityResponses)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

// RouterActivitiesData routes account activities data retrieval requests.
// It takes a context and an account activities request as input parameters.
// It returns the retrieved data or an error if any occurred.
func (d *Distributor) RouterActivitiesData(ctx context.Context, request dsl.AccountActivitiesRequest) ([]byte, error) {
	nodes := make([]model.NodeEndpointCache, 0, model.DefaultNodeCount)

	nodeAddresses, err := d.matchLightNodes(ctx, request)

	if err != nil {
		return nil, err
	}

	if len(nodeAddresses) > 0 {
		nodeStats, err := d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
			AddressList: nodeAddresses,
			Limit:       lo.ToPtr(model.DefaultNodeCount),
			PointsOrder: lo.ToPtr("DESC"),
		})

		if err != nil {
			return nil, err
		}

		num := lo.Ternary(len(nodeStats) > model.DefaultNodeCount, model.DefaultNodeCount, len(nodeStats))

		for i := 0; i < num; i++ {
			nodes = append(nodes, model.NodeEndpointCache{
				Address:  nodeStats[i].Address.String(),
				Endpoint: nodeStats[i].Endpoint,
			})
		}
	}

	if len(nodes) < model.DefaultNodeCount {
		fullNodes, err := d.retrieveNodes(ctx, model.FullNodeCacheKey)
		if err != nil {
			return nil, err
		}

		nodesNeeded := model.DefaultNodeCount - len(nodes)
		nodesToAdd := lo.Ternary(nodesNeeded > len(fullNodes), len(fullNodes), nodesNeeded)

		for i := 0; i < nodesToAdd; i++ {
			nodes = append(nodes, fullNodes[i])
		}
	}

	nodeMap, err := d.buildAccountActivitiesPath(request, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := d.simpleRouter.DistributeRequest(ctx, nodeMap, d.processActivitiesResponses)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

// retrieveNodes retrieves nodes from the cache or database.
// It takes a context and a cache key as input parameters.
// It returns the retrieved nodes or an error if any occurred.
func (d *Distributor) retrieveNodes(ctx context.Context, key string) ([]model.NodeEndpointCache, error) {
	var (
		nodesCache []model.NodeEndpointCache
		nodes      []*schema.Stat
	)

	err := d.cacheClient.Get(ctx, key, &nodesCache)
	if err == nil {
		return nodesCache, nil
	}

	zap.L().Info("not found nodes from cache", zap.String("key", key))

	if errors.Is(err, redis.Nil) {
		switch key {
		case model.RssNodeCacheKey:
			nodes, err = d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
				IsRssNode:    lo.ToPtr(true),
				Limit:        lo.ToPtr(model.DefaultNodeCount),
				ValidRequest: lo.ToPtr(model.DefaultSlashCount),
				PointsOrder:  lo.ToPtr("DESC"),
			})

			if err != nil {
				return nil, err
			}
		case model.FullNodeCacheKey:
			nodes, err = d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
				IsFullNode:   lo.ToPtr(true),
				Limit:        lo.ToPtr(model.DefaultNodeCount),
				ValidRequest: lo.ToPtr(model.DefaultSlashCount),
				PointsOrder:  lo.ToPtr("DESC"),
			})

			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown cache key: %s", key)
		}

		if err = d.setNodeCache(ctx, key, nodes); err != nil {
			return nil, err
		}

		zap.L().Info("set nodes to cache", zap.String("key", key))

		nodesCache = lo.Map(nodes, func(n *schema.Stat, _ int) model.NodeEndpointCache {
			return model.NodeEndpointCache{
				Address:  n.Address.String(),
				Endpoint: n.Endpoint,
			}
		})

		return nodesCache, nil
	}

	return nil, fmt.Errorf("get nodes from cache: %s, %w", key, err)
}

// processRSSHubResults processes the RSS Hub responses.
func (d *Distributor) processRSSHubResponses(responses []*model.DataResponse) {
	if err := d.simpleEnforcer.VerifyResponses(context.Background(), responses); err != nil {
		zap.L().Error("fail to verify rss hub responses", zap.Any("responses", len(responses)))
	} else {
		zap.L().Info("complete rss hub responses verify", zap.Any("responses", len(responses)))
	}
}

// processActivityResults processes activity data retrieval responses.
func (d *Distributor) processActivityResponses(responses []*model.DataResponse) {
	if err := d.simpleEnforcer.VerifyResponses(context.Background(), responses); err != nil {
		zap.L().Error("fail to verify activity id responses ", zap.Any("responses", len(responses)))
	} else {
		zap.L().Info("complete activity id responses verify", zap.Any("responses", len(responses)))
	}
}

// processActivitiesResults processes account activities data retrieval responses.
func (d *Distributor) processActivitiesResponses(responses []*model.DataResponse) {
	ctx := context.Background()

	if err := d.simpleEnforcer.VerifyResponses(ctx, responses); err != nil {
		zap.L().Error("fail to verify activity responses", zap.Any("responses", len(responses)))

		return
	}

	zap.L().Info("complete activity responses verify", zap.Any("responses", len(responses)))

	d.simpleEnforcer.VerifyPartialResponses(ctx, responses)
}

// setNodeCache sets nodes to the cache.
// It takes a context, a cache key, and a slice of stats as input parameters.
// It returns an error if any occurred.
func (d *Distributor) setNodeCache(ctx context.Context, key string, stats []*schema.Stat) error {
	nodesCache := lo.Map(stats, func(n *schema.Stat, _ int) model.NodeEndpointCache {
		return model.NodeEndpointCache{Address: n.Address.String(), Endpoint: n.Endpoint}
	})

	if err := d.cacheClient.Set(ctx, key, nodesCache); err != nil {
		return fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nil
}

// buildActivityPathByID builds the path for activity data retrieval by ID.
// It takes an activity request and a slice of cache nodes as input parameters.
// It returns a map of addresses to URLs or an error if any occurred.
func (d *Distributor) buildActivityPathByID(query dsl.ActivityRequest, nodes []model.NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.simpleRouter.BuildPath(fmt.Sprintf("/decentralized/tx/%s", query.ID), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// buildAccountActivitiesPath builds the path for account activities data retrieval.
// It takes an account activities request and a slice of cache nodes as input parameters.
// It returns a map of addresses to URLs or an error if any occurred.
func (d *Distributor) buildAccountActivitiesPath(query dsl.AccountActivitiesRequest, nodes []model.NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.simpleRouter.BuildPath(fmt.Sprintf("/decentralized/%s", query.Account), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// buildRSSHubPath builds the path for RSS Hub data retrieval.
// It takes a parameter, a query, and a slice of cache nodes as input parameters.
// It returns a map of addresses to URLs or an error if any occurred.
func (d *Distributor) buildRSSHubPath(param, query string, nodes []model.NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.simpleRouter.BuildPath(fmt.Sprintf("/rss/%s?%s", param, query), nil, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// NewDistributor creates a new distributor.
func NewDistributor(_ context.Context, database database.Client, cache cache.Client, httpClient httputil.Client, stakingContract *l2.Staking) *Distributor {
	return &Distributor{
		simpleEnforcer: enforcer.NewSimpleEnforcer(database, cache, stakingContract, httpClient),
		simpleRouter:   router.NewSimpleRouter(httpClient),
		databaseClient: database,
		cacheClient:    cache,
	}
}
