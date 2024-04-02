package settler

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

// calculateOperationRewards calculates the Operation Rewards for all Nodes
// For Alpha, there is no Operation Rewards, but a Special Rewards is calculated
// TODO: Implement the actual calculation logic
func calculateOperationRewards(nodes []*schema.Node, recentStackers map[common.Address]uint64, specialRewards *config.SpecialRewards) ([]*big.Int, error) {
	operationRewards, err := calculateAlphaSpecialRewards(nodes, recentStackers, specialRewards)
	if err != nil {
		return nil, err
	}

	// For Alpha, set the rewards to 0
	//for i := range operationRewards {
	//	operationRewards[i] = big.NewInt(0)
	//}

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
