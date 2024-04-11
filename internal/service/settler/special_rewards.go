package settler

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

// calculateAlphaSpecialRewards calculates the distribution of the Special Rewards used to replace the Operation Rewards
// the Special Rewards are used to incentivize staking in smaller Nodes
// currently, the amount is set to 30,000,000 / 486.6666666666667 * 0.2 ~= 12328
func calculateAlphaSpecialRewards(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, specialRewards *config.SpecialRewards) ([]*big.Int, []*big.Float, error) {
	var (
		totalStakeValue *big.Int
		totalPoolSize   = big.NewInt(0)
		totalScore      = big.NewFloat(0)
		totalRewards    float64
	)

	// Preprocessing step to avoid repeated parsing and condition checking.
	poolSizes, err := parsePoolSizes(nodes)
	if err != nil {
		return nil, nil, err
	}

	// Exclude unqualified nodes.
	excludeUnqualifiedNodes(nodes, recentStakers)

	// Calculate the total pool size.
	for i, poolSize := range poolSizes {
		if _, exist := recentStakers[nodes[i].Address]; exist {
			totalPoolSize.Add(totalPoolSize, poolSize)
		}
	}

	// Calculate total stake value.
	totalStakeValue = computeTotalStakeValue(nodes, recentStakers)

	// Calculate the ratio of active nodes to total nodes
	activeNodesRadio := float64(len(recentStakers)) / float64(len(nodes))

	totalRewards = specialRewards.Rewards
	// If the ratio is less than the threshold, reduce the rewards
	if activeNodesRadio < specialRewards.NodeThreshold {
		totalRewards = specialRewards.Rewards * activeNodesRadio
	}

	// Calculate scores for each node.
	scores := computeScores(nodes, recentStakers, poolSizes, totalPoolSize, totalStakeValue, specialRewards)

	for _, score := range scores {
		totalScore.Add(totalScore, score)
	}

	rewards, err := calculateFinalRewards(scores, totalScore, specialRewards, totalRewards)
	if err != nil {
		return nil, nil, err
	}

	return rewards, scores, nil
}

// updateNodeStakingData retrieves node information from a staking contract
// and updates the staking pool tokens for each node.
func (s *Server) updateNodeStakingData(nodeAddresses []common.Address, nodes []*schema.Node) error {
	nodeInfo, err := s.stakingContract.GetNodes(&bind.CallOpts{}, nodeAddresses)
	if err != nil {
		return fmt.Errorf("get nodes from chain: %w", err)
	}

	nodeInfoMap := lo.SliceToMap(nodeInfo, func(node l2.DataTypesNode) (common.Address, l2.DataTypesNode) {
		return node.Account, node
	})

	for _, node := range nodes {
		if nodeInfo, ok := nodeInfoMap[node.Address]; ok {
			stakePoolTokens := nodeInfo.StakingPoolTokens
			node.StakingPoolTokens = stakePoolTokens.String()
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

// parsePoolSizes extracts and parses staking pool sizes from nodes.
func parsePoolSizes(nodes []*schema.Node) ([]*big.Int, error) {
	poolSizes := make([]*big.Int, len(nodes))

	for i, node := range nodes {
		poolSize := new(big.Int)
		poolSize, ok := poolSize.SetString(node.StakingPoolTokens, 10)

		if !ok {
			return nil, fmt.Errorf("failed to parse staking pool tokens for node %s: invalid number", node.Address)
		}

		poolSizes[i] = poolSize
	}

	return poolSizes, nil
}

// computeTotalStakeValue calculates the total stake value.
func computeTotalStakeValue(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount) *big.Int {
	var totalStakeValue = big.NewInt(0)

	for _, node := range nodes {
		if _, exist := recentStakers[node.Address]; exist {
			totalStakeValue.Add(totalStakeValue, recentStakers[node.Address].StakeValue.BigInt())
		}
	}

	return totalStakeValue
}

// computeScores calculates the scores for each node based on various factors.
func computeScores(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, poolSizes []*big.Int, totalPoolSize *big.Int, totalStakeValue *big.Int, specialRewards *config.SpecialRewards) []*big.Float {
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

// calculateFinalRewards converts scores into reward amounts.
func calculateFinalRewards(scores []*big.Float, totalScore *big.Float, specialRewards *config.SpecialRewards, totalRewards float64) ([]*big.Int, error) {
	if totalScore.Cmp(big.NewFloat(0)) == 0 {
		return nil, fmt.Errorf("totalScore cannot be zero")
	}

	rewards := make([]*big.Int, len(scores))

	for i, score := range scores {
		// Calculate the ratio of score to totalScore
		scoreRatio := new(big.Float).Quo(score, totalScore)

		// Apply special rewards
		reward := new(big.Float).Mul(scoreRatio, big.NewFloat(0).SetFloat64(totalRewards))
		maxReward := big.NewFloat(0).SetFloat64(specialRewards.RewardsCeiling)

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
