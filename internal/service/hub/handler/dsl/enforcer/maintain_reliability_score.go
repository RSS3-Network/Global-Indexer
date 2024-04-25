package enforcer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
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

// getNodeStatsMap returns the current epoch.
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
	// Retrieve all Node addresses.
	nodeAddresses := extractNodeAddresses(stats)

	// Retrieve node information from the blockchain.
	nodesInfo, err := e.getNodesInfoFromBlockchain(nodeAddresses)
	if err != nil {
		return err
	}

	// Check if the length of NodesInfo and stats is the same.
	// TODO: If not, consider to process the queried nodes as much as possible
	if len(nodesInfo) != len(stats) {
		return fmt.Errorf("get Nodes info from blockchain: %d,%d", len(nodesInfo), len(stats))
	}

	// Retrieve node information from the database.
	nodes, err := e.getNodesInfoFromDatabase(ctx, nodeAddresses)
	if err != nil {
		return err
	}

	// Check if the length of Nodes and stats is the same.
	// TODO: If not, consider to process the queried nodes as much as possible
	if len(nodes) != len(stats) {
		return fmt.Errorf("get Nodes info from database: %d,%d", len(nodes), len(stats))
	}

	nodes = sortNodes(nodeAddresses, nodes)

	return updateStatsInPool(ctx, stats, nodesInfo, nodes, epoch)
}

func (e *SimpleEnforcer) getNodesInfoFromBlockchain(nodeAddresses []common.Address) ([]l2.DataTypesNode, error) {
	return e.stakingContract.GetNodes(&bind.CallOpts{}, nodeAddresses)
}

func (e *SimpleEnforcer) getNodesInfoFromDatabase(ctx context.Context, nodeAddresses []common.Address) ([]*schema.Node, error) {
	return e.databaseClient.FindNodes(ctx, schema.FindNodesQuery{NodeAddresses: nodeAddresses})
}

// sortNodes sorts Nodes by address.
func sortNodes(nodeAddresses []common.Address, nodes []*schema.Node) []*schema.Node {
	nodeMap := lo.SliceToMap(nodes, func(node *schema.Node) (common.Address, *schema.Node) {
		return node.Address, node
	})

	sortedNodes := make([]*schema.Node, len(nodeAddresses))

	for i, addr := range nodeAddresses {
		sortedNodes[i] = nodeMap[addr]
	}

	return sortedNodes
}

// updateStatsInPool concurrently updates the stats of the Nodes.
func updateStatsInPool(ctx context.Context, stats []*schema.Stat, nodesInfo []l2.DataTypesNode, nodes []*schema.Node, epoch int64) error {
	statsPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	for i, stat := range stats {
		i := i
		stat := stat

		statsPool.Go(func(_ context.Context) error {
			return updateNodeStat(stat, epoch, nodesInfo[i].StakingPoolTokens, nodes[i].Status)
		})
	}

	return statsPool.Wait()
}

// updateNodeStat updates Node's stat with Reliability Score.
func updateNodeStat(stat *schema.Stat, epoch int64, staking *big.Int, status schema.NodeStatus) error {
	// Convert the staking to float64.
	stat.Staking, _ = staking.Div(staking, big.NewInt(1e18)).Float64()

	if status == schema.NodeStatusOnline {
		// Reset the epoch request and invalid request if the epoch changes.
		if epoch != stat.Epoch {
			stat.EpochRequest = 0
			stat.EpochInvalidRequest = 0
			stat.Epoch = epoch
		}
	} else {
		// If Node's status is not online, then reset the alive time.
		stat.ResetAt = time.Now()
	}

	// calculate Reliability Score
	return calculateReliabilityScore(stat)
}

// calculateReliabilityScore calculates the Reliability Score σ of a given Node.
// σ is used to determine the probability of a Node receiving a request on DSL.
func calculateReliabilityScore(stat *schema.Stat) error {
	// staking pool tokens
	// maximum score is 0.2
	stat.Score = math.Min(math.Log(stat.Staking/stakingToScoreRate+1)/math.Log(stakingLogBase), stakingMaxScore)

	// public good node
	// If the Node is a public good node, then the score is 0
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

	// invalid request count in the current Epoch
	if stat.EpochInvalidRequest >= int64(model.DemotionCountBeforeSlashing) {
		// If the number of invalid requests in the epoch is greater than the threshold, then the score is 0.
		stat.Score = 0
	} else {
		stat.Score -= perSlashScore * float64(stat.EpochInvalidRequest)
	}

	return nil
}

// UpdateNodeCache updates the cache for all Nodes.
func (e *SimpleEnforcer) updateNodeCache(ctx context.Context, notify bool, epoch int64) error {
	if err := e.updateCacheForNodeType(ctx, model.RssNodeCacheKey); err != nil {
		return err
	}

	if err := e.updateCacheForNodeType(ctx, model.FullNodeCacheKey); err != nil {
		return err
	}

	if notify {
		return e.cacheClient.Set(ctx, model.NotifyKey, epoch)
	}

	return nil
}

// updateCacheForNodeType updates the cache for different types of Nodes.
func (e *SimpleEnforcer) updateCacheForNodeType(ctx context.Context, key string) error {
	nodesEndpointCache, err := retrieveNodeEndpointCache(ctx, key, e.databaseClient)
	if err != nil {
		return err
	}

	return e.setNodesEndpointCache(ctx, key, nodesEndpointCache)
}
