package distributor

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (d *Distributor) getNodes(ctx context.Context, request dsl.AccountActivitiesRequest) ([]model.NodeEndpointCache, error) {
	// Match light nodes.
	nodeAddresses, err := d.matchLightNodes(ctx, request)
	if err != nil {
		return nil, err
	}

	// Generate nodes.
	nodes, err := d.generateNodes(ctx, nodeAddresses)
	if err != nil {
		return nil, err
	}

	// If the number of nodes is less than the default node count, add full nodes.
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

	return nodes, nil
}

func (d *Distributor) generateNodes(ctx context.Context, nodeAddresses []common.Address) ([]model.NodeEndpointCache, error) {
	nodeStats, err := d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: nodeAddresses,
		Limit:       lo.ToPtr(model.DefaultNodeCount),
		PointsOrder: lo.ToPtr("DESC"),
	})

	if err != nil {
		return nil, err
	}

	nodes := make([]model.NodeEndpointCache, len(nodeStats))
	for i, stat := range nodeStats {
		nodes[i] = model.NodeEndpointCache{
			Address:  stat.Address.String(),
			Endpoint: stat.Endpoint,
		}
	}

	return nodes, nil
}

// retrieveNodes retrieves nodes from the cache or database.
// It takes a context and a cache key as input parameters.
// It returns the retrieved nodes or an error if any occurred.
func (d *Distributor) retrieveNodes(ctx context.Context, key string) ([]model.NodeEndpointCache, error) {
	var nodesCache []model.NodeEndpointCache

	if err := d.cacheClient.Get(ctx, key, &nodesCache); err == nil {
		return nodesCache, nil
	} else if !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("get nodes from cache: %s, %w", key, err)
	}

	zap.L().Info("not found nodes from cache", zap.String("key", key))
	nodesCache, err := d.retrieveNodesFromDB(ctx, key)

	if err != nil {
		return nil, err
	}

	if err = d.setNodeCache(ctx, key, nodesCache); err != nil {
		return nil, err
	}

	zap.L().Info("set nodes to cache", zap.String("key", key))

	return nodesCache, nil
}

// retrieveNodesFromDB retrieves nodes from the database.
func (d *Distributor) retrieveNodesFromDB(ctx context.Context, key string) ([]model.NodeEndpointCache, error) {
	var query schema.StatQuery

	switch key {
	case model.RssNodeCacheKey:
		query = schema.StatQuery{IsRssNode: lo.ToPtr(true), Limit: lo.ToPtr(model.DefaultNodeCount), ValidRequest: lo.ToPtr(model.DefaultSlashCount), PointsOrder: lo.ToPtr("DESC")}
	case model.FullNodeCacheKey:
		query = schema.StatQuery{IsFullNode: lo.ToPtr(true), Limit: lo.ToPtr(model.DefaultNodeCount), ValidRequest: lo.ToPtr(model.DefaultSlashCount), PointsOrder: lo.ToPtr("DESC")}
	default:
		return nil, fmt.Errorf("unknown cache key: %s", key)
	}

	nodes, err := d.databaseClient.FindNodeStats(ctx, &query)
	if err != nil {
		return nil, err
	}

	return lo.Map(nodes, func(n *schema.Stat, _ int) model.NodeEndpointCache {
		return model.NodeEndpointCache{Address: n.Address.String(), Endpoint: n.Endpoint}
	}), nil
}

// setNodeCache sets nodes to the cache.
// It takes a context, a cache key, and a slice of stats as input parameters.
// It returns an error if any occurred.
func (d *Distributor) setNodeCache(ctx context.Context, key string, nodesCache []model.NodeEndpointCache) error {
	if err := d.cacheClient.Set(ctx, key, nodesCache); err != nil {
		return fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nil
}
