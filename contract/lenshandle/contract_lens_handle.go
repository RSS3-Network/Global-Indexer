// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package lenshandle

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

// LensHandleMetaData contains all meta data concerning the LensHandle contract.
var LensHandleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lensHub\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenGuardianCooldown\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AlreadyEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DisablingAlreadyTriggered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GuardianEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HandleContainsInvalidCharacters\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HandleFirstCharInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HandleLengthInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidParameter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEOA\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotHub\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotOwnerNorWhitelisted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyOwnerOrHub\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"handle\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"namespace\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"handleId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"HandleMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenGuardianDisablingTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"TokenGuardianStateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DANGER__disableTokenGuardian\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LENS_HUB\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NAMESPACE\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NAMESPACE_HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OWNER\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TOKEN_GUARDIAN_COOLDOWN\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableTokenGuardian\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"exists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getHandle\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getLocalName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNamespace\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNamespaceHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"}],\"name\":\"getTokenGuardianDisablingTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"localName\",\"type\":\"string\"}],\"name\":\"getTokenId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"localName\",\"type\":\"string\"}],\"name\":\"migrateHandle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"localName\",\"type\":\"string\"}],\"name\":\"mintHandle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salePrice\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"royaltiesInBasisPoints\",\"type\":\"uint256\"}],\"name\":\"setRoyalty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// LensHandleABI is the input ABI used to generate the binding from.
// Deprecated: Use LensHandleMetaData.ABI instead.
var LensHandleABI = LensHandleMetaData.ABI

// LensHandle is an auto generated Go binding around an Ethereum contract.
type LensHandle struct {
	LensHandleCaller     // Read-only binding to the contract
	LensHandleTransactor // Write-only binding to the contract
	LensHandleFilterer   // Log filterer for contract events
}

// LensHandleCaller is an auto generated read-only Go binding around an Ethereum contract.
type LensHandleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LensHandleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LensHandleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LensHandleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LensHandleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LensHandleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LensHandleSession struct {
	Contract     *LensHandle       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LensHandleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LensHandleCallerSession struct {
	Contract *LensHandleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// LensHandleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LensHandleTransactorSession struct {
	Contract     *LensHandleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// LensHandleRaw is an auto generated low-level Go binding around an Ethereum contract.
type LensHandleRaw struct {
	Contract *LensHandle // Generic contract binding to access the raw methods on
}

// LensHandleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LensHandleCallerRaw struct {
	Contract *LensHandleCaller // Generic read-only contract binding to access the raw methods on
}

// LensHandleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LensHandleTransactorRaw struct {
	Contract *LensHandleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLensHandle creates a new instance of LensHandle, bound to a specific deployed contract.
func NewLensHandle(address common.Address, backend bind.ContractBackend) (*LensHandle, error) {
	contract, err := bindLensHandle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LensHandle{LensHandleCaller: LensHandleCaller{contract: contract}, LensHandleTransactor: LensHandleTransactor{contract: contract}, LensHandleFilterer: LensHandleFilterer{contract: contract}}, nil
}

// NewLensHandleCaller creates a new read-only instance of LensHandle, bound to a specific deployed contract.
func NewLensHandleCaller(address common.Address, caller bind.ContractCaller) (*LensHandleCaller, error) {
	contract, err := bindLensHandle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LensHandleCaller{contract: contract}, nil
}

// NewLensHandleTransactor creates a new write-only instance of LensHandle, bound to a specific deployed contract.
func NewLensHandleTransactor(address common.Address, transactor bind.ContractTransactor) (*LensHandleTransactor, error) {
	contract, err := bindLensHandle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LensHandleTransactor{contract: contract}, nil
}

// NewLensHandleFilterer creates a new log filterer instance of LensHandle, bound to a specific deployed contract.
func NewLensHandleFilterer(address common.Address, filterer bind.ContractFilterer) (*LensHandleFilterer, error) {
	contract, err := bindLensHandle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LensHandleFilterer{contract: contract}, nil
}

// bindLensHandle binds a generic wrapper to an already deployed contract.
func bindLensHandle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LensHandleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LensHandle *LensHandleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LensHandle.Contract.LensHandleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LensHandle *LensHandleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LensHandle.Contract.LensHandleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LensHandle *LensHandleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LensHandle.Contract.LensHandleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LensHandle *LensHandleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LensHandle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LensHandle *LensHandleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LensHandle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LensHandle *LensHandleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LensHandle.Contract.contract.Transact(opts, method, params...)
}

