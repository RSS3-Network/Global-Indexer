package staking

import "github.com/ethereum/go-ethereum/common"

//go:generate go run -mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/Staking.abi --pkg staking --type Staking --out contract_staking.go

var (
	AddressStaking = common.HexToAddress("0x952C35d168eefcCD0C328b0FEf8B7Db51285F039")
)
