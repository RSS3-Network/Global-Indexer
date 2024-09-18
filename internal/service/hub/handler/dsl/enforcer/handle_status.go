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
	currentEpoch, err := e.getCurrentEpoch(ctx)
	if err != nil {
		return fmt.Errorf("get current epoch: %w", err)
	}

	// Fixme: use pagination to get all nodes
	nodes, err := e.databaseClient.FindNodes(ctx, schema.FindNodesQuery{})
	if err != nil {
		return fmt.Errorf("find nodes: %w", err)
	}

	nodeVSLInfo, err := e.stakingContract.GetNodes(&bind.CallOpts{}, lo.Map(nodes, func(node *schema.Node, _ int) common.Address {
		return node.Address
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
		updatedNodes          []*schema.Node
		demotionNodeAddresses []common.Address
		reasons               []string
		reporters             []common.Address
	)

	for i := range nodes {
		// initial alpha version nodes are outdated
		// deprecated if there is no alpha version node
		if nodeVSLInfo[i].Status == uint8(schema.NodeStatusNone) && nodes[i].Version == schema.NodeVersionAlpha.String() {
			nodes[i].Status = schema.NodeStatusOutdated
			updatedNodes = append(updatedNodes, nodes[i])
		}

		if nodes[i].Version != schema.NodeVersionAlpha.String() {
			switch nodeVSLInfo[i].Status {
			case uint8(schema.NodeStatusNone), uint8(schema.NodeStatusRegistered), uint8(schema.NodeStatusOutdated), uint8(schema.NodeStatusInitializing):
				nodes[i].Status = schema.NodeStatus(nodeVSLInfo[i].Status)
				newStatus, errPath := e.determineStatus(ctx, nodes[i], minVersion)

				if nodes[i].Status != newStatus {
					nodes[i].Status = newStatus
					updatedNodes = append(updatedNodes, nodes[i])

					if newStatus == schema.NodeStatusOffline {
						responseValue, _ := json.Marshal(fmt.Sprintf(`{"error_message": "%s"}`, errPath))

						e.saveOfflineStatusToInvalidResponse(ctx, uint64(currentEpoch), nodes[i].Address, fmt.Sprintf("%s/%s", nodes[i].Endpoint, errPath), responseValue)
					}
				}
			case uint8(schema.NodeStatusOnline), uint8(schema.NodeStatusExiting):
				nodeDBInfo, err := e.databaseClient.FindNode(ctx, nodes[i].Address)
				if err != nil {
					zap.L().Error("find node", zap.Error(err), zap.Any("address", nodes[i].Address.String()))

					continue
				}

				if nodeDBInfo.Status == schema.NodeStatusOffline {
					nodes[i].Status = schema.NodeStatusOffline
					// TODO: slashing mechanism temporarily disabled.
					// demotionNodeAddresses = append(demotionNodeAddresses, stats[i].Address)
					// reasons = append(reasons, "offline")
					// reporters = append(reporters, ethereum.AddressGenesis)
					responseValue, _ := json.Marshal(fmt.Sprintf(`{"error_message": "%s"}`, "heartbeat"))
					e.saveOfflineStatusToInvalidResponse(ctx, uint64(currentEpoch), nodes[i].Address, "", responseValue)
					updatedNodes = append(updatedNodes, nodes[i])
				}
			}
		}
	}

	return e.updateNodeStatuses(ctx, updatedNodes, demotionNodeAddresses, reasons, reporters)
}

func (e *SimpleEnforcer) determineStatus(ctx context.Context, node *schema.Node, minVersion *version.Version) (schema.NodeStatus, string) {
	nodeInfo, err := e.getNodeInfo(ctx, node.Endpoint, node.AccessToken)
	if err != nil || nodeInfo == nil {
		zap.L().Error("get node info", zap.Error(err), zap.Any("address", node.Address.String()), zap.Any("endpoint", node.Endpoint), zap.Any("access_token", node.AccessToken))

		return schema.NodeStatusOffline, "info"
	}

	currentVersion, _ := version.NewVersion(nodeInfo.Data.Version.Tag)
	if currentVersion.LessThan(minVersion) {
		return schema.NodeStatusOutdated, ""
	}

	workersInfo, err := e.getNodeWorkerStatus(ctx, node.Endpoint, node.AccessToken)
	if err != nil || workersInfo == nil {
		zap.L().Error("get node worker status", zap.Error(err), zap.Any("address", node.Address.String()), zap.Any("endpoint", node.Endpoint), zap.Any("access_token", node.AccessToken))

		return schema.NodeStatusOffline, "workers_status"
	}

	for _, workerInfo := range workersInfo.Data.Decentralized {
		if workerInfo.Status != worker.StatusUnhealthy {
			return schema.NodeStatusInitializing, ""
		}
	}

	return schema.NodeStatusRegistered, ""
}

func (e *SimpleEnforcer) updateNodeStatuses(ctx context.Context, updatedNodes []*schema.Node, demotionNodeAddresses []common.Address, reasons []string, reporters []common.Address) error {
	nodeAddresses, nodeStatusList := make([]common.Address, 0, len(updatedNodes)), make([]uint8, 0, len(updatedNodes))
	for _, node := range updatedNodes {
		nodeAddresses = append(nodeAddresses, node.Address)
		nodeStatusList = append(nodeStatusList, uint8(node.Status))
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
