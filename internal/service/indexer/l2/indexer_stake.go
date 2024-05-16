package l2

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/sourcegraph/conc/pool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	case l2.EventHashStakingNodeUpdated:
		return s.indexStakingNodeUpdated(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingNodeUpdated2PublicGood:
		return s.indexStakingNodeUpdated2PublicGood(ctx, header, transaction, receipt, log, databaseTransaction)
	default: // Discard all unsupported events.
		return nil
	}
}

func (s *server) indexStakingDepositedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingDepositedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

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
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingWithdrawRequestedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

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
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingWithdrawalClaimedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

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
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingStakedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

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

	callOptions := bind.CallOpts{
		Context:     ctx,
		BlockNumber: header.Number,
	}

	resultPool := pool.
		NewWithResults[*schema.StakeChip]().
		WithContext(ctx).
		WithCancelOnError().
		WithFirstError()

	for _, chipID := range stakeTransaction.Chips {
		chipID := chipID

		resultPool.Go(func(_ context.Context) (*schema.StakeChip, error) {
			tokenURI, err := s.contractChips.TokenURI(&callOptions, chipID)
			if err != nil {
				return nil, fmt.Errorf("get #%d token uri", chipID)
			}

			encodedMetadata, found := strings.CutPrefix(tokenURI, "data:application/json;base64,")
			if !found {
				return nil, fmt.Errorf("invalid #%d token uri", chipID)
			}

			metadata, err := base64.StdEncoding.DecodeString(encodedMetadata)
			if err != nil {
				return nil, fmt.Errorf("decode #%d token metadata", chipID)
			}

			value, err := s.contractStaking.MinTokensToStake(&callOptions, stakeTransaction.Node)
			if err != nil {
				return nil, fmt.Errorf("get the minimum stake requirement for node %s", stakeTransaction.Node)
			}

			stakeChip := schema.StakeChip{
				ID:             chipID,
				Owner:          event.User,
				Node:           event.NodeAddr,
				Value:          decimal.NewFromBigInt(value, 0),
				Metadata:       metadata,
				BlockNumber:    header.Number,
				BlockTimestamp: header.Time,
			}

			return &stakeChip, nil
		})
	}

	stakeChips, err := resultPool.Wait()
	if err != nil {
		return fmt.Errorf("get chips: %w", err)
	}

	if err := databaseTransaction.SaveStakeChips(ctx, stakeChips...); err != nil {
		return fmt.Errorf("save stake chips: %w", err)
	}

	return nil
}

func (s *server) indexStakingUnstakeRequestedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingUnstakeRequestedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

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
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingUnstakeClaimedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

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
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingRewardDistributedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := s.contractStaking.ParseRewardDistributed(*log)
	if err != nil {
		return fmt.Errorf("parse RewardDistributed event: %w", err)
	}

	epoch := schema.Epoch{
		ID:                 event.Epoch.Uint64(),
		StartTimestamp:     event.StartTimestamp.Int64(),
		EndTimestamp:       event.EndTimestamp.Int64(),
		TransactionHash:    transaction.Hash(),
		BlockHash:          header.Hash(),
		BlockNumber:        header.Number,
		BlockTimestamp:     int64(header.Time),
		TransactionIndex:   receipt.TransactionIndex,
		TotalRewardedNodes: len(event.NodeAddrs),
		RewardedNodes:      make([]*schema.RewardedNode, len(event.NodeAddrs)),
	}

	if epoch.TotalRewardedNodes != len(event.StakingRewards) || epoch.TotalRewardedNodes != len(event.OperationRewards) || epoch.TotalRewardedNodes != len(event.TaxCollected) {
		zap.L().Error("indexRewardDistributedLog: length not match", zap.Int("length", epoch.TotalRewardedNodes), zap.String("transaction.hash", transaction.Hash().Hex()))

		return fmt.Errorf("length not match")
	}

	for i := 0; i < epoch.TotalRewardedNodes; i++ {
		epoch.RewardedNodes[i] = &schema.RewardedNode{
			EpochID:          event.Epoch.Uint64(),
			Index:            i,
			TransactionHash:  transaction.Hash(),
			NodeAddress:      event.NodeAddrs[i],
			OperationRewards: decimal.NewFromBigInt(event.OperationRewards[i], 0),
			StakingRewards:   decimal.NewFromBigInt(event.StakingRewards[i], 0),
			TaxCollected:     decimal.NewFromBigInt(event.TaxCollected[i], 0),
			RequestCount:     decimal.NewFromBigInt(event.RequestCounts[i], 0),
		}

		epoch.TotalOperationRewards = epoch.TotalOperationRewards.Add(epoch.RewardedNodes[i].OperationRewards)
		epoch.TotalStakingRewards = epoch.TotalStakingRewards.Add(epoch.RewardedNodes[i].StakingRewards)
		epoch.TotalRequestCounts = epoch.TotalRequestCounts.Add(epoch.RewardedNodes[i].RequestCount)
	}

	// Save epoch
	if err := databaseTransaction.SaveEpoch(ctx, &epoch); err != nil {
		zap.L().Error("indexRewardDistributedLog: save epoch", zap.Error(err), zap.String("transaction.hash", transaction.Hash().Hex()))

		return fmt.Errorf("save epoch: %w", err)
	}

	// Skip if no Nodes were rewarded in this Epoch.
	if epoch.TotalRewardedNodes == 0 {
		return nil
	}

	// Save Nodes
	if err := s.saveEpochRelatedNodes(ctx, databaseTransaction, &epoch); err != nil {
		zap.L().Error("indexRewardDistributedLog: save epoch related nodes", zap.Error(err), zap.String("transaction.hash", transaction.Hash().Hex()))

		return fmt.Errorf("save epoch related nodes: %w", err)
	}

	return nil
}

