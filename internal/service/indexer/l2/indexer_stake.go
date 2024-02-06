package l2

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"go.uber.org/zap"
)

func (s *server) indexStakingLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; eventHash {
	case l2.EventHashStakingDeposited:
		return s.indexStakingDepositedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingWithdrawRequested:
		return s.indexStakingWithdrawRequestedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingWithdrawalClaimed:
		return s.indexStakingWithdrawalClaimedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingStaked:
		return s.indexStakingStakedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingUnstakeRequested:
		return s.indexStakingUnstakeRequestedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingUnstakeClaimed:
		return s.indexStakingUnstakeClaimedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashRewardDistributed:
		return s.indexRewardDistributedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	default: // Discard all unsupported events.
		return nil
	}
}

func (s *server) indexStakingDepositedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseDeposited(*log)
	if err != nil {
		return fmt.Errorf("parse Deposited event: %w", err)
	}

	user, err := types.Sender(types.LatestSignerForChainID(transaction.ChainId()), transaction)
	if err != nil {
		return fmt.Errorf("invalid transaction signer: %w", err)
	}

	stakeTransaction := schema.StakeTransaction{
		ID:               transaction.Hash(),
		Type:             schema.StakeTransactionTypeDeposit,
		User:             user,
		Node:             event.NodeAddr,
		Value:            event.Amount,
		BlockTimestamp:   time.Unix(int64(header.Time), 0),
		BlockNumber:      header.Number.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
	}

	if err := databaseTransaction.SaveStakeTransaction(ctx, &stakeTransaction); err != nil {
		return fmt.Errorf("save stake transaction: %w", err)
	}

	stakeEvent := schema.StakeEvent{
		ID:                transaction.Hash(),
		Type:              schema.StakeEventTypeDepositDeposited,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}

func (s *server) indexStakingWithdrawRequestedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseWithdrawRequested(*log)
	if err != nil {
		return fmt.Errorf("parse WithdrawRequested event: %w", err)
	}

	user, err := types.Sender(types.LatestSignerForChainID(transaction.ChainId()), transaction)
	if err != nil {
		return fmt.Errorf("invalid transaction signer: %w", err)
	}

	stakeTransaction := schema.StakeTransaction{
		ID:               common.BigToHash(event.RequestId),
		Type:             schema.StakeTransactionTypeWithdraw,
		User:             user,
		Node:             event.NodeAddr,
		Value:            event.Amount,
		BlockTimestamp:   time.Unix(int64(header.Time), 0),
		BlockNumber:      header.Number.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
	}

	if err := databaseTransaction.SaveStakeTransaction(ctx, &stakeTransaction); err != nil {
		return fmt.Errorf("save stake transaction: %w", err)
	}

	stakeEvent := schema.StakeEvent{
		ID:                common.BigToHash(event.RequestId),
		Type:              schema.StakeEventTypeWithdrawRequested,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}

func (s *server) indexStakingWithdrawalClaimedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseWithdrawalClaimed(*log)
	if err != nil {
		return fmt.Errorf("parse WithdrawalClaimed event: %w", err)
	}

	stakeEvent := schema.StakeEvent{
		ID:                common.BigToHash(event.RequestId),
		Type:              schema.StakeEventTypeWithdrawClaimed,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}

func (s *server) indexStakingStakedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseStaked(*log)
	if err != nil {
		return fmt.Errorf("parse Staked event: %w", err)
	}

	stakeTransaction := schema.StakeTransaction{
		ID:               transaction.Hash(),
		Type:             schema.StakeTransactionTypeStake,
		User:             event.User,
		Node:             event.NodeAddr,
		Value:            event.Amount,
		BlockTimestamp:   time.Unix(int64(header.Time), 0),
		BlockNumber:      header.Number.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
	}

	for i := uint64(0); i+event.StartTokenId.Uint64() <= event.EndTokenId.Uint64(); i++ {
		stakeTransaction.Chips = append(stakeTransaction.Chips, new(big.Int).SetUint64(i+event.StartTokenId.Uint64()))
	}

	if err := databaseTransaction.SaveStakeTransaction(ctx, &stakeTransaction); err != nil {
		return fmt.Errorf("save stake transaction: %w", err)
	}

	stakeEvent := schema.StakeEvent{
		ID:                transaction.Hash(),
		Type:              schema.StakeEventTypeStakeStaked,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	stakeChips := make([]*schema.StakeChip, len(stakeTransaction.Chips))

	for index, chipID := range stakeTransaction.Chips {
		stakeChips[index] = &schema.StakeChip{
			ID:    chipID,
			Owner: event.User,
			Node:  event.NodeAddr,
		}
	}

	if err := databaseTransaction.SaveStakeChips(ctx, stakeChips...); err != nil {
		return fmt.Errorf("save stake chips: %w", err)
	}

	return nil
}

func (s *server) indexStakingUnstakeRequestedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseUnstakeRequested(*log)
	if err != nil {
		return fmt.Errorf("parse UnstakeRequested event: %w", err)
	}

	stakeTransaction := schema.StakeTransaction{
		ID:               common.BigToHash(event.RequestId),
		Type:             schema.StakeTransactionTypeUnstake,
		User:             event.User,
		Node:             event.NodeAddr,
		Value:            event.UnstakeAmount,
		Chips:            event.ChipsIds,
		BlockTimestamp:   time.Unix(int64(header.Time), 0),
		BlockNumber:      header.Number.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
	}

	if err := databaseTransaction.SaveStakeTransaction(ctx, &stakeTransaction); err != nil {
		return fmt.Errorf("save stake transaction: %w", err)
	}

	stakeEvent := schema.StakeEvent{
		ID:                common.BigToHash(event.RequestId),
		Type:              schema.StakeEventTypeUnstakeRequested,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}

func (s *server) indexStakingUnstakeClaimedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseUnstakeClaimed(*log)
	if err != nil {
		return fmt.Errorf("parse UnstakeClaimed event: %w", err)
	}

	stakeEvent := schema.StakeEvent{
		ID:                common.BigToHash(event.RequestId),
		Type:              schema.StakeEventTypeUnstakeClaimed,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}

func (s *server) indexRewardDistributedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseRewardDistributed(*log)
	if err != nil {
		return fmt.Errorf("parse RewardDistributed event: %w", err)
	}

	epoch := schema.Epoch{
		ID:               event.Epoch.Uint64(),
		StartTimestamp:   event.StartTimestamp.Int64(),
		EndTimestamp:     event.EndTimestamp.Int64(),
		TransactionHash:  transaction.Hash(),
		BlockNumber:      header.Number,
		TransactionIndex: receipt.TransactionIndex,
		Success:          receipt.Status == types.ReceiptStatusSuccessful,
		TotalRewardItems: len(event.NodeAddrs),
		RewardItems:      make([]*schema.EpochItem, 0, len(event.NodeAddrs)),
	}

	if epoch.TotalRewardItems != len(event.StakingRewards) || epoch.TotalRewardItems != len(event.OperationRewards) ||
		epoch.TotalRewardItems != len(event.RequestFees) || epoch.TotalRewardItems != len(event.TaxAmounts) {
		zap.L().Error("indexRewardDistributedLog: length not match", zap.Int("length", epoch.TotalRewardItems), zap.String("transaction.hash", transaction.Hash().Hex()))

		return fmt.Errorf("length not match")
	}

	if epoch.Success {
		var totalOperationRewards, totalStakingRewards big.Int

		for i := 0; i < epoch.TotalRewardItems; i++ {
			epoch.RewardItems = append(epoch.RewardItems, &schema.EpochItem{
				EpochID:          event.Epoch.Uint64(),
				Index:            i,
				NodeAddress:      event.NodeAddrs[i],
				RequestFees:      event.RequestFees[i].String(),
				OperationRewards: event.OperationRewards[i].String(),
				StakingRewards:   event.StakingRewards[i].String(),
				TaxAmounts:       event.TaxAmounts[i].String(),
			})

			totalOperationRewards.Add(&totalOperationRewards, event.OperationRewards[i])
			totalStakingRewards.Add(&totalStakingRewards, event.StakingRewards[i])
		}

		epoch.TotalOperationRewards = totalOperationRewards.String()
		epoch.TotalStakingRewards = totalStakingRewards.String()
	}

	if err := databaseTransaction.SaveEpoch(ctx, &epoch); err != nil {
		zap.L().Error("indexRewardDistributedLog: save epoch", zap.Error(err), zap.String("transaction.hash", transaction.Hash().Hex()))

		return fmt.Errorf("save epoch: %w", err)
	}

	return nil
}
