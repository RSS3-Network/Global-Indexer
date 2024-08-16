package settler

import (
	"math/big"
	"testing"
)

func TestCalculateFinalRewards(t *testing.T) {
	t.Parallel()

	correctRewards := []string{
		"1761000000000000000000",
		"3522000000000000000000",
		"5283000000000000000000",
		"1761000000000000000000",
	}

	expectedRewards := make([]*big.Int, len(correctRewards))

	// Convert strings to *big.Int
	for i, numStr := range correctRewards {
		expectedRewards[i] = new(big.Int)
		expectedRewards[i], _ = expectedRewards[i].SetString(numStr, 10) // Base 10 for decimal
	}

	tests := []struct {
		name            string
		requestCount    []*big.Int
		totalRewards    float64
		expectedRewards []*big.Int
	}{
		{
			name:            "no request counts",
			requestCount:    []*big.Int{},
			totalRewards:    12328,
			expectedRewards: []*big.Int{},
		},
		{
			name:            "zero request counts",
			requestCount:    []*big.Int{big.NewInt(0), big.NewInt(0)},
			totalRewards:    12328,
			expectedRewards: []*big.Int{big.NewInt(0), big.NewInt(0)},
		},
		{
			name:            "valid request counts",
			requestCount:    []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(1), big.NewInt(3)},
			totalRewards:    12328,
			expectedRewards: expectedRewards,
		},
		{
			name:            "exceeds ceiling request counts",
			requestCount:    []*big.Int{big.NewInt(1), big.NewInt(3)},
			totalRewards:    1,
			expectedRewards: []*big.Int{big.NewInt(0), big.NewInt(0)},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rewards, err := calculateFinalRewards(tt.requestCount, tt.totalRewards)

			if err != nil {
				t.Error(err)
			}

			for i := range rewards {
				if rewards[i].Cmp(tt.expectedRewards[i]) != 0 {
					t.Errorf("got = %v, want %v ", rewards[i], tt.expectedRewards[i])
				}
			}
		})
	}
}
