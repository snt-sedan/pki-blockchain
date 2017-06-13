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
const LuxUni_PKIABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_cert\",\"type\":\"bytes\"}],\"name\":\"populateCertificate\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deletedRegData\",\"outputs\":[{\"name\":\"nodeSender\",\"type\":\"address\"},{\"name\":\"deletionDate\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"caCertificate\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_regID\",\"type\":\"uint256\"}],\"name\":\"deleteRegDatum\",\"outputs\":[{\"name\":\"err\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"caAddr\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numRegData\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_dataHash\",\"type\":\"bytes\"},{\"name\":\"_contrAddr\",\"type\":\"address\"},{\"name\":\"_fileName\",\"type\":\"string\"},{\"name\":\"_description\",\"type\":\"string\"},{\"name\":\"_linkFile\",\"type\":\"string\"},{\"name\":\"_encrypted\",\"type\":\"uint256\"},{\"name\":\"_cryptoModule\",\"type\":\"address\"}],\"name\":\"newRegDatum\",\"outputs\":[{\"name\":\"_regID\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"regData\",\"outputs\":[{\"name\":\"nodeSender\",\"type\":\"address\"},{\"name\":\"dataHash\",\"type\":\"bytes\"},{\"name\":\"contrAddr\",\"type\":\"address\"},{\"name\":\"fileName\",\"type\":\"string\"},{\"name\":\"description\",\"type\":\"string\"},{\"name\":\"linkFile\",\"type\":\"string\"},{\"name\":\"creationDate\",\"type\":\"uint256\"},{\"name\":\"active\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"payable\":false,\"type\":\"constructor\"}]"

// LuxUni_PKIBin is the compiled bytecode used for deploying new contracts.
const LuxUni_PKIBin = `6060604052341561000c57fe5b604051602080610fe2833981016040528080519060200190919050505b60008173ffffffffffffffffffffffffffffffffffffffff16141561004c573390505b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b505b610f438061009f6000396000f3006060604052361561008c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632db9b7c51461008e578063313cea8b146100e8578063349ccd131461014f578063491a34f0146101e85780637f2cc11b1461021c57806383af428b1461026e5780639d5fd86f14610294578063ee16c4f514610412575bfe5b341561009657fe5b6100e6600480803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919050506106db565b005b34156100f057fe5b61010660048080359060200190919050506106f6565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b341561015757fe5b61015f61073a565b60405180806020018281038252838181518152602001915080519060200190808383600083146101ae575b8051825260208311156101ae5760208201915060208101905060208303925061018a565b505050905090810190601f1680156101da5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156101f057fe5b61020660048080359060200190919050506107d8565b6040518082815260200191505060405180910390f35b341561022457fe5b61022c610911565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b341561027657fe5b61027e610937565b6040518082815260200191505060405180910390f35b341561029c57fe5b6103fc600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001909190803573ffffffffffffffffffffffffffffffffffffffff1690602001909190505061093d565b6040518082815260200191505060405180910390f35b341561041a57fe5b6104306004808035906020019091905050610b41565b604051808973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001806020018873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018060200180602001806020018781526020018615151515815260200185810385528c8181546001816001161561010002031660029004815260200191508054600181600116156101000203166002900480156105395780601f1061050e57610100808354040283529160200191610539565b820191906000526020600020905b81548152906001019060200180831161051c57829003601f168201915b505085810384528a8181546001816001161561010002031660029004815260200191508054600181600116156101000203166002900480156105bc5780601f10610591576101008083540402835291602001916105bc565b820191906000526020600020905b81548152906001019060200180831161059f57829003601f168201915b505085810383528981815460018160011615610100020316600290048152602001915080546001816001161561010002031660029004801561063f5780601f106106145761010080835404028352916020019161063f565b820191906000526020600020905b81548152906001019060200180831161062257829003601f168201915b50508581038252888181546001816001161561010002031660029004815260200191508054600181600116156101000203166002900480156106c25780601f10610697576101008083540402835291602001916106c2565b820191906000526020600020905b8154815290600101906020018083116106a557829003601f168201915b50509c5050505050505050505050505060405180910390f35b80600090805190602001906106f1929190610be2565b505b50565b60046020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010154905082565b60008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107d05780601f106107a5576101008083540402835291602001916107d0565b820191906000526020600020905b8154815290600101906020018083116107b357829003601f168201915b505050505081565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156108375760006000fd5b6002548210151561084b576001905061090c565b60006004600084815260200190815260200160002060010154141515610874576002905061090c565b6040604051908101604052803373ffffffffffffffffffffffffffffffffffffffff168152602001428152506004600084815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160010155905050600090505b919050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60025481565b60006000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561099e5760006000fd5b600380548091906001016109b29190610c62565b91506003828154811015156109c357fe5b90600052602060002090600a020160005b509050338160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555088816001019080519060200190610a32929190610be2565b50878160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555086816003019080519060200190610a8e929190610c94565b5085816004019080519060200190610aa7929190610c94565b50878160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555084816005019080519060200190610b03929190610c94565b5042816006018190555060018160070160006101000a81548160ff021916908315150217905550600182016002819055505b50979650505050505050565b600381815481101515610b5057fe5b90600052602060002090600a020160005b915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169080600101908060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060030190806004019080600501908060060154908060070160009054906101000a900460ff16905088565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610c2357805160ff1916838001178555610c51565b82800160010185558215610c51579182015b82811115610c50578251825591602001919060010190610c35565b5b509050610c5e9190610d14565b5090565b815481835581811511610c8f57600a0281600a028360005260206000209182019101610c8e9190610d39565b5b505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610cd557805160ff1916838001178555610d03565b82800160010185558215610d03579182015b82811115610d02578251825591602001919060010190610ce7565b5b509050610d109190610d14565b5090565b610d3691905b80821115610d32576000816000905550600101610d1a565b5090565b90565b610e1091905b80821115610e0c5760006000820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600182016000610d809190610e13565b6002820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600382016000610db79190610e5b565b600482016000610dc79190610e5b565b600582016000610dd79190610e5b565b60068201600090556007820160006101000a81549060ff0219169055600882016000610e039190610ea3565b50600a01610d3f565b5090565b90565b50805460018160011615610100020316600290046000825580601f10610e395750610e58565b601f016020900490600052602060002090810190610e579190610d14565b5b50565b50805460018160011615610100020316600290046000825580601f10610e815750610ea0565b601f016020900490600052602060002090810190610e9f9190610d14565b5b50565b5080546000825560020290600052602060002090810190610ec49190610ec8565b5b50565b610f1491905b80821115610f1057600060008201600090556001820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905550600201610ece565b5090565b905600a165627a7a72305820a4a8a2dffc45446c2cedba30ca719a1a0e44720474feb2015a5764e1e5f4a8790029`

