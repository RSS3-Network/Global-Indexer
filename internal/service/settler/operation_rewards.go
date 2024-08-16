package settler

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/schema"
)

func calculateOperationRewards(nodes []*schema.Node, requestCount []*big.Int, operationRewards *config.OperationRewards) ([]*big.Int, error) {
	// If there are no nodes, return nil
	if len(nodes) == 0 {
		return nil, nil
	}

	rewards, err := calculateFinalRewards(requestCount, operationRewards.Rewards)

	if err != nil {
		return nil, fmt.Errorf("failed to calculate operation rewards: %w", err)
	}

	return rewards, nil
}

// calculateFinalRewards calculates the final rewards for each node based on the request count and total rewards.
func calculateFinalRewards(requestCount []*big.Int, totalRewards float64) ([]*big.Int, error) {
	// Calculate the total request count
	totalRequestCount := big.NewFloat(0)
	for _, count := range requestCount {
		totalRequestCount.Add(totalRequestCount, big.NewFloat(0).SetInt(count))
	}

	// Calculate the rewards for each node
	rewards := make([]*big.Int, len(requestCount))

	for i := range requestCount {
		if requestCount[i].Cmp(big.NewInt(0)) == 0 {
			rewards[i] = big.NewInt(0)

			continue
		}

		count := big.NewFloat(0).SetInt(requestCount[i])
		// Calculate the rewards for the node
		radio := new(big.Float).Quo(count, totalRequestCount)
		reward := new(big.Float).Mul(radio, big.NewFloat(totalRewards))

		// Convert to integer to truncate before scaling
		rewardFinal, _ := reward.Int(nil)

		// Apply gwei after truncation
		scaleGwei(rewardFinal)

		rewards[i] = rewardFinal
	}

	err := checkRewardsCeiling(rewards, totalRewards)
	if err != nil {
		return nil, err
	}

	return rewards, nil
}

// checkRewardsCeiling checks if the sum of rewards is less than or equal to specialRewards.Rewards.
func checkRewardsCeiling(rewards []*big.Int, totalRewards float64) error {
	sum := big.NewInt(0)
	for _, reward := range rewards {
		sum.Add(sum, reward)
	}

	// Scale the operationRewards by 10^18 to match the rewards scale
	operationRewardsBigInt := big.NewInt(0).SetUint64(uint64(totalRewards))
	scaleGwei(operationRewardsBigInt)

	if sum.Cmp(operationRewardsBigInt) > 0 {
		return fmt.Errorf("total rewards exceed the ceiling: %v > %v", sum, operationRewardsBigInt)
	}

	return nil
}

// prepareRequestCounts prepares the request counts for the nodes.
func (s *Server) prepareRequestCounts(ctx context.Context, nodes []common.Address) ([]*big.Int, []*big.Int, error) {
	if len(nodes) == 0 {
		return make([]*big.Int, 0), make([]*big.Int, 0), nil
	}

	stats, err := s.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		Addresses: nodes,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find node stats: %w", err)
	}

	statsMap := make(map[common.Address]*schema.Stat, len(stats))
	for _, stat := range stats {
		statsMap[stat.Address] = stat
	}

	requestCounts := make([]*big.Int, len(nodes))
	totalRequestCounts := make([]*big.Int, len(nodes))

	for i, node := range nodes {
		if stat, ok := statsMap[node]; ok {
			// set request counts for nodes from the epoch.
			requestCounts[i] = big.NewInt(stat.EpochRequest)
			// set total request counts for nodes.
			totalRequestCounts[i] = big.NewInt(stat.TotalRequest)
		} else {
			requestCounts[i] = big.NewInt(0)
			totalRequestCounts[i] = big.NewInt(0)
		}
	}

	return requestCounts, totalRequestCounts, nil
}
