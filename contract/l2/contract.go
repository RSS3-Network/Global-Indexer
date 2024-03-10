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
	AddressGovernanceTokenProxy        = predeploys.GovernanceTokenAddr        // https://scan.testnet.rss3.io/token/0x4200000000000000000000000000000000000042
	AddressL2StandardBridgeProxy       = predeploys.L2StandardBridgeAddr       // https://scan.testnet.rss3.io/address/0x4200000000000000000000000000000000000010
	AddressL2CrossDomainMessengerProxy = predeploys.L2CrossDomainMessengerAddr // https://scan.testnet.rss3.io/address/0x4200000000000000000000000000000000000007
	AddressL2ToL1MessagePasser         = predeploys.L2ToL1MessagePasserAddr    // https://scan.testnet.rss3.io/address/0x4200000000000000000000000000000000000007
)

var ContractMap = map[uint64]*struct {
	AddressStakingProxy    common.Address
	AddressChipsProxy      common.Address
	AddressSettlementProxy common.Address
}{
	2331: {
		AddressStakingProxy:    common.HexToAddress("0xb1b209Ee24272C7EE8076764DAa27563c5add9FF"), // https://scan.testnet.rss3.io/address/0xb1b209Ee24272C7EE8076764DAa27563c5add9FF
		AddressChipsProxy:      common.HexToAddress("0x305A3cD2E972ceE48C362ABca02DfA699161edd6"), // https://scan.testnet.rss3.io/token/0x305A3cD2E972ceE48C362ABca02DfA699161edd6
		AddressSettlementProxy: common.HexToAddress("0xA37a6Ef0c3635824be2b6c87A23F6Df5d0E2ba1b"), // https://scan.testnet.rss3.io/address/0xA37a6Ef0c3635824be2b6c87A23F6Df5d0E2ba1b
	},
	12553: {
		AddressStakingProxy:    common.HexToAddress("0x28F14d917fddbA0c1f2923C406952478DfDA5578"), // https://scan.rss3.io/address/0x28F14d917fddbA0c1f2923C406952478DfDA5578
		AddressChipsProxy:      common.HexToAddress("0x849f8F55078dCc69dD857b58Cc04631EBA54E4DE"), // https://scan.rss3.io/token/0x849f8F55078dCc69dD857b58Cc04631EBA54E4DE
		AddressSettlementProxy: common.HexToAddress("0x0cE3159BF19F3C55B648D04E8f0Ae1Ae118D2A0B"), // https://scan.rss3.io/address/0x0cE3159BF19F3C55B648D04E8f0Ae1Ae118D2A0B
	},
}

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
	EventHashStakingRewardDistributed = crypto.Keccak256Hash([]byte("RewardDistributed(uint256,uint256,uint256,address[],uint256[],uint256[],uint256[])"))
	EventHashStakingNodeCreated       = crypto.Keccak256Hash([]byte("NodeCreated(uint256,address,string,string,uint64,bool,bool)"))

	EventHashChipsTransfer = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
)

var (
	MethodDistributeRewards = "distributeRewards"
)

type ChipsTokenMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
