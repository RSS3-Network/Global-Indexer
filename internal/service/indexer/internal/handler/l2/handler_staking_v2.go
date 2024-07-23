package l2

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/common/ethereum"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
	"github.com/sourcegraph/conc/pool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (h *handler) indexStakingV2Log(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; eventHash {
	case l2.EventHashStakingV1Staked: // The event hash is the same as the staking v1 contract.
		return h.indexStakingV2StakedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingV2ChipsMerged:
		return h.indexStakingV2ChipsMergedLog(ctx, header, transaction, receipt, log, databaseTransaction)
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

	event, err := h.contractStakingV2.ParseStaked(*log)
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

	for _, chipID := range stakeTransaction.Chips {
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

func (h *handler) indexStakingV2ChipsMergedLog(ctx context.Context, header *types.Header, transaction *types.Transaction, _ *types.Receipt, log *types.Log, _ database.Client) error {
	_, span := otel.Tracer("").Start(ctx, "indexStakingV2ChipsMergedLog")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("block.number", header.Number.Int64()),
		attribute.Stringer("block.hash", header.Hash()),
		attribute.Stringer("transaction.hash", transaction.Hash()),
		attribute.Int("log.index", int(log.Index)),
	)

	if _, err := h.contractStakingV2.ParseChipsMerged(*log); err != nil {
		return fmt.Errorf("parse ChipsMerged event: %w", err)
	}

	return nil
}
