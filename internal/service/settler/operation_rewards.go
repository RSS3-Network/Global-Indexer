package settler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-version"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
)

func (s *Server) calculateOperationRewards(ctx context.Context, operationStats []*schema.Stat, rewards *config.Rewards) ([]*big.Int, error) {
	// If there are no nodes, return nil
	if len(operationStats) == 0 {
		return nil, nil
	}

	operationRewards, err := s.calculateFinalRewards(ctx, operationStats, rewards)

	if err != nil {
		return nil, fmt.Errorf("failed to calculate operation rewards: %w", err)
	}

	return operationRewards, nil
}

type StatValue struct {
	validCount, invalidCount, networkCount, indexerCount, activityCount, upTime *big.Float
	isLatestVersion                                                             bool
}

// calculateFinalRewards calculates the final rewards for each node based on the operation stats.
func (s *Server) calculateFinalRewards(ctx context.Context, operationStats []*schema.Stat, rewards *config.Rewards) ([]*big.Int, error) {
	operationRewards := make([]*big.Int, len(operationStats))
	// Separate RSSHub nodes and normal nodes
	var rsshubStats, normalStats []*schema.Stat

	var rsshubIndices, normalIndices []int

	for i, stat := range operationStats {
		if stat != nil && stat.IsRsshubNode {
			rsshubStats = append(rsshubStats, stat)
			rsshubIndices = append(rsshubIndices, i)
		} else if stat != nil {
			normalStats = append(normalStats, stat)
			normalIndices = append(normalIndices, i)
		}
	}

	// Calculate rewards for RSSHub nodes (50% of total rewards)
	if len(rsshubStats) > 0 {
		rsshubRewards := s.calculateRsshubRewards(rsshubStats, rewards.OperationRewards*0.5)
		for i, reward := range rsshubRewards {
			operationRewards[rsshubIndices[i]] = reward
		}
	}

	// Calculate rewards for normal nodes (50% of total rewards)
	if len(normalStats) > 0 {
		normalRewards := s.calculateNormalNodeRewards(ctx, normalStats, rewards.OperationRewards*0.5, rewards)
		for i, reward := range normalRewards {
			operationRewards[normalIndices[i]] = reward
		}
	}

	// Check if total rewards exceed the ceiling
	if err := checkRewardsCeiling(operationRewards, rewards.OperationRewards); err != nil {
		return nil, err
	}

	return operationRewards, nil
}

// calculateRsshubRewards calculates rewards for RSSHub nodes based on EpochRequest ratio
func (s *Server) calculateRsshubRewards(rsshubStats []*schema.Stat, totalRewards float64) []*big.Int {
	rewards := make([]*big.Int, len(rsshubStats))

	// Calculate total requests
	totalRequests := int64(0)
	for _, stat := range rsshubStats {
		totalRequests += stat.EpochRequest
	}

	// If no requests, all rewards are 0
	if totalRequests == 0 {
		for i := range rewards {
			rewards[i] = big.NewInt(0)
		}

		return rewards
	}

	// Distribute rewards based on request ratio
	for i, stat := range rsshubStats {
		ratio := new(big.Float).Quo(big.NewFloat(float64(stat.EpochRequest)), big.NewFloat(float64(totalRequests)))
		reward := new(big.Float).Mul(ratio, big.NewFloat(totalRewards))
		rewardFinal, _ := reward.Int(nil)
		scaleGwei(rewardFinal)
		rewards[i] = rewardFinal
	}

	return rewards
}

// calculateNormalNodeRewards calculates rewards for normal nodes using the original complex scoring logic
func (s *Server) calculateNormalNodeRewards(ctx context.Context, operationStats []*schema.Stat, totalRewards float64, rewards *config.Rewards) []*big.Int {
	operationRewards := make([]*big.Int, len(operationStats))
	maxStatValue := StatValue{
		validCount:    big.NewFloat(0),
		invalidCount:  big.NewFloat(0),
		networkCount:  big.NewFloat(0),
		indexerCount:  big.NewFloat(0),
		activityCount: big.NewFloat(0),
		upTime:        big.NewFloat(0),
	}
	statValues := make([]StatValue, len(operationStats))

	var mu sync.Mutex

	s.processStat(ctx, operationStats, &maxStatValue, &statValues, &mu)

	scores, totalScore := calculateScores(ctx, operationStats, statValues, maxStatValue, rewards, &mu)

	for i, score := range scores {
		if score.Cmp(big.NewFloat(0)) == 0 {
			operationRewards[i] = big.NewInt(0)
			continue
		}

		reward := new(big.Float).Mul(new(big.Float).Quo(score, totalScore), big.NewFloat(totalRewards))
		rewardFinal, _ := reward.Int(nil)
		scaleGwei(rewardFinal)
		operationRewards[i] = rewardFinal
	}

	return operationRewards
}

