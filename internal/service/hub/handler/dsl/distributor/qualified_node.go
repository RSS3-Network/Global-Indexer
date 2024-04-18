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

// A qualified Node has the capability to serve the incoming request.

// getQualifiedNodes retrieves all qualified Nodes from the cache or database.
func (d *Distributor) getQualifiedNodes(ctx context.Context, request dsl.ActivitiesRequest) ([]model.NodeEndpointCache, error) {
	// Match light Nodes.
	lightNodes, err := d.matchLightNodes(ctx, request)
	if err != nil {
		return nil, err
	}

	// Order Nodes and generate a cache.
	qualifiedNodeCache, err := d.generateQualifiedNodeCache(ctx, lightNodes)
	if err != nil {
		return nil, err
	}

	// Calculate the number of Nodes that still need to be added
	nodesNeeded := model.RequiredQualifiedNodeCount - len(qualifiedNodeCache)

	if nodesNeeded > 0 {
		// retrieve additional full Nodes.
		fullNodes, err := d.retrieveQualifiedNodes(ctx, model.FullNodeCacheKey)
		if err != nil {
			return nil, err
		}

		if nodesNeeded > len(fullNodes) {
			nodesNeeded = len(fullNodes)
		}

		// Append the required number of Nodes from fullNodes to qualifiedNodeCache
		qualifiedNodeCache = append(qualifiedNodeCache, fullNodes[:nodesNeeded]...)
	}

	return qualifiedNodeCache, nil
}

// generateQualifiedNodeCache generates an ordered qualified Node cache
func (d *Distributor) generateQualifiedNodeCache(ctx context.Context, nodeAddresses []common.Address) ([]model.NodeEndpointCache, error) {
	nodesOrderedByPoints, err := d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: nodeAddresses,
		Limit:       lo.ToPtr(model.RequiredQualifiedNodeCount),
		PointsOrder: lo.ToPtr("DESC"),
	})

	if err != nil {
		return nil, err
	}

	nodes := make([]model.NodeEndpointCache, len(nodesOrderedByPoints))
	for i, stat := range nodesOrderedByPoints {
		nodes[i] = model.NodeEndpointCache{
			Address:  stat.Address.String(),
			Endpoint: stat.Endpoint,
		}
	}

	return nodes, nil
}

// retrieveQualifiedNodes retrieves qualified Nodes from the cache or database.
func (d *Distributor) retrieveQualifiedNodes(ctx context.Context, key string) ([]model.NodeEndpointCache, error) {
	var nodesCache []model.NodeEndpointCache

	if err := d.cacheClient.Get(ctx, key, &nodesCache); err == nil {
		return nodesCache, nil
	} else if !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("get nodes from cache: %s, %w", key, err)
	}

	zap.L().Info("nodes not in cache", zap.String("key", key))
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
		query = schema.StatQuery{IsRssNode: lo.ToPtr(true), Limit: lo.ToPtr(model.RequiredQualifiedNodeCount), ValidRequest: lo.ToPtr(model.DefaultSlashCount), PointsOrder: lo.ToPtr("DESC")}
	case model.FullNodeCacheKey:
		query = schema.StatQuery{IsFullNode: lo.ToPtr(true), Limit: lo.ToPtr(model.RequiredQualifiedNodeCount), ValidRequest: lo.ToPtr(model.DefaultSlashCount), PointsOrder: lo.ToPtr("DESC")}
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

// setNodeCache sets Nodes to the cache.
func (d *Distributor) setNodeCache(ctx context.Context, key string, nodesCache []model.NodeEndpointCache) error {
	if err := d.cacheClient.Set(ctx, key, nodesCache); err != nil {
		return fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nil
}
