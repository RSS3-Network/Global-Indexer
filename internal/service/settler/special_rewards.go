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
func calculateAlphaSpecialRewards(nodes []*schema.Node, recentStackers map[common.Address]uint64, specialRewards *config.SpecialRewards) ([]*big.Int, error) {
	var (
		totalEffectiveStakers uint64
		maxPoolSize           *big.Float
		totalScore            = big.NewFloat(0)
	)

	// Preprocessing step to avoid repeated parsing and condition checking.
	poolSizes, err := parsePoolSizes(nodes)
	if err != nil {
		return nil, err
	}

	// Calculate total effective stakers and the maximum pool size.
	totalEffectiveStakers, maxPoolSize, err = computeEffectiveStakersAndMaxPoolSize(nodes, recentStackers, poolSizes, specialRewards)
	if err != nil {
		return nil, err
	}

	// Calculate scores for each node.
	scores, err := computeScores(nodes, recentStackers, poolSizes, totalEffectiveStakers, maxPoolSize, specialRewards)
	if err != nil {
		return nil, err
	}

	for _, score := range scores {
		totalScore.Add(totalScore, score)
	}

	return calculateFinalRewards(scores, totalScore, specialRewards), nil
}

// parsePoolSizes extracts and parses staking pool sizes from nodes.
func parsePoolSizes(nodes []*schema.Node) ([]*big.Float, error) {
	poolSizes := make([]*big.Float, len(nodes))

	for i, node := range nodes {
		poolSize := new(big.Float)
		poolSize, ok := poolSize.SetString(node.StakingPoolTokens)

		if !ok {
			return nil, fmt.Errorf("failed to parse staking pool tokens for node %s: invalid number", node.Address)
		}

		poolSizes[i] = poolSize
	}

	return poolSizes, nil
}

// computeEffectiveStakersAndMaxPoolSize calculates total effective stakers and the maximum pool size.
func computeEffectiveStakersAndMaxPoolSize(nodes []*schema.Node, recentStackers map[common.Address]uint64, poolSizes []*big.Float, specialRewards *config.SpecialRewards) (uint64, *big.Float, error) {
	var (
		totalEffectiveStakers uint64
		maxPoolSize           = big.NewFloat(0)
	)

	cliffPoint := big.NewFloat(0)
	_, ok := cliffPoint.SetString(specialRewards.CliffPoint)

	if !ok {
		return 0, nil, fmt.Errorf("CliffPoint conversion failed")
	}

	for i, node := range nodes {
		if poolSizes[i].Cmp(cliffPoint) <= 0 {
			totalEffectiveStakers += recentStackers[node.Address]
		}

		if poolSizes[i].Cmp(maxPoolSize) == 1 {
			maxPoolSize = poolSizes[i]
		}
	}

	return totalEffectiveStakers, maxPoolSize, nil
}

// computeScores calculates the scores for each node based on various factors.
func computeScores(nodes []*schema.Node, recentStackers map[common.Address]uint64, poolSizes []*big.Float, totalEffectiveStakers uint64, maxPoolSize *big.Float, specialRewards *config.SpecialRewards) ([]*big.Float, error) {
	scores := make([]*big.Float, len(nodes))

	for i, poolSize := range poolSizes {
		stakers := recentStackers[nodes[i].Address]
		if stakers == 0 {
			continue
		}

		score := applyGiniCoefficient(poolSize, specialRewards.GiniCoefficient)

		cliffPoint := big.NewFloat(0)
		_, ok := cliffPoint.SetString(specialRewards.CliffPoint)

		if !ok {
			return nil, fmt.Errorf("CliffPoint conversion failed")
		}

		if poolSize.Cmp(cliffPoint) == 1 {
			applyCliffFactor(poolSize, maxPoolSize, score, specialRewards.CliffFactor)
		}

		if totalEffectiveStakers > 0 {
			applyStakerFactor(stakers, totalEffectiveStakers, specialRewards.StakerFactor, score)
		}

		zero := big.NewFloat(0)

		if score.Cmp(zero) < 0 {
			return nil, fmt.Errorf("invalid score: %f", score)
		}

		scores[i] = score
	}

	return scores, nil
}

// calculateFinalRewards converts scores into reward amounts.
func calculateFinalRewards(scores []*big.Float, totalScore *big.Float, specialRewards *config.SpecialRewards) []*big.Int {
	rewards := make([]*big.Int, len(scores))
	// scale is 10^18
	scale := new(big.Float).SetInt(big.NewInt(1e18))

	for i, score := range scores {
		// Perform calculation: reward = score / totalScore * specialRewards.Rewards
		// truncate the reward to an integer to avoid floating point errors
		scoreRatio := new(big.Float).Quo(score, totalScore)
		reward := new(big.Float).Mul(scoreRatio, big.NewFloat(specialRewards.Rewards))
		scaledF := new(big.Float).Mul(reward, scale)
		rewardFinal, _ := scaledF.Int(nil)
		rewards[i] = rewardFinal
	}

	return rewards
}

// applyGiniCoefficient applies the Gini Coefficient to the score
func applyGiniCoefficient(poolSize *big.Float, giniCoefficient float64) *big.Float {
	// Perform calculation: score = 1 / (1 + giniCoefficient * poolSize)
	one := big.NewFloat(1)
	giniTimesPool := new(big.Float).Mul(new(big.Float).SetFloat64(giniCoefficient), poolSize)
	denominator := new(big.Float).Add(one, giniTimesPool)
	score := new(big.Float).Quo(one, denominator)

	return score
}

// applyCliffFactor applies the Cliff Factor to the score
func applyCliffFactor(poolSize *big.Float, maxPoolSize *big.Float, score *big.Float, cliffFactor float64) {
	// Calculate poolSize / maxPoolSize
	poolSizeRatio := new(big.Float).Quo(poolSize, maxPoolSize)

	// Calculate cliffFactor ** poolSizeRatio
	// As big.Float does not support exponentiation directly, using math.Pow after converting to float64
	// For the precision loss here is negligible
	poolSizeRatioFloat64, _ := poolSizeRatio.Float64()

	// Perform calculation: score *= cliffFactor ** poolSize / maxPoolSize
	score.Mul(score, big.NewFloat(math.Pow(cliffFactor, poolSizeRatioFloat64)))
}

// applyStakerFactor applies the Staker Factor to the score
func applyStakerFactor(stakers uint64, totalEffectiveStakers uint64, stakerFactor float64, score *big.Float) {
	totalEffectiveStakersFloat := new(big.Float).SetUint64(totalEffectiveStakers)

	// Perform calculation: score += (score * stakers * staker_factor) / total_stakers
	dividend := big.NewFloat(0)

	dividend.Mul(score, new(big.Float).SetUint64(stakers))
	dividend.Mul(dividend, big.NewFloat(stakerFactor))
	dividend.Quo(dividend, totalEffectiveStakersFloat)
	score.Add(score, dividend)
}
