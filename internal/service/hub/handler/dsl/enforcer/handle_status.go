package enforcer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (e *SimpleEnforcer) updateNodeStatusToVSL(ctx context.Context, nodeAddresses []common.Address, nodeStatusList []schema.NodeStatus) error {
	if len(nodeAddresses) == 0 {
		return fmt.Errorf("no node address provided")
	}

	if len(nodeStatusList) != len(nodeAddresses) {
		return fmt.Errorf("node status list length does not match node address list length")
	}

	if err := e.invokeSettlementContract(ctx, nodeAddresses, nodeStatusList); err != nil {
		return fmt.Errorf("invoke settlement contract: %w", err)
	}

	return nil
}

func (e *SimpleEnforcer) invokeSettlementContract(ctx context.Context, nodeAddresses []common.Address, nodeStatusList []schema.NodeStatus) error {
	input, err := prepareInputData(nodeAddresses, nodeStatusList)
	if err != nil {
		return err
	}

	if err = e.sendTransaction(ctx, input); err != nil {
		return err
	}

	return nil
}

func (e *SimpleEnforcer) sendTransaction(ctx context.Context, input []byte) error {
	txCandidate := txmgr.TxCandidate{
		TxData:   input,
		To:       lo.ToPtr(l2.ContractMap[e.chainID.Uint64()].AddressSettlementProxy),
		GasLimit: e.settlerConfig.GasLimit,
		Value:    big.NewInt(0),
	}

	receipt, err := e.txManager.Send(ctx, txCandidate)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		zap.L().Error("received an invalid transaction receipt", zap.String("tx", receipt.TxHash.String()))

		// Fixme: retry logic and error handling
		select {}
	}

	return nil
}

func prepareInputData(nodeAddresses []common.Address, nodeStatusList []schema.NodeStatus) ([]byte, error) {
	input, err := txmgr.EncodeInput(l2.SettlementMetaData.ABI, l2.MethodDistributeRewards, nodeAddresses, nodeStatusList)
	if err != nil {
		return nil, fmt.Errorf("encode input: %w", err)
	}

	return input, nil
}
