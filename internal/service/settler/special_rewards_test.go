package settler

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/shopspring/decimal"
)

var (
	specialRewards = config.SpecialRewards{
		GiniCoefficient: 2,
		StakerFactor:    0.05,
		NodeThreshold:   0.4,
		EpochLimit:      5,
		Rewards:         12328,
		RewardsCeiling:  1000,
	}
)

func TestCalculateOperationRewards(t *testing.T) {
	t.Parallel()

	correctRewards := [][]string{
		{
			"1000000000000000000000",
			"1000000000000000000000",
			"1000000000000000000000",
			"1000000000000000000000",
			"794000000000000000000",
			"1000000000000000000000",
			"1000000000000000000000",
			"1000000000000000000000",
			"1000000000000000000000",
			"711000000000000000000",
		},
		{
			"585000000000000000000",
			"1000000000000000000000",
			"245000000000000000000",
			"0",
			"0",
			"0",
			"0",
			"0",
			"0",
			"0",
		},
	}

	// Slice to hold pointers to big.Int
	expectedRewards := make([][]*big.Int, len(correctRewards))

	// Convert strings to *big.Int
	for i, numStr := range correctRewards {
		expectedRewards[i] = make([]*big.Int, len(numStr))
		for j, str := range numStr {
			expectedRewards[i][j] = new(big.Int)
			expectedRewards[i][j], _ = expectedRewards[i][j].SetString(str, 10) // Base 10 for decimal
		}
	}

	tests := []struct {
		name            string
		nodes           []*schema.Node
		recentStackers  map[common.Address]*schema.StakeRecentCount
		expectedRewards []*big.Int
	}{
		// Mock Nodes
		// [Pool size, Recent Stackers]
		{
			name: "case 1: reach the threshold",
			nodes: []*schema.Node{
				{
					Address:           common.Address{1},
					StakingPoolTokens: "254939021336715733204793",
				},
				{
					Address:           common.Address{2},
					StakingPoolTokens: "4504650447234721822705",
				},
				{
					Address:           common.Address{3},
					StakingPoolTokens: "2830103823431924402058258",
				},
				{
					Address:           common.Address{4},
					StakingPoolTokens: "1333245734400959927416363",
				},
				{
					Address:           common.Address{5},
					StakingPoolTokens: "8172497478991576157429545",
				},
				{
					Address:           common.Address{6},
					StakingPoolTokens: "3007716787095077937681957",
				},
				{
					Address:           common.Address{7},
					StakingPoolTokens: "1474974191505938360613272",
				},
				{
					Address:           common.Address{8},
					StakingPoolTokens: "262882511174109870226533",
				},
				{
					Address:           common.Address{9},
					StakingPoolTokens: "560896507154228577606539",
				},
				{
					Address:           common.Address{10},
					StakingPoolTokens: "10776091947611685629896941",
				},
			},
			recentStackers: map[common.Address]*schema.StakeRecentCount{
				common.Address{1}:  {StakerCount: 32, StakeValue: decimal.NewFromInt(600)},
				common.Address{2}:  {StakerCount: 32, StakeValue: decimal.NewFromInt(100000)},
				common.Address{3}:  {StakerCount: 6, StakeValue: decimal.NewFromInt(10000)},
				common.Address{4}:  {StakerCount: 1, StakeValue: decimal.NewFromInt(500)},
				common.Address{5}:  {StakerCount: 3, StakeValue: decimal.NewFromInt(1500)},
				common.Address{6}:  {StakerCount: 4, StakeValue: decimal.NewFromInt(80000)},
				common.Address{7}:  {StakerCount: 2, StakeValue: decimal.NewFromInt(10000)},
				common.Address{8}:  {StakerCount: 1, StakeValue: decimal.NewFromInt(500)},
				common.Address{9}:  {StakerCount: 1, StakeValue: decimal.NewFromInt(500)},
				common.Address{10}: {StakerCount: 1, StakeValue: decimal.NewFromInt(500)},
			},
			expectedRewards: expectedRewards[0],
		},
		{
			name: "case 2: below the threshold",
			nodes: []*schema.Node{
				{
					Address:           common.Address{1},
					StakingPoolTokens: "254939021336715733204793",
				},
				{
					Address:           common.Address{2},
					StakingPoolTokens: "4504650447234721822705",
				},
				{
					Address:           common.Address{3},
					StakingPoolTokens: "2830103823431924402058258",
				},
				{
					Address:           common.Address{4},
					StakingPoolTokens: "1333245734400959927416363",
				},
				{
					Address:           common.Address{5},
					StakingPoolTokens: "8172497478991576157429545",
				},
				{
					Address:           common.Address{6},
					StakingPoolTokens: "3007716787095077937681957",
				},
				{
					Address:           common.Address{7},
					StakingPoolTokens: "1474974191505938360613272",
				},
				{
					Address:           common.Address{8},
					StakingPoolTokens: "262882511174109870226533",
				},
				{
					Address:           common.Address{9},
					StakingPoolTokens: "560896507154228577606539",
				},
				{
					Address:           common.Address{10},
					StakingPoolTokens: "10776091947611685629896941",
				},
			},
			recentStackers: map[common.Address]*schema.StakeRecentCount{
				common.Address{1}: {StakerCount: 32, StakeValue: decimal.NewFromInt(600)},
				common.Address{2}: {StakerCount: 32, StakeValue: decimal.NewFromInt(100000)},
				common.Address{3}: {StakerCount: 6, StakeValue: decimal.NewFromInt(10000)},
			},
			expectedRewards: expectedRewards[1],
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rewards, _, err := calculateAlphaSpecialRewards(tt.nodes, tt.recentStackers, &specialRewards)

			if err != nil {
				t.Error(err)
			}

			if rewards != nil {
				totalRewards := big.NewInt(0)

				for i, reward := range rewards {
					totalRewards.Add(totalRewards, reward)
					diff := new(big.Int).Sub(reward, tt.expectedRewards[i])

					if reward.Cmp(tt.expectedRewards[i]) != 0 {
						t.Errorf("Reward got = %v, want %v, diff %v", reward, tt.expectedRewards[i], diff)
					}
				}

				// Convert specialRewards.Rewards to a *big.Int with 18 decimal places
				specialRewardsBigInt := new(big.Int).SetUint64(uint64(specialRewards.Rewards))
				rewardCeiling := new(big.Int).Mul(specialRewardsBigInt, big.NewInt(1e18))

				// totalRewards must be less than rewardCeiling
				if totalRewards.Cmp(rewardCeiling) >= 0 {
					t.Errorf("Total rewards is over the limit: %v, limit: %v", totalRewards, rewardCeiling)
				}
			} else {
				t.Error("Rewards is nil")
			}
		})
	}
}
