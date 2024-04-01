package settler

import (
	"fmt"
	"math/big"
	"testing"
)

// epsilon is the acceptable error margin for floating point comparisons
var epsilon = big.NewFloat(0.00000000001)

// TODO: refactor the tests to use calculateAlphaSpecialRewards() instead
// as it has more control flow and must be tested
func TestApplyGiniCoefficient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		poolSize      *big.Int
		expectedScore *big.Float
	}{
		{
			name:          "poolSize 50",
			poolSize:      big.NewInt(50),
			expectedScore: big.NewFloat(0.9852216748768474),
		},
		{
			name:          "poolSize 500",
			poolSize:      big.NewInt(500),
			expectedScore: big.NewFloat(0.8695652173913044),
		},
		{
			name:          "poolSize 133667",
			poolSize:      big.NewInt(133667),
			expectedScore: big.NewFloat(0.024330841044182375),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score := applyGiniCoefficient(tt.poolSize)

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
		poolSize      *big.Int
		maxPoolSize   *big.Int
		expectedScore *big.Float
	}{
		{
			name:          "poolSize 5000, maxPoolSize 50000",
			poolSize:      big.NewInt(5000),
			maxPoolSize:   big.NewInt(50000),
			expectedScore: big.NewFloat(0.37321319661472296),
		},
		{
			name:          "poolSize 50000, maxPoolSize 50000",
			poolSize:      big.NewInt(50000),
			maxPoolSize:   big.NewInt(50000),
			expectedScore: big.NewFloat(0.03125000000000001),
		},
		{
			name:          "poolSize 13756, maxPoolSize 50000",
			poolSize:      big.NewInt(13756),
			maxPoolSize:   big.NewInt(50000),
			expectedScore: big.NewFloat(0.19505344464383242),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score := applyGiniCoefficient(tt.poolSize)
			cliffCmp := tt.poolSize.Cmp(cliffPoint)

			if cliffCmp == 1 {
				applyCliffFactor(tt.poolSize, tt.maxPoolSize, score)
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
		poolSize              *big.Int
		maxPoolSize           *big.Int
		stakers               uint64
		totalEffectiveStakers uint64
		expectedScore         *big.Float
	}{
		{
			name:                  "10 stakers of 5 effective stakers",
			poolSize:              big.NewInt(50000),
			maxPoolSize:           big.NewInt(50000),
			stakers:               10,
			totalEffectiveStakers: 5,
			expectedScore:         big.NewFloat(0.03437500000000001),
		},
		{
			name:                  "15 stakers of 30 effective stakers",
			poolSize:              big.NewInt(500),
			maxPoolSize:           big.NewInt(50000),
			stakers:               15,
			totalEffectiveStakers: 30,
			expectedScore:         big.NewFloat(0.891304347826087),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			score := applyGiniCoefficient(tt.poolSize)

			cliffCmp := tt.poolSize.Cmp(cliffPoint)
			if cliffCmp == 1 {
				applyCliffFactor(tt.poolSize, tt.maxPoolSize, score)
			}

			applyStakerFactor(tt.stakers, tt.totalEffectiveStakers, score)

			fmt.Println(score, tt.expectedScore)

			diff := new(big.Float).Sub(score, tt.expectedScore)
			if diff.Abs(diff).Cmp(epsilon) > 0 {
				t.Errorf("Score got = %v, want %v, difference %v exceeds epsilon %v", score, tt.expectedScore, diff, epsilon)
			}
		})
	}
}
