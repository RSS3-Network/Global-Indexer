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
	AddressStakingProxy                = common.HexToAddress("0x6553a9971fe28DA69462613fb012b9c1c302Ce92") // https://scan.testnet.rss3.io/address/0x6553a9971fe28DA69462613fb012b9c1c302Ce92
	AddressChipsProxy                  = common.HexToAddress("0x63144882F6c43d7844e38AcEE55B528a5D883e34") // https://scan.testnet.rss3.io/token/0x63144882F6c43d7844e38AcEE55B528a5D883e34
	AddressSettlementProxy             = common.HexToAddress("0x4D7801d1f3da81A367C0C55Df49601Ee744D03Fe") // https://scan.testnet.rss3.io/address/0x4D7801d1f3da81A367C0C55Df49601Ee744D03Fe
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

var (
	MethodDistributeRewards = "distributeRewards"
)