// LENSHUB is a free data retrieval call binding the contract method 0x7bb9c89b.
//
// Solidity: function LENS_HUB() view returns(address)
func (_LensHandle *LensHandleCaller) LENSHUB(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "LENS_HUB")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LENSHUB is a free data retrieval call binding the contract method 0x7bb9c89b.
//
// Solidity: function LENS_HUB() view returns(address)
func (_LensHandle *LensHandleSession) LENSHUB() (common.Address, error) {
	return _LensHandle.Contract.LENSHUB(&_LensHandle.CallOpts)
}

// LENSHUB is a free data retrieval call binding the contract method 0x7bb9c89b.
//
// Solidity: function LENS_HUB() view returns(address)
func (_LensHandle *LensHandleCallerSession) LENSHUB() (common.Address, error) {
	return _LensHandle.Contract.LENSHUB(&_LensHandle.CallOpts)
}

// NAMESPACE is a free data retrieval call binding the contract method 0x44ba1fca.
//
// Solidity: function NAMESPACE() view returns(string)
func (_LensHandle *LensHandleCaller) NAMESPACE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "NAMESPACE")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NAMESPACE is a free data retrieval call binding the contract method 0x44ba1fca.
//
// Solidity: function NAMESPACE() view returns(string)
func (_LensHandle *LensHandleSession) NAMESPACE() (string, error) {
	return _LensHandle.Contract.NAMESPACE(&_LensHandle.CallOpts)
}

// NAMESPACE is a free data retrieval call binding the contract method 0x44ba1fca.
//
// Solidity: function NAMESPACE() view returns(string)
func (_LensHandle *LensHandleCallerSession) NAMESPACE() (string, error) {
	return _LensHandle.Contract.NAMESPACE(&_LensHandle.CallOpts)
}

// NAMESPACEHASH is a free data retrieval call binding the contract method 0x051fb218.
//
// Solidity: function NAMESPACE_HASH() view returns(bytes32)
func (_LensHandle *LensHandleCaller) NAMESPACEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "NAMESPACE_HASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NAMESPACEHASH is a free data retrieval call binding the contract method 0x051fb218.
//
// Solidity: function NAMESPACE_HASH() view returns(bytes32)
func (_LensHandle *LensHandleSession) NAMESPACEHASH() ([32]byte, error) {
	return _LensHandle.Contract.NAMESPACEHASH(&_LensHandle.CallOpts)
}

// NAMESPACEHASH is a free data retrieval call binding the contract method 0x051fb218.
//
// Solidity: function NAMESPACE_HASH() view returns(bytes32)
func (_LensHandle *LensHandleCallerSession) NAMESPACEHASH() ([32]byte, error) {
	return _LensHandle.Contract.NAMESPACEHASH(&_LensHandle.CallOpts)
}

// OWNER is a free data retrieval call binding the contract method 0x117803e3.
//
// Solidity: function OWNER() view returns(address)
func (_LensHandle *LensHandleCaller) OWNER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "OWNER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OWNER is a free data retrieval call binding the contract method 0x117803e3.
//
// Solidity: function OWNER() view returns(address)
func (_LensHandle *LensHandleSession) OWNER() (common.Address, error) {
	return _LensHandle.Contract.OWNER(&_LensHandle.CallOpts)
}

// OWNER is a free data retrieval call binding the contract method 0x117803e3.
//
// Solidity: function OWNER() view returns(address)
func (_LensHandle *LensHandleCallerSession) OWNER() (common.Address, error) {
	return _LensHandle.Contract.OWNER(&_LensHandle.CallOpts)
}

// TOKENGUARDIANCOOLDOWN is a free data retrieval call binding the contract method 0xa88fae83.
//
// Solidity: function TOKEN_GUARDIAN_COOLDOWN() view returns(uint256)
func (_LensHandle *LensHandleCaller) TOKENGUARDIANCOOLDOWN(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "TOKEN_GUARDIAN_COOLDOWN")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TOKENGUARDIANCOOLDOWN is a free data retrieval call binding the contract method 0xa88fae83.
//
// Solidity: function TOKEN_GUARDIAN_COOLDOWN() view returns(uint256)
func (_LensHandle *LensHandleSession) TOKENGUARDIANCOOLDOWN() (*big.Int, error) {
	return _LensHandle.Contract.TOKENGUARDIANCOOLDOWN(&_LensHandle.CallOpts)
}

