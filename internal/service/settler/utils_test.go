package settler

import (
	"fmt"
	"math/big"
	"testing"
)

func TestScaleGwei(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       *big.Int
		expected *big.Int
	}{
		{
			name:     "1",
			in:       big.NewInt(1),
			expected: big.NewInt(100000000000000000),
		},
		{
			name:     "50",
			in:       big.NewInt(50),
			expected: big.NewInt(5000000000000000000),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scaleGwei(tt.in)

			fmt.Println(tt.in, tt.expected)

			if tt.in.Cmp(tt.expected) != 1 {
				t.Errorf("got = %v, want %v ", tt.in, tt.expected)
			}
		})
	}
}
