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

// LuxUni_PKIABI is the input ABI used to generate the binding from.
const LuxUni_PKIABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_cryptoID\",\"type\":\"uint256\"},{\"name\":\"_bData\",\"type\":\"bytes\"}],\"name\":\"encryptCallBack\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deletedRegData\",\"outputs\":[{\"name\":\"nodeSender\",\"type\":\"address\"},{\"name\":\"deletionDate\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_regID\",\"type\":\"uint256\"}],\"name\":\"deleteRegDatum\",\"outputs\":[{\"name\":\"err\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"founder\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numRegData\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"encryptRegData\",\"outputs\":[{\"name\":\"nodeSender\",\"type\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"encryptDate\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_dataHash\",\"type\":\"bytes\"},{\"name\":\"_fileName\",\"type\":\"string\"},{\"name\":\"_description\",\"type\":\"string\"},{\"name\":\"_linkFile\",\"type\":\"string\"},{\"name\":\"_encrypted\",\"type\":\"uint256\"},{\"name\":\"_cryptoModule\",\"type\":\"address\"}],\"name\":\"newRegDatum\",\"outputs\":[{\"name\":\"_regID\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"regData\",\"outputs\":[{\"name\":\"nodeSender\",\"type\":\"address\"},{\"name\":\"dataHash\",\"type\":\"bytes\"},{\"name\":\"fileName\",\"type\":\"string\"},{\"name\":\"description\",\"type\":\"string\"},{\"name\":\"encrypted\",\"type\":\"uint256\"},{\"name\":\"cryptoModule\",\"type\":\"address\"},{\"name\":\"linkFile\",\"type\":\"string\"},{\"name\":\"creationDate\",\"type\":\"uint256\"},{\"name\":\"active\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_regID\",\"type\":\"uint256\"}],\"name\":\"evDataEncrypted\",\"type\":\"event\"}]"

// LuxUni_PKI is an auto generated Go binding around an Ethereum contract.
type LuxUni_PKI struct {
	LuxUni_PKICaller     // Read-only binding to the contract
	LuxUni_PKITransactor // Write-only binding to the contract
}

// LuxUni_PKICaller is an auto generated read-only Go binding around an Ethereum contract.
type LuxUni_PKICaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LuxUni_PKITransactor is an auto generated write-only Go binding around an Ethereum contract.
type LuxUni_PKITransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LuxUni_PKISession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LuxUni_PKISession struct {
	Contract     *LuxUni_PKI       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LuxUni_PKICallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LuxUni_PKICallerSession struct {
	Contract *LuxUni_PKICaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// LuxUni_PKITransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LuxUni_PKITransactorSession struct {
	Contract     *LuxUni_PKITransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// LuxUni_PKIRaw is an auto generated low-level Go binding around an Ethereum contract.
type LuxUni_PKIRaw struct {
	Contract *LuxUni_PKI // Generic contract binding to access the raw methods on
}

// LuxUni_PKICallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LuxUni_PKICallerRaw struct {
	Contract *LuxUni_PKICaller // Generic read-only contract binding to access the raw methods on
}

// LuxUni_PKITransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LuxUni_PKITransactorRaw struct {
	Contract *LuxUni_PKITransactor // Generic write-only contract binding to access the raw methods on
}

