// This file is an automatically generated Go binding. Do not modify as any
// change will likely be lost upon the next re-generation!

package main

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// LuxUni_EETABI is the input ABI used to generate the binding from.
const LuxUni_EETABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"receivers\",\"outputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"creationDate\",\"type\":\"uint256\"},{\"name\":\"numTransactions\",\"type\":\"uint256\"},{\"name\":\"numPoints\",\"type\":\"uint256\"},{\"name\":\"numPriorities\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"},{\"name\":\"isReceiver\",\"type\":\"bool\"}],\"name\":\"getTotPoints\",\"outputs\":[{\"name\":\"_total\",\"type\":\"int256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"},{\"name\":\"isReceiver\",\"type\":\"bool\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"addPoints\",\"outputs\":[{\"name\":\"error\",\"type\":\"int256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"senderList\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"receiverList\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_senderAddr\",\"type\":\"address\"},{\"name\":\"_receiverAddr\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_priorityID\",\"type\":\"int256\"},{\"name\":\"_certificate\",\"type\":\"bytes\"}],\"name\":\"makeDonation\",\"outputs\":[{\"name\":\"_ret\",\"type\":\"int256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"Eethiq\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_senderAddr\",\"type\":\"address\"},{\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"newSender\",\"outputs\":[{\"name\":\"_ret\",\"type\":\"int256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"},{\"name\":\"isReceiver\",\"type\":\"bool\"},{\"name\":\"isTotalScan\",\"type\":\"bool\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"_balance\",\"type\":\"int256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"errLog\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numErrLog\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numTransactions\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"senders\",\"outputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"creationDate\",\"type\":\"uint256\"},{\"name\":\"numTransactions\",\"type\":\"uint256\"},{\"name\":\"numPoints\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transactions\",\"outputs\":[{\"name\":\"senderAddr\",\"type\":\"address\"},{\"name\":\"receiverAddr\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"priorityQty\",\"type\":\"uint256\"},{\"name\":\"priorityID\",\"type\":\"uint256\"},{\"name\":\"transDate\",\"type\":\"uint256\"},{\"name\":\"certificate\",\"type\":\"bytes\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"priorities\",\"outputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"score\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"},{\"name\":\"isReceiver\",\"type\":\"bool\"},{\"name\":\"isFullScan\",\"type\":\"bool\"}],\"name\":\"getTotDonations\",\"outputs\":[{\"name\":\"_total\",\"type\":\"int256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isAdmin\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numReceivers\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_receiverAddr\",\"type\":\"address\"},{\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"newReceiver\",\"outputs\":[{\"name\":\"_ret\",\"type\":\"int256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numSenders\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"}]"

// LuxUni_EET is an auto generated Go binding around an Ethereum contract.
type LuxUni_EET struct {
	LuxUni_EETCaller     // Read-only binding to the contract
	LuxUni_EETTransactor // Write-only binding to the contract
}