// TOKENGUARDIANCOOLDOWN is a free data retrieval call binding the contract method 0xa88fae83.
//
// Solidity: function TOKEN_GUARDIAN_COOLDOWN() view returns(uint256)
func (_LensHandle *LensHandleCallerSession) TOKENGUARDIANCOOLDOWN() (*big.Int, error) {
	return _LensHandle.Contract.TOKENGUARDIANCOOLDOWN(&_LensHandle.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_LensHandle *LensHandleCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_LensHandle *LensHandleSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LensHandle.Contract.BalanceOf(&_LensHandle.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_LensHandle *LensHandleCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LensHandle.Contract.BalanceOf(&_LensHandle.CallOpts, owner)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 tokenId) view returns(bool)
func (_LensHandle *LensHandleCaller) Exists(opts *bind.CallOpts, tokenId *big.Int) (bool, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "exists", tokenId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 tokenId) view returns(bool)
func (_LensHandle *LensHandleSession) Exists(tokenId *big.Int) (bool, error) {
	return _LensHandle.Contract.Exists(&_LensHandle.CallOpts, tokenId)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 tokenId) view returns(bool)
func (_LensHandle *LensHandleCallerSession) Exists(tokenId *big.Int) (bool, error) {
	return _LensHandle.Contract.Exists(&_LensHandle.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_LensHandle *LensHandleCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_LensHandle *LensHandleSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _LensHandle.Contract.GetApproved(&_LensHandle.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_LensHandle *LensHandleCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _LensHandle.Contract.GetApproved(&_LensHandle.CallOpts, tokenId)
}

// GetHandle is a free data retrieval call binding the contract method 0xec81d194.
//
// Solidity: function getHandle(uint256 tokenId) view returns(string)
func (_LensHandle *LensHandleCaller) GetHandle(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "getHandle", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetHandle is a free data retrieval call binding the contract method 0xec81d194.
//
// Solidity: function getHandle(uint256 tokenId) view returns(string)
func (_LensHandle *LensHandleSession) GetHandle(tokenId *big.Int) (string, error) {
	return _LensHandle.Contract.GetHandle(&_LensHandle.CallOpts, tokenId)
}

// GetHandle is a free data retrieval call binding the contract method 0xec81d194.
//
// Solidity: function getHandle(uint256 tokenId) view returns(string)
func (_LensHandle *LensHandleCallerSession) GetHandle(tokenId *big.Int) (string, error) {
	return _LensHandle.Contract.GetHandle(&_LensHandle.CallOpts, tokenId)
}

// GetLocalName is a free data retrieval call binding the contract method 0x4985e504.
//
// Solidity: function getLocalName(uint256 tokenId) view returns(string)
func (_LensHandle *LensHandleCaller) GetLocalName(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "getLocalName", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetLocalName is a free data retrieval call binding the contract method 0x4985e504.
//
// Solidity: function getLocalName(uint256 tokenId) view returns(string)
func (_LensHandle *LensHandleSession) GetLocalName(tokenId *big.Int) (string, error) {
	return _LensHandle.Contract.GetLocalName(&_LensHandle.CallOpts, tokenId)
}

// GetLocalName is a free data retrieval call binding the contract method 0x4985e504.
//
// Solidity: function getLocalName(uint256 tokenId) view returns(string)
func (_LensHandle *LensHandleCallerSession) GetLocalName(tokenId *big.Int) (string, error) {
	return _LensHandle.Contract.GetLocalName(&_LensHandle.CallOpts, tokenId)
}

// GetNamespace is a free data retrieval call binding the contract method 0x27ac4b70.
//
// Solidity: function getNamespace() pure returns(string)
func (_LensHandle *LensHandleCaller) GetNamespace(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "getNamespace")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetNamespace is a free data retrieval call binding the contract method 0x27ac4b70.
//
// Solidity: function getNamespace() pure returns(string)
func (_LensHandle *LensHandleSession) GetNamespace() (string, error) {
	return _LensHandle.Contract.GetNamespace(&_LensHandle.CallOpts)
}

// GetNamespace is a free data retrieval call binding the contract method 0x27ac4b70.
//
// Solidity: function getNamespace() pure returns(string)
func (_LensHandle *LensHandleCallerSession) GetNamespace() (string, error) {
	return _LensHandle.Contract.GetNamespace(&_LensHandle.CallOpts)
}

// GetNamespaceHash is a free data retrieval call binding the contract method 0xb16f1eef.
//
// Solidity: function getNamespaceHash() pure returns(bytes32)
func (_LensHandle *LensHandleCaller) GetNamespaceHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "getNamespaceHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetNamespaceHash is a free data retrieval call binding the contract method 0xb16f1eef.
//
// Solidity: function getNamespaceHash() pure returns(bytes32)
func (_LensHandle *LensHandleSession) GetNamespaceHash() ([32]byte, error) {
	return _LensHandle.Contract.GetNamespaceHash(&_LensHandle.CallOpts)
}

// GetNamespaceHash is a free data retrieval call binding the contract method 0xb16f1eef.
//
// Solidity: function getNamespaceHash() pure returns(bytes32)
func (_LensHandle *LensHandleCallerSession) GetNamespaceHash() ([32]byte, error) {
	return _LensHandle.Contract.GetNamespaceHash(&_LensHandle.CallOpts)
}

// GetTokenGuardianDisablingTimestamp is a free data retrieval call binding the contract method 0xf3bc61f1.
//
// Solidity: function getTokenGuardianDisablingTimestamp(address wallet) view returns(uint256)
func (_LensHandle *LensHandleCaller) GetTokenGuardianDisablingTimestamp(opts *bind.CallOpts, wallet common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "getTokenGuardianDisablingTimestamp", wallet)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTokenGuardianDisablingTimestamp is a free data retrieval call binding the contract method 0xf3bc61f1.
//
// Solidity: function getTokenGuardianDisablingTimestamp(address wallet) view returns(uint256)
func (_LensHandle *LensHandleSession) GetTokenGuardianDisablingTimestamp(wallet common.Address) (*big.Int, error) {
	return _LensHandle.Contract.GetTokenGuardianDisablingTimestamp(&_LensHandle.CallOpts, wallet)
}

// GetTokenGuardianDisablingTimestamp is a free data retrieval call binding the contract method 0xf3bc61f1.
//
// Solidity: function getTokenGuardianDisablingTimestamp(address wallet) view returns(uint256)
func (_LensHandle *LensHandleCallerSession) GetTokenGuardianDisablingTimestamp(wallet common.Address) (*big.Int, error) {
	return _LensHandle.Contract.GetTokenGuardianDisablingTimestamp(&_LensHandle.CallOpts, wallet)
}

// GetTokenId is a free data retrieval call binding the contract method 0x1e7663bc.
//
// Solidity: function getTokenId(string localName) pure returns(uint256)
func (_LensHandle *LensHandleCaller) GetTokenId(opts *bind.CallOpts, localName string) (*big.Int, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "getTokenId", localName)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTokenId is a free data retrieval call binding the contract method 0x1e7663bc.
//
// Solidity: function getTokenId(string localName) pure returns(uint256)
func (_LensHandle *LensHandleSession) GetTokenId(localName string) (*big.Int, error) {
	return _LensHandle.Contract.GetTokenId(&_LensHandle.CallOpts, localName)
}

// GetTokenId is a free data retrieval call binding the contract method 0x1e7663bc.
//
// Solidity: function getTokenId(string localName) pure returns(uint256)
func (_LensHandle *LensHandleCallerSession) GetTokenId(localName string) (*big.Int, error) {
	return _LensHandle.Contract.GetTokenId(&_LensHandle.CallOpts, localName)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_LensHandle *LensHandleCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_LensHandle *LensHandleSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _LensHandle.Contract.IsApprovedForAll(&_LensHandle.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_LensHandle *LensHandleCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _LensHandle.Contract.IsApprovedForAll(&_LensHandle.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_LensHandle *LensHandleCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_LensHandle *LensHandleSession) Name() (string, error) {
	return _LensHandle.Contract.Name(&_LensHandle.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_LensHandle *LensHandleCallerSession) Name() (string, error) {
	return _LensHandle.Contract.Name(&_LensHandle.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_LensHandle *LensHandleCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_LensHandle *LensHandleSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _LensHandle.Contract.OwnerOf(&_LensHandle.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_LensHandle *LensHandleCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _LensHandle.Contract.OwnerOf(&_LensHandle.CallOpts, tokenId)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 tokenId, uint256 salePrice) view returns(address, uint256)
func (_LensHandle *LensHandleCaller) RoyaltyInfo(opts *bind.CallOpts, tokenId *big.Int, salePrice *big.Int) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "royaltyInfo", tokenId, salePrice)

	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 tokenId, uint256 salePrice) view returns(address, uint256)
func (_LensHandle *LensHandleSession) RoyaltyInfo(tokenId *big.Int, salePrice *big.Int) (common.Address, *big.Int, error) {
	return _LensHandle.Contract.RoyaltyInfo(&_LensHandle.CallOpts, tokenId, salePrice)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 tokenId, uint256 salePrice) view returns(address, uint256)
func (_LensHandle *LensHandleCallerSession) RoyaltyInfo(tokenId *big.Int, salePrice *big.Int) (common.Address, *big.Int, error) {
	return _LensHandle.Contract.RoyaltyInfo(&_LensHandle.CallOpts, tokenId, salePrice)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_LensHandle *LensHandleCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_LensHandle *LensHandleSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LensHandle.Contract.SupportsInterface(&_LensHandle.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_LensHandle *LensHandleCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LensHandle.Contract.SupportsInterface(&_LensHandle.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_LensHandle *LensHandleCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_LensHandle *LensHandleSession) Symbol() (string, error) {
	return _LensHandle.Contract.Symbol(&_LensHandle.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_LensHandle *LensHandleCallerSession) Symbol() (string, error) {
	return _LensHandle.Contract.Symbol(&_LensHandle.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_LensHandle *LensHandleCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_LensHandle *LensHandleSession) TokenURI(tokenId *big.Int) (string, error) {
	return _LensHandle.Contract.TokenURI(&_LensHandle.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_LensHandle *LensHandleCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _LensHandle.Contract.TokenURI(&_LensHandle.CallOpts, tokenId)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LensHandle *LensHandleCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LensHandle.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LensHandle *LensHandleSession) TotalSupply() (*big.Int, error) {
	return _LensHandle.Contract.TotalSupply(&_LensHandle.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LensHandle *LensHandleCallerSession) TotalSupply() (*big.Int, error) {
	return _LensHandle.Contract.TotalSupply(&_LensHandle.CallOpts)
}

// DANGERDisableTokenGuardian is a paid mutator transaction binding the contract method 0x2248f76d.
//
// Solidity: function DANGER__disableTokenGuardian() returns()
func (_LensHandle *LensHandleTransactor) DANGERDisableTokenGuardian(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "DANGER__disableTokenGuardian")
}

// DANGERDisableTokenGuardian is a paid mutator transaction binding the contract method 0x2248f76d.
//
// Solidity: function DANGER__disableTokenGuardian() returns()
func (_LensHandle *LensHandleSession) DANGERDisableTokenGuardian() (*types.Transaction, error) {
	return _LensHandle.Contract.DANGERDisableTokenGuardian(&_LensHandle.TransactOpts)
}

// DANGERDisableTokenGuardian is a paid mutator transaction binding the contract method 0x2248f76d.
//
// Solidity: function DANGER__disableTokenGuardian() returns()
func (_LensHandle *LensHandleTransactorSession) DANGERDisableTokenGuardian() (*types.Transaction, error) {
	return _LensHandle.Contract.DANGERDisableTokenGuardian(&_LensHandle.TransactOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_LensHandle *LensHandleTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_LensHandle *LensHandleSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.Contract.Approve(&_LensHandle.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_LensHandle *LensHandleTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.Contract.Approve(&_LensHandle.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_LensHandle *LensHandleTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_LensHandle *LensHandleSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.Contract.Burn(&_LensHandle.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_LensHandle *LensHandleTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.Contract.Burn(&_LensHandle.TransactOpts, tokenId)
}

// EnableTokenGuardian is a paid mutator transaction binding the contract method 0x1e9df673.
//
// Solidity: function enableTokenGuardian() returns()
func (_LensHandle *LensHandleTransactor) EnableTokenGuardian(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "enableTokenGuardian")
}

// EnableTokenGuardian is a paid mutator transaction binding the contract method 0x1e9df673.
//
// Solidity: function enableTokenGuardian() returns()
func (_LensHandle *LensHandleSession) EnableTokenGuardian() (*types.Transaction, error) {
	return _LensHandle.Contract.EnableTokenGuardian(&_LensHandle.TransactOpts)
}

// EnableTokenGuardian is a paid mutator transaction binding the contract method 0x1e9df673.
//
// Solidity: function enableTokenGuardian() returns()
func (_LensHandle *LensHandleTransactorSession) EnableTokenGuardian() (*types.Transaction, error) {
	return _LensHandle.Contract.EnableTokenGuardian(&_LensHandle.TransactOpts)
}

// MigrateHandle is a paid mutator transaction binding the contract method 0x5993bc26.
//
// Solidity: function migrateHandle(address to, string localName) returns(uint256)
func (_LensHandle *LensHandleTransactor) MigrateHandle(opts *bind.TransactOpts, to common.Address, localName string) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "migrateHandle", to, localName)
}

// MigrateHandle is a paid mutator transaction binding the contract method 0x5993bc26.
//
// Solidity: function migrateHandle(address to, string localName) returns(uint256)
func (_LensHandle *LensHandleSession) MigrateHandle(to common.Address, localName string) (*types.Transaction, error) {
	return _LensHandle.Contract.MigrateHandle(&_LensHandle.TransactOpts, to, localName)
}

// MigrateHandle is a paid mutator transaction binding the contract method 0x5993bc26.
//
// Solidity: function migrateHandle(address to, string localName) returns(uint256)
func (_LensHandle *LensHandleTransactorSession) MigrateHandle(to common.Address, localName string) (*types.Transaction, error) {
	return _LensHandle.Contract.MigrateHandle(&_LensHandle.TransactOpts, to, localName)
}

// MintHandle is a paid mutator transaction binding the contract method 0xf08e8f5e.
//
// Solidity: function mintHandle(address to, string localName) returns(uint256)
func (_LensHandle *LensHandleTransactor) MintHandle(opts *bind.TransactOpts, to common.Address, localName string) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "mintHandle", to, localName)
}

// MintHandle is a paid mutator transaction binding the contract method 0xf08e8f5e.
//
// Solidity: function mintHandle(address to, string localName) returns(uint256)
func (_LensHandle *LensHandleSession) MintHandle(to common.Address, localName string) (*types.Transaction, error) {
	return _LensHandle.Contract.MintHandle(&_LensHandle.TransactOpts, to, localName)
}

// MintHandle is a paid mutator transaction binding the contract method 0xf08e8f5e.
//
// Solidity: function mintHandle(address to, string localName) returns(uint256)
func (_LensHandle *LensHandleTransactorSession) MintHandle(to common.Address, localName string) (*types.Transaction, error) {
	return _LensHandle.Contract.MintHandle(&_LensHandle.TransactOpts, to, localName)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_LensHandle *LensHandleTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_LensHandle *LensHandleSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.Contract.SafeTransferFrom(&_LensHandle.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_LensHandle *LensHandleTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.Contract.SafeTransferFrom(&_LensHandle.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_LensHandle *LensHandleTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_LensHandle *LensHandleSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _LensHandle.Contract.SafeTransferFrom0(&_LensHandle.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_LensHandle *LensHandleTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _LensHandle.Contract.SafeTransferFrom0(&_LensHandle.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_LensHandle *LensHandleTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_LensHandle *LensHandleSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _LensHandle.Contract.SetApprovalForAll(&_LensHandle.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_LensHandle *LensHandleTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _LensHandle.Contract.SetApprovalForAll(&_LensHandle.TransactOpts, operator, approved)
}

// SetRoyalty is a paid mutator transaction binding the contract method 0x4209a2e1.
//
// Solidity: function setRoyalty(uint256 royaltiesInBasisPoints) returns()
func (_LensHandle *LensHandleTransactor) SetRoyalty(opts *bind.TransactOpts, royaltiesInBasisPoints *big.Int) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "setRoyalty", royaltiesInBasisPoints)
}

// SetRoyalty is a paid mutator transaction binding the contract method 0x4209a2e1.
//
// Solidity: function setRoyalty(uint256 royaltiesInBasisPoints) returns()
func (_LensHandle *LensHandleSession) SetRoyalty(royaltiesInBasisPoints *big.Int) (*types.Transaction, error) {
	return _LensHandle.Contract.SetRoyalty(&_LensHandle.TransactOpts, royaltiesInBasisPoints)
}

// SetRoyalty is a paid mutator transaction binding the contract method 0x4209a2e1.
//
// Solidity: function setRoyalty(uint256 royaltiesInBasisPoints) returns()
func (_LensHandle *LensHandleTransactorSession) SetRoyalty(royaltiesInBasisPoints *big.Int) (*types.Transaction, error) {
	return _LensHandle.Contract.SetRoyalty(&_LensHandle.TransactOpts, royaltiesInBasisPoints)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_LensHandle *LensHandleTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_LensHandle *LensHandleSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.Contract.TransferFrom(&_LensHandle.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_LensHandle *LensHandleTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LensHandle.Contract.TransferFrom(&_LensHandle.TransactOpts, from, to, tokenId)
}

// LensHandleApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the LensHandle contract.
type LensHandleApprovalIterator struct {
	Event *LensHandleApproval // Event containing the contract specifics and raw log

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
func (it *LensHandleApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LensHandleApproval)
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
		it.Event = new(LensHandleApproval)
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
func (it *LensHandleApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LensHandleApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LensHandleApproval represents a Approval event raised by the LensHandle contract.
type LensHandleApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_LensHandle *LensHandleFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*LensHandleApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LensHandle.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &LensHandleApprovalIterator{contract: _LensHandle.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_LensHandle *LensHandleFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LensHandleApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LensHandle.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LensHandleApproval)
				if err := _LensHandle.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_LensHandle *LensHandleFilterer) ParseApproval(log types.Log) (*LensHandleApproval, error) {
	event := new(LensHandleApproval)
	if err := _LensHandle.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LensHandleApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the LensHandle contract.
type LensHandleApprovalForAllIterator struct {
	Event *LensHandleApprovalForAll // Event containing the contract specifics and raw log

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
func (it *LensHandleApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LensHandleApprovalForAll)
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
		it.Event = new(LensHandleApprovalForAll)
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
func (it *LensHandleApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LensHandleApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LensHandleApprovalForAll represents a ApprovalForAll event raised by the LensHandle contract.
type LensHandleApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_LensHandle *LensHandleFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*LensHandleApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _LensHandle.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &LensHandleApprovalForAllIterator{contract: _LensHandle.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_LensHandle *LensHandleFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *LensHandleApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _LensHandle.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LensHandleApprovalForAll)
				if err := _LensHandle.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_LensHandle *LensHandleFilterer) ParseApprovalForAll(log types.Log) (*LensHandleApprovalForAll, error) {
	event := new(LensHandleApprovalForAll)
	if err := _LensHandle.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LensHandleHandleMintedIterator is returned from FilterHandleMinted and is used to iterate over the raw logs and unpacked data for HandleMinted events raised by the LensHandle contract.
type LensHandleHandleMintedIterator struct {
	Event *LensHandleHandleMinted // Event containing the contract specifics and raw log

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
func (it *LensHandleHandleMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LensHandleHandleMinted)
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
		it.Event = new(LensHandleHandleMinted)
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
func (it *LensHandleHandleMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LensHandleHandleMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LensHandleHandleMinted represents a HandleMinted event raised by the LensHandle contract.
type LensHandleHandleMinted struct {
	Handle    string
	Namespace string
	HandleId  *big.Int
	To        common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterHandleMinted is a free log retrieval operation binding the contract event 0x30a132e912787e50de6193fe56a96ea6188c0bbf676679d630a25d3293c3e19a.
//
// Solidity: event HandleMinted(string handle, string namespace, uint256 handleId, address to, uint256 timestamp)
func (_LensHandle *LensHandleFilterer) FilterHandleMinted(opts *bind.FilterOpts) (*LensHandleHandleMintedIterator, error) {

	logs, sub, err := _LensHandle.contract.FilterLogs(opts, "HandleMinted")
	if err != nil {
		return nil, err
	}
	return &LensHandleHandleMintedIterator{contract: _LensHandle.contract, event: "HandleMinted", logs: logs, sub: sub}, nil
}

// WatchHandleMinted is a free log subscription operation binding the contract event 0x30a132e912787e50de6193fe56a96ea6188c0bbf676679d630a25d3293c3e19a.
//
// Solidity: event HandleMinted(string handle, string namespace, uint256 handleId, address to, uint256 timestamp)
func (_LensHandle *LensHandleFilterer) WatchHandleMinted(opts *bind.WatchOpts, sink chan<- *LensHandleHandleMinted) (event.Subscription, error) {

	logs, sub, err := _LensHandle.contract.WatchLogs(opts, "HandleMinted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LensHandleHandleMinted)
				if err := _LensHandle.contract.UnpackLog(event, "HandleMinted", log); err != nil {
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

// ParseHandleMinted is a log parse operation binding the contract event 0x30a132e912787e50de6193fe56a96ea6188c0bbf676679d630a25d3293c3e19a.
//
// Solidity: event HandleMinted(string handle, string namespace, uint256 handleId, address to, uint256 timestamp)
func (_LensHandle *LensHandleFilterer) ParseHandleMinted(log types.Log) (*LensHandleHandleMinted, error) {
	event := new(LensHandleHandleMinted)
	if err := _LensHandle.contract.UnpackLog(event, "HandleMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LensHandleTokenGuardianStateChangedIterator is returned from FilterTokenGuardianStateChanged and is used to iterate over the raw logs and unpacked data for TokenGuardianStateChanged events raised by the LensHandle contract.
type LensHandleTokenGuardianStateChangedIterator struct {
	Event *LensHandleTokenGuardianStateChanged // Event containing the contract specifics and raw log

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
func (it *LensHandleTokenGuardianStateChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LensHandleTokenGuardianStateChanged)
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
		it.Event = new(LensHandleTokenGuardianStateChanged)
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
func (it *LensHandleTokenGuardianStateChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LensHandleTokenGuardianStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LensHandleTokenGuardianStateChanged represents a TokenGuardianStateChanged event raised by the LensHandle contract.
type LensHandleTokenGuardianStateChanged struct {
	Wallet                          common.Address
	Enabled                         bool
	TokenGuardianDisablingTimestamp *big.Int
	Timestamp                       *big.Int
	Raw                             types.Log // Blockchain specific contextual infos
}

// FilterTokenGuardianStateChanged is a free log retrieval operation binding the contract event 0x035adf3bbe16b317cf4a3e05c966ea6571d1af00147c5f121bd1514b1e322a06.
//
// Solidity: event TokenGuardianStateChanged(address indexed wallet, bool indexed enabled, uint256 tokenGuardianDisablingTimestamp, uint256 timestamp)
func (_LensHandle *LensHandleFilterer) FilterTokenGuardianStateChanged(opts *bind.FilterOpts, wallet []common.Address, enabled []bool) (*LensHandleTokenGuardianStateChangedIterator, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}
	var enabledRule []interface{}
	for _, enabledItem := range enabled {
		enabledRule = append(enabledRule, enabledItem)
	}

	logs, sub, err := _LensHandle.contract.FilterLogs(opts, "TokenGuardianStateChanged", walletRule, enabledRule)
	if err != nil {
		return nil, err
	}
	return &LensHandleTokenGuardianStateChangedIterator{contract: _LensHandle.contract, event: "TokenGuardianStateChanged", logs: logs, sub: sub}, nil
}

// WatchTokenGuardianStateChanged is a free log subscription operation binding the contract event 0x035adf3bbe16b317cf4a3e05c966ea6571d1af00147c5f121bd1514b1e322a06.
//
// Solidity: event TokenGuardianStateChanged(address indexed wallet, bool indexed enabled, uint256 tokenGuardianDisablingTimestamp, uint256 timestamp)
func (_LensHandle *LensHandleFilterer) WatchTokenGuardianStateChanged(opts *bind.WatchOpts, sink chan<- *LensHandleTokenGuardianStateChanged, wallet []common.Address, enabled []bool) (event.Subscription, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}
	var enabledRule []interface{}
	for _, enabledItem := range enabled {
		enabledRule = append(enabledRule, enabledItem)
	}

	logs, sub, err := _LensHandle.contract.WatchLogs(opts, "TokenGuardianStateChanged", walletRule, enabledRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LensHandleTokenGuardianStateChanged)
				if err := _LensHandle.contract.UnpackLog(event, "TokenGuardianStateChanged", log); err != nil {
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

// ParseTokenGuardianStateChanged is a log parse operation binding the contract event 0x035adf3bbe16b317cf4a3e05c966ea6571d1af00147c5f121bd1514b1e322a06.
//
// Solidity: event TokenGuardianStateChanged(address indexed wallet, bool indexed enabled, uint256 tokenGuardianDisablingTimestamp, uint256 timestamp)
func (_LensHandle *LensHandleFilterer) ParseTokenGuardianStateChanged(log types.Log) (*LensHandleTokenGuardianStateChanged, error) {
	event := new(LensHandleTokenGuardianStateChanged)
	if err := _LensHandle.contract.UnpackLog(event, "TokenGuardianStateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LensHandleTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the LensHandle contract.
type LensHandleTransferIterator struct {
	Event *LensHandleTransfer // Event containing the contract specifics and raw log

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
func (it *LensHandleTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LensHandleTransfer)
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
		it.Event = new(LensHandleTransfer)
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
func (it *LensHandleTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LensHandleTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LensHandleTransfer represents a Transfer event raised by the LensHandle contract.
type LensHandleTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_LensHandle *LensHandleFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*LensHandleTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LensHandle.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &LensHandleTransferIterator{contract: _LensHandle.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_LensHandle *LensHandleFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LensHandleTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LensHandle.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LensHandleTransfer)
				if err := _LensHandle.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_LensHandle *LensHandleFilterer) ParseTransfer(log types.Log) (*LensHandleTransfer, error) {
	event := new(LensHandleTransfer)
	if err := _LensHandle.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
