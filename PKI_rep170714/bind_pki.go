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
const LuxUni_PKIABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_i\",\"type\":\"uint256\"}],\"name\":\"getRegDescription\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_cert\",\"type\":\"bytes\"}],\"name\":\"populateCertificate\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getCaCertificate\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_i\",\"type\":\"uint256\"}],\"name\":\"getRegFileName\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_regID\",\"type\":\"uint256\"}],\"name\":\"deleteRegDatum\",\"outputs\":[{\"name\":\"err\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_dataHash\",\"type\":\"bytes\"},{\"name\":\"_contrAddr\",\"type\":\"address\"},{\"name\":\"_fileName\",\"type\":\"string\"},{\"name\":\"_description\",\"type\":\"string\"},{\"name\":\"_linkFile\",\"type\":\"string\"},{\"name\":\"_nodeSender\",\"type\":\"address\"}],\"name\":\"newRegDatum\",\"outputs\":[{\"name\":\"_regID\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_i\",\"type\":\"uint256\"}],\"name\":\"getDeletedRegDate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_i\",\"type\":\"uint256\"}],\"name\":\"getRegContrAddr\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_i\",\"type\":\"uint256\"}],\"name\":\"getRegLinkFile\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_i\",\"type\":\"uint256\"}],\"name\":\"getRegCreationDate\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_i\",\"type\":\"uint256\"}],\"name\":\"getRegDataHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_i\",\"type\":\"uint256\"}],\"name\":\"getRegNodeSender\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getNumRegData\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_i\",\"type\":\"uint256\"}],\"name\":\"getDeletedRegNodeSender\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"}]"

