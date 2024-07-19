package l2

import (
	"context"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/indexer/internal"
	"github.com/samber/lo"
)

var _ internal.Handler = (*handler)(nil)

type handler struct {
	chainID                        uint64
	finalized                      bool
	ethereumClient                 *ethclient.Client
	cacheClient                    cache.Client
	contractGovernanceToken        *bindings.GovernanceToken
	contractL2CrossDomainMessenger *bindings.L2CrossDomainMessenger
	contractL2StandardBridge       *bindings.L2StandardBridge
	contractL2ToL1MessagePasser    *bindings.L2ToL1MessagePasser
	contractStaking                *l2.Staking
	contractChips                  *l2.Chips
}

func (h *handler) Process(ctx context.Context, block *types.Block, receipts types.Receipts, databaseTransaction database.Client) error {
	if err := h.deleteUnfinalizedBlock(ctx, block.NumberU64(), databaseTransaction); err != nil {
		return fmt.Errorf("delete unfinalized block: %w", err)
	}

	if h.finalized {
		if err := h.confirmPreviousBlocks(ctx, block.NumberU64(), databaseTransaction); err != nil {
			return fmt.Errorf("confirm previous blocks: %w", err)
		}
	}

	header := block.Header()

	for _, receipt := range receipts {
		// Discard all contract creation transactions.
		if block.Transaction(receipt.TxHash).To() == nil {
			continue
		}

		// Discard all failed transactions.
		if receipt.Status != types.ReceiptStatusSuccessful {
			continue
		}

		for logIndex, log := range receipt.Logs {
			// Discard all removed logs.
			if log.Removed {
				continue
			}

			// Discard all anonymous logs.
			if len(log.Topics) == 0 {
				continue
			}

			switch log.Address {
			case l2.AddressL2StandardBridgeProxy:
				if err := h.indexBridgingLog(ctx, header, block.Transaction(log.TxHash), receipt, log, logIndex, databaseTransaction); err != nil {
					return fmt.Errorf("index bridge log: %w", err)
				}
			case l2.ContractMap[h.chainID].AddressStakingProxy:
				if err := h.indexStakingLog(ctx, header, block.Transaction(log.TxHash), receipt, log, databaseTransaction); err != nil {
					return fmt.Errorf("index staking v2 log: %w", err)
				}
			case l2.ContractMap[h.chainID].AddressChipsProxy:
				if err := h.indexChipsLog(ctx, header, block.Transaction(log.TxHash), receipt, log, databaseTransaction); err != nil {
					return fmt.Errorf("index staking log: %w", err)
				}
			}
		}
	}

	return nil
}

func (h *handler) deleteUnfinalizedBlock(ctx context.Context, blockNumber uint64, databaseTransaction database.Client) error {
	if err := databaseTransaction.DeleteBridgeTransactionsByBlockNumber(ctx, h.chainID, blockNumber); err != nil {
		return fmt.Errorf("delete bridge transactions by block number: %w", err)
	}

	if err := databaseTransaction.DeleteBridgeEventsByBlockNumber(ctx, h.chainID, blockNumber); err != nil {
		return fmt.Errorf("delete bridge events by block number: %w", err)
	}

	if err := databaseTransaction.DeleteStakeChipsByBlockNumber(ctx, blockNumber); err != nil {
		return fmt.Errorf("delete stake chips by block number: %w", err)
	}

	if err := databaseTransaction.DeleteStakeTransactionsByBlockNumber(ctx, blockNumber); err != nil {
		return fmt.Errorf("delete stake transactions by block number: %w", err)
	}

	if err := databaseTransaction.DeleteStakeEventsByBlockNumber(ctx, blockNumber); err != nil {
		return fmt.Errorf("delete stake events by block number: %w", err)
	}

	if err := databaseTransaction.DeleteNodeEventsByBlockNumber(ctx, blockNumber); err != nil {
		return fmt.Errorf("delete node events by block number: %w", err)
	}

	return nil
}

func (h *handler) confirmPreviousBlocks(ctx context.Context, blockNumber uint64, databaseTransaction database.Client) error {
	if err := databaseTransaction.UpdateStakeTransactionsFinalizedByBlockNumber(ctx, blockNumber); err != nil {
		return fmt.Errorf("update bridge transactions finalized by block number: %w", err)
	}

	if err := databaseTransaction.UpdateStakeEventsFinalizedByBlockNumber(ctx, blockNumber); err != nil {
		return fmt.Errorf("update bridge events finalized by block number: %w", err)
	}

	if err := databaseTransaction.UpdateStakeChipsFinalizedByBlockNumber(ctx, blockNumber); err != nil {
		return fmt.Errorf("update bridge chips finalized by block number: %w", err)
	}

	return nil
}

func NewHandler(chainID uint64, ethereumClient *ethclient.Client, cacheClient cache.Client, finalized bool) (internal.Handler, error) {
	contractAddresses := l2.ContractMap[chainID]
	if contractAddresses == nil {
		return nil, fmt.Errorf("chain id %d is not supported", chainID)
	}

	instance := handler{
		chainID:                        chainID,
		finalized:                      finalized,
		ethereumClient:                 ethereumClient,
		cacheClient:                    cacheClient,
		contractGovernanceToken:        lo.Must(bindings.NewGovernanceToken(l2.AddressGovernanceTokenProxy, ethereumClient)),
		contractL2CrossDomainMessenger: lo.Must(bindings.NewL2CrossDomainMessenger(l2.AddressL2CrossDomainMessengerProxy, ethereumClient)),
		contractL2StandardBridge:       lo.Must(bindings.NewL2StandardBridge(l2.AddressL2StandardBridgeProxy, ethereumClient)),
		contractL2ToL1MessagePasser:    lo.Must(bindings.NewL2ToL1MessagePasser(l2.AddressL2ToL1MessagePasser, ethereumClient)),
		contractStaking:                lo.Must(l2.NewStaking(contractAddresses.AddressStakingProxy, ethereumClient)),
		contractChips:                  lo.Must(l2.NewChips(contractAddresses.AddressChipsProxy, ethereumClient)),
	}

	return &instance, nil
}
