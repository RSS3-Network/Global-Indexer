package l1

import (
	"context"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rss3-network/global-indexer/contract/l1"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/indexer/internal"
	"github.com/samber/lo"
)

var _ internal.Handler = (*handler)(nil)

type handler struct {
	chainID                        uint64
	finalized                      bool
	ethereumClient                 *ethclient.Client
	contractGovernanceToken        *bindings.GovernanceToken
	contractOptimismPortal         *bindings.OptimismPortal
	contractL1CrossDomainMessenger *bindings.L1CrossDomainMessenger
	contractL1StandardBridge       *bindings.L1StandardBridge
}

func (h *handler) Process(ctx context.Context, block *types.Block, receipts types.Receipts, databaseTransaction database.Client) error {
	if err := databaseTransaction.DeleteBridgeTransactionsByBlockNumber(ctx, h.chainID, block.NumberU64()); err != nil {
		return fmt.Errorf("delete bridge transactions by block number: %w", err)
	}

	if err := databaseTransaction.DeleteBridgeEventsByBlockNumber(ctx, h.chainID, block.NumberU64()); err != nil {
		return fmt.Errorf("delete bridge events by block number: %w", err)
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

		for index, log := range receipt.Logs {
			// Discard all removed logs.
			if log.Removed {
				continue
			}

			// Discard all anonymous logs.
			if len(log.Topics) == 0 {
				continue
			}

			switch log.Address {
			case l1.ContractMap[h.chainID].AddressL1StandardBridgeProxy, l1.ContractMap[h.chainID].AddressOptimismPortalProxy:
				if err := h.indexBridgingLog(ctx, header, block.Transaction(log.TxHash), receipt, log, index, databaseTransaction); err != nil {
					return fmt.Errorf("index bridge log %s %d: %w", log.TxHash, log.Index, err)
				}
			}
		}
	}

	return nil
}

func NewHandler(chainID uint64, ethereumClient *ethclient.Client, finalized bool) (internal.Handler, error) {
	contractAddresses := l1.ContractMap[chainID]
	if contractAddresses == nil {
		return nil, fmt.Errorf("chain id %d is not supported", chainID)
	}

	instance := handler{
		chainID:                        chainID,
		finalized:                      finalized,
		ethereumClient:                 ethereumClient,
		contractGovernanceToken:        lo.Must(bindings.NewGovernanceToken(contractAddresses.AddressGovernanceTokenProxy, ethereumClient)),
		contractOptimismPortal:         lo.Must(bindings.NewOptimismPortal(contractAddresses.AddressOptimismPortalProxy, ethereumClient)),
		contractL1CrossDomainMessenger: lo.Must(bindings.NewL1CrossDomainMessenger(contractAddresses.AddressL1CrossDomainMessengerProxy, ethereumClient)),
		contractL1StandardBridge:       lo.Must(bindings.NewL1StandardBridge(contractAddresses.AddressL1StandardBridgeProxy, ethereumClient)),
	}

	return &instance, nil
}
