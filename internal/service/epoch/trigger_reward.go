package epoch

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

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
	transactor, err := s.prepareTransactor(ctx)
	if err != nil {
		zap.L().Error("prepare transactor", zap.Error(err))
		return fmt.Errorf("prepare transactor: %w", err)
	}

	tx, err := s.settlementContract.DistributeRewards(transactor, data.Epoch, data.NodeAddress, data.RequestFees, data.OperationRewards, data.IsFinal)
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
