// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package l2

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

// DataTypesNode is an auto generated low-level Go binding around an user-defined struct.
type DataTypesNode struct {
	NodeId              *big.Int
	Account             common.Address
	TaxRateBasisPoints  uint64
	PublicGood          bool
	Alpha               bool
	Name                string
	Description         string
	OperationPoolTokens *big.Int
	StakingPoolTokens   *big.Int
	TotalShares         *big.Int
	SlashedTokens       *big.Int
}

// DataTypesUnstakeRequest is an auto generated low-level Go binding around an user-defined struct.
type DataTypesUnstakeRequest struct {
	Owner         common.Address
	NodeAddr      common.Address
	Timestamp     *big.Int
	UnstakeAmount *big.Int
}

// DataTypesWithdrawalRequest is an auto generated low-level Go binding around an user-defined struct.
type DataTypesWithdrawalRequest struct {
	Owner     common.Address
	Timestamp *big.Int
	Amount    *big.Int
}

// StakingMetaData contains all meta data concerning the Staking contract.
var StakingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"treasury\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakeRatio\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakeUnbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"depositUnbondingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nodeSlashRateBasisPoints\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"userSlashRateBasisPoints\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minTaxRateBasisPoints\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEPOSIT_UNBONDING_PERIOD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_DEPOSIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_TAX_RATE_BASIS_POINTS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NODE_SLASH_RATE_BASIS_POINTS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ORACLE_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAUSE_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SHARES_PER_CHIP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STAKE_RATIO\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STAKE_UNBONDING_PERIOD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TREASURY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"USER_SLASH_RATE_BASIS_POINTS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"chipsContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimUnstake\",\"inputs\":[{\"name\":\"requestIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimWithdrawal\",\"inputs\":[{\"name\":\"requestIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createNode\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"disableAlphaPhase\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"distributeRewards\",\"inputs\":[{\"name\":\"epochInfo\",\"type\":\"uint256[3]\",\"internalType\":\"uint256[3]\"},{\"name\":\"nodeAddrs\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"operationRewards\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"stakingRewards\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"requestCounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"publicPoolRewards\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getChipsInfo\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNode\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDataTypes.Node\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"alpha\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"operationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashedTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeAvatar\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodes\",\"inputs\":[{\"name\":\"nodeAddrs\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structDataTypes.Node[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"alpha\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"operationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashedTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodesWithPagination\",\"inputs\":[{\"name\":\"offset\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"limit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"nodes\",\"type\":\"tuple[]\",\"internalType\":\"structDataTypes.Node[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"alpha\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"operationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashedTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingUnstake\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDataTypes.UnstakeRequest\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unstakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingWithdrawal\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDataTypes.WithdrawalRequest\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint40\",\"internalType\":\"uint40\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPoolInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"totalOperationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalStakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublicPool\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDataTypes.Node\",\"components\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"alpha\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"operationPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakingPoolTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashedTokens\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleMember\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleMemberCount\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"chips\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pauseAccount\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracleAccount\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isAlphaPhase\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSettlementPhase\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minTokensToStake\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestUnstake\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"chipsIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestWithdrawal\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSettlementPhase\",\"inputs\":[{\"name\":\"enabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTaxRateBasisPoints4Node\",\"inputs\":[{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTaxRateBasisPoints4PublicPool\",\"inputs\":[{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashNodes\",\"inputs\":[{\"name\":\"nodeAddrs\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"startTokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"endTokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeToPublicPool\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"startTokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"endTokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateToPublicGood\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw2Treasury\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Deposited\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeCreated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"alpha\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeSlashed\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"slashedOperationPool\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"slashedStakingPool\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeTaxRateBasisPointsSet\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeUpdated2PublicGood\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PublicGoodRewardDistributed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"publicPoolRewards\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"publicPoolTax\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PublicPoolTaxRateBasisPointsSet\",\"inputs\":[{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardDistributed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"nodeAddrs\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"operationRewards\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"stakingRewards\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"taxAmounts\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"requestCounts\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Staked\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnstakeClaimed\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"unstakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnstakeRequested\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"unstakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"chipsIds\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawRequested\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalClaimed\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlphaWithdrawNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyClaimed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AmountTooSmall\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"BatchSizeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CallerNotStaking\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckpointUnorderedInsertion\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChipNotAuthorized\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ChipNotPublicGood\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ChipNotValid\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ChipsIdOverflow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChipsNotSameOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ClaimIdNotExists\",\"inputs\":[{\"name\":\"claimId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ClaimTimeNotReady\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CreateNodeToZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyChipsIds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyNodeList\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExcessWithdrawalAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidArrayLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidEpoch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidEpochNumber\",\"inputs\":[{\"name\":\"current\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"got\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTraitId\",\"inputs\":[{\"name\":\"traitId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NodeAlreadyPublicGood\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NodeExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeNotExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeNotPublicGood\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NodeStakedOrDeposited\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OperationRewardsExceed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PublicGoodNodeNotDeposited\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PublicGoodNodeNotInAlphaPhase\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RewardsAlreadyDistributed\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SettlementPhase\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakeToPublicGoodNode\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StakingRewardsExceed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SubmissionIntervalNotElapsed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TaxRateBasisPointsTooLarge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TaxRateBasisPointsTooSmall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]}]",
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

// MINDEPOSIT is a free data retrieval call binding the contract method 0xe1e158a5.
//
// Solidity: function MIN_DEPOSIT() view returns(uint256)
func (_Staking *StakingCaller) MINDEPOSIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "MIN_DEPOSIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINDEPOSIT is a free data retrieval call binding the contract method 0xe1e158a5.
//
// Solidity: function MIN_DEPOSIT() view returns(uint256)
func (_Staking *StakingSession) MINDEPOSIT() (*big.Int, error) {
	return _Staking.Contract.MINDEPOSIT(&_Staking.CallOpts)
}

// MINDEPOSIT is a free data retrieval call binding the contract method 0xe1e158a5.
//
// Solidity: function MIN_DEPOSIT() view returns(uint256)
func (_Staking *StakingCallerSession) MINDEPOSIT() (*big.Int, error) {
	return _Staking.Contract.MINDEPOSIT(&_Staking.CallOpts)
}

// MINTAXRATEBASISPOINTS is a free data retrieval call binding the contract method 0x2fe3a2a0.
//
// Solidity: function MIN_TAX_RATE_BASIS_POINTS() view returns(uint256)
func (_Staking *StakingCaller) MINTAXRATEBASISPOINTS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "MIN_TAX_RATE_BASIS_POINTS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINTAXRATEBASISPOINTS is a free data retrieval call binding the contract method 0x2fe3a2a0.
//
// Solidity: function MIN_TAX_RATE_BASIS_POINTS() view returns(uint256)
func (_Staking *StakingSession) MINTAXRATEBASISPOINTS() (*big.Int, error) {
	return _Staking.Contract.MINTAXRATEBASISPOINTS(&_Staking.CallOpts)
}

// MINTAXRATEBASISPOINTS is a free data retrieval call binding the contract method 0x2fe3a2a0.
//
// Solidity: function MIN_TAX_RATE_BASIS_POINTS() view returns(uint256)
func (_Staking *StakingCallerSession) MINTAXRATEBASISPOINTS() (*big.Int, error) {
	return _Staking.Contract.MINTAXRATEBASISPOINTS(&_Staking.CallOpts)
}

// NODESLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0x3daf051f.
//
// Solidity: function NODE_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Staking *StakingCaller) NODESLASHRATEBASISPOINTS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "NODE_SLASH_RATE_BASIS_POINTS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NODESLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0x3daf051f.
//
// Solidity: function NODE_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Staking *StakingSession) NODESLASHRATEBASISPOINTS() (*big.Int, error) {
	return _Staking.Contract.NODESLASHRATEBASISPOINTS(&_Staking.CallOpts)
}

// NODESLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0x3daf051f.
//
// Solidity: function NODE_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Staking *StakingCallerSession) NODESLASHRATEBASISPOINTS() (*big.Int, error) {
	return _Staking.Contract.NODESLASHRATEBASISPOINTS(&_Staking.CallOpts)
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

// SHARESPERCHIP is a free data retrieval call binding the contract method 0x6b05f6dc.
//
// Solidity: function SHARES_PER_CHIP() view returns(uint256)
func (_Staking *StakingCaller) SHARESPERCHIP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "SHARES_PER_CHIP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SHARESPERCHIP is a free data retrieval call binding the contract method 0x6b05f6dc.
//
// Solidity: function SHARES_PER_CHIP() view returns(uint256)
func (_Staking *StakingSession) SHARESPERCHIP() (*big.Int, error) {
	return _Staking.Contract.SHARESPERCHIP(&_Staking.CallOpts)
}

// SHARESPERCHIP is a free data retrieval call binding the contract method 0x6b05f6dc.
//
// Solidity: function SHARES_PER_CHIP() view returns(uint256)
func (_Staking *StakingCallerSession) SHARESPERCHIP() (*big.Int, error) {
	return _Staking.Contract.SHARESPERCHIP(&_Staking.CallOpts)
}

// STAKERATIO is a free data retrieval call binding the contract method 0x736fcdf6.
//
// Solidity: function STAKE_RATIO() view returns(uint256)
func (_Staking *StakingCaller) STAKERATIO(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "STAKE_RATIO")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// STAKERATIO is a free data retrieval call binding the contract method 0x736fcdf6.
//
// Solidity: function STAKE_RATIO() view returns(uint256)
func (_Staking *StakingSession) STAKERATIO() (*big.Int, error) {
	return _Staking.Contract.STAKERATIO(&_Staking.CallOpts)
}

// STAKERATIO is a free data retrieval call binding the contract method 0x736fcdf6.
//
// Solidity: function STAKE_RATIO() view returns(uint256)
func (_Staking *StakingCallerSession) STAKERATIO() (*big.Int, error) {
	return _Staking.Contract.STAKERATIO(&_Staking.CallOpts)
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

// USERSLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0xb47d343c.
//
// Solidity: function USER_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Staking *StakingCaller) USERSLASHRATEBASISPOINTS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "USER_SLASH_RATE_BASIS_POINTS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// USERSLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0xb47d343c.
//
// Solidity: function USER_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Staking *StakingSession) USERSLASHRATEBASISPOINTS() (*big.Int, error) {
	return _Staking.Contract.USERSLASHRATEBASISPOINTS(&_Staking.CallOpts)
}

// USERSLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0xb47d343c.
//
// Solidity: function USER_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Staking *StakingCallerSession) USERSLASHRATEBASISPOINTS() (*big.Int, error) {
	return _Staking.Contract.USERSLASHRATEBASISPOINTS(&_Staking.CallOpts)
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

// GetChipsInfo is a free data retrieval call binding the contract method 0x90d3f47c.
//
// Solidity: function getChipsInfo(uint256 tokenId) view returns(address nodeAddr, uint256 tokens)
func (_Staking *StakingCaller) GetChipsInfo(opts *bind.CallOpts, tokenId *big.Int) (struct {
	NodeAddr common.Address
	Tokens   *big.Int
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getChipsInfo", tokenId)

	outstruct := new(struct {
		NodeAddr common.Address
		Tokens   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NodeAddr = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Tokens = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetChipsInfo is a free data retrieval call binding the contract method 0x90d3f47c.
//
// Solidity: function getChipsInfo(uint256 tokenId) view returns(address nodeAddr, uint256 tokens)
func (_Staking *StakingSession) GetChipsInfo(tokenId *big.Int) (struct {
	NodeAddr common.Address
	Tokens   *big.Int
}, error) {
	return _Staking.Contract.GetChipsInfo(&_Staking.CallOpts, tokenId)
}

// GetChipsInfo is a free data retrieval call binding the contract method 0x90d3f47c.
//
// Solidity: function getChipsInfo(uint256 tokenId) view returns(address nodeAddr, uint256 tokens)
func (_Staking *StakingCallerSession) GetChipsInfo(tokenId *big.Int) (struct {
	NodeAddr common.Address
	Tokens   *big.Int
}, error) {
	return _Staking.Contract.GetChipsInfo(&_Staking.CallOpts, tokenId)
}

// GetNode is a free data retrieval call binding the contract method 0x9d209048.
//
// Solidity: function getNode(address nodeAddr) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256))
func (_Staking *StakingCaller) GetNode(opts *bind.CallOpts, nodeAddr common.Address) (DataTypesNode, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getNode", nodeAddr)

	if err != nil {
		return *new(DataTypesNode), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesNode)).(*DataTypesNode)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0x9d209048.
//
// Solidity: function getNode(address nodeAddr) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256))
func (_Staking *StakingSession) GetNode(nodeAddr common.Address) (DataTypesNode, error) {
	return _Staking.Contract.GetNode(&_Staking.CallOpts, nodeAddr)
}

// GetNode is a free data retrieval call binding the contract method 0x9d209048.
//
// Solidity: function getNode(address nodeAddr) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256))
func (_Staking *StakingCallerSession) GetNode(nodeAddr common.Address) (DataTypesNode, error) {
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
// Solidity: function getNodes(address[] nodeAddrs) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Staking *StakingCaller) GetNodes(opts *bind.CallOpts, nodeAddrs []common.Address) ([]DataTypesNode, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getNodes", nodeAddrs)

	if err != nil {
		return *new([]DataTypesNode), err
	}

	out0 := *abi.ConvertType(out[0], new([]DataTypesNode)).(*[]DataTypesNode)

	return out0, err

}

// GetNodes is a free data retrieval call binding the contract method 0x38c96b14.
//
// Solidity: function getNodes(address[] nodeAddrs) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Staking *StakingSession) GetNodes(nodeAddrs []common.Address) ([]DataTypesNode, error) {
	return _Staking.Contract.GetNodes(&_Staking.CallOpts, nodeAddrs)
}

// GetNodes is a free data retrieval call binding the contract method 0x38c96b14.
//
// Solidity: function getNodes(address[] nodeAddrs) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Staking *StakingCallerSession) GetNodes(nodeAddrs []common.Address) ([]DataTypesNode, error) {
	return _Staking.Contract.GetNodes(&_Staking.CallOpts, nodeAddrs)
}

// GetNodesWithPagination is a free data retrieval call binding the contract method 0xd995415b.
//
// Solidity: function getNodesWithPagination(uint256 offset, uint256 limit) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Staking *StakingCaller) GetNodesWithPagination(opts *bind.CallOpts, offset *big.Int, limit *big.Int) ([]DataTypesNode, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getNodesWithPagination", offset, limit)

	if err != nil {
		return *new([]DataTypesNode), err
	}

	out0 := *abi.ConvertType(out[0], new([]DataTypesNode)).(*[]DataTypesNode)

	return out0, err

}

// GetNodesWithPagination is a free data retrieval call binding the contract method 0xd995415b.
//
// Solidity: function getNodesWithPagination(uint256 offset, uint256 limit) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Staking *StakingSession) GetNodesWithPagination(offset *big.Int, limit *big.Int) ([]DataTypesNode, error) {
	return _Staking.Contract.GetNodesWithPagination(&_Staking.CallOpts, offset, limit)
}

