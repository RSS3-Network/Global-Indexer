package l1

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/L1StandardBridge.abi --pkg l1 --type L1StandardBridge --out contract_l1_standard_bridge.go

var (
	AddressL1CrossDomainMessengerProxy = common.HexToAddress("0xf2aAAd7F0ec62f582891F9558dF5F953FEEcC1DA") // https://sepolia.etherscan.io/address/0xf2aAAd7F0ec62f582891F9558dF5F953FEEcC1DA
	AddressL1StandardBridgeProxy       = common.HexToAddress("0xdDD29bb63B0839FB1cE0eE439Ff027738595D07B") // https://sepolia.etherscan.io/address/0xdDD29bb63B0839FB1cE0eE439Ff027738595D07B
	AddressOptimismPortalProxy         = common.HexToAddress("0xcBD77E8E1E7F06B25baDe67142cdE82652Da7b57") // https://sepolia.etherscan.io/address/0xcBD77E8E1E7F06B25baDe67142cdE82652Da7b57
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
