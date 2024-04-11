package settler

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/schema"
)

// prepareInputData encodes input data for the transaction
func (s *Server) prepareInputData(data schema.SettlementData) ([]byte, error) {
	input, err := s.encodeInput(l2.SettlementMetaData.ABI, l2.MethodDistributeRewards, data.Epoch, data.NodeAddress, data.OperationRewards, data.RequestCounts, data.IsFinal)
	if err != nil {
		return nil, fmt.Errorf("encode input: %w", err)
	}

	return input, nil
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

func scaleGwei(in *big.Int) {
	in.Mul(in, big.NewInt(1e18))
}
