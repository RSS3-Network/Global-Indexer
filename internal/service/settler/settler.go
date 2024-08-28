package settler

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
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
		transactionData, nodes, scores, err := s.constructSettlementData(ctx, epoch, cursor)
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

		// Update the Node scores
		if len(nodes) > 0 {
			err = s.updateNodesScore(ctx, scores, nodes)
			if err != nil {
				zap.L().Error("failed to update node scores", zap.Error(err))
			}
		}

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
func (s *Server) constructSettlementData(ctx context.Context, epoch uint64, cursor *string) (*schema.SettlementData, []*schema.Node, []*big.Int, error) {
	// batchSize is the number of Nodes to process in each batch.
	// This is to prevent the contract call from running out of gas.
	batchSize := s.settlerConfig.BatchSize

	// Find qualified Nodes from the database
	nodes, err := s.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
		Status:  lo.ToPtr(schema.NodeStatusOnline),
		Version: lo.ToPtr(schema.NodeVersionNormal),
		Cursor:  cursor,
		Limit:   lo.ToPtr(batchSize + 1),
	})
	if err != nil {
		// No qualified Nodes found in the database
		if errors.Is(err, database.ErrorRowNotFound) {
			return nil, nil, nil, nil
		}

		zap.L().Error("No qualified Nodes found", zap.Error(err), zap.Any("cursor", cursor))

		return nil, nil, nil, err
	}

	// isFinal is true if it's the last batch of Nodes
	isFinal := len(nodes) <= batchSize
	if !isFinal {
		nodes = nodes[:batchSize]
	}

	// nodeAddresses is a slice of Node addresses
	nodeAddresses := make([]common.Address, 0, len(nodes))
	for _, node := range nodes {
		nodeAddresses = append(nodeAddresses, node.Address)
	}

	// Calculate the number of requests for the Nodes
	requestCount, totalRequestCount, err := s.prepareRequestCounts(ctx, nodeAddresses)
	if err != nil {
		return nil, nil, nil, err
	}

	// Calculate the Operation rewards for the Nodes
	operationRewards, err := calculateOperationRewards(nodes, requestCount, s.operationRewards)
	if err != nil {
		return nil, nil, nil, err
	}

	return &schema.SettlementData{
		Epoch:            big.NewInt(int64(epoch)),
		NodeAddress:      nodeAddresses,
		OperationRewards: operationRewards,
		RequestCount:     requestCount,
		IsFinal:          isFinal,
	}, nodes, totalRequestCount, nil
}

func (s *Server) updateNodesScore(ctx context.Context, scores []*big.Int, nodes []*schema.Node) error {
	for i, node := range nodes {
		node.ActiveScore = decimal.NewFromBigInt(scores[i], 0)
	}

	return s.databaseClient.UpdateNodesScore(ctx, nodes)
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
		GasLimit: s.settlerConfig.GasLimit,
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
