package l1

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/L1StandardBridge.abi --pkg l1 --type L1StandardBridge --out contract_l1_standard_bridge.go

var (
	AddressL1CrossDomainMessengerProxy = common.HexToAddress("0xb0496292B1A82B284898fA946BE34BdB43B2aDee") // https://sepolia.etherscan.io/address/0xb0496292B1A82B284898fA946BE34BdB43B2aDee
	AddressL1StandardBridgeProxy       = common.HexToAddress("0x30110496b378F5AaB1438E9cf48421f9173841A1") // https://sepolia.etherscan.io/address/0x30110496b378F5AaB1438E9cf48421f9173841A1
	AddressOptimismPortalProxy         = common.HexToAddress("0xB5143A98b600559398A12cC5F4ec8B9F97A7aD63") // https://sepolia.etherscan.io/address/0xB5143A98b600559398A12cC5F4ec8B9F97A7aD63
	AddressGovernanceTokenProxy        = common.HexToAddress("0x3Ef1D5be1E2Ce46c583a0c8e511f015706A0ab23") // https://sepolia.etherscan.io/address/0x3Ef1D5be1E2Ce46c583a0c8e511f015706A0ab23
)

var (
	EventHashL1CrossDomainMessengerSentMessage    = crypto.Keccak256Hash([]byte("SentMessage(address,address,bytes,uint256,uint256)"))
	EventHashL1CrossDomainMessengerRelayedMessage = crypto.Keccak256Hash([]byte("RelayedMessage(bytes32)"))

	EventHashL1StandardBridgeERC20DepositInitiated    = crypto.Keccak256Hash([]byte("ERC20DepositInitiated(address,address,address,address,uint256,bytes)"))
	EventHashL1StandardBridgeERC20WithdrawalFinalized = crypto.Keccak256Hash([]byte("ERC20WithdrawalFinalized(address,address,address,address,uint256,bytes)"))

	EventHashOptimismPortalWithdrawalProven    = crypto.Keccak256Hash([]byte("WithdrawalProven(bytes32,address,address)"))
	EventHashOptimismPortalWithdrawalFinalized = crypto.Keccak256Hash([]byte("WithdrawalFinalized(bytes32,bool)"))
)