func (s *server) indexStakingNodeCreated(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingNodeCreated")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := s.contractStaking.ParseNodeCreated(*log)
	if err != nil {
		return fmt.Errorf("parse NodeCreated event: %w", err)
	}

	addressTo := transaction.To()
	if addressTo == nil {
		addressTo = &l2.ContractMap[s.chainID.Uint64()].AddressStakingProxy
	}

	// save createdNode event
	nodeEvent := schema.NodeEvent{
		TransactionHash:  transaction.Hash(),
		TransactionIndex: receipt.TransactionIndex,
		NodeID:           event.NodeId,
		AddressFrom:      event.NodeAddr,
		AddressTo:        lo.FromPtr(addressTo),
		Type:             schema.NodeEventNodeCreated,
		LogIndex:         log.Index,
		ChainID:          s.chainID.Uint64(),
		BlockHash:        header.Hash(),
		BlockNumber:      header.Number,
		BlockTimestamp:   int64(header.Time),
		Metadata: schema.NodeEventMetadata{
			NodeCreatedMetadata: &schema.NodeCreatedMetadata{
				NodeID:             event.NodeId,
				Address:            event.NodeAddr,
				Name:               event.Name,
				Description:        event.Description,
				TaxRateBasisPoints: event.TaxRateBasisPoints,
				PublicGood:         event.PublicGood,
			},
		},
	}

	if err := databaseTransaction.SaveNodeEvent(ctx, &nodeEvent); err != nil {
		return fmt.Errorf("save Node event: %w", err)
	}

	// if node already exists, skip
	if node, _ := databaseTransaction.FindNode(ctx, event.NodeAddr); node != nil {
		return nil
	}

	// save node
	node := &schema.Node{
		Address:            event.NodeAddr,
		ID:                 event.NodeId,
		Name:               event.Name,
		Endpoint:           event.NodeAddr.String(), // initial endpoint
		Description:        event.Description,
		TaxRateBasisPoints: &event.TaxRateBasisPoints,
		IsPublicGood:       event.PublicGood,
		Status:             schema.NodeStatusRegistered,
	}

	// Get from redis if the tax rate of the Node needs to be hidden.
	if err := s.cacheClient.Get(ctx, s.buildNodeHideTaxRateKey(node.Address), &node.HideTaxRate); err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("get hide tax rate: %w", err)
	}

	// Get node minTokensToStake
	minTokensToStake, err := s.contractStaking.MinTokensToStake(&bind.CallOpts{BlockNumber: header.Number}, event.NodeAddr)
	if err != nil {
		return fmt.Errorf("get min tokens to stake: %w", err)
	}

	node.MinTokensToStake = decimal.NewFromBigInt(minTokensToStake, 0)

	// Save Node avatar
	avatar, err := s.contractStaking.GetNodeAvatar(&bind.CallOpts{}, event.NodeAddr)
	if err != nil {
		return fmt.Errorf("get Node avatar: %w", err)
	}

	encodedMetadata, ok := strings.CutPrefix(avatar, "data:application/json;base64,")
	if !ok {
		return fmt.Errorf("invalid avatar: %s", avatar)
	}

	metadata, err := base64.StdEncoding.DecodeString(encodedMetadata)
	if err != nil {
		return fmt.Errorf("decode avatar metadata: %w", err)
	}

	if err = json.Unmarshal(metadata, &node.Avatar); err != nil {
		return fmt.Errorf("unmarshal avatar metadata: %w", err)
	}

	if err := databaseTransaction.SaveNode(ctx, node); err != nil {
		return fmt.Errorf("save Node: %w", err)
	}

	return nil
}