// LuxUni_EETCaller is an auto generated read-only Go binding around an Ethereum contract.
type LuxUni_EETCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LuxUni_EETTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LuxUni_EETTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LuxUni_EETSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LuxUni_EETSession struct {
	Contract     *LuxUni_EET       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LuxUni_EETCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LuxUni_EETCallerSession struct {
	Contract *LuxUni_EETCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// LuxUni_EETTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LuxUni_EETTransactorSession struct {
	Contract     *LuxUni_EETTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// LuxUni_EETRaw is an auto generated low-level Go binding around an Ethereum contract.
type LuxUni_EETRaw struct {
	Contract *LuxUni_EET // Generic contract binding to access the raw methods on
}

// LuxUni_EETCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LuxUni_EETCallerRaw struct {
	Contract *LuxUni_EETCaller // Generic read-only contract binding to access the raw methods on
}

// LuxUni_EETTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LuxUni_EETTransactorRaw struct {
	Contract *LuxUni_EETTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLuxUni_EET creates a new instance of LuxUni_EET, bound to a specific deployed contract.
func NewLuxUni_EET(address common.Address, backend bind.ContractBackend) (*LuxUni_EET, error) {
	contract, err := bindLuxUni_EET(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LuxUni_EET{LuxUni_EETCaller: LuxUni_EETCaller{contract: contract}, LuxUni_EETTransactor: LuxUni_EETTransactor{contract: contract}}, nil
}

// NewLuxUni_EETCaller creates a new read-only instance of LuxUni_EET, bound to a specific deployed contract.
func NewLuxUni_EETCaller(address common.Address, caller bind.ContractCaller) (*LuxUni_EETCaller, error) {
	contract, err := bindLuxUni_EET(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &LuxUni_EETCaller{contract: contract}, nil
}

// NewLuxUni_EETTransactor creates a new write-only instance of LuxUni_EET, bound to a specific deployed contract.
func NewLuxUni_EETTransactor(address common.Address, transactor bind.ContractTransactor) (*LuxUni_EETTransactor, error) {
	contract, err := bindLuxUni_EET(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &LuxUni_EETTransactor{contract: contract}, nil
}

// bindLuxUni_EET binds a generic wrapper to an already deployed contract.
func bindLuxUni_EET(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LuxUni_EETABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LuxUni_EET *LuxUni_EETRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LuxUni_EET.Contract.LuxUni_EETCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LuxUni_EET *LuxUni_EETRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.LuxUni_EETTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LuxUni_EET *LuxUni_EETRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.LuxUni_EETTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LuxUni_EET *LuxUni_EETCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LuxUni_EET.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LuxUni_EET *LuxUni_EETTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LuxUni_EET *LuxUni_EETTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.contract.Transact(opts, method, params...)
}

// ErrLog is a free data retrieval call binding the contract method 0x84e9a230.
//
// Solidity: function errLog( uint256) constant returns(string)
func (_LuxUni_EET *LuxUni_EETCaller) ErrLog(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "errLog", arg0)
	return *ret0, err
}

// ErrLog is a free data retrieval call binding the contract method 0x84e9a230.
//
// Solidity: function errLog( uint256) constant returns(string)
func (_LuxUni_EET *LuxUni_EETSession) ErrLog(arg0 *big.Int) (string, error) {
	return _LuxUni_EET.Contract.ErrLog(&_LuxUni_EET.CallOpts, arg0)
}

// ErrLog is a free data retrieval call binding the contract method 0x84e9a230.
//
// Solidity: function errLog( uint256) constant returns(string)
func (_LuxUni_EET *LuxUni_EETCallerSession) ErrLog(arg0 *big.Int) (string, error) {
	return _LuxUni_EET.Contract.ErrLog(&_LuxUni_EET.CallOpts, arg0)
}

// GetBalance is a free data retrieval call binding the contract method 0x84d59a0e.
//
// Solidity: function getBalance(_addr address, isReceiver bool, isTotalScan bool) constant returns(_balance int256)
func (_LuxUni_EET *LuxUni_EETCaller) GetBalance(opts *bind.CallOpts, _addr common.Address, isReceiver bool, isTotalScan bool) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "getBalance", _addr, isReceiver, isTotalScan)
	return *ret0, err
}

// GetBalance is a free data retrieval call binding the contract method 0x84d59a0e.
//
// Solidity: function getBalance(_addr address, isReceiver bool, isTotalScan bool) constant returns(_balance int256)
func (_LuxUni_EET *LuxUni_EETSession) GetBalance(_addr common.Address, isReceiver bool, isTotalScan bool) (*big.Int, error) {
	return _LuxUni_EET.Contract.GetBalance(&_LuxUni_EET.CallOpts, _addr, isReceiver, isTotalScan)
}

// GetBalance is a free data retrieval call binding the contract method 0x84d59a0e.
//
// Solidity: function getBalance(_addr address, isReceiver bool, isTotalScan bool) constant returns(_balance int256)
func (_LuxUni_EET *LuxUni_EETCallerSession) GetBalance(_addr common.Address, isReceiver bool, isTotalScan bool) (*big.Int, error) {
	return _LuxUni_EET.Contract.GetBalance(&_LuxUni_EET.CallOpts, _addr, isReceiver, isTotalScan)
}

// GetTotDonations is a free data retrieval call binding the contract method 0xb6cdbaf5.
//
// Solidity: function getTotDonations(_addr address, isReceiver bool, isFullScan bool) constant returns(_total int256)
func (_LuxUni_EET *LuxUni_EETCaller) GetTotDonations(opts *bind.CallOpts, _addr common.Address, isReceiver bool, isFullScan bool) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "getTotDonations", _addr, isReceiver, isFullScan)
	return *ret0, err
}

