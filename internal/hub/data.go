package hub

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum/contract/staking"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/provider/node"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/naturalselectionlabs/rss3-node/config"
	"github.com/naturalselectionlabs/rss3-node/schema/filter"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

var (
	rssNodeCacheKey  = "nodes:rss"
	fullNodeCacheKey = "nodes:full"
)

func (h *Hub) getNode(ctx context.Context, address common.Address) (*schema.Node, error) {
	node, err := h.databaseClient.FindNode(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("get node %s: %w", address, err)
	}

	nodeInfo, err := h.stakingContract.GetNode(&bind.CallOpts{}, address)
	if err != nil {
		return nil, fmt.Errorf("get node from chain: %w", err)
	}

	node.Name = nodeInfo.Name
	node.Description = nodeInfo.Description
	node.TaxFraction = nodeInfo.TaxFraction
	node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
	node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
	node.TotalShares = nodeInfo.TotalShares.String()
	node.SlashedTokens = nodeInfo.SlashedTokens.String()

	return node, nil
}

func (h *Hub) getNodes(ctx context.Context, request *BatchNodeRequest) ([]*schema.Node, error) {
	nodes, err := h.databaseClient.FindNodes(ctx, request.NodeAddress, request.Cursor, request.Limit)
	if err != nil {
		return nil, fmt.Errorf("get nodes: %w", err)
	}

	addresses := lo.Map(nodes, func(node *schema.Node, _ int) common.Address {
		return node.Address
	})

	nodeInfo, err := h.stakingContract.GetNodes(&bind.CallOpts{}, addresses)
	if err != nil {
		return nil, fmt.Errorf("get nodes from chain: %w", err)
	}

	nodeInfoMap := lo.SliceToMap(nodeInfo, func(node staking.DataTypesNode) (common.Address, staking.DataTypesNode) {
		return node.Account, node
	})

	for _, node := range nodes {
		if nodeInfo, exists := nodeInfoMap[node.Address]; exists {
			node.Name = nodeInfo.Name
			node.Description = nodeInfo.Description
			node.TaxFraction = nodeInfo.TaxFraction
			node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
			node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
			node.TotalShares = nodeInfo.TotalShares.String()
			node.SlashedTokens = nodeInfo.SlashedTokens.String()
		}
	}

	return nodes, nil
}