func (s *server) indexStakingNodeUpdated(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingNodeUpdated")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	// Parse NodeUpdated event
	event, err := s.contractStaking.ParseNodeUpdated(*log)
	if err != nil {
		return fmt.Errorf("parse NodeUpdated event: %w", err)
	}

	// Query the Node from the contract
	node, err := s.contractStaking.GetNode(&bind.CallOpts{BlockNumber: header.Number}, event.NodeAddr)
	if err != nil {
		return fmt.Errorf("get Node: %w", err)
	}

	addressTo := transaction.To()
	if addressTo == nil {
		addressTo = &l2.ContractMap[s.chainID.Uint64()].AddressStakingProxy
	}

	// Save NodeUpdated event
	nodeEvent := schema.NodeEvent{
		TransactionHash:  transaction.Hash(),
		TransactionIndex: receipt.TransactionIndex,
		NodeID:           node.NodeId,
		AddressFrom:      event.NodeAddr,
		AddressTo:        lo.FromPtr(addressTo),
		Type:             schema.NodeEventNodeUpdated,
		LogIndex:         log.Index,
		ChainID:          s.chainID.Uint64(),
		BlockHash:        header.Hash(),
		BlockNumber:      header.Number,
		BlockTimestamp:   int64(header.Time),
		Metadata: schema.NodeEventMetadata{
			NodeUpdatedMetadata: &schema.NodeUpdatedMetadata{
				Address:     event.NodeAddr,
				Name:        event.Name,
				Description: event.Description,
			},
		},
	}

	// Only save the event
	// Don't need to update the NodeInfo, because the fields are not stored in the database
	if err := databaseTransaction.SaveNodeEvent(ctx, &nodeEvent); err != nil {
		return fmt.Errorf("save Node event: %w", err)
	}

	return nil
}

func (s *server) indexStakingNodeUpdated2PublicGood(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingNodeUpdated2PublicGood")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	// Parse NodeUpdated2PublicGood event
	event, err := s.contractStaking.ParseNodeUpdated2PublicGood(*log)
	if err != nil {
		return fmt.Errorf("parse NodeUpdated2PublicGood event: %w", err)
	}

	// Query the Node from the contract
	node, err := s.contractStaking.GetNode(&bind.CallOpts{BlockNumber: header.Number}, event.NodeAddr)
	if err != nil {
		return fmt.Errorf("get Node: %w", err)
	}

	addressTo := transaction.To()
	if addressTo == nil {
		addressTo = &l2.ContractMap[s.chainID.Uint64()].AddressStakingProxy
	}

	// Build Node event
	nodeEvent := schema.NodeEvent{
		TransactionHash:  transaction.Hash(),
		TransactionIndex: receipt.TransactionIndex,
		NodeID:           node.NodeId,
		AddressFrom:      event.NodeAddr,
		AddressTo:        lo.FromPtr(addressTo),
		Type:             schema.NodeEventNodeUpdated2PublicGood,
		LogIndex:         log.Index,
		ChainID:          s.chainID.Uint64(),
		BlockHash:        header.Hash(),
		BlockNumber:      header.Number,
		BlockTimestamp:   int64(header.Time),
		Metadata: schema.NodeEventMetadata{
			NodeUpdated2PublicGoodMetadata: &schema.NodeUpdated2PublicGoodMetadata{
				Address:    event.NodeAddr,
				PublicGood: true,
			},
		},
	}

	// Save NodeUpdated2PublicGood event to the database
	if err := databaseTransaction.SaveNodeEvent(ctx, &nodeEvent); err != nil {
		return fmt.Errorf("save Node event: %w", err)
	}

	// Update the Node's public good status
	if err := databaseTransaction.UpdateNodePublicGood(ctx, event.NodeAddr, true); err != nil {
		return fmt.Errorf("update Node public good status: %w", err)
	}

	return nil
}

