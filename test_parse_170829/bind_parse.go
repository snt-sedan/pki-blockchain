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

// LuxUni_ParseABI is the input ABI used to generate the binding from.
const LuxUni_ParseABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"hex_str\",\"type\":\"string\"}],\"name\":\"hexStrToBytes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_der\",\"type\":\"bytes\"}],\"name\":\"ParseCert\",\"outputs\":[{\"name\":\"_addrParent\",\"type\":\"address\"},{\"name\":\"_addrCA\",\"type\":\"address\"},{\"name\":\"errCode\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"b\",\"type\":\"bytes\"}],\"name\":\"bytesToAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_der\",\"type\":\"bytes\"}],\"name\":\"ParseAddrParent\",\"outputs\":[{\"name\":\"_addrParent\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_der\",\"type\":\"bytes\"}],\"name\":\"ParseAddrCA\",\"outputs\":[{\"name\":\"_addrCA\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_ba\",\"type\":\"bytes\"}],\"name\":\"bytestr_to_uint\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"}]"

// LuxUni_Parse is an auto generated Go binding around an Ethereum contract.
type LuxUni_Parse struct {
	LuxUni_ParseCaller     // Read-only binding to the contract
	LuxUni_ParseTransactor // Write-only binding to the contract
}

// LuxUni_ParseCaller is an auto generated read-only Go binding around an Ethereum contract.
type LuxUni_ParseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LuxUni_ParseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LuxUni_ParseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LuxUni_ParseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LuxUni_ParseSession struct {
	Contract     *LuxUni_Parse     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LuxUni_ParseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LuxUni_ParseCallerSession struct {
	Contract *LuxUni_ParseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// LuxUni_ParseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LuxUni_ParseTransactorSession struct {
	Contract     *LuxUni_ParseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// LuxUni_ParseRaw is an auto generated low-level Go binding around an Ethereum contract.
type LuxUni_ParseRaw struct {
	Contract *LuxUni_Parse // Generic contract binding to access the raw methods on
}

// LuxUni_ParseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LuxUni_ParseCallerRaw struct {
	Contract *LuxUni_ParseCaller // Generic read-only contract binding to access the raw methods on
}

// LuxUni_ParseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LuxUni_ParseTransactorRaw struct {
	Contract *LuxUni_ParseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLuxUni_Parse creates a new instance of LuxUni_Parse, bound to a specific deployed contract.
func NewLuxUni_Parse(address common.Address, backend bind.ContractBackend) (*LuxUni_Parse, error) {
	contract, err := bindLuxUni_Parse(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LuxUni_Parse{LuxUni_ParseCaller: LuxUni_ParseCaller{contract: contract}, LuxUni_ParseTransactor: LuxUni_ParseTransactor{contract: contract}}, nil
}

// NewLuxUni_ParseCaller creates a new read-only instance of LuxUni_Parse, bound to a specific deployed contract.
func NewLuxUni_ParseCaller(address common.Address, caller bind.ContractCaller) (*LuxUni_ParseCaller, error) {
	contract, err := bindLuxUni_Parse(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &LuxUni_ParseCaller{contract: contract}, nil
}

// NewLuxUni_ParseTransactor creates a new write-only instance of LuxUni_Parse, bound to a specific deployed contract.
func NewLuxUni_ParseTransactor(address common.Address, transactor bind.ContractTransactor) (*LuxUni_ParseTransactor, error) {
	contract, err := bindLuxUni_Parse(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &LuxUni_ParseTransactor{contract: contract}, nil
}

// bindLuxUni_Parse binds a generic wrapper to an already deployed contract.
func bindLuxUni_Parse(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LuxUni_ParseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LuxUni_Parse *LuxUni_ParseRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LuxUni_Parse.Contract.LuxUni_ParseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LuxUni_Parse *LuxUni_ParseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_Parse.Contract.LuxUni_ParseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LuxUni_Parse *LuxUni_ParseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LuxUni_Parse.Contract.LuxUni_ParseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LuxUni_Parse *LuxUni_ParseCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LuxUni_Parse.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LuxUni_Parse *LuxUni_ParseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_Parse.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LuxUni_Parse *LuxUni_ParseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LuxUni_Parse.Contract.contract.Transact(opts, method, params...)
}

// ParseAddrCA is a free data retrieval call binding the contract method 0xc42bdedb.
//
// Solidity: function ParseAddrCA(_der bytes) constant returns(_addrCA address)
func (_LuxUni_Parse *LuxUni_ParseCaller) ParseAddrCA(opts *bind.CallOpts, _der []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_Parse.contract.Call(opts, out, "ParseAddrCA", _der)
	return *ret0, err
}

// ParseAddrCA is a free data retrieval call binding the contract method 0xc42bdedb.
//
// Solidity: function ParseAddrCA(_der bytes) constant returns(_addrCA address)
func (_LuxUni_Parse *LuxUni_ParseSession) ParseAddrCA(_der []byte) (common.Address, error) {
	return _LuxUni_Parse.Contract.ParseAddrCA(&_LuxUni_Parse.CallOpts, _der)
}

// ParseAddrCA is a free data retrieval call binding the contract method 0xc42bdedb.
//
// Solidity: function ParseAddrCA(_der bytes) constant returns(_addrCA address)
func (_LuxUni_Parse *LuxUni_ParseCallerSession) ParseAddrCA(_der []byte) (common.Address, error) {
	return _LuxUni_Parse.Contract.ParseAddrCA(&_LuxUni_Parse.CallOpts, _der)
}

// ParseAddrParent is a free data retrieval call binding the contract method 0x6f12af62.
//
// Solidity: function ParseAddrParent(_der bytes) constant returns(_addrParent address)
func (_LuxUni_Parse *LuxUni_ParseCaller) ParseAddrParent(opts *bind.CallOpts, _der []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_Parse.contract.Call(opts, out, "ParseAddrParent", _der)
	return *ret0, err
}

// ParseAddrParent is a free data retrieval call binding the contract method 0x6f12af62.
//
// Solidity: function ParseAddrParent(_der bytes) constant returns(_addrParent address)
func (_LuxUni_Parse *LuxUni_ParseSession) ParseAddrParent(_der []byte) (common.Address, error) {
	return _LuxUni_Parse.Contract.ParseAddrParent(&_LuxUni_Parse.CallOpts, _der)
}

// ParseAddrParent is a free data retrieval call binding the contract method 0x6f12af62.
//
// Solidity: function ParseAddrParent(_der bytes) constant returns(_addrParent address)
func (_LuxUni_Parse *LuxUni_ParseCallerSession) ParseAddrParent(_der []byte) (common.Address, error) {
	return _LuxUni_Parse.Contract.ParseAddrParent(&_LuxUni_Parse.CallOpts, _der)
}

// ParseCert is a free data retrieval call binding the contract method 0x42488564.
//
// Solidity: function ParseCert(_der bytes) constant returns(_addrParent address, _addrCA address, errCode uint256)
func (_LuxUni_Parse *LuxUni_ParseCaller) ParseCert(opts *bind.CallOpts, _der []byte) (struct {
	_addrParent common.Address
	_addrCA     common.Address
	ErrCode     *big.Int
}, error) {
	ret := new(struct {
		_addrParent common.Address
		_addrCA     common.Address
		ErrCode     *big.Int
	})
	out := ret
	err := _LuxUni_Parse.contract.Call(opts, out, "ParseCert", _der)
	return *ret, err
}

// ParseCert is a free data retrieval call binding the contract method 0x42488564.
//
// Solidity: function ParseCert(_der bytes) constant returns(_addrParent address, _addrCA address, errCode uint256)
func (_LuxUni_Parse *LuxUni_ParseSession) ParseCert(_der []byte) (struct {
	_addrParent common.Address
	_addrCA     common.Address
	ErrCode     *big.Int
}, error) {
	return _LuxUni_Parse.Contract.ParseCert(&_LuxUni_Parse.CallOpts, _der)
}

// ParseCert is a free data retrieval call binding the contract method 0x42488564.
//
// Solidity: function ParseCert(_der bytes) constant returns(_addrParent address, _addrCA address, errCode uint256)
func (_LuxUni_Parse *LuxUni_ParseCallerSession) ParseCert(_der []byte) (struct {
	_addrParent common.Address
	_addrCA     common.Address
	ErrCode     *big.Int
}, error) {
	return _LuxUni_Parse.Contract.ParseCert(&_LuxUni_Parse.CallOpts, _der)
}

// BytesToAddress is a free data retrieval call binding the contract method 0x42526e4e.
//
// Solidity: function bytesToAddress(b bytes) constant returns(address)
func (_LuxUni_Parse *LuxUni_ParseCaller) BytesToAddress(opts *bind.CallOpts, b []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_Parse.contract.Call(opts, out, "bytesToAddress", b)
	return *ret0, err
}

// BytesToAddress is a free data retrieval call binding the contract method 0x42526e4e.
//
// Solidity: function bytesToAddress(b bytes) constant returns(address)
func (_LuxUni_Parse *LuxUni_ParseSession) BytesToAddress(b []byte) (common.Address, error) {
	return _LuxUni_Parse.Contract.BytesToAddress(&_LuxUni_Parse.CallOpts, b)
}

// BytesToAddress is a free data retrieval call binding the contract method 0x42526e4e.
//
// Solidity: function bytesToAddress(b bytes) constant returns(address)
func (_LuxUni_Parse *LuxUni_ParseCallerSession) BytesToAddress(b []byte) (common.Address, error) {
	return _LuxUni_Parse.Contract.BytesToAddress(&_LuxUni_Parse.CallOpts, b)
}

// Bytestr_to_uint is a free data retrieval call binding the contract method 0xcd58a2e9.
//
// Solidity: function bytestr_to_uint(_ba bytes) constant returns(uint256)
func (_LuxUni_Parse *LuxUni_ParseCaller) Bytestr_to_uint(opts *bind.CallOpts, _ba []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_Parse.contract.Call(opts, out, "bytestr_to_uint", _ba)
	return *ret0, err
}

// Bytestr_to_uint is a free data retrieval call binding the contract method 0xcd58a2e9.
//
// Solidity: function bytestr_to_uint(_ba bytes) constant returns(uint256)
func (_LuxUni_Parse *LuxUni_ParseSession) Bytestr_to_uint(_ba []byte) (*big.Int, error) {
	return _LuxUni_Parse.Contract.Bytestr_to_uint(&_LuxUni_Parse.CallOpts, _ba)
}

// Bytestr_to_uint is a free data retrieval call binding the contract method 0xcd58a2e9.
//
// Solidity: function bytestr_to_uint(_ba bytes) constant returns(uint256)
func (_LuxUni_Parse *LuxUni_ParseCallerSession) Bytestr_to_uint(_ba []byte) (*big.Int, error) {
	return _LuxUni_Parse.Contract.Bytestr_to_uint(&_LuxUni_Parse.CallOpts, _ba)
}

// HexStrToBytes is a free data retrieval call binding the contract method 0x08e8ac38.
//
// Solidity: function hexStrToBytes(hex_str string) constant returns(bytes)
func (_LuxUni_Parse *LuxUni_ParseCaller) HexStrToBytes(opts *bind.CallOpts, hex_str string) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _LuxUni_Parse.contract.Call(opts, out, "hexStrToBytes", hex_str)
	return *ret0, err
}

// HexStrToBytes is a free data retrieval call binding the contract method 0x08e8ac38.
//
// Solidity: function hexStrToBytes(hex_str string) constant returns(bytes)
func (_LuxUni_Parse *LuxUni_ParseSession) HexStrToBytes(hex_str string) ([]byte, error) {
	return _LuxUni_Parse.Contract.HexStrToBytes(&_LuxUni_Parse.CallOpts, hex_str)
}

// HexStrToBytes is a free data retrieval call binding the contract method 0x08e8ac38.
//
// Solidity: function hexStrToBytes(hex_str string) constant returns(bytes)
func (_LuxUni_Parse *LuxUni_ParseCallerSession) HexStrToBytes(hex_str string) ([]byte, error) {
	return _LuxUni_Parse.Contract.HexStrToBytes(&_LuxUni_Parse.CallOpts, hex_str)
}
