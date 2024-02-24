package crossbell

import "github.com/ethereum/go-ethereum/common"

//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/Character.abi --pkg crossbell --type Character --out contract_character.go

var (
	AddressCharacter = common.HexToAddress("0xa6f969045641Cf486a747A2688F3a5A6d43cd0D8")
)
