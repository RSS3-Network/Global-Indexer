package utils

import (
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum"
	"math/big"
)

func ParseAmount(rawAmount *big.Int) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).SetInt(rawAmount),
		new(big.Float).SetInt(big.NewInt(ethereum.BILLING_TOKEN_DECIMALS)),
	)
}