// NewLuxUni_PKI creates a new instance of LuxUni_PKI, bound to a specific deployed contract.
func NewLuxUni_PKI(address common.Address, backend bind.ContractBackend) (*LuxUni_PKI, error) {
	contract, err := bindLuxUni_PKI(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LuxUni_PKI{LuxUni_PKICaller: LuxUni_PKICaller{contract: contract}, LuxUni_PKITransactor: LuxUni_PKITransactor{contract: contract}}, nil
}

// NewLuxUni_PKICaller creates a new read-only instance of LuxUni_PKI, bound to a specific deployed contract.
func NewLuxUni_PKICaller(address common.Address, caller bind.ContractCaller) (*LuxUni_PKICaller, error) {
	contract, err := bindLuxUni_PKI(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &LuxUni_PKICaller{contract: contract}, nil
}

// NewLuxUni_PKITransactor creates a new write-only instance of LuxUni_PKI, bound to a specific deployed contract.
func NewLuxUni_PKITransactor(address common.Address, transactor bind.ContractTransactor) (*LuxUni_PKITransactor, error) {
	contract, err := bindLuxUni_PKI(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &LuxUni_PKITransactor{contract: contract}, nil
}

// bindLuxUni_PKI binds a generic wrapper to an already deployed contract.
func bindLuxUni_PKI(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LuxUni_PKIABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LuxUni_PKI *LuxUni_PKIRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LuxUni_PKI.Contract.LuxUni_PKICaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LuxUni_PKI *LuxUni_PKIRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.LuxUni_PKITransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LuxUni_PKI *LuxUni_PKIRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.LuxUni_PKITransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LuxUni_PKI *LuxUni_PKICallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LuxUni_PKI.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LuxUni_PKI *LuxUni_PKITransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LuxUni_PKI *LuxUni_PKITransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.contract.Transact(opts, method, params...)
}

// DeletedRegData is a free data retrieval call binding the contract method 0x313cea8b.
//
// Solidity: function deletedRegData( uint256) constant returns(nodeSender address, deletionDate uint256)
func (_LuxUni_PKI *LuxUni_PKICaller) DeletedRegData(opts *bind.CallOpts, arg0 *big.Int) (struct {
	NodeSender   common.Address
	DeletionDate *big.Int
}, error) {
	ret := new(struct {
		NodeSender   common.Address
		DeletionDate *big.Int
	})
	out := ret
	err := _LuxUni_PKI.contract.Call(opts, out, "deletedRegData", arg0)
	return *ret, err
}

// DeletedRegData is a free data retrieval call binding the contract method 0x313cea8b.
//
// Solidity: function deletedRegData( uint256) constant returns(nodeSender address, deletionDate uint256)
func (_LuxUni_PKI *LuxUni_PKISession) DeletedRegData(arg0 *big.Int) (struct {
	NodeSender   common.Address
	DeletionDate *big.Int
}, error) {
	return _LuxUni_PKI.Contract.DeletedRegData(&_LuxUni_PKI.CallOpts, arg0)
}

// DeletedRegData is a free data retrieval call binding the contract method 0x313cea8b.
//
// Solidity: function deletedRegData( uint256) constant returns(nodeSender address, deletionDate uint256)
func (_LuxUni_PKI *LuxUni_PKICallerSession) DeletedRegData(arg0 *big.Int) (struct {
	NodeSender   common.Address
	DeletionDate *big.Int
}, error) {
	return _LuxUni_PKI.Contract.DeletedRegData(&_LuxUni_PKI.CallOpts, arg0)
}

// EncryptRegData is a free data retrieval call binding the contract method 0xcc88d9a2.
//
// Solidity: function encryptRegData( uint256) constant returns(nodeSender address, data bytes, encryptDate uint256)
func (_LuxUni_PKI *LuxUni_PKICaller) EncryptRegData(opts *bind.CallOpts, arg0 *big.Int) (struct {
	NodeSender  common.Address
	Data        []byte
	EncryptDate *big.Int
}, error) {
	ret := new(struct {
		NodeSender  common.Address
		Data        []byte
		EncryptDate *big.Int
	})
	out := ret
	err := _LuxUni_PKI.contract.Call(opts, out, "encryptRegData", arg0)
	return *ret, err
}

// EncryptRegData is a free data retrieval call binding the contract method 0xcc88d9a2.
//
// Solidity: function encryptRegData( uint256) constant returns(nodeSender address, data bytes, encryptDate uint256)
func (_LuxUni_PKI *LuxUni_PKISession) EncryptRegData(arg0 *big.Int) (struct {
	NodeSender  common.Address
	Data        []byte
	EncryptDate *big.Int
}, error) {
	return _LuxUni_PKI.Contract.EncryptRegData(&_LuxUni_PKI.CallOpts, arg0)
}

// EncryptRegData is a free data retrieval call binding the contract method 0xcc88d9a2.
//
// Solidity: function encryptRegData( uint256) constant returns(nodeSender address, data bytes, encryptDate uint256)
func (_LuxUni_PKI *LuxUni_PKICallerSession) EncryptRegData(arg0 *big.Int) (struct {
	NodeSender  common.Address
	Data        []byte
	EncryptDate *big.Int
}, error) {
	return _LuxUni_PKI.Contract.EncryptRegData(&_LuxUni_PKI.CallOpts, arg0)
}

// Founder is a free data retrieval call binding the contract method 0x4d853ee5.
//
// Solidity: function founder() constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICaller) Founder(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "founder")
	return *ret0, err
}

// Founder is a free data retrieval call binding the contract method 0x4d853ee5.
//
// Solidity: function founder() constant returns(address)
func (_LuxUni_PKI *LuxUni_PKISession) Founder() (common.Address, error) {
	return _LuxUni_PKI.Contract.Founder(&_LuxUni_PKI.CallOpts)
}

// Founder is a free data retrieval call binding the contract method 0x4d853ee5.
//
// Solidity: function founder() constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICallerSession) Founder() (common.Address, error) {
	return _LuxUni_PKI.Contract.Founder(&_LuxUni_PKI.CallOpts)
}

// NumRegData is a free data retrieval call binding the contract method 0x83af428b.
//
// Solidity: function numRegData() constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKICaller) NumRegData(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "numRegData")
	return *ret0, err
}

