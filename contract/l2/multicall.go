package l2

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	v2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/contract/multicall3"
	"go.uber.org/zap"
)

type StakingV2MulticallClient struct {
	*v2.Staking
	chainID        uint64
	ethereumClient bind.ContractCaller
}

type ChipInfo struct {
	NodeAddr common.Address
	Tokens   *big.Int
	Shares   *big.Int
}

func (client *StakingV2MulticallClient) StakingV2GetChipsInfo(ctx context.Context, blockNumber *big.Int, chipIDs []*big.Int) ([]ChipInfo, error) {
	abi, err := v2.StakingMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("get staking contract abi: %w", err)
	}

	// Prepare the input data for the multicall
	calls := make([]multicall3.Multicall3Call3, 0, len(chipIDs))

	for _, chipID := range chipIDs {
		callData, err := abi.Pack("getChipInfo", chipID)
		if err != nil {
			return nil, fmt.Errorf("pack getChipInfo: %w", err)
		}

		calls = append(calls, multicall3.Multicall3Call3{
			Target:   ContractMap[client.chainID].AddressStakingProxy,
			CallData: callData,
		})
	}

	// Execute the multicall
	results, err := multicall3.Aggregate3(ctx, client.chainID, calls, blockNumber, client.ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("multicall failed: %w", err)
	}

	chipsInfo := make([]ChipInfo, len(chipIDs))

	// Process the response data
	for i, call := range results {
		if !call.Success {
			zap.L().Error("multicall failed", zap.String("chipID", chipIDs[i].String()), zap.Any("call", call))

			return nil, fmt.Errorf("multicall failed, chip id: %s", chipIDs[i].String())
		}

		var chipInfo ChipInfo

		// Unpack the returned data
		err := abi.UnpackIntoInterface(&chipInfo, "getChipInfo", call.ReturnData)
		if err != nil {
			zap.L().Error("unpack getChipInfo result", zap.Error(err), zap.String("chipID", chipIDs[i].String()), zap.Any("data", call.ReturnData))

			return nil, fmt.Errorf("unpack getChipInfo: %w", err)
		}

		chipsInfo[i] = chipInfo
	}

	return chipsInfo, nil
}

func NewStakingV2MulticallClient(chainID uint64, ethereumClient *ethclient.Client) (*StakingV2MulticallClient, error) {
	contractAddresses := ContractMap[chainID]
	if contractAddresses == nil {
		return nil, fmt.Errorf("contract address not found for chain id: %d", chainID)
	}

	staking, err := v2.NewStaking(contractAddresses.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("create staking contract: %w", err)
	}

	return &StakingV2MulticallClient{
		Staking:        staking,
		chainID:        chainID,
		ethereumClient: ethereumClient,
	}, nil
}
