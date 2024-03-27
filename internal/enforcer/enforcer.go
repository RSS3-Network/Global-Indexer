package enforcer

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/httpx"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

const (
	stakingToPointsRate           float64 = 100000
	stakingLogBase                        = 2
	stakingMaxPoints                      = 0.2
	hoursPerEpoch                         = 18
	activeTimeToPointsRate                = 120
	activeTimeMaxPoints                   = 0.3
	totalReqToPointsRate                  = 100000
	totalReqLogBase                       = 100
	totalReqMaxPoints                     = 0.3
	totalEpochReqToPointsRate             = 1000000
	totalEpochReqLogBase                  = 5000
	totalEpochReqMaxPoints                = 1
	perDecentralizedNetworkPoints         = 0.1
	perRssNetworkPoints                   = 0.3
	perFederatedNetworkPoints             = 0.1
	perIndexerPoints                      = 0.05
	indexerMaxPoints                      = 0.2
	perSlashPoints                        = 0.5
	nonExistPoints                float64 = 0
	existPoints                           = 1

	defaultLimit = 100
)

type Enforcer interface {
	Verify(ctx context.Context, results []model.DataResponse) error
	PartialVerify(ctx context.Context, results []model.DataResponse) error
	MaintainScore(ctx context.Context) error
	ChallengeStates(ctx context.Context) error
}

type SimpleEnforcer struct {
	databaseClient  database.Client
	cacheClient     cache.Client
	httpClient      httpx.Client
	stakingContract *l2.Staking
}

func (e *SimpleEnforcer) Verify(ctx context.Context, results []model.DataResponse) error {
	nodeStatsMap, err := e.getNodeStatsMap(ctx, results)
	if err != nil {
		return fmt.Errorf("find node stats: %w", err)
	}

	e.sortResults(results)

	if len(nodeStatsMap) < model.DefaultNodeCount {
		for i := range results {
			if _, exists := nodeStatsMap[results[i].Address]; exists {
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
			e.updateRequestsBasedOnDataCompare(results)
		}
	}

	e.updateStatsWithResults(nodeStatsMap, results)

	if err = e.databaseClient.SaveNodeStats(ctx, statsMapToSlice(nodeStatsMap)); err != nil {
		return fmt.Errorf("save node stats: %w", err)
	}

	return nil
}

func (e *SimpleEnforcer) getNodeStatsMap(ctx context.Context, results []model.DataResponse) (map[common.Address]*schema.Stat, error) {
	stats, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: lo.Map(results, func(result model.DataResponse, _ int) common.Address {
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

func (e *SimpleEnforcer) sortResults(results []model.DataResponse) {
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].First && !results[j].First
	})
}

func (e *SimpleEnforcer) updateStatsWithResults(statsMap map[common.Address]*schema.Stat, results []model.DataResponse) {
	for _, result := range results {
		if stat, exists := statsMap[result.Address]; exists {
			stat.TotalRequest += int64(result.Request)
			stat.EpochRequest += int64(result.Request)
			stat.EpochInvalidRequest += int64(result.InvalidRequest)
		}
	}
}

func (e *SimpleEnforcer) updateRequestsBasedOnDataCompare(results []model.DataResponse) {
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

func statsMapToSlice(statsMap map[common.Address]*schema.Stat) []*schema.Stat {
	statsSlice := make([]*schema.Stat, 0, len(statsMap))
	for _, stat := range statsMap {
		statsSlice = append(statsSlice, stat)
	}

	return statsSlice
}

func compareData(src, des []byte) bool {
	if src == nil || des == nil {
		return false
	}

	srcHash, destHash := sha256.Sum256(src), sha256.Sum256(des)

	return string(srcHash[:]) == string(destHash[:])
}

func (e *SimpleEnforcer) PartialVerify(ctx context.Context, results []model.DataResponse) error {
	activities, err := e.extractActivities(ctx, results)

	if err != nil {
		return err
	}

	workingNodes := lo.Map(results, func(result model.DataResponse, _ int) common.Address {
		return result.Address
	})

	e.verifyFeeds(ctx, activities.Data, workingNodes)

	return nil
}

func (e *SimpleEnforcer) extractActivities(_ context.Context, results []model.DataResponse) (*model.ActivitiesResponse, error) {
	if !results[0].First {
		// TODO return error
		return nil, nil
	}

	var activities *model.ActivitiesResponse

	data := results[0].Data

	if err := json.Unmarshal(data, &activities); err != nil {
		zap.L().Error("fail to unmarshall activities")

		return nil, err
	}

	// data is empty, no need to 2nd verify
	if activities.Data == nil {
		// TODO return error
		return nil, nil
	}

	return activities, nil
}

func (e *SimpleEnforcer) verifyFeeds(ctx context.Context, feeds []*model.Feed, workingNodes []common.Address) {
	platformMap := make(map[string]struct{})
	statMap := make(map[string]struct{})

	for _, feed := range feeds {
		if len(feed.Platform) == 0 {
			continue
		}

		_ = e.verifyPlatform(ctx, feed, platformMap, statMap, workingNodes)

		if _, exists := platformMap[feed.Platform]; !exists {
			if len(platformMap) == model.DefaultVerifyCount {
				break
			}
		}
	}
}

func (e *SimpleEnforcer) verifyPlatform(ctx context.Context, feed *model.Feed, platformMap, statMap map[string]struct{}, workingNodes []common.Address) error {
	pid, err := filter.PlatformString(feed.Platform)
	if err != nil {
		return err
	}

	worker := model.PlatformToWorkerMap[pid]

	indexers, err := e.databaseClient.FindNodeIndexers(ctx, nil, []string{feed.Network}, []string{worker})

	if err != nil {
		return err
	}

	nodeAddresses := lo.Map(indexers, func(indexer *schema.Indexer, _ int) common.Address {
		return indexer.Address
	})

	nodeAddresses = lo.Filter(nodeAddresses, func(item common.Address, _ int) bool {
		return !lo.Contains(workingNodes, item)
	})

	if len(nodeAddresses) == 0 {
		return nil
	}

	stats, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: nodeAddresses,
		PointsOrder: lo.ToPtr("DESC"),
	})

	if err != nil || len(stats) == 0 {
		return nil
	}

	_ = e.verifyStat(ctx, feed, stats, statMap)

	platformMap[feed.Platform] = struct{}{}

	return nil
}

