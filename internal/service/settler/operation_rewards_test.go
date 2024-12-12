package settler

import (
	"context"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/schema"
)

func TestCalculateScores(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex

	operationStats := []*schema.Stat{
		{Address: common.HexToAddress("0x0")},
		nil,
		{Address: common.HexToAddress("0x1")},
		{Address: common.HexToAddress("0x2")},
		{Address: common.HexToAddress("0x3")},
		{Address: common.HexToAddress("0x4")},
		{Address: common.HexToAddress("0x5")},
	}

	statsData := []StatValue{
		{
			validCount:      big.NewFloat(1300),
			invalidCount:    big.NewFloat(2),
			networkCount:    big.NewFloat(16),
			indexerCount:    big.NewFloat(75),
			activityCount:   big.NewFloat(1913144890),
			upTime:          big.NewFloat(3),
			isLatestVersion: true,
		},
		{},
		{
			validCount:      big.NewFloat(1272),
			invalidCount:    big.NewFloat(2),
			networkCount:    big.NewFloat(16),
			indexerCount:    big.NewFloat(75),
			activityCount:   big.NewFloat(1913144890),
			upTime:          big.NewFloat(3),
			isLatestVersion: true,
		},
		{
			validCount:      big.NewFloat(1206),
			invalidCount:    big.NewFloat(2),
			networkCount:    big.NewFloat(16),
			indexerCount:    big.NewFloat(75),
			activityCount:   big.NewFloat(1913144890),
			upTime:          big.NewFloat(3),
			isLatestVersion: true,
		},
		{
			validCount:      big.NewFloat(340),
			invalidCount:    big.NewFloat(3),
			networkCount:    big.NewFloat(14),
			indexerCount:    big.NewFloat(73),
			activityCount:   big.NewFloat(581482865),
			upTime:          big.NewFloat(3),
			isLatestVersion: true,
		},
		{
			validCount:      big.NewFloat(0),
			invalidCount:    big.NewFloat(0),
			networkCount:    big.NewFloat(1),
			indexerCount:    big.NewFloat(1),
			activityCount:   big.NewFloat(6901733),
			upTime:          big.NewFloat(21),
			isLatestVersion: false,
		},
		{
			validCount:      big.NewFloat(0),
			invalidCount:    big.NewFloat(0),
			networkCount:    big.NewFloat(1),
			indexerCount:    big.NewFloat(1),
			activityCount:   big.NewFloat(6901733),
			upTime:          big.NewFloat(21),
			isLatestVersion: false,
		},
	}

	maxValue := StatValue{
		validCount:    big.NewFloat(1300),
		invalidCount:  big.NewFloat(3),
		networkCount:  big.NewFloat(16),
		indexerCount:  big.NewFloat(75),
		activityCount: big.NewFloat(1913144890),
		upTime:        big.NewFloat(21),
	}

	rewards := &config.Rewards{
		OperationScore: &config.OperationScore{
			Distribution: &config.Distribution{
				Weight:        0.6,
				WeightInvalid: 0.5,
			},
			Data: &config.Data{
				Weight:         0.3,
				WeightNetwork:  0.3,
				WeightIndexer:  0.6,
				WeightActivity: 0.1,
			},
			Stability: &config.Stability{
				Weight:        0.1,
				WeightUptime:  0.7,
				WeightVersion: 0.3,
			},
		},
	}

	expectedScores := []*big.Float{
		big.NewFloat(0.74),
		big.NewFloat(0),
		big.NewFloat(0.7270769230769231),
		big.NewFloat(0.6966153846153846),
		big.NewFloat(0.1599913021245147),
		big.NewFloat(0.0781332259849122),
		big.NewFloat(0.0781332259849122),
	}

	tests := []struct {
		name           string
		operationStats []*schema.Stat
		statsData      []StatValue
		maxValue       StatValue
		rewards        *config.Rewards
		expectedScores []*big.Float
		totalScores    *big.Float
	}{
		{
			name:           "Test Calculate Scores",
			operationStats: operationStats,
			statsData:      statsData,
			maxValue:       maxValue,
			rewards:        rewards,
			expectedScores: expectedScores,
			totalScores:    big.NewFloat(2.4799500617866466),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scores, totalScore := calculateScores(context.Background(), tt.operationStats, tt.statsData, tt.maxValue, tt.rewards, &mu)

			for i := range scores {
				if scores[i].Cmp(tt.expectedScores[i]) != 0 {
					t.Errorf("[%d] got = %v, want %v ", i, scores[i], tt.expectedScores[i])
				}
			}

			if totalScore.Text('f', 3) != tt.totalScores.Text('f', 3) {
				t.Errorf("totalScore got = %v, want %v ", totalScore.Text('f', 3), tt.totalScores.Text('f', 3))
			}
		})
	}
}
