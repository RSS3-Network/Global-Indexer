package settler

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/schema"
)

// calculateOperationRewards calculates the Operation Rewards for all Nodes
// For Alpha, there is no Operation Rewards, but a Special Rewards is calculated
// TODO: Implement the actual calculation logic
func calculateOperationRewards(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, specialRewards *config.SpecialRewards) ([]*big.Int, []*big.Float, error) {
	// If there are no nodes, return nil
	if len(nodes) == 0 {
		return nil, nil, nil
	}

	operationRewards, activeScores, err := calculateAlphaSpecialRewards(nodes, recentStakers, specialRewards)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate special rewards: %w", err)
	}

	// Pause AlphaSpecialRewards
	for i := range operationRewards {
		operationRewards[i] = big.NewInt(0)
	}

	return operationRewards, activeScores, nil
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