// LuxUni_PKIBin is the compiled bytecode used for deploying new contracts.
const LuxUni_PKIBin = `60606040525b33600060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b5b611466806100576000396000f300606060405236156100e4576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806313af4035146100e6578063291d3deb1461011c5780632db9b7c5146101c35780633cb36a0a1461021d578063465bb0e0146102b6578063491a34f01461035d578063494dfd5514610391578063893d20e8146105065780638b63e8c6146105585780638df0835e1461058c578063a931411f146105ec578063aee47c2f14610693578063d6c56cc6146106c7578063e013b2af1461076e578063e904d953146107ce578063ed345a4b146107f4575bfe5b34156100ee57fe5b61011a600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610854565b005b341561012457fe5b61013a6004808035906020019091905050610899565b6040518080602001828103825283818151815260200191508051906020019080838360008314610189575b80518252602083111561018957602082019150602081019050602083039250610165565b505050905090810190601f1680156101b55780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156101cb57fe5b61021b600480803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091905050610966565b005b341561022557fe5b61022d6109df565b604051808060200182810382528381815181526020019150805190602001908083836000831461027c575b80518252602083111561027c57602082019150602081019050602083039250610258565b505050905090810190601f1680156102a85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156102be57fe5b6102d46004808035906020019091905050610a88565b6040518080602001828103825283818151815260200191508051906020019080838360008314610323575b805182526020831115610323576020820191506020810190506020830392506102ff565b505050905090810190601f16801561034f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561036557fe5b61037b6004808035906020019091905050610b55565b6040518082815260200191505060405180910390f35b341561039957fe5b6104f0600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610c8f565b6040518082815260200191505060405180910390f35b341561050e57fe5b610516610e76565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b341561056057fe5b6105766004808035906020019091905050610ea1565b6040518082815260200191505060405180910390f35b341561059457fe5b6105aa6004808035906020019091905050610ec2565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34156105f457fe5b61060a6004808035906020019091905050610f11565b6040518080602001828103825283818151815260200191508051906020019080838360008314610659575b80518252602083111561065957602082019150602081019050602083039250610635565b505050905090810190601f1680156106855780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561069b57fe5b6106b16004808035906020019091905050610fde565b6040518082815260200191505060405180910390f35b34156106cf57fe5b6106e5600480803590602001909190505061100d565b6040518080602001828103825283818151815260200191508051906020019080838360008314610734575b80518252602083111561073457602082019150602081019050602083039250610710565b505050905090810190601f1680156107605780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561077657fe5b61078c60048080359060200190919050506110da565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34156107d657fe5b6107de611129565b6040518082815260200191505060405180910390f35b34156107fc57fe5b6108126004808035906020019091905050611134565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b80600060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b50565b6108a1611175565b6003828154811015156108b057fe5b906000526020600020906007020160005b506004018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156109595780601f1061092e57610100808354040283529160200191610959565b820191906000526020600020905b81548152906001019060200180831161093c57829003601f168201915b505050505090505b919050565b600060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156109c35760006000fd5b80600190805190602001906109d9929190611189565b505b5b50565b6109e7611209565b60018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610a7d5780601f10610a5257610100808354040283529160200191610a7d565b820191906000526020600020905b815481529060010190602001808311610a6057829003601f168201915b505050505090505b90565b610a90611175565b600382815481101515610a9f57fe5b906000526020600020906007020160005b506003018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610b485780601f10610b1d57610100808354040283529160200191610b48565b820191906000526020600020905b815481529060010190602001808311610b2b57829003601f168201915b505050505090505b919050565b6000600060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610bb45760006000fd5b60025482101515610bc85760019050610c89565b60006004600084815260200190815260200160002060010154141515610bf15760029050610c89565b6040604051908101604052803373ffffffffffffffffffffffffffffffffffffffff168152602001428152506004600084815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160010155905050600090505b5b919050565b60006000600060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610cf05760006000fd5b60038054809190600101610d04919061121d565b9150600382815481101515610d1557fe5b906000526020600020906007020160005b509050828160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555087816001019080519060200190610d84929190611189565b50868160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555085816003019080519060200190610de092919061124f565b5084816004019080519060200190610df992919061124f565b50868160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555083816005019080519060200190610e5592919061124f565b50428160060181905550600182016002819055505b5b509695505050505050565b6000600060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690505b90565b6000600460008381526020019081526020016000206001015490505b919050565b6000600382815481101515610ed357fe5b906000526020600020906007020160005b5060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690505b919050565b610f19611175565b600382815481101515610f2857fe5b906000526020600020906007020160005b506005018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610fd15780601f10610fa657610100808354040283529160200191610fd1565b820191906000526020600020905b815481529060010190602001808311610fb457829003601f168201915b505050505090505b919050565b6000600382815481101515610fef57fe5b906000526020600020906007020160005b506006015490505b919050565b611015611209565b60038281548110151561102457fe5b906000526020600020906007020160005b506001018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156110cd5780601f106110a2576101008083540402835291602001916110cd565b820191906000526020600020905b8154815290600101906020018083116110b057829003601f168201915b505050505090505b919050565b60006003828154811015156110eb57fe5b906000526020600020906007020160005b5060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690505b919050565b600060025490505b90565b60006004600083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690505b919050565b602060405190810160405280600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106111ca57805160ff19168380011785556111f8565b828001600101855582156111f8579182015b828111156111f75782518255916020019190600101906111dc565b5b50905061120591906112cf565b5090565b602060405190810160405280600081525090565b81548183558181151161124a5760070281600702836000526020600020918201910161124991906112f4565b5b505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061129057805160ff19168380011785556112be565b828001600101855582156112be579182015b828111156112bd5782518255916020019190600101906112a2565b5b5090506112cb91906112cf565b5090565b6112f191905b808211156112ed5760008160009055506001016112d5565b5090565b90565b6113a791905b808211156113a35760006000820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905560018201600061133b91906113aa565b6002820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905560038201600061137291906113f2565b60048201600061138291906113f2565b60058201600061139291906113f2565b6006820160009055506007016112fa565b5090565b90565b50805460018160011615610100020316600290046000825580601f106113d057506113ef565b601f0160209004906000526020600020908101906113ee91906112cf565b5b50565b50805460018160011615610100020316600290046000825580601f106114185750611437565b601f01602090049060005260206000209081019061143691906112cf565b5b505600a165627a7a7230582072e6630240e0182cf3a56199b08738475308048e6e2aeaf9c280c0ae8e4263bb0029`

