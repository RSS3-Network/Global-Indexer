package settler

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/shopspring/decimal"
)

// calculateOperationRewards calculates the Operation Rewards for all Nodes
// For Alpha, there is no Operation Rewards, but a Special Rewards is calculated
// TODO: Implement the actual calculation logic
func calculateOperationRewards(nodes []*schema.Node) ([]*big.Int, error) {
	operationRewards := make([]*big.Int, len(nodes))

	// For Alpha, set the rewards to 0
	for i := range operationRewards {
		operationRewards[i] = big.NewInt(0)
	}

	return operationRewards, nil
}

// prepareRequestCounts prepares the Request Counts for all Nodes
// For Alpha, there is no actual calculation logic, the counts are set to 0
// TODO: Implement the actual logic to retrieve the counts from the database
func prepareRequestCounts(nodes []common.Address) []*big.Int {
	slice := make([]*big.Int, len(nodes))

	// For Alpha, set the counts to 0
	for i := range slice {
		slice[i] = big.NewInt(0)
	}

	return slice
}

func calculateNodeScore(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, specialRewards *config.SpecialRewards) ([]decimal.Decimal, error) {
	var (
		totalStakeValue *big.Int
		totalPoolSize   = big.NewInt(0)
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
	scoresFloat, err := computeScores(nodes, recentStakers, poolSizes, totalPoolSize, totalStakeValue, specialRewards)
	if err != nil {
		return nil, err
	}

	scores, _ := parseScores(scoresFloat)

	return scores, nil
}

func parseScores(scores []*big.Float) ([]decimal.Decimal, error) {
	scoreDecimals := make([]decimal.Decimal, len(scores))

	for i, score := range scores {
		strValue := score.Text('f', -1)

		decimalValue, err := decimal.NewFromString(strValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse score %d: %w", i, err)
		}

		scoreDecimals[i] = decimalValue
	}

	return scoreDecimals, nil
}