// GetTotDonations is a free data retrieval call binding the contract method 0xb6cdbaf5.
//
// Solidity: function getTotDonations(_addr address, isReceiver bool, isFullScan bool) constant returns(_total int256)
func (_LuxUni_EET *LuxUni_EETSession) GetTotDonations(_addr common.Address, isReceiver bool, isFullScan bool) (*big.Int, error) {
	return _LuxUni_EET.Contract.GetTotDonations(&_LuxUni_EET.CallOpts, _addr, isReceiver, isFullScan)
}

// GetTotDonations is a free data retrieval call binding the contract method 0xb6cdbaf5.
//
// Solidity: function getTotDonations(_addr address, isReceiver bool, isFullScan bool) constant returns(_total int256)
func (_LuxUni_EET *LuxUni_EETCallerSession) GetTotDonations(_addr common.Address, isReceiver bool, isFullScan bool) (*big.Int, error) {
	return _LuxUni_EET.Contract.GetTotDonations(&_LuxUni_EET.CallOpts, _addr, isReceiver, isFullScan)
}

// GetTotPoints is a free data retrieval call binding the contract method 0x13f38540.
//
// Solidity: function getTotPoints(_addr address, isReceiver bool) constant returns(_total int256)
func (_LuxUni_EET *LuxUni_EETCaller) GetTotPoints(opts *bind.CallOpts, _addr common.Address, isReceiver bool) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "getTotPoints", _addr, isReceiver)
	return *ret0, err
}

// GetTotPoints is a free data retrieval call binding the contract method 0x13f38540.
//
// Solidity: function getTotPoints(_addr address, isReceiver bool) constant returns(_total int256)
func (_LuxUni_EET *LuxUni_EETSession) GetTotPoints(_addr common.Address, isReceiver bool) (*big.Int, error) {
	return _LuxUni_EET.Contract.GetTotPoints(&_LuxUni_EET.CallOpts, _addr, isReceiver)
}

// GetTotPoints is a free data retrieval call binding the contract method 0x13f38540.
//
// Solidity: function getTotPoints(_addr address, isReceiver bool) constant returns(_total int256)
func (_LuxUni_EET *LuxUni_EETCallerSession) GetTotPoints(_addr common.Address, isReceiver bool) (*big.Int, error) {
	return _LuxUni_EET.Contract.GetTotPoints(&_LuxUni_EET.CallOpts, _addr, isReceiver)
}

// IsAdmin is a free data retrieval call binding the contract method 0xb6db75a0.
//
// Solidity: function isAdmin() constant returns(bool)
func (_LuxUni_EET *LuxUni_EETCaller) IsAdmin(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "isAdmin")
	return *ret0, err
}

// IsAdmin is a free data retrieval call binding the contract method 0xb6db75a0.
//
// Solidity: function isAdmin() constant returns(bool)
func (_LuxUni_EET *LuxUni_EETSession) IsAdmin() (bool, error) {
	return _LuxUni_EET.Contract.IsAdmin(&_LuxUni_EET.CallOpts)
}

