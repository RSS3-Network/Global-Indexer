package enforcer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
)

// retrieveNodeStatsFromDB retrieves the Node stats from the database.
func retrieveNodeStatsFromDB(ctx context.Context, key string, databaseClient database.Client) ([]*schema.Stat, error) {
	query := schema.StatQuery{
		Limit:        lo.ToPtr(defaultLimit),
		ValidRequest: lo.ToPtr(model.DemotionCountBeforeSlashing),
		PointsOrder:  lo.ToPtr("DESC"),
	}

	var nodeStats []*schema.Stat

	switch key {
	case model.RssNodeCacheKey:
		query.IsRssNode = lo.ToPtr(true)
	case model.FullNodeCacheKey:
		query.IsFullNode = lo.ToPtr(true)
	default:
		return nil, fmt.Errorf("unknown cache key: %s", key)
	}

	for {
		tempNodeStats, err := databaseClient.FindNodeStats(ctx, &query)
		if err != nil || len(tempNodeStats) == 0 {
			break
		}

		qualifiedNodeStats, err := getQualifiedNodes(ctx, tempNodeStats, databaseClient)
		if err != nil {
			return nil, err
		}

		nodeStats = append(nodeStats, qualifiedNodeStats...)
		query.Cursor = lo.ToPtr(tempNodeStats[len(tempNodeStats)-1].Address.String())

		if len(tempNodeStats) < defaultLimit {
			break
		}
	}

	if len(nodeStats) == 0 {
		tempNodeStats, err := databaseClient.FindNodeStats(ctx, &schema.StatQuery{
			Limit:        lo.ToPtr(defaultLimit),
			ValidRequest: lo.ToPtr(model.DemotionCountBeforeSlashing),
			PointsOrder:  lo.ToPtr("DESC"),
		})

		if err != nil {
			return nil, err
		}

		qualifiedNodeStats, err := getQualifiedNodes(ctx, tempNodeStats, databaseClient)
		if err != nil {
			return nil, err
		}

		nodeStats = qualifiedNodeStats[:model.RequiredQualifiedNodeCount]
	}

	return nodeStats, nil
}

// getQualifiedNodes filters the qualified nodes.
func getQualifiedNodes(ctx context.Context, stats []*schema.Stat, databaseClient database.Client) ([]*schema.Stat, error) {
	nodeAddresses := extractNodeAddresses(stats)

	// Retrieve the online Nodes from the database.
	nodes, err := databaseClient.FindNodes(ctx, schema.FindNodesQuery{
		NodeAddresses: nodeAddresses,
		Status:        lo.ToPtr(schema.NodeStatusOnline),
	})

	if err != nil {
		return nil, err
	}

	nodeMap := lo.SliceToMap(nodes, func(node *schema.Node) (common.Address, struct{}) {
		return node.Address, struct{}{}
	})

	var qualifiedNodes []*schema.Stat

	// Exclude the offline nodes.
	for _, stat := range stats {
		if _, exists := nodeMap[stat.Address]; exists {
			qualifiedNodes = append(qualifiedNodes, stat)
		}
	}

	return qualifiedNodes, nil
}

// extractNodeAddresses returns all Node addresses from stats.
func extractNodeAddresses(stats []*schema.Stat) []common.Address {
	return lo.Map(stats, func(stat *schema.Stat, _ int) common.Address {
		return stat.Address
	})
}
