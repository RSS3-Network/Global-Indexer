package settler

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
)

// calculateActiveScores calculates active scores for all Nodes based on various factors.
func calculateActiveScores(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, activeScores *config.ActiveScores) ([]*big.Float, error) {
	var (
		finalScores = make([]*big.Float, len(nodes))
		// sum of all active scores
		totalActiveScore = big.NewFloat(0)
		// sum of all Ps with recent stakers.
		totalActiveStake       = big.NewInt(0)
		totalOperationPoolSize = big.NewInt(0)
		// sum of all Ps sizes.
		totalStake *big.Int
	)

	// Preprocessing step to avoid repeated parsing and condition checking.
	stakingPoolSizes, operationPoolSizes, err := parsePoolSizes(nodes)
	if err != nil {
		return nil, err
	}

	// Exclude unqualified nodes.
	excludeUnqualifiedNodes(nodes, recentStakers)

	// Calculate the total pool size.
	for i, poolSize := range stakingPoolSizes {
		if _, exist := recentStakers[nodes[i].Address]; exist {
			totalActiveStake.Add(totalActiveStake, poolSize)
		}
	}

	// Calculate the total operation pool size.
	for _, poolSize := range operationPoolSizes {
		totalOperationPoolSize.Add(totalOperationPoolSize, poolSize)
	}

	// Calculate total stake.
	totalStake = sumTotalStake(nodes, recentStakers)

	// Calculate activeScores for each node.
	scores := computeActiveScores(nodes, recentStakers, stakingPoolSizes, totalActiveStake, totalStake, activeScores)

	// Calculate operationScores for each node.
	operationScores := computeOperationScores(operationPoolSizes, totalOperationPoolSize)

	// calculate the total finalScores for each node
	for i, activeScore := range scores {
		totalActiveScore.Add(totalActiveScore, activeScore)

		// finalScore = activeScore + operationScore
		finalScores[i] = new(big.Float).Add(activeScore, operationScores[i])
	}

	return finalScores, nil
}

// filter retrieves Node information from a staking contract.
func (s *Server) filter(nodeAddresses []common.Address, nodes []*schema.Node) ([]*schema.Node, []common.Address, error) {
	nodeInfoList, err := s.stakingContract.GetNodes(&bind.CallOpts{}, nodeAddresses)
	if err != nil {
		return nil, nil, fmt.Errorf("get Nodes from chain: %w", err)
	}

	nodeInfoMap := lo.SliceToMap(nodeInfoList, func(node stakingv2.Node) (common.Address, stakingv2.Node) {
		return node.Account, node
	})

	newNodes := make([]*schema.Node, 0, len(nodes))
	newNodeAddresses := make([]common.Address, 0, len(nodes))

	for i := range nodes {
		if nodeInfo, ok := nodeInfoMap[nodes[i].Address]; ok && isValidStatus(nodeInfo.Status) {
			nodes[i].StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
			nodes[i].OperationPoolTokens = nodeInfo.OperationPoolTokens.String()

			newNodes = append(newNodes, nodes[i])
			newNodeAddresses = append(newNodeAddresses, nodes[i].Address)
		}
	}

	return newNodes, newNodeAddresses, nil
}

// isValidStatus checks if the node status is valid.
func isValidStatus(status uint8) bool {
	return status == uint8(schema.NodeStatusInitializing) ||
		status == uint8(schema.NodeStatusOnline) ||
		status == uint8(schema.NodeStatusExiting)
}

// excludeUnqualifiedNodes excludes Nodes if they:
// 1. have no recent stakers
// 2. are offline
func excludeUnqualifiedNodes(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount) {
	onlineNodes := lo.SliceToMap(nodes, func(node *schema.Node) (common.Address, struct{}) {
		return node.Address, struct{}{}
	})

	for address := range recentStakers {
		if _, ok := onlineNodes[address]; !ok {
			delete(recentStakers, address)
		}
	}
}

// parsePoolSizes extracts and parses staking and operation pool sizes for all Nodes.
func parsePoolSizes(nodes []*schema.Node) ([]*big.Int, []*big.Int, error) {
	stakingPoolSizes := make([]*big.Int, len(nodes))
	operationPoolSizes := make([]*big.Int, len(nodes))

	for i, node := range nodes {
		stakingPoolSize, ok := new(big.Int).SetString(node.StakingPoolTokens, 10)

		if !ok {
			return nil, nil, fmt.Errorf("failed to parse staking pool tokens for node %s: invalid number", node.Address)
		}

		stakingPoolSizes[i] = stakingPoolSize

		operationPoolSize, ok := new(big.Int).SetString(node.OperationPoolTokens, 10)

		if !ok {
			return nil, nil, fmt.Errorf("failed to parse operation pool tokens for node %s: invalid number", node.Address)
		}

		operationPoolSizes[i] = operationPoolSize
	}

	return stakingPoolSizes, operationPoolSizes, nil
}

