package settler

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/schema"
)

// calculateActiveScores calculates the active scores for all Nodes
func calculateNodeActiveScores(nodes []*schema.Node, recentStakers map[common.Address]*schema.StakeRecentCount, activeScores *config.ActiveScores) ([]*big.Float, error) {
	// If there are no nodes, return nil
	if len(nodes) == 0 {
		return nil, nil
	}

	scores, err := calculateActiveScores(nodes, recentStakers, activeScores)

	if err != nil {
		return nil, fmt.Errorf("failed to calculate active scores: %w", err)
	}

	return scores, nil
}

// prepareRequestCounts prepares the request counts for all Nodes
func (s *Server) prepareRequestCounts(ctx context.Context, nodes []common.Address) ([]*big.Int, error) {
	if len(nodes) == 0 {
		return make([]*big.Int, 0), nil
	}

	stats, err := s.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		Addresses: nodes,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find node stats: %w", err)
	}

	statsMap := make(map[common.Address]*schema.Stat, len(stats))
	for _, stat := range stats {
		statsMap[stat.Address] = stat
	}

	requestCounts := make([]*big.Int, len(nodes))

	for i, node := range nodes {
		if stat, ok := statsMap[node]; ok {
			// set request counts for nodes from the epoch.
			requestCounts[i] = big.NewInt(stat.EpochRequest)
		} else {
			requestCounts[i] = big.NewInt(0)
		}
	}

	return requestCounts, nil
}
