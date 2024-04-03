package settler

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

var (
	// epsilon is the acceptable error margin for floating point comparisons
	epsilon        = big.NewFloat(0.00000000001)
	epsilonInt     = big.NewInt(500000)
	specialRewards = config.SpecialRewards{
		GiniCoefficient: 0.00005,
		CliffFactor:     300,
		CliffPoint:      "5000000",
		EpochLimit:      5,
		StakerFactor:    25,
		Rewards:         12328,
	}
)

func TestApplyGiniCoefficient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		poolSize      *big.Float
		expectedScore *big.Float
	}{
		{
			name:          "poolSize 50",
			poolSize:      big.NewFloat(50),
			expectedScore: big.NewFloat(0.9975062344139651),
		},
		{
			name:          "poolSize 500",
			poolSize:      big.NewFloat(500),
			expectedScore: big.NewFloat(0.9756097560975611),
		},
		{
			name:          "poolSize 133667",
			poolSize:      big.NewFloat(133667),
			expectedScore: big.NewFloat(0.13015156149335902),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score := applyGiniCoefficient(tt.poolSize, specialRewards.GiniCoefficient)

			fmt.Println(score, tt.expectedScore)

			diff := new(big.Float).Sub(score, tt.expectedScore)

			if diff.Abs(diff).Cmp(epsilon) > 0 {
				t.Errorf("Score got = %v, want %v, difference %v exceeds epsilon %v", score, tt.expectedScore, diff, epsilon)
			}
		})
	}
}

func TestApplyCliffFactor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		poolSize      *big.Float
		maxPoolSize   *big.Float
		expectedScore *big.Float
	}{
		{
			name:          "poolSize 5000, maxPoolSize 50000",
			poolSize:      big.NewFloat(5000),
			maxPoolSize:   big.NewFloat(50000),
			expectedScore: big.NewFloat(0.8),
		},
		{
			name:          "poolSize 50000, maxPoolSize 50000",
			poolSize:      big.NewFloat(50000),
			maxPoolSize:   big.NewFloat(50000),
			expectedScore: big.NewFloat(0.2857142857142857),
		},
		{
			name:          "poolSize 13756, maxPoolSize 50000",
			poolSize:      big.NewFloat(13756),
			maxPoolSize:   big.NewFloat(50000),
			expectedScore: big.NewFloat(0.5924872615238772),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score := applyGiniCoefficient(tt.poolSize, specialRewards.GiniCoefficient)

			fmt.Println(score)

			cliffPoint := big.NewFloat(0)
			cliffPoint.SetString(specialRewards.CliffPoint)

			if tt.poolSize.Cmp(cliffPoint) == 1 {
				applyCliffFactor(tt.poolSize, tt.maxPoolSize, score, specialRewards.CliffFactor)
			}

			fmt.Println(score, tt.expectedScore)

			diff := new(big.Float).Sub(score, tt.expectedScore)

			if diff.Abs(diff).Cmp(epsilon) > 0 {
				t.Errorf("Score got = %v, want %v, difference %v exceeds epsilon %v", score, tt.expectedScore, diff, epsilon)
			}
		})
	}
}

func TestApplyStakerFactor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		poolSize              *big.Float
		maxPoolSize           *big.Float
		stakers               uint64
		totalEffectiveStakers uint64
		expectedScore         *big.Float
	}{
		{
			name:                  "5 stakers of 20 effective stakers",
			poolSize:              big.NewFloat(50000),
			maxPoolSize:           big.NewFloat(50000),
			stakers:               5,
			totalEffectiveStakers: 20,
			expectedScore:         big.NewFloat(2.071428571428571),
		},
		{
			name:                  "15 stakers of 20 effective stakers",
			poolSize:              big.NewFloat(500),
			maxPoolSize:           big.NewFloat(50000),
			stakers:               15,
			totalEffectiveStakers: 20,
			expectedScore:         big.NewFloat(19.26829268292683),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score := applyGiniCoefficient(tt.poolSize, specialRewards.GiniCoefficient)

			cliffPoint := big.NewFloat(0)
			cliffPoint.SetString(specialRewards.CliffPoint)

			if tt.poolSize.Cmp(cliffPoint) == 1 {
				applyCliffFactor(tt.poolSize, tt.maxPoolSize, score, specialRewards.CliffFactor)
			}

			applyStakerFactor(tt.stakers, tt.totalEffectiveStakers, specialRewards.StakerFactor, score)

			fmt.Println(score, tt.expectedScore)

			diff := new(big.Float).Sub(score, tt.expectedScore)

			if diff.Abs(diff).Cmp(epsilon) > 0 {
				t.Errorf("Score got = %v, want %v, difference %v exceeds epsilon %v", score, tt.expectedScore, diff, epsilon)
			}
		})
	}
}

func TestCalculateOperationRewards(t *testing.T) {
	t.Parallel()

	//correctRewards := []string{
	//	"908820928063146295296",
	//	"1293616057212334505984",
	//	"1636483513676843974656",
	//	"3010321247805582082048",
	//	"2468583836050489606144",
	//	"3010174417191603011584",
	//}

	correctRewards := []string{
		"908000000000000000000",
		"1293000000000000000000",
		"1636000000000000000000",
		"3010000000000000000000",
		"2468000000000000000000",
		"3010000000000000000000",
	}

	// Slice to hold pointers to big.Int
	expectedRewards := make([]*big.Int, len(correctRewards))

	// Convert strings to *big.Int
	for i, numStr := range correctRewards {
		expectedRewards[i] = new(big.Int)
		expectedRewards[i], _ = expectedRewards[i].SetString(numStr, 10) // Base 10 for decimal
	}

	tests := []struct {
		name            string
		nodes           []*schema.Node
		recentStackers  map[common.Address]uint64
		expectedRewards []*big.Int
	}{
		// Mock Nodes
		// [Pool size, Recent Stackers]
		//[7209,3],
		//[3200,4],
		//[1568,5],
		//[501,10],
		//[5000,10],
		//[502,10],
		{
			name: "case 1",
			nodes: []*schema.Node{
				{
					Address:           common.Address{1},
					StakingPoolTokens: "7209",
				},
				{
					Address:           common.Address{2},
					StakingPoolTokens: "3200",
				},
				{
					Address:           common.Address{3},
					StakingPoolTokens: "1568",
				},
				{
					Address:           common.Address{4},
					StakingPoolTokens: "501",
				},
				{
					Address:           common.Address{5},
					StakingPoolTokens: "5000",
				},
				{
					Address:           common.Address{6},
					StakingPoolTokens: "502",
				},
			},
			recentStackers: map[common.Address]uint64{
				common.Address{1}: 3,
				common.Address{2}: 4,
				common.Address{3}: 5,
				common.Address{4}: 10,
				common.Address{5}: 10,
				common.Address{6}: 10,
			},
			expectedRewards: expectedRewards,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rewards, err := calculateAlphaSpecialRewards(tt.nodes, tt.recentStackers, &specialRewards)

			if err != nil {
				t.Error(err)
			}

			if rewards != nil {
				totalRewards := big.NewInt(0)

				for i, reward := range rewards {
					totalRewards.Add(totalRewards, reward)
					diff := new(big.Int).Sub(reward, tt.expectedRewards[i])

					if diff.Abs(diff).Cmp(epsilonInt) > 0 {
						t.Errorf("Reward got = %v, want %v , diff is %v > %v", reward, tt.expectedRewards[i], diff, epsilonInt)
					}
				}

				// Convert specialRewards.Rewards to a *big.Int with 18 decimal places
				specialRewardsBigInt := new(big.Int).SetUint64(specialRewards.Rewards)
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
