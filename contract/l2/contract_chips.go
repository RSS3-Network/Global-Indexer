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

// ChipsMetaData contains all meta data concerning the Chips contract.
var ChipsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"stakeRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakeUnbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositUnbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nodeSlashRateBasisPoints\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"userSlashRateBasisPoints\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minDeposit\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyClaimed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AmountTooSmall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchSizeZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CallerNotNodeOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CallerNotStaking\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CheckpointUnorderedInsertion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ChipNotAuthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ChipNotPublicGood\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"}],\"name\":\"ChipNotValid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ChipsIdOverflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"claimId\",\"type\":\"uint256\"}],\"name\":\"ClaimIdNotExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimTimeNotReady\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CreateNodeToZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DepositedTokensSlashedAll\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyNodeList\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidArrayLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"InvalidEpoch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NodeExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NodeNotExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"}],\"name\":\"NodeNotPublicGood\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NodeStakedOrDeposited\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PublicGoodNodeNotDeposited\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"}],\"name\":\"PublicGoodNodeNotStaked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RewardDistributionFailed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"bits\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"SafeCastOverflowedUintDowncast\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TaxRateBasisPointsTooLarge\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"publicGood\",\"type\":\"bool\"}],\"name\":\"NodeCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"}],\"name\":\"NodeDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"slashedOperationPool\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"slashedStakingPool\",\"type\":\"uint256\"}],\"name\":\"NodeSlashed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"}],\"name\":\"NodeTaxRateBasisPointsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"publicPoolReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"publicPoolTax\",\"type\":\"uint256\"}],\"name\":\"PublicGoodRewardDistributed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"}],\"name\":\"PublicPoolTaxRateBasisPointsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"nodeAddrs\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"requestFees\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"operationRewards\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"stakingRewards\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"taxAmounts\",\"type\":\"uint256[]\"}],\"name\":\"RewardDistributed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTokenId\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unstakeAmount\",\"type\":\"uint256\"}],\"name\":\"UnstakeClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"chipsIds\",\"type\":\"uint256[]\"}],\"name\":\"UnstakeRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"WithdrawRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"WithdrawalClaimed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPOSIT_UNBONDING_PERIOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_DEPOSIT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NODE_SLASH_RATE_BASIS_POINTS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ORACLE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SHARES_PER_CHIP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"STAKE_RATIO\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"STAKE_UNBONDING_PERIOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TOKEN\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TREASURY\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"USER_SLASH_RATE_BASIS_POINTS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chipsContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"requestIds\",\"type\":\"uint256[]\"}],\"name\":\"claimUnstake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"requestIds\",\"type\":\"uint256[]\"}],\"name\":\"claimWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"publicGood\",\"type\":\"bool\"}],\"name\":\"createNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"publicGood\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"createNodeAndDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"}],\"name\":\"deleteNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[3]\",\"name\":\"epochInfo\",\"type\":\"uint256[3]\"},{\"internalType\":\"address[]\",\"name\":\"nodeAddrs\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"requestFees\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"operationRewards\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingRewards\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"publicPoolReward\",\"type\":\"uint256\"}],\"name\":\"distributeRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getChipsInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"}],\"name\":\"getNode\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"publicGood\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"operationPoolTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakingPoolTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashedTokens\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.Node\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeAddrs\",\"type\":\"address[]\"}],\"name\":\"getNodes\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"publicGood\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"operationPoolTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakingPoolTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashedTokens\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getNodesWithPagination\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"publicGood\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"operationPoolTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakingPoolTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashedTokens\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.Node[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"getPendingUnstake\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unstakeAmount\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.UnstakeRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"getPendingWithdrawal\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint40\",\"name\":\"timestamp\",\"type\":\"uint40\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.WithdrawalRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPoolInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalOperationPoolTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalStakingPoolTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"treasuryAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPublicPool\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"publicGood\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"operationPoolTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakingPoolTokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashedTokens\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.Node\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"chips\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"pauseAccount\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAccount\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"}],\"name\":\"minTokensToStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"chipsIds\",\"type\":\"uint256[]\"}],\"name\":\"requestUnstake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"chipsIds\",\"type\":\"uint256[]\"}],\"name\":\"requestUnstakeFromPublicPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"requestWithdrawal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"}],\"name\":\"setTaxRateBasisPoints4Node\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\"}],\"name\":\"setTaxRateBasisPoints4PublicPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeAddrs\",\"type\":\"address[]\"}],\"name\":\"slashNodes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"startTokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTokenId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stakeToPublicPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"startTokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTokenId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakingToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw2Treasury\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ChipsABI is the input ABI used to generate the binding from.
// Deprecated: Use ChipsMetaData.ABI instead.
var ChipsABI = ChipsMetaData.ABI

// Chips is an auto generated Go binding around an Ethereum contract.
type Chips struct {
	ChipsCaller     // Read-only binding to the contract
	ChipsTransactor // Write-only binding to the contract
	ChipsFilterer   // Log filterer for contract events
}

// ChipsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChipsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChipsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChipsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChipsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChipsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChipsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChipsSession struct {
	Contract     *Chips            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChipsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChipsCallerSession struct {
	Contract *ChipsCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ChipsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChipsTransactorSession struct {
	Contract     *ChipsTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChipsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChipsRaw struct {
	Contract *Chips // Generic contract binding to access the raw methods on
}

// ChipsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChipsCallerRaw struct {
	Contract *ChipsCaller // Generic read-only contract binding to access the raw methods on
}

// ChipsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChipsTransactorRaw struct {
	Contract *ChipsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChips creates a new instance of Chips, bound to a specific deployed contract.
func NewChips(address common.Address, backend bind.ContractBackend) (*Chips, error) {
	contract, err := bindChips(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Chips{ChipsCaller: ChipsCaller{contract: contract}, ChipsTransactor: ChipsTransactor{contract: contract}, ChipsFilterer: ChipsFilterer{contract: contract}}, nil
}

// NewChipsCaller creates a new read-only instance of Chips, bound to a specific deployed contract.
func NewChipsCaller(address common.Address, caller bind.ContractCaller) (*ChipsCaller, error) {
	contract, err := bindChips(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChipsCaller{contract: contract}, nil
}

// NewChipsTransactor creates a new write-only instance of Chips, bound to a specific deployed contract.
func NewChipsTransactor(address common.Address, transactor bind.ContractTransactor) (*ChipsTransactor, error) {
	contract, err := bindChips(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChipsTransactor{contract: contract}, nil
}

// NewChipsFilterer creates a new log filterer instance of Chips, bound to a specific deployed contract.
func NewChipsFilterer(address common.Address, filterer bind.ContractFilterer) (*ChipsFilterer, error) {
	contract, err := bindChips(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChipsFilterer{contract: contract}, nil
}

// bindChips binds a generic wrapper to an already deployed contract.
func bindChips(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChipsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chips *ChipsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chips.Contract.ChipsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chips *ChipsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chips.Contract.ChipsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chips *ChipsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chips.Contract.ChipsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chips *ChipsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chips.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chips *ChipsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chips.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chips *ChipsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chips.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Chips *ChipsCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Chips *ChipsSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Chips.Contract.DEFAULTADMINROLE(&_Chips.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Chips *ChipsCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Chips.Contract.DEFAULTADMINROLE(&_Chips.CallOpts)
}

// DEPOSITUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x6bdc11d5.
//
// Solidity: function DEPOSIT_UNBONDING_PERIOD() view returns(uint256)
func (_Chips *ChipsCaller) DEPOSITUNBONDINGPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "DEPOSIT_UNBONDING_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEPOSITUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x6bdc11d5.
//
// Solidity: function DEPOSIT_UNBONDING_PERIOD() view returns(uint256)
func (_Chips *ChipsSession) DEPOSITUNBONDINGPERIOD() (*big.Int, error) {
	return _Chips.Contract.DEPOSITUNBONDINGPERIOD(&_Chips.CallOpts)
}

// DEPOSITUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x6bdc11d5.
//
// Solidity: function DEPOSIT_UNBONDING_PERIOD() view returns(uint256)
func (_Chips *ChipsCallerSession) DEPOSITUNBONDINGPERIOD() (*big.Int, error) {
	return _Chips.Contract.DEPOSITUNBONDINGPERIOD(&_Chips.CallOpts)
}

// MINDEPOSIT is a free data retrieval call binding the contract method 0xe1e158a5.
//
// Solidity: function MIN_DEPOSIT() view returns(uint256)
func (_Chips *ChipsCaller) MINDEPOSIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "MIN_DEPOSIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINDEPOSIT is a free data retrieval call binding the contract method 0xe1e158a5.
//
// Solidity: function MIN_DEPOSIT() view returns(uint256)
func (_Chips *ChipsSession) MINDEPOSIT() (*big.Int, error) {
	return _Chips.Contract.MINDEPOSIT(&_Chips.CallOpts)
}

// MINDEPOSIT is a free data retrieval call binding the contract method 0xe1e158a5.
//
// Solidity: function MIN_DEPOSIT() view returns(uint256)
func (_Chips *ChipsCallerSession) MINDEPOSIT() (*big.Int, error) {
	return _Chips.Contract.MINDEPOSIT(&_Chips.CallOpts)
}

// NODESLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0x3daf051f.
//
// Solidity: function NODE_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Chips *ChipsCaller) NODESLASHRATEBASISPOINTS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "NODE_SLASH_RATE_BASIS_POINTS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NODESLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0x3daf051f.
//
// Solidity: function NODE_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Chips *ChipsSession) NODESLASHRATEBASISPOINTS() (*big.Int, error) {
	return _Chips.Contract.NODESLASHRATEBASISPOINTS(&_Chips.CallOpts)
}

// NODESLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0x3daf051f.
//
// Solidity: function NODE_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Chips *ChipsCallerSession) NODESLASHRATEBASISPOINTS() (*big.Int, error) {
	return _Chips.Contract.NODESLASHRATEBASISPOINTS(&_Chips.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Chips *ChipsCaller) ORACLEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "ORACLE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Chips *ChipsSession) ORACLEROLE() ([32]byte, error) {
	return _Chips.Contract.ORACLEROLE(&_Chips.CallOpts)
}

// ORACLEROLE is a free data retrieval call binding the contract method 0x07e2cea5.
//
// Solidity: function ORACLE_ROLE() view returns(bytes32)
func (_Chips *ChipsCallerSession) ORACLEROLE() ([32]byte, error) {
	return _Chips.Contract.ORACLEROLE(&_Chips.CallOpts)
}

// PAUSEROLE is a free data retrieval call binding the contract method 0x389ed267.
//
// Solidity: function PAUSE_ROLE() view returns(bytes32)
func (_Chips *ChipsCaller) PAUSEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "PAUSE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSEROLE is a free data retrieval call binding the contract method 0x389ed267.
//
// Solidity: function PAUSE_ROLE() view returns(bytes32)
func (_Chips *ChipsSession) PAUSEROLE() ([32]byte, error) {
	return _Chips.Contract.PAUSEROLE(&_Chips.CallOpts)
}

// PAUSEROLE is a free data retrieval call binding the contract method 0x389ed267.
//
// Solidity: function PAUSE_ROLE() view returns(bytes32)
func (_Chips *ChipsCallerSession) PAUSEROLE() ([32]byte, error) {
	return _Chips.Contract.PAUSEROLE(&_Chips.CallOpts)
}

// SHARESPERCHIP is a free data retrieval call binding the contract method 0x6b05f6dc.
//
// Solidity: function SHARES_PER_CHIP() view returns(uint256)
func (_Chips *ChipsCaller) SHARESPERCHIP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "SHARES_PER_CHIP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SHARESPERCHIP is a free data retrieval call binding the contract method 0x6b05f6dc.
//
// Solidity: function SHARES_PER_CHIP() view returns(uint256)
func (_Chips *ChipsSession) SHARESPERCHIP() (*big.Int, error) {
	return _Chips.Contract.SHARESPERCHIP(&_Chips.CallOpts)
}

// SHARESPERCHIP is a free data retrieval call binding the contract method 0x6b05f6dc.
//
// Solidity: function SHARES_PER_CHIP() view returns(uint256)
func (_Chips *ChipsCallerSession) SHARESPERCHIP() (*big.Int, error) {
	return _Chips.Contract.SHARESPERCHIP(&_Chips.CallOpts)
}

// STAKERATIO is a free data retrieval call binding the contract method 0x736fcdf6.
//
// Solidity: function STAKE_RATIO() view returns(uint256)
func (_Chips *ChipsCaller) STAKERATIO(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "STAKE_RATIO")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// STAKERATIO is a free data retrieval call binding the contract method 0x736fcdf6.
//
// Solidity: function STAKE_RATIO() view returns(uint256)
func (_Chips *ChipsSession) STAKERATIO() (*big.Int, error) {
	return _Chips.Contract.STAKERATIO(&_Chips.CallOpts)
}

// STAKERATIO is a free data retrieval call binding the contract method 0x736fcdf6.
//
// Solidity: function STAKE_RATIO() view returns(uint256)
func (_Chips *ChipsCallerSession) STAKERATIO() (*big.Int, error) {
	return _Chips.Contract.STAKERATIO(&_Chips.CallOpts)
}

// STAKEUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x2606a44a.
//
// Solidity: function STAKE_UNBONDING_PERIOD() view returns(uint256)
func (_Chips *ChipsCaller) STAKEUNBONDINGPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "STAKE_UNBONDING_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// STAKEUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x2606a44a.
//
// Solidity: function STAKE_UNBONDING_PERIOD() view returns(uint256)
func (_Chips *ChipsSession) STAKEUNBONDINGPERIOD() (*big.Int, error) {
	return _Chips.Contract.STAKEUNBONDINGPERIOD(&_Chips.CallOpts)
}

// STAKEUNBONDINGPERIOD is a free data retrieval call binding the contract method 0x2606a44a.
//
// Solidity: function STAKE_UNBONDING_PERIOD() view returns(uint256)
func (_Chips *ChipsCallerSession) STAKEUNBONDINGPERIOD() (*big.Int, error) {
	return _Chips.Contract.STAKEUNBONDINGPERIOD(&_Chips.CallOpts)
}

// TOKEN is a free data retrieval call binding the contract method 0x82bfefc8.
//
// Solidity: function TOKEN() view returns(address)
func (_Chips *ChipsCaller) TOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "TOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TOKEN is a free data retrieval call binding the contract method 0x82bfefc8.
//
// Solidity: function TOKEN() view returns(address)
func (_Chips *ChipsSession) TOKEN() (common.Address, error) {
	return _Chips.Contract.TOKEN(&_Chips.CallOpts)
}

// TOKEN is a free data retrieval call binding the contract method 0x82bfefc8.
//
// Solidity: function TOKEN() view returns(address)
func (_Chips *ChipsCallerSession) TOKEN() (common.Address, error) {
	return _Chips.Contract.TOKEN(&_Chips.CallOpts)
}

// TREASURY is a free data retrieval call binding the contract method 0x2d2c5565.
//
// Solidity: function TREASURY() view returns(address)
func (_Chips *ChipsCaller) TREASURY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "TREASURY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TREASURY is a free data retrieval call binding the contract method 0x2d2c5565.
//
// Solidity: function TREASURY() view returns(address)
func (_Chips *ChipsSession) TREASURY() (common.Address, error) {
	return _Chips.Contract.TREASURY(&_Chips.CallOpts)
}

// TREASURY is a free data retrieval call binding the contract method 0x2d2c5565.
//
// Solidity: function TREASURY() view returns(address)
func (_Chips *ChipsCallerSession) TREASURY() (common.Address, error) {
	return _Chips.Contract.TREASURY(&_Chips.CallOpts)
}

// USERSLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0xb47d343c.
//
// Solidity: function USER_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Chips *ChipsCaller) USERSLASHRATEBASISPOINTS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "USER_SLASH_RATE_BASIS_POINTS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// USERSLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0xb47d343c.
//
// Solidity: function USER_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Chips *ChipsSession) USERSLASHRATEBASISPOINTS() (*big.Int, error) {
	return _Chips.Contract.USERSLASHRATEBASISPOINTS(&_Chips.CallOpts)
}

// USERSLASHRATEBASISPOINTS is a free data retrieval call binding the contract method 0xb47d343c.
//
// Solidity: function USER_SLASH_RATE_BASIS_POINTS() view returns(uint256)
func (_Chips *ChipsCallerSession) USERSLASHRATEBASISPOINTS() (*big.Int, error) {
	return _Chips.Contract.USERSLASHRATEBASISPOINTS(&_Chips.CallOpts)
}

// ChipsContract is a free data retrieval call binding the contract method 0xd13b19a3.
//
// Solidity: function chipsContract() view returns(address)
func (_Chips *ChipsCaller) ChipsContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "chipsContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ChipsContract is a free data retrieval call binding the contract method 0xd13b19a3.
//
// Solidity: function chipsContract() view returns(address)
func (_Chips *ChipsSession) ChipsContract() (common.Address, error) {
	return _Chips.Contract.ChipsContract(&_Chips.CallOpts)
}

// ChipsContract is a free data retrieval call binding the contract method 0xd13b19a3.
//
// Solidity: function chipsContract() view returns(address)
func (_Chips *ChipsCallerSession) ChipsContract() (common.Address, error) {
	return _Chips.Contract.ChipsContract(&_Chips.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Chips *ChipsCaller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "currentEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Chips *ChipsSession) CurrentEpoch() (*big.Int, error) {
	return _Chips.Contract.CurrentEpoch(&_Chips.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Chips *ChipsCallerSession) CurrentEpoch() (*big.Int, error) {
	return _Chips.Contract.CurrentEpoch(&_Chips.CallOpts)
}

// GetChipsInfo is a free data retrieval call binding the contract method 0x90d3f47c.
//
// Solidity: function getChipsInfo(uint256 tokenId) view returns(address nodeAddr, uint256 tokens)
func (_Chips *ChipsCaller) GetChipsInfo(opts *bind.CallOpts, tokenId *big.Int) (struct {
	NodeAddr common.Address
	Tokens   *big.Int
}, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getChipsInfo", tokenId)

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
func (_Chips *ChipsSession) GetChipsInfo(tokenId *big.Int) (struct {
	NodeAddr common.Address
	Tokens   *big.Int
}, error) {
	return _Chips.Contract.GetChipsInfo(&_Chips.CallOpts, tokenId)
}

// GetChipsInfo is a free data retrieval call binding the contract method 0x90d3f47c.
//
// Solidity: function getChipsInfo(uint256 tokenId) view returns(address nodeAddr, uint256 tokens)
func (_Chips *ChipsCallerSession) GetChipsInfo(tokenId *big.Int) (struct {
	NodeAddr common.Address
	Tokens   *big.Int
}, error) {
	return _Chips.Contract.GetChipsInfo(&_Chips.CallOpts, tokenId)
}

// GetMinDeposit is a free data retrieval call binding the contract method 0x0eaad3f1.
//
// Solidity: function getMinDeposit() view returns(uint256)
func (_Chips *ChipsCaller) GetMinDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getMinDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinDeposit is a free data retrieval call binding the contract method 0x0eaad3f1.
//
// Solidity: function getMinDeposit() view returns(uint256)
func (_Chips *ChipsSession) GetMinDeposit() (*big.Int, error) {
	return _Chips.Contract.GetMinDeposit(&_Chips.CallOpts)
}

// GetMinDeposit is a free data retrieval call binding the contract method 0x0eaad3f1.
//
// Solidity: function getMinDeposit() view returns(uint256)
func (_Chips *ChipsCallerSession) GetMinDeposit() (*big.Int, error) {
	return _Chips.Contract.GetMinDeposit(&_Chips.CallOpts)
}

// GetNode is a free data retrieval call binding the contract method 0x9d209048.
//
// Solidity: function getNode(address nodeAddr) view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256))
func (_Chips *ChipsCaller) GetNode(opts *bind.CallOpts, nodeAddr common.Address) (DataTypesNode, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getNode", nodeAddr)

	if err != nil {
		return *new(DataTypesNode), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesNode)).(*DataTypesNode)

	return out0, err

}

// GetNode is a free data retrieval call binding the contract method 0x9d209048.
//
// Solidity: function getNode(address nodeAddr) view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256))
func (_Chips *ChipsSession) GetNode(nodeAddr common.Address) (DataTypesNode, error) {
	return _Chips.Contract.GetNode(&_Chips.CallOpts, nodeAddr)
}

// GetNode is a free data retrieval call binding the contract method 0x9d209048.
//
// Solidity: function getNode(address nodeAddr) view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256))
func (_Chips *ChipsCallerSession) GetNode(nodeAddr common.Address) (DataTypesNode, error) {
	return _Chips.Contract.GetNode(&_Chips.CallOpts, nodeAddr)
}

// GetNodeCount is a free data retrieval call binding the contract method 0x39bf397e.
//
// Solidity: function getNodeCount() view returns(uint256)
func (_Chips *ChipsCaller) GetNodeCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getNodeCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeCount is a free data retrieval call binding the contract method 0x39bf397e.
//
// Solidity: function getNodeCount() view returns(uint256)
func (_Chips *ChipsSession) GetNodeCount() (*big.Int, error) {
	return _Chips.Contract.GetNodeCount(&_Chips.CallOpts)
}

// GetNodeCount is a free data retrieval call binding the contract method 0x39bf397e.
//
// Solidity: function getNodeCount() view returns(uint256)
func (_Chips *ChipsCallerSession) GetNodeCount() (*big.Int, error) {
	return _Chips.Contract.GetNodeCount(&_Chips.CallOpts)
}

// GetNodes is a free data retrieval call binding the contract method 0x38c96b14.
//
// Solidity: function getNodes(address[] nodeAddrs) view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Chips *ChipsCaller) GetNodes(opts *bind.CallOpts, nodeAddrs []common.Address) ([]DataTypesNode, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getNodes", nodeAddrs)

	if err != nil {
		return *new([]DataTypesNode), err
	}

	out0 := *abi.ConvertType(out[0], new([]DataTypesNode)).(*[]DataTypesNode)

	return out0, err

}

// GetNodes is a free data retrieval call binding the contract method 0x38c96b14.
//
// Solidity: function getNodes(address[] nodeAddrs) view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Chips *ChipsSession) GetNodes(nodeAddrs []common.Address) ([]DataTypesNode, error) {
	return _Chips.Contract.GetNodes(&_Chips.CallOpts, nodeAddrs)
}

// GetNodes is a free data retrieval call binding the contract method 0x38c96b14.
//
// Solidity: function getNodes(address[] nodeAddrs) view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Chips *ChipsCallerSession) GetNodes(nodeAddrs []common.Address) ([]DataTypesNode, error) {
	return _Chips.Contract.GetNodes(&_Chips.CallOpts, nodeAddrs)
}

// GetNodesWithPagination is a free data retrieval call binding the contract method 0xd995415b.
//
// Solidity: function getNodesWithPagination(uint256 offset, uint256 limit) view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Chips *ChipsCaller) GetNodesWithPagination(opts *bind.CallOpts, offset *big.Int, limit *big.Int) ([]DataTypesNode, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getNodesWithPagination", offset, limit)

	if err != nil {
		return *new([]DataTypesNode), err
	}

	out0 := *abi.ConvertType(out[0], new([]DataTypesNode)).(*[]DataTypesNode)

	return out0, err

}

// GetNodesWithPagination is a free data retrieval call binding the contract method 0xd995415b.
//
// Solidity: function getNodesWithPagination(uint256 offset, uint256 limit) view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Chips *ChipsSession) GetNodesWithPagination(offset *big.Int, limit *big.Int) ([]DataTypesNode, error) {
	return _Chips.Contract.GetNodesWithPagination(&_Chips.CallOpts, offset, limit)
}

