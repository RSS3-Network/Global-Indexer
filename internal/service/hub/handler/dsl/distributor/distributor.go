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

// RouterRSSHubData routes RSS Hub data retrieval requests.
// It takes a context, path, and query string as input parameters.
// It returns the retrieved data or an error if any occurred.
func (d *Distributor) RouterRSSHubData(ctx context.Context, path, query string) ([]byte, error) {
	nodes, err := d.retrieveNodes(ctx, model.RssNodeCacheKey)

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

// RouterActivityData routes activity data retrieval requests.
// It takes a context and an activity request as input parameters.
// It returns the retrieved data or an error if any occurred.
func (d *Distributor) RouterActivityData(ctx context.Context, request dsl.ActivityRequest) ([]byte, error) {
	nodes, err := d.retrieveNodes(ctx, model.FullNodeCacheKey)

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

// RouterActivitiesData routes account activities data retrieval requests.
// It takes a context and an account activities request as input parameters.
// It returns the retrieved data or an error if any occurred.
func (d *Distributor) RouterActivitiesData(ctx context.Context, request dsl.AccountActivitiesRequest) ([]byte, error) {
	nodes, err := d.getNodes(ctx, request)
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

// generateActivityPathByID builds the path for activity data retrieval by ID.
// It takes an activity request and a slice of cache nodes as input parameters.
// It returns a map of addresses to URLs or an error if any occurred.
func (d *Distributor) generateActivityPathByID(query dsl.ActivityRequest, nodes []model.NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.simpleRouter.BuildPath(fmt.Sprintf("/decentralized/tx/%s", query.ID), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// generateAccountActivitiesPath builds the path for account activities data retrieval.
// It takes an account activities request and a slice of cache nodes as input parameters.
// It returns a map of addresses to URLs or an error if any occurred.
func (d *Distributor) generateAccountActivitiesPath(query dsl.AccountActivitiesRequest, nodes []model.NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.simpleRouter.BuildPath(fmt.Sprintf("/decentralized/%s", query.Account), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// generateRSSHubPath builds the path for RSS Hub data retrieval.
// It takes a parameter, a query, and a slice of cache nodes as input parameters.
// It returns a map of addresses to URLs or an error if any occurred.
func (d *Distributor) generateRSSHubPath(param, query string, nodes []model.NodeEndpointCache) (map[common.Address]string, error) {
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
