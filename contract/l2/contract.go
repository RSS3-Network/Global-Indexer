package l2

import (
	"github.com/ethereum-optimism/optimism/op-bindings/predeploys"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/Staking.abi --pkg l2 --type Staking --out contract_staking.go
//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/Chips.abi --pkg l2 --type Chips --out contract_chips.go
//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/Settlement.abi --pkg l2 --type Settlement --out contract_settlement.go

var (
	AddressGovernanceTokenProxy        = predeploys.GovernanceTokenAddr                                    // https://scan.testnet.rss3.io/token/0x4200000000000000000000000000000000000042
	AddressL2StandardBridgeProxy       = predeploys.L2StandardBridgeAddr                                   // https://scan.testnet.rss3.io/address/0x4200000000000000000000000000000000000010
	AddressL2CrossDomainMessengerProxy = predeploys.L2CrossDomainMessengerAddr                             // https://scan.testnet.rss3.io/address/0x4200000000000000000000000000000000000007
	AddressL2ToL1MessagePasser         = predeploys.L2ToL1MessagePasserAddr                                // https://scan.testnet.rss3.io/address/0x4200000000000000000000000000000000000007
	AddressStakingProxy                = common.HexToAddress("0x8A312bC2dC1D9549e37f69A4922Da7Df8Bf239db") // https://scan.testnet.rss3.io/address/0x8A312bC2dC1D9549e37f69A4922Da7Df8Bf239db
	AddressChipsProxy                  = common.HexToAddress("0xFB5E5e6e4a90e17641af7EDc86412305E8e44b88") // https://scan.testnet.rss3.io/token/0xFB5E5e6e4a90e17641af7EDc86412305E8e44b88
	AddressSettlementProxy             = common.HexToAddress("0xFaF9d15Ab950220F7072db7B41DEEcA4616B15D9") // https://scan.testnet.rss3.io/address/0xFaF9d15Ab950220F7072db7B41DEEcA4616B15D9
)

var (
	EventHashL2StandardBridgeWithdrawalInitiated = crypto.Keccak256Hash([]byte("WithdrawalInitiated(address,address,address,address,uint256,bytes)"))
	EventHashL2StandardBridgeDepositFinalized    = crypto.Keccak256Hash([]byte("DepositFinalized(address,address,address,address,uint256,bytes)"))

	EventHashL2CrossDomainMessengerSentMessage    = crypto.Keccak256Hash([]byte("SentMessage(address,address,bytes,uint256,uint256)"))
	EventHashL2CrossDomainMessengerRelayedMessage = crypto.Keccak256Hash([]byte("RelayedMessage(bytes32)"))

	EventHashL2ToL1MessagePasserMessagePassed = crypto.Keccak256Hash([]byte("MessagePassed(uint256,address,address,uint256,uint256,bytes,bytes32)"))

	EventHashStakingDeposited         = crypto.Keccak256Hash([]byte("Deposited(address,uint256)"))
	EventHashStakingStaked            = crypto.Keccak256Hash([]byte("Staked(address,address,uint256,uint256,uint256)"))
	EventHashStakingUnstakeRequested  = crypto.Keccak256Hash([]byte("UnstakeRequested(address,address,uint256,uint256,uint256[])"))
	EventHashStakingUnstakeClaimed    = crypto.Keccak256Hash([]byte("UnstakeClaimed(uint256,address,address,uint256)"))
	EventHashStakingWithdrawRequested = crypto.Keccak256Hash([]byte("WithdrawRequested(address,uint256,uint256)"))
	EventHashStakingWithdrawalClaimed = crypto.Keccak256Hash([]byte("WithdrawalClaimed(uint256 indexed requestId)"))

	EventHashChipsTransfer = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

	EventHashRewardDistributed = crypto.Keccak256Hash([]byte("RewardDistributed(uint256,uint256,uint256,address[],uint256[],uint256[],uint256[],uint256[])"))
)

type ChipsTokenMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
