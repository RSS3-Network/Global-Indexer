package l2_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/contract/multicall3"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestStakingV2GetChipsInfoByMulticall(t *testing.T) {
	t.Parallel()

	type arguments struct {
		chainID     uint64
		blockNumber *big.Int
		ChipIDs     []*big.Int
	}

	testcases := []struct {
		name      string
		arguments arguments
		want      []l2.ChipInfo
	}{
		{
			name: "Get Mainnet Chips Info",
			arguments: arguments{
				chainID:     multicall3.ChainIDRSS3Mainnet,
				blockNumber: big.NewInt(6023346),
				ChipIDs: []*big.Int{
					big.NewInt(1869),
					big.NewInt(1870),
					big.NewInt(1671),
				},
			},
			want: []l2.ChipInfo{
				{
					NodeAddr: common.HexToAddress("0x39f9e912c1f696f533e7a2267ea233aec9742b35"),
					Tokens:   lo.Must(decimal.NewFromString("677912168357482297897")).BigInt(),
					Shares:   lo.Must(decimal.NewFromString("500000000000000000000")).BigInt(),
				},
				{
					NodeAddr: common.HexToAddress("0x3Ca85BD1eB958C0aBC9D06684C5Ac01f71029DD5"),
					Tokens:   lo.Must(decimal.NewFromString("648814753997254162748")).BigInt(),
					Shares:   lo.Must(decimal.NewFromString("500000000000000000000")).BigInt(),
				},
				{
					NodeAddr: common.HexToAddress("0xc8b960d09c0078c18dcbe7eb9ab9d816bcca8944"),
					Tokens:   lo.Must(decimal.NewFromString("639898391318730922465")).BigInt(),
					Shares:   lo.Must(decimal.NewFromString("500000000000000000000")).BigInt(),
				},
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		ctx := context.Background()

		ethereumClient, err := ethclient.DialContext(ctx, "https://rpc.rss3.io")
		require.NoError(t, err)

		client, err := l2.NewStakingV2MulticallClient(multicall3.ChainIDRSS3Mainnet, ethereumClient)
		require.NoError(t, err)

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			chipsInfo, err := client.StakingV2GetChipsInfoByMulticall(ctx, testcase.arguments.blockNumber, testcase.arguments.ChipIDs)
			require.NoError(t, err)

			require.Equal(t, testcase.want, chipsInfo)
		})
	}
}