// GetNodesWithPagination is a free data retrieval call binding the contract method 0xd995415b.
//
// Solidity: function getNodesWithPagination(uint256 offset, uint256 limit) view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Staking *StakingCallerSession) GetNodesWithPagination(offset *big.Int, limit *big.Int) ([]DataTypesNode, error) {
	return _Staking.Contract.GetNodesWithPagination(&_Staking.CallOpts, offset, limit)
}

// GetPendingUnstake is a free data retrieval call binding the contract method 0xadfd065f.
//
// Solidity: function getPendingUnstake(uint256 requestId) view returns((address,address,uint256,uint256))
func (_Staking *StakingCaller) GetPendingUnstake(opts *bind.CallOpts, requestId *big.Int) (DataTypesUnstakeRequest, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getPendingUnstake", requestId)

	if err != nil {
		return *new(DataTypesUnstakeRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesUnstakeRequest)).(*DataTypesUnstakeRequest)

	return out0, err

}

// GetPendingUnstake is a free data retrieval call binding the contract method 0xadfd065f.
//
// Solidity: function getPendingUnstake(uint256 requestId) view returns((address,address,uint256,uint256))
func (_Staking *StakingSession) GetPendingUnstake(requestId *big.Int) (DataTypesUnstakeRequest, error) {
	return _Staking.Contract.GetPendingUnstake(&_Staking.CallOpts, requestId)
}

// GetPendingUnstake is a free data retrieval call binding the contract method 0xadfd065f.
//
// Solidity: function getPendingUnstake(uint256 requestId) view returns((address,address,uint256,uint256))
func (_Staking *StakingCallerSession) GetPendingUnstake(requestId *big.Int) (DataTypesUnstakeRequest, error) {
	return _Staking.Contract.GetPendingUnstake(&_Staking.CallOpts, requestId)
}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x38a3c878.
//
// Solidity: function getPendingWithdrawal(uint256 requestId) view returns((address,uint40,uint256))
func (_Staking *StakingCaller) GetPendingWithdrawal(opts *bind.CallOpts, requestId *big.Int) (DataTypesWithdrawalRequest, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getPendingWithdrawal", requestId)

	if err != nil {
		return *new(DataTypesWithdrawalRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesWithdrawalRequest)).(*DataTypesWithdrawalRequest)

	return out0, err

}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x38a3c878.
//
// Solidity: function getPendingWithdrawal(uint256 requestId) view returns((address,uint40,uint256))
func (_Staking *StakingSession) GetPendingWithdrawal(requestId *big.Int) (DataTypesWithdrawalRequest, error) {
	return _Staking.Contract.GetPendingWithdrawal(&_Staking.CallOpts, requestId)
}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x38a3c878.
//
// Solidity: function getPendingWithdrawal(uint256 requestId) view returns((address,uint40,uint256))
func (_Staking *StakingCallerSession) GetPendingWithdrawal(requestId *big.Int) (DataTypesWithdrawalRequest, error) {
	return _Staking.Contract.GetPendingWithdrawal(&_Staking.CallOpts, requestId)
}

