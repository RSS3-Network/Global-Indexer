package settler

import (
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

// calculateAlphaSpecialRewards calculates the distribution of the Special Rewards used to replace the Operation Rewards
// the Special Rewards are used to incentivize staking in smaller Nodes
// currently, the amount is set to 30,000,000 / 486.6666666666667 * 0.2 ~= 12328
func calculateAlphaSpecialRewards(nodes []*schema.Node, recentStackers map[common.Address]uint64, specialRewards *config.SpecialRewards) ([]*big.Int, error) {
	var (
		totalEffectiveStakers, maxPoolSize uint64
		totalScore                         float64
	)

	// Preprocessing step to avoid repeated parsing and condition checking.
	poolSizes, err := parsePoolSizes(nodes)
	if err != nil {
		return nil, err // Error handling early exit.
	}

	totalEffectiveStakers, maxPoolSize = computeEffectiveStakersAndMaxPoolSize(nodes, recentStackers, poolSizes, specialRewards)

	scores, err := computeScores(nodes, recentStackers, poolSizes, totalEffectiveStakers, maxPoolSize, specialRewards)
	if err != nil {
		return nil, err // Centralized error handling.
	}

	for _, score := range scores {
		totalScore += score
	}

	return calculateFinalRewards(scores, totalScore, specialRewards), nil
}

// parsePoolSizes extracts and parses staking pool sizes from nodes.
func parsePoolSizes(nodes []*schema.Node) ([]uint64, error) {
	poolSizes := make([]uint64, len(nodes))

	for i, node := range nodes {
		poolSize, err := strconv.ParseUint(node.StakingPoolTokens, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse staking pool tokens for node %s: %w", node.Address, err)
		}

		poolSizes[i] = poolSize
	}

	return poolSizes, nil
}

// computeEffectiveStakersAndMaxPoolSize calculates total effective stakers and the maximum pool size.
func computeEffectiveStakersAndMaxPoolSize(nodes []*schema.Node, recentStackers map[common.Address]uint64, poolSizes []uint64, specialRewards *config.SpecialRewards) (uint64, uint64) {
	var totalEffectiveStakers, maxPoolSize uint64

	for i, node := range nodes {
		if poolSizes[i] <= specialRewards.CliffPoint {
			totalEffectiveStakers += recentStackers[node.Address]
		}

		if maxPoolSize < poolSizes[i] {
			maxPoolSize = poolSizes[i]
		}
	}

	return totalEffectiveStakers, maxPoolSize
}

// computeScores calculates the scores for each node based on various factors.
func computeScores(nodes []*schema.Node, recentStackers map[common.Address]uint64, poolSizes []uint64, totalEffectiveStakers, maxPoolSize uint64, specialRewards *config.SpecialRewards) ([]float64, error) {
	scores := make([]float64, len(nodes))

	for i, poolSize := range poolSizes {
		stakers := recentStackers[nodes[i].Address]
		if stakers == 0 {
			continue // Skip computation for nodes with no stakers.
		}

		score := applyGiniCoefficient(poolSize, specialRewards.GiniCoefficient)

		if poolSize > specialRewards.CliffPoint {
			applyCliffFactor(poolSize, maxPoolSize, &score, specialRewards.CliffFactor)
		}

		if totalEffectiveStakers > 0 {
			applyStakerFactor(stakers, totalEffectiveStakers, specialRewards.StakerFactor, &score)
		}

		if score < 0 || score >= 1 {
			return nil, fmt.Errorf("invalid score: %f", score)
		}

		scores[i] = score
	}

	return scores, nil
}

// calculateFinalRewards converts scores into reward amounts.
func calculateFinalRewards(scores []float64, totalScore float64, specialRewards *config.SpecialRewards) []*big.Int {
	rewards := make([]*big.Int, len(scores))
	scale := new(big.Float).SetInt(big.NewInt(1e18)) // Move outside loop to avoid repeated allocation.

	for i, score := range scores {
		reward := math.Trunc(score / totalScore * specialRewards.Rewards)
		rewardBigFloat := new(big.Float).SetFloat64(reward)
		scaledF := new(big.Float).Mul(rewardBigFloat, scale)
		rewardFinal, _ := scaledF.Int(nil) // Simplified conversion to big.Int.
		rewards[i] = rewardFinal
	}

	return rewards
}

// applyGiniCoefficient applies the Gini Coefficient to the score
func applyGiniCoefficient(poolSize uint64, giniCoefficient float64) float64 {
	// Perform calculation: score = 1 / (1 + giniCoefficient * poolSize)
	score := 1 / (1 + giniCoefficient*float64(poolSize))

	return score
}

// applyCliffFactor applies the Cliff Factor to the score
func applyCliffFactor(poolSize uint64, maxPoolSize uint64, score *float64, cliffFactor float64) {
	// Perform calculation: score *= cliffFactor ** poolSize / maxPoolSize
	*score *= math.Pow(cliffFactor, float64(poolSize)/float64(maxPoolSize))
}

// applyStakerFactor applies the Staker Factor to the score
func applyStakerFactor(stakers uint64, totalEffectiveStakers uint64, stakerFactor float64, score *float64) {
	// Perform calculation: score += (score * stakers * staker_factor) / total_stakers
	*score += (*score * float64(stakers) * stakerFactor) / float64(totalEffectiveStakers)
}
