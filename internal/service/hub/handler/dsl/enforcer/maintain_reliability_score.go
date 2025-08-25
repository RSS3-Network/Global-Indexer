package enforcer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
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
	perAINetworkScore                    = 0.3
	perRsshubScore                       = 0.3
	perFederatedNetworkScore             = 0.1
	perIndexerScore                      = 0.05
	indexerMaxScore                      = 0.2
	perSlashScore                        = 0.5
	nonExistScore                float64 = 0
	existScore                           = 1

	defaultLimit = 50
)

// getCurrentEpoch returns the current epoch.
func (e *SimpleEnforcer) getCurrentEpoch(ctx context.Context) (int64, error) {
	epochEvent, err := e.databaseClient.FindEpochs(ctx, &schema.FindEpochsQuery{Limit: lo.ToPtr(1)})
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		zap.L().Error("get latest epoch event from database", zap.Error(err))
		return 0, err
	}

	if len(epochEvent) > 0 {
		return int64(epochEvent[0].ID), nil
	}

	return 0, nil
}

// getAllNodeStats retrieves all node statistics matching the given query from the database.
func (e *SimpleEnforcer) getAllNodeStats(ctx context.Context, query *schema.StatQuery) ([]*schema.Stat, error) {
	stats := make([]*schema.Stat, 0)

	// Traverse the entire node.
	for {
		tempStats, err := e.databaseClient.FindNodeStats(ctx, query)
		if err != nil {
			return nil, err
		}

		// If there are no stats, exit the loop.
		if len(tempStats) == 0 {
			break
		}

		stats = append(stats, tempStats...)
		query.Cursor = lo.ToPtr(tempStats[len(tempStats)-1].Address.String())
	}

	return stats, nil
}

// processNodeStats processes the node statistics in parallel.
func (e *SimpleEnforcer) processNodeStats(ctx context.Context, stats []*schema.Stat, reset bool) error {
	if err := e.updateNodeStats(ctx, stats, reset); err != nil {
		return err
	}

	return e.databaseClient.SaveNodeStats(ctx, stats)
}

func (e *SimpleEnforcer) updateNodeStats(ctx context.Context, stats []*schema.Stat, reset bool) error {
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

	return e.updateStatsInPool(ctx, stats, nodesInfo, nodes, reset)
}

func (e *SimpleEnforcer) getNodesInfoFromBlockchain(nodeAddresses []common.Address) ([]stakingv2.Node, error) {
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
func (e *SimpleEnforcer) updateStatsInPool(ctx context.Context, stats []*schema.Stat, nodesInfo []stakingv2.Node, nodes []*schema.Node, reset bool) error {
	statsPool := pool.New().WithContext(ctx).WithMaxGoroutines(lo.Ternary(len(stats) < 20*runtime.NumCPU() && len(stats) > 0, len(stats), 20*runtime.NumCPU()))

	for i := range stats {
		i := i

		statsPool.Go(func(ctx context.Context) error {
			var validCount int64

			// Get the latest valid request count from the cache.
			if err := e.cacheClient.Get(ctx, formatNodeStatRedisKey(model.ValidRequestCount, stats[i].Address.String()), &validCount); err != nil && !errors.Is(err, redis.Nil) {
				return fmt.Errorf("get valid request count: %w", err)
			}

			// Update the total request count for the node.
			// This action occurs only if the valid request count from the cache surpasses the epoch request count from the db,
			// which indicates the receipt of new requests during the current epoch.
			if validCount >= stats[i].EpochRequest {
				stats[i].TotalRequest += validCount - stats[i].EpochRequest
				stats[i].EpochRequest = validCount
			}

			// Check the reset flag status for the node. If the reset flag is not set (false),
			// retrieve the invalid request count from the cache to maintain current epoch data.
			// If the reset flag is true, initialize the valid and invalid request counts to zero,
			// effectively resetting the node's counters for the new epoch.
			if !reset {
				if err := e.cacheClient.Get(ctx, formatNodeStatRedisKey(model.InvalidRequestCount, stats[i].Address.String()), &stats[i].EpochInvalidRequest); err != nil && !errors.Is(err, redis.Nil) {
					return fmt.Errorf("get invalid request count: %w", err)
				}
			} else {
				stats[i].EpochRequest = 0

				if err := e.cacheClient.Set(ctx, formatNodeStatRedisKey(model.ValidRequestCount, stats[i].Address.String()), 0, 0); err != nil {
					return fmt.Errorf("reset valid request count: %w", err)
				}

				if err := e.cacheClient.Set(ctx, formatNodeStatRedisKey(model.InvalidRequestCount, stats[i].Address.String()), stats[i].EpochInvalidRequest, 0); err != nil {
					return fmt.Errorf("reset invalid request count: %w", err)
				}
			}

			updateNodeStat(stats[i], nodesInfo[i].StakingPoolTokens, nodes[i].Status)

			return nil
		})
	}

	return statsPool.Wait()
}

// updateNodeStat updates Node's stat with Reliability Score.
func updateNodeStat(stat *schema.Stat, staking *big.Int, status schema.NodeStatus) {
	// Convert the staking to float64.
	stat.Staking, _ = staking.Div(staking, big.NewInt(1e18)).Float64()

	if status != schema.NodeStatusOnline {
		// If Node's status is not online, then reset the alive time.
		stat.ResetAt = time.Now()
	}

	// Calculate the Reliability Score.
	calculateReliabilityScore(stat)
}

// calculateReliabilityScore calculates the Reliability Score σ of a given Node.
// σ is used to determine the probability of a Node receiving a request on DSL.
func calculateReliabilityScore(stat *schema.Stat) {
	// baseline score
	baselineScore := math.Min(math.Log(stat.Staking/stakingToScoreRate+1)/math.Log(stakingLogBase), stakingMaxScore)

	// staking pool tokens
	// maximum score is 0.2
	stat.Score = baselineScore

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
	stat.Score += perDecentralizedNetworkScore*float64(stat.DecentralizedNetwork) +
		perRssNetworkScore*lo.Ternary(stat.IsRssNode, existScore, nonExistScore) +
		perAINetworkScore*lo.Ternary(stat.IsAINode, existScore, nonExistScore) +
		perRsshubScore*lo.Ternary(stat.IsRsshubNode, existScore, nonExistScore) +
		perFederatedNetworkScore*float64(stat.FederatedNetwork)

	// indexer count
	// maximum score is 0.2
	stat.Score += math.Min(float64(stat.Indexer)*perIndexerScore, indexerMaxScore)

	// invalid request count in the current Epoch
	if stat.EpochInvalidRequest >= int64(model.DemotionCountBeforeSlashing) {
		// If the number of invalid requests in the epoch is greater than the threshold, then the score is baseline score.
		stat.Score = baselineScore
	} else {
		// If the number of invalid requests in the epoch is less than the threshold, then the score is baseline score minus the number of invalid requests.
		stat.Score = math.Max(baselineScore, stat.Score-perSlashScore*float64(stat.EpochInvalidRequest))
	}
}
