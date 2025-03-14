package settler

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// submitEpochProof submits proof of this epoch on chain
// which calculates the Operation Rewards for the Nodes
// formats the data and invokes the contract
// a retry logic is implemented to handle possible failures
func (s *Server) submitEpochProof(ctx context.Context, epoch uint64) error {
	if err := s.mutex.Lock(); err != nil {
		zap.L().Error("lock error", zap.String("key", s.mutex.Name()), zap.Error(err))

		return nil
	}

	defer func() {
		if _, err := s.mutex.Unlock(); err != nil {
			zap.L().Error("release lock error", zap.String("key", s.mutex.Name()), zap.Error(err))
		}
	}()

	var (
		cursor      *string
		firstInvoke = true
	)

	for {
		msg := "construct Settlement data"
		// Construct transactionData as required by the Settlement contract
		transactionData, err := s.constructSettlementData(ctx, epoch, cursor)
		if err != nil {
			zap.L().Error(msg, zap.Error(err))

			return fmt.Errorf("%s: %w", msg, err)
		}

		// Finish processing when conditions are met
		if len(transactionData.NodeAddress) == 0 && !firstInvoke {
			zap.L().Info("finished processing transactionData.")

			break
		}

		zap.L().Info(msg, zap.Any("transactionData", transactionData))

		// Invoke the Settlement contract
		receipt, err := retry.DoWithData(
			func() (*types.Receipt, error) {
				return s.invokeSettlementContract(ctx, *transactionData)
			},
			retry.Delay(time.Second),
			retry.Attempts(5),
		)
		if err != nil {
			zap.L().Error("retry submitEpochProof invokeSettlementContract", zap.Error(err))

			return err
		}

		// Save the Settlement to the database, as the reference point for the next Epoch
		if err := s.saveSettlement(ctx, receipt, *transactionData); err != nil {
			return err
		}

		zap.L().Info("Settlement contracted invoked successfully", zap.String("tx", receipt.TxHash.String()), zap.Any("data", *transactionData))

		firstInvoke = false

		if len(transactionData.NodeAddress) > 0 {
			cursor = lo.ToPtr(transactionData.NodeAddress[len(transactionData.NodeAddress)-1].String())
		}
	}

	zap.L().Info("Epoch Proof submitted successfully", zap.Uint64("settler", epoch))

	return nil
}

// retryEpochProof retries the epoch proof submission.
// When a block reorganization occurs, the original epoch proof needs to be resubmitted.
func (s *Server) retryEpochProof(ctx context.Context, epochID uint64) error {
	if err := s.mutex.Lock(); err != nil {
		zap.L().Error("lock error", zap.String("key", s.mutex.Name()), zap.Error(err))

		return nil
	}

	defer func() {
		if _, err := s.mutex.Unlock(); err != nil {
			zap.L().Error("release lock error", zap.String("key", s.mutex.Name()), zap.Error(err))
		}
	}()

	// Find the EpochTrigger by the epochID
	epochTriggers, err := s.databaseClient.FindEpochTriggers(ctx, epochID)
	if err != nil {
		zap.L().Error("find epoch triggers", zap.Error(err))

		return err
	}

	for _, trigger := range epochTriggers {
		// Invoke the Settlement contract
		receipt, err := retry.DoWithData(func() (*types.Receipt, error) {
			return s.invokeSettlementContract(ctx, trigger.Data)
		}, retry.Delay(time.Second), retry.Attempts(5))

		if err != nil {
			zap.L().Error("retry submitEpochProof invokeSettlementContract", zap.Error(err))

			return err
		}

		// Skip saving the Settlement to the database, just log the result
		zap.L().Info("Settlement contracted invoked successfully", zap.Uint64("epoch_id", epochID), zap.String("tx", receipt.TxHash.String()), zap.Any("data", trigger.Data))
	}

	return nil
}

