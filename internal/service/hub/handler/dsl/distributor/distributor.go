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

// DistributeData distributes requests to qualified Nodes.
func (d *Distributor) DistributeData(ctx context.Context, requestType, component string, request interface{}, params url.Values, workers, networks []string) ([]byte, error) {
	var (
		nodes          []*model.NodeEndpointCache
		processResults = d.processDecentralizedActivitiesResponses

		err error
	)

	switch requestType {
	case model.DistributorRequestActivity:
		if component == model.ComponentDecentralized {
			nodes, err = d.simpleEnforcer.RetrieveQualifiedNodes(ctx, model.FullNodeCacheKey)
			processResults = d.processDecentralizedActivityResponses
		}

		if component == model.ComponentFederated {
			id := request.(dsl.ActivityRequest).ID
			account, err := extractFediverseAddress(id)
			if err != nil {
				return nil, fmt.Errorf("extract fediverse address: %w", err)
			}
			nodes, err = d.getFederatedQualifiedNodes(ctx, account)
			processResults = d.processFederatedActivityResponses
		}
	case model.DistributorRequestAccountActivities:
		if component == model.ComponentDecentralized {
			nodes, err = d.getQualifiedNodes(ctx, workers, networks)
		}

		if component == model.ComponentFederated {
			account := request.(dsl.ActivitiesRequest).Account
			nodes, err = d.getFederatedQualifiedNodes(ctx, account)
			processResults = d.processFederatedActivitiesResponses
		}
	case model.DistributorRequestBatchAccountActivities:
		nodes, err = d.getQualifiedNodes(ctx, workers, networks)
	case model.DistributorRequestNetworkActivities:
		nodes, err = d.getQualifiedNodes(ctx, workers, networks)
	case model.DistributorRequestPlatformActivities:
		nodes, err = d.getQualifiedNodes(ctx, workers, networks)
	default:
		return nil, fmt.Errorf("invalid request type: %s", requestType)
	}

	if err != nil {
		return nil, fmt.Errorf("get qualified nodes: %w", err)
	}

	nodeMap, err := d.generatePath(requestType, component, request, params, nodes)
	if err != nil {
		return nil, fmt.Errorf("generate path: %w", err)
	}

	nodeResponse, err := d.simpleRouter.DistributeRequest(ctx, nodeMap, processResults)
	if err != nil {
		return nil, fmt.Errorf("distribute request: %w", err)
	}

	zap.L().Info("first node return", zap.Any("address", nodeResponse.Address.String()))

	if nodeResponse.Err != nil {
		return nil, nodeResponse.Err
	}

	return nodeResponse.Data, nil
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