// DeployLuxUni_PKI deploys a new Ethereum contract, binding an instance of LuxUni_PKI to it.
func DeployLuxUni_PKI(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LuxUni_PKI, error) {
	parsed, err := abi.JSON(strings.NewReader(LuxUni_PKIABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LuxUni_PKIBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LuxUni_PKI{LuxUni_PKICaller: LuxUni_PKICaller{contract: contract}, LuxUni_PKITransactor: LuxUni_PKITransactor{contract: contract}}, nil
}

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

// GetCaCertificate is a free data retrieval call binding the contract method 0x3cb36a0a.
//
// Solidity: function getCaCertificate() constant returns(bytes)
func (_LuxUni_PKI *LuxUni_PKICaller) GetCaCertificate(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getCaCertificate")
	return *ret0, err
}

// GetCaCertificate is a free data retrieval call binding the contract method 0x3cb36a0a.
//
// Solidity: function getCaCertificate() constant returns(bytes)
func (_LuxUni_PKI *LuxUni_PKISession) GetCaCertificate() ([]byte, error) {
	return _LuxUni_PKI.Contract.GetCaCertificate(&_LuxUni_PKI.CallOpts)
}

// GetCaCertificate is a free data retrieval call binding the contract method 0x3cb36a0a.
//
// Solidity: function getCaCertificate() constant returns(bytes)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetCaCertificate() ([]byte, error) {
	return _LuxUni_PKI.Contract.GetCaCertificate(&_LuxUni_PKI.CallOpts)
}

// GetDeletedRegDate is a free data retrieval call binding the contract method 0x8b63e8c6.
//
// Solidity: function getDeletedRegDate(_i uint256) constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKICaller) GetDeletedRegDate(opts *bind.CallOpts, _i *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getDeletedRegDate", _i)
	return *ret0, err
}

// GetDeletedRegDate is a free data retrieval call binding the contract method 0x8b63e8c6.
//
// Solidity: function getDeletedRegDate(_i uint256) constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKISession) GetDeletedRegDate(_i *big.Int) (*big.Int, error) {
	return _LuxUni_PKI.Contract.GetDeletedRegDate(&_LuxUni_PKI.CallOpts, _i)
}

// GetDeletedRegDate is a free data retrieval call binding the contract method 0x8b63e8c6.
//
// Solidity: function getDeletedRegDate(_i uint256) constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetDeletedRegDate(_i *big.Int) (*big.Int, error) {
	return _LuxUni_PKI.Contract.GetDeletedRegDate(&_LuxUni_PKI.CallOpts, _i)
}

// GetDeletedRegNodeSender is a free data retrieval call binding the contract method 0xed345a4b.
//
// Solidity: function getDeletedRegNodeSender(_i uint256) constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICaller) GetDeletedRegNodeSender(opts *bind.CallOpts, _i *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getDeletedRegNodeSender", _i)
	return *ret0, err
}

// GetDeletedRegNodeSender is a free data retrieval call binding the contract method 0xed345a4b.
//
// Solidity: function getDeletedRegNodeSender(_i uint256) constant returns(address)
func (_LuxUni_PKI *LuxUni_PKISession) GetDeletedRegNodeSender(_i *big.Int) (common.Address, error) {
	return _LuxUni_PKI.Contract.GetDeletedRegNodeSender(&_LuxUni_PKI.CallOpts, _i)
}

// GetDeletedRegNodeSender is a free data retrieval call binding the contract method 0xed345a4b.
//
// Solidity: function getDeletedRegNodeSender(_i uint256) constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetDeletedRegNodeSender(_i *big.Int) (common.Address, error) {
	return _LuxUni_PKI.Contract.GetDeletedRegNodeSender(&_LuxUni_PKI.CallOpts, _i)
}

