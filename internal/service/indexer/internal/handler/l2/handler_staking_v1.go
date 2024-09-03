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
	"github.com/rss3-network/global-indexer/common/ethereum"
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

func (h *handler) indexStakingV1Log(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; {
	case eventHash == l2.EventHashStakingV1Deposited:
		return h.indexStakingV1DepositedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case eventHash == l2.EventHashStakingV1WithdrawRequested:
		return h.indexStakingV1WithdrawRequestedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case eventHash == l2.EventHashStakingV1WithdrawalClaimed:
		return h.indexStakingV1WithdrawalClaimedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case eventHash == l2.EventHashStakingV1Staked:
		return h.indexStakingV1StakedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case eventHash == l2.EventHashStakingV1UnstakeRequested:
		return h.indexStakingV1UnstakeRequestedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case eventHash == l2.EventHashStakingV1UnstakeClaimed:
		return h.indexStakingV1UnstakeClaimedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case eventHash == l2.EventHashStakingV1RewardDistributed:
		return h.indexStakingV1RewardDistributedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case eventHash == l2.EventHashStakingV1NodeCreated:
		return h.indexStakingV1NodeCreated(ctx, header, transaction, receipt, log, databaseTransaction)
	case eventHash == l2.EventHashStakingV1NodeUpdated:
		return h.indexStakingV1NodeUpdated(ctx, header, transaction, receipt, log, databaseTransaction)
	default: // Discard all unsupported events.
		return nil
	}
}

func (h *handler) indexStakingV1DepositedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV1DepositedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingV1.ParseDeposited(*log)
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
		Finalized:        h.finalized,
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
		LogIndex:          log.Index,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
		Finalized:         h.finalized,
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}

func (h *handler) indexStakingV1WithdrawRequestedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV1WithdrawRequestedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingV1.ParseWithdrawRequested(*log)
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
		Finalized:        h.finalized,
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
		LogIndex:          log.Index,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
		Finalized:         h.finalized,
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}

func (h *handler) indexStakingV1WithdrawalClaimedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV1WithdrawalClaimedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingV1.ParseWithdrawalClaimed(*log)
	if err != nil {
		return fmt.Errorf("parse WithdrawalClaimed event: %w", err)
	}

	stakeEvent := schema.StakeEvent{
		ID:                common.BigToHash(event.RequestId),
		Type:              schema.StakeEventTypeWithdrawClaimed,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		LogIndex:          log.Index,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
		Finalized:         h.finalized,
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}

