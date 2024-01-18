package l1

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/L1StandardBridge.abi --pkg l1 --type L1StandardBridge --out contract_l1_standard_bridge.go

var (
	AddressL1CrossDomainMessengerProxy = common.HexToAddress("0x17C635E784B0f098Ab57A39d6dDeA0C786A3AfC1") // https://sepolia.etherscan.io/address/0x17C635E784B0f098Ab57A39d6dDeA0C786A3AfC1
	AddressL1StandardBridgeProxy       = common.HexToAddress("0xc575bd904D16a433624db98D01d5AbD5c92D0F38") // https://sepolia.etherscan.io/address/0xc575bd904D16a433624db98D01d5AbD5c92D0F38
	AddressOptimismPortalProxy         = common.HexToAddress("0xb58f3f17Ef3fAF6cd1C4Fa87b6e15A97B653993E") // https://sepolia.etherscan.io/address/0xb58f3f17Ef3fAF6cd1C4Fa87b6e15A97B653993E
	AddressGovernanceTokenProxy        = common.HexToAddress("0x568F64582A377ea52d0067c4E430B9aE22A60473") // https://sepolia.etherscan.io/address/0x568F64582A377ea52d0067c4E430B9aE22A60473
)

var (
	EventHashL1CrossDomainMessengerSentMessage    = crypto.Keccak256Hash([]byte("SentMessage(address,address,bytes,uint256,uint256)"))
	EventHashL1CrossDomainMessengerRelayedMessage = crypto.Keccak256Hash([]byte("RelayedMessage(bytes32)"))

	EventHashL1StandardBridgeERC20DepositInitiated    = crypto.Keccak256Hash([]byte("ERC20DepositInitiated(address,address,address,address,uint256,bytes)"))
	EventHashL1StandardBridgeERC20WithdrawalFinalized = crypto.Keccak256Hash([]byte("ERC20WithdrawalFinalized(address,address,address,address,uint256,bytes)"))

	EventHashOptimismPortalWithdrawalProven    = crypto.Keccak256Hash([]byte("WithdrawalProven(bytes32,address,address)"))
	EventHashOptimismPortalWithdrawalFinalized = crypto.Keccak256Hash([]byte("WithdrawalFinalized(bytes32,bool)"))
)
