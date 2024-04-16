package enforcer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/distributor"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

const (
	stakingToScoreRate           float64 = 100000
	stakingLogBase                       = 2
	stakingMaxScore                      = 0.2
	hoursPerEpoch                        = 18
	activeTimeToScoreRate                = 120
	activeTimeMaxScore                   = 0.3
	totalReqToScoreRate                  = 100000
	totalReqLogBase                      = 100
	totalReqMaxScore                     = 0.3
	totalEpochReqToScoreRate             = 1000000
	totalEpochReqLogBase                 = 5000
	totalEpochReqMaxScore                = 1
	perDecentralizedNetworkScore         = 0.1
	perRssNetworkScore                   = 0.3
	perFederatedNetworkScore             = 0.1
	perIndexerScore                      = 0.05
	indexerMaxScore                      = 0.2
	perSlashScore                        = 0.5
	nonExistScore                float64 = 0
	existScore                           = 1

	defaultLimit = 50
)

func (e *SimpleEnforcer) getCurrentEpoch(ctx context.Context) (int64, error) {
	epochEvent, err := e.databaseClient.FindEpochs(ctx, 1, nil)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		zap.L().Error("get latest epoch event from database", zap.Error(err))
		return 0, err
	}

	if len(epochEvent) > 0 {
		return int64(epochEvent[0].ID), nil
	}

	return 0, nil
}

func (e *SimpleEnforcer) processNodeStats(ctx context.Context, stats []*schema.Stat, currentEpoch int64) error {
	if err := e.updateNodeStats(ctx, stats, currentEpoch); err != nil {
		return err
	}

	return e.databaseClient.SaveNodeStats(ctx, stats)
}

func (e *SimpleEnforcer) updateNodeStats(ctx context.Context, stats []*schema.Stat, epoch int64) error {
	// Retrieve all node addresses.
	nodeAddresses := extractNodeAddresses(stats)

	// Retrieve node information from the blockchain.
	nodesInfo, err := e.getNodesInfoFromBlockchain(nodeAddresses)
	if err != nil {
		return err
	}

	// Check if the length of nodesInfo and stats is the same.
	// TODO: If not, consider to process the queried nodes as much as possible
	if len(nodesInfo) != len(stats) {
		return fmt.Errorf("get nodes info from blockchain: %d,%d", len(nodesInfo), len(stats))
	}

	// Retrieve node information from the database.
	nodes, err := e.getNodesInfoFromDatabase(ctx, nodeAddresses)
	if err != nil {
		return err
	}

	// Check if the length of nodes and stats is the same.
	// TODO: If not, consider to process the queried nodes as much as possible
	if len(nodes) != len(stats) {
		return fmt.Errorf("get nodes info from database: %d,%d", len(nodes), len(stats))
	}

	return updateStatsInPool(ctx, stats, nodesInfo, nodes, epoch)
}

func (e *SimpleEnforcer) getNodesInfoFromBlockchain(nodeAddresses []common.Address) ([]l2.DataTypesNode, error) {
	return e.stakingContract.GetNodes(&bind.CallOpts{}, nodeAddresses)
}

func (e *SimpleEnforcer) getNodesInfoFromDatabase(ctx context.Context, nodeAddresses []common.Address) ([]*schema.Node, error) {
	return e.databaseClient.FindNodes(ctx, schema.FindNodesQuery{NodeAddresses: nodeAddresses})
}

// updateStatsInPool concurrently updates the stats of the nodes.
func updateStatsInPool(ctx context.Context, stats []*schema.Stat, nodesInfo []l2.DataTypesNode, nodes []*schema.Node, epoch int64) error {
	statsPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	for i, stat := range stats {
		staking := float64(nodesInfo[i].StakingPoolTokens.Uint64())
		status := nodes[i].Status

		statsPool.Go(func(_ context.Context) error {
			return updateNodeStat(stat, epoch, staking, status)
		})
	}

	return statsPool.Wait()
}

func updateNodeStat(stat *schema.Stat, epoch int64, staking float64, status schema.NodeStatus) error {
	stat.Staking = staking

	if status == schema.NodeStatusOnline {
		// Update node epoch.
		if epoch != stat.Epoch {
			stat.EpochRequest = 0
			stat.EpochInvalidRequest = 0
			stat.Epoch = epoch
		}
	} else {
		// If node's status is not online, then reset the time.
		stat.ResetAt = time.Now()
	}

	// calculate score
	return calculateScore(stat)
}

