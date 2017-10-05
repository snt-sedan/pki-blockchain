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

// LuxUni_PKI_webABI is the input ABI used to generate the binding from.
const LuxUni_PKI_webABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_addrContr\",\"type\":\"address\"},{\"name\":\"_index\",\"type\":\"uint96\"}],\"name\":\"getRegCreationDate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_parentAddr\",\"type\":\"address\"},{\"name\":\"_arrInd\",\"type\":\"uint256\"},{\"name\":\"_ethAccCA\",\"type\":\"address\"},{\"name\":\"_contrAddr\",\"type\":\"address\"},{\"name\":\"_fileName\",\"type\":\"string\"},{\"name\":\"_description\",\"type\":\"string\"}],\"name\":\"newRegDatum\",\"outputs\":[{\"name\":\"err\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addrContr\",\"type\":\"address\"},{\"name\":\"_index\",\"type\":\"uint96\"}],\"name\":\"getRegEthAccCA\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addrContr\",\"type\":\"address\"},{\"name\":\"_index\",\"type\":\"uint96\"}],\"name\":\"getRegContrAddr\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addrContr\",\"type\":\"address\"},{\"name\":\"_index\",\"type\":\"uint96\"}],\"name\":\"getRegDescription\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addrContr\",\"type\":\"address\"},{\"name\":\"_index\",\"type\":\"uint96\"}],\"name\":\"getRegFileName\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addrContr\",\"type\":\"address\"},{\"name\":\"_index\",\"type\":\"uint96\"}],\"name\":\"EncodeMapID\",\"outputs\":[{\"name\":\"_res\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"}]"

// LuxUni_PKI_web is an auto generated Go binding around an Ethereum contract.
type LuxUni_PKI_web struct {
	LuxUni_PKI_webCaller     // Read-only binding to the contract
	LuxUni_PKI_webTransactor // Write-only binding to the contract
}