// sumTotalStake sum the total stake.
func sumTotalStake(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount) *big.Int {
	var totalStake = big.NewInt(0)

	for _, node := range nodes {
		if _, exist := recentStakers[node.Address]; exist {
			totalStake.Add(totalStake, recentStakers[node.Address].StakeValue.BigInt())
		}
	}

	return totalStake
}

// computeActiveScores calculates active scores for all Nodes based on various factors.
// Active score is used to calculate Alpha Operation Rewards, and will be deprecated in the future.
func computeActiveScores(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, poolSizes []*big.Int, totalPoolSize *big.Int, totalStakeValue *big.Int, activeScores *config.ActiveScores) []*big.Float {
	scores := make([]*big.Float, len(nodes))

	for i, poolSize := range poolSizes {
		// If the Node has no recent stakers, set the score to 0
		if _, exist := recentStakers[nodes[i].Address]; !exist {
			scores[i] = big.NewFloat(0)
			continue
		}

		stakeInfo := recentStakers[nodes[i].Address]

		// Calculate the ratio of poolSize to totalPoolSize
		// poolSizeRatio = poolSize / totalPoolSize
		poolSizeRatio := new(big.Float).Quo(new(big.Float).SetInt(poolSize), new(big.Float).SetInt(totalPoolSize))
		// Calculate the ratio of stakeValue to totalStakeValue
		// stakeRadio = stakeValue / totalStakeValue
		stakeRadio := new(big.Float).Quo(new(big.Float).SetInt(stakeInfo.StakeValue.BigInt()), new(big.Float).SetInt(totalStakeValue))

		score := applyGiniCoefficient(poolSizeRatio, activeScores.GiniCoefficient)

		applyStakerFactor(stakeInfo.StakerCount, stakeRadio, activeScores.StakerFactor, score)

		scores[i] = score
	}

	return scores
}

// computeOperationScores calculates the scores for each Node based on the operation pool size.
func computeOperationScores(poolSizes []*big.Int, totalPoolSize *big.Int) []*big.Float {
	scores := make([]*big.Float, len(poolSizes))

	for i, poolSize := range poolSizes {
		// Calculate the ratio(scores) of poolSize to totalPoolSize
		scores[i] = new(big.Float).Quo(new(big.Float).SetInt(poolSize), new(big.Float).SetInt(totalPoolSize))
	}

	return scores
}

// applyGiniCoefficient applies the Gini Coefficient to the score
func applyGiniCoefficient(poolSizeRatio *big.Float, giniCoefficient float64) *big.Float {
	// Perform calculation: score = 1 / (1 + giniCoefficient * poolSizeRatio)
	one := big.NewFloat(1)
	giniTimesPool := new(big.Float).Mul(new(big.Float).SetFloat64(giniCoefficient), poolSizeRatio)
	denominator := new(big.Float).Add(one, giniTimesPool)
	score := new(big.Float).Quo(one, denominator)

	return score
}

// applyStakerFactor applies the Staker Factor to the score
func applyStakerFactor(stakers uint64, stakeRadio *big.Float, stakerFactor float64, score *big.Float) {
	// Convert stakers to big.Float for calculation.
	stakersFloat := new(big.Float).SetUint64(stakers)
	// Ensure stakerFactor is in big.Float for consistency in operations.
	stakerFactorFloat := big.NewFloat(stakerFactor)
	// stakerFactor * stakeRatio
	stakerFactorCalculated := new(big.Float).Mul(stakerFactorFloat, stakersFloat)
	// Calculate the exponent for the exponential function.
	exponentFloat := new(big.Float).Mul(stakerFactorCalculated, stakeRadio)

	exponent, _ := exponentFloat.Float64()
	// Perform calculation: score *= math.exp(stakerFactorCalculated * stakeRatio)
	expResultBigFloat := new(big.Float).SetFloat64(math.Exp(exponent))

	score.Mul(score, expResultBigFloat)
}