// constructSettlementData constructs Settlement data as required by the Settlement contract
func (s *Server) constructSettlementData(ctx context.Context, epoch uint64, cursor *string) (*schema.SettlementData, error) {
	// batchSize is the number of Nodes to process in each batch.
	// This is to prevent the contract call from running out of gas.
	// TODO: This method needs to be refactored when the number of nodes exceeds the batch size value.
	batchSize := s.config.Settler.BatchSize

	// Find qualified Nodes from the database
	query := schema.FindNodesQuery{
		Status: lo.ToPtr(schema.NodeStatusOnline),
		Cursor: cursor,
		Limit:  lo.ToPtr(batchSize + 1),
	}

	// Set the Node version to Normal after the grace period
	if epoch >= uint64(s.config.Settler.ProductionStartEpoch+s.config.Settler.GracePeriodEpochs) {
		query.Type = lo.ToPtr(schema.NodeTypeProduction)
	}

	nodes, err := s.databaseClient.FindNodes(ctx, query)

	if err != nil {
		// No qualified Nodes found in the database
		if errors.Is(err, database.ErrorRowNotFound) {
			return nil, nil
		}

		zap.L().Error("No qualified Nodes found", zap.Error(err), zap.Any("cursor", cursor))

		return nil, err
	}

	// isFinal is true if it's the last batch of Nodes
	isFinal := len(nodes) <= batchSize
	if !isFinal {
		nodes = nodes[:batchSize]
	}

	filterNodeAddresses, filterNodes, err := s.filter(nodes)
	if err != nil {
		return nil, err
	}

	// Calculate the number of requests for the Nodes
	requestCount, operationStats, err := s.prepareRequestCounts(ctx, filterNodeAddresses, filterNodes)
	if err != nil {
		return nil, err
	}

	// Calculate the Operation rewards for the Nodes
	operationRewards, err := s.calculateOperationRewards(ctx, operationStats, s.config.Rewards)
	if err != nil {
		return nil, err
	}

	return &schema.SettlementData{
		Epoch:            big.NewInt(int64(epoch)),
		NodeAddress:      filterNodeAddresses,
		OperationRewards: operationRewards,
		RequestCount:     requestCount,
		IsFinal:          isFinal,
	}, nil
}

// filter retrieves Node information from a staking contract.
func (s *Server) filter(nodes []*schema.Node) ([]common.Address, []*schema.Node, error) {
	nodeAddresses := lo.Map(nodes, func(node *schema.Node, _ int) common.Address {
		return node.Address
	})

	nodeInfoList, err := s.stakingContract.GetNodes(&bind.CallOpts{}, nodeAddresses)
	if err != nil {
		return nil, nil, fmt.Errorf("get Nodes from chain: %w", err)
	}

	nodeInfoMap := lo.SliceToMap(nodeInfoList, func(node stakingv2.Node) (common.Address, stakingv2.Node) {
		return node.Account, node
	})

	newNodeAddresses := make([]common.Address, 0, len(nodeAddresses))
	newNodes := make([]*schema.Node, 0, len(nodeAddresses))

	for i := range nodeAddresses {
		if nodeInfo, ok := nodeInfoMap[nodeAddresses[i]]; ok && isValidStatus(nodeInfo.Status) {
			newNodeAddresses = append(newNodeAddresses, nodeAddresses[i])
			newNodes = append(newNodes, nodes[i])
		}
	}

	return newNodeAddresses, newNodes, nil
}

// isValidStatus checks if the node status is valid.
func isValidStatus(status uint8) bool {
	return status == uint8(schema.NodeStatusInitializing) ||
		status == uint8(schema.NodeStatusOnline) ||
		status == uint8(schema.NodeStatusExiting) ||
		status == uint8(schema.NodeStatusSlashing)
}

// invokeSettlementContract invokes the Settlement contract with prepared data
// and saves the Settlement to the database
func (s *Server) invokeSettlementContract(ctx context.Context, data schema.SettlementData) (*types.Receipt, error) {
	input, err := s.prepareInputData(data)
	if err != nil {
		return nil, err
	}

	receipt, err := s.sendTransaction(ctx, input)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// sendTransaction sends the transaction and returns the receipt if successful
func (s *Server) sendTransaction(ctx context.Context, input []byte) (*types.Receipt, error) {
	txCandidate := txmgr.TxCandidate{
		TxData:   input,
		To:       lo.ToPtr(l2.ContractMap[s.chainID.Uint64()].AddressSettlementProxy),
		GasLimit: s.config.Settler.GasLimit,
		Value:    big.NewInt(0),
	}

	receipt, err := s.txManager.Send(ctx, txCandidate)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		zap.L().Error("received an invalid transaction receipt", zap.String("tx", receipt.TxHash.String()))

		// select {} purposely block the process as it is a critical error and meaningless to continue
		// if panic() is called, the process will be restarted by the supervisor
		// we do not want that as it will be stuck in the same state
		select {}
	}

	// return the receipt if the transaction is successful
	return receipt, nil
}

// saveSettlement saves the Settlement data to the database
func (s *Server) saveSettlement(ctx context.Context, receipt *types.Receipt, data schema.SettlementData) error {
	if err := s.databaseClient.SaveEpochTrigger(ctx, &schema.EpochTrigger{
		TransactionHash: receipt.TxHash,
		EpochID:         data.Epoch.Uint64(),
		Data:            data,
	}); err != nil {
		return fmt.Errorf("save settler submitEpochProof: %w", err)
	}

	return nil
}
