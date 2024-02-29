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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
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

func (s *Server) buildDistributeRewards(ctx context.Context, epoch uint64, cursor *string) (*schema.DistributeRewardsData, error) {
	nodes, err := s.databaseClient.FindNodes(ctx, nil, lo.ToPtr(schema.StatusOnline), cursor, BatchSize+1)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return nil, nil
		}

		zap.L().Error("find online nodes", zap.Error(err), zap.Any("cursor", cursor))

		return nil, err
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
	nonce, err := s.ethereumClient.PendingNonceAt(ctx, s.fromAddress)
	if err != nil {
		return fmt.Errorf("get pending nonce: %w", err)
	}

	gasPrice, err := s.ethereumClient.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("get gas price: %w", err)
	}

	input, err := s.encodeInput(l2.SettlementMetaData.ABI, l2.MethodDistributeRewards, data.Epoch, data.NodeAddress, data.RequestFees, data.OperationRewards, data.IsFinal)
	if err != nil {
		return fmt.Errorf("encode input: %w", err)
	}

	unsignedTX := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      s.gasLimit,
		To:       lo.ToPtr(l2.ContractMap[s.chainID.Uint64()].AddressSettlementProxy),
		Value:    big.NewInt(0),
		Data:     input,
	})

	args := s.newTransactionArgsFromTransaction(s.chainID, s.fromAddress, unsignedTX)

	var result hexutil.Bytes
	if err = s.rpcClient.CallContext(ctx, &result, "eth_signTransaction", args); err != nil {
		return fmt.Errorf("eth_signTransaction failed: %w", err)
	}

	signedTX := &types.Transaction{}
	if err = signedTX.UnmarshalBinary(result); err != nil {
		return err
	}

	if err = s.ethereumClient.SendTransaction(ctx, signedTX); err != nil {
		zap.L().Error("distribute rewards", zap.Error(err), zap.Any("data", data))

		return fmt.Errorf("distribute rewards: %w", err)
	}

	// Save epoch trigger to database.
	if err = s.databaseClient.SaveEpochTrigger(ctx, &schema.EpochTrigger{
		TransactionHash: signedTX.Hash(),
		EpochID:         data.Epoch.Uint64(),
		Data:            data,
	}); err != nil {
		return fmt.Errorf("save epoch trigger: %w", err)
	}

	// Wait for transaction receipt.
	if err = s.transactionReceipt(ctx, signedTX.Hash()); err != nil {
		zap.L().Error("wait for transaction receipt", zap.Error(err), zap.Any("data", data))

		return fmt.Errorf("wait for transaction receipt: %w", err)
	}

	zap.L().Info("distribute rewards successfully", zap.String("tx", signedTX.Hash().String()), zap.Any("data", data))

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

type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
}

func (s *Server) newTransactionArgsFromTransaction(chainID *big.Int, from common.Address, tx *types.Transaction) *TransactionArgs {
	data := hexutil.Bytes(tx.Data())
	nonce := hexutil.Uint64(tx.Nonce())
	gas := hexutil.Uint64(tx.Gas())
	accesses := tx.AccessList()

	return &TransactionArgs{
		From:                 &from,
		Input:                &data,
		Nonce:                &nonce,
		Value:                (*hexutil.Big)(tx.Value()),
		Gas:                  &gas,
		To:                   tx.To(),
		ChainID:              (*hexutil.Big)(chainID),
		MaxFeePerGas:         (*hexutil.Big)(tx.GasFeeCap()),
		MaxPriorityFeePerGas: (*hexutil.Big)(tx.GasTipCap()),
		AccessList:           &accesses,
	}
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