// IsAdmin is a free data retrieval call binding the contract method 0xb6db75a0.
//
// Solidity: function isAdmin() constant returns(bool)
func (_LuxUni_EET *LuxUni_EETCallerSession) IsAdmin() (bool, error) {
	return _LuxUni_EET.Contract.IsAdmin(&_LuxUni_EET.CallOpts)
}

// NumErrLog is a free data retrieval call binding the contract method 0x8b614a93.
//
// Solidity: function numErrLog() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETCaller) NumErrLog(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "numErrLog")
	return *ret0, err
}

// NumErrLog is a free data retrieval call binding the contract method 0x8b614a93.
//
// Solidity: function numErrLog() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETSession) NumErrLog() (*big.Int, error) {
	return _LuxUni_EET.Contract.NumErrLog(&_LuxUni_EET.CallOpts)
}

// NumErrLog is a free data retrieval call binding the contract method 0x8b614a93.
//
// Solidity: function numErrLog() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETCallerSession) NumErrLog() (*big.Int, error) {
	return _LuxUni_EET.Contract.NumErrLog(&_LuxUni_EET.CallOpts)
}

// NumReceivers is a free data retrieval call binding the contract method 0xcbc06484.
//
// Solidity: function numReceivers() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETCaller) NumReceivers(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "numReceivers")
	return *ret0, err
}

// NumReceivers is a free data retrieval call binding the contract method 0xcbc06484.
//
// Solidity: function numReceivers() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETSession) NumReceivers() (*big.Int, error) {
	return _LuxUni_EET.Contract.NumReceivers(&_LuxUni_EET.CallOpts)
}

// NumReceivers is a free data retrieval call binding the contract method 0xcbc06484.
//
// Solidity: function numReceivers() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETCallerSession) NumReceivers() (*big.Int, error) {
	return _LuxUni_EET.Contract.NumReceivers(&_LuxUni_EET.CallOpts)
}

// NumSenders is a free data retrieval call binding the contract method 0xda40c30e.
//
// Solidity: function numSenders() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETCaller) NumSenders(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "numSenders")
	return *ret0, err
}

// NumSenders is a free data retrieval call binding the contract method 0xda40c30e.
//
// Solidity: function numSenders() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETSession) NumSenders() (*big.Int, error) {
	return _LuxUni_EET.Contract.NumSenders(&_LuxUni_EET.CallOpts)
}

// NumSenders is a free data retrieval call binding the contract method 0xda40c30e.
//
// Solidity: function numSenders() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETCallerSession) NumSenders() (*big.Int, error) {
	return _LuxUni_EET.Contract.NumSenders(&_LuxUni_EET.CallOpts)
}

// NumTransactions is a free data retrieval call binding the contract method 0x90b4cc72.
//
// Solidity: function numTransactions() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETCaller) NumTransactions(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "numTransactions")
	return *ret0, err
}

// NumTransactions is a free data retrieval call binding the contract method 0x90b4cc72.
//
// Solidity: function numTransactions() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETSession) NumTransactions() (*big.Int, error) {
	return _LuxUni_EET.Contract.NumTransactions(&_LuxUni_EET.CallOpts)
}

// NumTransactions is a free data retrieval call binding the contract method 0x90b4cc72.
//
// Solidity: function numTransactions() constant returns(uint256)
func (_LuxUni_EET *LuxUni_EETCallerSession) NumTransactions() (*big.Int, error) {
	return _LuxUni_EET.Contract.NumTransactions(&_LuxUni_EET.CallOpts)
}

// Priorities is a free data retrieval call binding the contract method 0xb36f923e.
//
// Solidity: function priorities( uint256) constant returns(name string, score uint256)
func (_LuxUni_EET *LuxUni_EETCaller) Priorities(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Name  string
	Score *big.Int
}, error) {
	ret := new(struct {
		Name  string
		Score *big.Int
	})
	out := ret
	err := _LuxUni_EET.contract.Call(opts, out, "priorities", arg0)
	return *ret, err
}