// GetNodesWithPagination is a free data retrieval call binding the contract method 0xd995415b.
//
// Solidity: function getNodesWithPagination(uint256 offset, uint256 limit) view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256)[] nodes)
func (_Chips *ChipsCallerSession) GetNodesWithPagination(offset *big.Int, limit *big.Int) ([]DataTypesNode, error) {
	return _Chips.Contract.GetNodesWithPagination(&_Chips.CallOpts, offset, limit)
}

// GetPendingUnstake is a free data retrieval call binding the contract method 0xadfd065f.
//
// Solidity: function getPendingUnstake(uint256 requestId) view returns((address,address,uint256,uint256))
func (_Chips *ChipsCaller) GetPendingUnstake(opts *bind.CallOpts, requestId *big.Int) (DataTypesUnstakeRequest, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getPendingUnstake", requestId)

	if err != nil {
		return *new(DataTypesUnstakeRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesUnstakeRequest)).(*DataTypesUnstakeRequest)

	return out0, err

}

// GetPendingUnstake is a free data retrieval call binding the contract method 0xadfd065f.
//
// Solidity: function getPendingUnstake(uint256 requestId) view returns((address,address,uint256,uint256))
func (_Chips *ChipsSession) GetPendingUnstake(requestId *big.Int) (DataTypesUnstakeRequest, error) {
	return _Chips.Contract.GetPendingUnstake(&_Chips.CallOpts, requestId)
}

// GetPendingUnstake is a free data retrieval call binding the contract method 0xadfd065f.
//
// Solidity: function getPendingUnstake(uint256 requestId) view returns((address,address,uint256,uint256))
func (_Chips *ChipsCallerSession) GetPendingUnstake(requestId *big.Int) (DataTypesUnstakeRequest, error) {
	return _Chips.Contract.GetPendingUnstake(&_Chips.CallOpts, requestId)
}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x38a3c878.
//
// Solidity: function getPendingWithdrawal(uint256 requestId) view returns((address,uint40,uint256))
func (_Chips *ChipsCaller) GetPendingWithdrawal(opts *bind.CallOpts, requestId *big.Int) (DataTypesWithdrawalRequest, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getPendingWithdrawal", requestId)

	if err != nil {
		return *new(DataTypesWithdrawalRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesWithdrawalRequest)).(*DataTypesWithdrawalRequest)

	return out0, err

}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x38a3c878.
//
// Solidity: function getPendingWithdrawal(uint256 requestId) view returns((address,uint40,uint256))
func (_Chips *ChipsSession) GetPendingWithdrawal(requestId *big.Int) (DataTypesWithdrawalRequest, error) {
	return _Chips.Contract.GetPendingWithdrawal(&_Chips.CallOpts, requestId)
}

// GetPendingWithdrawal is a free data retrieval call binding the contract method 0x38a3c878.
//
// Solidity: function getPendingWithdrawal(uint256 requestId) view returns((address,uint40,uint256))
func (_Chips *ChipsCallerSession) GetPendingWithdrawal(requestId *big.Int) (DataTypesWithdrawalRequest, error) {
	return _Chips.Contract.GetPendingWithdrawal(&_Chips.CallOpts, requestId)
}

