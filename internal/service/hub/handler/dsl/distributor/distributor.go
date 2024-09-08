package distributor

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
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

// DistributeDecentralizedData distributes decentralized requests to qualified Nodes.
func (d *Distributor) DistributeDecentralizedData(ctx context.Context, requestType string, request interface{}, params url.Values, workers, networks []string) ([]byte, error) {
	var (
		nodes          []*model.NodeEndpointCache
		processResults = d.processActivitiesResponses

		err error
	)

	switch requestType {
	case model.DistributorRequestActivity:
		nodes, err = d.simpleEnforcer.RetrieveQualifiedNodes(ctx, model.FullNodeCacheKey)
		processResults = d.processActivityResponses
	case model.DistributorRequestAccountActivities:
		nodes, err = d.getQualifiedNodes(ctx, workers, networks)
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
		return nil, err
	}

	nodeMap, err := d.generateDecentralizedPath(requestType, request, params, nodes)
	if err != nil {
		return nil, err
	}

	nodeResponse, err := d.simpleRouter.DistributeRequest(ctx, nodeMap, processResults)
	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeResponse.Address.String()))

	if nodeResponse.Err != nil {
		return nil, nodeResponse.Err
	}

	return nodeResponse.Data, nil
}

// generateDecentralizedPath builds the path for decentralized requests.
func (d *Distributor) generateDecentralizedPath(requestType string, request interface{}, params url.Values, nodes []*model.NodeEndpointCache) (map[common.Address]model.RequestMeta, error) {
	var (
		path   string
		method = http.MethodGet

		body []byte
		err  error
	)

	switch req := request.(type) {
	case dsl.ActivityRequest:
		path = fmt.Sprintf("/decentralized/tx/%s", req.ID)
	case dsl.ActivitiesRequest:
		path = fmt.Sprintf("/decentralized/%s", req.Account)
	case dsl.AccountsActivitiesRequest:
		path = "/decentralized/accounts"
		method = http.MethodPost
		body, err = json.Marshal(req)

		if err != nil {
			return nil, fmt.Errorf("marshal request data: %w", err)
		}
	case dsl.NetworkActivitiesRequest:
		path = fmt.Sprintf("/decentralized/network/%s", req.Network)
	case dsl.PlatformActivitiesRequest:
		path = fmt.Sprintf("/decentralized/platform/%s", req.Platform)
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
func NewDistributor(ctx context.Context, database database.Client, cache cache.Client, httpClient httputil.Client, stakingContract *stakingv2.Staking, networkParamsContract *l2.NetworkParams, txManager *txmgr.SimpleTxManager, settlerConfig *config.Settler, chainID *big.Int) (*Distributor, error) {
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