// Priorities is a free data retrieval call binding the contract method 0xb36f923e.
//
// Solidity: function priorities( uint256) constant returns(name string, score uint256)
func (_LuxUni_EET *LuxUni_EETSession) Priorities(arg0 *big.Int) (struct {
	Name  string
	Score *big.Int
}, error) {
	return _LuxUni_EET.Contract.Priorities(&_LuxUni_EET.CallOpts, arg0)
}

// Priorities is a free data retrieval call binding the contract method 0xb36f923e.
//
// Solidity: function priorities( uint256) constant returns(name string, score uint256)
func (_LuxUni_EET *LuxUni_EETCallerSession) Priorities(arg0 *big.Int) (struct {
	Name  string
	Score *big.Int
}, error) {
	return _LuxUni_EET.Contract.Priorities(&_LuxUni_EET.CallOpts, arg0)
}

// ReceiverList is a free data retrieval call binding the contract method 0x5b6f7b1d.
//
// Solidity: function receiverList( uint256) constant returns(address)
func (_LuxUni_EET *LuxUni_EETCaller) ReceiverList(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "receiverList", arg0)
	return *ret0, err
}

// ReceiverList is a free data retrieval call binding the contract method 0x5b6f7b1d.
//
// Solidity: function receiverList( uint256) constant returns(address)
func (_LuxUni_EET *LuxUni_EETSession) ReceiverList(arg0 *big.Int) (common.Address, error) {
	return _LuxUni_EET.Contract.ReceiverList(&_LuxUni_EET.CallOpts, arg0)
}

// ReceiverList is a free data retrieval call binding the contract method 0x5b6f7b1d.
//
// Solidity: function receiverList( uint256) constant returns(address)
func (_LuxUni_EET *LuxUni_EETCallerSession) ReceiverList(arg0 *big.Int) (common.Address, error) {
	return _LuxUni_EET.Contract.ReceiverList(&_LuxUni_EET.CallOpts, arg0)
}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers( address) constant returns(name string, creationDate uint256, numTransactions uint256, numPoints uint256, numPriorities uint256)
func (_LuxUni_EET *LuxUni_EETCaller) Receivers(opts *bind.CallOpts, arg0 common.Address) (struct {
	Name            string
	CreationDate    *big.Int
	NumTransactions *big.Int
	NumPoints       *big.Int
	NumPriorities   *big.Int
}, error) {
	ret := new(struct {
		Name            string
		CreationDate    *big.Int
		NumTransactions *big.Int
		NumPoints       *big.Int
		NumPriorities   *big.Int
	})
	out := ret
	err := _LuxUni_EET.contract.Call(opts, out, "receivers", arg0)
	return *ret, err
}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers( address) constant returns(name string, creationDate uint256, numTransactions uint256, numPoints uint256, numPriorities uint256)
func (_LuxUni_EET *LuxUni_EETSession) Receivers(arg0 common.Address) (struct {
	Name            string
	CreationDate    *big.Int
	NumTransactions *big.Int
	NumPoints       *big.Int
	NumPriorities   *big.Int
}, error) {
	return _LuxUni_EET.Contract.Receivers(&_LuxUni_EET.CallOpts, arg0)
}

// Receivers is a free data retrieval call binding the contract method 0x0cb8150f.
//
// Solidity: function receivers( address) constant returns(name string, creationDate uint256, numTransactions uint256, numPoints uint256, numPriorities uint256)
func (_LuxUni_EET *LuxUni_EETCallerSession) Receivers(arg0 common.Address) (struct {
	Name            string
	CreationDate    *big.Int
	NumTransactions *big.Int
	NumPoints       *big.Int
	NumPriorities   *big.Int
}, error) {
	return _LuxUni_EET.Contract.Receivers(&_LuxUni_EET.CallOpts, arg0)
}

// SenderList is a free data retrieval call binding the contract method 0x492ddfb1.
//
// Solidity: function senderList( uint256) constant returns(address)
func (_LuxUni_EET *LuxUni_EETCaller) SenderList(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_EET.contract.Call(opts, out, "senderList", arg0)
	return *ret0, err
}

