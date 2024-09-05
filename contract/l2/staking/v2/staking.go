// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package v2

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// Demotion is an auto generated low-level Go binding around an user-defined struct.
type Demotion struct {
	DemotionId *big.Int
	NodeAddr   common.Address
	Epoch      *big.Int
	Reason     string
	Reporter   common.Address
}

// Node is an auto generated low-level Go binding around an user-defined struct.
type Node struct {
	NodeId                     *big.Int
	Account                    common.Address
	TaxRateBasisPoints         uint64
	PublicGood                 bool
	Alpha                      bool
	Name                       string
	Description                string
	OperationPoolTokens        *big.Int
	StakingPoolTokens          *big.Int
	TotalShares                *big.Int
	SlashedOperationPoolTokens *big.Int
	SlashedStakingPoolTokens   *big.Int
	ExitTime                   *big.Int
	Status                     uint8
}

// UnstakeRequest is an auto generated low-level Go binding around an user-defined struct.
type UnstakeRequest struct {
	Owner         common.Address
	NodeAddr      common.Address
	Timestamp     *big.Int
	UnstakeAmount *big.Int
}

// WithdrawalRequest is an auto generated low-level Go binding around an user-defined struct.
type WithdrawalRequest struct {
	Owner     common.Address
	Timestamp *big.Int
	Amount    *big.Int
}

// StakingMetaData contains all meta data concerning the Staking contract.
var StakingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"treasury\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakeUnbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"depositUnbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"paymentProcessor\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEPOSIT_UNBONDING_PERIOD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ORACLE_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAUSE_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAYMENT_PROCESSOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STAKE_UNBONDING_PERIOD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TREASURY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"chipsContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimUnstake\",\"inputs\":[{\"name\":\"requestIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimWithdrawal\",\"inputs\":[{\"name\":\"requestIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitSlashing\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createNode\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"distributeRewards\",\"inputs\":[{\"name\":\"epochInfo\",\"type\":\"uint256[3]\",\"internalType\":\"uint256[3]\"},{\"name\":\"nodeAddrs\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"operationRewards\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"stakingRewards\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"requestCounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"publicPoolRewards\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"exit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getChipInfo\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDemotions\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"demotions\",\"type\":\"tuple[]\",\"internalType\":\"structDemotion[]\",\"components\":[{\"name\":\"demotionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"reporter\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNode\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNode\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"alpha\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"operationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashedOperationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashedStakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exitTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumNodeStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeAvatar\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodes\",\"inputs\":[{\"name\":\"nodeAddrs\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structNode[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"alpha\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"operationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashedOperationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashedStakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exitTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumNodeStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingUnstake\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structUnstakeRequest\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingWithdrawal\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structWithdrawalRequest\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint40\",\"internalType\":\"uint40\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPoolInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"totalOperationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalStakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalSlashingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublicPool\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNode\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"alpha\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"operationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashedOperationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashedStakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exitTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumNodeStatus\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleMember\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleMemberCount\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"chips\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pauseAccount\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracleAccount\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isAlphaPhase_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isAlphaPhase\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSettlementPhase\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mergeChips\",\"inputs\":[{\"name\":\"chipIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"newTokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"multicall\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"results\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"online\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestUnstake\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"chipIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestWithdrawal\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeDemotions\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"demotionIdsToRevoke\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNodeStatus\",\"inputs\":[{\"name\":\"nodeAddrs\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"status\",\"type\":\"uint8[]\",\"internalType\":\"enumNodeStatus[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSettlementPhase\",\"inputs\":[{\"name\":\"enabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTaxRateBasisPoints4Node\",\"inputs\":[{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTaxRateBasisPoints4PublicPool\",\"inputs\":[{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeToPublicPool\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"submitDemotions\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nodeAddrs\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"reasons\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"reporters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateNode\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw2Treasury\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PublicGoodRewardDistributed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"publicPoolRewards\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"publicPoolTax\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardDistributed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"nodeAddrs\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"operationRewards\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"stakingRewards\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"taxCollected\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"requestCounts\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExcessWithdrawalAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidArrayLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeNotExists\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NodeNotPublicGood\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SettlementPhase\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakeToPublicGoodNode\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"WithdrawalAmountExceedsOperationPoolTokens\",\"inputs\":[]}]",
}

// StakingABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingMetaData.ABI instead.
var StakingABI = StakingMetaData.ABI

// Staking is an auto generated Go binding around an Ethereum contract.
type Staking struct {
	StakingCaller     // Read-only binding to the contract
	StakingTransactor // Write-only binding to the contract
	StakingFilterer   // Log filterer for contract events
}

// StakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingSession struct {
	Contract     *Staking          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingCallerSession struct {
	Contract *StakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTransactorSession struct {
	Contract     *StakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingRaw struct {
	Contract *Staking // Generic contract binding to access the raw methods on
}

// StakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingCallerRaw struct {
	Contract *StakingCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTransactorRaw struct {
	Contract *StakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStaking creates a new instance of Staking, bound to a specific deployed contract.
func NewStaking(address common.Address, backend bind.ContractBackend) (*Staking, error) {
	contract, err := bindStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// NewStakingCaller creates a new read-only instance of Staking, bound to a specific deployed contract.
func NewStakingCaller(address common.Address, caller bind.ContractCaller) (*StakingCaller, error) {
	contract, err := bindStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingCaller{contract: contract}, nil
}

// NewStakingTransactor creates a new write-only instance of Staking, bound to a specific deployed contract.
func NewStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTransactor, error) {
	contract, err := bindStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTransactor{contract: contract}, nil
}

// NewStakingFilterer creates a new log filterer instance of Staking, bound to a specific deployed contract.
func NewStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingFilterer, error) {
	contract, err := bindStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingFilterer{contract: contract}, nil
}

// bindStaking binds a generic wrapper to an already deployed contract.
func bindStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.StakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Staking *StakingCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Staking *StakingSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Staking.Contract.DEFAULTADMINROLE(&_Staking.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Staking *StakingCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Staking.Contract.DEFAULTADMINROLE(&_Staking.CallOpts)
}

// DEPOSITUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x6bdc11d5.
//
// Solidity: function DEPOSIT_UNBONDING_PERIOD() view returns(uint256)
func (_Staking *StakingCaller) DEPOSITUNBONDINGPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "DEPOSIT_UNBONDING_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEPOSITUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x6bdc11d5.
//
// Solidity: function DEPOSIT_UNBONDING_PERIOD() view returns(uint256)
func (_Staking *StakingSession) DEPOSITUNBONDINGPERIOD() (*big.Int, error) {
	return _Staking.Contract.DEPOSITUNBONDINGPERIOD(&_Staking.CallOpts)
}

// DEPOSITUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x6bdc11d5.
//
// Solidity: function DEPOSIT_UNBONDING_PERIOD() view returns(uint256)
func (_Staking *StakingCallerSession) DEPOSITUNBONDINGPERIOD() (*big.Int, error) {
	return _Staking.Contract.DEPOSITUNBONDINGPERIOD(&_Staking.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Staking *StakingCaller) ORACLEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "ORACLE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Staking *StakingSession) ORACLEROLE() ([32]byte, error) {
	return _Staking.Contract.ORACLEROLE(&_Staking.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Staking *StakingCallerSession) ORACLEROLE() ([32]byte, error) {
	return _Staking.Contract.ORACLEROLE(&_Staking.CallOpts)
}

// PAUSEROLE is a free data retrieval call binding the contract method 0x389ed267.
//
// Solidity: function PAUSE_ROLE() view returns(bytes32)
func (_Staking *StakingCaller) PAUSEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "PAUSE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSEROLE is a free data retrieval call binding the contract method 0x389ed267.
//
// Solidity: function PAUSE_ROLE() view returns(bytes32)
func (_Staking *StakingSession) PAUSEROLE() ([32]byte, error) {
	return _Staking.Contract.PAUSEROLE(&_Staking.CallOpts)
}

// PAUSEROLE is a free data retrieval call binding the contract method 0x389ed267.
//
// Solidity: function PAUSE_ROLE() view returns(bytes32)
func (_Staking *StakingCallerSession) PAUSEROLE() ([32]byte, error) {
	return _Staking.Contract.PAUSEROLE(&_Staking.CallOpts)
}

// PAYMENTPROCESSOR is a free data retrieval call binding the contract method 0xaf2ae425.
//
// Solidity: function PAYMENT_PROCESSOR() view returns(address)
func (_Staking *StakingCaller) PAYMENTPROCESSOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "PAYMENT_PROCESSOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PAYMENTPROCESSOR is a free data retrieval call binding the contract method 0xaf2ae425.
//
// Solidity: function PAYMENT_PROCESSOR() view returns(address)
func (_Staking *StakingSession) PAYMENTPROCESSOR() (common.Address, error) {
	return _Staking.Contract.PAYMENTPROCESSOR(&_Staking.CallOpts)
}

// PAYMENTPROCESSOR is a free data retrieval call binding the contract method 0xaf2ae425.
//
// Solidity: function PAYMENT_PROCESSOR() view returns(address)
func (_Staking *StakingCallerSession) PAYMENTPROCESSOR() (common.Address, error) {
	return _Staking.Contract.PAYMENTPROCESSOR(&_Staking.CallOpts)
}

// STAKEUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x2606a44a.
//
// Solidity: function STAKE_UNBONDING_PERIOD() view returns(uint256)
func (_Staking *StakingCaller) STAKEUNBONDINGPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "STAKE_UNBONDING_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// STAKEUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x2606a44a.
//
// Solidity: function STAKE_UNBONDING_PERIOD() view returns(uint256)
func (_Staking *StakingSession) STAKEUNBONDINGPERIOD() (*big.Int, error) {
	return _Staking.Contract.STAKEUNBONDINGPERIOD(&_Staking.CallOpts)
}

// STAKEUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x2606a44a.
//
// Solidity: function STAKE_UNBONDING_PERIOD() view returns(uint256)
func (_Staking *StakingCallerSession) STAKEUNBONDINGPERIOD() (*big.Int, error) {
	return _Staking.Contract.STAKEUNBONDINGPERIOD(&_Staking.CallOpts)
}

// TREASURY is a free data retrieval call binding the contract method 0x2d2c5565.
//
// Solidity: function TREASURY() view returns(address)
func (_Staking *StakingCaller) TREASURY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "TREASURY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TREASURY is a free data retrieval call binding the contract method 0x2d2c5565.
//
// Solidity: function TREASURY() view returns(address)
func (_Staking *StakingSession) TREASURY() (common.Address, error) {
	return _Staking.Contract.TREASURY(&_Staking.CallOpts)
}

// TREASURY is a free data retrieval call binding the contract method 0x2d2c5565.
//
// Solidity: function TREASURY() view returns(address)
func (_Staking *StakingCallerSession) TREASURY() (common.Address, error) {
	return _Staking.Contract.TREASURY(&_Staking.CallOpts)
}

// ChipsContract is a free data retrieval call binding the contract method 0xd13b19a3.
//
// Solidity: function chipsContract() view returns(address)
func (_Staking *StakingCaller) ChipsContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "chipsContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ChipsContract is a free data retrieval call binding the contract method 0xd13b19a3.
//
// Solidity: function chipsContract() view returns(address)
func (_Staking *StakingSession) ChipsContract() (common.Address, error) {
	return _Staking.Contract.ChipsContract(&_Staking.CallOpts)
}

// ChipsContract is a free data retrieval call binding the contract method 0xd13b19a3.
//
// Solidity: function chipsContract() view returns(address)
func (_Staking *StakingCallerSession) ChipsContract() (common.Address, error) {
	return _Staking.Contract.ChipsContract(&_Staking.CallOpts)
}

// GetChipInfo is a free data retrieval call binding the contract method 0x391f9418.
//
// Solidity: function getChipInfo(uint256 tokenId) view returns(address nodeAddr, uint256 tokens, uint256 shares)
func (_Staking *StakingCaller) GetChipInfo(opts *bind.CallOpts, tokenId *big.Int) (struct {
	NodeAddr common.Address
	Tokens   *big.Int
	Shares   *big.Int
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getChipInfo", tokenId)

	outstruct := new(struct {
		NodeAddr common.Address
		Tokens   *big.Int
		Shares   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NodeAddr = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Tokens = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Shares = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetChipInfo is a free data retrieval call binding the contract method 0x391f9418.
//
// Solidity: function getChipInfo(uint256 tokenId) view returns(address nodeAddr, uint256 tokens, uint256 shares)
func (_Staking *StakingSession) GetChipInfo(tokenId *big.Int) (struct {
	NodeAddr common.Address
	Tokens   *big.Int
	Shares   *big.Int
}, error) {
	return _Staking.Contract.GetChipInfo(&_Staking.CallOpts, tokenId)
}

// GetChipInfo is a free data retrieval call binding the contract method 0x391f9418.
//
// Solidity: function getChipInfo(uint256 tokenId) view returns(address nodeAddr, uint256 tokens, uint256 shares)
func (_Staking *StakingCallerSession) GetChipInfo(tokenId *big.Int) (struct {
	NodeAddr common.Address
	Tokens   *big.Int
	Shares   *big.Int
}, error) {
	return _Staking.Contract.GetChipInfo(&_Staking.CallOpts, tokenId)
}

// GetDemotions is a free data retrieval call binding the contract method 0x7378fa9e.
//
// Solidity: function getDemotions(address nodeAddr, uint256 epoch) view returns((uint256,address,uint256,string,address)[] demotions)
func (_Staking *StakingCaller) GetDemotions(opts *bind.CallOpts, nodeAddr common.Address, epoch *big.Int) ([]Demotion, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDemotions", nodeAddr, epoch)

	if err != nil {
		return *new([]Demotion), err
	}

	out0 := *abi.ConvertType(out[0], new([]Demotion)).(*[]Demotion)

	return out0, err

}

// GetDemotions is a free data retrieval call binding the contract method 0x7378fa9e.
//
// Solidity: function getDemotions(address nodeAddr, uint256 epoch) view returns((uint256,address,uint256,string,address)[] demotions)
func (_Staking *StakingSession) GetDemotions(nodeAddr common.Address, epoch *big.Int) ([]Demotion, error) {
	return _Staking.Contract.GetDemotions(&_Staking.CallOpts, nodeAddr, epoch)
}

// GetDemotions is a free data retrieval call binding the contract method 0x7378fa9e.
//
// Solidity: function getDemotions(address nodeAddr, uint256 epoch) view returns((uint256,address,uint256,string,address)[] demotions)
func (_Staking *StakingCallerSession) GetDemotions(nodeAddr common.Address, epoch *big.Int) ([]Demotion, error) {
	return _Staking.Contract.GetDemotions(&_Staking.CallOpts, nodeAddr, epoch)
}

// GetNode is a free data retrieval call binding the contract method 0x9d209048.
//
// Solidity: function getNode(address nodeAddr) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256,uint256,uint256,uint8))
func (_Staking *StakingCaller) GetNode(opts *bind.CallOpts, nodeAddr common.Address) (Node, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getNode", nodeAddr)

	if err != nil {
		return *new(Node), err
	}

	out0 := *abi.ConvertType(out[0], new(Node)).(*Node)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0x9d209048.
//
// Solidity: function getNode(address nodeAddr) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256,uint256,uint256,uint8))
func (_Staking *StakingSession) GetNode(nodeAddr common.Address) (Node, error) {
	return _Staking.Contract.GetNode(&_Staking.CallOpts, nodeAddr)
}

// GetNode is a free data retrieval call binding the contract method 0x9d209048.
//
// Solidity: function getNode(address nodeAddr) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256,uint256,uint256,uint8))
func (_Staking *StakingCallerSession) GetNode(nodeAddr common.Address) (Node, error) {
	return _Staking.Contract.GetNode(&_Staking.CallOpts, nodeAddr)
}

// GetNodeAvatar is a free data retrieval call binding the contract method 0x1474deaa.
//
// Solidity: function getNodeAvatar(address nodeAddr) view returns(string)
func (_Staking *StakingCaller) GetNodeAvatar(opts *bind.CallOpts, nodeAddr common.Address) (string, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getNodeAvatar", nodeAddr)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetNodeAvatar is a free data retrieval call binding the contract method 0x1474deaa.
//
// Solidity: function getNodeAvatar(address nodeAddr) view returns(string)
func (_Staking *StakingSession) GetNodeAvatar(nodeAddr common.Address) (string, error) {
	return _Staking.Contract.GetNodeAvatar(&_Staking.CallOpts, nodeAddr)
}

// GetNodeAvatar is a free data retrieval call binding the contract method 0x1474deaa.
//
// Solidity: function getNodeAvatar(address nodeAddr) view returns(string)
func (_Staking *StakingCallerSession) GetNodeAvatar(nodeAddr common.Address) (string, error) {
	return _Staking.Contract.GetNodeAvatar(&_Staking.CallOpts, nodeAddr)
}

// GetNodeCount is a free data retrieval call binding the contract method 0x39bf397e.
//
// Solidity: function getNodeCount() view returns(uint256)
func (_Staking *StakingCaller) GetNodeCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getNodeCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeCount is a free data retrieval call binding the contract method 0x39bf397e.
//
// Solidity: function getNodeCount() view returns(uint256)
func (_Staking *StakingSession) GetNodeCount() (*big.Int, error) {
	return _Staking.Contract.GetNodeCount(&_Staking.CallOpts)
}

// GetNodeCount is a free data retrieval call binding the contract method 0x39bf397e.
//
// Solidity: function getNodeCount() view returns(uint256)
func (_Staking *StakingCallerSession) GetNodeCount() (*big.Int, error) {
	return _Staking.Contract.GetNodeCount(&_Staking.CallOpts)
}

// GetNodes is a free data retrieval call binding the contract method 0x38c96b14.
//
// Solidity: function getNodes(address[] nodeAddrs) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256,uint256,uint256,uint8)[] nodes)
func (_Staking *StakingCaller) GetNodes(opts *bind.CallOpts, nodeAddrs []common.Address) ([]Node, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getNodes", nodeAddrs)

	if err != nil {
		return *new([]Node), err
	}

	out0 := *abi.ConvertType(out[0], new([]Node)).(*[]Node)

	return out0, err

}

// GetNodes is a free data retrieval call binding the contract method 0x38c96b14.
//
// Solidity: function getNodes(address[] nodeAddrs) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256,uint256,uint256,uint8)[] nodes)
func (_Staking *StakingSession) GetNodes(nodeAddrs []common.Address) ([]Node, error) {
	return _Staking.Contract.GetNodes(&_Staking.CallOpts, nodeAddrs)
}

// GetNodes is a free data retrieval call binding the contract method 0x38c96b14.
//
// Solidity: function getNodes(address[] nodeAddrs) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256,uint256,uint256,uint8)[] nodes)
func (_Staking *StakingCallerSession) GetNodes(nodeAddrs []common.Address) ([]Node, error) {
	return _Staking.Contract.GetNodes(&_Staking.CallOpts, nodeAddrs)
}

// GetPendingUnstake is a free data retrieval call binding the contract method 0xadfd065f.
//
// Solidity: function getPendingUnstake(uint256 requestId) view returns((address,address,uint256,uint256))
func (_Staking *StakingCaller) GetPendingUnstake(opts *bind.CallOpts, requestId *big.Int) (UnstakeRequest, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getPendingUnstake", requestId)

	if err != nil {
		return *new(UnstakeRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(UnstakeRequest)).(*UnstakeRequest)

	return out0, err

}

// GetPendingUnstake is a free data retrieval call binding the contract method 0xadfd065f.
//
// Solidity: function getPendingUnstake(uint256 requestId) view returns((address,address,uint256,uint256))
func (_Staking *StakingSession) GetPendingUnstake(requestId *big.Int) (UnstakeRequest, error) {
	return _Staking.Contract.GetPendingUnstake(&_Staking.CallOpts, requestId)
}

// GetPendingUnstake is a free data retrieval call binding the contract method 0xadfd065f.
//
// Solidity: function getPendingUnstake(uint256 requestId) view returns((address,address,uint256,uint256))
func (_Staking *StakingCallerSession) GetPendingUnstake(requestId *big.Int) (UnstakeRequest, error) {
	return _Staking.Contract.GetPendingUnstake(&_Staking.CallOpts, requestId)
}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x38a3c878.
//
// Solidity: function getPendingWithdrawal(uint256 requestId) view returns((address,uint40,uint256))
func (_Staking *StakingCaller) GetPendingWithdrawal(opts *bind.CallOpts, requestId *big.Int) (WithdrawalRequest, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getPendingWithdrawal", requestId)

	if err != nil {
		return *new(WithdrawalRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(WithdrawalRequest)).(*WithdrawalRequest)

	return out0, err

}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x38a3c878.
//
// Solidity: function getPendingWithdrawal(uint256 requestId) view returns((address,uint40,uint256))
func (_Staking *StakingSession) GetPendingWithdrawal(requestId *big.Int) (WithdrawalRequest, error) {
	return _Staking.Contract.GetPendingWithdrawal(&_Staking.CallOpts, requestId)
}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x38a3c878.
//
// Solidity: function getPendingWithdrawal(uint256 requestId) view returns((address,uint40,uint256))
func (_Staking *StakingCallerSession) GetPendingWithdrawal(requestId *big.Int) (WithdrawalRequest, error) {
	return _Staking.Contract.GetPendingWithdrawal(&_Staking.CallOpts, requestId)
}

// GetPoolInfo is a free data retrieval call binding the contract method 0x60246c88.
//
// Solidity: function getPoolInfo() view returns(uint256 totalOperationPoolTokens, uint256 totalStakingPoolTokens, uint256 totalSlashingPoolTokens)
func (_Staking *StakingCaller) GetPoolInfo(opts *bind.CallOpts) (struct {
	TotalOperationPoolTokens *big.Int
	TotalStakingPoolTokens   *big.Int
	TotalSlashingPoolTokens  *big.Int
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getPoolInfo")

	outstruct := new(struct {
		TotalOperationPoolTokens *big.Int
		TotalStakingPoolTokens   *big.Int
		TotalSlashingPoolTokens  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalOperationPoolTokens = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalStakingPoolTokens = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TotalSlashingPoolTokens = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPoolInfo is a free data retrieval call binding the contract method 0x60246c88.
//
// Solidity: function getPoolInfo() view returns(uint256 totalOperationPoolTokens, uint256 totalStakingPoolTokens, uint256 totalSlashingPoolTokens)
func (_Staking *StakingSession) GetPoolInfo() (struct {
	TotalOperationPoolTokens *big.Int
	TotalStakingPoolTokens   *big.Int
	TotalSlashingPoolTokens  *big.Int
}, error) {
	return _Staking.Contract.GetPoolInfo(&_Staking.CallOpts)
}

// GetPoolInfo is a free data retrieval call binding the contract method 0x60246c88.
//
// Solidity: function getPoolInfo() view returns(uint256 totalOperationPoolTokens, uint256 totalStakingPoolTokens, uint256 totalSlashingPoolTokens)
func (_Staking *StakingCallerSession) GetPoolInfo() (struct {
	TotalOperationPoolTokens *big.Int
	TotalStakingPoolTokens   *big.Int
	TotalSlashingPoolTokens  *big.Int
}, error) {
	return _Staking.Contract.GetPoolInfo(&_Staking.CallOpts)
}

// GetPublicPool is a free data retrieval call binding the contract method 0xc84c42a3.
//
// Solidity: function getPublicPool() pure returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256,uint256,uint256,uint8))
func (_Staking *StakingCaller) GetPublicPool(opts *bind.CallOpts) (Node, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getPublicPool")

	if err != nil {
		return *new(Node), err
	}

	out0 := *abi.ConvertType(out[0], new(Node)).(*Node)

	return out0, err

}

// GetPublicPool is a free data retrieval call binding the contract method 0xc84c42a3.
//
// Solidity: function getPublicPool() pure returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256,uint256,uint256,uint8))
func (_Staking *StakingSession) GetPublicPool() (Node, error) {
	return _Staking.Contract.GetPublicPool(&_Staking.CallOpts)
}

// GetPublicPool is a free data retrieval call binding the contract method 0xc84c42a3.
//
// Solidity: function getPublicPool() pure returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256,uint256,uint256,uint8))
func (_Staking *StakingCallerSession) GetPublicPool() (Node, error) {
	return _Staking.Contract.GetPublicPool(&_Staking.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Staking *StakingCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Staking *StakingSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Staking.Contract.GetRoleAdmin(&_Staking.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Staking *StakingCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Staking.Contract.GetRoleAdmin(&_Staking.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Staking *StakingCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Staking *StakingSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Staking.Contract.GetRoleMember(&_Staking.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Staking *StakingCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Staking.Contract.GetRoleMember(&_Staking.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Staking *StakingCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Staking *StakingSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Staking.Contract.GetRoleMemberCount(&_Staking.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Staking *StakingCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Staking.Contract.GetRoleMemberCount(&_Staking.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Staking *StakingCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Staking *StakingSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Staking.Contract.HasRole(&_Staking.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Staking *StakingCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Staking.Contract.HasRole(&_Staking.CallOpts, role, account)
}

// IsAlphaPhase is a free data retrieval call binding the contract method 0x69ff71f6.
//
// Solidity: function isAlphaPhase() view returns(bool)
func (_Staking *StakingCaller) IsAlphaPhase(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "isAlphaPhase")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAlphaPhase is a free data retrieval call binding the contract method 0x69ff71f6.
//
// Solidity: function isAlphaPhase() view returns(bool)
func (_Staking *StakingSession) IsAlphaPhase() (bool, error) {
	return _Staking.Contract.IsAlphaPhase(&_Staking.CallOpts)
}

// IsAlphaPhase is a free data retrieval call binding the contract method 0x69ff71f6.
//
// Solidity: function isAlphaPhase() view returns(bool)
func (_Staking *StakingCallerSession) IsAlphaPhase() (bool, error) {
	return _Staking.Contract.IsAlphaPhase(&_Staking.CallOpts)
}

// IsSettlementPhase is a free data retrieval call binding the contract method 0x2e75fd59.
//
// Solidity: function isSettlementPhase() view returns(bool)
func (_Staking *StakingCaller) IsSettlementPhase(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "isSettlementPhase")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSettlementPhase is a free data retrieval call binding the contract method 0x2e75fd59.
//
// Solidity: function isSettlementPhase() view returns(bool)
func (_Staking *StakingSession) IsSettlementPhase() (bool, error) {
	return _Staking.Contract.IsSettlementPhase(&_Staking.CallOpts)
}

// IsSettlementPhase is a free data retrieval call binding the contract method 0x2e75fd59.
//
// Solidity: function isSettlementPhase() view returns(bool)
func (_Staking *StakingCallerSession) IsSettlementPhase() (bool, error) {
	return _Staking.Contract.IsSettlementPhase(&_Staking.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Staking *StakingCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Staking *StakingSession) Paused() (bool, error) {
	return _Staking.Contract.Paused(&_Staking.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Staking *StakingCallerSession) Paused() (bool, error) {
	return _Staking.Contract.Paused(&_Staking.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Staking *StakingCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Staking *StakingSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Staking.Contract.SupportsInterface(&_Staking.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Staking *StakingCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Staking.Contract.SupportsInterface(&_Staking.CallOpts, interfaceId)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Staking *StakingCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Staking *StakingSession) Version() (string, error) {
	return _Staking.Contract.Version(&_Staking.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Staking *StakingCallerSession) Version() (string, error) {
	return _Staking.Contract.Version(&_Staking.CallOpts)
}

// ClaimUnstake is a paid mutator transaction binding the contract method 0x04a4fb10.
//
// Solidity: function claimUnstake(uint256[] requestIds) returns()
func (_Staking *StakingTransactor) ClaimUnstake(opts *bind.TransactOpts, requestIds []*big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "claimUnstake", requestIds)
}

// ClaimUnstake is a paid mutator transaction binding the contract method 0x04a4fb10.
//
// Solidity: function claimUnstake(uint256[] requestIds) returns()
func (_Staking *StakingSession) ClaimUnstake(requestIds []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.ClaimUnstake(&_Staking.TransactOpts, requestIds)
}

// ClaimUnstake is a paid mutator transaction binding the contract method 0x04a4fb10.
//
// Solidity: function claimUnstake(uint256[] requestIds) returns()
func (_Staking *StakingTransactorSession) ClaimUnstake(requestIds []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.ClaimUnstake(&_Staking.TransactOpts, requestIds)
}

// ClaimWithdrawal is a paid mutator transaction binding the contract method 0x3c256b98.
//
// Solidity: function claimWithdrawal(uint256[] requestIds) returns()
func (_Staking *StakingTransactor) ClaimWithdrawal(opts *bind.TransactOpts, requestIds []*big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "claimWithdrawal", requestIds)
}

// ClaimWithdrawal is a paid mutator transaction binding the contract method 0x3c256b98.
//
// Solidity: function claimWithdrawal(uint256[] requestIds) returns()
func (_Staking *StakingSession) ClaimWithdrawal(requestIds []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.ClaimWithdrawal(&_Staking.TransactOpts, requestIds)
}

// ClaimWithdrawal is a paid mutator transaction binding the contract method 0x3c256b98.
//
// Solidity: function claimWithdrawal(uint256[] requestIds) returns()
func (_Staking *StakingTransactorSession) ClaimWithdrawal(requestIds []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.ClaimWithdrawal(&_Staking.TransactOpts, requestIds)
}

// CommitSlashing is a paid mutator transaction binding the contract method 0xbe168625.
//
// Solidity: function commitSlashing(address nodeAddr, uint256 epoch) returns()
func (_Staking *StakingTransactor) CommitSlashing(opts *bind.TransactOpts, nodeAddr common.Address, epoch *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "commitSlashing", nodeAddr, epoch)
}

// CommitSlashing is a paid mutator transaction binding the contract method 0xbe168625.
//
// Solidity: function commitSlashing(address nodeAddr, uint256 epoch) returns()
func (_Staking *StakingSession) CommitSlashing(nodeAddr common.Address, epoch *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.CommitSlashing(&_Staking.TransactOpts, nodeAddr, epoch)
}

// CommitSlashing is a paid mutator transaction binding the contract method 0xbe168625.
//
// Solidity: function commitSlashing(address nodeAddr, uint256 epoch) returns()
func (_Staking *StakingTransactorSession) CommitSlashing(nodeAddr common.Address, epoch *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.CommitSlashing(&_Staking.TransactOpts, nodeAddr, epoch)
}

// CreateNode is a paid mutator transaction binding the contract method 0x96531623.
//
// Solidity: function createNode(string name, string description, uint64 taxRateBasisPoints, bool publicGood) payable returns()
func (_Staking *StakingTransactor) CreateNode(opts *bind.TransactOpts, name string, description string, taxRateBasisPoints uint64, publicGood bool) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "createNode", name, description, taxRateBasisPoints, publicGood)
}

// CreateNode is a paid mutator transaction binding the contract method 0x96531623.
//
// Solidity: function createNode(string name, string description, uint64 taxRateBasisPoints, bool publicGood) payable returns()
func (_Staking *StakingSession) CreateNode(name string, description string, taxRateBasisPoints uint64, publicGood bool) (*types.Transaction, error) {
	return _Staking.Contract.CreateNode(&_Staking.TransactOpts, name, description, taxRateBasisPoints, publicGood)
}

// CreateNode is a paid mutator transaction binding the contract method 0x96531623.
//
// Solidity: function createNode(string name, string description, uint64 taxRateBasisPoints, bool publicGood) payable returns()
func (_Staking *StakingTransactorSession) CreateNode(name string, description string, taxRateBasisPoints uint64, publicGood bool) (*types.Transaction, error) {
	return _Staking.Contract.CreateNode(&_Staking.TransactOpts, name, description, taxRateBasisPoints, publicGood)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Staking *StakingTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Staking *StakingSession) Deposit() (*types.Transaction, error) {
	return _Staking.Contract.Deposit(&_Staking.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Staking *StakingTransactorSession) Deposit() (*types.Transaction, error) {
	return _Staking.Contract.Deposit(&_Staking.TransactOpts)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x8e3e6174.
//
// Solidity: function distributeRewards(uint256[3] epochInfo, address[] nodeAddrs, uint256[] operationRewards, uint256[] stakingRewards, uint256[] requestCounts, uint256 publicPoolRewards) payable returns()
func (_Staking *StakingTransactor) DistributeRewards(opts *bind.TransactOpts, epochInfo [3]*big.Int, nodeAddrs []common.Address, operationRewards []*big.Int, stakingRewards []*big.Int, requestCounts []*big.Int, publicPoolRewards *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "distributeRewards", epochInfo, nodeAddrs, operationRewards, stakingRewards, requestCounts, publicPoolRewards)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x8e3e6174.
//
// Solidity: function distributeRewards(uint256[3] epochInfo, address[] nodeAddrs, uint256[] operationRewards, uint256[] stakingRewards, uint256[] requestCounts, uint256 publicPoolRewards) payable returns()
func (_Staking *StakingSession) DistributeRewards(epochInfo [3]*big.Int, nodeAddrs []common.Address, operationRewards []*big.Int, stakingRewards []*big.Int, requestCounts []*big.Int, publicPoolRewards *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.DistributeRewards(&_Staking.TransactOpts, epochInfo, nodeAddrs, operationRewards, stakingRewards, requestCounts, publicPoolRewards)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x8e3e6174.
//
// Solidity: function distributeRewards(uint256[3] epochInfo, address[] nodeAddrs, uint256[] operationRewards, uint256[] stakingRewards, uint256[] requestCounts, uint256 publicPoolRewards) payable returns()
func (_Staking *StakingTransactorSession) DistributeRewards(epochInfo [3]*big.Int, nodeAddrs []common.Address, operationRewards []*big.Int, stakingRewards []*big.Int, requestCounts []*big.Int, publicPoolRewards *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.DistributeRewards(&_Staking.TransactOpts, epochInfo, nodeAddrs, operationRewards, stakingRewards, requestCounts, publicPoolRewards)
}

// Exit is a paid mutator transaction binding the contract method 0xe9fad8ee.
//
// Solidity: function exit() returns()
func (_Staking *StakingTransactor) Exit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "exit")
}

// Exit is a paid mutator transaction binding the contract method 0xe9fad8ee.
//
// Solidity: function exit() returns()
func (_Staking *StakingSession) Exit() (*types.Transaction, error) {
	return _Staking.Contract.Exit(&_Staking.TransactOpts)
}

// Exit is a paid mutator transaction binding the contract method 0xe9fad8ee.
//
// Solidity: function exit() returns()
func (_Staking *StakingTransactorSession) Exit() (*types.Transaction, error) {
	return _Staking.Contract.Exit(&_Staking.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Staking *StakingTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Staking *StakingSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Staking.Contract.GrantRole(&_Staking.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Staking *StakingTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Staking.Contract.GrantRole(&_Staking.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xfecf9734.
//
// Solidity: function initialize(address chips, address pauseAccount, address oracleAccount, bool isAlphaPhase_) returns()
func (_Staking *StakingTransactor) Initialize(opts *bind.TransactOpts, chips common.Address, pauseAccount common.Address, oracleAccount common.Address, isAlphaPhase_ bool) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "initialize", chips, pauseAccount, oracleAccount, isAlphaPhase_)
}

// Initialize is a paid mutator transaction binding the contract method 0xfecf9734.
//
// Solidity: function initialize(address chips, address pauseAccount, address oracleAccount, bool isAlphaPhase_) returns()
func (_Staking *StakingSession) Initialize(chips common.Address, pauseAccount common.Address, oracleAccount common.Address, isAlphaPhase_ bool) (*types.Transaction, error) {
	return _Staking.Contract.Initialize(&_Staking.TransactOpts, chips, pauseAccount, oracleAccount, isAlphaPhase_)
}

// Initialize is a paid mutator transaction binding the contract method 0xfecf9734.
//
// Solidity: function initialize(address chips, address pauseAccount, address oracleAccount, bool isAlphaPhase_) returns()
func (_Staking *StakingTransactorSession) Initialize(chips common.Address, pauseAccount common.Address, oracleAccount common.Address, isAlphaPhase_ bool) (*types.Transaction, error) {
	return _Staking.Contract.Initialize(&_Staking.TransactOpts, chips, pauseAccount, oracleAccount, isAlphaPhase_)
}

// MergeChips is a paid mutator transaction binding the contract method 0xca91c816.
//
// Solidity: function mergeChips(uint256[] chipIds) returns(uint256 newTokenId)
func (_Staking *StakingTransactor) MergeChips(opts *bind.TransactOpts, chipIds []*big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "mergeChips", chipIds)
}

// MergeChips is a paid mutator transaction binding the contract method 0xca91c816.
//
// Solidity: function mergeChips(uint256[] chipIds) returns(uint256 newTokenId)
func (_Staking *StakingSession) MergeChips(chipIds []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.MergeChips(&_Staking.TransactOpts, chipIds)
}

// MergeChips is a paid mutator transaction binding the contract method 0xca91c816.
//
// Solidity: function mergeChips(uint256[] chipIds) returns(uint256 newTokenId)
func (_Staking *StakingTransactorSession) MergeChips(chipIds []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.MergeChips(&_Staking.TransactOpts, chipIds)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Staking *StakingTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Staking *StakingSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Staking.Contract.Multicall(&_Staking.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Staking *StakingTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Staking.Contract.Multicall(&_Staking.TransactOpts, data)
}

// Online is a paid mutator transaction binding the contract method 0x5cd9f87a.
//
// Solidity: function online() returns()
func (_Staking *StakingTransactor) Online(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "online")
}

// Online is a paid mutator transaction binding the contract method 0x5cd9f87a.
//
// Solidity: function online() returns()
func (_Staking *StakingSession) Online() (*types.Transaction, error) {
	return _Staking.Contract.Online(&_Staking.TransactOpts)
}

// Online is a paid mutator transaction binding the contract method 0x5cd9f87a.
//
// Solidity: function online() returns()
func (_Staking *StakingTransactorSession) Online() (*types.Transaction, error) {
	return _Staking.Contract.Online(&_Staking.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Staking *StakingTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Staking *StakingSession) Pause() (*types.Transaction, error) {
	return _Staking.Contract.Pause(&_Staking.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Staking *StakingTransactorSession) Pause() (*types.Transaction, error) {
	return _Staking.Contract.Pause(&_Staking.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_Staking *StakingTransactor) Register(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "register")
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_Staking *StakingSession) Register() (*types.Transaction, error) {
	return _Staking.Contract.Register(&_Staking.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_Staking *StakingTransactorSession) Register() (*types.Transaction, error) {
	return _Staking.Contract.Register(&_Staking.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Staking *StakingTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Staking *StakingSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Staking.Contract.RenounceRole(&_Staking.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Staking *StakingTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Staking.Contract.RenounceRole(&_Staking.TransactOpts, role, callerConfirmation)
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xbcdd4190.
//
// Solidity: function requestUnstake(address nodeAddr, uint256[] chipIds) returns(uint256 requestId)
func (_Staking *StakingTransactor) RequestUnstake(opts *bind.TransactOpts, nodeAddr common.Address, chipIds []*big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "requestUnstake", nodeAddr, chipIds)
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xbcdd4190.
//
// Solidity: function requestUnstake(address nodeAddr, uint256[] chipIds) returns(uint256 requestId)
func (_Staking *StakingSession) RequestUnstake(nodeAddr common.Address, chipIds []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.RequestUnstake(&_Staking.TransactOpts, nodeAddr, chipIds)
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xbcdd4190.
//
// Solidity: function requestUnstake(address nodeAddr, uint256[] chipIds) returns(uint256 requestId)
func (_Staking *StakingTransactorSession) RequestUnstake(nodeAddr common.Address, chipIds []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.RequestUnstake(&_Staking.TransactOpts, nodeAddr, chipIds)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0x9ee679e8.
//
// Solidity: function requestWithdrawal(uint256 amount) returns(uint256 requestId)
func (_Staking *StakingTransactor) RequestWithdrawal(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "requestWithdrawal", amount)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0x9ee679e8.
//
// Solidity: function requestWithdrawal(uint256 amount) returns(uint256 requestId)
func (_Staking *StakingSession) RequestWithdrawal(amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.RequestWithdrawal(&_Staking.TransactOpts, amount)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0x9ee679e8.
//
// Solidity: function requestWithdrawal(uint256 amount) returns(uint256 requestId)
func (_Staking *StakingTransactorSession) RequestWithdrawal(amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.RequestWithdrawal(&_Staking.TransactOpts, amount)
}

// RevokeDemotions is a paid mutator transaction binding the contract method 0x6d623207.
//
// Solidity: function revokeDemotions(address nodeAddr, uint256 epoch, uint256[] demotionIdsToRevoke) returns()
func (_Staking *StakingTransactor) RevokeDemotions(opts *bind.TransactOpts, nodeAddr common.Address, epoch *big.Int, demotionIdsToRevoke []*big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "revokeDemotions", nodeAddr, epoch, demotionIdsToRevoke)
}

// RevokeDemotions is a paid mutator transaction binding the contract method 0x6d623207.
//
// Solidity: function revokeDemotions(address nodeAddr, uint256 epoch, uint256[] demotionIdsToRevoke) returns()
func (_Staking *StakingSession) RevokeDemotions(nodeAddr common.Address, epoch *big.Int, demotionIdsToRevoke []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.RevokeDemotions(&_Staking.TransactOpts, nodeAddr, epoch, demotionIdsToRevoke)
}

// RevokeDemotions is a paid mutator transaction binding the contract method 0x6d623207.
//
// Solidity: function revokeDemotions(address nodeAddr, uint256 epoch, uint256[] demotionIdsToRevoke) returns()
func (_Staking *StakingTransactorSession) RevokeDemotions(nodeAddr common.Address, epoch *big.Int, demotionIdsToRevoke []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.RevokeDemotions(&_Staking.TransactOpts, nodeAddr, epoch, demotionIdsToRevoke)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Staking *StakingTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Staking *StakingSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Staking.Contract.RevokeRole(&_Staking.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Staking *StakingTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Staking.Contract.RevokeRole(&_Staking.TransactOpts, role, account)
}

// SetNodeStatus is a paid mutator transaction binding the contract method 0x03b0d316.
//
// Solidity: function setNodeStatus(address[] nodeAddrs, uint8[] status) returns()
func (_Staking *StakingTransactor) SetNodeStatus(opts *bind.TransactOpts, nodeAddrs []common.Address, status []uint8) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "setNodeStatus", nodeAddrs, status)
}

// SetNodeStatus is a paid mutator transaction binding the contract method 0x03b0d316.
//
// Solidity: function setNodeStatus(address[] nodeAddrs, uint8[] status) returns()
func (_Staking *StakingSession) SetNodeStatus(nodeAddrs []common.Address, status []uint8) (*types.Transaction, error) {
	return _Staking.Contract.SetNodeStatus(&_Staking.TransactOpts, nodeAddrs, status)
}

// SetNodeStatus is a paid mutator transaction binding the contract method 0x03b0d316.
//
// Solidity: function setNodeStatus(address[] nodeAddrs, uint8[] status) returns()
func (_Staking *StakingTransactorSession) SetNodeStatus(nodeAddrs []common.Address, status []uint8) (*types.Transaction, error) {
	return _Staking.Contract.SetNodeStatus(&_Staking.TransactOpts, nodeAddrs, status)
}

// SetSettlementPhase is a paid mutator transaction binding the contract method 0x4e7a1286.
//
// Solidity: function setSettlementPhase(bool enabled) returns()
func (_Staking *StakingTransactor) SetSettlementPhase(opts *bind.TransactOpts, enabled bool) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "setSettlementPhase", enabled)
}

// SetSettlementPhase is a paid mutator transaction binding the contract method 0x4e7a1286.
//
// Solidity: function setSettlementPhase(bool enabled) returns()
func (_Staking *StakingSession) SetSettlementPhase(enabled bool) (*types.Transaction, error) {
	return _Staking.Contract.SetSettlementPhase(&_Staking.TransactOpts, enabled)
}

// SetSettlementPhase is a paid mutator transaction binding the contract method 0x4e7a1286.
//
// Solidity: function setSettlementPhase(bool enabled) returns()
func (_Staking *StakingTransactorSession) SetSettlementPhase(enabled bool) (*types.Transaction, error) {
	return _Staking.Contract.SetSettlementPhase(&_Staking.TransactOpts, enabled)
}

// SetTaxRateBasisPoints4Node is a paid mutator transaction binding the contract method 0xc7057c1f.
//
// Solidity: function setTaxRateBasisPoints4Node(uint64 taxRateBasisPoints) returns()
func (_Staking *StakingTransactor) SetTaxRateBasisPoints4Node(opts *bind.TransactOpts, taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "setTaxRateBasisPoints4Node", taxRateBasisPoints)
}

// SetTaxRateBasisPoints4Node is a paid mutator transaction binding the contract method 0xc7057c1f.
//
// Solidity: function setTaxRateBasisPoints4Node(uint64 taxRateBasisPoints) returns()
func (_Staking *StakingSession) SetTaxRateBasisPoints4Node(taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Staking.Contract.SetTaxRateBasisPoints4Node(&_Staking.TransactOpts, taxRateBasisPoints)
}

// SetTaxRateBasisPoints4Node is a paid mutator transaction binding the contract method 0xc7057c1f.
//
// Solidity: function setTaxRateBasisPoints4Node(uint64 taxRateBasisPoints) returns()
func (_Staking *StakingTransactorSession) SetTaxRateBasisPoints4Node(taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Staking.Contract.SetTaxRateBasisPoints4Node(&_Staking.TransactOpts, taxRateBasisPoints)
}

// SetTaxRateBasisPoints4PublicPool is a paid mutator transaction binding the contract method 0xe3fb8dca.
//
// Solidity: function setTaxRateBasisPoints4PublicPool(uint64 taxRateBasisPoints) returns()
func (_Staking *StakingTransactor) SetTaxRateBasisPoints4PublicPool(opts *bind.TransactOpts, taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "setTaxRateBasisPoints4PublicPool", taxRateBasisPoints)
}

// SetTaxRateBasisPoints4PublicPool is a paid mutator transaction binding the contract method 0xe3fb8dca.
//
// Solidity: function setTaxRateBasisPoints4PublicPool(uint64 taxRateBasisPoints) returns()
func (_Staking *StakingSession) SetTaxRateBasisPoints4PublicPool(taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Staking.Contract.SetTaxRateBasisPoints4PublicPool(&_Staking.TransactOpts, taxRateBasisPoints)
}

// SetTaxRateBasisPoints4PublicPool is a paid mutator transaction binding the contract method 0xe3fb8dca.
//
// Solidity: function setTaxRateBasisPoints4PublicPool(uint64 taxRateBasisPoints) returns()
func (_Staking *StakingTransactorSession) SetTaxRateBasisPoints4PublicPool(taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Staking.Contract.SetTaxRateBasisPoints4PublicPool(&_Staking.TransactOpts, taxRateBasisPoints)
}

// Stake is a paid mutator transaction binding the contract method 0x26476204.
//
// Solidity: function stake(address nodeAddr) payable returns(uint256 tokenId)
func (_Staking *StakingTransactor) Stake(opts *bind.TransactOpts, nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "stake", nodeAddr)
}

// Stake is a paid mutator transaction binding the contract method 0x26476204.
//
// Solidity: function stake(address nodeAddr) payable returns(uint256 tokenId)
func (_Staking *StakingSession) Stake(nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.Stake(&_Staking.TransactOpts, nodeAddr)
}

// Stake is a paid mutator transaction binding the contract method 0x26476204.
//
// Solidity: function stake(address nodeAddr) payable returns(uint256 tokenId)
func (_Staking *StakingTransactorSession) Stake(nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.Stake(&_Staking.TransactOpts, nodeAddr)
}

// StakeToPublicPool is a paid mutator transaction binding the contract method 0x379f8100.
//
// Solidity: function stakeToPublicPool(address nodeAddr) payable returns(uint256 tokenId)
func (_Staking *StakingTransactor) StakeToPublicPool(opts *bind.TransactOpts, nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "stakeToPublicPool", nodeAddr)
}

// StakeToPublicPool is a paid mutator transaction binding the contract method 0x379f8100.
//
// Solidity: function stakeToPublicPool(address nodeAddr) payable returns(uint256 tokenId)
func (_Staking *StakingSession) StakeToPublicPool(nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.StakeToPublicPool(&_Staking.TransactOpts, nodeAddr)
}

// StakeToPublicPool is a paid mutator transaction binding the contract method 0x379f8100.
//
// Solidity: function stakeToPublicPool(address nodeAddr) payable returns(uint256 tokenId)
func (_Staking *StakingTransactorSession) StakeToPublicPool(nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.StakeToPublicPool(&_Staking.TransactOpts, nodeAddr)
}

// SubmitDemotions is a paid mutator transaction binding the contract method 0x287e9ac8.
//
// Solidity: function submitDemotions(uint256 epoch, address[] nodeAddrs, string[] reasons, address[] reporters) returns()
func (_Staking *StakingTransactor) SubmitDemotions(opts *bind.TransactOpts, epoch *big.Int, nodeAddrs []common.Address, reasons []string, reporters []common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "submitDemotions", epoch, nodeAddrs, reasons, reporters)
}

// SubmitDemotions is a paid mutator transaction binding the contract method 0x287e9ac8.
//
// Solidity: function submitDemotions(uint256 epoch, address[] nodeAddrs, string[] reasons, address[] reporters) returns()
func (_Staking *StakingSession) SubmitDemotions(epoch *big.Int, nodeAddrs []common.Address, reasons []string, reporters []common.Address) (*types.Transaction, error) {
	return _Staking.Contract.SubmitDemotions(&_Staking.TransactOpts, epoch, nodeAddrs, reasons, reporters)
}

// SubmitDemotions is a paid mutator transaction binding the contract method 0x287e9ac8.
//
// Solidity: function submitDemotions(uint256 epoch, address[] nodeAddrs, string[] reasons, address[] reporters) returns()
func (_Staking *StakingTransactorSession) SubmitDemotions(epoch *big.Int, nodeAddrs []common.Address, reasons []string, reporters []common.Address) (*types.Transaction, error) {
	return _Staking.Contract.SubmitDemotions(&_Staking.TransactOpts, epoch, nodeAddrs, reasons, reporters)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Staking *StakingTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Staking *StakingSession) Unpause() (*types.Transaction, error) {
	return _Staking.Contract.Unpause(&_Staking.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Staking *StakingTransactorSession) Unpause() (*types.Transaction, error) {
	return _Staking.Contract.Unpause(&_Staking.TransactOpts)
}

// UpdateNode is a paid mutator transaction binding the contract method 0x517fb203.
//
// Solidity: function updateNode(string name, string description) returns()
func (_Staking *StakingTransactor) UpdateNode(opts *bind.TransactOpts, name string, description string) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "updateNode", name, description)
}

// UpdateNode is a paid mutator transaction binding the contract method 0x517fb203.
//
// Solidity: function updateNode(string name, string description) returns()
func (_Staking *StakingSession) UpdateNode(name string, description string) (*types.Transaction, error) {
	return _Staking.Contract.UpdateNode(&_Staking.TransactOpts, name, description)
}

// UpdateNode is a paid mutator transaction binding the contract method 0x517fb203.
//
// Solidity: function updateNode(string name, string description) returns()
func (_Staking *StakingTransactorSession) UpdateNode(name string, description string) (*types.Transaction, error) {
	return _Staking.Contract.UpdateNode(&_Staking.TransactOpts, name, description)
}

// Withdraw2Treasury is a paid mutator transaction binding the contract method 0x4a7dfc90.
//
// Solidity: function withdraw2Treasury() returns()
func (_Staking *StakingTransactor) Withdraw2Treasury(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdraw2Treasury")
}

// Withdraw2Treasury is a paid mutator transaction binding the contract method 0x4a7dfc90.
//
// Solidity: function withdraw2Treasury() returns()
func (_Staking *StakingSession) Withdraw2Treasury() (*types.Transaction, error) {
	return _Staking.Contract.Withdraw2Treasury(&_Staking.TransactOpts)
}

// Withdraw2Treasury is a paid mutator transaction binding the contract method 0x4a7dfc90.
//
// Solidity: function withdraw2Treasury() returns()
func (_Staking *StakingTransactorSession) Withdraw2Treasury() (*types.Transaction, error) {
	return _Staking.Contract.Withdraw2Treasury(&_Staking.TransactOpts)
}

// StakingInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Staking contract.
type StakingInitializedIterator struct {
	Event *StakingInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingInitialized represents a Initialized event raised by the Staking contract.
type StakingInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Staking *StakingFilterer) FilterInitialized(opts *bind.FilterOpts) (*StakingInitializedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StakingInitializedIterator{contract: _Staking.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Staking *StakingFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StakingInitialized) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingInitialized)
				if err := _Staking.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Staking *StakingFilterer) ParseInitialized(log types.Log) (*StakingInitialized, error) {
	event := new(StakingInitialized)
	if err := _Staking.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Staking contract.
type StakingPausedIterator struct {
	Event *StakingPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPaused represents a Paused event raised by the Staking contract.
type StakingPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Staking *StakingFilterer) FilterPaused(opts *bind.FilterOpts) (*StakingPausedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &StakingPausedIterator{contract: _Staking.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Staking *StakingFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *StakingPaused) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPaused)
				if err := _Staking.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Staking *StakingFilterer) ParsePaused(log types.Log) (*StakingPaused, error) {
	event := new(StakingPaused)
	if err := _Staking.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPublicGoodRewardDistributedIterator is returned from FilterPublicGoodRewardDistributed and is used to iterate over the raw logs and unpacked data for PublicGoodRewardDistributed events raised by the Staking contract.
type StakingPublicGoodRewardDistributedIterator struct {
	Event *StakingPublicGoodRewardDistributed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingPublicGoodRewardDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPublicGoodRewardDistributed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingPublicGoodRewardDistributed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingPublicGoodRewardDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPublicGoodRewardDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPublicGoodRewardDistributed represents a PublicGoodRewardDistributed event raised by the Staking contract.
type StakingPublicGoodRewardDistributed struct {
	Epoch             *big.Int
	StartTimestamp    *big.Int
	EndTimestamp      *big.Int
	PublicPoolRewards *big.Int
	PublicPoolTax     *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPublicGoodRewardDistributed is a free log retrieval operation binding the contract event 0xab7d25a2f6206ef56c88807f2474ddcd97e1a6323cb25149cde3a607fed6f2d7.
//
// Solidity: event PublicGoodRewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, uint256 publicPoolRewards, uint256 publicPoolTax)
func (_Staking *StakingFilterer) FilterPublicGoodRewardDistributed(opts *bind.FilterOpts, epoch []*big.Int) (*StakingPublicGoodRewardDistributedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "PublicGoodRewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return &StakingPublicGoodRewardDistributedIterator{contract: _Staking.contract, event: "PublicGoodRewardDistributed", logs: logs, sub: sub}, nil
}

// WatchPublicGoodRewardDistributed is a free log subscription operation binding the contract event 0xab7d25a2f6206ef56c88807f2474ddcd97e1a6323cb25149cde3a607fed6f2d7.
//
// Solidity: event PublicGoodRewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, uint256 publicPoolRewards, uint256 publicPoolTax)
func (_Staking *StakingFilterer) WatchPublicGoodRewardDistributed(opts *bind.WatchOpts, sink chan<- *StakingPublicGoodRewardDistributed, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "PublicGoodRewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPublicGoodRewardDistributed)
				if err := _Staking.contract.UnpackLog(event, "PublicGoodRewardDistributed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePublicGoodRewardDistributed is a log parse operation binding the contract event 0xab7d25a2f6206ef56c88807f2474ddcd97e1a6323cb25149cde3a607fed6f2d7.
//
// Solidity: event PublicGoodRewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, uint256 publicPoolRewards, uint256 publicPoolTax)
func (_Staking *StakingFilterer) ParsePublicGoodRewardDistributed(log types.Log) (*StakingPublicGoodRewardDistributed, error) {
	event := new(StakingPublicGoodRewardDistributed)
	if err := _Staking.contract.UnpackLog(event, "PublicGoodRewardDistributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRewardDistributedIterator is returned from FilterRewardDistributed and is used to iterate over the raw logs and unpacked data for RewardDistributed events raised by the Staking contract.
type StakingRewardDistributedIterator struct {
	Event *StakingRewardDistributed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingRewardDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRewardDistributed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingRewardDistributed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingRewardDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRewardDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRewardDistributed represents a RewardDistributed event raised by the Staking contract.
type StakingRewardDistributed struct {
	Epoch            *big.Int
	StartTimestamp   *big.Int
	EndTimestamp     *big.Int
	NodeAddrs        []common.Address
	OperationRewards []*big.Int
	StakingRewards   []*big.Int
	TaxCollected     []*big.Int
	RequestCounts    []*big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRewardDistributed is a free log retrieval operation binding the contract event 0x8ea79f19e90b084c2009d3490a097547d8bbb315a883b9efec0996502c1dd7ae.
//
// Solidity: event RewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, address[] nodeAddrs, uint256[] operationRewards, uint256[] stakingRewards, uint256[] taxCollected, uint256[] requestCounts)
func (_Staking *StakingFilterer) FilterRewardDistributed(opts *bind.FilterOpts, epoch []*big.Int) (*StakingRewardDistributedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "RewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return &StakingRewardDistributedIterator{contract: _Staking.contract, event: "RewardDistributed", logs: logs, sub: sub}, nil
}

// WatchRewardDistributed is a free log subscription operation binding the contract event 0x8ea79f19e90b084c2009d3490a097547d8bbb315a883b9efec0996502c1dd7ae.
//
// Solidity: event RewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, address[] nodeAddrs, uint256[] operationRewards, uint256[] stakingRewards, uint256[] taxCollected, uint256[] requestCounts)
func (_Staking *StakingFilterer) WatchRewardDistributed(opts *bind.WatchOpts, sink chan<- *StakingRewardDistributed, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "RewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRewardDistributed)
				if err := _Staking.contract.UnpackLog(event, "RewardDistributed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardDistributed is a log parse operation binding the contract event 0x8ea79f19e90b084c2009d3490a097547d8bbb315a883b9efec0996502c1dd7ae.
//
// Solidity: event RewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, address[] nodeAddrs, uint256[] operationRewards, uint256[] stakingRewards, uint256[] taxCollected, uint256[] requestCounts)
func (_Staking *StakingFilterer) ParseRewardDistributed(log types.Log) (*StakingRewardDistributed, error) {
	event := new(StakingRewardDistributed)
	if err := _Staking.contract.UnpackLog(event, "RewardDistributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Staking contract.
type StakingRoleAdminChangedIterator struct {
	Event *StakingRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRoleAdminChanged represents a RoleAdminChanged event raised by the Staking contract.
type StakingRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Staking *StakingFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*StakingRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &StakingRoleAdminChangedIterator{contract: _Staking.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Staking *StakingFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *StakingRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRoleAdminChanged)
				if err := _Staking.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Staking *StakingFilterer) ParseRoleAdminChanged(log types.Log) (*StakingRoleAdminChanged, error) {
	event := new(StakingRoleAdminChanged)
	if err := _Staking.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Staking contract.
type StakingRoleGrantedIterator struct {
	Event *StakingRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRoleGranted represents a RoleGranted event raised by the Staking contract.
type StakingRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Staking *StakingFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*StakingRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingRoleGrantedIterator{contract: _Staking.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Staking *StakingFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *StakingRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRoleGranted)
				if err := _Staking.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Staking *StakingFilterer) ParseRoleGranted(log types.Log) (*StakingRoleGranted, error) {
	event := new(StakingRoleGranted)
	if err := _Staking.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Staking contract.
type StakingRoleRevokedIterator struct {
	Event *StakingRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRoleRevoked represents a RoleRevoked event raised by the Staking contract.
type StakingRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Staking *StakingFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*StakingRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingRoleRevokedIterator{contract: _Staking.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Staking *StakingFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *StakingRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRoleRevoked)
				if err := _Staking.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Staking *StakingFilterer) ParseRoleRevoked(log types.Log) (*StakingRoleRevoked, error) {
	event := new(StakingRoleRevoked)
	if err := _Staking.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Staking contract.
type StakingUnpausedIterator struct {
	Event *StakingUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUnpaused represents a Unpaused event raised by the Staking contract.
type StakingUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Staking *StakingFilterer) FilterUnpaused(opts *bind.FilterOpts) (*StakingUnpausedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &StakingUnpausedIterator{contract: _Staking.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Staking *StakingFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *StakingUnpaused) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUnpaused)
				if err := _Staking.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Staking *StakingFilterer) ParseUnpaused(log types.Log) (*StakingUnpaused, error) {
	event := new(StakingUnpaused)
	if err := _Staking.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
