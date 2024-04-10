package averagetax

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/txmgr"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// submitAverageTax submits an average tax to the chain, and saves the submission record.
func (s *Server) submitAverageTax(ctx context.Context, epochID uint64) error {
	// Calculate the average tax.
	averageTax, err := s.calculateAverageTax(ctx)
	if err != nil {
		return fmt.Errorf("calculate average tax: %w", err)
	}

	// Submit the average tax to the chain.
	transactionHash, err := s.invokeSettlementContract(ctx, *averageTax)
	if err != nil {
		return fmt.Errorf("invoke settlement contract: %w", err)
	}

	// Save the submission record to the database.
	submission := &schema.AverageTaxSubmission{
		EpochID:         epochID,
		AverageTax:      *averageTax,
		TransactionHash: *transactionHash,
	}
	if err := s.databaseClient.SaveAverageTaxSubmission(ctx, submission); err != nil {
		return fmt.Errorf("save average tax submission: %w", err)
	}

	return nil
}

// calculateAverageTax calculates the average tax for all non-public good nodes.
func (s *Server) calculateAverageTax(ctx context.Context) (*decimal.Decimal, error) {
	var (
		nodesTaxes, nodesAmount decimal.Decimal
		cursor                  *string
	)

	for {
		// Find nodes from the database.
		nodes, err := s.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
			Cursor: cursor,
			Limit:  lo.ToPtr(100),
		})
		if err != nil {
			zap.L().Error("find nodes", zap.Error(err))

			return nil, err
		}

		if len(nodes) == 0 {
			break
		}

		// Filter public good nodes
		var nodeAddresses []common.Address

		for _, node := range nodes {
			if node.IsPublicGood {
				continue
			}

			nodeAddresses = append(nodeAddresses, node.Address)
		}

		// Query the nodes on the chain by staking contract.
		chainNodes, err := s.stakingContract.GetNodes(&bind.CallOpts{Context: ctx}, nodeAddresses)
		if err != nil {
			zap.L().Error("get nodes on the chain by staking contract", zap.Error(err))

			return nil, err
		}

		// Accumulate the tax and number of all nodes.
		for _, node := range chainNodes {
			nodesTaxes = nodesTaxes.Add(decimal.NewFromInt(int64(node.TaxRateBasisPoints)))
			nodesAmount = nodesAmount.Add(decimal.NewFromInt(1))
		}

		cursor = lo.ToPtr(nodes[len(nodes)-1].Address.String())
	}

	// Calculate the average tax.
	if nodesAmount.IsZero() {
		return lo.ToPtr(decimal.Zero), nil
	}

	return lo.ToPtr(nodesTaxes.Div(nodesAmount)), nil
}

// invokeSettlementContract invokes the settlement contract to submit the average tax.
func (s *Server) invokeSettlementContract(ctx context.Context, tax decimal.Decimal) (*common.Hash, error) {
	// Prepare the input data for the settlement contract.
	input, err := s.prepareInputData(tax.BigInt().Uint64())
	if err != nil {
		return nil, fmt.Errorf("prepare input data: %w", err)
	}

	// Send the transaction to the chain.
	transactionHash, err := s.sendTransaction(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("send transaction: %w", err)
	}

	return transactionHash, nil
}

// prepareInputData encodes the input data for the settlement contract.
func (s *Server) prepareInputData(taxRateBasisPoints uint64) ([]byte, error) {
	input, err := txmgr.EncodeInput(l2.SettlementMetaData.ABI, l2.MethodSetTaxRateBasisPoints4PublicPool, taxRateBasisPoints)
	if err != nil {
		return nil, fmt.Errorf("encode input: %w", err)
	}

	return input, nil
}

// sendTransaction sends the transaction to the chain.
func (s *Server) sendTransaction(ctx context.Context, input []byte) (*common.Hash, error) {
	receipt, err := s.txManager.SendTransaction(ctx, input, lo.ToPtr(l2.ContractMap[s.chainID.Uint64()].AddressSettlementProxy), s.settlerConfig.GasLimit)
	if err != nil {
		return nil, fmt.Errorf("send transaction: %w", err)
	}

	return lo.ToPtr(receipt.TxHash), nil
}