// SenderList is a free data retrieval call binding the contract method 0x492ddfb1.
//
// Solidity: function senderList( uint256) constant returns(address)
func (_LuxUni_EET *LuxUni_EETSession) SenderList(arg0 *big.Int) (common.Address, error) {
	return _LuxUni_EET.Contract.SenderList(&_LuxUni_EET.CallOpts, arg0)
}

// SenderList is a free data retrieval call binding the contract method 0x492ddfb1.
//
// Solidity: function senderList( uint256) constant returns(address)
func (_LuxUni_EET *LuxUni_EETCallerSession) SenderList(arg0 *big.Int) (common.Address, error) {
	return _LuxUni_EET.Contract.SenderList(&_LuxUni_EET.CallOpts, arg0)
}

// Senders is a free data retrieval call binding the contract method 0x982fb9d8.
//
// Solidity: function senders( address) constant returns(name string, creationDate uint256, numTransactions uint256, numPoints uint256)
func (_LuxUni_EET *LuxUni_EETCaller) Senders(opts *bind.CallOpts, arg0 common.Address) (struct {
	Name            string
	CreationDate    *big.Int
	NumTransactions *big.Int
	NumPoints       *big.Int
}, error) {
	ret := new(struct {
		Name            string
		CreationDate    *big.Int
		NumTransactions *big.Int
		NumPoints       *big.Int
	})
	out := ret
	err := _LuxUni_EET.contract.Call(opts, out, "senders", arg0)
	return *ret, err
}

// Senders is a free data retrieval call binding the contract method 0x982fb9d8.
//
// Solidity: function senders( address) constant returns(name string, creationDate uint256, numTransactions uint256, numPoints uint256)
func (_LuxUni_EET *LuxUni_EETSession) Senders(arg0 common.Address) (struct {
	Name            string
	CreationDate    *big.Int
	NumTransactions *big.Int
	NumPoints       *big.Int
}, error) {
	return _LuxUni_EET.Contract.Senders(&_LuxUni_EET.CallOpts, arg0)
}

// Senders is a free data retrieval call binding the contract method 0x982fb9d8.
//
// Solidity: function senders( address) constant returns(name string, creationDate uint256, numTransactions uint256, numPoints uint256)
func (_LuxUni_EET *LuxUni_EETCallerSession) Senders(arg0 common.Address) (struct {
	Name            string
	CreationDate    *big.Int
	NumTransactions *big.Int
	NumPoints       *big.Int
}, error) {
	return _LuxUni_EET.Contract.Senders(&_LuxUni_EET.CallOpts, arg0)
}

// Transactions is a free data retrieval call binding the contract method 0x9ace38c2.
//
// Solidity: function transactions( uint256) constant returns(senderAddr address, receiverAddr address, amount uint256, priorityQty uint256, priorityID uint256, transDate uint256, certificate bytes)
func (_LuxUni_EET *LuxUni_EETCaller) Transactions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	SenderAddr   common.Address
	ReceiverAddr common.Address
	Amount       *big.Int
	PriorityQty  *big.Int
	PriorityID   *big.Int
	TransDate    *big.Int
	Certificate  []byte
}, error) {
	ret := new(struct {
		SenderAddr   common.Address
		ReceiverAddr common.Address
		Amount       *big.Int
		PriorityQty  *big.Int
		PriorityID   *big.Int
		TransDate    *big.Int
		Certificate  []byte
	})
	out := ret
	err := _LuxUni_EET.contract.Call(opts, out, "transactions", arg0)
	return *ret, err
}

// Transactions is a free data retrieval call binding the contract method 0x9ace38c2.
//
// Solidity: function transactions( uint256) constant returns(senderAddr address, receiverAddr address, amount uint256, priorityQty uint256, priorityID uint256, transDate uint256, certificate bytes)
func (_LuxUni_EET *LuxUni_EETSession) Transactions(arg0 *big.Int) (struct {
	SenderAddr   common.Address
	ReceiverAddr common.Address
	Amount       *big.Int
	PriorityQty  *big.Int
	PriorityID   *big.Int
	TransDate    *big.Int
	Certificate  []byte
}, error) {
	return _LuxUni_EET.Contract.Transactions(&_LuxUni_EET.CallOpts, arg0)
}