func (e *SimpleEnforcer) verifyStat(ctx context.Context, feed *model.Feed, stats []*schema.Stat, statMap map[string]struct{}) error {
	for _, stat := range stats {
		if stat.EpochInvalidRequest >= int64(model.DefaultSlashCount) {
			continue
		}

		if _, exists := statMap[stat.Address.String()]; !exists {
			statMap[stat.Address.String()] = struct{}{}

			res, err := e.fetchActivityByTxID(ctx, stat.Endpoint, feed.ID)

			if err != nil || res == nil {
				stat.EpochInvalidRequest++
			} else {
				if !e.compareFeeds(feed, res.Data) {
					stat.EpochInvalidRequest++
				} else {
					stat.TotalRequest++
					stat.EpochRequest++
				}
			}

			_ = e.databaseClient.SaveNodeStat(ctx, stat)

			break
		}
	}

	return nil
}

func (e *SimpleEnforcer) fetchActivityByTxID(ctx context.Context, endpoint, txID string) (*model.ActivityResponse, error) {
	fullURL, err := url.JoinPath(endpoint, fmt.Sprintf("/decentralized/tx/%s", txID))

	if err != nil {
		return nil, fmt.Errorf("failed to join path for node %s: %w", endpoint, err)
	}

	decodedURL, err := url.QueryUnescape(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to unescape url for node %s: %w", endpoint, err)
	}

	body, err := e.httpClient.Fetch(ctx, decodedURL)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var (
		res      model.ActivityResponse
		errRes   model.ErrResponse
		notFound model.NotFoundResponse
	)

	if err = json.Unmarshal(data, &errRes); err != nil {
		return nil, err
	}

	if errRes.ErrorCode != "" {
		return nil, nil
	}

	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &notFound); err != nil {
		return nil, err
	}

	if notFound.Message != "" {
		return nil, nil
	}

	return &res, nil
}

func (e *SimpleEnforcer) compareFeeds(src, des *model.Feed) bool {
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

func (e *SimpleEnforcer) MaintainScore(ctx context.Context) error {
	var currentEpoch int64

	epochEvent, err := e.databaseClient.FindEpochs(ctx, 1, nil)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		zap.L().Error("get latest epoch event from database", zap.Error(err))

		return err
	}

	if len(epochEvent) > 0 {
		currentEpoch = int64(epochEvent[0].ID)
	}

	query := &schema.StatQuery{
		Limit: lo.ToPtr(defaultLimit),
	}

	for first := true; query.Cursor != nil || first; first = false {
		stats, err := e.databaseClient.FindNodeStats(ctx, query)

		if err != nil {
			return err
		}

		statsPool := pool.New().
			WithContext(ctx).
			WithCancelOnError().
			WithFirstError()

		for _, stat := range stats {
			stat := stat

			statsPool.Go(func(ctx context.Context) error {
				if err = e.updateNodeEpochStats(stat, currentEpoch); err != nil {
					return err
				}

				return e.updateNodePoints(stat)
			})
		}

		if err := statsPool.Wait(); err != nil {
			return fmt.Errorf("wait stats pool: %w", err)
		}

		if err = e.databaseClient.SaveNodeStats(ctx, stats); err != nil {
			return err
		}

		if len(stats) == 0 {
			break
		}

		lastStat, _ := lo.Last(stats)
		query.Cursor = lo.ToPtr(lastStat.Address.String())
	}

	return e.updateNodeCache(ctx)
}

