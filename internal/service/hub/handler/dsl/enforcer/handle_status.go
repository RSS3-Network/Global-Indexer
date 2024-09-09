package enforcer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/go-version"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/node/schema/worker"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (e *SimpleEnforcer) maintainNodeStatus(ctx context.Context) error {
	stats, err := e.getAllNodeStats(ctx, &schema.StatQuery{
		Limit: lo.ToPtr(defaultLimit),
	})
	if err != nil {
		return err
	}

	nodeInfo, err := e.stakingContract.GetNodes(&bind.CallOpts{}, lo.Map(stats, func(stat *schema.Stat, _ int) common.Address {
		return stat.Address
	}))
	if err != nil {
		return fmt.Errorf("get Nodes from chain: %w", err)
	}

	minVersionStr, err := e.getNodeMinVersion()
	if err != nil {
		zap.L().Error("get node min version", zap.Error(err))

		return err
	}

	minVersion, _ := version.NewVersion(minVersionStr)

	var updatedStats []*schema.Stat

	for i := range stats {
		switch nodeInfo[i].Status {
		case uint8(schema.NodeStatusRegistered), uint8(schema.NodeStatusOutdated):
			workersInfo, _ := e.getNodeWorkerStatus(ctx, stats[i].Endpoint, stats[i].AccessToken)

			if workersInfo != nil {
				for _, workerInfo := range workersInfo.Data.Decentralized {
					if workerInfo.Status != worker.StatusReady {
						stats[i].Status = schema.NodeStatusInitializing

						updatedStats = append(updatedStats, stats[i])

						break
					}
				}
			}
		case uint8(schema.NodeStatusInitializing):
			nodeIfo, _ := e.getNodeInfo(ctx, stats[i].Endpoint, stats[i].AccessToken)
			if nodeIfo != nil {
				currentVersion, _ := version.NewVersion(nodeIfo.Data.Version.Tag)

				if currentVersion.LessThan(minVersion) {
					stats[i].Status = schema.NodeStatusOutdated

					updatedStats = append(updatedStats, stats[i])
				}
			}
		case uint8(schema.NodeStatusOnline), uint8(schema.NodeStatusExiting):
			if err = e.getNodeHealth(ctx, stats[i].Endpoint, stats[i].Address); err != nil {
				stats[i].Status = schema.NodeStatusOffline
				// TODO: add offline to invalid response
				updatedStats = append(updatedStats, stats[i])
			}
		}
	}

	return e.updateNodeStatusToVSL(ctx, updatedStats)
}

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

// getNodeHealth retrieves node health.
func (e *SimpleEnforcer) getNodeHealth(ctx context.Context, endpoint string, address common.Address) error {
	response, err := e.httpClient.FetchWithMethod(ctx, http.MethodGet, endpoint, "", nil)
	if err != nil {
		return fmt.Errorf("fetch node endpoint %s: %w", endpoint, err)
	}

	defer lo.Try(response.Close)

	// Use a limited reader to avoid reading too much data.
	content, err := io.ReadAll(io.LimitReader(response, 4096))
	if err != nil {
		return fmt.Errorf("parse node response: %w", err)
	}

	// Check if the node's address is in the response.
	// This is a simple check to ensure the node is responding correctly.
	// The content sample is: "This is an RSS3 Node operated by 0x0000000000000000000000000000000000000000.".
	if !strings.Contains(string(content), address.String()) {
		return fmt.Errorf("invalid node response")
	}

	return nil
}