// Transactions is a free data retrieval call binding the contract method 0x9ace38c2.
//
// Solidity: function transactions( uint256) constant returns(senderAddr address, receiverAddr address, amount uint256, priorityQty uint256, priorityID uint256, transDate uint256, certificate bytes)
func (_LuxUni_EET *LuxUni_EETCallerSession) Transactions(arg0 *big.Int) (struct {
	SenderAddr   common.Address
	ReceiverAddr common.Address
	Amount       *big.Int
	PriorityQty  *big.Int
	PriorityID   *big.Int
	TransDate    *big.Int
	Certificate  []byte
}, error) {
	return _LuxUni_EET.Contract.Transactions(&_LuxUni_EET.CallOpts, arg0)
}

// Eethiq is a paid mutator transaction binding the contract method 0x5f319d66.
//
// Solidity: function Eethiq() returns()
func (_LuxUni_EET *LuxUni_EETTransactor) Eethiq(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_EET.contract.Transact(opts, "Eethiq")
}

// Eethiq is a paid mutator transaction binding the contract method 0x5f319d66.
//
// Solidity: function Eethiq() returns()
func (_LuxUni_EET *LuxUni_EETSession) Eethiq() (*types.Transaction, error) {
	return _LuxUni_EET.Contract.Eethiq(&_LuxUni_EET.TransactOpts)
}

// Eethiq is a paid mutator transaction binding the contract method 0x5f319d66.
//
// Solidity: function Eethiq() returns()
func (_LuxUni_EET *LuxUni_EETTransactorSession) Eethiq() (*types.Transaction, error) {
	return _LuxUni_EET.Contract.Eethiq(&_LuxUni_EET.TransactOpts)
}

// AddPoints is a paid mutator transaction binding the contract method 0x193a08d2.
//
// Solidity: function addPoints(_addr address, isReceiver bool, _amount uint256) returns(error int256)
func (_LuxUni_EET *LuxUni_EETTransactor) AddPoints(opts *bind.TransactOpts, _addr common.Address, isReceiver bool, _amount *big.Int) (*types.Transaction, error) {
	return _LuxUni_EET.contract.Transact(opts, "addPoints", _addr, isReceiver, _amount)
}

// AddPoints is a paid mutator transaction binding the contract method 0x193a08d2.
//
// Solidity: function addPoints(_addr address, isReceiver bool, _amount uint256) returns(error int256)
func (_LuxUni_EET *LuxUni_EETSession) AddPoints(_addr common.Address, isReceiver bool, _amount *big.Int) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.AddPoints(&_LuxUni_EET.TransactOpts, _addr, isReceiver, _amount)
}

// AddPoints is a paid mutator transaction binding the contract method 0x193a08d2.
//
// Solidity: function addPoints(_addr address, isReceiver bool, _amount uint256) returns(error int256)
func (_LuxUni_EET *LuxUni_EETTransactorSession) AddPoints(_addr common.Address, isReceiver bool, _amount *big.Int) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.AddPoints(&_LuxUni_EET.TransactOpts, _addr, isReceiver, _amount)
}

// MakeDonation is a paid mutator transaction binding the contract method 0x5d63e7e3.
//
// Solidity: function makeDonation(_senderAddr address, _receiverAddr address, _amount uint256, _priorityID int256, _certificate bytes) returns(_ret int256)
func (_LuxUni_EET *LuxUni_EETTransactor) MakeDonation(opts *bind.TransactOpts, _senderAddr common.Address, _receiverAddr common.Address, _amount *big.Int, _priorityID *big.Int, _certificate []byte) (*types.Transaction, error) {
	return _LuxUni_EET.contract.Transact(opts, "makeDonation", _senderAddr, _receiverAddr, _amount, _priorityID, _certificate)
}

