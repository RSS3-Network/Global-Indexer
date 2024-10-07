package distributor

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
)

// A qualified Node has the capability to serve the incoming request.

// getQualifiedNodes retrieves all qualified Nodes from the cache or database.
func (d *Distributor) getQualifiedNodes(ctx context.Context, workers, networks []string) ([]*model.NodeEndpointCache, error) {
	// Match light Nodes.
	lightNodes, err := d.matchLightNodes(ctx, workers, networks)
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
		fullNodes, err := d.simpleEnforcer.RetrieveQualifiedNodes(ctx, model.FullNodeCacheKey)
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
func (d *Distributor) generateQualifiedNodeCache(ctx context.Context, nodeAddresses []common.Address) ([]*model.NodeEndpointCache, error) {
	if len(nodeAddresses) == 0 {
		return nil, nil
	}

	nodesOrderedByPoints, err := d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		Addresses:    nodeAddresses,
		ValidRequest: lo.ToPtr(model.DemotionCountBeforeSlashing),
		Limit:        lo.ToPtr(model.RequiredQualifiedNodeCount),
		PointsOrder:  lo.ToPtr("DESC"),
	})

	if err != nil {
		return nil, err
	}

	nodes := make([]*model.NodeEndpointCache, len(nodesOrderedByPoints))
	for i, stat := range nodesOrderedByPoints {
		nodes[i] = &model.NodeEndpointCache{
			Address:  stat.Address.String(),
			Endpoint: stat.Endpoint,
		}
	}

	return nodes, nil
}

// getFederatedQualifiedNodes retrieves qualified Nodes associated with a federated handle.
func (d *Distributor) getFederatedQualifiedNodes(ctx context.Context, account string) ([]*model.NodeEndpointCache, error) {
	cacheKey := fmt.Sprintf("%s%s", model.FederatedHandlesPrefixCacheKey, account)

	var addresses []string

	if err := d.cacheClient.Get(ctx, cacheKey, &addresses); err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("get federated handles: %w", err)
		}

		var since uint64
		if err = d.cacheClient.Get(ctx, fmt.Sprintf("%s%s", model.FederatedHandlesPrefixCacheKey, "since"), &since); err != nil {
			if !errors.Is(err, redis.Nil) {
				return nil, fmt.Errorf("get federated handles since: %w", err)
			}
		}

		// If the cache is empty, no Nodes are qualified.
		if since > 0 && len(addresses) == 0 {
			return nil, nil
		}
	}

	nodeAddresses := make([]common.Address, len(addresses))
	for i := range addresses {
		nodeAddresses[i] = common.HexToAddress(addresses[i])
	}

	nodeStats, err := d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		Addresses:   nodeAddresses,
		PointsOrder: lo.ToPtr("DESC"),
	})
	if err != nil {
		return nil, fmt.Errorf("find node stats: %w", err)
	}

	nodeEndpointCaches := make([]*model.NodeEndpointCache, len(nodeStats))
	for i, stat := range nodeStats {
		nodeEndpointCaches[i] = &model.NodeEndpointCache{
			Address:     stat.Address.String(),
			Endpoint:    stat.Endpoint,
			AccessToken: stat.AccessToken,
		}
	}

	return nodeEndpointCaches, nil
}
