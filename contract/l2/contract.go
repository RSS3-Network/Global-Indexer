package l2

import (
	"math/big"

	"github.com/ethereum-optimism/optimism/op-bindings/predeploys"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/Staking.abi --pkg v1 --type Staking --out ./staking/v1/staking.go
//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/StakingV2.abi --pkg v2 --type Staking --out ./staking/v2/staking.go
//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/Chips.abi --pkg l2 --type Chips --out contract_chips.go
//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/Settlement.abi --pkg l2 --type Settlement --out contract_settlement.go
//go:generate go run --mod=mod github.com/ethereum/go-ethereum/cmd/abigen@v1.13.5 --abi ./abi/NetworkParams.abi --pkg l2 --type NetworkParams --out contract_network_params.go

const (
	ChainIDMainnet = 12553
	ChainIDTestnet = 2331
)

var (
	AddressGovernanceTokenProxy        = predeploys.GovernanceTokenAddr        // https://scan.testnet.rss3.io/token/0x4200000000000000000000000000000000000042
	AddressL2StandardBridgeProxy       = predeploys.L2StandardBridgeAddr       // https://scan.testnet.rss3.io/address/0x4200000000000000000000000000000000000010
	AddressL2CrossDomainMessengerProxy = predeploys.L2CrossDomainMessengerAddr // https://scan.testnet.rss3.io/address/0x4200000000000000000000000000000000000007
	AddressL2ToL1MessagePasser         = predeploys.L2ToL1MessagePasserAddr    // https://scan.testnet.rss3.io/address/0x4200000000000000000000000000000000000007
)

var ContractMap = map[uint64]*struct {
	AddressStakingProxy       common.Address
	AddressChipsProxy         common.Address
	AddressSettlementProxy    common.Address
	AddressNetworkParamsProxy common.Address
}{
	ChainIDTestnet: {
		AddressStakingProxy:       common.HexToAddress("0xb1b209Ee24272C7EE8076764DAa27563c5add9FF"), // https://scan.testnet.rss3.io/address/0xb1b209Ee24272C7EE8076764DAa27563c5add9FF
		AddressChipsProxy:         common.HexToAddress("0x305A3cD2E972ceE48C362ABca02DfA699161edd6"), // https://scan.testnet.rss3.io/token/0x305A3cD2E972ceE48C362ABca02DfA699161edd6
		AddressSettlementProxy:    common.HexToAddress("0xA37a6Ef0c3635824be2b6c87A23F6Df5d0E2ba1b"), // https://scan.testnet.rss3.io/address/0xA37a6Ef0c3635824be2b6c87A23F6Df5d0E2ba1b
		AddressNetworkParamsProxy: common.HexToAddress("0x5d768cAef86d3DA8eC6009eE4B3d9b7Fe26A43CB"), // https://scan.testnet.rss3.io/address/0x5d768cAef86d3DA8eC6009eE4B3d9b7Fe26A43CB
	},
	ChainIDMainnet: {
		AddressStakingProxy:       common.HexToAddress("0x28F14d917fddbA0c1f2923C406952478DfDA5578"), // https://scan.rss3.io/address/0x28F14d917fddbA0c1f2923C406952478DfDA5578
		AddressChipsProxy:         common.HexToAddress("0x849f8F55078dCc69dD857b58Cc04631EBA54E4DE"), // https://scan.rss3.io/token/0x849f8F55078dCc69dD857b58Cc04631EBA54E4DE
		AddressSettlementProxy:    common.HexToAddress("0x0cE3159BF19F3C55B648D04E8f0Ae1Ae118D2A0B"), // https://scan.rss3.io/address/0x0cE3159BF19F3C55B648D04E8f0Ae1Ae118D2A0B
		AddressNetworkParamsProxy: common.HexToAddress("0x15176Aabdc4836c38947a67313d209204051C502"), // https://scan.rss3.io/address/0x15176Aabdc4836c38947a67313d209204051C502
	},
}

var GenesisEpochMap = map[uint64]int64{
	ChainIDMainnet: 1710208800, // 2024-03-12 02:00:00 UTC
}

var (
	EventHashL2StandardBridgeWithdrawalInitiated = crypto.Keccak256Hash([]byte("WithdrawalInitiated(address,address,address,address,uint256,bytes)"))
	EventHashL2StandardBridgeDepositFinalized    = crypto.Keccak256Hash([]byte("DepositFinalized(address,address,address,address,uint256,bytes)"))

	EventHashL2CrossDomainMessengerSentMessage    = crypto.Keccak256Hash([]byte("SentMessage(address,address,bytes,uint256,uint256)"))
	EventHashL2CrossDomainMessengerRelayedMessage = crypto.Keccak256Hash([]byte("RelayedMessage(bytes32)"))

	EventHashL2ToL1MessagePasserMessagePassed = crypto.Keccak256Hash([]byte("MessagePassed(uint256,address,address,uint256,uint256,bytes,bytes32)"))

	EventHashStakingV1Deposited         = crypto.Keccak256Hash([]byte("Deposited(address,uint256)"))
	EventHashStakingV1Staked            = crypto.Keccak256Hash([]byte("Staked(address,address,uint256,uint256,uint256)"))
	EventHashStakingV1UnstakeRequested  = crypto.Keccak256Hash([]byte("UnstakeRequested(address,address,uint256,uint256,uint256[])"))
	EventHashStakingV1UnstakeClaimed    = crypto.Keccak256Hash([]byte("UnstakeClaimed(uint256,address,address,uint256)"))
	EventHashStakingV1WithdrawRequested = crypto.Keccak256Hash([]byte("WithdrawRequested(address,uint256,uint256)"))
	EventHashStakingV1WithdrawalClaimed = crypto.Keccak256Hash([]byte("WithdrawalClaimed(uint256)"))
	EventHashStakingV1RewardDistributed = crypto.Keccak256Hash([]byte("RewardDistributed(uint256,uint256,uint256,address[],uint256[],uint256[],uint256[],uint256[])"))
	EventHashStakingV1NodeCreated       = crypto.Keccak256Hash([]byte("NodeCreated(uint256,address,string,string,uint64,bool,bool)"))
	EventHashStakingV1NodeUpdated       = crypto.Keccak256Hash([]byte("NodeUpdated(address,string,string)"))

	EventHashStakingV2ChipsMerged       = crypto.Keccak256Hash([]byte("ChipsMerged(address,address,uint256,uint256[])"))
	EventHashStakingV2WithdrawalClaimed = crypto.Keccak256Hash([]byte("WithdrawalClaimed(uint256,address,uint256)"))

	EventHashChipsTransfer = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
)

var (
	MethodDistributeRewards                = "distributeRewards"
	MethodSetTaxRateBasisPoints4PublicPool = "setTaxRateBasisPoints4PublicPool"
)

type ChipsTokenMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func IsStakingV2Deployed(chainID *big.Int, blockNumber *big.Int, transactionIndex uint) bool {
	switch chainID.Uint64() {
	case ChainIDMainnet:
		// https://scan.rss3.io/tx/0x0360cc8c8c91063412551f2c6e97dbf5c0b4a352a77971f9ecaf75a128dcd2d2
		return blockNumber.Uint64() >= 6023345 && transactionIndex >= 0 // nolint:staticcheck // False positive.
	case ChainIDTestnet:
		// https://scan.testnet.rss3.io/tx/0xdfe5e81939f2183cb99076b0dd860b95718ceb42d163267e94c35f079761db93
		return blockNumber.Uint64() >= 6516895 && transactionIndex >= 0 // nolint:staticcheck // False positive.
	default:
		return false
	}
}
