package ethereum

import "github.com/ethereum/go-ethereum/common"

var (
	AddressGenesis = common.HexToAddress("0x0000000000000000000000000000000000000000")
	HashGenesis    = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
)

const (
	BillingTokenDecimals = 1e18
)
