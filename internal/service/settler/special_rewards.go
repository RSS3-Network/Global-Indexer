package settler

import (
	"fmt"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
)

// calculateAlphaSpecialRewards calculates the distribution of the Special Rewards used to replace the Operation Rewards
// the Special Rewards are used to incentivize staking in smaller Nodes
// currently, the amount is set to 30,000,000 / 486.6666666666667 * 0.2 ~= 12328
func calculateAlphaSpecialRewards(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, specialRewards *config.SpecialRewards) ([]*big.Int, []*big.Float, error) {
	var (
		finalScores = make([]*big.Float, len(nodes))
		// sum of all active scores
		totalActiveScore = big.NewFloat(0)
		// sum of all Ps with recent stakers.
		totalActiveStake       = big.NewInt(0)
		totalActiveRewards     float64
		totalOperationPoolSize = big.NewInt(0)
		totalOperationRewards  float64
		// sum of all Ps sizes.
		totalStake *big.Int
	)

	// Preprocessing step to avoid repeated parsing and condition checking.
	stakingPoolSizes, operationPoolSizes, err := parsePoolSizes(nodes)
	if err != nil {
		return nil, nil, err
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

	// Calculate the ratio of active nodes to total nodes.
	activeNodesRadio := float64(len(recentStakers)) / float64(len(nodes))

	totalActiveRewards = specialRewards.Rewards * specialRewards.RewardsRatioActive
	// If the ratio is less than the threshold, reduce the rewards
	if activeNodesRadio < specialRewards.NodeThreshold {
		totalActiveRewards *= activeNodesRadio
	}

	// Calculate the total operation rewards.
	totalOperationRewards = specialRewards.Rewards * specialRewards.RewardsRatioOperation

	// Calculate activeScores for each node.
	activeScores := computeActiveScores(nodes, recentStakers, stakingPoolSizes, totalActiveStake, totalStake, specialRewards)

	// Calculate operationScores for each node.
	operationScores := computeOperationScores(operationPoolSizes, totalOperationPoolSize)

	// calculate the total activeScore for each node
	for i, activeScore := range activeScores {
		totalActiveScore.Add(totalActiveScore, activeScore)

		// finalScore = activeScore + operationScore
		finalScores[i] = new(big.Float).Add(activeScore, operationScores[i])
	}

	rewards, err := calculateFinalRewards(activeScores, totalActiveScore, operationScores, specialRewards, totalActiveRewards, totalOperationRewards)
	if err != nil {
		return nil, nil, err
	}

	return rewards, finalScores, nil
}

// fetchNodePoolSizes retrieves Node information from a staking contract
// and updates the staking and operation pool sizes for each Node.
func (s *Server) fetchNodePoolSizes(nodeAddresses []common.Address, nodes []*schema.Node) error {
	nodeInfo, err := s.stakingContract.GetNodes(&bind.CallOpts{}, nodeAddresses)
	if err != nil {
		return fmt.Errorf("get Nodes from chain: %w", err)
	}

	nodeInfoMap := lo.SliceToMap(nodeInfo, func(node stakingv2.DataTypesNode) (common.Address, stakingv2.DataTypesNode) {
		return node.Account, node
	})

	for _, node := range nodes {
		if nodeInfo, ok := nodeInfoMap[node.Address]; ok {
			node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
			node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
		}
	}

	return nil
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
func computeActiveScores(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, poolSizes []*big.Int, totalPoolSize *big.Int, totalStakeValue *big.Int, specialRewards *config.SpecialRewards) []*big.Float {
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

		score := applyGiniCoefficient(poolSizeRatio, specialRewards.GiniCoefficient)

		applyStakerFactor(stakeInfo.StakerCount, stakeRadio, specialRewards.StakerFactor, score)

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

// calculateFinalRewards converts scores into reward amounts.
func calculateFinalRewards(activeScores []*big.Float, totalActiveScore *big.Float, operationScores []*big.Float, specialRewards *config.SpecialRewards, totalRewards, totalOperationRewards float64) ([]*big.Int, error) {
	rewards := make([]*big.Int, len(activeScores))
	maxReward := big.NewFloat(0).SetFloat64(specialRewards.RewardsCeiling)

	for i, activeScore := range activeScores {
		// Apply active rewards
		activeReward := big.NewFloat(0)

		// If totalActiveScore is greater than 0, calculate the active reward
		if totalActiveScore.Cmp(big.NewFloat(0)) > 0 {
			// Calculate the ratio of activeScore to totalScore
			scoreRatio := new(big.Float).Quo(activeScore, totalActiveScore)
			activeReward = new(big.Float).Mul(scoreRatio, big.NewFloat(0).SetFloat64(totalRewards))
		}
		// Apply operation rewards
		operationReward := new(big.Float).Mul(operationScores[i], big.NewFloat(0).SetFloat64(totalOperationRewards))
		reward := new(big.Float).Add(activeReward, operationReward)

		if reward.Cmp(maxReward) == 1 {
			reward = maxReward
		}

		// Convert to integer to truncate before scaling
		rewardFinal, _ := reward.Int(nil)

		// Apply gwei after truncation
		scaleGwei(rewardFinal)

		rewards[i] = rewardFinal
	}

	err := checkRewardsCeiling(rewards, specialRewards.Rewards)
	if err != nil {
		return nil, err
	}

	return rewards, nil
}

// checkRewardsCeiling checks if the sum of rewards is less than or equal to specialRewards.Rewards.
func checkRewardsCeiling(rewards []*big.Int, specialRewards float64) error {
	sum := big.NewInt(0)
	for _, reward := range rewards {
		sum.Add(sum, reward)
	}

	// Scale the specialRewards by 10^18 to match the rewards scale
	specialRewardsBigInt := big.NewInt(0).SetUint64(uint64(specialRewards))
	scaleGwei(specialRewardsBigInt)

	if sum.Cmp(specialRewardsBigInt) > 0 {
		return fmt.Errorf("total rewards exceed the ceiling: %v > %v", sum, specialRewardsBigInt)
	}

	return nil
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