// NumRegData is a free data retrieval call binding the contract method 0x83af428b.
//
// Solidity: function numRegData() constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKISession) NumRegData() (*big.Int, error) {
	return _LuxUni_PKI.Contract.NumRegData(&_LuxUni_PKI.CallOpts)
}

// NumRegData is a free data retrieval call binding the contract method 0x83af428b.
//
// Solidity: function numRegData() constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKICallerSession) NumRegData() (*big.Int, error) {
	return _LuxUni_PKI.Contract.NumRegData(&_LuxUni_PKI.CallOpts)
}

// RegData is a free data retrieval call binding the contract method 0xee16c4f5.
//
// Solidity: function regData( uint256) constant returns(nodeSender address, dataHash bytes, fileName string, description string, encrypted uint256, cryptoModule address, linkFile string, creationDate uint256, active bool)
func (_LuxUni_PKI *LuxUni_PKICaller) RegData(opts *bind.CallOpts, arg0 *big.Int) (struct {
	NodeSender   common.Address
	DataHash     []byte
	FileName     string
	Description  string
	Encrypted    *big.Int
	CryptoModule common.Address
	LinkFile     string
	CreationDate *big.Int
	Active       bool
}, error) {
	ret := new(struct {
		NodeSender   common.Address
		DataHash     []byte
		FileName     string
		Description  string
		Encrypted    *big.Int
		CryptoModule common.Address
		LinkFile     string
		CreationDate *big.Int
		Active       bool
	})
	out := ret
	err := _LuxUni_PKI.contract.Call(opts, out, "regData", arg0)
	return *ret, err
}

// RegData is a free data retrieval call binding the contract method 0xee16c4f5.
//
// Solidity: function regData( uint256) constant returns(nodeSender address, dataHash bytes, fileName string, description string, encrypted uint256, cryptoModule address, linkFile string, creationDate uint256, active bool)
func (_LuxUni_PKI *LuxUni_PKISession) RegData(arg0 *big.Int) (struct {
	NodeSender   common.Address
	DataHash     []byte
	FileName     string
	Description  string
	Encrypted    *big.Int
	CryptoModule common.Address
	LinkFile     string
	CreationDate *big.Int
	Active       bool
}, error) {
	return _LuxUni_PKI.Contract.RegData(&_LuxUni_PKI.CallOpts, arg0)
}

// RegData is a free data retrieval call binding the contract method 0xee16c4f5.
//
// Solidity: function regData( uint256) constant returns(nodeSender address, dataHash bytes, fileName string, description string, encrypted uint256, cryptoModule address, linkFile string, creationDate uint256, active bool)
func (_LuxUni_PKI *LuxUni_PKICallerSession) RegData(arg0 *big.Int) (struct {
	NodeSender   common.Address
	DataHash     []byte
	FileName     string
	Description  string
	Encrypted    *big.Int
	CryptoModule common.Address
	LinkFile     string
	CreationDate *big.Int
	Active       bool
}, error) {
	return _LuxUni_PKI.Contract.RegData(&_LuxUni_PKI.CallOpts, arg0)
}

