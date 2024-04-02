package settler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/txmgr"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// submitEpochProof submits proof of this epoch on chain
// which calculates the Operation Rewards for the Nodes
// formats the data and invokes the contract
// a retry logic is implemented to handle possible failures
func (s *Server) submitEpochProof(ctx context.Context, epoch uint64) error {
	if err := s.mutex.Lock(); err != nil {
		zap.L().Error("lock error", zap.String("key", s.mutex.Name()), zap.Error(err))

		return nil
	}

	defer func() {
		if _, err := s.mutex.Unlock(); err != nil {
			zap.L().Error("release lock error", zap.String("key", s.mutex.Name()), zap.Error(err))
		}
	}()

	var cursor *string

	for {
		msg := "construct Settlement data"
		// Construct transactionData as required by the Settlement contract
		transactionData, err := s.constructSettlementData(ctx, epoch, cursor)
		if err != nil {
			zap.L().Error(msg, zap.Error(err))

			return fmt.Errorf("%s: %w", msg, err)
		}

		// Finish processing when conditions are met
		if len(transactionData.NodeAddress) == 0 && cursor != nil {
			zap.L().Info("finished processing transactionData.")

			break
		}

		zap.L().Info(msg, zap.Any("transactionData", transactionData))

		// Invoke the Settlement contract
		if err = retry.Do(func() error {
			return s.invokeSettlementContract(ctx, *transactionData)
		}, retry.Delay(time.Second), retry.Attempts(5)); err != nil {
			zap.L().Error("retry submitEpochProof invokeSettlementContract", zap.Error(err))

			return err
		}

		if len(transactionData.NodeAddress) > 0 {
			cursor = lo.ToPtr(transactionData.NodeAddress[len(transactionData.NodeAddress)-1].String())
		}
	}

	zap.L().Info("Epoch Proof submitted successfully", zap.Uint64("settler", epoch))

	return nil
}

// constructSettlementData constructs Settlement data as required by the Settlement contract
func (s *Server) constructSettlementData(ctx context.Context, epoch uint64, cursor *string) (*schema.SettlementData, error) {
	batchSize := s.settlerConfig.BatchSize

	// Find qualified Nodes from the database
	nodes, err := s.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
		Status: lo.ToPtr(schema.NodeStatusOnline),
		Cursor: cursor,
		Limit:  lo.ToPtr(batchSize + 1),
	})
	if err != nil {
		// No qualified Nodes found in the database
		if errors.Is(err, database.ErrorRowNotFound) {
			return nil, nil
		}

		zap.L().Error("No qualified Nodes found", zap.Error(err), zap.Any("cursor", cursor))

		return nil, err
	}

	// isFinal is true if it's the last batch of Nodes
	isFinal := len(nodes) <= batchSize
	if !isFinal {
		nodes = nodes[:batchSize]
	}

	// nodeAddresses is a slice of Node addresses
	nodeAddresses := make([]common.Address, 0, len(nodes))
	for _, node := range nodes {
		nodeAddresses = append(nodeAddresses, node.Address)
	}

	// Get the number of stakers in the last 5 epochs for all nodes.
	recentStakers, err := s.databaseClient.FindStakerCountRecentEpochs(ctx, s.specialRewards.EpochLimit)
	if err != nil {
		return nil, fmt.Errorf("find recent stakers: %w", err)
	}

	// Calculate the operation rewards for the Nodes
	operationRewards, err := calculateOperationRewards(nodes, recentStakers, s.specialRewards)
	if err != nil {
		return nil, err
	}

	// Calculate the operation rewards for the Nodes
	requestCounts := prepareRequestCounts(nodeAddresses)

	return &schema.SettlementData{
		Epoch:            big.NewInt(int64(epoch)),
		NodeAddress:      nodeAddresses,
		OperationRewards: operationRewards,
		RequestCounts:    requestCounts,
		IsFinal:          isFinal,
	}, nil
}

// calculateOperationRewards calculates the Operation Rewards for all Nodes
// For Alpha, there is no Operation Rewards, but a Special Rewards is calculated
// TODO: Implement the actual calculation logic
func calculateOperationRewards(nodes []*schema.Node, recentStackers map[common.Address]uint64, specialRewards *config.SpecialRewards) ([]*big.Int, error) {
	operationRewards, err := calculateAlphaSpecialRewards(nodes, recentStackers, specialRewards)
	if err != nil {
		return nil, err
	}

	// For Alpha, set the rewards to 0
	//for i := range operationRewards {
	//	operationRewards[i] = big.NewInt(0)
	//}

	return operationRewards, nil
}

// prepareRequestCounts prepares the Request Counts for all Nodes
// For Alpha, there is no actual calculation logic, the counts are set to 0
// TODO: Implement the actual logic to retrieve the counts from the database
func prepareRequestCounts(nodes []common.Address) []*big.Int {
	slice := make([]*big.Int, len(nodes))

	// For Alpha, set the counts to 0
	for i := range slice {
		slice[i] = big.NewInt(0)
	}

	return slice
}

