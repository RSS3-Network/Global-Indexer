package epoch

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (s *Server) trigger(ctx context.Context, epoch uint64) error {
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

		// Trigger distributeReward contract.
		if err := s.triggerDistributeRewards(ctx, lo.FromPtr(data)); err != nil {
			return fmt.Errorf("trigger distributeReward: %w", err)
		}
	}
}

func (s *Server) buildDistributeRewards(ctx context.Context, epoch uint64, cursor *string) (*schema.DistributeRewardsData, error) {
	nodes, err := s.databaseClient.FindNodes(ctx, nil, lo.ToPtr(schema.StatusOnline), cursor, 200)
	if err != nil {
		zap.L().Error("find online nodes", zap.Error(err), zap.Any("cursor", cursor))

		return nil, err
	}

	if len(nodes) == 0 {
		return nil, nil
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

	settlement, err := l2.NewSettlement(l2.AddressSettlementProxy, s.ethereumClient)
	if err != nil {
		return fmt.Errorf("new settlement: %w", err)
	}

	tx, err := settlement.DistributeRewards(transactor, data.Epoch, data.NodeAddress, data.RequestFees, data.OperationRewards)
	if err != nil {
		zap.L().Error("distribute rewards", zap.Error(err), zap.String("tx", tx.Hash().String()), zap.Any("data", data))

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

	zap.L().Info("distribute rewards successfully", zap.String("tx", tx.Hash().String()), zap.Any("data", data))

	return nil
}
