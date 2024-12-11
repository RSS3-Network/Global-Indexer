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
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-version"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
)

func (s *Server) calculateOperationRewards(nodeAddresses []common.Address, operationStats []*schema.Stat, rewards *config.Rewards) ([]*big.Int, error) {
	// If there are no nodes, return nil
	if len(nodeAddresses) == 0 {
		return nil, nil
	}

	operationRewards, err := s.calculateFinalRewards(nodeAddresses, operationStats, rewards)

	if err != nil {
		return nil, fmt.Errorf("failed to calculate operation rewards: %w", err)
	}

	return operationRewards, nil
}

// calculateFinalRewards calculates the final rewards for each node based on the request count and total rewards.
func (s *Server) calculateFinalRewards(nodeAddresses []common.Address, operationStats []*schema.Stat, rewards *config.Rewards) ([]*big.Int, error) {
	latestVersionStr, _ := s.getNodeLatestVersion()
	latestVersion := version.Must(version.NewVersion(latestVersionStr))

	operationRewards := make([]*big.Int, len(nodeAddresses))
	now := time.Now()

	validCounts := make([]*big.Float, len(nodeAddresses))
	invalidCounts := make([]*big.Float, len(nodeAddresses))
	networkCounts := make([]*big.Float, len(nodeAddresses))
	indexerCounts := make([]*big.Float, len(nodeAddresses))
	activityCounts := make([]*big.Float, len(nodeAddresses))
	upTimes := make([]*big.Float, len(nodeAddresses))
	isLatestVersions := make([]bool, len(nodeAddresses))

	validCountMax := big.NewFloat(0)
	invalidCountMax := big.NewFloat(0)
	networkCountMax := big.NewFloat(0)
	indexerCountMax := big.NewFloat(0)
	activityCountMax := big.NewFloat(0)
	upTimeMax := big.NewFloat(0)

	for i := range operationStats {
		if operationStats[i] == nil {
			continue
		}

		validCounts[i] = big.NewFloat(float64(operationStats[i].EpochRequest))
		invalidCounts[i] = big.NewFloat(float64(operationStats[i].EpochInvalidRequest))
		networkCounts[i] = big.NewFloat(float64(operationStats[i].DecentralizedNetwork + operationStats[i].FederatedNetwork))
		indexerCounts[i] = big.NewFloat(float64(operationStats[i].Indexer))

		activityCountResp, err := s.getNodeActivityCount(context.Background(), operationStats[i].Version, operationStats[i].Endpoint, operationStats[i].AccessToken)
		if err != nil {
			activityCounts[i] = big.NewFloat(0)
		}

		activityCounts[i] = big.NewFloat(float64(activityCountResp.Count))

		upTimes[i] = big.NewFloat(now.Sub(operationStats[i].ResetAt).Seconds())
		isLatestVersions[i] = version.Must(version.NewVersion(operationStats[i].Version)).GreaterThanOrEqual(latestVersion)

		if validCounts[i].Cmp(validCountMax) > 0 {
			validCountMax = validCounts[i]
		}

		if invalidCounts[i].Cmp(invalidCountMax) > 0 {
			invalidCountMax = invalidCounts[i]
		}

		if networkCounts[i].Cmp(networkCountMax) > 0 {
			networkCountMax = networkCounts[i]
		}

		if indexerCounts[i].Cmp(indexerCountMax) > 0 {
			indexerCountMax = indexerCounts[i]
		}

		if activityCounts[i].Cmp(activityCountMax) > 0 {
			activityCountMax = activityCounts[i]
		}

		if upTimes[i].Cmp(upTimeMax) > 0 {
			upTimeMax = upTimes[i]
		}
	}

	scores := make([]*big.Float, len(operationStats))
	// Calculate the total score
	totalScore := big.NewFloat(0)

	for i := range operationStats {
		if operationStats[i] == nil {
			scores[i] = big.NewFloat(0)

			continue
		}

		validCountRadio := big.NewFloat(0).Quo(validCounts[i], validCountMax)
		invalidCountRadio := big.NewFloat(0).Quo(invalidCounts[i], invalidCountMax)
		distributionValidScore := big.NewFloat(rewards.OperationScore.Distribution.Weight).Mul(validCountRadio, big.NewFloat(1))
		distributionInValidScore := big.NewFloat(rewards.OperationScore.Distribution.Weight).Mul(invalidCountRadio, big.NewFloat(rewards.OperationScore.Distribution.WeightInvalid))
		distributionScore := big.NewFloat(0).Sub(distributionValidScore, distributionInValidScore)

		networkCountRadio := big.NewFloat(0).Quo(networkCounts[i], networkCountMax)
		indexerCountRadio := big.NewFloat(0).Quo(indexerCounts[i], indexerCountMax)
		activityCountRadio := big.NewFloat(0).Quo(activityCounts[i], activityCountMax)
		networkScore := big.NewFloat(rewards.OperationScore.Data.Weight).Mul(networkCountRadio, big.NewFloat(rewards.OperationScore.Data.WeightNetwork))
		indexerScore := big.NewFloat(rewards.OperationScore.Data.Weight).Mul(indexerCountRadio, big.NewFloat(rewards.OperationScore.Data.WeightIndexer))
		activityScore := big.NewFloat(rewards.OperationScore.Data.Weight).Mul(activityCountRadio, big.NewFloat(rewards.OperationScore.Data.WeightActivity))
		dataScore := networkScore.Add(indexerScore, activityScore)

		upTimeRadio := big.NewFloat(0).Quo(upTimes[i], upTimeMax)
		upTimesScore := big.NewFloat(rewards.OperationScore.Stability.Weight).Mul(upTimeRadio, big.NewFloat(rewards.OperationScore.Stability.WeightUptime))
		versionScore := big.NewFloat(rewards.OperationScore.Stability.Weight).Mul(big.NewFloat(float64(lo.Ternary(isLatestVersions[i], 1, 0))), big.NewFloat(rewards.OperationScore.Stability.WeightVersion))
		stabilityScore := big.NewFloat(0).Add(upTimesScore, versionScore)

		scores[i] = distributionScore.Add(dataScore, stabilityScore)
		totalScore = totalScore.Add(totalScore, scores[i])
	}

	//Calculate the rewards for each node
	for i := range scores {
		if scores[i].Cmp(big.NewFloat(0)) == 0 {
			operationRewards[i] = big.NewInt(0)

			continue
		}

		// Calculate the rewards for the node
		radio := new(big.Float).Quo(scores[i], totalScore)
		reward := new(big.Float).Mul(radio, big.NewFloat(rewards.OperationRewards))

		// Convert to integer to truncate before scaling
		rewardFinal, _ := reward.Int(nil)

		// Apply gwei after truncation
		scaleGwei(rewardFinal)

		operationRewards[i] = rewardFinal
	}

	err := checkRewardsCeiling(operationRewards, rewards.OperationRewards)
	if err != nil {
		return nil, err
	}

	return operationRewards, nil
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

	body, err := s.httpClient.FetchWithMethod(ctx, http.MethodGet, fullURL, accessToken, nil)
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
