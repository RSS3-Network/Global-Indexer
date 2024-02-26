package epoch

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

	// billing
	err := s.billingFlow(ctx)
	if err != nil {
		return err
	}

	// distribute rewards

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

func (s *Server) prepareTransactor(ctx context.Context) (*bind.TransactOpts, error) {
	fromAddress := crypto.PubkeyToAddress(s.privateKey.PublicKey)

	nonce, err := s.ethereumClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("get pending nonce: %w", err)
	}

	transactor, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainID)
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}

	transactor.Nonce = big.NewInt(int64(nonce))
	transactor.Value = big.NewInt(0)
	transactor.GasLimit = s.gasLimit

	transactor.GasPrice, err = s.ethereumClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("get gas price: %w", err)
	}

	return transactor, nil
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