func calculateScore(stat *schema.Stat) error {
	// staking pool tokens
	// maximum score is 0.2
	stat.Score = math.Min(math.Log(stat.Staking/stakingToScoreRate+1)/math.Log(stakingLogBase), stakingMaxScore)

	// public good node
	// If the node is a public good node, then the score is 0
	// Otherwise, the score is 1
	stat.Score += lo.Ternary(stat.IsPublicGood, nonExistScore, existScore)

	// node active time
	// maximum score is 0.3
	// If node is active for about 2 epochs, then the score is 0.3
	stat.Score += math.Min(math.Ceil(time.Since(stat.ResetAt).Hours()/hoursPerEpoch)/activeTimeToScoreRate, activeTimeMaxScore)

	// total requests
	// maximum score is 0.3
	stat.Score += math.Min(math.Log(float64(stat.TotalRequest)/totalReqToScoreRate+1)/math.Log(totalReqLogBase), totalReqMaxScore)

	// epoch requests
	// maximum score is 1
	stat.Score += math.Min(math.Log(float64(stat.EpochRequest)/totalEpochReqToScoreRate+1)/math.Log(totalEpochReqLogBase), totalEpochReqMaxScore)

	// network count
	stat.Score += perDecentralizedNetworkScore*float64(stat.DecentralizedNetwork+stat.FederatedNetwork) + perRssNetworkScore*lo.Ternary(stat.IsRssNode, existScore, nonExistScore)

	// indexer count
	// maximum score is 0.2
	stat.Score += math.Min(float64(stat.Indexer)*perIndexerScore, indexerMaxScore)

	// epoch failure requests

	if stat.EpochInvalidRequest >= int64(distributor.DefaultSlashCount) {
		// If the number of invalid requests in the epoch is greater than the threshold, then the score is 0.
		stat.Score = 0
	} else {
		stat.Score -= perSlashScore * float64(stat.EpochInvalidRequest)
	}

	return nil
}

// UpdateNodeCache updates the cache for the node type.
func (e *SimpleEnforcer) updateNodeCache(ctx context.Context) error {
	if err := e.updateCacheForNodeType(ctx, distributor.RssNodeCacheKey); err != nil {
		return err
	}

	return e.updateCacheForNodeType(ctx, distributor.FullNodeCacheKey)
}

func (e *SimpleEnforcer) updateCacheForNodeType(ctx context.Context, key string) error {
	query := &schema.StatQuery{PointsOrder: lo.ToPtr("DESC")}

	switch key {
	case distributor.FullNodeCacheKey:
		query.IsFullNode = lo.ToPtr(true)
	case distributor.RssNodeCacheKey:
		query.IsRssNode = lo.ToPtr(true)
	}

	nodes, err := e.databaseClient.FindNodeStats(ctx, query)
	if err != nil {
		return err
	}

	qualifiedNodes, err := e.getQualifiedNodes(ctx, nodes)
	if err != nil {
		return err
	}

	return e.setNodeCache(ctx, key, qualifiedNodes)
}

// getQualifiedNodes filters the qualified nodes.
func (e *SimpleEnforcer) getQualifiedNodes(ctx context.Context, stats []*schema.Stat) ([]*schema.Stat, error) {
	nodeAddresses := extractNodeAddresses(stats)

	// Retrieve the online nodes from the database.
	nodes, err := e.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
		NodeAddresses: nodeAddresses,
		Status:        lo.ToPtr(schema.NodeStatusOnline),
	})

	if err != nil {
		return nil, err
	}

	nodeMap := lo.SliceToMap(nodes, func(node *schema.Node) (common.Address, struct{}) {
		return node.Address, struct{}{}
	})

	var qualifiedNodes []*schema.Stat

	// Exclude the offline nodes.
	for _, stat := range stats {
		if _, exists := nodeMap[stat.Address]; exists {
			qualifiedNodes = append(qualifiedNodes, stat)
		}

		if len(qualifiedNodes) >= distributor.DefaultNodeCount {
			break
		}
	}

	return qualifiedNodes, nil
}

// setNodeCache sets the cache for the nodes.
func (e *SimpleEnforcer) setNodeCache(ctx context.Context, key string, stats []*schema.Stat) error {
	nodesCache := lo.Map(stats, func(n *schema.Stat, _ int) distributor.NodeEndpointCache {
		return distributor.NodeEndpointCache{Address: n.Address.String(), Endpoint: n.Endpoint}
	})

	return e.cacheClient.Set(ctx, key, nodesCache)
}

// extractNodeAddresses returns all Node addresses from stats.
func extractNodeAddresses(stats []*schema.Stat) []common.Address {
	return lo.Map(stats, func(stat *schema.Stat, _ int) common.Address {
		return stat.Address
	})
}
