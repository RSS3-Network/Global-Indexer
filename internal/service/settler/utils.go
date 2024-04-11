package settler

import (
	"fmt"
	"math/big"

	"github.com/naturalselectionlabs/rss3-global-indexer/common/txmgr"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

// prepareInputData encodes input data for the transaction
func (s *Server) prepareInputData(data schema.SettlementData) ([]byte, error) {
	input, err := txmgr.EncodeInput(l2.SettlementMetaData.ABI, l2.MethodDistributeRewards, data.Epoch, data.NodeAddress, data.OperationRewards, data.RequestCounts, data.IsFinal)
	if err != nil {
		return nil, fmt.Errorf("encode input: %w", err)
	}

	return input, nil
}

func scaleGwei(in *big.Int) {
	in.Mul(in, big.NewInt(1e18))
}
