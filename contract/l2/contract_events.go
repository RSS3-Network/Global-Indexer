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

// EventsMetaData contains all meta data concerning the Events contract.
var EventsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"event\",\"name\":\"ChipsMerged\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newTokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"burnedTokenIds\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DemotionRevoked\",\"inputs\":[{\"name\":\"demotionId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DemotionSubmitted\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"demotionId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"reporter\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposited\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeCreated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"publicGood\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"alpha\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeStatusChanged\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"curStatus\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumNodeStatus\"},{\"name\":\"newStatus\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumNodeStatus\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeTaxRateBasisPointsSet\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeUpdated\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PublicGoodRewardDistributed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"publicPoolRewards\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"publicPoolTax\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PublicPoolTaxRateBasisPointsSet\",\"inputs\":[{\"name\":\"taxRateBasisPoints\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardDistributed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"nodeAddrs\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"operationRewards\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"stakingRewards\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"taxCollected\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"requestCounts\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashCommitted\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashRecorded\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"slashedOperationPool\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"slashedStakingPool\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashRevoked\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Staked\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"startTokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnstakeClaimed\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"unstakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnstakeRequested\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"unstakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"chipsIds\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawRequested\",\"inputs\":[{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalClaimed\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"nodeAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// EventsABI is the input ABI used to generate the binding from.
// Deprecated: Use EventsMetaData.ABI instead.
var EventsABI = EventsMetaData.ABI

// Events is an auto generated Go binding around an Ethereum contract.
type Events struct {
	EventsCaller     // Read-only binding to the contract
	EventsTransactor // Write-only binding to the contract
	EventsFilterer   // Log filterer for contract events
}