// MakeDonation is a paid mutator transaction binding the contract method 0x5d63e7e3.
//
// Solidity: function makeDonation(_senderAddr address, _receiverAddr address, _amount uint256, _priorityID int256, _certificate bytes) returns(_ret int256)
func (_LuxUni_EET *LuxUni_EETSession) MakeDonation(_senderAddr common.Address, _receiverAddr common.Address, _amount *big.Int, _priorityID *big.Int, _certificate []byte) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.MakeDonation(&_LuxUni_EET.TransactOpts, _senderAddr, _receiverAddr, _amount, _priorityID, _certificate)
}

// MakeDonation is a paid mutator transaction binding the contract method 0x5d63e7e3.
//
// Solidity: function makeDonation(_senderAddr address, _receiverAddr address, _amount uint256, _priorityID int256, _certificate bytes) returns(_ret int256)
func (_LuxUni_EET *LuxUni_EETTransactorSession) MakeDonation(_senderAddr common.Address, _receiverAddr common.Address, _amount *big.Int, _priorityID *big.Int, _certificate []byte) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.MakeDonation(&_LuxUni_EET.TransactOpts, _senderAddr, _receiverAddr, _amount, _priorityID, _certificate)
}

// NewReceiver is a paid mutator transaction binding the contract method 0xcce19bcc.
//
// Solidity: function newReceiver(_receiverAddr address, _name string) returns(_ret int256)
func (_LuxUni_EET *LuxUni_EETTransactor) NewReceiver(opts *bind.TransactOpts, _receiverAddr common.Address, _name string) (*types.Transaction, error) {
	return _LuxUni_EET.contract.Transact(opts, "newReceiver", _receiverAddr, _name)
}

// NewReceiver is a paid mutator transaction binding the contract method 0xcce19bcc.
//
// Solidity: function newReceiver(_receiverAddr address, _name string) returns(_ret int256)
func (_LuxUni_EET *LuxUni_EETSession) NewReceiver(_receiverAddr common.Address, _name string) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.NewReceiver(&_LuxUni_EET.TransactOpts, _receiverAddr, _name)
}

// NewReceiver is a paid mutator transaction binding the contract method 0xcce19bcc.
//
// Solidity: function newReceiver(_receiverAddr address, _name string) returns(_ret int256)
func (_LuxUni_EET *LuxUni_EETTransactorSession) NewReceiver(_receiverAddr common.Address, _name string) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.NewReceiver(&_LuxUni_EET.TransactOpts, _receiverAddr, _name)
}

// NewSender is a paid mutator transaction binding the contract method 0x67b2f8e1.
//
// Solidity: function newSender(_senderAddr address, _name string) returns(_ret int256)
func (_LuxUni_EET *LuxUni_EETTransactor) NewSender(opts *bind.TransactOpts, _senderAddr common.Address, _name string) (*types.Transaction, error) {
	return _LuxUni_EET.contract.Transact(opts, "newSender", _senderAddr, _name)
}

// NewSender is a paid mutator transaction binding the contract method 0x67b2f8e1.
//
// Solidity: function newSender(_senderAddr address, _name string) returns(_ret int256)
func (_LuxUni_EET *LuxUni_EETSession) NewSender(_senderAddr common.Address, _name string) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.NewSender(&_LuxUni_EET.TransactOpts, _senderAddr, _name)
}

// NewSender is a paid mutator transaction binding the contract method 0x67b2f8e1.
//
// Solidity: function newSender(_senderAddr address, _name string) returns(_ret int256)
func (_LuxUni_EET *LuxUni_EETTransactorSession) NewSender(_senderAddr common.Address, _name string) (*types.Transaction, error) {
	return _LuxUni_EET.Contract.NewSender(&_LuxUni_EET.TransactOpts, _senderAddr, _name)
}