func (h *handler) indexStakingV1StakedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV1StakedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingV1.ParseStaked(*log)
	if err != nil {
		return fmt.Errorf("parse Staked event: %w", err)
	}

	callOptions := bind.CallOpts{
		Context:     ctx,
		BlockNumber: header.Number,
	}

	// If user staked token to a public good node, the event will be emitted with the genesis address.
	// So we need to get the actual node address from the stake contract by the token ID.
	if event.NodeAddr == ethereum.AddressGenesis {
		chipsInfo, err := h.contractStakingV1.GetChipsInfo(&callOptions, event.StartTokenId)
		if err != nil {
			return fmt.Errorf("get the info of chips %s: %w", event.StartTokenId, err)
		}

		// Override the node address with the actual node address.
		event.NodeAddr = chipsInfo.NodeAddr
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
		Finalized:        h.finalized,
	}

	for i := uint64(0); i+event.StartTokenId.Uint64() <= event.EndTokenId.Uint64(); i++ {
		stakeTransaction.ChipIDs = append(stakeTransaction.ChipIDs, new(big.Int).SetUint64(i+event.StartTokenId.Uint64()))
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
		LogIndex:          log.Index,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
		Finalized:         h.finalized,
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	resultPool := pool.
		NewWithResults[*schema.StakeChip]().
		WithContext(ctx).
		WithCancelOnError().
		WithFirstError()

	for _, chipID := range stakeTransaction.ChipIDs {
		chipID := chipID

		resultPool.Go(func(_ context.Context) (*schema.StakeChip, error) {
			tokenURI, err := h.contractChips.TokenURI(&callOptions, chipID)
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

			chipInfo, err := h.contractStakingV2.GetChipInfo(&callOptions, chipID)
			if err != nil {
				return nil, fmt.Errorf("get chip info from rpc: %w", err)
			}

			stakeChip := schema.StakeChip{
				ID:             chipID,
				Owner:          event.User,
				Node:           event.NodeAddr,
				Value:          decimal.NewFromBigInt(chipInfo.Tokens, 0),
				Metadata:       metadata,
				BlockNumber:    header.Number,
				BlockTimestamp: header.Time,
				Finalized:      h.finalized,
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

func (h *handler) indexStakingV1UnstakeRequestedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV1UnstakeRequestedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingV1.ParseUnstakeRequested(*log)
	if err != nil {
		return fmt.Errorf("parse UnstakeRequested event: %w", err)
	}

	stakeTransaction := schema.StakeTransaction{
		ID:               common.BigToHash(event.RequestId),
		Type:             schema.StakeTransactionTypeUnstake,
		User:             event.User,
		Node:             event.NodeAddr,
		Value:            event.UnstakeAmount,
		ChipIDs:          event.ChipsIds,
		BlockTimestamp:   time.Unix(int64(header.Time), 0),
		BlockNumber:      header.Number.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
		Finalized:        h.finalized,
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
		LogIndex:          log.Index,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
		Finalized:         h.finalized,
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}

func (h *handler) indexStakingV1UnstakeClaimedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV1UnstakeClaimedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingV1.ParseUnstakeClaimed(*log)
	if err != nil {
		return fmt.Errorf("parse UnstakeClaimed event: %w", err)
	}

	stakeEvent := schema.StakeEvent{
		ID:                common.BigToHash(event.RequestId),
		Type:              schema.StakeEventTypeUnstakeClaimed,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		LogIndex:          log.Index,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
		Finalized:         h.finalized,
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	return nil
}

func (h *handler) indexStakingV1RewardDistributedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV1RewardDistributedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingV1.ParseRewardDistributed(*log)
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
	if err := h.saveEpochRelatedNodes(ctx, databaseTransaction, &epoch); err != nil {
		zap.L().Error("indexRewardDistributedLog: save epoch related nodes", zap.Error(err), zap.String("transaction.hash", transaction.Hash().Hex()))

		return fmt.Errorf("save epoch related nodes: %w", err)
	}

	return nil
}

func (h *handler) indexStakingV1NodeCreated(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV1NodeCreated")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingV1.ParseNodeCreated(*log)
	if err != nil {
		return fmt.Errorf("parse NodeCreated event: %w", err)
	}

	addressTo := transaction.To()
	if addressTo == nil {
		addressTo = &l2.ContractMap[h.chainID].AddressStakingProxy
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
		ChainID:          h.chainID,
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
		Finalized: h.finalized,
	}

	if err := databaseTransaction.SaveNodeEvent(ctx, &nodeEvent); err != nil {
		return fmt.Errorf("save Node event: %w", err)
	}

	// Skip save node info if the block is not finalized.
	if !h.finalized {
		return nil
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
	if err := h.cacheClient.Get(ctx, h.buildNodeHideTaxRateKey(node.Address), &node.HideTaxRate); err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("get hide tax rate: %w", err)
	}

	// Save Node avatar
	avatar, err := h.contractStakingV1.GetNodeAvatar(&bind.CallOpts{}, event.NodeAddr)
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

func (h *handler) indexStakingV1NodeUpdated(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV1NodeUpdated")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	// Parse NodeUpdated event
	event, err := h.contractStakingV1.ParseNodeUpdated(*log)
	if err != nil {
		return fmt.Errorf("parse NodeUpdated event: %w", err)
	}

	// Query the Node from the contract
	node, err := h.contractStakingV1.GetNode(&bind.CallOpts{BlockNumber: header.Number}, event.NodeAddr)
	if err != nil {
		return fmt.Errorf("get Node: %w", err)
	}

	addressTo := transaction.To()
	if addressTo == nil {
		addressTo = &l2.ContractMap[h.chainID].AddressStakingProxy
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
		ChainID:          h.chainID,
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
		Finalized: true,
	}

	// Only save the event
	// Don't need to update the NodeInfo, because the fields are not stored in the database
	if err := databaseTransaction.SaveNodeEvent(ctx, &nodeEvent); err != nil {
		return fmt.Errorf("save Node event: %w", err)
	}

	return nil
}

func (h *handler) buildNodeHideTaxRateKey(address common.Address) string {
	return fmt.Sprintf("node::%s::hideTaxRate", strings.ToLower(address.String()))
}

func (h *handler) saveEpochRelatedNodes(ctx context.Context, databaseTransaction database.Client, epoch *schema.Epoch) error {
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
				apy     decimal.Decimal
				address = epoch.RewardedNodes[i].NodeAddress
			)

			// Calculate node APY
			node, err := h.contractStakingV1.GetNode(&bind.CallOpts{BlockNumber: epoch.BlockNumber}, address)
			if err != nil {
				zap.L().Error("indexRewardDistributedLog: Get node from rpc", zap.Error(err), zap.String("address", address.String()))

				return fmt.Errorf("get Node: %w", err)
			}

			// APY = (operationRewards + StakingRewards) / (StakingPoolTokens) * (1 - tax) * number of epochs in a year
			// number of epochs in a year = 365 * 24 / 18 = 486.6666666666667
			if node.StakingPoolTokens.Cmp(big.NewInt(0)) > 0 {
				tax := 1 - float64(node.TaxRateBasisPoints)/10000

				apy = epoch.RewardedNodes[i].OperationRewards.Add(epoch.RewardedNodes[i].StakingRewards).
					Div(decimal.NewFromBigInt(node.StakingPoolTokens, 0)).
					Mul(decimal.NewFromFloat(tax)).
					Mul(decimal.NewFromFloat(486.6666666666667))
			}

			data[i] = &schema.BatchUpdateNode{
				Address: address,
				Apy:     apy,
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
