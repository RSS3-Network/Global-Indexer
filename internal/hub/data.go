package hub

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum/contract/staking"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/provider/node"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/naturalselectionlabs/rss3-node/config"
	"github.com/naturalselectionlabs/rss3-node/schema/filter"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	rssNodeCacheKey  = "nodes:rss"
	fullNodeCacheKey = "nodes:full"
)

var message = "I, %s, am signing this message for registering my intention to operate an RSS3 Serving Node."

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
	// Check signature.
	if err := h.checkSignature(ctx, request.Address, request.Signature); err != nil {
		return err
	}

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
		return fmt.Errorf("node: %s has not been registered on the chain", strings.ToLower(request.Address.String()))
	}

	if strings.Compare(nodeInfo.OperationPoolTokens.String(), decimal.NewFromInt(10000).Mul(decimal.NewFromInt(1e18)).String()) < 0 {
		return fmt.Errorf("insufficient operation pool tokens")
	}

	node.IsPublicGood = nodeInfo.PublicGood

	stat := &schema.Stat{
		Address:      request.Address,
		Endpoint:     request.Endpoint,
		IsPublicGood: nodeInfo.PublicGood,
		ResetAt:      time.Now(),
		IsFullNode:   h.isFullNode(request.Config.Decentralized),
		IsRssNode:    len(request.Config.RSS) > 0,
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

// Check if node is full node
func (h *Hub) isFullNode(indexers []*config.Module) bool {
	// TODO: A more comprehensive check should also include the network and parameters.
	if len(indexers) < len(filter.NameStrings())-1 {
		return false
	}

	workerMap := make(map[string]struct{}, 0)

	for _, indexer := range indexers {
		workerMap[indexer.Worker.String()] = struct{}{}
	}

	for _, worker := range filter.NameStrings() {
		if worker == filter.Unknown.String() {
			continue
		}

		if _, exists := workerMap[worker]; !exists {
			return false
		}
	}

	return true
}

func (h *Hub) routerRSSHubData(ctx context.Context, path, query string) ([]byte, error) {
	nodes, err := h.filterNodes(ctx, rssNodeCacheKey)

	if err != nil {
		return nil, err
	}

	// generate path
	nodeMap, err := h.pathBuilder.GetRSSHubPath(path, query, nodes)

	if err != nil {
		return nil, err
	}

	// batch request
	node, err := h.batchRequest(ctx, nodeMap, h.processRSSHubResults)

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
	node, err := h.batchRequest(ctx, nodeMap, h.processActivityResults)

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
	node, err := h.batchRequest(ctx, nodeMap, h.processActivitiesResults)

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

	// Get nodes from cache.
	exists, err := cache.Get(ctx, key, &nodesCache)
	if err != nil {
		return nil, fmt.Errorf("get nodes from cache: %s, %w", key, err)
	}

	if exists {
		return nodesCache, nil
	}

	// Get nodes from database.
	switch key {
	case rssNodeCacheKey:
		nodes, err = h.databaseClient.FindNodeStatsByType(ctx, nil, lo.ToPtr(true), 3)

		if err != nil {
			return nil, err
		}
	case fullNodeCacheKey:
		nodes, err = h.databaseClient.FindNodeStatsByType(ctx, lo.ToPtr(true), nil, 3)

		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown cache key: %s", key)
	}

	if err = h.setNodeCache(ctx, key, nodes); err != nil {
		return nil, err
	}

	return nodesCache, nil
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
				zap.L().Error("fetch request error", zap.Any("node", address.String()))
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

func (h *Hub) processRSSHubResults(results []node.DataResponse) {
	zap.L().Info("begin rss request verify", zap.Any("results", len(results)))

	ctx := context.Background()

	h.processData(ctx, results)
}

func (h *Hub) processActivityResults(results []node.DataResponse) {
	zap.L().Info("begin feed id request verify", zap.Any("results", len(results)))

	ctx := context.Background()

	h.processData(ctx, results)
}

func (h *Hub) processActivitiesResults(results []node.DataResponse) {
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

func (h *Hub) processData(ctx context.Context, results []node.DataResponse) error {
	stats, err := h.databaseClient.FindNodeStats(ctx, lo.Map(results, func(result node.DataResponse, _ int) common.Address {
		return result.Address
	}))

	if err != nil {
		// TODO error handling
		return err
	}

	verify01 := h.verifyData(results[0].Data, results[1].Data)
	verify02 := h.verifyData(results[0].Data, results[2].Data)
	verify12 := h.verifyData(results[1].Data, results[2].Data)

	if verify01 && verify02 {
		results[0].Request = 2
		results[1].Request = 1
		results[2].Request = 1
	}

	if !verify01 && verify12 {
		results[0].InvalidRequest = 1
		results[1].Request = 1
		results[2].Request = 1
	}

	if verify01 && !verify12 {
		results[0].Request = 2
		results[1].Request = 1
		results[2].InvalidRequest = 1
	}

	if verify02 && !verify12 {
		results[0].Request = 2
		results[1].InvalidRequest = 1
		results[2].Request = 1
	}

	for _, result := range results {
		for _, stat := range stats {
			if result.Address == stat.Address {
				stat.TotalRequest += int64(result.Request)
				stat.EpochRequest += int64(result.Request)
				stat.EpochInvalidRequest += int64(result.InvalidRequest)
			}
		}
	}

	if err = h.databaseClient.SaveNodeStats(ctx, stats); err != nil {
		// TODO error handling
		return err
	}

	return nil
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

// cron task
func (h *Hub) sortNodesTask() error {
	var (
		stats []*schema.Stat

		err error
	)

	ctx := context.Background()

	stats, err = h.databaseClient.FindNodeStats(ctx, []common.Address{})

	if err != nil {
		return err
	}

	// TODO: parallel
	for _, stat := range stats {
		if err = h.updateNodeEpochStats(stat); err != nil {
			return err
		}

		h.calcPoints(stat)

		if err = h.databaseClient.SaveNodeStat(ctx, stat); err != nil {
			return err
		}
	}

	// Update node cache.
	rssNodes, err := h.databaseClient.FindNodeStatsByType(ctx, nil, lo.ToPtr(true), 3)

	if err != nil {
		return err
	}

	if err = h.setNodeCache(ctx, rssNodeCacheKey, rssNodes); err != nil {
		return err
	}

	fullNodes, err := h.databaseClient.FindNodeStatsByType(ctx, lo.ToPtr(true), nil, 3)

	if err != nil {
		return err
	}

	if err = h.setNodeCache(ctx, fullNodeCacheKey, fullNodes); err != nil {
		return err
	}

	return nil
}

func (h *Hub) setNodeCache(ctx context.Context, key string, stats []*schema.Stat) error {
	nodesCache := lo.Map(stats, func(n *schema.Stat, _ int) node.Cache {
		return node.Cache{Address: n.Address, Endpoint: n.Endpoint}
	})

	if err := cache.Set(ctx, key, nodesCache); err != nil {
		return fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nil
}

func (h *Hub) updateNodeEpochStats(stat *schema.Stat) error {
	nodeInfo, err := h.stakingContract.GetNode(&bind.CallOpts{}, stat.Address)

	if err != nil {
		return fmt.Errorf("get node info: %s,%w", stat.Address.String(), err)
	}

	stat.Staking = float64(nodeInfo.StakingPoolTokens.Uint64())
	stat.EpochRequest = 0
	stat.EpochInvalidRequest = 0

	return nil
}

// calculation rule https://docs.google.com/spreadsheets/d/1N7zEwUooiOjCIHzhoHuf8aM_lbF5bS0ZC-4luxc2qNU/edit?pli=1#gid=0
func (h *Hub) calcPoints(stat *schema.Stat) {
	// staking pool tokens
	stat.Points = math.Min(math.Log2(stat.Staking/100000)+1, 0.2)

	// public good
	stat.Points += float64(lo.Ternary(stat.IsPublicGood, 0, 1))

	// running time
	stat.Points += math.Min(math.Ceil(time.Since(stat.ResetAt).Hours()/18)/120, 0.3)

	// total requests
	stat.Points += math.Min(math.Log(float64(stat.TotalRequest)/100000+1)/math.Log(100), 0.3)

	// epoch requests
	stat.Points += math.Min(math.Log(float64(stat.EpochRequest)/1000000+1)/math.Log(5000), 1)

	// network count
	stat.Points += 0.1*float64(stat.DecentralizedNetwork+stat.FederatedNetwork) + 0.3*float64(lo.Ternary(stat.IsRssNode, 1, 0))

	// indexer count
	stat.Points += math.Min(float64(stat.Indexer)*0.05, 0.2)

	// epoch failure requests
	stat.Points -= 0.5 * float64(stat.EpochInvalidRequest)
}

func (h *Hub) checkSignature(_ context.Context, address common.Address, signature string) error {
	message := fmt.Sprintf(message, strings.ToLower(address.Hex()))

	pubKey, err := crypto.SigToPub(crypto.Keccak256Hash([]byte(message)).Bytes(), []byte(signature))
	if err != nil {
		return fmt.Errorf("failed to parse signature: %w", err)
	}

	result := crypto.PubkeyToAddress(*pubKey)

	if address != result {
		return fmt.Errorf("invalid signature")
	}

	return nil
}