func (h *Hub) registerNode(ctx context.Context, request *RegisterNodeRequest) error {
	node := &schema.Node{
		Address:  request.Address,
		Endpoint: request.Endpoint,
		Stream:   request.Stream,
		Config:   request.Config,
	}

	// Check node from chain.
	nodeInfo, err := h.stakingContract.GetNode(&bind.CallOpts{}, request.Address)
	if err != nil {
		return fmt.Errorf("get node from chain: %w", err)
	}

	if nodeInfo.Account == ethereum.AddressGenesis {
		return fmt.Errorf("node: %s has not been registered on the chain", request.Address.String())
	}

	node.IsPublicGood = nodeInfo.PublicGood

	stat := &schema.Stat{
		Address:      request.Address,
		Endpoint:     request.Endpoint,
		IsPublicGood: nodeInfo.PublicGood,
		ResetAt:      time.Now(),
		// todo: check if node is full node
		IsFullNode: true,
		IsRssNode:  len(request.Config.RSS) > 0,
		DecentralizedNetwork: len(lo.UniqBy(request.Config.Decentralized, func(module *config.Module) filter.Network {
			return module.Network
		})),
		FederatedNetwork: len(request.Config.Federated),
		Indexer:          len(request.Config.Decentralized),
	}

	indexers := make([]*schema.Indexer, 0, len(request.Config.Decentralized))

	for _, indexer := range request.Config.Decentralized {
		indexers = append(indexers, &schema.Indexer{
			Address: request.Address,
			Network: indexer.Network.String(),
			Worker:  indexer.Worker.String(),
		})
	}

	// Save node info to the database.
	if err := h.databaseClient.WithTransaction(ctx, func(ctx context.Context, client database.Client) error {
		// Save node to database.
		if err := h.databaseClient.SaveNode(ctx, node); err != nil {
			return fmt.Errorf("save node: %s, %w", node.Address.String(), err)
		}

		zap.L().Info("save node", zap.Any("node", node.Address.String()))

		// Save node stat to database
		if err := h.databaseClient.SaveNodeStat(ctx, stat); err != nil {
			return fmt.Errorf("save node stat: %s, %w", node.Address.String(), err)
		}

		zap.L().Info("save node stat", zap.Any("node", node.Address.String()))

		// Save node indexers to database
		if err := h.databaseClient.SaveNodeIndexers(ctx, indexers); err != nil {
			return fmt.Errorf("save node indexers: %s, %w", node.Address.String(), err)
		}

		zap.L().Info("save node indexer", zap.Any("node", node.Address.String()))

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (h *Hub) routerRSSHubData(ctx context.Context, path, query string) ([]byte, error) {
	nodes, err := h.filterNodes(ctx, rssNodeCacheKey)

	if err != nil {
		return nil, err
	}

	// TODO generate path
	nodeMap, err := h.pathBuilder.GetRSSHubPath(path, query, nodes)

	if err != nil {
		return nil, err
	}

	// TODO batch request
	node, err := h.batchRequest(ctx, nodeMap, processRSSHubResults)

	if err != nil {
		return nil, err
	}

	fmt.Printf("first node return %s\n", node.Address.String())

	return node.Data, nil
}

func (h *Hub) routerActivityData(ctx context.Context, request node.ActivityRequest) ([]byte, error) {
	nodes, err := h.filterNodes(ctx, fullNodeCacheKey)

	if err != nil {
		return nil, err
	}

	// TODO generate path
	nodeMap, err := h.pathBuilder.GetActivityByIDPath(request, nodes)

	if err != nil {
		return nil, err
	}

	// TODO batch request
	node, err := h.batchRequest(ctx, nodeMap, processActivityResults)

	if err != nil {
		return nil, err
	}

	fmt.Printf("first node return %s\n", node.Address.String())

	return node.Data, nil
}

func (h *Hub) routerActivitiesData(ctx context.Context, request node.AccountActivitiesRequest) ([]byte, error) {
	nodes, err := h.filterNodes(ctx, fullNodeCacheKey)

	if err != nil {
		return nil, err
	}

	// TODO generate path
	nodeMap, err := h.pathBuilder.GetAccountActivitiesPath(request, nodes)

	if err != nil {
		return nil, err
	}

	// TODO batch request
	node, err := h.batchRequest(ctx, nodeMap, processActivitiesResults)

	if err != nil {
		return nil, err
	}

	fmt.Printf("first node return %s\n", node.Address.String())

	return node.Data, nil
}

func (h *Hub) filterNodes(ctx context.Context, key string) ([]node.Cache, error) {
	var (
		nodesCache []node.Cache
		nodes      []*schema.Stat
	)

	// TODO: get nodes from cache
	exists, err := cache.Get(ctx, key, &nodesCache)
	if err != nil {
		return nil, fmt.Errorf("get nodes from cache: %s, %w", key, err)
	}

	if exists {
		return nodesCache, nil
	}

	// TODO: get nodes from database
	switch key {
	case rssNodeCacheKey:
		nodes, err = h.databaseClient.FindNodeStats(ctx)

		if err != nil {
			return nil, err
		}
	case fullNodeCacheKey:
		nodes, err = h.databaseClient.FindNodeStats(ctx)

		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown cache key: %s", key)
	}

	nodesCache = lo.Map(nodes, func(n *schema.Stat, _ int) node.Cache {
		return node.Cache{Address: n.Address, Endpoint: n.Endpoint}
	})

	if err := cache.Set(ctx, key, nodesCache); err != nil {
		return nil, fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nodesCache, nil
}

func processRSSHubResults(results []node.DataResponse) {
	fmt.Printf("rss finish, %d\n", len(results))

	for _, node := range results {
		fmt.Println(node.Address.String())
	}

	// TODO
	time.Sleep(time.Second * 1)
	fmt.Println("step1: data verify")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step2: data statistic")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step3: 2nd request、verify")
}

func processActivityResults(results []node.DataResponse) {
	fmt.Printf("feed id finish, %d\n", len(results))

	for _, node := range results {
		fmt.Println(node.Address.String())
	}

	// TODO
	time.Sleep(time.Second * 1)
	fmt.Println("step1: data verify")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step2: data statistic")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step3: 2nd request、verify")
}

func processActivitiesResults(results []node.DataResponse) {
	fmt.Printf("feeds finish, %d\n", len(results))

	for _, node := range results {
		fmt.Println(node.Address.String())
	}

	// TODO
	time.Sleep(time.Second * 1)
	fmt.Println("step1: data verify")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step2: data statistic")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step3: 2nd request、verify")
}

func (h *Hub) batchRequest(_ context.Context, nodeMap map[common.Address]string, processResults func([]node.DataResponse)) (node.DataResponse, error) {
	var (
		waitGroup   sync.WaitGroup
		firstResult = make(chan node.DataResponse, 1)
		results     []node.DataResponse
		mu          sync.Mutex
	)

	for address, endpoint := range nodeMap {
		waitGroup.Add(1)

		go func(address common.Address, endpoint string) {
			defer waitGroup.Done()

			data, err := h.fetch(context.Background(), endpoint)
			if err != nil {
				fmt.Printf("fetch error: %s, %v\n", address.String(), err)

				return
			}

			mu.Lock()
			results = append(results, node.DataResponse{Address: address, Data: data})
			mu.Unlock()

			select {
			case firstResult <- node.DataResponse{Address: address, Data: data}:
			default:
			}
		}(address, endpoint)
	}

	go func() {
		waitGroup.Wait()
		close(firstResult)

		processResults(results)
	}()

	select {
	case result := <-firstResult:
		return result, nil
	case <-time.After(time.Second * 3):
		return node.DataResponse{}, fmt.Errorf("timeout waiting for results")
	}
}

func (h *Hub) fetch(ctx context.Context, decodedURI string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, decodedURI, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	res, err := h.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return data, nil
}

func (h *Hub) verifyData(src, des []byte) bool {
	srcHash, destHash := sha256.Sum256(src), sha256.Sum256(des)

	return string(srcHash[:]) == string(destHash[:])
}

func (h *Hub) updateStats() error {
	return nil
}

// cron task
func (h *Hub) sortNodesTask() {
	// TODO nodes sort based on rules

	// TODO save to cache
}