// LuxUni_PKI_webCaller is an auto generated read-only Go binding around an Ethereum contract.
type LuxUni_PKI_webCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LuxUni_PKI_webTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LuxUni_PKI_webTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LuxUni_PKI_webSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LuxUni_PKI_webSession struct {
	Contract     *LuxUni_PKI_web   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LuxUni_PKI_webCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LuxUni_PKI_webCallerSession struct {
	Contract *LuxUni_PKI_webCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// LuxUni_PKI_webTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LuxUni_PKI_webTransactorSession struct {
	Contract     *LuxUni_PKI_webTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// LuxUni_PKI_webRaw is an auto generated low-level Go binding around an Ethereum contract.
type LuxUni_PKI_webRaw struct {
	Contract *LuxUni_PKI_web // Generic contract binding to access the raw methods on
}

// LuxUni_PKI_webCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LuxUni_PKI_webCallerRaw struct {
	Contract *LuxUni_PKI_webCaller // Generic read-only contract binding to access the raw methods on
}

// LuxUni_PKI_webTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LuxUni_PKI_webTransactorRaw struct {
	Contract *LuxUni_PKI_webTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLuxUni_PKI_web creates a new instance of LuxUni_PKI_web, bound to a specific deployed contract.
func NewLuxUni_PKI_web(address common.Address, backend bind.ContractBackend) (*LuxUni_PKI_web, error) {
	contract, err := bindLuxUni_PKI_web(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LuxUni_PKI_web{LuxUni_PKI_webCaller: LuxUni_PKI_webCaller{contract: contract}, LuxUni_PKI_webTransactor: LuxUni_PKI_webTransactor{contract: contract}}, nil
}

// NewLuxUni_PKI_webCaller creates a new read-only instance of LuxUni_PKI_web, bound to a specific deployed contract.
func NewLuxUni_PKI_webCaller(address common.Address, caller bind.ContractCaller) (*LuxUni_PKI_webCaller, error) {
	contract, err := bindLuxUni_PKI_web(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &LuxUni_PKI_webCaller{contract: contract}, nil
}

// NewLuxUni_PKI_webTransactor creates a new write-only instance of LuxUni_PKI_web, bound to a specific deployed contract.
func NewLuxUni_PKI_webTransactor(address common.Address, transactor bind.ContractTransactor) (*LuxUni_PKI_webTransactor, error) {
	contract, err := bindLuxUni_PKI_web(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &LuxUni_PKI_webTransactor{contract: contract}, nil
}

// bindLuxUni_PKI_web binds a generic wrapper to an already deployed contract.
func bindLuxUni_PKI_web(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LuxUni_PKI_webABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LuxUni_PKI_web *LuxUni_PKI_webRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LuxUni_PKI_web.Contract.LuxUni_PKI_webCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LuxUni_PKI_web *LuxUni_PKI_webRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_PKI_web.Contract.LuxUni_PKI_webTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LuxUni_PKI_web *LuxUni_PKI_webRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LuxUni_PKI_web.Contract.LuxUni_PKI_webTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LuxUni_PKI_web *LuxUni_PKI_webCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LuxUni_PKI_web.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LuxUni_PKI_web *LuxUni_PKI_webTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_PKI_web.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LuxUni_PKI_web *LuxUni_PKI_webTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LuxUni_PKI_web.Contract.contract.Transact(opts, method, params...)
}

// EncodeMapID is a free data retrieval call binding the contract method 0xe2ecd637.
//
// Solidity: function EncodeMapID(_addrContr address, _index uint96) constant returns(_res uint256)
func (_LuxUni_PKI_web *LuxUni_PKI_webCaller) EncodeMapID(opts *bind.CallOpts, _addrContr common.Address, _index *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_PKI_web.contract.Call(opts, out, "EncodeMapID", _addrContr, _index)
	return *ret0, err
}

// EncodeMapID is a free data retrieval call binding the contract method 0xe2ecd637.
//
// Solidity: function EncodeMapID(_addrContr address, _index uint96) constant returns(_res uint256)
func (_LuxUni_PKI_web *LuxUni_PKI_webSession) EncodeMapID(_addrContr common.Address, _index *big.Int) (*big.Int, error) {
	return _LuxUni_PKI_web.Contract.EncodeMapID(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// EncodeMapID is a free data retrieval call binding the contract method 0xe2ecd637.
//
// Solidity: function EncodeMapID(_addrContr address, _index uint96) constant returns(_res uint256)
func (_LuxUni_PKI_web *LuxUni_PKI_webCallerSession) EncodeMapID(_addrContr common.Address, _index *big.Int) (*big.Int, error) {
	return _LuxUni_PKI_web.Contract.EncodeMapID(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// GetRegContrAddr is a free data retrieval call binding the contract method 0xa46eb2c5.
//
// Solidity: function getRegContrAddr(_addrContr address, _index uint96) constant returns(address)
func (_LuxUni_PKI_web *LuxUni_PKI_webCaller) GetRegContrAddr(opts *bind.CallOpts, _addrContr common.Address, _index *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_PKI_web.contract.Call(opts, out, "getRegContrAddr", _addrContr, _index)
	return *ret0, err
}

// GetRegContrAddr is a free data retrieval call binding the contract method 0xa46eb2c5.
//
// Solidity: function getRegContrAddr(_addrContr address, _index uint96) constant returns(address)
func (_LuxUni_PKI_web *LuxUni_PKI_webSession) GetRegContrAddr(_addrContr common.Address, _index *big.Int) (common.Address, error) {
	return _LuxUni_PKI_web.Contract.GetRegContrAddr(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// GetRegContrAddr is a free data retrieval call binding the contract method 0xa46eb2c5.
//
// Solidity: function getRegContrAddr(_addrContr address, _index uint96) constant returns(address)
func (_LuxUni_PKI_web *LuxUni_PKI_webCallerSession) GetRegContrAddr(_addrContr common.Address, _index *big.Int) (common.Address, error) {
	return _LuxUni_PKI_web.Contract.GetRegContrAddr(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// GetRegCreationDate is a free data retrieval call binding the contract method 0x06862d3e.
//
// Solidity: function getRegCreationDate(_addrContr address, _index uint96) constant returns(uint256)
func (_LuxUni_PKI_web *LuxUni_PKI_webCaller) GetRegCreationDate(opts *bind.CallOpts, _addrContr common.Address, _index *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_PKI_web.contract.Call(opts, out, "getRegCreationDate", _addrContr, _index)
	return *ret0, err
}

// GetRegCreationDate is a free data retrieval call binding the contract method 0x06862d3e.
//
// Solidity: function getRegCreationDate(_addrContr address, _index uint96) constant returns(uint256)
func (_LuxUni_PKI_web *LuxUni_PKI_webSession) GetRegCreationDate(_addrContr common.Address, _index *big.Int) (*big.Int, error) {
	return _LuxUni_PKI_web.Contract.GetRegCreationDate(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// GetRegCreationDate is a free data retrieval call binding the contract method 0x06862d3e.
//
// Solidity: function getRegCreationDate(_addrContr address, _index uint96) constant returns(uint256)
func (_LuxUni_PKI_web *LuxUni_PKI_webCallerSession) GetRegCreationDate(_addrContr common.Address, _index *big.Int) (*big.Int, error) {
	return _LuxUni_PKI_web.Contract.GetRegCreationDate(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// GetRegDescription is a free data retrieval call binding the contract method 0xd2b9596b.
//
// Solidity: function getRegDescription(_addrContr address, _index uint96) constant returns(string)
func (_LuxUni_PKI_web *LuxUni_PKI_webCaller) GetRegDescription(opts *bind.CallOpts, _addrContr common.Address, _index *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LuxUni_PKI_web.contract.Call(opts, out, "getRegDescription", _addrContr, _index)
	return *ret0, err
}

// GetRegDescription is a free data retrieval call binding the contract method 0xd2b9596b.
//
// Solidity: function getRegDescription(_addrContr address, _index uint96) constant returns(string)
func (_LuxUni_PKI_web *LuxUni_PKI_webSession) GetRegDescription(_addrContr common.Address, _index *big.Int) (string, error) {
	return _LuxUni_PKI_web.Contract.GetRegDescription(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// GetRegDescription is a free data retrieval call binding the contract method 0xd2b9596b.
//
// Solidity: function getRegDescription(_addrContr address, _index uint96) constant returns(string)
func (_LuxUni_PKI_web *LuxUni_PKI_webCallerSession) GetRegDescription(_addrContr common.Address, _index *big.Int) (string, error) {
	return _LuxUni_PKI_web.Contract.GetRegDescription(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// GetRegEthAccCA is a free data retrieval call binding the contract method 0x72994190.
//
// Solidity: function getRegEthAccCA(_addrContr address, _index uint96) constant returns(address)
func (_LuxUni_PKI_web *LuxUni_PKI_webCaller) GetRegEthAccCA(opts *bind.CallOpts, _addrContr common.Address, _index *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_PKI_web.contract.Call(opts, out, "getRegEthAccCA", _addrContr, _index)
	return *ret0, err
}

// GetRegEthAccCA is a free data retrieval call binding the contract method 0x72994190.
//
// Solidity: function getRegEthAccCA(_addrContr address, _index uint96) constant returns(address)
func (_LuxUni_PKI_web *LuxUni_PKI_webSession) GetRegEthAccCA(_addrContr common.Address, _index *big.Int) (common.Address, error) {
	return _LuxUni_PKI_web.Contract.GetRegEthAccCA(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// GetRegEthAccCA is a free data retrieval call binding the contract method 0x72994190.
//
// Solidity: function getRegEthAccCA(_addrContr address, _index uint96) constant returns(address)
func (_LuxUni_PKI_web *LuxUni_PKI_webCallerSession) GetRegEthAccCA(_addrContr common.Address, _index *big.Int) (common.Address, error) {
	return _LuxUni_PKI_web.Contract.GetRegEthAccCA(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// GetRegFileName is a free data retrieval call binding the contract method 0xd5735f5c.
//
// Solidity: function getRegFileName(_addrContr address, _index uint96) constant returns(string)
func (_LuxUni_PKI_web *LuxUni_PKI_webCaller) GetRegFileName(opts *bind.CallOpts, _addrContr common.Address, _index *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LuxUni_PKI_web.contract.Call(opts, out, "getRegFileName", _addrContr, _index)
	return *ret0, err
}

// GetRegFileName is a free data retrieval call binding the contract method 0xd5735f5c.
//
// Solidity: function getRegFileName(_addrContr address, _index uint96) constant returns(string)
func (_LuxUni_PKI_web *LuxUni_PKI_webSession) GetRegFileName(_addrContr common.Address, _index *big.Int) (string, error) {
	return _LuxUni_PKI_web.Contract.GetRegFileName(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// GetRegFileName is a free data retrieval call binding the contract method 0xd5735f5c.
//
// Solidity: function getRegFileName(_addrContr address, _index uint96) constant returns(string)
func (_LuxUni_PKI_web *LuxUni_PKI_webCallerSession) GetRegFileName(_addrContr common.Address, _index *big.Int) (string, error) {
	return _LuxUni_PKI_web.Contract.GetRegFileName(&_LuxUni_PKI_web.CallOpts, _addrContr, _index)
}

// NewRegDatum is a paid mutator transaction binding the contract method 0x6fd4c3b9.
//
// Solidity: function newRegDatum(_parentAddr address, _arrInd uint256, _ethAccCA address, _contrAddr address, _fileName string, _description string) returns(err uint256)
func (_LuxUni_PKI_web *LuxUni_PKI_webTransactor) NewRegDatum(opts *bind.TransactOpts, _parentAddr common.Address, _arrInd *big.Int, _ethAccCA common.Address, _contrAddr common.Address, _fileName string, _description string) (*types.Transaction, error) {
	return _LuxUni_PKI_web.contract.Transact(opts, "newRegDatum", _parentAddr, _arrInd, _ethAccCA, _contrAddr, _fileName, _description)
}

// NewRegDatum is a paid mutator transaction binding the contract method 0x6fd4c3b9.
//
// Solidity: function newRegDatum(_parentAddr address, _arrInd uint256, _ethAccCA address, _contrAddr address, _fileName string, _description string) returns(err uint256)
func (_LuxUni_PKI_web *LuxUni_PKI_webSession) NewRegDatum(_parentAddr common.Address, _arrInd *big.Int, _ethAccCA common.Address, _contrAddr common.Address, _fileName string, _description string) (*types.Transaction, error) {
	return _LuxUni_PKI_web.Contract.NewRegDatum(&_LuxUni_PKI_web.TransactOpts, _parentAddr, _arrInd, _ethAccCA, _contrAddr, _fileName, _description)
}

// NewRegDatum is a paid mutator transaction binding the contract method 0x6fd4c3b9.
//
// Solidity: function newRegDatum(_parentAddr address, _arrInd uint256, _ethAccCA address, _contrAddr address, _fileName string, _description string) returns(err uint256)
func (_LuxUni_PKI_web *LuxUni_PKI_webTransactorSession) NewRegDatum(_parentAddr common.Address, _arrInd *big.Int, _ethAccCA common.Address, _contrAddr common.Address, _fileName string, _description string) (*types.Transaction, error) {
	return _LuxUni_PKI_web.Contract.NewRegDatum(&_LuxUni_PKI_web.TransactOpts, _parentAddr, _arrInd, _ethAccCA, _contrAddr, _fileName, _description)
}