func (e *SimpleEnforcer) updateNodeEpochStats(stat *schema.Stat, epoch int64) error {
	nodeInfo, err := e.stakingContract.GetNode(&bind.CallOpts{}, stat.Address)

	if err != nil {
		return fmt.Errorf("get node info: %s,%w", stat.Address.String(), err)
	}

	stat.Staking = float64(nodeInfo.StakingPoolTokens.Uint64())
	if epoch != stat.Epoch {
		stat.EpochRequest = 0
		stat.EpochInvalidRequest = 0
		stat.Epoch = epoch
	}

	return nil
}

func (e *SimpleEnforcer) updateNodePoints(stat *schema.Stat) error {
	node, err := e.databaseClient.FindNode(context.Background(), stat.Address)

	if err != nil {
		return fmt.Errorf("find node: %s, %w", stat.Address.String(), err)
	}

	if node.Status == schema.NodeStatusOffline {
		stat.ResetAt = time.Now()
		stat.EpochInvalidRequest = int64(model.DefaultSlashCount)

		return nil
	}

	e.calcPoints(stat)

	return nil
}

func (e *SimpleEnforcer) updateNodeCache(ctx context.Context) error {
	rssNodes, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		IsRssNode:    lo.ToPtr(true),
		ValidRequest: lo.ToPtr(model.DefaultSlashCount),
		Limit:        lo.ToPtr(model.DefaultNodeCount),
	})

	if err != nil {
		return err
	}

	if err = e.setNodeCache(ctx, model.RssNodeCacheKey, rssNodes); err != nil {
		return err
	}

	fullNodes, err := e.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		IsFullNode:   lo.ToPtr(true),
		ValidRequest: lo.ToPtr(model.DefaultSlashCount),
		Limit:        lo.ToPtr(model.DefaultNodeCount),
	})
	if err != nil {
		return err
	}

	return e.setNodeCache(ctx, model.FullNodeCacheKey, fullNodes)
}

func (e *SimpleEnforcer) setNodeCache(ctx context.Context, key string, stats []*schema.Stat) error {
	nodesCache := lo.Map(stats, func(n *schema.Stat, _ int) model.Cache {
		return model.Cache{Address: n.Address.String(), Endpoint: n.Endpoint}
	})

	if err := e.cacheClient.Set(ctx, key, nodesCache); err != nil {
		return fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nil
}

func (e *SimpleEnforcer) calcPoints(stat *schema.Stat) {
	// staking pool tokens
	stat.Points = math.Min(math.Log(stat.Staking/stakingToPointsRate+1)/math.Log(stakingLogBase), stakingMaxPoints)

	// public good node
	stat.Points += lo.Ternary(stat.IsPublicGood, nonExistPoints, existPoints)

	// node active time
	stat.Points += math.Min(math.Ceil(time.Since(stat.ResetAt).Hours()/hoursPerEpoch)/activeTimeToPointsRate, activeTimeMaxPoints)

	// total requests
	stat.Points += math.Min(math.Log(float64(stat.TotalRequest)/totalReqToPointsRate+1)/math.Log(totalReqLogBase), totalReqMaxPoints)

	// epoch requests
	stat.Points += math.Min(math.Log(float64(stat.EpochRequest)/totalEpochReqToPointsRate+1)/math.Log(totalEpochReqLogBase), totalEpochReqMaxPoints)

	// network count
	stat.Points += perDecentralizedNetworkPoints*float64(stat.DecentralizedNetwork+stat.FederatedNetwork) + perRssNetworkPoints*lo.Ternary(stat.IsRssNode, existPoints, nonExistPoints)

	// indexer count
	stat.Points += math.Min(float64(stat.Indexer)*perIndexerPoints, indexerMaxPoints)

	// epoch failure requests
	stat.Points -= perSlashPoints * float64(stat.EpochInvalidRequest)
}

func (e *SimpleEnforcer) ChallengeStates(_ context.Context) error {
	return nil
}

func NewSimpleEnforcer(databaseClient database.Client, httpClient httpx.Client, cacheClient cache.Client, stakingContract *l2.Staking) (*SimpleEnforcer, error) {
	return &SimpleEnforcer{
		databaseClient:  databaseClient,
		httpClient:      httpClient,
		cacheClient:     cacheClient,
		stakingContract: stakingContract,
	}, nil
}
