package l2

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/shopspring/decimal"
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
	case l2.EventHashStakingRewardDistributed:
		return s.indexStakingRewardDistributedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingNodeCreated:
		return s.indexStakingNodeCreated(ctx, header, transaction, receipt, log, databaseTransaction)
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

	callOptions := bind.CallOpts{
		Context:     ctx,
		BlockNumber: header.Number,
	}

	for index, chipID := range stakeTransaction.Chips {
		tokenURI, err := s.contractChips.TokenURI(&callOptions, chipID)
		if err != nil {
			return fmt.Errorf("get #%d token uri", chipID)
		}

		encodedMetadata, found := strings.CutPrefix(tokenURI, "data:application/json;base64,")
		if !found {
			return fmt.Errorf("invalid #%d token uri", chipID)
		}

		metadata, err := base64.StdEncoding.DecodeString(encodedMetadata)
		if err != nil {
			return fmt.Errorf("decode #%d token metadata", chipID)
		}

		value, err := s.contractStaking.MinTokensToStake(&callOptions, stakeTransaction.Node)
		if err != nil {
			return fmt.Errorf("get the minimum stake requirement for node %s", stakeTransaction.Node)
		}

		stakeChips[index] = &schema.StakeChip{
			ID:             chipID,
			Owner:          event.User,
			Node:           event.NodeAddr,
			Value:          decimal.NewFromBigInt(value, 0),
			Metadata:       metadata,
			BlockNumber:    header.Number,
			BlockTimestamp: header.Time,
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

func (s *server) indexStakingRewardDistributedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseRewardDistributed(*log)
	if err != nil {
		return fmt.Errorf("parse RewardDistributed event: %w", err)
	}

	epoch := schema.Epoch{
		ID:               event.Epoch.Uint64(),
		StartTimestamp:   event.StartTimestamp.Int64(),
		EndTimestamp:     event.EndTimestamp.Int64(),
		TransactionHash:  transaction.Hash(),
		BlockHash:        header.Hash(),
		BlockNumber:      header.Number,
		BlockTimestamp:   int64(header.Time),
		TransactionIndex: receipt.TransactionIndex,
		TotalRewardItems: len(event.NodeAddrs),
		RewardItems:      make([]*schema.EpochItem, 0, len(event.NodeAddrs)),
	}

	if epoch.TotalRewardItems != len(event.StakingRewards) || epoch.TotalRewardItems != len(event.OperationRewards) || epoch.TotalRewardItems != len(event.TaxAmounts) {
		zap.L().Error("indexRewardDistributedLog: length not match", zap.Int("length", epoch.TotalRewardItems), zap.String("transaction.hash", transaction.Hash().Hex()))

		return fmt.Errorf("length not match")
	}

	var totalOperationRewards, totalStakingRewards big.Int

	for i := 0; i < epoch.TotalRewardItems; i++ {
		epoch.RewardItems = append(epoch.RewardItems, &schema.EpochItem{
			EpochID:          event.Epoch.Uint64(),
			Index:            i,
			TransactionHash:  transaction.Hash(),
			NodeAddress:      event.NodeAddrs[i],
			OperationRewards: event.OperationRewards[i].String(),
			StakingRewards:   event.StakingRewards[i].String(),
			TaxAmounts:       event.TaxAmounts[i].String(),
		})

		totalOperationRewards.Add(&totalOperationRewards, event.OperationRewards[i])
		totalStakingRewards.Add(&totalStakingRewards, event.StakingRewards[i])
	}

	epoch.TotalOperationRewards = totalOperationRewards.String()
	epoch.TotalStakingRewards = totalStakingRewards.String()

	if err := databaseTransaction.SaveEpoch(ctx, &epoch); err != nil {
		zap.L().Error("indexRewardDistributedLog: save epoch", zap.Error(err), zap.String("transaction.hash", transaction.Hash().Hex()))

		return fmt.Errorf("save epoch: %w", err)
	}

	return nil
}

func (s *server) indexStakingNodeCreated(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseNodeCreated(*log)
	if err != nil {
		return fmt.Errorf("parse NodeCreated event: %w", err)
	}

	// save createdNode event
	nodeEvent := schema.NodeEvent{
		TransactionHash:  transaction.Hash(),
		TransactionIndex: receipt.TransactionIndex,
		AddressFrom:      event.NodeAddr,
		AddressTo:        receipt.ContractAddress,
		Type:             schema.NodeEventNodeCreated,
		LogIndex:         log.Index,
		ChainID:          s.chainID.Uint64(),
		BlockHash:        header.Hash(),
		BlockNumber:      header.Number.Uint64(),
		BlockTimestamp:   int64(header.Time),
		Metadata: schema.NodeEventMetadata{
			NodeCreatedMetadata: &schema.NodeCreatedMetadata{
				Address:            event.NodeAddr,
				Name:               event.Name,
				Description:        event.Description,
				TaxRateBasisPoints: event.TaxRateBasisPoints,
				PublicGood:         event.PublicGood,
			},
		},
	}

	if err := databaseTransaction.SaveNodeEvent(ctx, &nodeEvent); err != nil {
		return fmt.Errorf("save node event: %w", err)
	}

	// save node
	node := schema.Node{
		Address:            event.NodeAddr,
		Name:               event.Name,
		Endpoint:           event.NodeAddr.String(), // initial endpoint
		Description:        event.Description,
		TaxRateBasisPoints: event.TaxRateBasisPoints,
		IsPublicGood:       event.PublicGood,
		Status:             schema.NodeStatusRegistered,
		Stream:             json.RawMessage{},
		Config:             json.RawMessage{},
	}

	if err := databaseTransaction.SaveNode(ctx, &node); err != nil {
		return fmt.Errorf("save node: %w", err)
	}

	return nil
}
