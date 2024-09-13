package enforcer

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/go-version"
	"github.com/rss3-network/global-indexer/common/ethereum"
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
		return fmt.Errorf("get nodes from chain: %w", err)
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
		case uint8(schema.NodeStatusRegistered), uint8(schema.NodeStatusOutdated), uint8(schema.NodeStatusInitializing):
			stats[i].Status = schema.NodeStatus(nodeVSLInfo[i].Status)
			newStatus, errPath := e.determineStatus(ctx, stats[i], minVersion)

			if stats[i].Status != newStatus {
				stats[i].Status = newStatus
				updatedStats = append(updatedStats, stats[i])

				if newStatus == schema.NodeStatusOffline {
					responseValue, _ := json.Marshal(fmt.Sprintf(`{"error_message": "%s"}`, errPath))

					e.saveOfflineStatusToInvalidResponse(ctx, uint64(stats[i].Epoch), stats[i].Address, fmt.Sprintf("%s/%s", stats[i].Endpoint, errPath), responseValue)
				}
			}
		case uint8(schema.NodeStatusOnline), uint8(schema.NodeStatusExiting):
			nodeDBInfo, err := e.databaseClient.FindNode(ctx, stats[i].Address)
			if err != nil {
				zap.L().Error("find node", zap.Error(err), zap.Any("address", stats[i].Address.String()))

				continue
			}

			if nodeDBInfo.Status == schema.NodeStatusOffline {
				stats[i].Status = schema.NodeStatusOffline
				// TODO: slashing mechanism temporarily disabled.
				// demotionNodeAddresses = append(demotionNodeAddresses, stats[i].Address)
				// reasons = append(reasons, "offline")
				// reporters = append(reporters, ethereum.AddressGenesis)
				responseValue, _ := json.Marshal(fmt.Sprintf(`{"error_message": "%s"}`, "heartbeat"))
				e.saveOfflineStatusToInvalidResponse(ctx, uint64(stats[i].Epoch), stats[i].Address, "", responseValue)
				updatedStats = append(updatedStats, stats[i])
			}
		}
	}

	return e.updateNodeStatuses(ctx, updatedStats, demotionNodeAddresses, reasons, reporters)
}

func (e *SimpleEnforcer) determineStatus(ctx context.Context, stat *schema.Stat, minVersion *version.Version) (schema.NodeStatus, string) {
	nodeInfo, err := e.getNodeInfo(ctx, stat.Endpoint, stat.AccessToken)
	if err != nil || nodeInfo == nil {
		zap.L().Error("get node info", zap.Error(err), zap.Any("address", stat.Address.String()), zap.Any("endpoint", stat.Endpoint), zap.Any("access_token", stat.AccessToken))

		return schema.NodeStatusOffline, "info"
	}

	currentVersion, _ := version.NewVersion(nodeInfo.Data.Version.Tag)
	if currentVersion.LessThan(minVersion) {
		return schema.NodeStatusOutdated, ""
	}

	workersInfo, err := e.getNodeWorkerStatus(ctx, stat.Endpoint, stat.AccessToken)
	if err != nil || workersInfo == nil {
		zap.L().Error("get node worker status", zap.Error(err), zap.Any("address", stat.Address.String()), zap.Any("endpoint", stat.Endpoint), zap.Any("access_token", stat.AccessToken))

		return schema.NodeStatusOffline, "workers_status"
	}

	for _, workerInfo := range workersInfo.Data.Decentralized {
		if workerInfo.Status != worker.StatusUnhealthy {
			return schema.NodeStatusInitializing, ""
		}
	}

	return schema.NodeStatusRegistered, ""
}

func (e *SimpleEnforcer) updateNodeStatuses(ctx context.Context, updatedStats []*schema.Stat, demotionNodeAddresses []common.Address, reasons []string, reporters []common.Address) error {
	nodeAddresses, nodeStatusList := make([]common.Address, 0, len(updatedStats)), make([]uint8, 0, len(updatedStats))
	for _, stat := range updatedStats {
		nodeAddresses = append(nodeAddresses, stat.Address)
		nodeStatusList = append(nodeStatusList, uint8(stat.Status))
	}

	return e.updateNodeStatusAndSubmitDemotionToVSL(ctx, nodeAddresses, nodeStatusList, demotionNodeAddresses, reasons, reporters)
}

func (e *SimpleEnforcer) saveOfflineStatusToInvalidResponse(ctx context.Context, epochID uint64, nodeAddress common.Address, request string, response json.RawMessage) {
	nodeInvalidResponse := &schema.NodeInvalidResponse{
		EpochID:          epochID,
		Type:             schema.NodeInvalidResponseTypeOffline,
		VerifierNodes:    []common.Address{ethereum.AddressGenesis},
		Request:          request,
		VerifierResponse: json.RawMessage{},
		Node:             nodeAddress,
		Response:         response,
	}

	if err := e.databaseClient.SaveNodeInvalidResponses(ctx, []*schema.NodeInvalidResponse{nodeInvalidResponse}); err != nil {
		zap.L().Error("save node invalid response", zap.Error(err))
	}
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
