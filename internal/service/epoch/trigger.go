package epoch

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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

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
			return fmt.Errorf("find online nodes: %w", err)
		}

		if data == nil {
			return nil
		}

		zap.L().Info("build distributeRewards", zap.Any("data", data))

		// Trigger distributeReward contract.
		if err = retry.Do(func() error {
			return s.triggerDistributeRewards(ctx, lo.FromPtr(data))
		}, retry.Delay(time.Second), retry.Attempts(5)); err != nil {
			zap.L().Error("retry trigger distributeReward", zap.Error(err))

			return err
		}

		cursor = lo.ToPtr(data.NodeAddress[len(data.NodeAddress)-1].String())
	}
}

func (s *Server) buildDistributeRewards(ctx context.Context, epoch uint64, cursor *string) (*schema.DistributeRewardsData, error) {
	nodes, err := s.databaseClient.FindNodes(ctx, nil, lo.ToPtr(schema.StatusOnline), cursor, BatchSize+1)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return nil, nil
		}

		zap.L().Error("find online nodes", zap.Error(err), zap.Any("cursor", cursor))

		return nil, err
	}

	if len(nodes) == 0 {
		return nil, nil
	}

	var isFinal = true

	if len(nodes) > BatchSize {
		nodes = nodes[:BatchSize]
		isFinal = false
	}

	nodeAddress := make([]common.Address, 0, len(nodes))

	for _, node := range nodes {
		nodeAddress = append(nodeAddress, node.Address)
	}

	zeroRewards := make([]*big.Int, len(nodes))

	for i := range zeroRewards {
		zeroRewards[i] = big.NewInt(0)
	}

	return &schema.DistributeRewardsData{
		Epoch:            big.NewInt(int64(epoch)),
		NodeAddress:      nodeAddress,
		RequestFees:      zeroRewards,
		OperationRewards: zeroRewards,
		IsFinal:          isFinal,
	}, nil
}

func (s *Server) triggerDistributeRewards(ctx context.Context, data schema.DistributeRewardsData) error {
	// Trigger distributeReward contract.
	fromAddress := crypto.PubkeyToAddress(s.privateKey.PublicKey)

	nonce, err := s.ethereumClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return fmt.Errorf("get pending nonce: %w", err)
	}

	transactor, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainID)
	if err != nil {
		return fmt.Errorf("create transactor: %w", err)
	}

	transactor.Nonce = big.NewInt(int64(nonce))
	transactor.Value = big.NewInt(0)
	transactor.GasLimit = s.gasLimit

	transactor.GasPrice, err = s.ethereumClient.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("get gas price: %w", err)
	}

	tx, err := s.settlementContract.DistributeRewards(transactor, data.Epoch, data.NodeAddress, data.RequestFees, data.OperationRewards)
	if err != nil {
		zap.L().Error("distribute rewards", zap.Error(err), zap.Any("data", data))

		return fmt.Errorf("distribute rewards: %w", err)
	}

	// Save epoch trigger to database.
	if err = s.databaseClient.SaveEpochTrigger(ctx, &schema.EpochTrigger{
		TransactionHash: tx.Hash(),
		EpochID:         data.Epoch.Uint64(),
		Data:            data,
	}); err != nil {
		return fmt.Errorf("save epoch trigger: %w", err)
	}

	// Wait for transaction receipt.
	if err = s.transactionReceipt(ctx, tx.Hash()); err != nil {
		zap.L().Error("wait for transaction receipt", zap.Error(err), zap.Any("data", data))

		return fmt.Errorf("wait for transaction receipt: %w", err)
	}

	zap.L().Info("distribute rewards successfully", zap.String("tx", tx.Hash().String()), zap.Any("data", data))

	return nil
}

func (s *Server) transactionReceipt(ctx context.Context, txHash common.Hash) error {
	for {
		receipt, err := s.ethereumClient.TransactionReceipt(ctx, txHash)
		if err != nil {
			zap.L().Warn("wait for transaction", zap.Error(err), zap.String("tx", txHash.String()))

			continue
		}

		if receipt.Status == types.ReceiptStatusSuccessful {
			return nil
		}

		if receipt.Status == types.ReceiptStatusFailed {
			return fmt.Errorf("transaction failed: %s", receipt.TxHash.String())
		}
	}
}
