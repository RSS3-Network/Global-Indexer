package distributor

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
)

// processSecondVerify processes the second verification stage by verifying Activity against working nodes.
// It takes a list of Activity and a list of working nodes' addresses as input parameters.
func (d *Distributor) processSecondVerify(activities []*Activity, workingNodes []common.Address) {
	ctx := context.Background()
	platformMap := make(map[string]struct{})
	statMap := make(map[string]struct{})

	for _, activity := range activities {
		if len(activity.Platform) == 0 {
			continue
		}

		d.verifyPlatform(ctx, activity, platformMap, statMap, workingNodes)

		if _, exists := platformMap[activity.Platform]; !exists {
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
					results[i].InvalidPoint = 1
				} else {
					results[i].ValidPoint = 1
				}
			}
		}
	} else {
		if !results[0].Valid {
			for i := range results {
				results[i].InvalidPoint = 1
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

// verifyPlatform verifies activity against nodes associated with the activity's platform.
// It takes a context, an activity pointer, platform and stat maps, and a list of working nodes' addresses as input parameters.
func (d *Distributor) verifyPlatform(ctx context.Context, activity *Activity, platformMap, statMap map[string]struct{}, workingNodes []common.Address) {
	pid, err := filter.PlatformString(activity.Platform)
	if err != nil {
		return
	}

	worker := PlatformToWorkerMap[pid]

	indexers, err := d.databaseClient.FindNodeIndexers(ctx, nil, []string{activity.Network}, []string{worker})

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

	d.verifyStat(ctx, activity, stats, statMap)

	platformMap[activity.Platform] = struct{}{}
}

// verifyStat verifies Activity statistics and updates them based on comparison with Activity data.
// It takes a context, a Activity pointer, a slice of Node statistics, and a stat map as input parameters.
func (d *Distributor) verifyStat(ctx context.Context, activity *Activity, stats []*schema.Stat, statMap map[string]struct{}) {
	for _, stat := range stats {
		if stat.EpochInvalidRequest >= int64(DefaultSlashCount) {
			continue
		}

		if _, exists := statMap[stat.Address.String()]; !exists {
			statMap[stat.Address.String()] = struct{}{}

			request := dsl.ActivityRequest{
				ID: activity.ID,
			}

			nodeMap, err := d.buildActivityPathByID(
				request,
				[]NodeEndpointCache{
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
				if !d.compareActivities(activity, res.Data) {
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

// compareActivities returns true if two Activity are identical.
// It takes two Activity pointers as input parameters.
// Deprecated: replaced by isActivityIdentical, as inner metadata cannot be efficiently compared.
func (d *Distributor) compareActivities(src, des *Activity) bool {
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
		return results[i].Valid && !results[j].Valid
	})
}

// updateStatsWithResults updates node statistics based on the results of data responses.
// It takes a map of node addresses to their statistics and a slice of data responses as input parameters.
func (d *Distributor) updateStatsWithResults(statsMap map[common.Address]*schema.Stat, results []DataResponse) {
	for _, result := range results {
		if stat, exists := statsMap[result.Address]; exists {
			stat.TotalRequest += int64(result.ValidPoint)
			stat.EpochRequest += int64(result.ValidPoint)
			stat.EpochInvalidRequest += int64(result.InvalidPoint)
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
		results[0].ValidPoint = 2
		results[1].ValidPoint = 1
		results[2].ValidPoint = 1
	} else if !diff01 && diff12 {
		results[0].InvalidPoint = 1
		results[1].ValidPoint = 1
		results[2].ValidPoint = 1
	} else if !diff01 && diff02 {
		results[0].ValidPoint = 2
		results[1].InvalidPoint = 1
		results[2].ValidPoint = 1
	} else if diff01 && !diff02 {
		results[0].ValidPoint = 2
		results[1].ValidPoint = 1
		results[2].InvalidPoint = 1
	} else if !diff01 && !diff02 && !diff12 {
		for i := range results {
			if results[i].Data == nil && results[i].Err != nil {
				results[i].InvalidPoint = 1
			}

			if results[i].Data != nil && results[i].Err == nil {
				results[i].ValidPoint = 1
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