func (s *server) buildNodeHideTaxRateKey(address common.Address) string {
	return fmt.Sprintf("node::%s::hideTaxRate", strings.ToLower(address.String()))
}

func (s *server) saveEpochRelatedNodes(ctx context.Context, databaseTransaction database.Client, epoch *schema.Epoch) error {
	ctx, span := otel.Tracer("").Start(ctx, "saveEpochRelatedNodes")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", epoch.BlockNumber.Int64()),
		attribute.Stringer("block.hash", epoch.BlockHash),
		attribute.Stringer("transaction.hash", epoch.TransactionHash),
		attribute.Int64("epoch.id", int64(epoch.ID)),
	)

	var (
		data      = make([]*schema.BatchUpdateNode, len(epoch.RewardedNodes))
		errorPool = pool.New().WithContext(ctx).WithMaxGoroutines(50).WithCancelOnError().WithFirstError()
	)

	for i := 0; i < epoch.TotalRewardedNodes; i++ {
		i := i

		errorPool.Go(func(_ context.Context) error {
			var (
				apy, minTokensToStake decimal.Decimal
				address               = epoch.RewardedNodes[i].NodeAddress
			)

			// Calculate node APY
			node, err := s.contractStaking.GetNode(&bind.CallOpts{BlockNumber: epoch.BlockNumber}, address)
			if err != nil {
				zap.L().Error("indexRewardDistributedLog: Get node from rpc", zap.Error(err), zap.String("address", address.String()))

				return fmt.Errorf("get Node: %w", err)
			}

			// APY = (operationRewards + stakingRewards) / (stakingPoolTokens) * (1 - tax) * number of epochs in a year
			// number of epochs in a year = 365 * 24 / 18 = 486.6666666666667
			if node.StakingPoolTokens.Cmp(big.NewInt(0)) > 0 {
				tax := 1 - float64(node.TaxRateBasisPoints)/10000

				apy = epoch.RewardedNodes[i].OperationRewards.Add(epoch.RewardedNodes[i].StakingRewards).
					Div(decimal.NewFromBigInt(node.StakingPoolTokens, 0)).
					Mul(decimal.NewFromFloat(tax)).
					Mul(decimal.NewFromFloat(486.6666666666667))
			}

			// Query the minTokensToStake of the Node
			minTokens, err := s.contractStaking.MinTokensToStake(&bind.CallOpts{BlockNumber: epoch.BlockNumber}, address)
			if err != nil {
				zap.L().Error("indexRewardDistributedLog: get min tokens to stake", zap.Error(err), zap.String("address", address.String()))

				return fmt.Errorf("get min tokens to stake: %w", err)
			}

			minTokensToStake = decimal.NewFromBigInt(minTokens, 0)

			data[i] = &schema.BatchUpdateNode{
				Address:          address,
				Apy:              apy,
				MinTokensToStake: minTokensToStake,
			}

			return nil
		})
	}

	if err := errorPool.Wait(); err != nil {
		return fmt.Errorf("wait error pool: %w", err)
	}

	// Save Nodes
	if err := databaseTransaction.BatchUpdateNodes(ctx, data); err != nil {
		zap.L().Error("batch update epoch-related nodes", zap.Error(err), zap.Any("epoch", epoch), zap.Any("data", data))

		return fmt.Errorf("batch update epoch-related nodes: %w", err)
	}

	return nil
}