// calculateAlphaSpecialRewards calculates the distribution of the Special Rewards used to replace the Operation Rewards
// the Special Rewards are used to incentivize staking in smaller Nodes
// currently, the amount is set to 30,000,000 / 486.6666666666667 * 0.2 ~= 12328
func calculateAlphaSpecialRewards(nodes []*schema.Node, recentStackers map[common.Address]uint64, specialRewards *config.SpecialRewards) ([]*big.Int, error) {
	var (
		totalEffectiveStakers uint64
		maxPoolSize           uint64
		scores                []float64
		rewards               []*big.Int
		totalScore            float64
	)

	rewards = make([]*big.Int, len(nodes))
	scores = make([]float64, len(nodes))

	for _, node := range nodes {
		poolSize, err := strconv.ParseUint(node.StakingPoolTokens, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse staking pool tokens: %s", node.StakingPoolTokens)
		}

		// calculate the number of effective stakers
		// which is the number of stakers for poolSize <= cliffPoint
		if poolSize <= specialRewards.CliffPoint {
			// If the node has no recent stakers, the map will return 0.
			totalEffectiveStakers += recentStackers[node.Address]
		}

		// store the max pool size
		if maxPoolSize < poolSize {
			maxPoolSize = poolSize
		}
	}

	for i, node := range nodes {
		stakers := recentStackers[node.Address]

		// no stakers, no rewards
		if stakers == 0 {
			rewards[i] = big.NewInt(0)
			continue
		}

		poolSize, err := strconv.ParseUint(node.StakingPoolTokens, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse staking pool tokens: %s", node.StakingPoolTokens)
		}

		// apply the Gini Coefficient
		// giniCoefficient balances rewards in favor of Nodes with smaller P_s
		// higher values will favor smaller P_s
		score := applyGiniCoefficient(poolSize, specialRewards.GiniCoefficient)

		// cliffPoint, P_s beyond this point will have the rewards reduced
		// higher values will favor smaller P_s
		if poolSize > specialRewards.CliffPoint {
			// apply the Cliff Factor only when poolSize > cliffPoint
			applyCliffFactor(poolSize, maxPoolSize, &score, specialRewards.CliffFactor)
		}

		if totalEffectiveStakers > 0 {
			// apply the Staker Factor
			applyStakerFactor(stakers, totalEffectiveStakers, specialRewards.StakerFactor, &score)
		}

		if score < 0 || score >= 1 {
			return nil, fmt.Errorf("AlphaSpecialRewards: invalid score: %f", score)
		}

		totalScore += score
		scores[i] = score
	}

	// final loop to calculate the rewards
	for i := range nodes {
		// Deliberately do it step-by-step to make it easier to understand
		// Truncate the floating point number to an integer
		reward := math.Trunc(scores[i] / totalScore * specialRewards.Rewards)

		// Represent the float64 as big.Float
		rewardBigFloat := new(big.Float).SetFloat64(reward)

		// Create a big.Float representation of 10^18 for scaling
		scale := new(big.Float).SetInt(big.NewInt(1e18))

		// Multiply the original number by the scaling factor
		scaledF := new(big.Float).Mul(rewardBigFloat, scale)

		// Convert the scaled big.Float to big.Int
		rewardFinal := new(big.Int)
		scaledF.Int(rewardFinal)

		rewards[i] = rewardFinal
	}

	return rewards, nil
}

// applyGiniCoefficient applies the Gini Coefficient to the score
func applyGiniCoefficient(poolSize uint64, giniCoefficient float64) float64 {
	// Perform calculation: score = 1 / (1 + giniCoefficient * poolSize)
	score := 1 / (1 + giniCoefficient*float64(poolSize))

	return score
}

// applyCliffFactor applies the Cliff Factor to the score
func applyCliffFactor(poolSize uint64, maxPoolSize uint64, score *float64, cliffFactor float64) {
	// Perform calculation: score *= cliffFactor ** poolSize / maxPoolSize
	*score *= math.Pow(cliffFactor, float64(poolSize)/float64(maxPoolSize))
}

// applyStakerFactor applies the Staker Factor to the score
func applyStakerFactor(stakers uint64, totalEffectiveStakers uint64, stakerFactor float64, score *float64) {
	// Perform calculation: score += (score * stakers * staker_factor) / total_stakers
	*score += (*score * float64(stakers) * stakerFactor) / float64(totalEffectiveStakers)
}

// invokeSettlementContract invokes the Settlement contract with prepared data
// and saves the Settlement to the database
func (s *Server) invokeSettlementContract(ctx context.Context, data schema.SettlementData) error {
	input, err := s.prepareInputData(data)
	if err != nil {
		return err
	}

	receipt, err := s.sendTransaction(ctx, input)
	if err != nil {
		return err
	}

	// Save the Settlement to the database, as the reference point for the next Epoch
	if err := s.saveSettlement(ctx, receipt, data); err != nil {
		return err
	}

	zap.L().Info("Settlement contracted invoked successfully", zap.String("tx", receipt.TxHash.String()), zap.Any("data", data))

	return nil
}

// prepareInputData encodes input data for the transaction
func (s *Server) prepareInputData(data schema.SettlementData) ([]byte, error) {
	input, err := s.encodeInput(l2.SettlementMetaData.ABI, l2.MethodDistributeRewards, data.Epoch, data.NodeAddress, data.OperationRewards, data.IsFinal)
	if err != nil {
		return nil, fmt.Errorf("encode input: %w", err)
	}

	return input, nil
}

// sendTransaction sends the transaction and returns the receipt if successful
func (s *Server) sendTransaction(ctx context.Context, input []byte) (*types.Receipt, error) {
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

	// return the receipt if the transaction is successful
	return receipt, nil
}

// saveSettlement saves the Settlement data to the database
func (s *Server) saveSettlement(ctx context.Context, receipt *types.Receipt, data schema.SettlementData) error {
	if err := s.databaseClient.SaveEpochTrigger(ctx, &schema.EpochTrigger{
		TransactionHash: receipt.TxHash,
		EpochID:         data.Epoch.Uint64(),
		Data:            data,
	}); err != nil {
		return fmt.Errorf("save settler submitEpochProof: %w", err)
	}

	return nil
}

// encodeInput encodes the input data according to the contract ABI
func (s *Server) encodeInput(contractABI, methodName string, args ...interface{}) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}

	encodedArgs, err := parsedABI.Pack(methodName, args...)
	if err != nil {
		return nil, err
	}

	return encodedArgs, nil
}