// processStat processes the stat for the operation rewards calculation.
func (s *Server) processStat(ctx context.Context, operationStats []*schema.Stat, maxValues *StatValue, statsData *[]StatValue, mu *sync.Mutex) {
	latestVersionStr, _ := s.getNodeLatestVersion()
	latestVersion := version.Must(version.NewVersion(latestVersionStr))
	now := time.Now()

	errorPool := pool.New().WithContext(ctx).WithMaxGoroutines(30).WithCancelOnError().WithFirstError()

	for i := range operationStats {
		i := i

		errorPool.Go(func(_ context.Context) error {
			if operationStats[i] == nil || operationStats[i].EpochInvalidRequest >= int64(model.DemotionCountBeforeSlashing) {
				return nil
			}

			(*statsData)[i].validCount = big.NewFloat(float64(operationStats[i].EpochRequest))
			(*statsData)[i].invalidCount = big.NewFloat(float64(operationStats[i].EpochInvalidRequest))
			(*statsData)[i].networkCount = big.NewFloat(float64(operationStats[i].DecentralizedNetwork + operationStats[i].FederatedNetwork))
			(*statsData)[i].indexerCount = big.NewFloat(float64(operationStats[i].Indexer))

			activityCountResp, err := s.getNodeActivityCount(context.Background(), operationStats[i].Version, operationStats[i].Endpoint, operationStats[i].AccessToken)
			if err != nil {
				(*statsData)[i].activityCount = big.NewFloat(0)
			} else {
				(*statsData)[i].activityCount = big.NewFloat(float64(activityCountResp.Count))
			}

			(*statsData)[i].upTime = big.NewFloat(now.Sub(operationStats[i].ResetAt).Seconds())
			(*statsData)[i].isLatestVersion = version.Must(version.NewVersion(operationStats[i].Version)).GreaterThanOrEqual(latestVersion)

			mu.Lock()
			maxValues.validCount = maxFloat(maxValues.validCount, (*statsData)[i].validCount)
			maxValues.invalidCount = maxFloat(maxValues.invalidCount, (*statsData)[i].invalidCount)
			maxValues.networkCount = maxFloat(maxValues.networkCount, (*statsData)[i].networkCount)
			maxValues.indexerCount = maxFloat(maxValues.indexerCount, (*statsData)[i].indexerCount)
			maxValues.activityCount = maxFloat(maxValues.activityCount, (*statsData)[i].activityCount)
			maxValues.upTime = maxFloat(maxValues.upTime, (*statsData)[i].upTime)
			mu.Unlock()

			return nil
		})
	}

	_ = errorPool.Wait()
}

// calculateScores calculates the scores for the operation rewards calculation.
func calculateScores(ctx context.Context, operationStats []*schema.Stat, statsData []StatValue, maxValues StatValue, rewards *config.Rewards, mu *sync.Mutex) ([]*big.Float, *big.Float) {
	scores := make([]*big.Float, len(operationStats))
	totalScore := big.NewFloat(0)

	errorPool := pool.New().WithContext(ctx).WithMaxGoroutines(30).WithCancelOnError().WithFirstError()

	for i := range statsData {
		i := i

		errorPool.Go(func(_ context.Context) error {
			if operationStats[i] == nil || operationStats[i].EpochInvalidRequest >= int64(model.DemotionCountBeforeSlashing) {
				scores[i] = big.NewFloat(0)

				return nil
			}

			distributionScore := new(big.Float).
				Sub(
					calculateScore(statsData[i].validCount, maxValues.validCount, rewards.OperationScore.Distribution.Weight, 1),
					calculateScore(statsData[i].invalidCount, maxValues.invalidCount, rewards.OperationScore.Distribution.Weight, rewards.OperationScore.Distribution.WeightInvalid),
				)

			dataScore := new(big.Float).Add(calculateScore(statsData[i].networkCount, maxValues.networkCount, rewards.OperationScore.Data.Weight, rewards.OperationScore.Data.WeightNetwork),
				new(big.Float).Add(
					calculateScore(statsData[i].indexerCount, maxValues.indexerCount, rewards.OperationScore.Data.Weight, rewards.OperationScore.Data.WeightIndexer),
					calculateScore(statsData[i].activityCount, maxValues.activityCount, rewards.OperationScore.Data.Weight, rewards.OperationScore.Data.WeightActivity),
				),
			)

			stabilityScore := new(big.Float).
				Add(
					calculateScore(statsData[i].upTime, maxValues.upTime, rewards.OperationScore.Stability.Weight, rewards.OperationScore.Stability.WeightUptime),
					calculateScore(big.NewFloat(float64(lo.Ternary(statsData[i].isLatestVersion, 1, 0))), big.NewFloat(1), rewards.OperationScore.Stability.Weight, rewards.OperationScore.Stability.WeightVersion),
				)

			scores[i] = new(big.Float).Add(distributionScore, new(big.Float).Add(dataScore, stabilityScore))

			mu.Lock()
			// If the score is less than 0, set it to 0
			if scores[i].Cmp(big.NewFloat(0)) < 0 {
				scores[i].Set(big.NewFloat(0))
			}

			totalScore = totalScore.Add(totalScore, scores[i])
			mu.Unlock()

			return nil
		})
	}

	_ = errorPool.Wait()

	return scores, totalScore
}

