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

// LuxUni_PKI_validABI is the input ABI used to generate the binding from.
const LuxUni_PKI_validABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_newHash\",\"type\":\"bytes32\"},{\"name\":\"_addrCA\",\"type\":\"address\"},{\"name\":\"_addrRoot\",\"type\":\"address\"}],\"name\":\"CheckCert\",\"outputs\":[{\"name\":\"_result\",\"type\":\"int256\"}],\"payable\":false,\"type\":\"function\"}]"

// LuxUni_PKI_valid is an auto generated Go binding around an Ethereum contract.
type LuxUni_PKI_valid struct {
	LuxUni_PKI_validCaller     // Read-only binding to the contract
	LuxUni_PKI_validTransactor // Write-only binding to the contract
}

// LuxUni_PKI_validCaller is an auto generated read-only Go binding around an Ethereum contract.
type LuxUni_PKI_validCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LuxUni_PKI_validTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LuxUni_PKI_validTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LuxUni_PKI_validSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LuxUni_PKI_validSession struct {
	Contract     *LuxUni_PKI_valid // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LuxUni_PKI_validCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LuxUni_PKI_validCallerSession struct {
	Contract *LuxUni_PKI_validCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// LuxUni_PKI_validTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LuxUni_PKI_validTransactorSession struct {
	Contract     *LuxUni_PKI_validTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// LuxUni_PKI_validRaw is an auto generated low-level Go binding around an Ethereum contract.
type LuxUni_PKI_validRaw struct {
	Contract *LuxUni_PKI_valid // Generic contract binding to access the raw methods on
}

// LuxUni_PKI_validCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LuxUni_PKI_validCallerRaw struct {
	Contract *LuxUni_PKI_validCaller // Generic read-only contract binding to access the raw methods on
}

// LuxUni_PKI_validTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LuxUni_PKI_validTransactorRaw struct {
	Contract *LuxUni_PKI_validTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLuxUni_PKI_valid creates a new instance of LuxUni_PKI_valid, bound to a specific deployed contract.
func NewLuxUni_PKI_valid(address common.Address, backend bind.ContractBackend) (*LuxUni_PKI_valid, error) {
	contract, err := bindLuxUni_PKI_valid(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LuxUni_PKI_valid{LuxUni_PKI_validCaller: LuxUni_PKI_validCaller{contract: contract}, LuxUni_PKI_validTransactor: LuxUni_PKI_validTransactor{contract: contract}}, nil
}

// NewLuxUni_PKI_validCaller creates a new read-only instance of LuxUni_PKI_valid, bound to a specific deployed contract.
func NewLuxUni_PKI_validCaller(address common.Address, caller bind.ContractCaller) (*LuxUni_PKI_validCaller, error) {
	contract, err := bindLuxUni_PKI_valid(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &LuxUni_PKI_validCaller{contract: contract}, nil
}

// NewLuxUni_PKI_validTransactor creates a new write-only instance of LuxUni_PKI_valid, bound to a specific deployed contract.
func NewLuxUni_PKI_validTransactor(address common.Address, transactor bind.ContractTransactor) (*LuxUni_PKI_validTransactor, error) {
	contract, err := bindLuxUni_PKI_valid(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &LuxUni_PKI_validTransactor{contract: contract}, nil
}

// bindLuxUni_PKI_valid binds a generic wrapper to an already deployed contract.
func bindLuxUni_PKI_valid(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LuxUni_PKI_validABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LuxUni_PKI_valid *LuxUni_PKI_validRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LuxUni_PKI_valid.Contract.LuxUni_PKI_validCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LuxUni_PKI_valid *LuxUni_PKI_validRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_PKI_valid.Contract.LuxUni_PKI_validTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LuxUni_PKI_valid *LuxUni_PKI_validRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LuxUni_PKI_valid.Contract.LuxUni_PKI_validTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LuxUni_PKI_valid *LuxUni_PKI_validCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LuxUni_PKI_valid.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LuxUni_PKI_valid *LuxUni_PKI_validTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_PKI_valid.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LuxUni_PKI_valid *LuxUni_PKI_validTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LuxUni_PKI_valid.Contract.contract.Transact(opts, method, params...)
}

// CheckCert is a free data retrieval call binding the contract method 0xb31c6071.
//
// Solidity: function CheckCert(_newHash bytes32, _addrCA address, _addrRoot address) constant returns(_result int256)
func (_LuxUni_PKI_valid *LuxUni_PKI_validCaller) CheckCert(opts *bind.CallOpts, _newHash [32]byte, _addrCA common.Address, _addrRoot common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_PKI_valid.contract.Call(opts, out, "CheckCert", _newHash, _addrCA, _addrRoot)
	return *ret0, err
}

// CheckCert is a free data retrieval call binding the contract method 0xb31c6071.
//
// Solidity: function CheckCert(_newHash bytes32, _addrCA address, _addrRoot address) constant returns(_result int256)
func (_LuxUni_PKI_valid *LuxUni_PKI_validSession) CheckCert(_newHash [32]byte, _addrCA common.Address, _addrRoot common.Address) (*big.Int, error) {
	return _LuxUni_PKI_valid.Contract.CheckCert(&_LuxUni_PKI_valid.CallOpts, _newHash, _addrCA, _addrRoot)
}

// CheckCert is a free data retrieval call binding the contract method 0xb31c6071.
//
// Solidity: function CheckCert(_newHash bytes32, _addrCA address, _addrRoot address) constant returns(_result int256)
func (_LuxUni_PKI_valid *LuxUni_PKI_validCallerSession) CheckCert(_newHash [32]byte, _addrCA common.Address, _addrRoot common.Address) (*big.Int, error) {
	return _LuxUni_PKI_valid.Contract.CheckCert(&_LuxUni_PKI_valid.CallOpts, _newHash, _addrCA, _addrRoot)
}