// DeployLuxUni_PKI deploys a new Ethereum contract, binding an instance of LuxUni_PKI to it.
func DeployLuxUni_PKI(auth *bind.TransactOpts, backend bind.ContractBackend, _addr common.Address) (common.Address, *types.Transaction, *LuxUni_PKI, error) {
	parsed, err := abi.JSON(strings.NewReader(LuxUni_PKIABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LuxUni_PKIBin), backend, _addr)
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

// CaAddr is a free data retrieval call binding the contract method 0x7f2cc11b.
//
// Solidity: function caAddr() constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICaller) CaAddr(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "caAddr")
	return *ret0, err
}

// CaAddr is a free data retrieval call binding the contract method 0x7f2cc11b.
//
// Solidity: function caAddr() constant returns(address)
func (_LuxUni_PKI *LuxUni_PKISession) CaAddr() (common.Address, error) {
	return _LuxUni_PKI.Contract.CaAddr(&_LuxUni_PKI.CallOpts)
}

// CaAddr is a free data retrieval call binding the contract method 0x7f2cc11b.
//
// Solidity: function caAddr() constant returns(address)
func (_LuxUni_PKI *LuxUni_PKICallerSession) CaAddr() (common.Address, error) {
	return _LuxUni_PKI.Contract.CaAddr(&_LuxUni_PKI.CallOpts)
}

// CaCertificate is a free data retrieval call binding the contract method 0x349ccd13.
//
// Solidity: function caCertificate() constant returns(bytes)
func (_LuxUni_PKI *LuxUni_PKICaller) CaCertificate(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _LuxUni_PKI.contract.Call(opts, out, "caCertificate")
	return *ret0, err
}

// CaCertificate is a free data retrieval call binding the contract method 0x349ccd13.
//
// Solidity: function caCertificate() constant returns(bytes)
func (_LuxUni_PKI *LuxUni_PKISession) CaCertificate() ([]byte, error) {
	return _LuxUni_PKI.Contract.CaCertificate(&_LuxUni_PKI.CallOpts)
}