// maxFloat returns the maximum of two big.Float values.
func maxFloat(a, b *big.Float) *big.Float {
	if a.Cmp(b) > 0 {
		return a
	}

	return b
}

// calculateScore calculates the score for the operation rewards calculation.
func calculateScore(value, maxValue *big.Float, weight, factor float64) *big.Float {
	if maxValue.Cmp(big.NewFloat(0)) == 0 {
		return big.NewFloat(0)
	}

	radio := new(big.Float).Quo(value, maxValue)

	// weight * radio * factor
	return new(big.Float).Mul(big.NewFloat(weight), new(big.Float).Mul(radio, big.NewFloat(factor)))
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
func (s *Server) prepareRequestCounts(ctx context.Context, nodeAddresses []common.Address, nodes []*schema.Node) ([]*big.Int, []*schema.Stat, error) {
	if len(nodeAddresses) == 0 {
		return make([]*big.Int, 0), make([]*schema.Stat, 0), nil
	}

	stats, err := s.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		Addresses: nodeAddresses,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find node stats: %w", err)
	}

	statsMap := make(map[common.Address]*schema.Stat, len(stats))
	for _, stat := range stats {
		statsMap[stat.Address] = stat
	}

	requestCounts := make([]*big.Int, len(nodes))
	operationStats := make([]*schema.Stat, len(nodes))

	for i, nodeAddress := range nodeAddresses {
		if stat, ok := statsMap[nodeAddress]; ok {
			// set request counts for nodes from the epoch.
			requestCounts[i] = big.NewInt(stat.EpochRequest)
			stat.Version = nodes[i].Version
			operationStats[i] = stat
		} else {
			requestCounts[i] = big.NewInt(0)
			operationStats[i] = nil
		}
	}

	return requestCounts, operationStats, nil
}

type ActivityCountResponse struct {
	Count      int64     `json:"count"`
	LastUpdate time.Time `json:"last_update"`
}

// getNodeActivityCount retrieves the s for the node.
func (s *Server) getNodeActivityCount(ctx context.Context, versionStr, endpoint, accessToken string) (*ActivityCountResponse, error) {
	curVersion, _ := version.NewVersion(versionStr)

	var prefix string
	if minVersion, _ := version.NewVersion("1.1.2"); curVersion.GreaterThanOrEqual(minVersion) {
		prefix = "operators/"
	}

	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	fullURL := endpoint + prefix + "activity_count"

	body, _, err := s.httpClient.FetchWithMethod(ctx, http.MethodGet, fullURL, accessToken, nil)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	response := &ActivityCountResponse{}

	if err = json.Unmarshal(data, response); err != nil {
		return nil, err
	}

	return response, nil
}

// getNodeLatestVersion retrieves the latest node version from the network params contract
func (s *Server) getNodeLatestVersion() (string, error) {
	params, err := s.networkParamsContract.GetParams(&bind.CallOpts{}, math.MaxUint64)

	if err != nil {
		return "", fmt.Errorf("failed to get params for lastest epoch %w", err)
	}

	var networkParam struct {
		LatestNodeVersion string `json:"latest_node_version"`
	}

	if err = json.Unmarshal([]byte(params), &networkParam); err != nil {
		return "", fmt.Errorf("failed to unmarshal network params %w", err)
	}

	return networkParam.LatestNodeVersion, nil
}
