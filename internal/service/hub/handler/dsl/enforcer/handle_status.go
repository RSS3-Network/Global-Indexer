package enforcer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (e *SimpleEnforcer) updateNodeStatusToVSL(ctx context.Context, stats []*schema.Stat) error {
	nodeAddresses := make([]common.Address, len(stats))
	nodeStatusList := make([]schema.NodeStatus, len(stats))

	for i := range stats {
		nodeAddresses[i], nodeStatusList[i] = stats[i].Address, stats[i].Status
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

func (e *SimpleEnforcer) getNodeMinVersion() (string, error) {
	params, err := e.networkParamsContract.GetParams(&bind.CallOpts{}, math.MaxUint64)

	if err != nil {
		return "", fmt.Errorf("failed to get params for lastest epoch %w", err)
	}

	var networkParam struct {
		MinNodeVersion string `json:"minimal_node_version"`
	}

	if err = json.Unmarshal([]byte(params), &networkParam); err != nil {
		return "", fmt.Errorf("failed to unmarshal network params %w", err)
	}

	return networkParam.MinNodeVersion, nil
}

// getNodeInfo retrieves node info.
func (e *SimpleEnforcer) getNodeInfo(ctx context.Context, endpoint, accessToken string) (*InfoResponse, error) {
	fullURL := endpoint + "/info"

	body, err := e.httpClient.FetchWithMethod(ctx, http.MethodGet, fullURL, accessToken, nil)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	response := &InfoResponse{}

	if err = json.Unmarshal(data, response); err != nil {
		return nil, err
	}

	return response, nil
}
