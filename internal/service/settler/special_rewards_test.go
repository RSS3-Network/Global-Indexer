package settler

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

var (
	// epsilon is the acceptable error margin for floating point comparisons
	epsilon        = 0.00000000001
	specialRewards = config.SpecialRewards{
		GiniCoefficient: 0.0003,
		CliffFactor:     0.5,
		CliffPoint:      500,
		EpochLimit:      5,
		StakerFactor:    0.05,
		Rewards:         12328,
	}
)

func TestApplyGiniCoefficient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		poolSize      uint64
		expectedScore float64
	}{
		{
			name:          "poolSize 50",
			poolSize:      50,
			expectedScore: 0.9852216748768474,
		},
		{
			name:          "poolSize 500",
			poolSize:      500,
			expectedScore: 0.8695652173913044,
		},
		{
			name:          "poolSize 133667",
			poolSize:      133667,
			expectedScore: 0.024330841044182375,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score := applyGiniCoefficient(tt.poolSize, specialRewards.GiniCoefficient)

			fmt.Println(score, tt.expectedScore)

			diff := math.Abs(score - tt.expectedScore)

			if diff > epsilon {
				t.Errorf("Score got = %v, want %v, difference %v exceeds epsilon %v", score, tt.expectedScore, diff, epsilon)
			}
		})
	}
}

func TestApplyCliffFactor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		poolSize      uint64
		maxPoolSize   uint64
		expectedScore float64
	}{
		{
			name:          "poolSize 5000, maxPoolSize 50000",
			poolSize:      5000,
			maxPoolSize:   50000,
			expectedScore: 0.37321319661472296,
		},
		{
			name:          "poolSize 50000, maxPoolSize 50000",
			poolSize:      50000,
			maxPoolSize:   50000,
			expectedScore: 0.03125000000000001,
		},
		{
			name:          "poolSize 13756, maxPoolSize 50000",
			poolSize:      13756,
			maxPoolSize:   50000,
			expectedScore: 0.16118857353672705,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score := applyGiniCoefficient(tt.poolSize, specialRewards.GiniCoefficient)

			fmt.Println(score)

			if tt.poolSize > specialRewards.CliffPoint {
				applyCliffFactor(tt.poolSize, tt.maxPoolSize, &score, specialRewards.CliffFactor)
			}

			fmt.Println(score, tt.expectedScore)

			diff := math.Abs(score - tt.expectedScore)

			if diff > epsilon {
				t.Errorf("Score got = %v, want %v, difference %v exceeds epsilon %v", score, tt.expectedScore, diff, epsilon)
			}
		})
	}
}

func TestApplyStakerFactor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		poolSize              uint64
		maxPoolSize           uint64
		stakers               uint64
		totalEffectiveStakers uint64
		expectedScore         float64
	}{
		{
			name:                  "10 stakers of 5 effective stakers",
			poolSize:              50000,
			maxPoolSize:           50000,
			stakers:               10,
			totalEffectiveStakers: 5,
			expectedScore:         0.03437500000000001,
		},
		{
			name:                  "15 stakers of 30 effective stakers",
			poolSize:              500,
			maxPoolSize:           50000,
			stakers:               15,
			totalEffectiveStakers: 30,
			expectedScore:         0.891304347826087,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score := applyGiniCoefficient(tt.poolSize, specialRewards.GiniCoefficient)

			if tt.poolSize > specialRewards.CliffPoint {
				applyCliffFactor(tt.poolSize, tt.maxPoolSize, &score, specialRewards.CliffFactor)
			}

			applyStakerFactor(tt.stakers, tt.totalEffectiveStakers, specialRewards.StakerFactor, &score)

			fmt.Println(score, tt.expectedScore)

			diff := math.Abs(score - tt.expectedScore)

			if diff > epsilon {
				t.Errorf("Score got = %v, want %v, difference %v exceeds epsilon %v", score, tt.expectedScore, diff, epsilon)
			}
		})
	}
}

func TestCalculateOperationRewards(t *testing.T) {
	t.Parallel()

	correctRewards := []string{
		"644000000000000000000",
		"1530000000000000000000",
		"2386000000000000000000",
		"3379000000000000000000",
		"1008000000000000000000",
		"3378000000000000000000",
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
				fmt.Print(rewards)

				for i, reward := range rewards {
					if reward.Cmp(tt.expectedRewards[i]) != 0 {
						t.Errorf("Reward got = %v, want %v", reward, tt.expectedRewards[i])
					}
				}
			} else {
				t.Error("Rewards is nil")
			}
		})
	}
}