// CaCertificate is a free data retrieval call binding the contract method 0x349ccd13.
//
// Solidity: function caCertificate() constant returns(bytes)
func (_LuxUni_PKI *LuxUni_PKICallerSession) CaCertificate() ([]byte, error) {
	return _LuxUni_PKI.Contract.CaCertificate(&_LuxUni_PKI.CallOpts)
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
// Solidity: function regData( uint256) constant returns(nodeSender address, dataHash bytes, contrAddr address, fileName string, description string, linkFile string, creationDate uint256, active bool)
func (_LuxUni_PKI *LuxUni_PKICaller) RegData(opts *bind.CallOpts, arg0 *big.Int) (struct {
	NodeSender   common.Address
	DataHash     []byte
	ContrAddr    common.Address
	FileName     string
	Description  string
	LinkFile     string
	CreationDate *big.Int
	Active       bool
}, error) {
	ret := new(struct {
		NodeSender   common.Address
		DataHash     []byte
		ContrAddr    common.Address
		FileName     string
		Description  string
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
// Solidity: function regData( uint256) constant returns(nodeSender address, dataHash bytes, contrAddr address, fileName string, description string, linkFile string, creationDate uint256, active bool)
func (_LuxUni_PKI *LuxUni_PKISession) RegData(arg0 *big.Int) (struct {
	NodeSender   common.Address
	DataHash     []byte
	ContrAddr    common.Address
	FileName     string
	Description  string
	LinkFile     string
	CreationDate *big.Int
	Active       bool
}, error) {
	return _LuxUni_PKI.Contract.RegData(&_LuxUni_PKI.CallOpts, arg0)
}

// RegData is a free data retrieval call binding the contract method 0xee16c4f5.
//
// Solidity: function regData( uint256) constant returns(nodeSender address, dataHash bytes, contrAddr address, fileName string, description string, linkFile string, creationDate uint256, active bool)
func (_LuxUni_PKI *LuxUni_PKICallerSession) RegData(arg0 *big.Int) (struct {
	NodeSender   common.Address
	DataHash     []byte
	ContrAddr    common.Address
	FileName     string
	Description  string
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

// NewRegDatum is a paid mutator transaction binding the contract method 0x9d5fd86f.
//
// Solidity: function newRegDatum(_dataHash bytes, _contrAddr address, _fileName string, _description string, _linkFile string, _encrypted uint256, _cryptoModule address) returns(_regID uint256)
func (_LuxUni_PKI *LuxUni_PKITransactor) NewRegDatum(opts *bind.TransactOpts, _dataHash []byte, _contrAddr common.Address, _fileName string, _description string, _linkFile string, _encrypted *big.Int, _cryptoModule common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.contract.Transact(opts, "newRegDatum", _dataHash, _contrAddr, _fileName, _description, _linkFile, _encrypted, _cryptoModule)
}

// NewRegDatum is a paid mutator transaction binding the contract method 0x9d5fd86f.
//
// Solidity: function newRegDatum(_dataHash bytes, _contrAddr address, _fileName string, _description string, _linkFile string, _encrypted uint256, _cryptoModule address) returns(_regID uint256)
func (_LuxUni_PKI *LuxUni_PKISession) NewRegDatum(_dataHash []byte, _contrAddr common.Address, _fileName string, _description string, _linkFile string, _encrypted *big.Int, _cryptoModule common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.NewRegDatum(&_LuxUni_PKI.TransactOpts, _dataHash, _contrAddr, _fileName, _description, _linkFile, _encrypted, _cryptoModule)
}

// NewRegDatum is a paid mutator transaction binding the contract method 0x9d5fd86f.
//
// Solidity: function newRegDatum(_dataHash bytes, _contrAddr address, _fileName string, _description string, _linkFile string, _encrypted uint256, _cryptoModule address) returns(_regID uint256)
func (_LuxUni_PKI *LuxUni_PKITransactorSession) NewRegDatum(_dataHash []byte, _contrAddr common.Address, _fileName string, _description string, _linkFile string, _encrypted *big.Int, _cryptoModule common.Address) (*types.Transaction, error) {
	return _LuxUni_PKI.Contract.NewRegDatum(&_LuxUni_PKI.TransactOpts, _dataHash, _contrAddr, _fileName, _description, _linkFile, _encrypted, _cryptoModule)
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
