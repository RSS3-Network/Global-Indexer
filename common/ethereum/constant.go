package ethereum

import "github.com/ethereum/go-ethereum/common"

var (
	AddressGenesis = common.HexToAddress("0x0000000000000000000000000000000000000000")
	HashGenesis    = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
)

const (
	BILLING_TOKEN_DECIMALS = 1e18
)