// DeleteRegDatum is a paid mutator transaction binding the contract method 0x491a34f0.
//
// Solidity: function deleteRegDatum(_regID uint256) returns(err uint256)
func (_LuxUni_PKI *LuxUni_PKITransactor) DeleteRegDatum(opts *bind.TransactOpts, _regID *big.Int) (*types.Transaction, error) {
	return _LuxUni_PKI.contract.Transact(opts, "deleteRegDatum", _regID)
}

// DeleteRegDatum is a paid mutator transaction binding the contract method 0x491a34f0.
//
// Solidity: function deleteRegDatum(_regID uint256) returns(err uint256)
func (_LuxUni_PKI *LuxUni_PKISession) DeleteRegDatum(_regID *big.Int) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.DeleteRegDatum(&_LuxUni_PKI.TransactOpts, _regID)
}

// DeleteRegDatum is a paid mutator transaction binding the contract method 0x491a34f0.
//
// Solidity: function deleteRegDatum(_regID uint256) returns(err uint256)
func (_LuxUni_PKI *LuxUni_PKITransactorSession) DeleteRegDatum(_regID *big.Int) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.DeleteRegDatum(&_LuxUni_PKI.TransactOpts, _regID)
}

// EncryptCallBack is a paid mutator transaction binding the contract method 0x26659612.
//
// Solidity: function encryptCallBack(_cryptoID uint256, _bData bytes) returns()
func (_LuxUni_PKI *LuxUni_PKITransactor) EncryptCallBack(opts *bind.TransactOpts, _cryptoID *big.Int, _bData []byte) (*types.Transaction, error) {
	return _LuxUni_PKI.contract.Transact(opts, "encryptCallBack", _cryptoID, _bData)
}

// EncryptCallBack is a paid mutator transaction binding the contract method 0x26659612.
//
// Solidity: function encryptCallBack(_cryptoID uint256, _bData bytes) returns()
func (_LuxUni_PKI *LuxUni_PKISession) EncryptCallBack(_cryptoID *big.Int, _bData []byte) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.EncryptCallBack(&_LuxUni_PKI.TransactOpts, _cryptoID, _bData)
}

// EncryptCallBack is a paid mutator transaction binding the contract method 0x26659612.
//
// Solidity: function encryptCallBack(_cryptoID uint256, _bData bytes) returns()
func (_LuxUni_PKI *LuxUni_PKITransactorSession) EncryptCallBack(_cryptoID *big.Int, _bData []byte) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.EncryptCallBack(&_LuxUni_PKI.TransactOpts, _cryptoID, _bData)
}

// NewRegDatum is a paid mutator transaction binding the contract method 0xea5e2e20.
//
// Solidity: function newRegDatum(_dataHash bytes, _fileName string, _description string, _linkFile string, _encrypted uint256, _cryptoModule address) returns(_regID uint256)
func (_LuxUni_PKI *LuxUni_PKITransactor) NewRegDatum(opts *bind.TransactOpts, _dataHash []byte, _fileName string, _description string, _linkFile string, _encrypted *big.Int, _cryptoModule common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.contract.Transact(opts, "newRegDatum", _dataHash, _fileName, _description, _linkFile, _encrypted, _cryptoModule)
}

// NewRegDatum is a paid mutator transaction binding the contract method 0xea5e2e20.
//
// Solidity: function newRegDatum(_dataHash bytes, _fileName string, _description string, _linkFile string, _encrypted uint256, _cryptoModule address) returns(_regID uint256)
func (_LuxUni_PKI *LuxUni_PKISession) NewRegDatum(_dataHash []byte, _fileName string, _description string, _linkFile string, _encrypted *big.Int, _cryptoModule common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.NewRegDatum(&_LuxUni_PKI.TransactOpts, _dataHash, _fileName, _description, _linkFile, _encrypted, _cryptoModule)
}

// NewRegDatum is a paid mutator transaction binding the contract method 0xea5e2e20.
//
// Solidity: function newRegDatum(_dataHash bytes, _fileName string, _description string, _linkFile string, _encrypted uint256, _cryptoModule address) returns(_regID uint256)
func (_LuxUni_PKI *LuxUni_PKITransactorSession) NewRegDatum(_dataHash []byte, _fileName string, _description string, _linkFile string, _encrypted *big.Int, _cryptoModule common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.NewRegDatum(&_LuxUni_PKI.TransactOpts, _dataHash, _fileName, _description, _linkFile, _encrypted, _cryptoModule)
}
