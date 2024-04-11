package distributor

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/hub/model/dsl"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
)

// processSecondVerify processes the second verification stage by verifying feeds against working nodes.
// It takes a list of feeds and a list of working nodes' addresses as input parameters.
func (d *Distributor) processSecondVerify(feeds []*Feed, workingNodes []common.Address) {
	ctx := context.Background()
	platformMap := make(map[string]struct{})
	statMap := make(map[string]struct{})

	for _, feed := range feeds {
		if len(feed.Platform) == 0 {
			continue
		}

		d.verifyPlatform(ctx, feed, platformMap, statMap, workingNodes)

		if _, exists := platformMap[feed.Platform]; !exists {
			if len(platformMap) == DefaultVerifyCount {
				break
			}
		}
	}
}

// verifyData verifies the data responses and updates node statistics accordingly.
// It takes a context and a slice of data responses as input parameters.
// It returns an error if any occurred during the verification process.
func (d *Distributor) verifyData(ctx context.Context, results []DataResponse) error {
	statsMap, err := d.getNodeStatsMap(ctx, results)
	if err != nil {
		return fmt.Errorf("find node stats: %w", err)
	}

	d.sortResults(results)

	if len(statsMap) < DefaultNodeCount {
		for i := range results {
			if _, exists := statsMap[results[i].Address]; exists {
				if results[i].Err != nil {
					results[i].InvalidRequest = 1
				} else {
					results[i].Request = 1
				}
			}
		}
	} else {
		if !results[0].First {
			for i := range results {
				results[i].InvalidRequest = 1
			}
		} else {
			d.updateRequestsBasedOnDataCompare(results)
		}
	}

	d.updateStatsWithResults(statsMap, results)

	if err = d.databaseClient.SaveNodeStats(ctx, lo.MapToSlice(statsMap, func(_ common.Address, value *schema.Stat) *schema.Stat {
		return value
	})); err != nil {
		return fmt.Errorf("save node stats: %w", err)
	}

	return nil
}

// verifyPlatform verifies feeds against nodes associated with the feed's platform.
// It takes a context, a feed pointer, platform and stat maps, and a list of working nodes' addresses as input parameters.
func (d *Distributor) verifyPlatform(ctx context.Context, feed *Feed, platformMap, statMap map[string]struct{}, workingNodes []common.Address) {
	pid, err := filter.PlatformString(feed.Platform)
	if err != nil {
		return
	}

	worker := PlatformToWorkerMap[pid]

	indexers, err := d.databaseClient.FindNodeIndexers(ctx, nil, []string{feed.Network}, []string{worker})

	if err != nil {
		return
	}

	nodeAddresses := lo.Map(indexers, func(indexer *schema.Indexer, _ int) common.Address {
		return indexer.Address
	})

	nodeAddresses = lo.Filter(nodeAddresses, func(item common.Address, _ int) bool {
		return !lo.Contains(workingNodes, item)
	})

	if len(nodeAddresses) == 0 {
		return
	}

	stats, err := d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: nodeAddresses,
		PointsOrder: lo.ToPtr("DESC"),
	})

	if err != nil || len(stats) == 0 {
		return
	}

	d.verifyStat(ctx, feed, stats, statMap)

	platformMap[feed.Platform] = struct{}{}
}

// verifyStat verifies feed statistics and updates them based on comparison with feed data.
// It takes a context, a feed pointer, a slice of feed statistics, and a stat map as input parameters.
func (d *Distributor) verifyStat(ctx context.Context, feed *Feed, stats []*schema.Stat, statMap map[string]struct{}) {
	for _, stat := range stats {
		if stat.EpochInvalidRequest >= int64(DefaultSlashCount) {
			continue
		}

		if _, exists := statMap[stat.Address.String()]; !exists {
			statMap[stat.Address.String()] = struct{}{}

			request := dsl.ActivityRequest{
				ID: feed.ID,
			}

			nodeMap, err := d.buildActivityPathByID(
				request,
				[]Cache{
					{
						Address:  stat.Address.String(),
						Endpoint: stat.Endpoint,
					},
				},
			)

			if err != nil {
				continue
			}

			data, err := d.fetch(ctx, nodeMap[stat.Address])

			flag, res := d.validateActivity(data)

			if err != nil || !flag {
				stat.EpochInvalidRequest++
			} else {
				if !d.compareFeeds(feed, res.Data) {
					stat.EpochInvalidRequest++
				} else {
					stat.TotalRequest++
					stat.EpochRequest++
				}
			}

			_ = d.databaseClient.SaveNodeStat(ctx, stat)

			break
		}
	}
}

