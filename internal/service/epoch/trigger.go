package epoch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/txmgr"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// BatchSize is the number of Nodes to process in each batch.
// This is to prevent the contract call from running out of gas.
const BatchSize = 200

func (s *Server) trigger(ctx context.Context, epoch uint64) error {
	if err := s.mutex.Lock(); err != nil {
		zap.L().Error("lock error", zap.String("key", s.mutex.Name()), zap.Error(err))

		return nil
	}

	defer func() {
		if _, err := s.mutex.Unlock(); err != nil {
			zap.L().Error("release lock error", zap.String("key", s.mutex.Name()), zap.Error(err))
		}
	}()

	var cursor *string

	for {
		// Build distribute rewards data.
		data, err := s.buildDistributeRewards(ctx, epoch, cursor)
		if err != nil {
			zap.L().Error("finding online nodes", zap.Error(err))

			return fmt.Errorf("find online nodes: %w", err)
		}

		// Check data existence.
		if len(data.NodeAddress) == 0 && cursor != nil {
			zap.L().Info("no more data to process. exiting")

			break
		}

		zap.L().Info("build distributeRewards", zap.Any("data", data))

		// Trigger distributeReward contract.
		if err = retry.Do(func() error {
			return s.triggerDistributeRewards(ctx, *data)
		}, retry.Delay(time.Second), retry.Attempts(5)); err != nil {
			zap.L().Error("retry trigger distributeReward", zap.Error(err))

			return err
		}

		if len(data.NodeAddress) > 0 {
			cursor = lo.ToPtr(data.NodeAddress[len(data.NodeAddress)-1].String())
		}
	}

	zap.L().Info("Reward distribution completed")

	return nil
}

// buildDistributeRewards builds the distribute rewards data struct
func (s *Server) buildDistributeRewards(ctx context.Context, epoch uint64, cursor *string) (*schema.DistributeRewardsData, error) {
	nodes, err := s.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
		Status: lo.ToPtr(schema.NodeStatusOnline),
		Cursor: cursor,
		Limit:  lo.ToPtr(BatchSize + 1),
	})
	if err != nil {
		// No Nodes in the database.
		if errors.Is(err, database.ErrorRowNotFound) {
			return nil, nil
		}

		zap.L().Error("find online nodes", zap.Error(err), zap.Any("cursor", cursor))

		return nil, err
	}

	// isFinal is true if it's the last batch of Nodes
	isFinal := len(nodes) <= BatchSize
	if !isFinal {
		nodes = nodes[:BatchSize]
	}

	// nodeAddresses is a slice of Node addresses.
	nodeAddresses := make([]common.Address, 0, len(nodes))
	for _, node := range nodes {
		nodeAddresses = append(nodeAddresses, node.Address)
	}

	operationRewards := calculateOperationRewards(nodeAddresses)

	return &schema.DistributeRewardsData{
		Epoch:            big.NewInt(int64(epoch)),
		NodeAddress:      nodeAddresses,
		OperationRewards: operationRewards,
		IsFinal:          isFinal,
	}, nil
}

// calculateOperationRewards calculates the operation rewards for Nodes
// For Alpha, there is no actual calculation logic.
// TODO: Implement the actual calculation logic
func calculateOperationRewards(nodes []common.Address) []*big.Int {
	slice := make([]*big.Int, len(nodes))

	// For Alpha, set the rewards to 0
	for i := range slice {
		slice[i] = big.NewInt(0)
	}

	return slice
}

func (s *Server) triggerDistributeRewards(ctx context.Context, data schema.DistributeRewardsData) error {
	input, err := s.prepareInputData(data)
	if err != nil {
		return err
	}

	receipt, err := s.sendTransaction(ctx, input)
	if err != nil {
		return err
	}

	// Save epoch trigger to database.
	if err := s.handleReceipt(ctx, receipt, data); err != nil {
		return err
	}

	zap.L().Info("rewards distributed successfully", zap.String("tx", receipt.TxHash.String()), zap.Any("data", data))

	return nil
}

// prepareInputData encodes input data for the transaction.
func (s *Server) prepareInputData(data schema.DistributeRewardsData) ([]byte, error) {
	input, err := s.encodeInput(l2.SettlementMetaData.ABI, l2.MethodDistributeRewards, data.Epoch, data.NodeAddress, data.OperationRewards, data.IsFinal)
	if err != nil {
		return nil, fmt.Errorf("encode input: %w", err)
	}

	return input, nil
}

// sendTransaction sends the transaction and returns the receipt.
func (s *Server) sendTransaction(ctx context.Context, input []byte) (*types.Receipt, error) {
	txCandidate := txmgr.TxCandidate{
		TxData:   input,
		To:       lo.ToPtr(l2.ContractMap[s.chainID.Uint64()].AddressSettlementProxy),
		GasLimit: s.gasLimit,
		Value:    big.NewInt(0),
	}

	receipt, err := s.txManager.Send(ctx, txCandidate)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		zap.L().Error("invalid transaction receipt", zap.String("tx", receipt.TxHash.String()))

		// select {} purposely block the process as it is a critical error and meaningless to continue
		// if panic() is called, the process will be restarted by the supervisor
		// we do not want that as it will be stuck in the same state
		select {}
	}

	return receipt, nil
}

// handleReceipt processes the transaction receipt and updates the database
func (s *Server) handleReceipt(ctx context.Context, receipt *types.Receipt, data schema.DistributeRewardsData) error {
	if err := s.databaseClient.SaveEpochTrigger(ctx, &schema.EpochTrigger{
		TransactionHash: receipt.TxHash,
		EpochID:         data.Epoch.Uint64(),
		Data:            data,
	}); err != nil {
		return fmt.Errorf("save epoch trigger: %w", err)
	}

	return nil
}

func (s *Server) encodeInput(contractABI, methodName string, args ...interface{}) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}

	encodedArgs, err := parsedABI.Pack(methodName, args...)
	if err != nil {
		return nil, err
	}

	return encodedArgs, nil
}