// GetPoolInfo is a free data retrieval call binding the contract method 0x60246c88.
//
// Solidity: function getPoolInfo() view returns(uint256 totalOperationPoolTokens, uint256 totalStakingPoolTokens)
func (_Staking *StakingCaller) GetPoolInfo(opts *bind.CallOpts) (struct {
	TotalOperationPoolTokens *big.Int
	TotalStakingPoolTokens   *big.Int
}, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getPoolInfo")

	outstruct := new(struct {
		TotalOperationPoolTokens *big.Int
		TotalStakingPoolTokens   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalOperationPoolTokens = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalStakingPoolTokens = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPoolInfo is a free data retrieval call binding the contract method 0x60246c88.
//
// Solidity: function getPoolInfo() view returns(uint256 totalOperationPoolTokens, uint256 totalStakingPoolTokens)
func (_Staking *StakingSession) GetPoolInfo() (struct {
	TotalOperationPoolTokens *big.Int
	TotalStakingPoolTokens   *big.Int
}, error) {
	return _Staking.Contract.GetPoolInfo(&_Staking.CallOpts)
}

// GetPoolInfo is a free data retrieval call binding the contract method 0x60246c88.
//
// Solidity: function getPoolInfo() view returns(uint256 totalOperationPoolTokens, uint256 totalStakingPoolTokens)
func (_Staking *StakingCallerSession) GetPoolInfo() (struct {
	TotalOperationPoolTokens *big.Int
	TotalStakingPoolTokens   *big.Int
}, error) {
	return _Staking.Contract.GetPoolInfo(&_Staking.CallOpts)
}

// GetPublicPool is a free data retrieval call binding the contract method 0xc84c42a3.
//
// Solidity: function getPublicPool() view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256))
func (_Staking *StakingCaller) GetPublicPool(opts *bind.CallOpts) (DataTypesNode, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getPublicPool")

	if err != nil {
		return *new(DataTypesNode), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesNode)).(*DataTypesNode)

	return out0, err

}

// GetPublicPool is a free data retrieval call binding the contract method 0xc84c42a3.
//
// Solidity: function getPublicPool() view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256))
func (_Staking *StakingSession) GetPublicPool() (DataTypesNode, error) {
	return _Staking.Contract.GetPublicPool(&_Staking.CallOpts)
}

// GetPublicPool is a free data retrieval call binding the contract method 0xc84c42a3.
//
// Solidity: function getPublicPool() view returns((uint256,address,uint64,bool,bool,string,string,uint256,uint256,uint256,uint256))
func (_Staking *StakingCallerSession) GetPublicPool() (DataTypesNode, error) {
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

// MinTokensToStake is a free data retrieval call binding the contract method 0x14936b13.
//
// Solidity: function minTokensToStake(address nodeAddr) view returns(uint256)
func (_Staking *StakingCaller) MinTokensToStake(opts *bind.CallOpts, nodeAddr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "minTokensToStake", nodeAddr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinTokensToStake is a free data retrieval call binding the contract method 0x14936b13.
//
// Solidity: function minTokensToStake(address nodeAddr) view returns(uint256)
func (_Staking *StakingSession) MinTokensToStake(nodeAddr common.Address) (*big.Int, error) {
	return _Staking.Contract.MinTokensToStake(&_Staking.CallOpts, nodeAddr)
}

// MinTokensToStake is a free data retrieval call binding the contract method 0x14936b13.
//
// Solidity: function minTokensToStake(address nodeAddr) view returns(uint256)
func (_Staking *StakingCallerSession) MinTokensToStake(nodeAddr common.Address) (*big.Int, error) {
	return _Staking.Contract.MinTokensToStake(&_Staking.CallOpts, nodeAddr)
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

// DisableAlphaPhase is a paid mutator transaction binding the contract method 0x81001e60.
//
// Solidity: function disableAlphaPhase() returns()
func (_Staking *StakingTransactor) DisableAlphaPhase(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "disableAlphaPhase")
}

// DisableAlphaPhase is a paid mutator transaction binding the contract method 0x81001e60.
//
// Solidity: function disableAlphaPhase() returns()
func (_Staking *StakingSession) DisableAlphaPhase() (*types.Transaction, error) {
	return _Staking.Contract.DisableAlphaPhase(&_Staking.TransactOpts)
}

// DisableAlphaPhase is a paid mutator transaction binding the contract method 0x81001e60.
//
// Solidity: function disableAlphaPhase() returns()
func (_Staking *StakingTransactorSession) DisableAlphaPhase() (*types.Transaction, error) {
	return _Staking.Contract.DisableAlphaPhase(&_Staking.TransactOpts)
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

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address chips, address pauseAccount, address oracleAccount) returns()
func (_Staking *StakingTransactor) Initialize(opts *bind.TransactOpts, chips common.Address, pauseAccount common.Address, oracleAccount common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "initialize", chips, pauseAccount, oracleAccount)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address chips, address pauseAccount, address oracleAccount) returns()
func (_Staking *StakingSession) Initialize(chips common.Address, pauseAccount common.Address, oracleAccount common.Address) (*types.Transaction, error) {
	return _Staking.Contract.Initialize(&_Staking.TransactOpts, chips, pauseAccount, oracleAccount)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address chips, address pauseAccount, address oracleAccount) returns()
func (_Staking *StakingTransactorSession) Initialize(chips common.Address, pauseAccount common.Address, oracleAccount common.Address) (*types.Transaction, error) {
	return _Staking.Contract.Initialize(&_Staking.TransactOpts, chips, pauseAccount, oracleAccount)
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
// Solidity: function requestUnstake(address nodeAddr, uint256[] chipsIds) returns(uint256 requestId)
func (_Staking *StakingTransactor) RequestUnstake(opts *bind.TransactOpts, nodeAddr common.Address, chipsIds []*big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "requestUnstake", nodeAddr, chipsIds)
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xbcdd4190.
//
// Solidity: function requestUnstake(address nodeAddr, uint256[] chipsIds) returns(uint256 requestId)
func (_Staking *StakingSession) RequestUnstake(nodeAddr common.Address, chipsIds []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.RequestUnstake(&_Staking.TransactOpts, nodeAddr, chipsIds)
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xbcdd4190.
//
// Solidity: function requestUnstake(address nodeAddr, uint256[] chipsIds) returns(uint256 requestId)
func (_Staking *StakingTransactorSession) RequestUnstake(nodeAddr common.Address, chipsIds []*big.Int) (*types.Transaction, error) {
	return _Staking.Contract.RequestUnstake(&_Staking.TransactOpts, nodeAddr, chipsIds)
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

// SlashNodes is a paid mutator transaction binding the contract method 0xa2f641c3.
//
// Solidity: function slashNodes(address[] nodeAddrs) returns()
func (_Staking *StakingTransactor) SlashNodes(opts *bind.TransactOpts, nodeAddrs []common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "slashNodes", nodeAddrs)
}

// SlashNodes is a paid mutator transaction binding the contract method 0xa2f641c3.
//
// Solidity: function slashNodes(address[] nodeAddrs) returns()
func (_Staking *StakingSession) SlashNodes(nodeAddrs []common.Address) (*types.Transaction, error) {
	return _Staking.Contract.SlashNodes(&_Staking.TransactOpts, nodeAddrs)
}

// SlashNodes is a paid mutator transaction binding the contract method 0xa2f641c3.
//
// Solidity: function slashNodes(address[] nodeAddrs) returns()
func (_Staking *StakingTransactorSession) SlashNodes(nodeAddrs []common.Address) (*types.Transaction, error) {
	return _Staking.Contract.SlashNodes(&_Staking.TransactOpts, nodeAddrs)
}

// Stake is a paid mutator transaction binding the contract method 0x26476204.
//
// Solidity: function stake(address nodeAddr) payable returns(uint256 startTokenId, uint256 endTokenId)
func (_Staking *StakingTransactor) Stake(opts *bind.TransactOpts, nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "stake", nodeAddr)
}

// Stake is a paid mutator transaction binding the contract method 0x26476204.
//
// Solidity: function stake(address nodeAddr) payable returns(uint256 startTokenId, uint256 endTokenId)
func (_Staking *StakingSession) Stake(nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.Stake(&_Staking.TransactOpts, nodeAddr)
}

// Stake is a paid mutator transaction binding the contract method 0x26476204.
//
// Solidity: function stake(address nodeAddr) payable returns(uint256 startTokenId, uint256 endTokenId)
func (_Staking *StakingTransactorSession) Stake(nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.Stake(&_Staking.TransactOpts, nodeAddr)
}

// StakeToPublicPool is a paid mutator transaction binding the contract method 0x379f8100.
//
// Solidity: function stakeToPublicPool(address nodeAddr) payable returns(uint256 startTokenId, uint256 endTokenId)
func (_Staking *StakingTransactor) StakeToPublicPool(opts *bind.TransactOpts, nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "stakeToPublicPool", nodeAddr)
}

// StakeToPublicPool is a paid mutator transaction binding the contract method 0x379f8100.
//
// Solidity: function stakeToPublicPool(address nodeAddr) payable returns(uint256 startTokenId, uint256 endTokenId)
func (_Staking *StakingSession) StakeToPublicPool(nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.StakeToPublicPool(&_Staking.TransactOpts, nodeAddr)
}

// StakeToPublicPool is a paid mutator transaction binding the contract method 0x379f8100.
//
// Solidity: function stakeToPublicPool(address nodeAddr) payable returns(uint256 startTokenId, uint256 endTokenId)
func (_Staking *StakingTransactorSession) StakeToPublicPool(nodeAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.StakeToPublicPool(&_Staking.TransactOpts, nodeAddr)
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

// UpdateToPublicGood is a paid mutator transaction binding the contract method 0xc9af094c.
//
// Solidity: function updateToPublicGood() returns()
func (_Staking *StakingTransactor) UpdateToPublicGood(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "updateToPublicGood")
}

// UpdateToPublicGood is a paid mutator transaction binding the contract method 0xc9af094c.
//
// Solidity: function updateToPublicGood() returns()
func (_Staking *StakingSession) UpdateToPublicGood() (*types.Transaction, error) {
	return _Staking.Contract.UpdateToPublicGood(&_Staking.TransactOpts)
}

// UpdateToPublicGood is a paid mutator transaction binding the contract method 0xc9af094c.
//
// Solidity: function updateToPublicGood() returns()
func (_Staking *StakingTransactorSession) UpdateToPublicGood() (*types.Transaction, error) {
	return _Staking.Contract.UpdateToPublicGood(&_Staking.TransactOpts)
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

// StakingDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the Staking contract.
type StakingDepositedIterator struct {
	Event *StakingDeposited // Event containing the contract specifics and raw log

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
func (it *StakingDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDeposited)
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
		it.Event = new(StakingDeposited)
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
func (it *StakingDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDeposited represents a Deposited event raised by the Staking contract.
type StakingDeposited struct {
	NodeAddr common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed nodeAddr, uint256 indexed amount)
func (_Staking *StakingFilterer) FilterDeposited(opts *bind.FilterOpts, nodeAddr []common.Address, amount []*big.Int) (*StakingDepositedIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Deposited", nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakingDepositedIterator{contract: _Staking.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed nodeAddr, uint256 indexed amount)
func (_Staking *StakingFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *StakingDeposited, nodeAddr []common.Address, amount []*big.Int) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Deposited", nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDeposited)
				if err := _Staking.contract.UnpackLog(event, "Deposited", log); err != nil {
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

// ParseDeposited is a log parse operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed nodeAddr, uint256 indexed amount)
func (_Staking *StakingFilterer) ParseDeposited(log types.Log) (*StakingDeposited, error) {
	event := new(StakingDeposited)
	if err := _Staking.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

// StakingNodeCreatedIterator is returned from FilterNodeCreated and is used to iterate over the raw logs and unpacked data for NodeCreated events raised by the Staking contract.
type StakingNodeCreatedIterator struct {
	Event *StakingNodeCreated // Event containing the contract specifics and raw log

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
func (it *StakingNodeCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingNodeCreated)
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
		it.Event = new(StakingNodeCreated)
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
func (it *StakingNodeCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingNodeCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingNodeCreated represents a NodeCreated event raised by the Staking contract.
type StakingNodeCreated struct {
	NodeId             *big.Int
	NodeAddr           common.Address
	Name               string
	Description        string
	TaxRateBasisPoints uint64
	PublicGood         bool
	Alpha              bool
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterNodeCreated is a free log retrieval operation binding the contract event 0x37570f68d94fd46cd4009b3823da2b2bc1a9a7e38f824f311ede9e876816e321.
//
// Solidity: event NodeCreated(uint256 indexed nodeId, address indexed nodeAddr, string name, string description, uint64 taxRateBasisPoints, bool publicGood, bool alpha)
func (_Staking *StakingFilterer) FilterNodeCreated(opts *bind.FilterOpts, nodeId []*big.Int, nodeAddr []common.Address) (*StakingNodeCreatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "NodeCreated", nodeIdRule, nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return &StakingNodeCreatedIterator{contract: _Staking.contract, event: "NodeCreated", logs: logs, sub: sub}, nil
}

// WatchNodeCreated is a free log subscription operation binding the contract event 0x37570f68d94fd46cd4009b3823da2b2bc1a9a7e38f824f311ede9e876816e321.
//
// Solidity: event NodeCreated(uint256 indexed nodeId, address indexed nodeAddr, string name, string description, uint64 taxRateBasisPoints, bool publicGood, bool alpha)
func (_Staking *StakingFilterer) WatchNodeCreated(opts *bind.WatchOpts, sink chan<- *StakingNodeCreated, nodeId []*big.Int, nodeAddr []common.Address) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "NodeCreated", nodeIdRule, nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingNodeCreated)
				if err := _Staking.contract.UnpackLog(event, "NodeCreated", log); err != nil {
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

// ParseNodeCreated is a log parse operation binding the contract event 0x37570f68d94fd46cd4009b3823da2b2bc1a9a7e38f824f311ede9e876816e321.
//
// Solidity: event NodeCreated(uint256 indexed nodeId, address indexed nodeAddr, string name, string description, uint64 taxRateBasisPoints, bool publicGood, bool alpha)
func (_Staking *StakingFilterer) ParseNodeCreated(log types.Log) (*StakingNodeCreated, error) {
	event := new(StakingNodeCreated)
	if err := _Staking.contract.UnpackLog(event, "NodeCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingNodeSlashedIterator is returned from FilterNodeSlashed and is used to iterate over the raw logs and unpacked data for NodeSlashed events raised by the Staking contract.
type StakingNodeSlashedIterator struct {
	Event *StakingNodeSlashed // Event containing the contract specifics and raw log

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
func (it *StakingNodeSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingNodeSlashed)
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
		it.Event = new(StakingNodeSlashed)
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
func (it *StakingNodeSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingNodeSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingNodeSlashed represents a NodeSlashed event raised by the Staking contract.
type StakingNodeSlashed struct {
	NodeAddr             common.Address
	SlashedOperationPool *big.Int
	SlashedStakingPool   *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterNodeSlashed is a free log retrieval operation binding the contract event 0xa8d720d0a0a2e7c96bf9eb87433901ebb6331356c8f3283b2568de34478703cc.
//
// Solidity: event NodeSlashed(address indexed nodeAddr, uint256 indexed slashedOperationPool, uint256 indexed slashedStakingPool)
func (_Staking *StakingFilterer) FilterNodeSlashed(opts *bind.FilterOpts, nodeAddr []common.Address, slashedOperationPool []*big.Int, slashedStakingPool []*big.Int) (*StakingNodeSlashedIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var slashedOperationPoolRule []interface{}
	for _, slashedOperationPoolItem := range slashedOperationPool {
		slashedOperationPoolRule = append(slashedOperationPoolRule, slashedOperationPoolItem)
	}
	var slashedStakingPoolRule []interface{}
	for _, slashedStakingPoolItem := range slashedStakingPool {
		slashedStakingPoolRule = append(slashedStakingPoolRule, slashedStakingPoolItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "NodeSlashed", nodeAddrRule, slashedOperationPoolRule, slashedStakingPoolRule)
	if err != nil {
		return nil, err
	}
	return &StakingNodeSlashedIterator{contract: _Staking.contract, event: "NodeSlashed", logs: logs, sub: sub}, nil
}

// WatchNodeSlashed is a free log subscription operation binding the contract event 0xa8d720d0a0a2e7c96bf9eb87433901ebb6331356c8f3283b2568de34478703cc.
//
// Solidity: event NodeSlashed(address indexed nodeAddr, uint256 indexed slashedOperationPool, uint256 indexed slashedStakingPool)
func (_Staking *StakingFilterer) WatchNodeSlashed(opts *bind.WatchOpts, sink chan<- *StakingNodeSlashed, nodeAddr []common.Address, slashedOperationPool []*big.Int, slashedStakingPool []*big.Int) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var slashedOperationPoolRule []interface{}
	for _, slashedOperationPoolItem := range slashedOperationPool {
		slashedOperationPoolRule = append(slashedOperationPoolRule, slashedOperationPoolItem)
	}
	var slashedStakingPoolRule []interface{}
	for _, slashedStakingPoolItem := range slashedStakingPool {
		slashedStakingPoolRule = append(slashedStakingPoolRule, slashedStakingPoolItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "NodeSlashed", nodeAddrRule, slashedOperationPoolRule, slashedStakingPoolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingNodeSlashed)
				if err := _Staking.contract.UnpackLog(event, "NodeSlashed", log); err != nil {
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

// ParseNodeSlashed is a log parse operation binding the contract event 0xa8d720d0a0a2e7c96bf9eb87433901ebb6331356c8f3283b2568de34478703cc.
//
// Solidity: event NodeSlashed(address indexed nodeAddr, uint256 indexed slashedOperationPool, uint256 indexed slashedStakingPool)
func (_Staking *StakingFilterer) ParseNodeSlashed(log types.Log) (*StakingNodeSlashed, error) {
	event := new(StakingNodeSlashed)
	if err := _Staking.contract.UnpackLog(event, "NodeSlashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingNodeTaxRateBasisPointsSetIterator is returned from FilterNodeTaxRateBasisPointsSet and is used to iterate over the raw logs and unpacked data for NodeTaxRateBasisPointsSet events raised by the Staking contract.
type StakingNodeTaxRateBasisPointsSetIterator struct {
	Event *StakingNodeTaxRateBasisPointsSet // Event containing the contract specifics and raw log

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
func (it *StakingNodeTaxRateBasisPointsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingNodeTaxRateBasisPointsSet)
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
		it.Event = new(StakingNodeTaxRateBasisPointsSet)
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
func (it *StakingNodeTaxRateBasisPointsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingNodeTaxRateBasisPointsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingNodeTaxRateBasisPointsSet represents a NodeTaxRateBasisPointsSet event raised by the Staking contract.
type StakingNodeTaxRateBasisPointsSet struct {
	NodeAddr           common.Address
	TaxRateBasisPoints uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterNodeTaxRateBasisPointsSet is a free log retrieval operation binding the contract event 0xb8e5551053b871a40f7c7382e5bd3af5a62dd737d059d3838cf3aa7c325bd479.
//
// Solidity: event NodeTaxRateBasisPointsSet(address indexed nodeAddr, uint64 indexed taxRateBasisPoints)
func (_Staking *StakingFilterer) FilterNodeTaxRateBasisPointsSet(opts *bind.FilterOpts, nodeAddr []common.Address, taxRateBasisPoints []uint64) (*StakingNodeTaxRateBasisPointsSetIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "NodeTaxRateBasisPointsSet", nodeAddrRule, taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return &StakingNodeTaxRateBasisPointsSetIterator{contract: _Staking.contract, event: "NodeTaxRateBasisPointsSet", logs: logs, sub: sub}, nil
}

// WatchNodeTaxRateBasisPointsSet is a free log subscription operation binding the contract event 0xb8e5551053b871a40f7c7382e5bd3af5a62dd737d059d3838cf3aa7c325bd479.
//
// Solidity: event NodeTaxRateBasisPointsSet(address indexed nodeAddr, uint64 indexed taxRateBasisPoints)
func (_Staking *StakingFilterer) WatchNodeTaxRateBasisPointsSet(opts *bind.WatchOpts, sink chan<- *StakingNodeTaxRateBasisPointsSet, nodeAddr []common.Address, taxRateBasisPoints []uint64) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "NodeTaxRateBasisPointsSet", nodeAddrRule, taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingNodeTaxRateBasisPointsSet)
				if err := _Staking.contract.UnpackLog(event, "NodeTaxRateBasisPointsSet", log); err != nil {
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

// ParseNodeTaxRateBasisPointsSet is a log parse operation binding the contract event 0xb8e5551053b871a40f7c7382e5bd3af5a62dd737d059d3838cf3aa7c325bd479.
//
// Solidity: event NodeTaxRateBasisPointsSet(address indexed nodeAddr, uint64 indexed taxRateBasisPoints)
func (_Staking *StakingFilterer) ParseNodeTaxRateBasisPointsSet(log types.Log) (*StakingNodeTaxRateBasisPointsSet, error) {
	event := new(StakingNodeTaxRateBasisPointsSet)
	if err := _Staking.contract.UnpackLog(event, "NodeTaxRateBasisPointsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingNodeUpdated2PublicGoodIterator is returned from FilterNodeUpdated2PublicGood and is used to iterate over the raw logs and unpacked data for NodeUpdated2PublicGood events raised by the Staking contract.
type StakingNodeUpdated2PublicGoodIterator struct {
	Event *StakingNodeUpdated2PublicGood // Event containing the contract specifics and raw log

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
func (it *StakingNodeUpdated2PublicGoodIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingNodeUpdated2PublicGood)
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
		it.Event = new(StakingNodeUpdated2PublicGood)
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
func (it *StakingNodeUpdated2PublicGoodIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingNodeUpdated2PublicGoodIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingNodeUpdated2PublicGood represents a NodeUpdated2PublicGood event raised by the Staking contract.
type StakingNodeUpdated2PublicGood struct {
	NodeAddr common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNodeUpdated2PublicGood is a free log retrieval operation binding the contract event 0x86538b79bef9c52dbe4e888742cfbd70114655c47ef30ede4997791fe79a9376.
//
// Solidity: event NodeUpdated2PublicGood(address indexed nodeAddr)
func (_Staking *StakingFilterer) FilterNodeUpdated2PublicGood(opts *bind.FilterOpts, nodeAddr []common.Address) (*StakingNodeUpdated2PublicGoodIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "NodeUpdated2PublicGood", nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return &StakingNodeUpdated2PublicGoodIterator{contract: _Staking.contract, event: "NodeUpdated2PublicGood", logs: logs, sub: sub}, nil
}

// WatchNodeUpdated2PublicGood is a free log subscription operation binding the contract event 0x86538b79bef9c52dbe4e888742cfbd70114655c47ef30ede4997791fe79a9376.
//
// Solidity: event NodeUpdated2PublicGood(address indexed nodeAddr)
func (_Staking *StakingFilterer) WatchNodeUpdated2PublicGood(opts *bind.WatchOpts, sink chan<- *StakingNodeUpdated2PublicGood, nodeAddr []common.Address) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "NodeUpdated2PublicGood", nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingNodeUpdated2PublicGood)
				if err := _Staking.contract.UnpackLog(event, "NodeUpdated2PublicGood", log); err != nil {
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

// ParseNodeUpdated2PublicGood is a log parse operation binding the contract event 0x86538b79bef9c52dbe4e888742cfbd70114655c47ef30ede4997791fe79a9376.
//
// Solidity: event NodeUpdated2PublicGood(address indexed nodeAddr)
func (_Staking *StakingFilterer) ParseNodeUpdated2PublicGood(log types.Log) (*StakingNodeUpdated2PublicGood, error) {
	event := new(StakingNodeUpdated2PublicGood)
	if err := _Staking.contract.UnpackLog(event, "NodeUpdated2PublicGood", log); err != nil {
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

// StakingPublicPoolTaxRateBasisPointsSetIterator is returned from FilterPublicPoolTaxRateBasisPointsSet and is used to iterate over the raw logs and unpacked data for PublicPoolTaxRateBasisPointsSet events raised by the Staking contract.
type StakingPublicPoolTaxRateBasisPointsSetIterator struct {
	Event *StakingPublicPoolTaxRateBasisPointsSet // Event containing the contract specifics and raw log

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
func (it *StakingPublicPoolTaxRateBasisPointsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPublicPoolTaxRateBasisPointsSet)
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
		it.Event = new(StakingPublicPoolTaxRateBasisPointsSet)
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
func (it *StakingPublicPoolTaxRateBasisPointsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPublicPoolTaxRateBasisPointsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPublicPoolTaxRateBasisPointsSet represents a PublicPoolTaxRateBasisPointsSet event raised by the Staking contract.
type StakingPublicPoolTaxRateBasisPointsSet struct {
	TaxRateBasisPoints uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterPublicPoolTaxRateBasisPointsSet is a free log retrieval operation binding the contract event 0x948cf2302b029d76db2ac06e4ef2625e6c687335de349317468f47942a44e8b0.
//
// Solidity: event PublicPoolTaxRateBasisPointsSet(uint64 indexed taxRateBasisPoints)
func (_Staking *StakingFilterer) FilterPublicPoolTaxRateBasisPointsSet(opts *bind.FilterOpts, taxRateBasisPoints []uint64) (*StakingPublicPoolTaxRateBasisPointsSetIterator, error) {

	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "PublicPoolTaxRateBasisPointsSet", taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return &StakingPublicPoolTaxRateBasisPointsSetIterator{contract: _Staking.contract, event: "PublicPoolTaxRateBasisPointsSet", logs: logs, sub: sub}, nil
}

// WatchPublicPoolTaxRateBasisPointsSet is a free log subscription operation binding the contract event 0x948cf2302b029d76db2ac06e4ef2625e6c687335de349317468f47942a44e8b0.
//
// Solidity: event PublicPoolTaxRateBasisPointsSet(uint64 indexed taxRateBasisPoints)
func (_Staking *StakingFilterer) WatchPublicPoolTaxRateBasisPointsSet(opts *bind.WatchOpts, sink chan<- *StakingPublicPoolTaxRateBasisPointsSet, taxRateBasisPoints []uint64) (event.Subscription, error) {

	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "PublicPoolTaxRateBasisPointsSet", taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPublicPoolTaxRateBasisPointsSet)
				if err := _Staking.contract.UnpackLog(event, "PublicPoolTaxRateBasisPointsSet", log); err != nil {
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

// ParsePublicPoolTaxRateBasisPointsSet is a log parse operation binding the contract event 0x948cf2302b029d76db2ac06e4ef2625e6c687335de349317468f47942a44e8b0.
//
// Solidity: event PublicPoolTaxRateBasisPointsSet(uint64 indexed taxRateBasisPoints)
func (_Staking *StakingFilterer) ParsePublicPoolTaxRateBasisPointsSet(log types.Log) (*StakingPublicPoolTaxRateBasisPointsSet, error) {
	event := new(StakingPublicPoolTaxRateBasisPointsSet)
	if err := _Staking.contract.UnpackLog(event, "PublicPoolTaxRateBasisPointsSet", log); err != nil {
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
	TaxAmounts       []*big.Int
	RequestCounts    []*big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRewardDistributed is a free log retrieval operation binding the contract event 0x8ea79f19e90b084c2009d3490a097547d8bbb315a883b9efec0996502c1dd7ae.
//
// Solidity: event RewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, address[] nodeAddrs, uint256[] operationRewards, uint256[] stakingRewards, uint256[] taxAmounts, uint256[] requestCounts)
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
// Solidity: event RewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, address[] nodeAddrs, uint256[] operationRewards, uint256[] stakingRewards, uint256[] taxAmounts, uint256[] requestCounts)
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
// Solidity: event RewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, address[] nodeAddrs, uint256[] operationRewards, uint256[] stakingRewards, uint256[] taxAmounts, uint256[] requestCounts)
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

// StakingStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Staking contract.
type StakingStakedIterator struct {
	Event *StakingStaked // Event containing the contract specifics and raw log

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
func (it *StakingStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingStaked)
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
		it.Event = new(StakingStaked)
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
func (it *StakingStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingStaked represents a Staked event raised by the Staking contract.
type StakingStaked struct {
	User         common.Address
	NodeAddr     common.Address
	Amount       *big.Int
	StartTokenId *big.Int
	EndTokenId   *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0xad3fa07f4195b47e64892eb944ecbfc253384053c119852bb2bcae484c2fcb69.
//
// Solidity: event Staked(address indexed user, address indexed nodeAddr, uint256 indexed amount, uint256 startTokenId, uint256 endTokenId)
func (_Staking *StakingFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address, nodeAddr []common.Address, amount []*big.Int) (*StakingStakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Staked", userRule, nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakingStakedIterator{contract: _Staking.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0xad3fa07f4195b47e64892eb944ecbfc253384053c119852bb2bcae484c2fcb69.
//
// Solidity: event Staked(address indexed user, address indexed nodeAddr, uint256 indexed amount, uint256 startTokenId, uint256 endTokenId)
func (_Staking *StakingFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *StakingStaked, user []common.Address, nodeAddr []common.Address, amount []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Staked", userRule, nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingStaked)
				if err := _Staking.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0xad3fa07f4195b47e64892eb944ecbfc253384053c119852bb2bcae484c2fcb69.
//
// Solidity: event Staked(address indexed user, address indexed nodeAddr, uint256 indexed amount, uint256 startTokenId, uint256 endTokenId)
func (_Staking *StakingFilterer) ParseStaked(log types.Log) (*StakingStaked, error) {
	event := new(StakingStaked)
	if err := _Staking.contract.UnpackLog(event, "Staked", log); err != nil {
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

// StakingUnstakeClaimedIterator is returned from FilterUnstakeClaimed and is used to iterate over the raw logs and unpacked data for UnstakeClaimed events raised by the Staking contract.
type StakingUnstakeClaimedIterator struct {
	Event *StakingUnstakeClaimed // Event containing the contract specifics and raw log

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
func (it *StakingUnstakeClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUnstakeClaimed)
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
		it.Event = new(StakingUnstakeClaimed)
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
func (it *StakingUnstakeClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUnstakeClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUnstakeClaimed represents a UnstakeClaimed event raised by the Staking contract.
type StakingUnstakeClaimed struct {
	RequestId     *big.Int
	NodeAddr      common.Address
	User          common.Address
	UnstakeAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUnstakeClaimed is a free log retrieval operation binding the contract event 0x2769ece66eadb650afd8c08c7a8772e39381dddd7230f9e039669e631044d47c.
//
// Solidity: event UnstakeClaimed(uint256 indexed requestId, address indexed nodeAddr, address indexed user, uint256 unstakeAmount)
func (_Staking *StakingFilterer) FilterUnstakeClaimed(opts *bind.FilterOpts, requestId []*big.Int, nodeAddr []common.Address, user []common.Address) (*StakingUnstakeClaimedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "UnstakeClaimed", requestIdRule, nodeAddrRule, userRule)
	if err != nil {
		return nil, err
	}
	return &StakingUnstakeClaimedIterator{contract: _Staking.contract, event: "UnstakeClaimed", logs: logs, sub: sub}, nil
}

// WatchUnstakeClaimed is a free log subscription operation binding the contract event 0x2769ece66eadb650afd8c08c7a8772e39381dddd7230f9e039669e631044d47c.
//
// Solidity: event UnstakeClaimed(uint256 indexed requestId, address indexed nodeAddr, address indexed user, uint256 unstakeAmount)
func (_Staking *StakingFilterer) WatchUnstakeClaimed(opts *bind.WatchOpts, sink chan<- *StakingUnstakeClaimed, requestId []*big.Int, nodeAddr []common.Address, user []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "UnstakeClaimed", requestIdRule, nodeAddrRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUnstakeClaimed)
				if err := _Staking.contract.UnpackLog(event, "UnstakeClaimed", log); err != nil {
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

// ParseUnstakeClaimed is a log parse operation binding the contract event 0x2769ece66eadb650afd8c08c7a8772e39381dddd7230f9e039669e631044d47c.
//
// Solidity: event UnstakeClaimed(uint256 indexed requestId, address indexed nodeAddr, address indexed user, uint256 unstakeAmount)
func (_Staking *StakingFilterer) ParseUnstakeClaimed(log types.Log) (*StakingUnstakeClaimed, error) {
	event := new(StakingUnstakeClaimed)
	if err := _Staking.contract.UnpackLog(event, "UnstakeClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUnstakeRequestedIterator is returned from FilterUnstakeRequested and is used to iterate over the raw logs and unpacked data for UnstakeRequested events raised by the Staking contract.
type StakingUnstakeRequestedIterator struct {
	Event *StakingUnstakeRequested // Event containing the contract specifics and raw log

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
func (it *StakingUnstakeRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUnstakeRequested)
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
		it.Event = new(StakingUnstakeRequested)
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
func (it *StakingUnstakeRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUnstakeRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUnstakeRequested represents a UnstakeRequested event raised by the Staking contract.
type StakingUnstakeRequested struct {
	User          common.Address
	NodeAddr      common.Address
	RequestId     *big.Int
	UnstakeAmount *big.Int
	ChipsIds      []*big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUnstakeRequested is a free log retrieval operation binding the contract event 0x2808f92d5a0fada467cbe4e766f62f323e78271a7471058a87ef63a9e3e4c5c5.
//
// Solidity: event UnstakeRequested(address indexed user, address indexed nodeAddr, uint256 indexed requestId, uint256 unstakeAmount, uint256[] chipsIds)
func (_Staking *StakingFilterer) FilterUnstakeRequested(opts *bind.FilterOpts, user []common.Address, nodeAddr []common.Address, requestId []*big.Int) (*StakingUnstakeRequestedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "UnstakeRequested", userRule, nodeAddrRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &StakingUnstakeRequestedIterator{contract: _Staking.contract, event: "UnstakeRequested", logs: logs, sub: sub}, nil
}

// WatchUnstakeRequested is a free log subscription operation binding the contract event 0x2808f92d5a0fada467cbe4e766f62f323e78271a7471058a87ef63a9e3e4c5c5.
//
// Solidity: event UnstakeRequested(address indexed user, address indexed nodeAddr, uint256 indexed requestId, uint256 unstakeAmount, uint256[] chipsIds)
func (_Staking *StakingFilterer) WatchUnstakeRequested(opts *bind.WatchOpts, sink chan<- *StakingUnstakeRequested, user []common.Address, nodeAddr []common.Address, requestId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "UnstakeRequested", userRule, nodeAddrRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUnstakeRequested)
				if err := _Staking.contract.UnpackLog(event, "UnstakeRequested", log); err != nil {
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

// ParseUnstakeRequested is a log parse operation binding the contract event 0x2808f92d5a0fada467cbe4e766f62f323e78271a7471058a87ef63a9e3e4c5c5.
//
// Solidity: event UnstakeRequested(address indexed user, address indexed nodeAddr, uint256 indexed requestId, uint256 unstakeAmount, uint256[] chipsIds)
func (_Staking *StakingFilterer) ParseUnstakeRequested(log types.Log) (*StakingUnstakeRequested, error) {
	event := new(StakingUnstakeRequested)
	if err := _Staking.contract.UnpackLog(event, "UnstakeRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingWithdrawRequestedIterator is returned from FilterWithdrawRequested and is used to iterate over the raw logs and unpacked data for WithdrawRequested events raised by the Staking contract.
type StakingWithdrawRequestedIterator struct {
	Event *StakingWithdrawRequested // Event containing the contract specifics and raw log

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
func (it *StakingWithdrawRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingWithdrawRequested)
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
		it.Event = new(StakingWithdrawRequested)
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
func (it *StakingWithdrawRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingWithdrawRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingWithdrawRequested represents a WithdrawRequested event raised by the Staking contract.
type StakingWithdrawRequested struct {
	NodeAddr  common.Address
	Amount    *big.Int
	RequestId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawRequested is a free log retrieval operation binding the contract event 0xd72eb5d043f24a0168ae744d5c44f9596fd673a26bf74d9646bff4b844882d14.
//
// Solidity: event WithdrawRequested(address indexed nodeAddr, uint256 indexed amount, uint256 indexed requestId)
func (_Staking *StakingFilterer) FilterWithdrawRequested(opts *bind.FilterOpts, nodeAddr []common.Address, amount []*big.Int, requestId []*big.Int) (*StakingWithdrawRequestedIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "WithdrawRequested", nodeAddrRule, amountRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &StakingWithdrawRequestedIterator{contract: _Staking.contract, event: "WithdrawRequested", logs: logs, sub: sub}, nil
}

// WatchWithdrawRequested is a free log subscription operation binding the contract event 0xd72eb5d043f24a0168ae744d5c44f9596fd673a26bf74d9646bff4b844882d14.
//
// Solidity: event WithdrawRequested(address indexed nodeAddr, uint256 indexed amount, uint256 indexed requestId)
func (_Staking *StakingFilterer) WatchWithdrawRequested(opts *bind.WatchOpts, sink chan<- *StakingWithdrawRequested, nodeAddr []common.Address, amount []*big.Int, requestId []*big.Int) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "WithdrawRequested", nodeAddrRule, amountRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingWithdrawRequested)
				if err := _Staking.contract.UnpackLog(event, "WithdrawRequested", log); err != nil {
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

// ParseWithdrawRequested is a log parse operation binding the contract event 0xd72eb5d043f24a0168ae744d5c44f9596fd673a26bf74d9646bff4b844882d14.
//
// Solidity: event WithdrawRequested(address indexed nodeAddr, uint256 indexed amount, uint256 indexed requestId)
func (_Staking *StakingFilterer) ParseWithdrawRequested(log types.Log) (*StakingWithdrawRequested, error) {
	event := new(StakingWithdrawRequested)
	if err := _Staking.contract.UnpackLog(event, "WithdrawRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingWithdrawalClaimedIterator is returned from FilterWithdrawalClaimed and is used to iterate over the raw logs and unpacked data for WithdrawalClaimed events raised by the Staking contract.
type StakingWithdrawalClaimedIterator struct {
	Event *StakingWithdrawalClaimed // Event containing the contract specifics and raw log

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
func (it *StakingWithdrawalClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingWithdrawalClaimed)
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
		it.Event = new(StakingWithdrawalClaimed)
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
func (it *StakingWithdrawalClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingWithdrawalClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingWithdrawalClaimed represents a WithdrawalClaimed event raised by the Staking contract.
type StakingWithdrawalClaimed struct {
	RequestId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalClaimed is a free log retrieval operation binding the contract event 0x8772d6f79a1845a0c0e90ef18d99f91242bbc0ba98c9ca780feaad42b81f02ba.
//
// Solidity: event WithdrawalClaimed(uint256 indexed requestId)
func (_Staking *StakingFilterer) FilterWithdrawalClaimed(opts *bind.FilterOpts, requestId []*big.Int) (*StakingWithdrawalClaimedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "WithdrawalClaimed", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &StakingWithdrawalClaimedIterator{contract: _Staking.contract, event: "WithdrawalClaimed", logs: logs, sub: sub}, nil
}

// WatchWithdrawalClaimed is a free log subscription operation binding the contract event 0x8772d6f79a1845a0c0e90ef18d99f91242bbc0ba98c9ca780feaad42b81f02ba.
//
// Solidity: event WithdrawalClaimed(uint256 indexed requestId)
func (_Staking *StakingFilterer) WatchWithdrawalClaimed(opts *bind.WatchOpts, sink chan<- *StakingWithdrawalClaimed, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "WithdrawalClaimed", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingWithdrawalClaimed)
				if err := _Staking.contract.UnpackLog(event, "WithdrawalClaimed", log); err != nil {
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

// ParseWithdrawalClaimed is a log parse operation binding the contract event 0x8772d6f79a1845a0c0e90ef18d99f91242bbc0ba98c9ca780feaad42b81f02ba.
//
// Solidity: event WithdrawalClaimed(uint256 indexed requestId)
func (_Staking *StakingFilterer) ParseWithdrawalClaimed(log types.Log) (*StakingWithdrawalClaimed, error) {
	event := new(StakingWithdrawalClaimed)
	if err := _Staking.contract.UnpackLog(event, "WithdrawalClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
