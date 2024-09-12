package multicall3

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
)

// Multicall https://github.com/mds1/multicall
// https://etherscan.io/address/0xcA11bde05977b3631167028862bE2a173976CA11
//go:generate go run -mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/Multicall3.abi --pkg multicall3 --type Multicall3 --out contract_multicall3.go

var (
	AddressMulticall3         = common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")
	ChainIDRSS3Mainnet uint64 = 12553
	ChainIDRSS3Testnet uint64 = 2331
)

func Aggregate3(ctx context.Context, chainID uint64, calls []Multicall3Call3, blockNumber *big.Int, contractBackend bind.ContractCaller) ([]*Multicall3Result, error) {
	if !lo.Contains([]uint64{ChainIDRSS3Mainnet, ChainIDRSS3Testnet}, chainID) {
		return nil, fmt.Errorf("unsupported chain id: %d", chainID)
	}

	errorPool := pool.New().WithContext(ctx).WithCancelOnError()

	results := make([]*Multicall3Result, len(calls))

	for index, call := range calls {
		index := index
		call := call

		errorPool.Go(func(ctx context.Context) error {
			message := ethereum.CallMsg{
				To:   lo.ToPtr(call.Target),
				Data: call.CallData,
			}

			data, err := contractBackend.CallContract(ctx, message, blockNumber)
			if err != nil && !call.AllowFailure {
				return err
			}

			result := Multicall3Result{
				Success:    err == nil,
				ReturnData: data,
			}

			results[index] = &result

			return nil
		})
	}

	if err := errorPool.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}
