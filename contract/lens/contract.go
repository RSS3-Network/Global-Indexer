package lens

import "github.com/ethereum/go-ethereum/common"

//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/LensHandle.abi --pkg lens --type LensHandle --out contract_lens_handle.go

var (
	AddressLensHandle = common.HexToAddress("0xe7E7EaD361f3AaCD73A61A9bD6C10cA17F38E945")
)