// compareFeeds compares two feed objects and returns true if they are identical, otherwise false.
// It takes two feed pointers as input parameters.
func (d *Distributor) compareFeeds(src, des *Feed) bool {
	var flag bool

	if src.ID != des.ID ||
		src.Network != des.Network ||
		src.Index != des.Index ||
		src.From != des.From ||
		src.To != des.To ||
		src.Tag != des.Tag ||
		src.Type != des.Type ||
		src.Platform != des.Platform ||
		len(src.Actions) != len(des.Actions) {
		return false
	}

	if len(src.Actions) > 0 {
		srcAction := src.Actions[0]

		for _, action := range des.Actions {
			if srcAction.From == action.From &&
				srcAction.To == action.To &&
				srcAction.Tag == action.Tag &&
				srcAction.Type == action.Type {
				desMetadata, _ := json.Marshal(action.Metadata)
				srcMetadata, _ := json.Marshal(srcAction.Metadata)

				if compareData(srcMetadata, desMetadata) {
					flag = true
				}
			}
		}
	}

	return flag
}

// getNodeStatsMap retrieves and returns node statistics mapped by their addresses.
// It takes a context and a slice of data responses as input parameters.
// It returns a map of node addresses to their statistics and an error if any occurred.
func (d *Distributor) getNodeStatsMap(ctx context.Context, results []DataResponse) (map[common.Address]*schema.Stat, error) {
	stats, err := d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: lo.Map(results, func(result DataResponse, _ int) common.Address {
			return result.Address
		}),
		PointsOrder: lo.ToPtr("DESC"),
	})

	if err != nil {
		return nil, err
	}

	statsMap := make(map[common.Address]*schema.Stat)

	for _, stat := range stats {
		statsMap[stat.Address] = stat
	}

	return statsMap, nil
}

// sortResults sorts data responses based on their first status.
// It takes a slice of data responses as an input parameter.
func (d *Distributor) sortResults(results []DataResponse) {
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].First && !results[j].First
	})
}

// updateStatsWithResults updates node statistics based on the results of data responses.
// It takes a map of node addresses to their statistics and a slice of data responses as input parameters.
func (d *Distributor) updateStatsWithResults(statsMap map[common.Address]*schema.Stat, results []DataResponse) {
	for _, result := range results {
		if stat, exists := statsMap[result.Address]; exists {
			stat.TotalRequest += int64(result.Request)
			stat.EpochRequest += int64(result.Request)
			stat.EpochInvalidRequest += int64(result.InvalidRequest)
		}
	}
}

// updateRequestsBasedOnDataCompare updates data response requests based on comparison results.
// It takes a slice of data responses as input parameters.
func (d *Distributor) updateRequestsBasedOnDataCompare(results []DataResponse) {
	diff01 := compareData(results[0].Data, results[1].Data)
	diff02 := compareData(results[0].Data, results[2].Data)
	diff12 := compareData(results[1].Data, results[2].Data)

	if diff01 && diff02 {
		results[0].Request = 2
		results[1].Request = 1
		results[2].Request = 1
	} else if !diff01 && diff12 {
		results[0].InvalidRequest = 1
		results[1].Request = 1
		results[2].Request = 1
	} else if !diff01 && diff02 {
		results[0].Request = 2
		results[1].InvalidRequest = 1
		results[2].Request = 1
	} else if diff01 && !diff02 {
		results[0].Request = 2
		results[1].Request = 1
		results[2].InvalidRequest = 1
	} else if !diff01 && !diff02 && !diff12 {
		for i := range results {
			if results[i].Data == nil && results[i].Err != nil {
				results[i].InvalidRequest = 1
			}

			if results[i].Data != nil && results[i].Err == nil {
				results[i].Request = 1
			}
		}
	}
}

// compareData compares two byte slices and returns true if they are identical, otherwise false.
// It takes two byte slices as input parameters.
func compareData(src, des []byte) bool {
	if src == nil || des == nil {
		return false
	}

	srcHash, destHash := sha256.Sum256(src), sha256.Sum256(des)

	return string(srcHash[:]) == string(destHash[:])
}