// EventsCaller is an auto generated read-only Go binding around an Ethereum contract.
type EventsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EventsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EventsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EventsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EventsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EventsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EventsSession struct {
	Contract     *Events           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EventsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EventsCallerSession struct {
	Contract *EventsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EventsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EventsTransactorSession struct {
	Contract     *EventsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EventsRaw is an auto generated low-level Go binding around an Ethereum contract.
type EventsRaw struct {
	Contract *Events // Generic contract binding to access the raw methods on
}

// EventsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EventsCallerRaw struct {
	Contract *EventsCaller // Generic read-only contract binding to access the raw methods on
}

// EventsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EventsTransactorRaw struct {
	Contract *EventsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEvents creates a new instance of Events, bound to a specific deployed contract.
func NewEvents(address common.Address, backend bind.ContractBackend) (*Events, error) {
	contract, err := bindEvents(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Events{EventsCaller: EventsCaller{contract: contract}, EventsTransactor: EventsTransactor{contract: contract}, EventsFilterer: EventsFilterer{contract: contract}}, nil
}

// NewEventsCaller creates a new read-only instance of Events, bound to a specific deployed contract.
func NewEventsCaller(address common.Address, caller bind.ContractCaller) (*EventsCaller, error) {
	contract, err := bindEvents(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EventsCaller{contract: contract}, nil
}

// NewEventsTransactor creates a new write-only instance of Events, bound to a specific deployed contract.
func NewEventsTransactor(address common.Address, transactor bind.ContractTransactor) (*EventsTransactor, error) {
	contract, err := bindEvents(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EventsTransactor{contract: contract}, nil
}

// NewEventsFilterer creates a new log filterer instance of Events, bound to a specific deployed contract.
func NewEventsFilterer(address common.Address, filterer bind.ContractFilterer) (*EventsFilterer, error) {
	contract, err := bindEvents(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EventsFilterer{contract: contract}, nil
}

// bindEvents binds a generic wrapper to an already deployed contract.
func bindEvents(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EventsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Events *EventsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Events.Contract.EventsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Events *EventsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Events.Contract.EventsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Events *EventsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Events.Contract.EventsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Events *EventsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Events.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Events *EventsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Events.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Events *EventsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Events.Contract.contract.Transact(opts, method, params...)
}

// EventsChipsMergedIterator is returned from FilterChipsMerged and is used to iterate over the raw logs and unpacked data for ChipsMerged events raised by the Events contract.
type EventsChipsMergedIterator struct {
	Event *EventsChipsMerged // Event containing the contract specifics and raw log

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
func (it *EventsChipsMergedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsChipsMerged)
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
		it.Event = new(EventsChipsMerged)
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
func (it *EventsChipsMergedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsChipsMergedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsChipsMerged represents a ChipsMerged event raised by the Events contract.
type EventsChipsMerged struct {
	User           common.Address
	NodeAddr       common.Address
	NewTokenId     *big.Int
	BurnedTokenIds []*big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterChipsMerged is a free log retrieval operation binding the contract event 0x18f7a8d9091ae36e0b14acaa3e2ac8a6672a389ef8e18c8e3523fcae05d24f18.
//
// Solidity: event ChipsMerged(address indexed user, address indexed nodeAddr, uint256 indexed newTokenId, uint256[] burnedTokenIds)
func (_Events *EventsFilterer) FilterChipsMerged(opts *bind.FilterOpts, user []common.Address, nodeAddr []common.Address, newTokenId []*big.Int) (*EventsChipsMergedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var newTokenIdRule []interface{}
	for _, newTokenIdItem := range newTokenId {
		newTokenIdRule = append(newTokenIdRule, newTokenIdItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "ChipsMerged", userRule, nodeAddrRule, newTokenIdRule)
	if err != nil {
		return nil, err
	}
	return &EventsChipsMergedIterator{contract: _Events.contract, event: "ChipsMerged", logs: logs, sub: sub}, nil
}

// WatchChipsMerged is a free log subscription operation binding the contract event 0x18f7a8d9091ae36e0b14acaa3e2ac8a6672a389ef8e18c8e3523fcae05d24f18.
//
// Solidity: event ChipsMerged(address indexed user, address indexed nodeAddr, uint256 indexed newTokenId, uint256[] burnedTokenIds)
func (_Events *EventsFilterer) WatchChipsMerged(opts *bind.WatchOpts, sink chan<- *EventsChipsMerged, user []common.Address, nodeAddr []common.Address, newTokenId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var newTokenIdRule []interface{}
	for _, newTokenIdItem := range newTokenId {
		newTokenIdRule = append(newTokenIdRule, newTokenIdItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "ChipsMerged", userRule, nodeAddrRule, newTokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsChipsMerged)
				if err := _Events.contract.UnpackLog(event, "ChipsMerged", log); err != nil {
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

// ParseChipsMerged is a log parse operation binding the contract event 0x18f7a8d9091ae36e0b14acaa3e2ac8a6672a389ef8e18c8e3523fcae05d24f18.
//
// Solidity: event ChipsMerged(address indexed user, address indexed nodeAddr, uint256 indexed newTokenId, uint256[] burnedTokenIds)
func (_Events *EventsFilterer) ParseChipsMerged(log types.Log) (*EventsChipsMerged, error) {
	event := new(EventsChipsMerged)
	if err := _Events.contract.UnpackLog(event, "ChipsMerged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsDemotionRevokedIterator is returned from FilterDemotionRevoked and is used to iterate over the raw logs and unpacked data for DemotionRevoked events raised by the Events contract.
type EventsDemotionRevokedIterator struct {
	Event *EventsDemotionRevoked // Event containing the contract specifics and raw log

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
func (it *EventsDemotionRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsDemotionRevoked)
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
		it.Event = new(EventsDemotionRevoked)
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
func (it *EventsDemotionRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsDemotionRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsDemotionRevoked represents a DemotionRevoked event raised by the Events contract.
type EventsDemotionRevoked struct {
	DemotionId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDemotionRevoked is a free log retrieval operation binding the contract event 0xe07fc9955958541e451cd076ec11a41f7d3d69411be92715d1bd412050ad5918.
//
// Solidity: event DemotionRevoked(uint256 indexed demotionId)
func (_Events *EventsFilterer) FilterDemotionRevoked(opts *bind.FilterOpts, demotionId []*big.Int) (*EventsDemotionRevokedIterator, error) {

	var demotionIdRule []interface{}
	for _, demotionIdItem := range demotionId {
		demotionIdRule = append(demotionIdRule, demotionIdItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "DemotionRevoked", demotionIdRule)
	if err != nil {
		return nil, err
	}
	return &EventsDemotionRevokedIterator{contract: _Events.contract, event: "DemotionRevoked", logs: logs, sub: sub}, nil
}

// WatchDemotionRevoked is a free log subscription operation binding the contract event 0xe07fc9955958541e451cd076ec11a41f7d3d69411be92715d1bd412050ad5918.
//
// Solidity: event DemotionRevoked(uint256 indexed demotionId)
func (_Events *EventsFilterer) WatchDemotionRevoked(opts *bind.WatchOpts, sink chan<- *EventsDemotionRevoked, demotionId []*big.Int) (event.Subscription, error) {

	var demotionIdRule []interface{}
	for _, demotionIdItem := range demotionId {
		demotionIdRule = append(demotionIdRule, demotionIdItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "DemotionRevoked", demotionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsDemotionRevoked)
				if err := _Events.contract.UnpackLog(event, "DemotionRevoked", log); err != nil {
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

// ParseDemotionRevoked is a log parse operation binding the contract event 0xe07fc9955958541e451cd076ec11a41f7d3d69411be92715d1bd412050ad5918.
//
// Solidity: event DemotionRevoked(uint256 indexed demotionId)
func (_Events *EventsFilterer) ParseDemotionRevoked(log types.Log) (*EventsDemotionRevoked, error) {
	event := new(EventsDemotionRevoked)
	if err := _Events.contract.UnpackLog(event, "DemotionRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsDemotionSubmittedIterator is returned from FilterDemotionSubmitted and is used to iterate over the raw logs and unpacked data for DemotionSubmitted events raised by the Events contract.
type EventsDemotionSubmittedIterator struct {
	Event *EventsDemotionSubmitted // Event containing the contract specifics and raw log

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
func (it *EventsDemotionSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsDemotionSubmitted)
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
		it.Event = new(EventsDemotionSubmitted)
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
func (it *EventsDemotionSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsDemotionSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsDemotionSubmitted represents a DemotionSubmitted event raised by the Events contract.
type EventsDemotionSubmitted struct {
	Epoch      *big.Int
	NodeAddr   common.Address
	DemotionId *big.Int
	Reason     string
	Reporter   common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDemotionSubmitted is a free log retrieval operation binding the contract event 0x72cf50a64fcd4d3043cad57179bf0016d93d285e61b4fd0565a0e92b763deee9.
//
// Solidity: event DemotionSubmitted(uint256 indexed epoch, address indexed nodeAddr, uint256 indexed demotionId, string reason, address reporter)
func (_Events *EventsFilterer) FilterDemotionSubmitted(opts *bind.FilterOpts, epoch []*big.Int, nodeAddr []common.Address, demotionId []*big.Int) (*EventsDemotionSubmittedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var demotionIdRule []interface{}
	for _, demotionIdItem := range demotionId {
		demotionIdRule = append(demotionIdRule, demotionIdItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "DemotionSubmitted", epochRule, nodeAddrRule, demotionIdRule)
	if err != nil {
		return nil, err
	}
	return &EventsDemotionSubmittedIterator{contract: _Events.contract, event: "DemotionSubmitted", logs: logs, sub: sub}, nil
}

// WatchDemotionSubmitted is a free log subscription operation binding the contract event 0x72cf50a64fcd4d3043cad57179bf0016d93d285e61b4fd0565a0e92b763deee9.
//
// Solidity: event DemotionSubmitted(uint256 indexed epoch, address indexed nodeAddr, uint256 indexed demotionId, string reason, address reporter)
func (_Events *EventsFilterer) WatchDemotionSubmitted(opts *bind.WatchOpts, sink chan<- *EventsDemotionSubmitted, epoch []*big.Int, nodeAddr []common.Address, demotionId []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var demotionIdRule []interface{}
	for _, demotionIdItem := range demotionId {
		demotionIdRule = append(demotionIdRule, demotionIdItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "DemotionSubmitted", epochRule, nodeAddrRule, demotionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsDemotionSubmitted)
				if err := _Events.contract.UnpackLog(event, "DemotionSubmitted", log); err != nil {
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

// ParseDemotionSubmitted is a log parse operation binding the contract event 0x72cf50a64fcd4d3043cad57179bf0016d93d285e61b4fd0565a0e92b763deee9.
//
// Solidity: event DemotionSubmitted(uint256 indexed epoch, address indexed nodeAddr, uint256 indexed demotionId, string reason, address reporter)
func (_Events *EventsFilterer) ParseDemotionSubmitted(log types.Log) (*EventsDemotionSubmitted, error) {
	event := new(EventsDemotionSubmitted)
	if err := _Events.contract.UnpackLog(event, "DemotionSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the Events contract.
type EventsDepositedIterator struct {
	Event *EventsDeposited // Event containing the contract specifics and raw log

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
func (it *EventsDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsDeposited)
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
		it.Event = new(EventsDeposited)
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
func (it *EventsDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsDeposited represents a Deposited event raised by the Events contract.
type EventsDeposited struct {
	NodeAddr common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed nodeAddr, uint256 indexed amount)
func (_Events *EventsFilterer) FilterDeposited(opts *bind.FilterOpts, nodeAddr []common.Address, amount []*big.Int) (*EventsDepositedIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "Deposited", nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &EventsDepositedIterator{contract: _Events.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed nodeAddr, uint256 indexed amount)
func (_Events *EventsFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *EventsDeposited, nodeAddr []common.Address, amount []*big.Int) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "Deposited", nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsDeposited)
				if err := _Events.contract.UnpackLog(event, "Deposited", log); err != nil {
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
func (_Events *EventsFilterer) ParseDeposited(log types.Log) (*EventsDeposited, error) {
	event := new(EventsDeposited)
	if err := _Events.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsNodeCreatedIterator is returned from FilterNodeCreated and is used to iterate over the raw logs and unpacked data for NodeCreated events raised by the Events contract.
type EventsNodeCreatedIterator struct {
	Event *EventsNodeCreated // Event containing the contract specifics and raw log

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
func (it *EventsNodeCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsNodeCreated)
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
		it.Event = new(EventsNodeCreated)
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
func (it *EventsNodeCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsNodeCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsNodeCreated represents a NodeCreated event raised by the Events contract.
type EventsNodeCreated struct {
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
func (_Events *EventsFilterer) FilterNodeCreated(opts *bind.FilterOpts, nodeId []*big.Int, nodeAddr []common.Address) (*EventsNodeCreatedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "NodeCreated", nodeIdRule, nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return &EventsNodeCreatedIterator{contract: _Events.contract, event: "NodeCreated", logs: logs, sub: sub}, nil
}

// WatchNodeCreated is a free log subscription operation binding the contract event 0x37570f68d94fd46cd4009b3823da2b2bc1a9a7e38f824f311ede9e876816e321.
//
// Solidity: event NodeCreated(uint256 indexed nodeId, address indexed nodeAddr, string name, string description, uint64 taxRateBasisPoints, bool publicGood, bool alpha)
func (_Events *EventsFilterer) WatchNodeCreated(opts *bind.WatchOpts, sink chan<- *EventsNodeCreated, nodeId []*big.Int, nodeAddr []common.Address) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "NodeCreated", nodeIdRule, nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsNodeCreated)
				if err := _Events.contract.UnpackLog(event, "NodeCreated", log); err != nil {
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
func (_Events *EventsFilterer) ParseNodeCreated(log types.Log) (*EventsNodeCreated, error) {
	event := new(EventsNodeCreated)
	if err := _Events.contract.UnpackLog(event, "NodeCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsNodeStatusChangedIterator is returned from FilterNodeStatusChanged and is used to iterate over the raw logs and unpacked data for NodeStatusChanged events raised by the Events contract.
type EventsNodeStatusChangedIterator struct {
	Event *EventsNodeStatusChanged // Event containing the contract specifics and raw log

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
func (it *EventsNodeStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsNodeStatusChanged)
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
		it.Event = new(EventsNodeStatusChanged)
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
func (it *EventsNodeStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsNodeStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsNodeStatusChanged represents a NodeStatusChanged event raised by the Events contract.
type EventsNodeStatusChanged struct {
	NodeAddr  common.Address
	CurStatus uint8
	NewStatus uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNodeStatusChanged is a free log retrieval operation binding the contract event 0xfb5740b379943f137d27260c0f7bd5f908f4d60a4507fd1c4824d264a00f0a72.
//
// Solidity: event NodeStatusChanged(address indexed nodeAddr, uint8 indexed curStatus, uint8 indexed newStatus)
func (_Events *EventsFilterer) FilterNodeStatusChanged(opts *bind.FilterOpts, nodeAddr []common.Address, curStatus []uint8, newStatus []uint8) (*EventsNodeStatusChangedIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var curStatusRule []interface{}
	for _, curStatusItem := range curStatus {
		curStatusRule = append(curStatusRule, curStatusItem)
	}
	var newStatusRule []interface{}
	for _, newStatusItem := range newStatus {
		newStatusRule = append(newStatusRule, newStatusItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "NodeStatusChanged", nodeAddrRule, curStatusRule, newStatusRule)
	if err != nil {
		return nil, err
	}
	return &EventsNodeStatusChangedIterator{contract: _Events.contract, event: "NodeStatusChanged", logs: logs, sub: sub}, nil
}

// WatchNodeStatusChanged is a free log subscription operation binding the contract event 0xfb5740b379943f137d27260c0f7bd5f908f4d60a4507fd1c4824d264a00f0a72.
//
// Solidity: event NodeStatusChanged(address indexed nodeAddr, uint8 indexed curStatus, uint8 indexed newStatus)
func (_Events *EventsFilterer) WatchNodeStatusChanged(opts *bind.WatchOpts, sink chan<- *EventsNodeStatusChanged, nodeAddr []common.Address, curStatus []uint8, newStatus []uint8) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var curStatusRule []interface{}
	for _, curStatusItem := range curStatus {
		curStatusRule = append(curStatusRule, curStatusItem)
	}
	var newStatusRule []interface{}
	for _, newStatusItem := range newStatus {
		newStatusRule = append(newStatusRule, newStatusItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "NodeStatusChanged", nodeAddrRule, curStatusRule, newStatusRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsNodeStatusChanged)
				if err := _Events.contract.UnpackLog(event, "NodeStatusChanged", log); err != nil {
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

// ParseNodeStatusChanged is a log parse operation binding the contract event 0xfb5740b379943f137d27260c0f7bd5f908f4d60a4507fd1c4824d264a00f0a72.
//
// Solidity: event NodeStatusChanged(address indexed nodeAddr, uint8 indexed curStatus, uint8 indexed newStatus)
func (_Events *EventsFilterer) ParseNodeStatusChanged(log types.Log) (*EventsNodeStatusChanged, error) {
	event := new(EventsNodeStatusChanged)
	if err := _Events.contract.UnpackLog(event, "NodeStatusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsNodeTaxRateBasisPointsSetIterator is returned from FilterNodeTaxRateBasisPointsSet and is used to iterate over the raw logs and unpacked data for NodeTaxRateBasisPointsSet events raised by the Events contract.
type EventsNodeTaxRateBasisPointsSetIterator struct {
	Event *EventsNodeTaxRateBasisPointsSet // Event containing the contract specifics and raw log

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
func (it *EventsNodeTaxRateBasisPointsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsNodeTaxRateBasisPointsSet)
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
		it.Event = new(EventsNodeTaxRateBasisPointsSet)
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
func (it *EventsNodeTaxRateBasisPointsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsNodeTaxRateBasisPointsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsNodeTaxRateBasisPointsSet represents a NodeTaxRateBasisPointsSet event raised by the Events contract.
type EventsNodeTaxRateBasisPointsSet struct {
	NodeAddr           common.Address
	TaxRateBasisPoints uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterNodeTaxRateBasisPointsSet is a free log retrieval operation binding the contract event 0xb8e5551053b871a40f7c7382e5bd3af5a62dd737d059d3838cf3aa7c325bd479.
//
// Solidity: event NodeTaxRateBasisPointsSet(address indexed nodeAddr, uint64 indexed taxRateBasisPoints)
func (_Events *EventsFilterer) FilterNodeTaxRateBasisPointsSet(opts *bind.FilterOpts, nodeAddr []common.Address, taxRateBasisPoints []uint64) (*EventsNodeTaxRateBasisPointsSetIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "NodeTaxRateBasisPointsSet", nodeAddrRule, taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return &EventsNodeTaxRateBasisPointsSetIterator{contract: _Events.contract, event: "NodeTaxRateBasisPointsSet", logs: logs, sub: sub}, nil
}

// WatchNodeTaxRateBasisPointsSet is a free log subscription operation binding the contract event 0xb8e5551053b871a40f7c7382e5bd3af5a62dd737d059d3838cf3aa7c325bd479.
//
// Solidity: event NodeTaxRateBasisPointsSet(address indexed nodeAddr, uint64 indexed taxRateBasisPoints)
func (_Events *EventsFilterer) WatchNodeTaxRateBasisPointsSet(opts *bind.WatchOpts, sink chan<- *EventsNodeTaxRateBasisPointsSet, nodeAddr []common.Address, taxRateBasisPoints []uint64) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "NodeTaxRateBasisPointsSet", nodeAddrRule, taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsNodeTaxRateBasisPointsSet)
				if err := _Events.contract.UnpackLog(event, "NodeTaxRateBasisPointsSet", log); err != nil {
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
func (_Events *EventsFilterer) ParseNodeTaxRateBasisPointsSet(log types.Log) (*EventsNodeTaxRateBasisPointsSet, error) {
	event := new(EventsNodeTaxRateBasisPointsSet)
	if err := _Events.contract.UnpackLog(event, "NodeTaxRateBasisPointsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsNodeUpdatedIterator is returned from FilterNodeUpdated and is used to iterate over the raw logs and unpacked data for NodeUpdated events raised by the Events contract.
type EventsNodeUpdatedIterator struct {
	Event *EventsNodeUpdated // Event containing the contract specifics and raw log

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
func (it *EventsNodeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsNodeUpdated)
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
		it.Event = new(EventsNodeUpdated)
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
func (it *EventsNodeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsNodeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsNodeUpdated represents a NodeUpdated event raised by the Events contract.
type EventsNodeUpdated struct {
	NodeAddr    common.Address
	Name        string
	Description string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNodeUpdated is a free log retrieval operation binding the contract event 0x8dd72b3e159c2446f32a80f24459ad76e9f8fbb74165952a01c27adb16aba725.
//
// Solidity: event NodeUpdated(address indexed nodeAddr, string name, string description)
func (_Events *EventsFilterer) FilterNodeUpdated(opts *bind.FilterOpts, nodeAddr []common.Address) (*EventsNodeUpdatedIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "NodeUpdated", nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return &EventsNodeUpdatedIterator{contract: _Events.contract, event: "NodeUpdated", logs: logs, sub: sub}, nil
}

// WatchNodeUpdated is a free log subscription operation binding the contract event 0x8dd72b3e159c2446f32a80f24459ad76e9f8fbb74165952a01c27adb16aba725.
//
// Solidity: event NodeUpdated(address indexed nodeAddr, string name, string description)
func (_Events *EventsFilterer) WatchNodeUpdated(opts *bind.WatchOpts, sink chan<- *EventsNodeUpdated, nodeAddr []common.Address) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "NodeUpdated", nodeAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsNodeUpdated)
				if err := _Events.contract.UnpackLog(event, "NodeUpdated", log); err != nil {
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

// ParseNodeUpdated is a log parse operation binding the contract event 0x8dd72b3e159c2446f32a80f24459ad76e9f8fbb74165952a01c27adb16aba725.
//
// Solidity: event NodeUpdated(address indexed nodeAddr, string name, string description)
func (_Events *EventsFilterer) ParseNodeUpdated(log types.Log) (*EventsNodeUpdated, error) {
	event := new(EventsNodeUpdated)
	if err := _Events.contract.UnpackLog(event, "NodeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsPublicGoodRewardDistributedIterator is returned from FilterPublicGoodRewardDistributed and is used to iterate over the raw logs and unpacked data for PublicGoodRewardDistributed events raised by the Events contract.
type EventsPublicGoodRewardDistributedIterator struct {
	Event *EventsPublicGoodRewardDistributed // Event containing the contract specifics and raw log

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
func (it *EventsPublicGoodRewardDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsPublicGoodRewardDistributed)
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
		it.Event = new(EventsPublicGoodRewardDistributed)
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
func (it *EventsPublicGoodRewardDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsPublicGoodRewardDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsPublicGoodRewardDistributed represents a PublicGoodRewardDistributed event raised by the Events contract.
type EventsPublicGoodRewardDistributed struct {
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
func (_Events *EventsFilterer) FilterPublicGoodRewardDistributed(opts *bind.FilterOpts, epoch []*big.Int) (*EventsPublicGoodRewardDistributedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "PublicGoodRewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return &EventsPublicGoodRewardDistributedIterator{contract: _Events.contract, event: "PublicGoodRewardDistributed", logs: logs, sub: sub}, nil
}

// WatchPublicGoodRewardDistributed is a free log subscription operation binding the contract event 0xab7d25a2f6206ef56c88807f2474ddcd97e1a6323cb25149cde3a607fed6f2d7.
//
// Solidity: event PublicGoodRewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, uint256 publicPoolRewards, uint256 publicPoolTax)
func (_Events *EventsFilterer) WatchPublicGoodRewardDistributed(opts *bind.WatchOpts, sink chan<- *EventsPublicGoodRewardDistributed, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "PublicGoodRewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsPublicGoodRewardDistributed)
				if err := _Events.contract.UnpackLog(event, "PublicGoodRewardDistributed", log); err != nil {
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
func (_Events *EventsFilterer) ParsePublicGoodRewardDistributed(log types.Log) (*EventsPublicGoodRewardDistributed, error) {
	event := new(EventsPublicGoodRewardDistributed)
	if err := _Events.contract.UnpackLog(event, "PublicGoodRewardDistributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsPublicPoolTaxRateBasisPointsSetIterator is returned from FilterPublicPoolTaxRateBasisPointsSet and is used to iterate over the raw logs and unpacked data for PublicPoolTaxRateBasisPointsSet events raised by the Events contract.
type EventsPublicPoolTaxRateBasisPointsSetIterator struct {
	Event *EventsPublicPoolTaxRateBasisPointsSet // Event containing the contract specifics and raw log

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
func (it *EventsPublicPoolTaxRateBasisPointsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsPublicPoolTaxRateBasisPointsSet)
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
		it.Event = new(EventsPublicPoolTaxRateBasisPointsSet)
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
func (it *EventsPublicPoolTaxRateBasisPointsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsPublicPoolTaxRateBasisPointsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsPublicPoolTaxRateBasisPointsSet represents a PublicPoolTaxRateBasisPointsSet event raised by the Events contract.
type EventsPublicPoolTaxRateBasisPointsSet struct {
	TaxRateBasisPoints uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterPublicPoolTaxRateBasisPointsSet is a free log retrieval operation binding the contract event 0x948cf2302b029d76db2ac06e4ef2625e6c687335de349317468f47942a44e8b0.
//
// Solidity: event PublicPoolTaxRateBasisPointsSet(uint64 indexed taxRateBasisPoints)
func (_Events *EventsFilterer) FilterPublicPoolTaxRateBasisPointsSet(opts *bind.FilterOpts, taxRateBasisPoints []uint64) (*EventsPublicPoolTaxRateBasisPointsSetIterator, error) {

	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "PublicPoolTaxRateBasisPointsSet", taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return &EventsPublicPoolTaxRateBasisPointsSetIterator{contract: _Events.contract, event: "PublicPoolTaxRateBasisPointsSet", logs: logs, sub: sub}, nil
}

// WatchPublicPoolTaxRateBasisPointsSet is a free log subscription operation binding the contract event 0x948cf2302b029d76db2ac06e4ef2625e6c687335de349317468f47942a44e8b0.
//
// Solidity: event PublicPoolTaxRateBasisPointsSet(uint64 indexed taxRateBasisPoints)
func (_Events *EventsFilterer) WatchPublicPoolTaxRateBasisPointsSet(opts *bind.WatchOpts, sink chan<- *EventsPublicPoolTaxRateBasisPointsSet, taxRateBasisPoints []uint64) (event.Subscription, error) {

	var taxRateBasisPointsRule []interface{}
	for _, taxRateBasisPointsItem := range taxRateBasisPoints {
		taxRateBasisPointsRule = append(taxRateBasisPointsRule, taxRateBasisPointsItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "PublicPoolTaxRateBasisPointsSet", taxRateBasisPointsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsPublicPoolTaxRateBasisPointsSet)
				if err := _Events.contract.UnpackLog(event, "PublicPoolTaxRateBasisPointsSet", log); err != nil {
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
func (_Events *EventsFilterer) ParsePublicPoolTaxRateBasisPointsSet(log types.Log) (*EventsPublicPoolTaxRateBasisPointsSet, error) {
	event := new(EventsPublicPoolTaxRateBasisPointsSet)
	if err := _Events.contract.UnpackLog(event, "PublicPoolTaxRateBasisPointsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsRewardDistributedIterator is returned from FilterRewardDistributed and is used to iterate over the raw logs and unpacked data for RewardDistributed events raised by the Events contract.
type EventsRewardDistributedIterator struct {
	Event *EventsRewardDistributed // Event containing the contract specifics and raw log

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
func (it *EventsRewardDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsRewardDistributed)
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
		it.Event = new(EventsRewardDistributed)
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
func (it *EventsRewardDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsRewardDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsRewardDistributed represents a RewardDistributed event raised by the Events contract.
type EventsRewardDistributed struct {
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
func (_Events *EventsFilterer) FilterRewardDistributed(opts *bind.FilterOpts, epoch []*big.Int) (*EventsRewardDistributedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "RewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return &EventsRewardDistributedIterator{contract: _Events.contract, event: "RewardDistributed", logs: logs, sub: sub}, nil
}

// WatchRewardDistributed is a free log subscription operation binding the contract event 0x8ea79f19e90b084c2009d3490a097547d8bbb315a883b9efec0996502c1dd7ae.
//
// Solidity: event RewardDistributed(uint256 indexed epoch, uint256 startTimestamp, uint256 endTimestamp, address[] nodeAddrs, uint256[] operationRewards, uint256[] stakingRewards, uint256[] taxCollected, uint256[] requestCounts)
func (_Events *EventsFilterer) WatchRewardDistributed(opts *bind.WatchOpts, sink chan<- *EventsRewardDistributed, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "RewardDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsRewardDistributed)
				if err := _Events.contract.UnpackLog(event, "RewardDistributed", log); err != nil {
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
func (_Events *EventsFilterer) ParseRewardDistributed(log types.Log) (*EventsRewardDistributed, error) {
	event := new(EventsRewardDistributed)
	if err := _Events.contract.UnpackLog(event, "RewardDistributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsSlashCommittedIterator is returned from FilterSlashCommitted and is used to iterate over the raw logs and unpacked data for SlashCommitted events raised by the Events contract.
type EventsSlashCommittedIterator struct {
	Event *EventsSlashCommitted // Event containing the contract specifics and raw log

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
func (it *EventsSlashCommittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsSlashCommitted)
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
		it.Event = new(EventsSlashCommitted)
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
func (it *EventsSlashCommittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsSlashCommittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsSlashCommitted represents a SlashCommitted event raised by the Events contract.
type EventsSlashCommitted struct {
	NodeAddr common.Address
	Epoch    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSlashCommitted is a free log retrieval operation binding the contract event 0x9520d039b6a33a0b9fb3d5c2c5c505f1d07dd471d64492553a5b10ffe05e1f9c.
//
// Solidity: event SlashCommitted(address nodeAddr, uint256 epoch)
func (_Events *EventsFilterer) FilterSlashCommitted(opts *bind.FilterOpts) (*EventsSlashCommittedIterator, error) {

	logs, sub, err := _Events.contract.FilterLogs(opts, "SlashCommitted")
	if err != nil {
		return nil, err
	}
	return &EventsSlashCommittedIterator{contract: _Events.contract, event: "SlashCommitted", logs: logs, sub: sub}, nil
}

// WatchSlashCommitted is a free log subscription operation binding the contract event 0x9520d039b6a33a0b9fb3d5c2c5c505f1d07dd471d64492553a5b10ffe05e1f9c.
//
// Solidity: event SlashCommitted(address nodeAddr, uint256 epoch)
func (_Events *EventsFilterer) WatchSlashCommitted(opts *bind.WatchOpts, sink chan<- *EventsSlashCommitted) (event.Subscription, error) {

	logs, sub, err := _Events.contract.WatchLogs(opts, "SlashCommitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsSlashCommitted)
				if err := _Events.contract.UnpackLog(event, "SlashCommitted", log); err != nil {
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

// ParseSlashCommitted is a log parse operation binding the contract event 0x9520d039b6a33a0b9fb3d5c2c5c505f1d07dd471d64492553a5b10ffe05e1f9c.
//
// Solidity: event SlashCommitted(address nodeAddr, uint256 epoch)
func (_Events *EventsFilterer) ParseSlashCommitted(log types.Log) (*EventsSlashCommitted, error) {
	event := new(EventsSlashCommitted)
	if err := _Events.contract.UnpackLog(event, "SlashCommitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsSlashRecordedIterator is returned from FilterSlashRecorded and is used to iterate over the raw logs and unpacked data for SlashRecorded events raised by the Events contract.
type EventsSlashRecordedIterator struct {
	Event *EventsSlashRecorded // Event containing the contract specifics and raw log

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
func (it *EventsSlashRecordedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsSlashRecorded)
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
		it.Event = new(EventsSlashRecorded)
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
func (it *EventsSlashRecordedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsSlashRecordedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsSlashRecorded represents a SlashRecorded event raised by the Events contract.
type EventsSlashRecorded struct {
	NodeAddr             common.Address
	Epoch                *big.Int
	SlashedOperationPool *big.Int
	SlashedStakingPool   *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSlashRecorded is a free log retrieval operation binding the contract event 0x364424c102a33f9705e7393b69481cfcc8bb4f448148681f2eca17247f017f7d.
//
// Solidity: event SlashRecorded(address indexed nodeAddr, uint256 indexed epoch, uint256 slashedOperationPool, uint256 slashedStakingPool)
func (_Events *EventsFilterer) FilterSlashRecorded(opts *bind.FilterOpts, nodeAddr []common.Address, epoch []*big.Int) (*EventsSlashRecordedIterator, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "SlashRecorded", nodeAddrRule, epochRule)
	if err != nil {
		return nil, err
	}
	return &EventsSlashRecordedIterator{contract: _Events.contract, event: "SlashRecorded", logs: logs, sub: sub}, nil
}

// WatchSlashRecorded is a free log subscription operation binding the contract event 0x364424c102a33f9705e7393b69481cfcc8bb4f448148681f2eca17247f017f7d.
//
// Solidity: event SlashRecorded(address indexed nodeAddr, uint256 indexed epoch, uint256 slashedOperationPool, uint256 slashedStakingPool)
func (_Events *EventsFilterer) WatchSlashRecorded(opts *bind.WatchOpts, sink chan<- *EventsSlashRecorded, nodeAddr []common.Address, epoch []*big.Int) (event.Subscription, error) {

	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "SlashRecorded", nodeAddrRule, epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsSlashRecorded)
				if err := _Events.contract.UnpackLog(event, "SlashRecorded", log); err != nil {
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

// ParseSlashRecorded is a log parse operation binding the contract event 0x364424c102a33f9705e7393b69481cfcc8bb4f448148681f2eca17247f017f7d.
//
// Solidity: event SlashRecorded(address indexed nodeAddr, uint256 indexed epoch, uint256 slashedOperationPool, uint256 slashedStakingPool)
func (_Events *EventsFilterer) ParseSlashRecorded(log types.Log) (*EventsSlashRecorded, error) {
	event := new(EventsSlashRecorded)
	if err := _Events.contract.UnpackLog(event, "SlashRecorded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsSlashRevokedIterator is returned from FilterSlashRevoked and is used to iterate over the raw logs and unpacked data for SlashRevoked events raised by the Events contract.
type EventsSlashRevokedIterator struct {
	Event *EventsSlashRevoked // Event containing the contract specifics and raw log

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
func (it *EventsSlashRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsSlashRevoked)
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
		it.Event = new(EventsSlashRevoked)
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
func (it *EventsSlashRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsSlashRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsSlashRevoked represents a SlashRevoked event raised by the Events contract.
type EventsSlashRevoked struct {
	NodeAddr common.Address
	Epoch    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSlashRevoked is a free log retrieval operation binding the contract event 0x15e7cd41a21eea06cee0aeb26dbe9c430c42463fc7f8e0124688bcde5bf04a31.
//
// Solidity: event SlashRevoked(address nodeAddr, uint256 epoch)
func (_Events *EventsFilterer) FilterSlashRevoked(opts *bind.FilterOpts) (*EventsSlashRevokedIterator, error) {

	logs, sub, err := _Events.contract.FilterLogs(opts, "SlashRevoked")
	if err != nil {
		return nil, err
	}
	return &EventsSlashRevokedIterator{contract: _Events.contract, event: "SlashRevoked", logs: logs, sub: sub}, nil
}

// WatchSlashRevoked is a free log subscription operation binding the contract event 0x15e7cd41a21eea06cee0aeb26dbe9c430c42463fc7f8e0124688bcde5bf04a31.
//
// Solidity: event SlashRevoked(address nodeAddr, uint256 epoch)
func (_Events *EventsFilterer) WatchSlashRevoked(opts *bind.WatchOpts, sink chan<- *EventsSlashRevoked) (event.Subscription, error) {

	logs, sub, err := _Events.contract.WatchLogs(opts, "SlashRevoked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsSlashRevoked)
				if err := _Events.contract.UnpackLog(event, "SlashRevoked", log); err != nil {
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

// ParseSlashRevoked is a log parse operation binding the contract event 0x15e7cd41a21eea06cee0aeb26dbe9c430c42463fc7f8e0124688bcde5bf04a31.
//
// Solidity: event SlashRevoked(address nodeAddr, uint256 epoch)
func (_Events *EventsFilterer) ParseSlashRevoked(log types.Log) (*EventsSlashRevoked, error) {
	event := new(EventsSlashRevoked)
	if err := _Events.contract.UnpackLog(event, "SlashRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Events contract.
type EventsStakedIterator struct {
	Event *EventsStaked // Event containing the contract specifics and raw log

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
func (it *EventsStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsStaked)
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
		it.Event = new(EventsStaked)
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
func (it *EventsStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsStaked represents a Staked event raised by the Events contract.
type EventsStaked struct {
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
func (_Events *EventsFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address, nodeAddr []common.Address, amount []*big.Int) (*EventsStakedIterator, error) {

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

	logs, sub, err := _Events.contract.FilterLogs(opts, "Staked", userRule, nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &EventsStakedIterator{contract: _Events.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0xad3fa07f4195b47e64892eb944ecbfc253384053c119852bb2bcae484c2fcb69.
//
// Solidity: event Staked(address indexed user, address indexed nodeAddr, uint256 indexed amount, uint256 startTokenId, uint256 endTokenId)
func (_Events *EventsFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *EventsStaked, user []common.Address, nodeAddr []common.Address, amount []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Events.contract.WatchLogs(opts, "Staked", userRule, nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsStaked)
				if err := _Events.contract.UnpackLog(event, "Staked", log); err != nil {
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
func (_Events *EventsFilterer) ParseStaked(log types.Log) (*EventsStaked, error) {
	event := new(EventsStaked)
	if err := _Events.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsUnstakeClaimedIterator is returned from FilterUnstakeClaimed and is used to iterate over the raw logs and unpacked data for UnstakeClaimed events raised by the Events contract.
type EventsUnstakeClaimedIterator struct {
	Event *EventsUnstakeClaimed // Event containing the contract specifics and raw log

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
func (it *EventsUnstakeClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsUnstakeClaimed)
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
		it.Event = new(EventsUnstakeClaimed)
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
func (it *EventsUnstakeClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsUnstakeClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsUnstakeClaimed represents a UnstakeClaimed event raised by the Events contract.
type EventsUnstakeClaimed struct {
	RequestId     *big.Int
	NodeAddr      common.Address
	User          common.Address
	UnstakeAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUnstakeClaimed is a free log retrieval operation binding the contract event 0x2769ece66eadb650afd8c08c7a8772e39381dddd7230f9e039669e631044d47c.
//
// Solidity: event UnstakeClaimed(uint256 indexed requestId, address indexed nodeAddr, address indexed user, uint256 unstakeAmount)
func (_Events *EventsFilterer) FilterUnstakeClaimed(opts *bind.FilterOpts, requestId []*big.Int, nodeAddr []common.Address, user []common.Address) (*EventsUnstakeClaimedIterator, error) {

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

	logs, sub, err := _Events.contract.FilterLogs(opts, "UnstakeClaimed", requestIdRule, nodeAddrRule, userRule)
	if err != nil {
		return nil, err
	}
	return &EventsUnstakeClaimedIterator{contract: _Events.contract, event: "UnstakeClaimed", logs: logs, sub: sub}, nil
}

// WatchUnstakeClaimed is a free log subscription operation binding the contract event 0x2769ece66eadb650afd8c08c7a8772e39381dddd7230f9e039669e631044d47c.
//
// Solidity: event UnstakeClaimed(uint256 indexed requestId, address indexed nodeAddr, address indexed user, uint256 unstakeAmount)
func (_Events *EventsFilterer) WatchUnstakeClaimed(opts *bind.WatchOpts, sink chan<- *EventsUnstakeClaimed, requestId []*big.Int, nodeAddr []common.Address, user []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Events.contract.WatchLogs(opts, "UnstakeClaimed", requestIdRule, nodeAddrRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsUnstakeClaimed)
				if err := _Events.contract.UnpackLog(event, "UnstakeClaimed", log); err != nil {
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
func (_Events *EventsFilterer) ParseUnstakeClaimed(log types.Log) (*EventsUnstakeClaimed, error) {
	event := new(EventsUnstakeClaimed)
	if err := _Events.contract.UnpackLog(event, "UnstakeClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsUnstakeRequestedIterator is returned from FilterUnstakeRequested and is used to iterate over the raw logs and unpacked data for UnstakeRequested events raised by the Events contract.
type EventsUnstakeRequestedIterator struct {
	Event *EventsUnstakeRequested // Event containing the contract specifics and raw log

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
func (it *EventsUnstakeRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsUnstakeRequested)
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
		it.Event = new(EventsUnstakeRequested)
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
func (it *EventsUnstakeRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsUnstakeRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsUnstakeRequested represents a UnstakeRequested event raised by the Events contract.
type EventsUnstakeRequested struct {
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
func (_Events *EventsFilterer) FilterUnstakeRequested(opts *bind.FilterOpts, user []common.Address, nodeAddr []common.Address, requestId []*big.Int) (*EventsUnstakeRequestedIterator, error) {

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

	logs, sub, err := _Events.contract.FilterLogs(opts, "UnstakeRequested", userRule, nodeAddrRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &EventsUnstakeRequestedIterator{contract: _Events.contract, event: "UnstakeRequested", logs: logs, sub: sub}, nil
}

// WatchUnstakeRequested is a free log subscription operation binding the contract event 0x2808f92d5a0fada467cbe4e766f62f323e78271a7471058a87ef63a9e3e4c5c5.
//
// Solidity: event UnstakeRequested(address indexed user, address indexed nodeAddr, uint256 indexed requestId, uint256 unstakeAmount, uint256[] chipsIds)
func (_Events *EventsFilterer) WatchUnstakeRequested(opts *bind.WatchOpts, sink chan<- *EventsUnstakeRequested, user []common.Address, nodeAddr []common.Address, requestId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Events.contract.WatchLogs(opts, "UnstakeRequested", userRule, nodeAddrRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsUnstakeRequested)
				if err := _Events.contract.UnpackLog(event, "UnstakeRequested", log); err != nil {
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
func (_Events *EventsFilterer) ParseUnstakeRequested(log types.Log) (*EventsUnstakeRequested, error) {
	event := new(EventsUnstakeRequested)
	if err := _Events.contract.UnpackLog(event, "UnstakeRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsWithdrawRequestedIterator is returned from FilterWithdrawRequested and is used to iterate over the raw logs and unpacked data for WithdrawRequested events raised by the Events contract.
type EventsWithdrawRequestedIterator struct {
	Event *EventsWithdrawRequested // Event containing the contract specifics and raw log

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
func (it *EventsWithdrawRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsWithdrawRequested)
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
		it.Event = new(EventsWithdrawRequested)
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
func (it *EventsWithdrawRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsWithdrawRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsWithdrawRequested represents a WithdrawRequested event raised by the Events contract.
type EventsWithdrawRequested struct {
	NodeAddr  common.Address
	Amount    *big.Int
	RequestId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawRequested is a free log retrieval operation binding the contract event 0xd72eb5d043f24a0168ae744d5c44f9596fd673a26bf74d9646bff4b844882d14.
//
// Solidity: event WithdrawRequested(address indexed nodeAddr, uint256 indexed amount, uint256 indexed requestId)
func (_Events *EventsFilterer) FilterWithdrawRequested(opts *bind.FilterOpts, nodeAddr []common.Address, amount []*big.Int, requestId []*big.Int) (*EventsWithdrawRequestedIterator, error) {

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

	logs, sub, err := _Events.contract.FilterLogs(opts, "WithdrawRequested", nodeAddrRule, amountRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &EventsWithdrawRequestedIterator{contract: _Events.contract, event: "WithdrawRequested", logs: logs, sub: sub}, nil
}

// WatchWithdrawRequested is a free log subscription operation binding the contract event 0xd72eb5d043f24a0168ae744d5c44f9596fd673a26bf74d9646bff4b844882d14.
//
// Solidity: event WithdrawRequested(address indexed nodeAddr, uint256 indexed amount, uint256 indexed requestId)
func (_Events *EventsFilterer) WatchWithdrawRequested(opts *bind.WatchOpts, sink chan<- *EventsWithdrawRequested, nodeAddr []common.Address, amount []*big.Int, requestId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Events.contract.WatchLogs(opts, "WithdrawRequested", nodeAddrRule, amountRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsWithdrawRequested)
				if err := _Events.contract.UnpackLog(event, "WithdrawRequested", log); err != nil {
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
func (_Events *EventsFilterer) ParseWithdrawRequested(log types.Log) (*EventsWithdrawRequested, error) {
	event := new(EventsWithdrawRequested)
	if err := _Events.contract.UnpackLog(event, "WithdrawRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EventsWithdrawalClaimedIterator is returned from FilterWithdrawalClaimed and is used to iterate over the raw logs and unpacked data for WithdrawalClaimed events raised by the Events contract.
type EventsWithdrawalClaimedIterator struct {
	Event *EventsWithdrawalClaimed // Event containing the contract specifics and raw log

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
func (it *EventsWithdrawalClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EventsWithdrawalClaimed)
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
		it.Event = new(EventsWithdrawalClaimed)
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
func (it *EventsWithdrawalClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EventsWithdrawalClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EventsWithdrawalClaimed represents a WithdrawalClaimed event raised by the Events contract.
type EventsWithdrawalClaimed struct {
	RequestId *big.Int
	NodeAddr  common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalClaimed is a free log retrieval operation binding the contract event 0x8adb7a84b2998a8d11cd9284395f95d5a99f160be785ae79998c654979bd3d9a.
//
// Solidity: event WithdrawalClaimed(uint256 indexed requestId, address indexed nodeAddr, uint256 indexed amount)
func (_Events *EventsFilterer) FilterWithdrawalClaimed(opts *bind.FilterOpts, requestId []*big.Int, nodeAddr []common.Address, amount []*big.Int) (*EventsWithdrawalClaimedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Events.contract.FilterLogs(opts, "WithdrawalClaimed", requestIdRule, nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &EventsWithdrawalClaimedIterator{contract: _Events.contract, event: "WithdrawalClaimed", logs: logs, sub: sub}, nil
}

// WatchWithdrawalClaimed is a free log subscription operation binding the contract event 0x8adb7a84b2998a8d11cd9284395f95d5a99f160be785ae79998c654979bd3d9a.
//
// Solidity: event WithdrawalClaimed(uint256 indexed requestId, address indexed nodeAddr, uint256 indexed amount)
func (_Events *EventsFilterer) WatchWithdrawalClaimed(opts *bind.WatchOpts, sink chan<- *EventsWithdrawalClaimed, requestId []*big.Int, nodeAddr []common.Address, amount []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var nodeAddrRule []interface{}
	for _, nodeAddrItem := range nodeAddr {
		nodeAddrRule = append(nodeAddrRule, nodeAddrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Events.contract.WatchLogs(opts, "WithdrawalClaimed", requestIdRule, nodeAddrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EventsWithdrawalClaimed)
				if err := _Events.contract.UnpackLog(event, "WithdrawalClaimed", log); err != nil {
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

// ParseWithdrawalClaimed is a log parse operation binding the contract event 0x8adb7a84b2998a8d11cd9284395f95d5a99f160be785ae79998c654979bd3d9a.
//
// Solidity: event WithdrawalClaimed(uint256 indexed requestId, address indexed nodeAddr, uint256 indexed amount)
func (_Events *EventsFilterer) ParseWithdrawalClaimed(log types.Log) (*EventsWithdrawalClaimed, error) {
	event := new(EventsWithdrawalClaimed)
	if err := _Events.contract.UnpackLog(event, "WithdrawalClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
