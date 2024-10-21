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

// maintainNodeStatus updates the node statuses based on the node information retrieved from the VSL.
func (e *SimpleEnforcer) maintainNodeStatus(ctx context.Context) error {
	// get the current epoch
	currentEpoch, err := e.getCurrentEpoch(ctx)
	if err != nil {
		return fmt.Errorf("get current epoch: %w", err)
	}

	// Fixme: use pagination to get all nodes
	// or use node stats to get all nodes
	nodes, err := e.databaseClient.FindNodes(ctx, schema.FindNodesQuery{})
	if err != nil {
		return fmt.Errorf("find nodes: %w", err)
	}

	// retrieve the node info from the VSL
	nodeVSLInfo, err := e.stakingContract.GetNodes(&bind.CallOpts{}, lo.Map(nodes, func(node *schema.Node, _ int) common.Address {
		return node.Address
	}))
	if err != nil {
		return fmt.Errorf("get nodes from chain: %w", err)
	}

	// get the min version of the node in rss3 network
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
		switch nodeVSLInfo[i].Status {
		// Handle cases for None, Registered, Outdated, and Initializing statuses
		case uint8(schema.NodeStatusNone),
			uint8(schema.NodeStatusRegistered),
			uint8(schema.NodeStatusOutdated),
			uint8(schema.NodeStatusInitializing),
			uint8(schema.NodeStatusOffline):
			// It indicates that the node has not started.
			// Keep the status as registered.
			if nodes[i].Status == schema.NodeStatusRegistered {
				if nodeVSLInfo[i].Status == uint8(schema.NodeStatusNone) {
					updatedNodes = append(updatedNodes, nodes[i])
				}

				continue
			}

			// Determine new status and potential error path
			newStatus, errPath := e.determineStatus(ctx, nodes[i], minVersion)

			// If status has changed, update and handle accordingly
			if schema.NodeStatus(nodeVSLInfo[i].Status) != newStatus {
				nodes[i].Status = newStatus
				updatedNodes = append(updatedNodes, nodes[i])

				// If new status is offline, save error information
				if newStatus == schema.NodeStatusOffline {
					responseValue, _ := json.Marshal(fmt.Sprintf(`{"error_message": "%s"}`, errPath))

					e.saveOfflineStatusToInvalidResponse(ctx, uint64(currentEpoch), nodes[i].Address, fmt.Sprintf("%s/%s", nodes[i].Endpoint, errPath), responseValue)
				}
			}
		// Handle cases for Online and Exiting statuses
		case uint8(schema.NodeStatusOnline),
			uint8(schema.NodeStatusExiting):
			// If node status from heartbeat is offline, update node status
			if nodes[i].Status == schema.NodeStatusOffline {
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

	return e.updateNodeStatuses(ctx, updatedNodes, demotionNodeAddresses, reasons, reporters)
}

// determineStatus checks the node's status and version to determine its current state
func (e *SimpleEnforcer) determineStatus(ctx context.Context, node *schema.Node, minVersion *version.Version) (schema.NodeStatus, string) {
	// Check if node version meets minimum requirements
	currentVersion, _ := version.NewVersion(node.Version)
	if currentVersion.LessThan(minVersion) {
		// Return outdated status if version is below minimum
		return schema.NodeStatusOutdated, ""
	}

	// Get worker status
	workersInfo, err := e.getNodeWorkerStatus(ctx, node.Endpoint, node.AccessToken)
	if err != nil || workersInfo == nil {
		zap.L().Error("get node worker status", zap.Error(err), zap.Any("address", node.Address.String()), zap.Any("endpoint", node.Endpoint), zap.Any("access_token", node.AccessToken))

		// Fixme: deprecated if the last node version is released
		if node.Status == schema.NodeStatusOffline {
			return schema.NodeStatusOffline, "workers_status"
		}

		return schema.NodeStatusInitializing, ""
	}

	// Check if any decentralized worker is unhealthy
	for _, workerInfo := range workersInfo.Data.Decentralized {
		if workerInfo.Status != worker.StatusUnhealthy {
			// Return initializing status if any worker is not unhealthy
			return schema.NodeStatusInitializing, ""
		}
	}

	// Check if any federated worker is unhealthy
	for _, workerInfo := range workersInfo.Data.Federated {
		if workerInfo.Status != worker.StatusUnhealthy {
			// Return initializing status if any worker is not unhealthy
			return schema.NodeStatusInitializing, ""
		}
	}

	// Check if any RSS worker is unhealthy
	if workersInfo.Data.RSS != nil {
		if workersInfo.Data.RSS.Status != worker.StatusUnhealthy {
			// Return initializing status if any worker is not unhealthy
			return schema.NodeStatusInitializing, ""
		}
	}

	// Return registered status if all checks pass
	return schema.NodeStatusRegistered, ""
}

// updateNodeStatuses updates node statuses and submits demotion information to VSL
func (e *SimpleEnforcer) updateNodeStatuses(ctx context.Context, updatedNodes []*schema.Node, demotionNodeAddresses []common.Address, reasons []string, reporters []common.Address) error {
	// Initialize node address and status lists
	nodeAddresses, nodeStatusList := make([]common.Address, 0, len(updatedNodes)), make([]uint8, 0, len(updatedNodes))

	// Iterate through updated nodes, collecting addresses and statuses
	for _, node := range updatedNodes {
		nodeAddresses = append(nodeAddresses, node.Address)
		nodeStatusList = append(nodeStatusList, uint8(node.Status))
	}

	// If there are no node statuses to update or nodes to demote, return immediately
	if len(nodeAddresses) == 0 && len(demotionNodeAddresses) == 0 {
		zap.L().Info("No node statuses need to be updated")
		return nil
	}

	// Call the method to update node statuses and submit demotions to VSL
	return e.updateNodeStatusAndSubmitDemotionToVSL(ctx, nodeAddresses, nodeStatusList, demotionNodeAddresses, reasons, reporters)
}

// saveOfflineStatusToInvalidResponse saves the offline status to the invalid response table
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

// updateNodeStatusAndSubmitDemotionToVSL updates node statuses and submits demotion information to VSL
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

// prepareSetNodeStatusAndSubmitDemotionsData prepares the data for setting node statuses and submitting demotions
func prepareSetNodeStatusAndSubmitDemotionsData(nodeAddresses []common.Address, nodeStatusList []uint8, demotionNodeAddresses []common.Address, reasons []string, reporters []common.Address) ([][]byte, error) {
	data := make([][]byte, 0)

	// Prepare data for setting node statuses
	if len(nodeAddresses) > 0 {
		inputSetNodeStatus, err := txmgr.EncodeInput(l2.SettlementMetaData.ABI, l2.MethodSetNodeStatus, nodeAddresses, nodeStatusList)
		if err != nil {
			return nil, fmt.Errorf("encode setNodeStatus input: %w", err)
		}

		data = append(data, inputSetNodeStatus)
	}

	// Prepare data for submitting demotions
	if len(demotionNodeAddresses) > 0 {
		inputSubmitDemotions, err := txmgr.EncodeInput(l2.SettlementMetaData.ABI, l2.MethodSubmitDemotions, demotionNodeAddresses, reasons, reporters)
		if err != nil {
			return nil, fmt.Errorf("encode submitDemotions input: %w", err)
		}

		data = append(data, inputSubmitDemotions)
	}

	return data, nil
}

// invokeSettlementMultiContract invokes the multicall contract on the VSL
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

// sendTransaction sends a transaction to the VSL
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

// getNodeMinVersion retrieves the minimum node version from the network params contract
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
