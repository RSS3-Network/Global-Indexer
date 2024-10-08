package l1

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/L1StandardBridge.abi --pkg l1 --type L1StandardBridge --out contract_l1_standard_bridge.go

const (
	ChainIDMainnet = 1
	ChainIDTestnet = 11155111
)

var ContractMap = map[uint64]*struct {
	AddressL1CrossDomainMessengerProxy common.Address
	AddressL1StandardBridgeProxy       common.Address
	AddressOptimismPortalProxy         common.Address
	AddressGovernanceTokenProxy        common.Address
	AddressUSDCToken                   common.Address
	AddressUSDTToken                   common.Address
	AddressWETHToken                   common.Address
}{
	ChainIDTestnet: {
		AddressL1CrossDomainMessengerProxy: common.HexToAddress("0xf2aAAd7F0ec62f582891F9558dF5F953FEEcC1DA"), // https://sepolia.etherscan.io/address/0xf2aAAd7F0ec62f582891F9558dF5F953FEEcC1DA
		AddressL1StandardBridgeProxy:       common.HexToAddress("0xdDD29bb63B0839FB1cE0eE439Ff027738595D07B"), // https://sepolia.etherscan.io/address/0xdDD29bb63B0839FB1cE0eE439Ff027738595D07B
		AddressOptimismPortalProxy:         common.HexToAddress("0xcBD77E8E1E7F06B25baDe67142cdE82652Da7b57"), // https://sepolia.etherscan.io/address/0xcBD77E8E1E7F06B25baDe67142cdE82652Da7b57
		AddressGovernanceTokenProxy:        common.HexToAddress("0x3Ef1D5be1E2Ce46c583a0c8e511f015706A0ab23"), // https://sepolia.etherscan.io/address/0x3Ef1D5be1E2Ce46c583a0c8e511f015706A0ab23
		AddressUSDCToken:                   common.HexToAddress("0xee32784ef64FDF5F61B2723F7725AdED4F896aD1"), // https://sepolia.etherscan.io/address/0xee32784ef64FDF5F61B2723F7725AdED4F896aD1
		AddressUSDTToken:                   common.HexToAddress("0x259a691d27Da1b81A881c2a7BcD2F4d8c430552C"), // https://sepolia.etherscan.io/address/0x259a691d27Da1b81A881c2a7BcD2F4d8c430552C
		AddressWETHToken:                   common.HexToAddress("0x85dBa244224580D919C6EF1951c5869C604c6C84"), // https://sepolia.etherscan.io/address/0x85dBa244224580D919C6EF1951c5869C604c6C84
	},
	ChainIDMainnet: {
		AddressL1CrossDomainMessengerProxy: common.HexToAddress("0x892CAa506c86C5101f5eC11C6f09589c9dC8A85C"), // https://etherscan.io/address/0x892CAa506c86C5101f5eC11C6f09589c9dC8A85C
		AddressL1StandardBridgeProxy:       common.HexToAddress("0x4cbab69108Aa72151EDa5A3c164eA86845f18438"), // https://etherscan.io/address/0x4cbab69108Aa72151EDa5A3c164eA86845f18438
		AddressOptimismPortalProxy:         common.HexToAddress("0x6A12432491bbbE8d3babf75F759766774C778Db4"), // https://etherscan.io/address/0x6A12432491bbbE8d3babf75F759766774C778Db4
		AddressGovernanceTokenProxy:        common.HexToAddress("0xc98D64DA73a6616c42117b582e832812e7B8D57F"), // https://etherscan.io/address/0xc98D64DA73a6616c42117b582e832812e7B8D57F
		AddressUSDCToken:                   common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), // https://etherscan.io/address/0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
		AddressUSDTToken:                   common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), // https://etherscan.io/address/0xdAC17F958D2ee523a2206206994597C13D831ec7
		AddressWETHToken:                   common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"), // https://etherscan.io/address/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2
	},
}

var (
	EventHashL1CrossDomainMessengerSentMessage    = crypto.Keccak256Hash([]byte("SentMessage(address,address,bytes,uint256,uint256)"))
	EventHashL1CrossDomainMessengerRelayedMessage = crypto.Keccak256Hash([]byte("RelayedMessage(bytes32)"))

	EventHashL1StandardBridgeERC20DepositInitiated    = crypto.Keccak256Hash([]byte("ERC20DepositInitiated(address,address,address,address,uint256,bytes)"))
	EventHashL1StandardBridgeERC20WithdrawalFinalized = crypto.Keccak256Hash([]byte("ERC20WithdrawalFinalized(address,address,address,address,uint256,bytes)"))

	EventHashOptimismPortalWithdrawalProven    = crypto.Keccak256Hash([]byte("WithdrawalProven(bytes32,address,address)"))
	EventHashOptimismPortalWithdrawalFinalized = crypto.Keccak256Hash([]byte("WithdrawalFinalized(bytes32,bool)"))
)