// GetPoolInfo is a free data retrieval call binding the contract method 0x60246c88.
//
// Solidity: function getPoolInfo() view returns(uint256 totalOperationPoolTokens, uint256 totalStakingPoolTokens, uint256 treasuryAmount)
func (_Chips *ChipsCaller) GetPoolInfo(opts *bind.CallOpts) (struct {
	TotalOperationPoolTokens *big.Int
	TotalStakingPoolTokens   *big.Int
	TreasuryAmount           *big.Int
}, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getPoolInfo")

	outstruct := new(struct {
		TotalOperationPoolTokens *big.Int
		TotalStakingPoolTokens   *big.Int
		TreasuryAmount           *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalOperationPoolTokens = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalStakingPoolTokens = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TreasuryAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPoolInfo is a free data retrieval call binding the contract method 0x60246c88.
//
// Solidity: function getPoolInfo() view returns(uint256 totalOperationPoolTokens, uint256 totalStakingPoolTokens, uint256 treasuryAmount)
func (_Chips *ChipsSession) GetPoolInfo() (struct {
	TotalOperationPoolTokens *big.Int
	TotalStakingPoolTokens   *big.Int
	TreasuryAmount           *big.Int
}, error) {
	return _Chips.Contract.GetPoolInfo(&_Chips.CallOpts)
}

// GetPoolInfo is a free data retrieval call binding the contract method 0x60246c88.
//
// Solidity: function getPoolInfo() view returns(uint256 totalOperationPoolTokens, uint256 totalStakingPoolTokens, uint256 treasuryAmount)
func (_Chips *ChipsCallerSession) GetPoolInfo() (struct {
	TotalOperationPoolTokens *big.Int
	TotalStakingPoolTokens   *big.Int
	TreasuryAmount           *big.Int
}, error) {
	return _Chips.Contract.GetPoolInfo(&_Chips.CallOpts)
}

// GetPublicPool is a free data retrieval call binding the contract method 0xc84c42a3.
//
// Solidity: function getPublicPool() view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256))
func (_Chips *ChipsCaller) GetPublicPool(opts *bind.CallOpts) (DataTypesNode, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getPublicPool")

	if err != nil {
		return *new(DataTypesNode), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesNode)).(*DataTypesNode)

	return out0, err

}

// GetPublicPool is a free data retrieval call binding the contract method 0xc84c42a3.
//
// Solidity: function getPublicPool() view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256))
func (_Chips *ChipsSession) GetPublicPool() (DataTypesNode, error) {
	return _Chips.Contract.GetPublicPool(&_Chips.CallOpts)
}

// GetPublicPool is a free data retrieval call binding the contract method 0xc84c42a3.
//
// Solidity: function getPublicPool() view returns((address,uint64,bool,string,string,uint256,uint256,uint256,uint256))
func (_Chips *ChipsCallerSession) GetPublicPool() (DataTypesNode, error) {
	return _Chips.Contract.GetPublicPool(&_Chips.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Chips *ChipsCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Chips *ChipsSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Chips.Contract.GetRoleAdmin(&_Chips.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Chips *ChipsCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Chips.Contract.GetRoleAdmin(&_Chips.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Chips *ChipsCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Chips *ChipsSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Chips.Contract.GetRoleMember(&_Chips.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_Chips *ChipsCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _Chips.Contract.GetRoleMember(&_Chips.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Chips *ChipsCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Chips *ChipsSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Chips.Contract.GetRoleMemberCount(&_Chips.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_Chips *ChipsCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _Chips.Contract.GetRoleMemberCount(&_Chips.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Chips *ChipsCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Chips *ChipsSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Chips.Contract.HasRole(&_Chips.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Chips *ChipsCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Chips.Contract.HasRole(&_Chips.CallOpts, role, account)
}

// MinTokensToStake is a free data retrieval call binding the contract method 0x14936b13.
//
// Solidity: function minTokensToStake(address nodeAddr) view returns(uint256)
func (_Chips *ChipsCaller) MinTokensToStake(opts *bind.CallOpts, nodeAddr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "minTokensToStake", nodeAddr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinTokensToStake is a free data retrieval call binding the contract method 0x14936b13.
//
// Solidity: function minTokensToStake(address nodeAddr) view returns(uint256)
func (_Chips *ChipsSession) MinTokensToStake(nodeAddr common.Address) (*big.Int, error) {
	return _Chips.Contract.MinTokensToStake(&_Chips.CallOpts, nodeAddr)
}

// MinTokensToStake is a free data retrieval call binding the contract method 0x14936b13.
//
// Solidity: function minTokensToStake(address nodeAddr) view returns(uint256)
func (_Chips *ChipsCallerSession) MinTokensToStake(nodeAddr common.Address) (*big.Int, error) {
	return _Chips.Contract.MinTokensToStake(&_Chips.CallOpts, nodeAddr)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Chips *ChipsCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Chips *ChipsSession) Paused() (bool, error) {
	return _Chips.Contract.Paused(&_Chips.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Chips *ChipsCallerSession) Paused() (bool, error) {
	return _Chips.Contract.Paused(&_Chips.CallOpts)
}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_Chips *ChipsCaller) StakingToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "stakingToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_Chips *ChipsSession) StakingToken() (common.Address, error) {
	return _Chips.Contract.StakingToken(&_Chips.CallOpts)
}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_Chips *ChipsCallerSession) StakingToken() (common.Address, error) {
	return _Chips.Contract.StakingToken(&_Chips.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Chips *ChipsCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Chips.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Chips *ChipsSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Chips.Contract.SupportsInterface(&_Chips.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Chips *ChipsCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Chips.Contract.SupportsInterface(&_Chips.CallOpts, interfaceId)
}

// ClaimUnstake is a paid mutator transaction binding the contract method 0x04a4fb10.
//
// Solidity: function claimUnstake(uint256[] requestIds) returns()
func (_Chips *ChipsTransactor) ClaimUnstake(opts *bind.TransactOpts, requestIds []*big.Int) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "claimUnstake", requestIds)
}

// ClaimUnstake is a paid mutator transaction binding the contract method 0x04a4fb10.
//
// Solidity: function claimUnstake(uint256[] requestIds) returns()
func (_Chips *ChipsSession) ClaimUnstake(requestIds []*big.Int) (*types.Transaction, error) {
	return _Chips.Contract.ClaimUnstake(&_Chips.TransactOpts, requestIds)
}

// ClaimUnstake is a paid mutator transaction binding the contract method 0x04a4fb10.
//
// Solidity: function claimUnstake(uint256[] requestIds) returns()
func (_Chips *ChipsTransactorSession) ClaimUnstake(requestIds []*big.Int) (*types.Transaction, error) {
	return _Chips.Contract.ClaimUnstake(&_Chips.TransactOpts, requestIds)
}

// ClaimWithdrawal is a paid mutator transaction binding the contract method 0x3c256b98.
//
// Solidity: function claimWithdrawal(uint256[] requestIds) returns()
func (_Chips *ChipsTransactor) ClaimWithdrawal(opts *bind.TransactOpts, requestIds []*big.Int) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "claimWithdrawal", requestIds)
}

// ClaimWithdrawal is a paid mutator transaction binding the contract method 0x3c256b98.
//
// Solidity: function claimWithdrawal(uint256[] requestIds) returns()
func (_Chips *ChipsSession) ClaimWithdrawal(requestIds []*big.Int) (*types.Transaction, error) {
	return _Chips.Contract.ClaimWithdrawal(&_Chips.TransactOpts, requestIds)
}

// ClaimWithdrawal is a paid mutator transaction binding the contract method 0x3c256b98.
//
// Solidity: function claimWithdrawal(uint256[] requestIds) returns()
func (_Chips *ChipsTransactorSession) ClaimWithdrawal(requestIds []*big.Int) (*types.Transaction, error) {
	return _Chips.Contract.ClaimWithdrawal(&_Chips.TransactOpts, requestIds)
}

// CreateNode is a paid mutator transaction binding the contract method 0xec3f4783.
//
// Solidity: function createNode(address to, string name, string description, uint64 taxRateBasisPoints, bool publicGood) returns()
func (_Chips *ChipsTransactor) CreateNode(opts *bind.TransactOpts, to common.Address, name string, description string, taxRateBasisPoints uint64, publicGood bool) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "createNode", to, name, description, taxRateBasisPoints, publicGood)
}

// CreateNode is a paid mutator transaction binding the contract method 0xec3f4783.
//
// Solidity: function createNode(address to, string name, string description, uint64 taxRateBasisPoints, bool publicGood) returns()
func (_Chips *ChipsSession) CreateNode(to common.Address, name string, description string, taxRateBasisPoints uint64, publicGood bool) (*types.Transaction, error) {
	return _Chips.Contract.CreateNode(&_Chips.TransactOpts, to, name, description, taxRateBasisPoints, publicGood)
}

// CreateNode is a paid mutator transaction binding the contract method 0xec3f4783.
//
// Solidity: function createNode(address to, string name, string description, uint64 taxRateBasisPoints, bool publicGood) returns()
func (_Chips *ChipsTransactorSession) CreateNode(to common.Address, name string, description string, taxRateBasisPoints uint64, publicGood bool) (*types.Transaction, error) {
	return _Chips.Contract.CreateNode(&_Chips.TransactOpts, to, name, description, taxRateBasisPoints, publicGood)
}

// CreateNodeAndDeposit is a paid mutator transaction binding the contract method 0x52ada782.
//
// Solidity: function createNodeAndDeposit(string name, string description, uint64 taxRateBasisPoints, bool publicGood, uint256 amount) returns()
func (_Chips *ChipsTransactor) CreateNodeAndDeposit(opts *bind.TransactOpts, name string, description string, taxRateBasisPoints uint64, publicGood bool, amount *big.Int) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "createNodeAndDeposit", name, description, taxRateBasisPoints, publicGood, amount)
}

// CreateNodeAndDeposit is a paid mutator transaction binding the contract method 0x52ada782.
//
// Solidity: function createNodeAndDeposit(string name, string description, uint64 taxRateBasisPoints, bool publicGood, uint256 amount) returns()
func (_Chips *ChipsSession) CreateNodeAndDeposit(name string, description string, taxRateBasisPoints uint64, publicGood bool, amount *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.CreateNodeAndDeposit(&_Chips.TransactOpts, name, description, taxRateBasisPoints, publicGood, amount)
}

// CreateNodeAndDeposit is a paid mutator transaction binding the contract method 0x52ada782.
//
// Solidity: function createNodeAndDeposit(string name, string description, uint64 taxRateBasisPoints, bool publicGood, uint256 amount) returns()
func (_Chips *ChipsTransactorSession) CreateNodeAndDeposit(name string, description string, taxRateBasisPoints uint64, publicGood bool, amount *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.CreateNodeAndDeposit(&_Chips.TransactOpts, name, description, taxRateBasisPoints, publicGood, amount)
}

// DeleteNode is a paid mutator transaction binding the contract method 0x2d4ede93.
//
// Solidity: function deleteNode(address nodeAddr) returns()
func (_Chips *ChipsTransactor) DeleteNode(opts *bind.TransactOpts, nodeAddr common.Address) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "deleteNode", nodeAddr)
}

// DeleteNode is a paid mutator transaction binding the contract method 0x2d4ede93.
//
// Solidity: function deleteNode(address nodeAddr) returns()
func (_Chips *ChipsSession) DeleteNode(nodeAddr common.Address) (*types.Transaction, error) {
	return _Chips.Contract.DeleteNode(&_Chips.TransactOpts, nodeAddr)
}

// DeleteNode is a paid mutator transaction binding the contract method 0x2d4ede93.
//
// Solidity: function deleteNode(address nodeAddr) returns()
func (_Chips *ChipsTransactorSession) DeleteNode(nodeAddr common.Address) (*types.Transaction, error) {
	return _Chips.Contract.DeleteNode(&_Chips.TransactOpts, nodeAddr)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Chips *ChipsTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Chips *ChipsSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.Deposit(&_Chips.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Chips *ChipsTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.Deposit(&_Chips.TransactOpts, amount)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x8e3e6174.
//
// Solidity: function distributeRewards(uint256[3] epochInfo, address[] nodeAddrs, uint256[] requestFees, uint256[] operationRewards, uint256[] stakingRewards, uint256 publicPoolReward) returns()
func (_Chips *ChipsTransactor) DistributeRewards(opts *bind.TransactOpts, epochInfo [3]*big.Int, nodeAddrs []common.Address, requestFees []*big.Int, operationRewards []*big.Int, stakingRewards []*big.Int, publicPoolReward *big.Int) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "distributeRewards", epochInfo, nodeAddrs, requestFees, operationRewards, stakingRewards, publicPoolReward)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x8e3e6174.
//
// Solidity: function distributeRewards(uint256[3] epochInfo, address[] nodeAddrs, uint256[] requestFees, uint256[] operationRewards, uint256[] stakingRewards, uint256 publicPoolReward) returns()
func (_Chips *ChipsSession) DistributeRewards(epochInfo [3]*big.Int, nodeAddrs []common.Address, requestFees []*big.Int, operationRewards []*big.Int, stakingRewards []*big.Int, publicPoolReward *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.DistributeRewards(&_Chips.TransactOpts, epochInfo, nodeAddrs, requestFees, operationRewards, stakingRewards, publicPoolReward)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x8e3e6174.
//
// Solidity: function distributeRewards(uint256[3] epochInfo, address[] nodeAddrs, uint256[] requestFees, uint256[] operationRewards, uint256[] stakingRewards, uint256 publicPoolReward) returns()
func (_Chips *ChipsTransactorSession) DistributeRewards(epochInfo [3]*big.Int, nodeAddrs []common.Address, requestFees []*big.Int, operationRewards []*big.Int, stakingRewards []*big.Int, publicPoolReward *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.DistributeRewards(&_Chips.TransactOpts, epochInfo, nodeAddrs, requestFees, operationRewards, stakingRewards, publicPoolReward)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Chips *ChipsTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Chips *ChipsSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Chips.Contract.GrantRole(&_Chips.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Chips *ChipsTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Chips.Contract.GrantRole(&_Chips.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address chips, address pauseAccount, address oracleAccount) returns()
func (_Chips *ChipsTransactor) Initialize(opts *bind.TransactOpts, chips common.Address, pauseAccount common.Address, oracleAccount common.Address) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "initialize", chips, pauseAccount, oracleAccount)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address chips, address pauseAccount, address oracleAccount) returns()
func (_Chips *ChipsSession) Initialize(chips common.Address, pauseAccount common.Address, oracleAccount common.Address) (*types.Transaction, error) {
	return _Chips.Contract.Initialize(&_Chips.TransactOpts, chips, pauseAccount, oracleAccount)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address chips, address pauseAccount, address oracleAccount) returns()
func (_Chips *ChipsTransactorSession) Initialize(chips common.Address, pauseAccount common.Address, oracleAccount common.Address) (*types.Transaction, error) {
	return _Chips.Contract.Initialize(&_Chips.TransactOpts, chips, pauseAccount, oracleAccount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Chips *ChipsTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Chips *ChipsSession) Pause() (*types.Transaction, error) {
	return _Chips.Contract.Pause(&_Chips.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Chips *ChipsTransactorSession) Pause() (*types.Transaction, error) {
	return _Chips.Contract.Pause(&_Chips.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Chips *ChipsTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Chips *ChipsSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Chips.Contract.RenounceRole(&_Chips.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Chips *ChipsTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Chips.Contract.RenounceRole(&_Chips.TransactOpts, role, callerConfirmation)
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xbcdd4190.
//
// Solidity: function requestUnstake(address nodeAddr, uint256[] chipsIds) returns(uint256 requestId)
func (_Chips *ChipsTransactor) RequestUnstake(opts *bind.TransactOpts, nodeAddr common.Address, chipsIds []*big.Int) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "requestUnstake", nodeAddr, chipsIds)
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xbcdd4190.
//
// Solidity: function requestUnstake(address nodeAddr, uint256[] chipsIds) returns(uint256 requestId)
func (_Chips *ChipsSession) RequestUnstake(nodeAddr common.Address, chipsIds []*big.Int) (*types.Transaction, error) {
	return _Chips.Contract.RequestUnstake(&_Chips.TransactOpts, nodeAddr, chipsIds)
}

// RequestUnstake is a paid mutator transaction binding the contract method 0xbcdd4190.
//
// Solidity: function requestUnstake(address nodeAddr, uint256[] chipsIds) returns(uint256 requestId)
func (_Chips *ChipsTransactorSession) RequestUnstake(nodeAddr common.Address, chipsIds []*big.Int) (*types.Transaction, error) {
	return _Chips.Contract.RequestUnstake(&_Chips.TransactOpts, nodeAddr, chipsIds)
}

// RequestUnstakeFromPublicPool is a paid mutator transaction binding the contract method 0x22ce5370.
//
// Solidity: function requestUnstakeFromPublicPool(uint256[] chipsIds) returns(uint256 requestId)
func (_Chips *ChipsTransactor) RequestUnstakeFromPublicPool(opts *bind.TransactOpts, chipsIds []*big.Int) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "requestUnstakeFromPublicPool", chipsIds)
}

// RequestUnstakeFromPublicPool is a paid mutator transaction binding the contract method 0x22ce5370.
//
// Solidity: function requestUnstakeFromPublicPool(uint256[] chipsIds) returns(uint256 requestId)
func (_Chips *ChipsSession) RequestUnstakeFromPublicPool(chipsIds []*big.Int) (*types.Transaction, error) {
	return _Chips.Contract.RequestUnstakeFromPublicPool(&_Chips.TransactOpts, chipsIds)
}

// RequestUnstakeFromPublicPool is a paid mutator transaction binding the contract method 0x22ce5370.
//
// Solidity: function requestUnstakeFromPublicPool(uint256[] chipsIds) returns(uint256 requestId)
func (_Chips *ChipsTransactorSession) RequestUnstakeFromPublicPool(chipsIds []*big.Int) (*types.Transaction, error) {
	return _Chips.Contract.RequestUnstakeFromPublicPool(&_Chips.TransactOpts, chipsIds)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0x9ee679e8.
//
// Solidity: function requestWithdrawal(uint256 amount) returns(uint256 requestId)
func (_Chips *ChipsTransactor) RequestWithdrawal(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "requestWithdrawal", amount)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0x9ee679e8.
//
// Solidity: function requestWithdrawal(uint256 amount) returns(uint256 requestId)
func (_Chips *ChipsSession) RequestWithdrawal(amount *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.RequestWithdrawal(&_Chips.TransactOpts, amount)
}

// RequestWithdrawal is a paid mutator transaction binding the contract method 0x9ee679e8.
//
// Solidity: function requestWithdrawal(uint256 amount) returns(uint256 requestId)
func (_Chips *ChipsTransactorSession) RequestWithdrawal(amount *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.RequestWithdrawal(&_Chips.TransactOpts, amount)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Chips *ChipsTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Chips *ChipsSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Chips.Contract.RevokeRole(&_Chips.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Chips *ChipsTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Chips.Contract.RevokeRole(&_Chips.TransactOpts, role, account)
}

// SetTaxRateBasisPoints4Node is a paid mutator transaction binding the contract method 0xb65d660c.
//
// Solidity: function setTaxRateBasisPoints4Node(address nodeAddr, uint64 taxRateBasisPoints) returns()
func (_Chips *ChipsTransactor) SetTaxRateBasisPoints4Node(opts *bind.TransactOpts, nodeAddr common.Address, taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "setTaxRateBasisPoints4Node", nodeAddr, taxRateBasisPoints)
}

// SetTaxRateBasisPoints4Node is a paid mutator transaction binding the contract method 0xb65d660c.
//
// Solidity: function setTaxRateBasisPoints4Node(address nodeAddr, uint64 taxRateBasisPoints) returns()
func (_Chips *ChipsSession) SetTaxRateBasisPoints4Node(nodeAddr common.Address, taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Chips.Contract.SetTaxRateBasisPoints4Node(&_Chips.TransactOpts, nodeAddr, taxRateBasisPoints)
}

// SetTaxRateBasisPoints4Node is a paid mutator transaction binding the contract method 0xb65d660c.
//
// Solidity: function setTaxRateBasisPoints4Node(address nodeAddr, uint64 taxRateBasisPoints) returns()
func (_Chips *ChipsTransactorSession) SetTaxRateBasisPoints4Node(nodeAddr common.Address, taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Chips.Contract.SetTaxRateBasisPoints4Node(&_Chips.TransactOpts, nodeAddr, taxRateBasisPoints)
}

// SetTaxRateBasisPoints4PublicPool is a paid mutator transaction binding the contract method 0xe3fb8dca.
//
// Solidity: function setTaxRateBasisPoints4PublicPool(uint64 taxRateBasisPoints) returns()
func (_Chips *ChipsTransactor) SetTaxRateBasisPoints4PublicPool(opts *bind.TransactOpts, taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "setTaxRateBasisPoints4PublicPool", taxRateBasisPoints)
}

// SetTaxRateBasisPoints4PublicPool is a paid mutator transaction binding the contract method 0xe3fb8dca.
//
// Solidity: function setTaxRateBasisPoints4PublicPool(uint64 taxRateBasisPoints) returns()
func (_Chips *ChipsSession) SetTaxRateBasisPoints4PublicPool(taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Chips.Contract.SetTaxRateBasisPoints4PublicPool(&_Chips.TransactOpts, taxRateBasisPoints)
}

// SetTaxRateBasisPoints4PublicPool is a paid mutator transaction binding the contract method 0xe3fb8dca.
//
// Solidity: function setTaxRateBasisPoints4PublicPool(uint64 taxRateBasisPoints) returns()
func (_Chips *ChipsTransactorSession) SetTaxRateBasisPoints4PublicPool(taxRateBasisPoints uint64) (*types.Transaction, error) {
	return _Chips.Contract.SetTaxRateBasisPoints4PublicPool(&_Chips.TransactOpts, taxRateBasisPoints)
}

// SlashNodes is a paid mutator transaction binding the contract method 0xa2f641c3.
//
// Solidity: function slashNodes(address[] nodeAddrs) returns()
func (_Chips *ChipsTransactor) SlashNodes(opts *bind.TransactOpts, nodeAddrs []common.Address) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "slashNodes", nodeAddrs)
}

// SlashNodes is a paid mutator transaction binding the contract method 0xa2f641c3.
//
// Solidity: function slashNodes(address[] nodeAddrs) returns()
func (_Chips *ChipsSession) SlashNodes(nodeAddrs []common.Address) (*types.Transaction, error) {
	return _Chips.Contract.SlashNodes(&_Chips.TransactOpts, nodeAddrs)
}

// SlashNodes is a paid mutator transaction binding the contract method 0xa2f641c3.
//
// Solidity: function slashNodes(address[] nodeAddrs) returns()
func (_Chips *ChipsTransactorSession) SlashNodes(nodeAddrs []common.Address) (*types.Transaction, error) {
	return _Chips.Contract.SlashNodes(&_Chips.TransactOpts, nodeAddrs)
}

// Stake is a paid mutator transaction binding the contract method 0xadc9772e.
//
// Solidity: function stake(address nodeAddr, uint256 amount) returns(uint256 startTokenId, uint256 endTokenId)
func (_Chips *ChipsTransactor) Stake(opts *bind.TransactOpts, nodeAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "stake", nodeAddr, amount)
}

// Stake is a paid mutator transaction binding the contract method 0xadc9772e.
//
// Solidity: function stake(address nodeAddr, uint256 amount) returns(uint256 startTokenId, uint256 endTokenId)
func (_Chips *ChipsSession) Stake(nodeAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.Stake(&_Chips.TransactOpts, nodeAddr, amount)
}

// Stake is a paid mutator transaction binding the contract method 0xadc9772e.
//
// Solidity: function stake(address nodeAddr, uint256 amount) returns(uint256 startTokenId, uint256 endTokenId)
func (_Chips *ChipsTransactorSession) Stake(nodeAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.Stake(&_Chips.TransactOpts, nodeAddr, amount)
}

// StakeToPublicPool is a paid mutator transaction binding the contract method 0x53db41ac.
//
// Solidity: function stakeToPublicPool(address nodeAddr, uint256 amount) returns(uint256 startTokenId, uint256 endTokenId)
func (_Chips *ChipsTransactor) StakeToPublicPool(opts *bind.TransactOpts, nodeAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "stakeToPublicPool", nodeAddr, amount)
}

// StakeToPublicPool is a paid mutator transaction binding the contract method 0x53db41ac.
//
// Solidity: function stakeToPublicPool(address nodeAddr, uint256 amount) returns(uint256 startTokenId, uint256 endTokenId)
func (_Chips *ChipsSession) StakeToPublicPool(nodeAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.StakeToPublicPool(&_Chips.TransactOpts, nodeAddr, amount)
}

// StakeToPublicPool is a paid mutator transaction binding the contract method 0x53db41ac.
//
// Solidity: function stakeToPublicPool(address nodeAddr, uint256 amount) returns(uint256 startTokenId, uint256 endTokenId)
func (_Chips *ChipsTransactorSession) StakeToPublicPool(nodeAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Chips.Contract.StakeToPublicPool(&_Chips.TransactOpts, nodeAddr, amount)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Chips *ChipsTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Chips *ChipsSession) Unpause() (*types.Transaction, error) {
	return _Chips.Contract.Unpause(&_Chips.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Chips *ChipsTransactorSession) Unpause() (*types.Transaction, error) {
	return _Chips.Contract.Unpause(&_Chips.TransactOpts)
}

// Withdraw2Treasury is a paid mutator transaction binding the contract method 0x4a7dfc90.
//
// Solidity: function withdraw2Treasury() returns()
func (_Chips *ChipsTransactor) Withdraw2Treasury(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chips.contract.Transact(opts, "withdraw2Treasury")
}

// Withdraw2Treasury is a paid mutator transaction binding the contract method 0x4a7dfc90.
//
// Solidity: function withdraw2Treasury() returns()
func (_Chips *ChipsSession) Withdraw2Treasury() (*types.Transaction, error) {
	return _Chips.Contract.Withdraw2Treasury(&_Chips.TransactOpts)
}

// Withdraw2Treasury is a paid mutator transaction binding the contract method 0x4a7dfc90.
//
// Solidity: function withdraw2Treasury() returns()
func (_Chips *ChipsTransactorSession) Withdraw2Treasury() (*types.Transaction, error) {
	return _Chips.Contract.Withdraw2Treasury(&_Chips.TransactOpts)
}

// ChipsDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the Chips contract.
type ChipsDepositedIterator struct {
	Event *ChipsDeposited // Event containing the contract specifics and raw log

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
func (it *ChipsDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsDeposited)
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
		it.Event = new(ChipsDeposited)
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
func (it *ChipsDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsDeposited represents a Deposited event raised by the Chips contract.
type ChipsDeposited struct {
	NodeAddr common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed nodeAddr, uint256 indexed amount)
func (_Chips *ChipsFilterer) FilterDeposited(opts *bind.FilterOpts, nodeAddr []common.Address, amount []*big.Int) (*ChipsDepositedIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Chips.contract.FilterLogs(opts, "Deposited", nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &ChipsDepositedIterator{contract: _Chips.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed nodeAddr, uint256 indexed amount)
func (_Chips *ChipsFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *ChipsDeposited, nodeAddr []common.Address, amount []*big.Int) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Chips.contract.WatchLogs(opts, "Deposited", nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsDeposited)
				if err := _Chips.contract.UnpackLog(event, "Deposited", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseDeposited(log types.Log) (*ChipsDeposited, error) {
	event := new(ChipsDeposited)
	if err := _Chips.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Chips contract.
type ChipsInitializedIterator struct {
	Event *ChipsInitialized // Event containing the contract specifics and raw log

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
func (it *ChipsInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsInitialized)
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
		it.Event = new(ChipsInitialized)
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
func (it *ChipsInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsInitialized represents a Initialized event raised by the Chips contract.
type ChipsInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Chips *ChipsFilterer) FilterInitialized(opts *bind.FilterOpts) (*ChipsInitializedIterator, error) {

	logs, sub, err := _Chips.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ChipsInitializedIterator{contract: _Chips.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Chips *ChipsFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ChipsInitialized) (event.Subscription, error) {

	logs, sub, err := _Chips.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsInitialized)
				if err := _Chips.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseInitialized(log types.Log) (*ChipsInitialized, error) {
	event := new(ChipsInitialized)
	if err := _Chips.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsNodeCreatedIterator is returned from FilterNodeCreated and is used to iterate over the raw logs and unpacked data for NodeCreated events raised by the Chips contract.
type ChipsNodeCreatedIterator struct {
	Event *ChipsNodeCreated // Event containing the contract specifics and raw log

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
func (it *ChipsNodeCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsNodeCreated)
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
		it.Event = new(ChipsNodeCreated)
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
func (it *ChipsNodeCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsNodeCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsNodeCreated represents a NodeCreated event raised by the Chips contract.
type ChipsNodeCreated struct {
	NodeAddr           common.Address
	Name               string
	Description        string
	TaxRateBasisPoints uint64
	PublicGood         bool
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterNodeCreated is a free log retrieval operation binding the contract event 0x6ae2420fc4e19a67c13bace8d75372b34c462f70f67dc8c54cb587d4f493c044.
//
// Solidity: event NodeCreated(address indexed nodeAddr, string name, string description, uint64 taxRateBasisPoints, bool publicGood)
func (_Chips *ChipsFilterer) FilterNodeCreated(opts *bind.FilterOpts, nodeAddr []common.Address) (*ChipsNodeCreatedIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Chips.contract.FilterLogs(opts, "NodeCreated", nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return &ChipsNodeCreatedIterator{contract: _Chips.contract, event: "NodeCreated", logs: logs, sub: sub}, nil
}

// WatchNodeCreated is a free log subscription operation binding the contract event 0x6ae2420fc4e19a67c13bace8d75372b34c462f70f67dc8c54cb587d4f493c044.
//
// Solidity: event NodeCreated(address indexed nodeAddr, string name, string description, uint64 taxRateBasisPoints, bool publicGood)
func (_Chips *ChipsFilterer) WatchNodeCreated(opts *bind.WatchOpts, sink chan<- *ChipsNodeCreated, nodeAddr []common.Address) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Chips.contract.WatchLogs(opts, "NodeCreated", nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsNodeCreated)
				if err := _Chips.contract.UnpackLog(event, "NodeCreated", log); err != nil {
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

// ParseNodeCreated is a log parse operation binding the contract event 0x6ae2420fc4e19a67c13bace8d75372b34c462f70f67dc8c54cb587d4f493c044.
//
// Solidity: event NodeCreated(address indexed nodeAddr, string name, string description, uint64 taxRateBasisPoints, bool publicGood)
func (_Chips *ChipsFilterer) ParseNodeCreated(log types.Log) (*ChipsNodeCreated, error) {
	event := new(ChipsNodeCreated)
	if err := _Chips.contract.UnpackLog(event, "NodeCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsNodeDeletedIterator is returned from FilterNodeDeleted and is used to iterate over the raw logs and unpacked data for NodeDeleted events raised by the Chips contract.
type ChipsNodeDeletedIterator struct {
	Event *ChipsNodeDeleted // Event containing the contract specifics and raw log

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
func (it *ChipsNodeDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsNodeDeleted)
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
		it.Event = new(ChipsNodeDeleted)
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
func (it *ChipsNodeDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsNodeDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsNodeDeleted represents a NodeDeleted event raised by the Chips contract.
type ChipsNodeDeleted struct {
	NodeAddr common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNodeDeleted is a free log retrieval operation binding the contract event 0x1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad.
//
// Solidity: event NodeDeleted(address indexed nodeAddr)
func (_Chips *ChipsFilterer) FilterNodeDeleted(opts *bind.FilterOpts, nodeAddr []common.Address) (*ChipsNodeDeletedIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Chips.contract.FilterLogs(opts, "NodeDeleted", nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return &ChipsNodeDeletedIterator{contract: _Chips.contract, event: "NodeDeleted", logs: logs, sub: sub}, nil
}

// WatchNodeDeleted is a free log subscription operation binding the contract event 0x1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad.
//
// Solidity: event NodeDeleted(address indexed nodeAddr)
func (_Chips *ChipsFilterer) WatchNodeDeleted(opts *bind.WatchOpts, sink chan<- *ChipsNodeDeleted, nodeAddr []common.Address) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Chips.contract.WatchLogs(opts, "NodeDeleted", nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsNodeDeleted)
				if err := _Chips.contract.UnpackLog(event, "NodeDeleted", log); err != nil {
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

// ParseNodeDeleted is a log parse operation binding the contract event 0x1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad.
//
// Solidity: event NodeDeleted(address indexed nodeAddr)
func (_Chips *ChipsFilterer) ParseNodeDeleted(log types.Log) (*ChipsNodeDeleted, error) {
	event := new(ChipsNodeDeleted)
	if err := _Chips.contract.UnpackLog(event, "NodeDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsNodeSlashedIterator is returned from FilterNodeSlashed and is used to iterate over the raw logs and unpacked data for NodeSlashed events raised by the Chips contract.
type ChipsNodeSlashedIterator struct {
	Event *ChipsNodeSlashed // Event containing the contract specifics and raw log

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
func (it *ChipsNodeSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsNodeSlashed)
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
		it.Event = new(ChipsNodeSlashed)
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
func (it *ChipsNodeSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsNodeSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsNodeSlashed represents a NodeSlashed event raised by the Chips contract.
type ChipsNodeSlashed struct {
	NodeAddr             common.Address
	SlashedOperationPool *big.Int
	SlashedStakingPool   *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterNodeSlashed is a free log retrieval operation binding the contract event 0xa8d720d0a0a2e7c96bf9eb87433901ebb6331356c8f3283b2568de34478703cc.
//
// Solidity: event NodeSlashed(address indexed nodeAddr, uint256 indexed slashedOperationPool, uint256 indexed slashedStakingPool)
func (_Chips *ChipsFilterer) FilterNodeSlashed(opts *bind.FilterOpts, nodeAddr []common.Address, slashedOperationPool []*big.Int, slashedStakingPool []*big.Int) (*ChipsNodeSlashedIterator, error) {

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

	logs, sub, err := _Chips.contract.FilterLogs(opts, "NodeSlashed", nodeAddrRule, slashedOperationPoolRule, slashedStakingPoolRule)
	if err != nil {
		return nil, err
	}
	return &ChipsNodeSlashedIterator{contract: _Chips.contract, event: "NodeSlashed", logs: logs, sub: sub}, nil
}

// WatchNodeSlashed is a free log subscription operation binding the contract event 0xa8d720d0a0a2e7c96bf9eb87433901ebb6331356c8f3283b2568de34478703cc.
//
// Solidity: event NodeSlashed(address indexed nodeAddr, uint256 indexed slashedOperationPool, uint256 indexed slashedStakingPool)
func (_Chips *ChipsFilterer) WatchNodeSlashed(opts *bind.WatchOpts, sink chan<- *ChipsNodeSlashed, nodeAddr []common.Address, slashedOperationPool []*big.Int, slashedStakingPool []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Chips.contract.WatchLogs(opts, "NodeSlashed", nodeAddrRule, slashedOperationPoolRule, slashedStakingPoolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsNodeSlashed)
				if err := _Chips.contract.UnpackLog(event, "NodeSlashed", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseNodeSlashed(log types.Log) (*ChipsNodeSlashed, error) {
	event := new(ChipsNodeSlashed)
	if err := _Chips.contract.UnpackLog(event, "NodeSlashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsNodeTaxRateBasisPointsSetIterator is returned from FilterNodeTaxRateBasisPointsSet and is used to iterate over the raw logs and unpacked data for NodeTaxRateBasisPointsSet events raised by the Chips contract.
type ChipsNodeTaxRateBasisPointsSetIterator struct {
	Event *ChipsNodeTaxRateBasisPointsSet // Event containing the contract specifics and raw log

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
func (it *ChipsNodeTaxRateBasisPointsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsNodeTaxRateBasisPointsSet)
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
		it.Event = new(ChipsNodeTaxRateBasisPointsSet)
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
func (it *ChipsNodeTaxRateBasisPointsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsNodeTaxRateBasisPointsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsNodeTaxRateBasisPointsSet represents a NodeTaxRateBasisPointsSet event raised by the Chips contract.
type ChipsNodeTaxRateBasisPointsSet struct {
	NodeAddr           common.Address
	TaxRateBasisPoints uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterNodeTaxRateBasisPointsSet is a free log retrieval operation binding the contract event 0xb8e5551053b871a40f7c7382e5bd3af5a62dd737d059d3838cf3aa7c325bd479.
//
// Solidity: event NodeTaxRateBasisPointsSet(address indexed nodeAddr, uint64 indexed taxRateBasisPoints)
func (_Chips *ChipsFilterer) FilterNodeTaxRateBasisPointsSet(opts *bind.FilterOpts, nodeAddr []common.Address, taxRateBasisPoints []uint64) (*ChipsNodeTaxRateBasisPointsSetIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Chips.contract.FilterLogs(opts, "NodeTaxRateBasisPointsSet", nodeAddrRule, taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return &ChipsNodeTaxRateBasisPointsSetIterator{contract: _Chips.contract, event: "NodeTaxRateBasisPointsSet", logs: logs, sub: sub}, nil
}

// WatchNodeTaxRateBasisPointsSet is a free log subscription operation binding the contract event 0xb8e5551053b871a40f7c7382e5bd3af5a62dd737d059d3838cf3aa7c325bd479.
//
// Solidity: event NodeTaxRateBasisPointsSet(address indexed nodeAddr, uint64 indexed taxRateBasisPoints)
func (_Chips *ChipsFilterer) WatchNodeTaxRateBasisPointsSet(opts *bind.WatchOpts, sink chan<- *ChipsNodeTaxRateBasisPointsSet, nodeAddr []common.Address, taxRateBasisPoints []uint64) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Chips.contract.WatchLogs(opts, "NodeTaxRateBasisPointsSet", nodeAddrRule, taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsNodeTaxRateBasisPointsSet)
				if err := _Chips.contract.UnpackLog(event, "NodeTaxRateBasisPointsSet", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseNodeTaxRateBasisPointsSet(log types.Log) (*ChipsNodeTaxRateBasisPointsSet, error) {
	event := new(ChipsNodeTaxRateBasisPointsSet)
	if err := _Chips.contract.UnpackLog(event, "NodeTaxRateBasisPointsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Chips contract.
type ChipsPausedIterator struct {
	Event *ChipsPaused // Event containing the contract specifics and raw log

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
func (it *ChipsPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsPaused)
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
		it.Event = new(ChipsPaused)
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
func (it *ChipsPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsPaused represents a Paused event raised by the Chips contract.
type ChipsPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Chips *ChipsFilterer) FilterPaused(opts *bind.FilterOpts) (*ChipsPausedIterator, error) {

	logs, sub, err := _Chips.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &ChipsPausedIterator{contract: _Chips.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Chips *ChipsFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ChipsPaused) (event.Subscription, error) {

	logs, sub, err := _Chips.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsPaused)
				if err := _Chips.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Chips *ChipsFilterer) ParsePaused(log types.Log) (*ChipsPaused, error) {
	event := new(ChipsPaused)
	if err := _Chips.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsPublicGoodRewardDistributedIterator is returned from FilterPublicGoodRewardDistributed and is used to iterate over the raw logs and unpacked data for PublicGoodRewardDistributed events raised by the Chips contract.
type ChipsPublicGoodRewardDistributedIterator struct {
	Event *ChipsPublicGoodRewardDistributed // Event containing the contract specifics and raw log

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
func (it *ChipsPublicGoodRewardDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsPublicGoodRewardDistributed)
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
		it.Event = new(ChipsPublicGoodRewardDistributed)
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
func (it *ChipsPublicGoodRewardDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsPublicGoodRewardDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsPublicGoodRewardDistributed represents a PublicGoodRewardDistributed event raised by the Chips contract.
type ChipsPublicGoodRewardDistributed struct {
	Epoch            *big.Int
	StartTimestamp   *big.Int
	EndTimestamp     *big.Int
	PublicPoolReward *big.Int
	PublicPoolTax    *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPublicGoodRewardDistributed is a free log retrieval operation binding the contract event 0xab7d25a2f6206ef56c88807f2474ddcd97e1a6323cb25149cde3a607fed6f2d7.
//
// Solidity: event PublicGoodRewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, uint256 publicPoolReward, uint256 publicPoolTax)
func (_Chips *ChipsFilterer) FilterPublicGoodRewardDistributed(opts *bind.FilterOpts, epoch []*big.Int) (*ChipsPublicGoodRewardDistributedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Chips.contract.FilterLogs(opts, "PublicGoodRewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return &ChipsPublicGoodRewardDistributedIterator{contract: _Chips.contract, event: "PublicGoodRewardDistributed", logs: logs, sub: sub}, nil
}

// WatchPublicGoodRewardDistributed is a free log subscription operation binding the contract event 0xab7d25a2f6206ef56c88807f2474ddcd97e1a6323cb25149cde3a607fed6f2d7.
//
// Solidity: event PublicGoodRewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, uint256 publicPoolReward, uint256 publicPoolTax)
func (_Chips *ChipsFilterer) WatchPublicGoodRewardDistributed(opts *bind.WatchOpts, sink chan<- *ChipsPublicGoodRewardDistributed, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Chips.contract.WatchLogs(opts, "PublicGoodRewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsPublicGoodRewardDistributed)
				if err := _Chips.contract.UnpackLog(event, "PublicGoodRewardDistributed", log); err != nil {
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
// Solidity: event PublicGoodRewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, uint256 publicPoolReward, uint256 publicPoolTax)
func (_Chips *ChipsFilterer) ParsePublicGoodRewardDistributed(log types.Log) (*ChipsPublicGoodRewardDistributed, error) {
	event := new(ChipsPublicGoodRewardDistributed)
	if err := _Chips.contract.UnpackLog(event, "PublicGoodRewardDistributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsPublicPoolTaxRateBasisPointsSetIterator is returned from FilterPublicPoolTaxRateBasisPointsSet and is used to iterate over the raw logs and unpacked data for PublicPoolTaxRateBasisPointsSet events raised by the Chips contract.
type ChipsPublicPoolTaxRateBasisPointsSetIterator struct {
	Event *ChipsPublicPoolTaxRateBasisPointsSet // Event containing the contract specifics and raw log

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
func (it *ChipsPublicPoolTaxRateBasisPointsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsPublicPoolTaxRateBasisPointsSet)
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
		it.Event = new(ChipsPublicPoolTaxRateBasisPointsSet)
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
func (it *ChipsPublicPoolTaxRateBasisPointsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsPublicPoolTaxRateBasisPointsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsPublicPoolTaxRateBasisPointsSet represents a PublicPoolTaxRateBasisPointsSet event raised by the Chips contract.
type ChipsPublicPoolTaxRateBasisPointsSet struct {
	TaxRateBasisPoints uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterPublicPoolTaxRateBasisPointsSet is a free log retrieval operation binding the contract event 0x948cf2302b029d76db2ac06e4ef2625e6c687335de349317468f47942a44e8b0.
//
// Solidity: event PublicPoolTaxRateBasisPointsSet(uint64 indexed taxRateBasisPoints)
func (_Chips *ChipsFilterer) FilterPublicPoolTaxRateBasisPointsSet(opts *bind.FilterOpts, taxRateBasisPoints []uint64) (*ChipsPublicPoolTaxRateBasisPointsSetIterator, error) {

	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Chips.contract.FilterLogs(opts, "PublicPoolTaxRateBasisPointsSet", taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return &ChipsPublicPoolTaxRateBasisPointsSetIterator{contract: _Chips.contract, event: "PublicPoolTaxRateBasisPointsSet", logs: logs, sub: sub}, nil
}

// WatchPublicPoolTaxRateBasisPointsSet is a free log subscription operation binding the contract event 0x948cf2302b029d76db2ac06e4ef2625e6c687335de349317468f47942a44e8b0.
//
// Solidity: event PublicPoolTaxRateBasisPointsSet(uint64 indexed taxRateBasisPoints)
func (_Chips *ChipsFilterer) WatchPublicPoolTaxRateBasisPointsSet(opts *bind.WatchOpts, sink chan<- *ChipsPublicPoolTaxRateBasisPointsSet, taxRateBasisPoints []uint64) (event.Subscription, error) {

	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Chips.contract.WatchLogs(opts, "PublicPoolTaxRateBasisPointsSet", taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsPublicPoolTaxRateBasisPointsSet)
				if err := _Chips.contract.UnpackLog(event, "PublicPoolTaxRateBasisPointsSet", log); err != nil {
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
func (_Chips *ChipsFilterer) ParsePublicPoolTaxRateBasisPointsSet(log types.Log) (*ChipsPublicPoolTaxRateBasisPointsSet, error) {
	event := new(ChipsPublicPoolTaxRateBasisPointsSet)
	if err := _Chips.contract.UnpackLog(event, "PublicPoolTaxRateBasisPointsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsRewardDistributedIterator is returned from FilterRewardDistributed and is used to iterate over the raw logs and unpacked data for RewardDistributed events raised by the Chips contract.
type ChipsRewardDistributedIterator struct {
	Event *ChipsRewardDistributed // Event containing the contract specifics and raw log

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
func (it *ChipsRewardDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsRewardDistributed)
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
		it.Event = new(ChipsRewardDistributed)
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
func (it *ChipsRewardDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsRewardDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsRewardDistributed represents a RewardDistributed event raised by the Chips contract.
type ChipsRewardDistributed struct {
	Epoch            *big.Int
	StartTimestamp   *big.Int
	EndTimestamp     *big.Int
	NodeAddrs        []common.Address
	RequestFees      []*big.Int
	OperationRewards []*big.Int
	StakingRewards   []*big.Int
	TaxAmounts       []*big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRewardDistributed is a free log retrieval operation binding the contract event 0x8ea79f19e90b084c2009d3490a097547d8bbb315a883b9efec0996502c1dd7ae.
//
// Solidity: event RewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, address[] nodeAddrs, uint256[] requestFees, uint256[] operationRewards, uint256[] stakingRewards, uint256[] taxAmounts)
func (_Chips *ChipsFilterer) FilterRewardDistributed(opts *bind.FilterOpts, epoch []*big.Int) (*ChipsRewardDistributedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Chips.contract.FilterLogs(opts, "RewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return &ChipsRewardDistributedIterator{contract: _Chips.contract, event: "RewardDistributed", logs: logs, sub: sub}, nil
}

// WatchRewardDistributed is a free log subscription operation binding the contract event 0x8ea79f19e90b084c2009d3490a097547d8bbb315a883b9efec0996502c1dd7ae.
//
// Solidity: event RewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, address[] nodeAddrs, uint256[] requestFees, uint256[] operationRewards, uint256[] stakingRewards, uint256[] taxAmounts)
func (_Chips *ChipsFilterer) WatchRewardDistributed(opts *bind.WatchOpts, sink chan<- *ChipsRewardDistributed, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Chips.contract.WatchLogs(opts, "RewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsRewardDistributed)
				if err := _Chips.contract.UnpackLog(event, "RewardDistributed", log); err != nil {
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
// Solidity: event RewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, address[] nodeAddrs, uint256[] requestFees, uint256[] operationRewards, uint256[] stakingRewards, uint256[] taxAmounts)
func (_Chips *ChipsFilterer) ParseRewardDistributed(log types.Log) (*ChipsRewardDistributed, error) {
	event := new(ChipsRewardDistributed)
	if err := _Chips.contract.UnpackLog(event, "RewardDistributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Chips contract.
type ChipsRoleAdminChangedIterator struct {
	Event *ChipsRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ChipsRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsRoleAdminChanged)
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
		it.Event = new(ChipsRoleAdminChanged)
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
func (it *ChipsRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsRoleAdminChanged represents a RoleAdminChanged event raised by the Chips contract.
type ChipsRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Chips *ChipsFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ChipsRoleAdminChangedIterator, error) {

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

	logs, sub, err := _Chips.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ChipsRoleAdminChangedIterator{contract: _Chips.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Chips *ChipsFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ChipsRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Chips.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsRoleAdminChanged)
				if err := _Chips.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseRoleAdminChanged(log types.Log) (*ChipsRoleAdminChanged, error) {
	event := new(ChipsRoleAdminChanged)
	if err := _Chips.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Chips contract.
type ChipsRoleGrantedIterator struct {
	Event *ChipsRoleGranted // Event containing the contract specifics and raw log

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
func (it *ChipsRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsRoleGranted)
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
		it.Event = new(ChipsRoleGranted)
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
func (it *ChipsRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsRoleGranted represents a RoleGranted event raised by the Chips contract.
type ChipsRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Chips *ChipsFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ChipsRoleGrantedIterator, error) {

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

	logs, sub, err := _Chips.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ChipsRoleGrantedIterator{contract: _Chips.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Chips *ChipsFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ChipsRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Chips.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsRoleGranted)
				if err := _Chips.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseRoleGranted(log types.Log) (*ChipsRoleGranted, error) {
	event := new(ChipsRoleGranted)
	if err := _Chips.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Chips contract.
type ChipsRoleRevokedIterator struct {
	Event *ChipsRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ChipsRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsRoleRevoked)
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
		it.Event = new(ChipsRoleRevoked)
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
func (it *ChipsRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsRoleRevoked represents a RoleRevoked event raised by the Chips contract.
type ChipsRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Chips *ChipsFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ChipsRoleRevokedIterator, error) {

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

	logs, sub, err := _Chips.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ChipsRoleRevokedIterator{contract: _Chips.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Chips *ChipsFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ChipsRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Chips.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsRoleRevoked)
				if err := _Chips.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseRoleRevoked(log types.Log) (*ChipsRoleRevoked, error) {
	event := new(ChipsRoleRevoked)
	if err := _Chips.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Chips contract.
type ChipsStakedIterator struct {
	Event *ChipsStaked // Event containing the contract specifics and raw log

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
func (it *ChipsStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsStaked)
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
		it.Event = new(ChipsStaked)
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
func (it *ChipsStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsStaked represents a Staked event raised by the Chips contract.
type ChipsStaked struct {
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
func (_Chips *ChipsFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address, nodeAddr []common.Address, amount []*big.Int) (*ChipsStakedIterator, error) {

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

	logs, sub, err := _Chips.contract.FilterLogs(opts, "Staked", userRule, nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &ChipsStakedIterator{contract: _Chips.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0xad3fa07f4195b47e64892eb944ecbfc253384053c119852bb2bcae484c2fcb69.
//
// Solidity: event Staked(address indexed user, address indexed nodeAddr, uint256 indexed amount, uint256 startTokenId, uint256 endTokenId)
func (_Chips *ChipsFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *ChipsStaked, user []common.Address, nodeAddr []common.Address, amount []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Chips.contract.WatchLogs(opts, "Staked", userRule, nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsStaked)
				if err := _Chips.contract.UnpackLog(event, "Staked", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseStaked(log types.Log) (*ChipsStaked, error) {
	event := new(ChipsStaked)
	if err := _Chips.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Chips contract.
type ChipsUnpausedIterator struct {
	Event *ChipsUnpaused // Event containing the contract specifics and raw log

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
func (it *ChipsUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsUnpaused)
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
		it.Event = new(ChipsUnpaused)
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
func (it *ChipsUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsUnpaused represents a Unpaused event raised by the Chips contract.
type ChipsUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Chips *ChipsFilterer) FilterUnpaused(opts *bind.FilterOpts) (*ChipsUnpausedIterator, error) {

	logs, sub, err := _Chips.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &ChipsUnpausedIterator{contract: _Chips.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Chips *ChipsFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ChipsUnpaused) (event.Subscription, error) {

	logs, sub, err := _Chips.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsUnpaused)
				if err := _Chips.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseUnpaused(log types.Log) (*ChipsUnpaused, error) {
	event := new(ChipsUnpaused)
	if err := _Chips.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsUnstakeClaimedIterator is returned from FilterUnstakeClaimed and is used to iterate over the raw logs and unpacked data for UnstakeClaimed events raised by the Chips contract.
type ChipsUnstakeClaimedIterator struct {
	Event *ChipsUnstakeClaimed // Event containing the contract specifics and raw log

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
func (it *ChipsUnstakeClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsUnstakeClaimed)
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
		it.Event = new(ChipsUnstakeClaimed)
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
func (it *ChipsUnstakeClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsUnstakeClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsUnstakeClaimed represents a UnstakeClaimed event raised by the Chips contract.
type ChipsUnstakeClaimed struct {
	RequestId     *big.Int
	NodeAddr      common.Address
	User          common.Address
	UnstakeAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUnstakeClaimed is a free log retrieval operation binding the contract event 0x2769ece66eadb650afd8c08c7a8772e39381dddd7230f9e039669e631044d47c.
//
// Solidity: event UnstakeClaimed(uint256 indexed requestId, address indexed nodeAddr, address indexed user, uint256 unstakeAmount)
func (_Chips *ChipsFilterer) FilterUnstakeClaimed(opts *bind.FilterOpts, requestId []*big.Int, nodeAddr []common.Address, user []common.Address) (*ChipsUnstakeClaimedIterator, error) {

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

	logs, sub, err := _Chips.contract.FilterLogs(opts, "UnstakeClaimed", requestIdRule, nodeAddrRule, userRule)
	if err != nil {
		return nil, err
	}
	return &ChipsUnstakeClaimedIterator{contract: _Chips.contract, event: "UnstakeClaimed", logs: logs, sub: sub}, nil
}

// WatchUnstakeClaimed is a free log subscription operation binding the contract event 0x2769ece66eadb650afd8c08c7a8772e39381dddd7230f9e039669e631044d47c.
//
// Solidity: event UnstakeClaimed(uint256 indexed requestId, address indexed nodeAddr, address indexed user, uint256 unstakeAmount)
func (_Chips *ChipsFilterer) WatchUnstakeClaimed(opts *bind.WatchOpts, sink chan<- *ChipsUnstakeClaimed, requestId []*big.Int, nodeAddr []common.Address, user []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Chips.contract.WatchLogs(opts, "UnstakeClaimed", requestIdRule, nodeAddrRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsUnstakeClaimed)
				if err := _Chips.contract.UnpackLog(event, "UnstakeClaimed", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseUnstakeClaimed(log types.Log) (*ChipsUnstakeClaimed, error) {
	event := new(ChipsUnstakeClaimed)
	if err := _Chips.contract.UnpackLog(event, "UnstakeClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsUnstakeRequestedIterator is returned from FilterUnstakeRequested and is used to iterate over the raw logs and unpacked data for UnstakeRequested events raised by the Chips contract.
type ChipsUnstakeRequestedIterator struct {
	Event *ChipsUnstakeRequested // Event containing the contract specifics and raw log

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
func (it *ChipsUnstakeRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsUnstakeRequested)
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
		it.Event = new(ChipsUnstakeRequested)
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
func (it *ChipsUnstakeRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsUnstakeRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsUnstakeRequested represents a UnstakeRequested event raised by the Chips contract.
type ChipsUnstakeRequested struct {
	User      common.Address
	NodeAddr  common.Address
	RequestId *big.Int
	ChipsIds  []*big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnstakeRequested is a free log retrieval operation binding the contract event 0xdb16a7c2edf86059b1faf94dd3d8c2144c14ff49e54690aaa6eb633c796b6b0c.
//
// Solidity: event UnstakeRequested(address indexed user, address indexed nodeAddr, uint256 indexed requestId, uint256[] chipsIds)
func (_Chips *ChipsFilterer) FilterUnstakeRequested(opts *bind.FilterOpts, user []common.Address, nodeAddr []common.Address, requestId []*big.Int) (*ChipsUnstakeRequestedIterator, error) {

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

	logs, sub, err := _Chips.contract.FilterLogs(opts, "UnstakeRequested", userRule, nodeAddrRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &ChipsUnstakeRequestedIterator{contract: _Chips.contract, event: "UnstakeRequested", logs: logs, sub: sub}, nil
}

// WatchUnstakeRequested is a free log subscription operation binding the contract event 0xdb16a7c2edf86059b1faf94dd3d8c2144c14ff49e54690aaa6eb633c796b6b0c.
//
// Solidity: event UnstakeRequested(address indexed user, address indexed nodeAddr, uint256 indexed requestId, uint256[] chipsIds)
func (_Chips *ChipsFilterer) WatchUnstakeRequested(opts *bind.WatchOpts, sink chan<- *ChipsUnstakeRequested, user []common.Address, nodeAddr []common.Address, requestId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Chips.contract.WatchLogs(opts, "UnstakeRequested", userRule, nodeAddrRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsUnstakeRequested)
				if err := _Chips.contract.UnpackLog(event, "UnstakeRequested", log); err != nil {
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

// ParseUnstakeRequested is a log parse operation binding the contract event 0xdb16a7c2edf86059b1faf94dd3d8c2144c14ff49e54690aaa6eb633c796b6b0c.
//
// Solidity: event UnstakeRequested(address indexed user, address indexed nodeAddr, uint256 indexed requestId, uint256[] chipsIds)
func (_Chips *ChipsFilterer) ParseUnstakeRequested(log types.Log) (*ChipsUnstakeRequested, error) {
	event := new(ChipsUnstakeRequested)
	if err := _Chips.contract.UnpackLog(event, "UnstakeRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsWithdrawRequestedIterator is returned from FilterWithdrawRequested and is used to iterate over the raw logs and unpacked data for WithdrawRequested events raised by the Chips contract.
type ChipsWithdrawRequestedIterator struct {
	Event *ChipsWithdrawRequested // Event containing the contract specifics and raw log

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
func (it *ChipsWithdrawRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsWithdrawRequested)
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
		it.Event = new(ChipsWithdrawRequested)
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
func (it *ChipsWithdrawRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsWithdrawRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsWithdrawRequested represents a WithdrawRequested event raised by the Chips contract.
type ChipsWithdrawRequested struct {
	NodeAddr  common.Address
	Amount    *big.Int
	RequestId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawRequested is a free log retrieval operation binding the contract event 0xd72eb5d043f24a0168ae744d5c44f9596fd673a26bf74d9646bff4b844882d14.
//
// Solidity: event WithdrawRequested(address indexed nodeAddr, uint256 indexed amount, uint256 indexed requestId)
func (_Chips *ChipsFilterer) FilterWithdrawRequested(opts *bind.FilterOpts, nodeAddr []common.Address, amount []*big.Int, requestId []*big.Int) (*ChipsWithdrawRequestedIterator, error) {

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

	logs, sub, err := _Chips.contract.FilterLogs(opts, "WithdrawRequested", nodeAddrRule, amountRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &ChipsWithdrawRequestedIterator{contract: _Chips.contract, event: "WithdrawRequested", logs: logs, sub: sub}, nil
}

// WatchWithdrawRequested is a free log subscription operation binding the contract event 0xd72eb5d043f24a0168ae744d5c44f9596fd673a26bf74d9646bff4b844882d14.
//
// Solidity: event WithdrawRequested(address indexed nodeAddr, uint256 indexed amount, uint256 indexed requestId)
func (_Chips *ChipsFilterer) WatchWithdrawRequested(opts *bind.WatchOpts, sink chan<- *ChipsWithdrawRequested, nodeAddr []common.Address, amount []*big.Int, requestId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Chips.contract.WatchLogs(opts, "WithdrawRequested", nodeAddrRule, amountRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsWithdrawRequested)
				if err := _Chips.contract.UnpackLog(event, "WithdrawRequested", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseWithdrawRequested(log types.Log) (*ChipsWithdrawRequested, error) {
	event := new(ChipsWithdrawRequested)
	if err := _Chips.contract.UnpackLog(event, "WithdrawRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChipsWithdrawalClaimedIterator is returned from FilterWithdrawalClaimed and is used to iterate over the raw logs and unpacked data for WithdrawalClaimed events raised by the Chips contract.
type ChipsWithdrawalClaimedIterator struct {
	Event *ChipsWithdrawalClaimed // Event containing the contract specifics and raw log

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
func (it *ChipsWithdrawalClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChipsWithdrawalClaimed)
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
		it.Event = new(ChipsWithdrawalClaimed)
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
func (it *ChipsWithdrawalClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChipsWithdrawalClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChipsWithdrawalClaimed represents a WithdrawalClaimed event raised by the Chips contract.
type ChipsWithdrawalClaimed struct {
	RequestId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalClaimed is a free log retrieval operation binding the contract event 0x8772d6f79a1845a0c0e90ef18d99f91242bbc0ba98c9ca780feaad42b81f02ba.
//
// Solidity: event WithdrawalClaimed(uint256 indexed requestId)
func (_Chips *ChipsFilterer) FilterWithdrawalClaimed(opts *bind.FilterOpts, requestId []*big.Int) (*ChipsWithdrawalClaimedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Chips.contract.FilterLogs(opts, "WithdrawalClaimed", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &ChipsWithdrawalClaimedIterator{contract: _Chips.contract, event: "WithdrawalClaimed", logs: logs, sub: sub}, nil
}

// WatchWithdrawalClaimed is a free log subscription operation binding the contract event 0x8772d6f79a1845a0c0e90ef18d99f91242bbc0ba98c9ca780feaad42b81f02ba.
//
// Solidity: event WithdrawalClaimed(uint256 indexed requestId)
func (_Chips *ChipsFilterer) WatchWithdrawalClaimed(opts *bind.WatchOpts, sink chan<- *ChipsWithdrawalClaimed, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _Chips.contract.WatchLogs(opts, "WithdrawalClaimed", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChipsWithdrawalClaimed)
				if err := _Chips.contract.UnpackLog(event, "WithdrawalClaimed", log); err != nil {
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
func (_Chips *ChipsFilterer) ParseWithdrawalClaimed(log types.Log) (*ChipsWithdrawalClaimed, error) {
	event := new(ChipsWithdrawalClaimed)
	if err := _Chips.contract.UnpackLog(event, "WithdrawalClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
