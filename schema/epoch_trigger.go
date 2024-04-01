package schema

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type EpochTrigger struct {
	TransactionHash common.Hash           `json:"transactionHash"`
	EpochID         uint64                `json:"epochID"`
	Data            DistributeRewardsData `json:"data"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
}

type DistributeRewardsData struct {
	Epoch            *big.Int         `json:"epoch"`
	NodeAddress      []common.Address `json:"nodeAddrs"`
	OperationRewards []*big.Int       `json:"operationRewards"`
	RequestCounts    []*big.Int       `json:"requestCounts"`
	IsFinal          bool             `json:"isFinal"`
}
