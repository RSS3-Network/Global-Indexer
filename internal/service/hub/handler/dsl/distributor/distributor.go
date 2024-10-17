package distributor

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/config"
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

// generateRSSHubPath builds the path for RSSHub requests.
func (d *Distributor) generateRSSHubPath(param, query string, nodes []*model.NodeEndpointCache) (map[common.Address]model.RequestMeta, error) {
	endpointMap, err := d.simpleRouter.BuildPath(http.MethodGet, fmt.Sprintf("/rss/%s?%s", param, query), nil, nodes, nil)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

type nodeRetriever func(ctx context.Context, workers, networks []string) ([]*model.NodeEndpointCache, error)
type responseProcessor func([]*model.DataResponse)

// DistributeData distributes requests to qualified Nodes.
func (d *Distributor) DistributeData(ctx context.Context, requestType, component string, request interface{}, params url.Values, workers, networks []string) ([]byte, error) {
	retriever, processor, err := d.getStrategyForRequest(requestType, component, request)
	if err != nil {
		return nil, fmt.Errorf("get strategy for request: %w", err)
	}

	nodes, err := retriever(ctx, workers, networks)
	if err != nil {
		return nil, fmt.Errorf("retrieving nodes: %w", err)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes available")
	}

	nodeMap, err := d.generatePath(requestType, component, request, params, nodes)
	if err != nil {
		return nil, fmt.Errorf("generate path: %w", err)
	}

	nodeResponse, err := d.simpleRouter.DistributeRequest(ctx, nodeMap, processor)
	if err != nil {
		return nil, fmt.Errorf("distribute request: %w", err)
	}

	zap.L().Info("first node return", zap.Any("address", nodeResponse.Address.String()))

	if nodeResponse.Err != nil {
		return nil, nodeResponse.Err
	}

	return nodeResponse.Data, nil
}

// getStrategyForRequest returns the node retriever and response processor for the request.
func (d *Distributor) getStrategyForRequest(requestType, component string, request interface{}) (nodeRetriever, responseProcessor, error) {
	switch requestType {
	case model.DistributorRequestActivity:
		return d.getActivityStrategy(component, request)
	case model.DistributorRequestAccountActivities,
		model.DistributorRequestBatchAccountActivities,
		model.DistributorRequestNetworkActivities,
		model.DistributorRequestPlatformActivities:
		return d.getAccountActivitiesStrategy(component, request)
	default:
		return nil, nil, fmt.Errorf("invalid request type: %s", requestType)
	}
}

// getActivityStrategy returns the node retriever and response processor for activity requests.
func (d *Distributor) getActivityStrategy(component string, request interface{}) (nodeRetriever, responseProcessor, error) {
	switch component {
	case model.ComponentDecentralized:
		return func(ctx context.Context, _, _ []string) ([]*model.NodeEndpointCache, error) {
			return d.simpleEnforcer.RetrieveQualifiedNodes(ctx, model.FullNodeCacheKey)
		}, d.processDecentralizedActivityResponses, nil
	case model.ComponentFederated:
		return d.getFederatedNodeRetriever(request), d.processFederatedActivityResponses, nil
	default:
		return nil, nil, fmt.Errorf("invalid component: %s", component)
	}
}

// getAccountActivitiesStrategy returns the node retriever and response processor for account activities requests.
func (d *Distributor) getAccountActivitiesStrategy(component string, request interface{}) (nodeRetriever, responseProcessor, error) {
	switch component {
	case model.ComponentDecentralized:
		return d.getQualifiedNodes, d.processDecentralizedActivitiesResponses, nil
	case model.ComponentFederated:
		return d.getFederatedNodeRetriever(request), d.processFederatedActivitiesResponses, nil
	default:
		return nil, nil, fmt.Errorf("invalid component: %s", component)
	}
}

// getFederatedNodeRetriever returns a node retriever for federated requests.
func (d *Distributor) getFederatedNodeRetriever(request interface{}) nodeRetriever {
	return func(ctx context.Context, _, _ []string) ([]*model.NodeEndpointCache, error) {
		var account string

		switch r := request.(type) {
		case dsl.ActivityRequest:
			var err error
			r.ID, err = url.QueryUnescape(r.ID)

			if err != nil {
				return nil, fmt.Errorf("unescape ID: %w", err)
			}

			account, err = extractFediverseAddress(r.ID)

			if err != nil {
				return nil, fmt.Errorf("extract fediverse address: %w", err)
			}
		case dsl.ActivitiesRequest:
			var err error
			account, err = url.QueryUnescape(r.Account)

			if err != nil {
				return nil, fmt.Errorf("unescape account: %w", err)
			}
		default:
			return d.getFederatedDefaultNodes(ctx)
		}

		return d.getFederatedQualifiedNodes(ctx, account)
	}
}

// extractFediverseAddress extracts the fediverse address from an activity URL.
func extractFediverseAddress(activityURL string) (string, error) {
	// parse url
	parsedURL, err := url.Parse(activityURL)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}

	// split path
	pathParts := strings.Split(parsedURL.Path, "/")
	if len(pathParts) < 3 {
		return "", fmt.Errorf("invalid url path")
	}

	// extract username and domain
	username := pathParts[2]
	domain := parsedURL.Host

	// construct fediverse address
	fediverseAddress := fmt.Sprintf("@%s@%s", username, domain)

	return fediverseAddress, nil
}

// generatePath builds the path for distributor requests.
func (d *Distributor) generatePath(requestType, component string, request interface{}, params url.Values, nodes []*model.NodeEndpointCache) (map[common.Address]model.RequestMeta, error) {
	var (
		path   string
		method = http.MethodGet

		body []byte
		err  error
	)

	switch req := request.(type) {
	case dsl.ActivityRequest:
		path = fmt.Sprintf("/%s/tx/%s", component, req.ID)
	case dsl.ActivitiesRequest:
		path = fmt.Sprintf("/%s/%s", component, req.Account)
	case dsl.AccountsActivitiesRequest:
		path = "/" + component + "/accounts"
		method = http.MethodPost
		body, err = json.Marshal(req)

		if err != nil {
			return nil, fmt.Errorf("marshal request data: %w", err)
		}
	case dsl.NetworkActivitiesRequest:
		path = fmt.Sprintf("/%s/network/%s", component, req.Network)
	case dsl.PlatformActivitiesRequest:
		path = fmt.Sprintf("/%s/platform/%s", component, req.Platform)
	default:
		return nil, fmt.Errorf("invalid request type: %s", requestType)
	}

	endpointMap, err := d.simpleRouter.BuildPath(method, path, params, nodes, body)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// NewDistributor creates a new distributor.
func NewDistributor(ctx context.Context, database database.Client, cache cache.Client, httpClient httputil.Client, stakingContract *l2.StakingV2MulticallClient, networkParamsContract *l2.NetworkParams, txManager *txmgr.SimpleTxManager, settlerConfig *config.Settler, chainID *big.Int) (*Distributor, error) {
	simpleEnforcer, err := enforcer.NewSimpleEnforcer(ctx, database, cache, stakingContract, networkParamsContract, httpClient, txManager, settlerConfig, chainID, true)

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
