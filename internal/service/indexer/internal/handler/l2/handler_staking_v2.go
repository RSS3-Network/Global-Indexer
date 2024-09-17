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
	"github.com/rss3-network/global-indexer/common/ethereum"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
	"github.com/sourcegraph/conc/pool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (h *handler) indexStakingV2Log(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	if eventHash := log.Topics[0]; eventHash == common.HexToHash("0x8adb7a84b2998a8d11cd9284395f95d5a99f160be785ae79998c654979bd3d9a") {
		zap.L().Info("staking v2 log", zap.String("event", "Staked"), zap.String("hash", eventHash.Hex()))
	}

	switch eventHash := log.Topics[0]; eventHash {
	case l2.EventHashStakingV1Staked: // The event hash is the same as the staking v1 contract.
		return h.indexStakingV2StakedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingV2ChipsMerged:
		return h.indexStakingV2ChipsMergedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingV2WithdrawalClaimed:
		return h.indexStakingV2WithdrawalClaimedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashNodeStatusChanged:
		return h.indexNodeStatusChangedLog(ctx, header, transaction, log, databaseTransaction)
	default:
		return h.indexStakingV1Log(ctx, header, transaction, receipt, log, databaseTransaction)
	}
}

func (h *handler) indexStakingV2StakedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV2StakedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingEvents.ParseStaked(*log)
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
		chipsInfo, err := h.contractStakingV2.GetChipInfo(&callOptions, event.StartTokenId)
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
				return nil, fmt.Errorf("get chips #%d info", chipID)
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

func (h *handler) indexStakingV2ChipsMergedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	_, span := otel.Tracer("").Start(ctx, "indexStakingV2ChipsMergedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingEvents.ParseChipsMerged(*log)
	if err != nil {
		return fmt.Errorf("parse ChipsMerged event: %w", err)
	}

	callOptions := bind.CallOpts{
		Context:     ctx,
		BlockNumber: header.Number,
	}

	stakeTransaction := schema.StakeTransaction{
		ID:               transaction.Hash(),
		Type:             schema.StakeTransactionTypeMergeChips,
		User:             event.User,
		Node:             event.NodeAddr,
		ChipIDs:          append(event.BurnedTokenIds, event.NewTokenId),
		BlockTimestamp:   time.Unix(int64(header.Time), 0),
		BlockNumber:      header.Number.Uint64(),
		TransactionIndex: receipt.TransactionIndex,
		Value:            big.NewInt(0),
		Finalized:        h.finalized,
	}

	if err := databaseTransaction.SaveStakeTransaction(ctx, &stakeTransaction); err != nil {
		return fmt.Errorf("save stake transaction: %w", err)
	}

	metadata, err := json.Marshal(schema.StakeEventChipsMergedMetadata{
		BurnedTokenIDs: event.BurnedTokenIds,
		NewTokenID:     event.NewTokenId,
	})
	if err != nil {
		return fmt.Errorf("marshal chips merged metadata: %w", err)
	}

	stakeEvent := schema.StakeEvent{
		ID:                transaction.Hash(),
		Type:              schema.StakeEventTypeChipsMerged,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		LogIndex:          log.Index,
		Metadata:          metadata,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
		Finalized:         h.finalized,
	}

	if err := databaseTransaction.SaveStakeEvent(ctx, &stakeEvent); err != nil {
		return fmt.Errorf("save stake event: %w", err)
	}

	// Save New Chip
	tokenURI, err := h.contractChips.TokenURI(&callOptions, event.NewTokenId)
	if err != nil {
		return fmt.Errorf("get #%d token uri", event.NewTokenId)
	}

	encodedMetadata, found := strings.CutPrefix(tokenURI, "data:application/json;base64,")
	if !found {
		return fmt.Errorf("invalid #%d token uri", event.NewTokenId)
	}

	chipMetadata, err := base64.StdEncoding.DecodeString(encodedMetadata)
	if err != nil {
		return fmt.Errorf("decode #%d token metadata", event.NewTokenId)
	}

	chipInfo, err := h.contractStakingV2.GetChipInfo(&callOptions, event.NewTokenId)
	if err != nil {
		return fmt.Errorf("get chips #%d info", event.NewTokenId)
	}

	stakeChip := &schema.StakeChip{
		ID:             event.NewTokenId,
		Owner:          event.User,
		Node:           event.NodeAddr,
		Value:          decimal.NewFromBigInt(chipInfo.Tokens, 0),
		Metadata:       chipMetadata,
		BlockNumber:    header.Number,
		BlockTimestamp: header.Time,
		Finalized:      h.finalized,
	}

	if err := databaseTransaction.SaveStakeChips(ctx, stakeChip); err != nil {
		return fmt.Errorf("save stake chips: %w", err)
	}

	return nil
}

func (h *handler) indexStakingV2WithdrawalClaimedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexStakingV2WithdrawalClaimedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingEvents.ParseWithdrawalClaimed(*log)
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

func (h *handler) indexNodeStatusChangedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, log *types.Log, databaseTransaction database.Client) error {
	ctx, span := otel.Tracer("").Start(ctx, "indexNodeStatusChangedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	event, err := h.contractStakingEvents.ParseNodeStatusChanged(*log)
	if err != nil {
		return fmt.Errorf("parse NodeStatusChanged event: %w", err)
	}

	nodeAddress := event.NodeAddr
	nodeCurrentStatus := event.CurStatus
	nodeNewStatus := event.NewStatus

	switch nodeNewStatus {
	case uint8(schema.NodeStatusSlashing):
		return h.handleNodeSlashing(ctx, nodeAddress, nodeCurrentStatus, databaseTransaction)
	// TODO: node status reverted to online from slashing
	// case uint8(schema.NodeStatusOnline):
	//	 return h.handleNodeOnline(ctx, nodeAddress, nodeCurrentStatus, databaseTransaction)
	default:
		return nil
	}
}

func (h *handler) handleNodeSlashing(ctx context.Context, nodeAddress common.Address, nodeCurrentStatus uint8, databaseTransaction database.Client) error {
	if nodeCurrentStatus == uint8(schema.NodeStatusOnline) || nodeCurrentStatus == uint8(schema.NodeStatusExiting) {
		zap.L().Info("node status changed", zap.Stringer("node", nodeAddress), zap.String("new status", "Slashing"))

		if err := h.removeNodeFromCache(ctx, nodeAddress); err != nil {
			return err
		}

		return h.markNodeAsSlashed(ctx, nodeAddress, databaseTransaction)
	}

	return nil
}

func (h *handler) removeNodeFromCache(ctx context.Context, nodeAddress common.Address) error {
	if err := h.cacheClient.ZRem(ctx, model.FullNodeCacheKey, nodeAddress.String()); err != nil {
		return err
	}

	return h.cacheClient.ZRem(ctx, model.RssNodeCacheKey, nodeAddress.String())
}

func (h *handler) markNodeAsSlashed(ctx context.Context, nodeAddress common.Address, databaseTransaction database.Client) error {
	nodeStat, err := databaseTransaction.FindNodeStat(ctx, nodeAddress)
	if err != nil {
		return fmt.Errorf("find node stat: %w", err)
	}

	nodeStat.EpochInvalidRequest = int64(model.DemotionCountBeforeSlashing)

	return databaseTransaction.SaveNodeStat(ctx, nodeStat)
}
