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

	nodeVSLInfo, err := e.stakingContract.GetNodes(&bind.CallOpts{}, lo.Map(stats, func(stat *schema.Stat, _ int) common.Address {
		return stat.Address
	}))
	if err != nil {
		return fmt.Errorf("get Nodes from chain: %w", err)
	}

	minVersionStr, err := e.getNodeMinVersion()
	if err != nil {
		return fmt.Errorf("get node min version: %w", err)
	}

	minVersion, _ := version.NewVersion(minVersionStr)

	var (
		updatedStats          []*schema.Stat
		demotionNodeAddresses []common.Address
		reasons               []string
		reporters             []common.Address
	)

	for i := range stats {
		switch nodeVSLInfo[i].Status {
		case uint8(schema.NodeStatusRegistered):
			nodeInfo, err := e.getNodeInfo(ctx, stats[i].Endpoint, stats[i].AccessToken)
			if err != nil || nodeInfo == nil {
				zap.L().Error("get node info", zap.Error(err), zap.Any("address", stats[i].Address.String()), zap.Any("endpoint", stats[i].Endpoint), zap.Any("access_token", stats[i].AccessToken))

				continue
			}

			currentVersion, _ := version.NewVersion(nodeInfo.Data.Version.Tag)
			if currentVersion.LessThan(minVersion) {
				stats[i].Status = schema.NodeStatusOutdated

				updatedStats = append(updatedStats, stats[i])

				continue
			}

			workersInfo, err := e.getNodeWorkerStatus(ctx, stats[i].Endpoint, stats[i].AccessToken)
			if err != nil || workersInfo == nil {
				zap.L().Error("get node worker status", zap.Error(err), zap.Any("address", stats[i].Address.String()), zap.Any("endpoint", stats[i].Endpoint), zap.Any("access_token", stats[i].AccessToken))

				continue
			}

			for _, workerInfo := range workersInfo.Data.Decentralized {
				if workerInfo.Status != worker.StatusUnhealthy {
					stats[i].Status = schema.NodeStatusInitializing

					updatedStats = append(updatedStats, stats[i])

					break
				}
			}
		case uint8(schema.NodeStatusOutdated):
			nodeInfo, err := e.getNodeInfo(ctx, stats[i].Endpoint, stats[i].AccessToken)
			if err != nil || nodeInfo == nil {
				zap.L().Error("get node info", zap.Error(err), zap.Any("address", stats[i].Address.String()), zap.Any("endpoint", stats[i].Endpoint), zap.Any("access_token", stats[i].AccessToken))

				continue
			}

			currentVersion, _ := version.NewVersion(nodeInfo.Data.Version.Tag)
			if currentVersion.LessThan(minVersion) {
				continue
			}

			stats[i].Status = schema.NodeStatusRegistered

			workersInfo, err := e.getNodeWorkerStatus(ctx, stats[i].Endpoint, stats[i].AccessToken)
			if err != nil || workersInfo == nil {
				zap.L().Error("get node worker status", zap.Error(err), zap.Any("address", stats[i].Address.String()), zap.Any("endpoint", stats[i].Endpoint), zap.Any("access_token", stats[i].AccessToken))

				updatedStats = append(updatedStats, stats[i])

				continue
			}

			for _, workerInfo := range workersInfo.Data.Decentralized {
				if workerInfo.Status != worker.StatusUnhealthy {
					stats[i].Status = schema.NodeStatusInitializing

					updatedStats = append(updatedStats, stats[i])

					break
				}
			}
		case uint8(schema.NodeStatusInitializing):
			nodeInfo, err := e.getNodeInfo(ctx, stats[i].Endpoint, stats[i].AccessToken)
			if err != nil || nodeInfo == nil {
				zap.L().Error("get node info", zap.Error(err), zap.Any("address", stats[i].Address.String()), zap.Any("endpoint", stats[i].Endpoint), zap.Any("access_token", stats[i].AccessToken))

				continue
			}

			currentVersion, _ := version.NewVersion(nodeInfo.Data.Version.Tag)
			if currentVersion.LessThan(minVersion) {
				stats[i].Status = schema.NodeStatusOutdated

				updatedStats = append(updatedStats, stats[i])
			}
		case uint8(schema.NodeStatusOnline), uint8(schema.NodeStatusExiting):
			nodeDBInfo, err := e.databaseClient.FindNode(ctx, stats[i].Address)

			if err != nil {
				zap.L().Error("find node", zap.Error(err), zap.Any("address", stats[i].Address.String()))

				continue
			}

			if nodeDBInfo.Status == schema.NodeStatusOffline {
				stats[i].Status = schema.NodeStatusOffline
				// TODO: add offline to invalid response
				updatedStats = append(updatedStats, stats[i])
				// TODO: slashing mechanism temporarily disabled.
				//demotionNodeAddresses = append(demotionNodeAddresses, stats[i].Address)
				//reasons = append(reasons, "offline")
				//reporters = append(reporters, ethereum.AddressGenesis)
			}
		}
	}

	nodeAddresses := make([]common.Address, len(updatedStats))
	nodeStatusList := make([]uint8, len(updatedStats))

	for i := range updatedStats {
		nodeAddresses[i], nodeStatusList[i] = updatedStats[i].Address, uint8(updatedStats[i].Status)
	}

	return e.updateNodeStatusAndSubmitDemotionToVSL(ctx, nodeAddresses, nodeStatusList, demotionNodeAddresses, reasons, reporters)
}

func (e *SimpleEnforcer) updateNodeStatusAndSubmitDemotionToVSL(ctx context.Context, nodeAddresses []common.Address, nodeStatusList []uint8, demotionNodeAddresses []common.Address, reasons []string, reporters []common.Address) error {
	data, err := prepareSetNodeStatusAndSubmitDemotionsData(nodeAddresses, nodeStatusList, demotionNodeAddresses, reasons, reporters)

	if err != nil {
		return err
	}

	if err = e.invokeSettlementMultiContract(ctx, data); err != nil {
		return fmt.Errorf("invoke settlement contract: %w", err)
	}

	return nil
}

func prepareSetNodeStatusAndSubmitDemotionsData(nodeAddresses []common.Address, nodeStatusList []uint8, demotionNodeAddresses []common.Address, reasons []string, reporters []common.Address) ([][]byte, error) {
	data := make([][]byte, 0)

	if len(nodeAddresses) > 0 {
		inputSetNodeStatus, err := txmgr.EncodeInput(l2.SettlementMetaData.ABI, l2.MethodSetNodeStatus, nodeAddresses, nodeStatusList)
		if err != nil {
			return nil, fmt.Errorf("encode setNodeStatus input: %w", err)
		}

		data = append(data, inputSetNodeStatus)
	}

	if len(demotionNodeAddresses) > 0 {
		inputSubmitDemotions, err := txmgr.EncodeInput(l2.SettlementMetaData.ABI, l2.MethodSubmitDemotions, demotionNodeAddresses, reasons, reporters)
		if err != nil {
			return nil, fmt.Errorf("encode submitDemotions input: %w", err)
		}

		data = append(data, inputSubmitDemotions)
	}

	return data, nil
}

func (e *SimpleEnforcer) invokeSettlementMultiContract(ctx context.Context, data [][]byte) error {
	input, err := txmgr.EncodeInput(l2.SettlementMetaData.ABI, l2.MethodMulticall, data)
	if err != nil {
		return fmt.Errorf("encode input: %w", err)
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