// GetNumRegData is a free data retrieval call binding the contract method 0xe904d953.
//
// Solidity: function getNumRegData() constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKICaller) GetNumRegData(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getNumRegData")
	return *ret0, err
}

// GetNumRegData is a free data retrieval call binding the contract method 0xe904d953.
//
// Solidity: function getNumRegData() constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKISession) GetNumRegData() (*big.Int, error) {
	return _LuxUni_PKI.Contract.GetNumRegData(&_LuxUni_PKI.CallOpts)
}

// GetNumRegData is a free data retrieval call binding the contract method 0xe904d953.
//
// Solidity: function getNumRegData() constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetNumRegData() (*big.Int, error) {
	return _LuxUni_PKI.Contract.GetNumRegData(&_LuxUni_PKI.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICaller) GetOwner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getOwner")
	return *ret0, err
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_LuxUni_PKI *LuxUni_PKISession) GetOwner() (common.Address, error) {
	return _LuxUni_PKI.Contract.GetOwner(&_LuxUni_PKI.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetOwner() (common.Address, error) {
	return _LuxUni_PKI.Contract.GetOwner(&_LuxUni_PKI.CallOpts)
}

// GetRegContrAddr is a free data retrieval call binding the contract method 0x8df0835e.
//
// Solidity: function getRegContrAddr(_i uint256) constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICaller) GetRegContrAddr(opts *bind.CallOpts, _i *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getRegContrAddr", _i)
	return *ret0, err
}

// GetRegContrAddr is a free data retrieval call binding the contract method 0x8df0835e.
//
// Solidity: function getRegContrAddr(_i uint256) constant returns(address)
func (_LuxUni_PKI *LuxUni_PKISession) GetRegContrAddr(_i *big.Int) (common.Address, error) {
	return _LuxUni_PKI.Contract.GetRegContrAddr(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegContrAddr is a free data retrieval call binding the contract method 0x8df0835e.
//
// Solidity: function getRegContrAddr(_i uint256) constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetRegContrAddr(_i *big.Int) (common.Address, error) {
	return _LuxUni_PKI.Contract.GetRegContrAddr(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegCreationDate is a free data retrieval call binding the contract method 0xaee47c2f.
//
// Solidity: function getRegCreationDate(_i uint256) constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKICaller) GetRegCreationDate(opts *bind.CallOpts, _i *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getRegCreationDate", _i)
	return *ret0, err
}

// GetRegCreationDate is a free data retrieval call binding the contract method 0xaee47c2f.
//
// Solidity: function getRegCreationDate(_i uint256) constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKISession) GetRegCreationDate(_i *big.Int) (*big.Int, error) {
	return _LuxUni_PKI.Contract.GetRegCreationDate(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegCreationDate is a free data retrieval call binding the contract method 0xaee47c2f.
//
// Solidity: function getRegCreationDate(_i uint256) constant returns(uint256)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetRegCreationDate(_i *big.Int) (*big.Int, error) {
	return _LuxUni_PKI.Contract.GetRegCreationDate(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegDataHash is a free data retrieval call binding the contract method 0xd6c56cc6.
//
// Solidity: function getRegDataHash(_i uint256) constant returns(bytes)
func (_LuxUni_PKI *LuxUni_PKICaller) GetRegDataHash(opts *bind.CallOpts, _i *big.Int) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getRegDataHash", _i)
	return *ret0, err
}

// GetRegDataHash is a free data retrieval call binding the contract method 0xd6c56cc6.
//
// Solidity: function getRegDataHash(_i uint256) constant returns(bytes)
func (_LuxUni_PKI *LuxUni_PKISession) GetRegDataHash(_i *big.Int) ([]byte, error) {
	return _LuxUni_PKI.Contract.GetRegDataHash(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegDataHash is a free data retrieval call binding the contract method 0xd6c56cc6.
//
// Solidity: function getRegDataHash(_i uint256) constant returns(bytes)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetRegDataHash(_i *big.Int) ([]byte, error) {
	return _LuxUni_PKI.Contract.GetRegDataHash(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegDescription is a free data retrieval call binding the contract method 0x291d3deb.
//
// Solidity: function getRegDescription(_i uint256) constant returns(string)
func (_LuxUni_PKI *LuxUni_PKICaller) GetRegDescription(opts *bind.CallOpts, _i *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getRegDescription", _i)
	return *ret0, err
}

// GetRegDescription is a free data retrieval call binding the contract method 0x291d3deb.
//
// Solidity: function getRegDescription(_i uint256) constant returns(string)
func (_LuxUni_PKI *LuxUni_PKISession) GetRegDescription(_i *big.Int) (string, error) {
	return _LuxUni_PKI.Contract.GetRegDescription(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegDescription is a free data retrieval call binding the contract method 0x291d3deb.
//
// Solidity: function getRegDescription(_i uint256) constant returns(string)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetRegDescription(_i *big.Int) (string, error) {
	return _LuxUni_PKI.Contract.GetRegDescription(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegFileName is a free data retrieval call binding the contract method 0x465bb0e0.
//
// Solidity: function getRegFileName(_i uint256) constant returns(string)
func (_LuxUni_PKI *LuxUni_PKICaller) GetRegFileName(opts *bind.CallOpts, _i *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getRegFileName", _i)
	return *ret0, err
}

// GetRegFileName is a free data retrieval call binding the contract method 0x465bb0e0.
//
// Solidity: function getRegFileName(_i uint256) constant returns(string)
func (_LuxUni_PKI *LuxUni_PKISession) GetRegFileName(_i *big.Int) (string, error) {
	return _LuxUni_PKI.Contract.GetRegFileName(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegFileName is a free data retrieval call binding the contract method 0x465bb0e0.
//
// Solidity: function getRegFileName(_i uint256) constant returns(string)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetRegFileName(_i *big.Int) (string, error) {
	return _LuxUni_PKI.Contract.GetRegFileName(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegLinkFile is a free data retrieval call binding the contract method 0xa931411f.
//
// Solidity: function getRegLinkFile(_i uint256) constant returns(string)
func (_LuxUni_PKI *LuxUni_PKICaller) GetRegLinkFile(opts *bind.CallOpts, _i *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getRegLinkFile", _i)
	return *ret0, err
}

// GetRegLinkFile is a free data retrieval call binding the contract method 0xa931411f.
//
// Solidity: function getRegLinkFile(_i uint256) constant returns(string)
func (_LuxUni_PKI *LuxUni_PKISession) GetRegLinkFile(_i *big.Int) (string, error) {
	return _LuxUni_PKI.Contract.GetRegLinkFile(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegLinkFile is a free data retrieval call binding the contract method 0xa931411f.
//
// Solidity: function getRegLinkFile(_i uint256) constant returns(string)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetRegLinkFile(_i *big.Int) (string, error) {
	return _LuxUni_PKI.Contract.GetRegLinkFile(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegNodeSender is a free data retrieval call binding the contract method 0xe013b2af.
//
// Solidity: function getRegNodeSender(_i uint256) constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICaller) GetRegNodeSender(opts *bind.CallOpts, _i *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "getRegNodeSender", _i)
	return *ret0, err
}

// GetRegNodeSender is a free data retrieval call binding the contract method 0xe013b2af.
//
// Solidity: function getRegNodeSender(_i uint256) constant returns(address)
func (_LuxUni_PKI *LuxUni_PKISession) GetRegNodeSender(_i *big.Int) (common.Address, error) {
	return _LuxUni_PKI.Contract.GetRegNodeSender(&_LuxUni_PKI.CallOpts, _i)
}

// GetRegNodeSender is a free data retrieval call binding the contract method 0xe013b2af.
//
// Solidity: function getRegNodeSender(_i uint256) constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICallerSession) GetRegNodeSender(_i *big.Int) (common.Address, error) {
	return _LuxUni_PKI.Contract.GetRegNodeSender(&_LuxUni_PKI.CallOpts, _i)
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

// NewRegDatum is a paid mutator transaction binding the contract method 0x494dfd55.
//
// Solidity: function newRegDatum(_dataHash bytes, _contrAddr address, _fileName string, _description string, _linkFile string, _nodeSender address) returns(_regID uint256)
func (_LuxUni_PKI *LuxUni_PKITransactor) NewRegDatum(opts *bind.TransactOpts, _dataHash []byte, _contrAddr common.Address, _fileName string, _description string, _linkFile string, _nodeSender common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.contract.Transact(opts, "newRegDatum", _dataHash, _contrAddr, _fileName, _description, _linkFile, _nodeSender)
}

// NewRegDatum is a paid mutator transaction binding the contract method 0x494dfd55.
//
// Solidity: function newRegDatum(_dataHash bytes, _contrAddr address, _fileName string, _description string, _linkFile string, _nodeSender address) returns(_regID uint256)
func (_LuxUni_PKI *LuxUni_PKISession) NewRegDatum(_dataHash []byte, _contrAddr common.Address, _fileName string, _description string, _linkFile string, _nodeSender common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.NewRegDatum(&_LuxUni_PKI.TransactOpts, _dataHash, _contrAddr, _fileName, _description, _linkFile, _nodeSender)
}

// NewRegDatum is a paid mutator transaction binding the contract method 0x494dfd55.
//
// Solidity: function newRegDatum(_dataHash bytes, _contrAddr address, _fileName string, _description string, _linkFile string, _nodeSender address) returns(_regID uint256)
func (_LuxUni_PKI *LuxUni_PKITransactorSession) NewRegDatum(_dataHash []byte, _contrAddr common.Address, _fileName string, _description string, _linkFile string, _nodeSender common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.NewRegDatum(&_LuxUni_PKI.TransactOpts, _dataHash, _contrAddr, _fileName, _description, _linkFile, _nodeSender)
}

// PopulateCertificate is a paid mutator transaction binding the contract method 0x2db9b7c5.
//
// Solidity: function populateCertificate(_cert bytes) returns()
func (_LuxUni_PKI *LuxUni_PKITransactor) PopulateCertificate(opts *bind.TransactOpts, _cert []byte) (*types.Transaction, error) {
	return _LuxUni_PKI.contract.Transact(opts, "populateCertificate", _cert)
}

// PopulateCertificate is a paid mutator transaction binding the contract method 0x2db9b7c5.
//
// Solidity: function populateCertificate(_cert bytes) returns()
func (_LuxUni_PKI *LuxUni_PKISession) PopulateCertificate(_cert []byte) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.PopulateCertificate(&_LuxUni_PKI.TransactOpts, _cert)
}

// PopulateCertificate is a paid mutator transaction binding the contract method 0x2db9b7c5.
//
// Solidity: function populateCertificate(_cert bytes) returns()
func (_LuxUni_PKI *LuxUni_PKITransactorSession) PopulateCertificate(_cert []byte) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.PopulateCertificate(&_LuxUni_PKI.TransactOpts, _cert)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(_addr address) returns()
func (_LuxUni_PKI *LuxUni_PKITransactor) SetOwner(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.contract.Transact(opts, "setOwner", _addr)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(_addr address) returns()
func (_LuxUni_PKI *LuxUni_PKISession) SetOwner(_addr common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.SetOwner(&_LuxUni_PKI.TransactOpts, _addr)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(_addr address) returns()
func (_LuxUni_PKI *LuxUni_PKITransactorSession) SetOwner(_addr common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.SetOwner(&_LuxUni_PKI.TransactOpts, _addr)
}
