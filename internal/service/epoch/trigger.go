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
	"github.com/naturalselectionlabs/rss3-global-indexer/common/txmgr"
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
	nodes, err := s.databaseClient.FindNodes(ctx, nil, lo.ToPtr(schema.NodeStatusOnline), cursor, BatchSize+1)
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
		OperationRewards: zeroRewards,
		IsFinal:          isFinal,
	}, nil
}

func (s *Server) triggerDistributeRewards(ctx context.Context, data schema.DistributeRewardsData) error {
	input, err := s.encodeInput(l2.SettlementMetaData.ABI, l2.MethodDistributeRewards, data.Epoch, data.NodeAddress, data.OperationRewards, data.IsFinal)
	if err != nil {
		return fmt.Errorf("encode input: %w", err)
	}

	txCandidate := txmgr.TxCandidate{
		TxData:   input,
		To:       lo.ToPtr(l2.ContractMap[s.chainID.Uint64()].AddressSettlementProxy),
		GasLimit: s.gasLimit,
		Value:    big.NewInt(0),
	}

	receipt, err := s.txManager.Send(ctx, txCandidate)

	if err != nil {
		return fmt.Errorf("send tx failed %w", err)
	}

	// Save epoch trigger to database.
	if err = s.databaseClient.SaveEpochTrigger(ctx, &schema.EpochTrigger{
		TransactionHash: receipt.TxHash,
		EpochID:         data.Epoch.Uint64(),
		Data:            data,
	}); err != nil {
		return fmt.Errorf("save epoch trigger: %w", err)
	}

	zap.L().Info("distribute rewards successfully", zap.String("tx", receipt.TxHash.String()), zap.Any("data", data))

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
