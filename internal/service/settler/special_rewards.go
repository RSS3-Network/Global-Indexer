package settler

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

// calculateAlphaSpecialRewards calculates the distribution of the Special Rewards used to replace the Operation Rewards
// the Special Rewards are used to incentivize staking in smaller Nodes
// currently, the amount is set to 30,000,000 / 486.6666666666667 * 0.2 ~= 12328
func calculateAlphaSpecialRewards(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, specialRewards *config.SpecialRewards) ([]*big.Int, error) {
	var (
		totalStakeValue *big.Int
		totalPoolSize   = big.NewInt(0)
		totalScore      = big.NewFloat(0)
	)

	// Preprocessing step to avoid repeated parsing and condition checking.
	poolSizes, err := parsePoolSizes(nodes)
	if err != nil {
		return nil, err
	}

	// Calculate the total pool size.
	for _, poolSize := range poolSizes {
		totalPoolSize.Add(totalPoolSize, poolSize)
	}

	// Calculate total stake value.
	totalStakeValue, err = computeTotalStakeValue(nodes, recentStakers)
	if err != nil {
		return nil, err
	}

	// Calculate scores for each node.
	scores, err := computeScores(nodes, recentStakers, poolSizes, totalPoolSize, totalStakeValue, specialRewards)
	if err != nil {
		return nil, err
	}

	for _, score := range scores {
		totalScore.Add(totalScore, score)
	}

	rewards, err := calculateFinalRewards(scores, totalScore, specialRewards)
	if err != nil {
		return nil, err
	}

	return rewards, nil
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

// computeTotalStakersAndTotalStakeValue calculates the total effective stakers and the stake value.
func computeTotalStakeValue(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount) (*big.Int, error) {
	var totalStakeValue = big.NewInt(0)

	for _, node := range nodes {
		if _, exist := recentStakers[node.Address]; exist {
			totalStakeValue.Add(totalStakeValue, recentStakers[node.Address].StakeValue.BigInt())
		}
	}

	return totalStakeValue, nil
}

// computeScores calculates the scores for each node based on various factors.
func computeScores(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, poolSizes []*big.Int, totalPoolSize *big.Int, totalStakeValue *big.Int, specialRewards *config.SpecialRewards) ([]*big.Float, error) {
	scores := make([]*big.Float, len(nodes))

	for i, poolSize := range poolSizes {
		stakers := recentStakers[nodes[i].Address]

		if stakers == nil || stakers.StakerCount == 0 {
			scores[i] = big.NewFloat(0)
			continue
		}

		poolSizeRatio := new(big.Float).Quo(new(big.Float).SetInt(poolSize), new(big.Float).SetInt(totalPoolSize))
		stakeRadio := new(big.Float).Quo(new(big.Float).SetInt(stakers.StakeValue.BigInt()), new(big.Float).SetInt(totalStakeValue))

		score := applyGiniCoefficient(poolSizeRatio, specialRewards.GiniCoefficient)

		applyStakerFactor(stakers.StakerCount, stakeRadio, specialRewards.StakerFactor, score)

		scores[i] = score
	}

	return scores, nil
}

// calculateFinalRewards converts scores into reward amounts.
func calculateFinalRewards(scores []*big.Float, totalScore *big.Float, specialRewards *config.SpecialRewards) ([]*big.Int, error) {
	if totalScore.Cmp(big.NewFloat(0)) == 0 {
		return nil, fmt.Errorf("totalScore cannot be zero")
	}

	rewards := make([]*big.Int, len(scores))

	for i, score := range scores {
		// Calculate the ratio of score to totalScore
		scoreRatio := new(big.Float).Quo(score, totalScore)

		// Apply special rewards
		reward := new(big.Float).Mul(scoreRatio, big.NewFloat(0).SetFloat64(specialRewards.Rewards))

		if reward.Cmp(big.NewFloat(0).SetFloat64(specialRewards.RewardsCeiling)) == 1 {
			reward = big.NewFloat(0).SetFloat64(specialRewards.RewardsCeiling)
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
	stakersFloat := new(big.Float).SetUint64(stakers)                             // Convert stakers to big.Float for calculation.
	stakerFactorFloat := big.NewFloat(stakerFactor)                               // Ensure stakerFactor is in big.Float for consistency in operations.
	stakerFactorCalculated := new(big.Float).Mul(stakerFactorFloat, stakersFloat) // stakerFactor * stakeRatio

	e, _ := new(big.Float).Mul(stakerFactorCalculated, stakeRadio).Float64()
	// Perform calculation: score *= math.exp(stakerFactorCalculated * stakeRatio)
	expResultBigFloat := new(big.Float).SetFloat64(math.Exp(e))
	score.Mul(score, expResultBigFloat)
}
