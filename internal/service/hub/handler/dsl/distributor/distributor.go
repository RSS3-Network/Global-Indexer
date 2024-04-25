package distributor

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/enforcer"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/router"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"go.uber.org/zap"
)

type Distributor struct {
	simpleEnforcer *enforcer.SimpleEnforcer
	simpleRouter   *router.SimpleRouter
	databaseClient database.Client
	cacheClient    cache.Client
}

// DistributeRSSHubData distributes RSSHub requests to qualified Nodes.
func (d *Distributor) DistributeRSSHubData(ctx context.Context, path, query string) ([]byte, error) {
	nodes, err := d.simpleEnforcer.RetrieveQualifiedNodes(ctx, model.RssNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeMap, err := d.generateRSSHubPath(path, query, nodes)

	if err != nil {
		return nil, err
	}

	nodeResponse, err := d.simpleRouter.DistributeRequest(ctx, nodeMap, d.processRSSHubResponses)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeResponse.Address.String()))

	if nodeResponse.Err != nil {
		return nil, nodeResponse.Err
	}

	return nodeResponse.Data, nil
}

// DistributeActivityRequest distributes Activity requests to qualified Nodes.
func (d *Distributor) DistributeActivityRequest(ctx context.Context, request dsl.ActivityRequest) ([]byte, error) {
	nodes, err := d.simpleEnforcer.RetrieveQualifiedNodes(ctx, model.FullNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeMap, err := d.generateActivityPathByID(request, nodes)

	if err != nil {
		return nil, err
	}

	nodeResponse, err := d.simpleRouter.DistributeRequest(ctx, nodeMap, d.processActivityResponses)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeResponse.Address.String()))

	if nodeResponse.Err != nil {
		return nil, nodeResponse.Err
	}

	return nodeResponse.Data, nil
}

// DistributeActivitiesData distributes Activities requests to qualified Nodes.
func (d *Distributor) DistributeActivitiesData(ctx context.Context, request dsl.ActivitiesRequest) ([]byte, error) {
	nodes, err := d.getQualifiedNodes(ctx, request)
	if err != nil {
		return nil, err
	}

	nodeMap, err := d.generateAccountActivitiesPath(request, nodes)

	if err != nil {
		return nil, err
	}

	nodeResponse, err := d.simpleRouter.DistributeRequest(ctx, nodeMap, d.processActivitiesResponses)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeResponse.Address.String()))

	if nodeResponse.Err != nil {
		return nil, nodeResponse.Err
	}

	return nodeResponse.Data, nil
}

// generateActivityPathByID builds the path for Activity requests.
func (d *Distributor) generateActivityPathByID(query dsl.ActivityRequest, nodes []*model.NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.simpleRouter.BuildPath(fmt.Sprintf("/decentralized/tx/%s", query.ID), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// generateAccountActivitiesPath builds the path for Activities requests.
func (d *Distributor) generateAccountActivitiesPath(query dsl.ActivitiesRequest, nodes []*model.NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.simpleRouter.BuildPath(fmt.Sprintf("/decentralized/%s", query.Account), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// generateRSSHubPath builds the path for RSSHub requests.
func (d *Distributor) generateRSSHubPath(param, query string, nodes []*model.NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.simpleRouter.BuildPath(fmt.Sprintf("/rss/%s?%s", param, query), nil, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// NewDistributor creates a new distributor.
func NewDistributor(_ context.Context, database database.Client, cache cache.Client, httpClient httputil.Client, stakingContract *l2.Staking) (*Distributor, error) {
	simpleEnforcer, err := enforcer.NewSimpleEnforcer(database, cache, stakingContract, httpClient, true)

	if err != nil {
		return nil, err
	}

	return &Distributor{
		simpleEnforcer: simpleEnforcer,
		simpleRouter:   router.NewSimpleRouter(httpClient),
		databaseClient: database,
		cacheClient:    cache,
	}, nil
}
