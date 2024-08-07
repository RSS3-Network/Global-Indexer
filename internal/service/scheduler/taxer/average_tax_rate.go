package taxer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// checkAndSubmitAverageTaxRate checks the average tax rate and submits it to the VSL if necessary.
func (s *Server) checkAndSubmitAverageTaxRate(ctx context.Context) error {
	// Query the submission record of the average tax rate
	submissions, err := s.databaseClient.FindAverageTaxSubmissions(ctx, schema.AverageTaxRateSubmissionQuery{
		Limit: lo.ToPtr(1),
	})
	if err != nil {
		zap.L().Error("find average tax submissions", zap.Error(err))

		return err
	}

	// Query the latest of epoch events
	latestEvent, err := s.databaseClient.FindEpochs(ctx, &schema.FindEpochsQuery{Limit: lo.ToPtr(1)})
	if err != nil {
		zap.L().Error("find epochs", zap.Error(err))

		return err
	}

	// If there is no latest epoch event, do nothing.
	if len(latestEvent) == 0 {
		return nil
	}

	// If the latest submission record is the same as the latest epoch event, do nothing.
	if len(submissions) > 0 && submissions[0].EpochID == latestEvent[0].ID {
		return nil
	}

	// Submit a new average tax rate and save record
	if err = s.submitAverageTaxRate(ctx, latestEvent[0].ID); err != nil {
		zap.L().Error("submit average tax", zap.Error(err))

		return err
	}

	return nil
}

// submitAverageTaxRate submits the average tax rate to the VSL, and saves the submission record.
func (s *Server) submitAverageTaxRate(ctx context.Context, epochID uint64) error {
	// Calculate the average tax.
	averageTax, err := s.calculateAverageTaxRate(ctx)
	if err != nil {
		return fmt.Errorf("calculate average tax rate: %w", err)
	}

	// Submit the average tax to the VSL.
	transactionHash, err := s.invokeSettlementContract(ctx, *averageTax)
	if err != nil {
		return fmt.Errorf("invoke settlement contract: %w", err)
	}

	// Save the submission record to the database.
	submission := &schema.AverageTaxRateSubmission{
		EpochID:         epochID,
		AverageTaxRate:  *averageTax,
		TransactionHash: *transactionHash,
	}
	if err := s.databaseClient.SaveAverageTaxSubmission(ctx, submission); err != nil {
		return fmt.Errorf("save average tax submission: %w", err)
	}

	return nil
}

// calculateAverageTaxRate calculates the average tax rate based on all non-public good nodes.
func (s *Server) calculateAverageTaxRate(ctx context.Context) (*decimal.Decimal, error) {
	var (
		sumTaxRate, nodeCount decimal.Decimal
		cursor                *string
	)

	for {
		// Find nodes from the database.
		nodes, err := s.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
			Cursor: cursor,
			Limit:  lo.ToPtr(100),
		})
		if err != nil {
			zap.L().Error("find Nodes", zap.Error(err))

			return nil, err
		}

		if len(nodes) == 0 {
			break
		}

		// Exclude public good nodes
		var nodeAddresses []common.Address

		for _, node := range nodes {
			if node.IsPublicGood {
				continue
			}

			nodeAddresses = append(nodeAddresses, node.Address)
		}

		// Query the VSL to complement the Node info.
		nodeInfo, err := s.stakingContract.GetNodes(&bind.CallOpts{Context: ctx}, nodeAddresses)
		if err != nil {
			zap.L().Error("get Nodes on the VSL by staking contract", zap.Error(err))

			return nil, err
		}

		// Accumulate the tax and number of all nodes.
		for _, node := range nodeInfo {
			sumTaxRate = sumTaxRate.Add(decimal.NewFromInt(int64(node.TaxRateBasisPoints)))
			nodeCount = nodeCount.Add(decimal.NewFromInt(1))
		}

		cursor = lo.ToPtr(nodes[len(nodes)-1].Address.String())
	}

	// If there are no Nodes, return 0.
	if nodeCount.IsZero() {
		return lo.ToPtr(decimal.Zero), nil
	}

	// Calculate and return the average tax rate.
	return lo.ToPtr(sumTaxRate.Div(nodeCount)), nil
}

// invokeSettlementContract invokes the settlement contract to submit the average tax.
func (s *Server) invokeSettlementContract(ctx context.Context, tax decimal.Decimal) (*common.Hash, error) {
	// Prepare the input data for the settlement contract.
	input, err := s.prepareInputData(tax.BigInt().Uint64())
	if err != nil {
		return nil, fmt.Errorf("prepare input data: %w", err)
	}

	// Send the transaction to the VSL.
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

// sendTransaction sends the transaction to the VSL.
func (s *Server) sendTransaction(ctx context.Context, input []byte) (*common.Hash, error) {
	txCandidate := txmgr.TxCandidate{
		TxData:   input,
		To:       lo.ToPtr(l2.ContractMap[s.chainID.Uint64()].AddressSettlementProxy),
		GasLimit: s.settlerConfig.GasLimit,
		Value:    big.NewInt(0),
	}

	receipt, err := s.txManager.Send(ctx, txCandidate)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		zap.L().Error("received an invalid transaction receipt", zap.String("tx", receipt.TxHash.String()))

		// select {} purposely block the process as it is a critical error and meaningless to continue
		// if panic() is called, the process will be restarted by the supervisor
		// we do not want that as it will be stuck in the same state
		select {}
	}

	return lo.ToPtr(receipt.TxHash), nil
}
