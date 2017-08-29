package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	//"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/ethclient"

	"encoding/json"
)

var session *LuxUni_Parse
var gClient *ethclient.Client

const gConfigFile = "./pki-conf.json"

/*
CORE PARAMETERS ARE STORED IN gConfig

const gCryptoModulHash = "0x3f2ed40488d0a9586013faa415718f3f64644fa1";
const gContractHash = "0xf1918d06d7e66e60153d7109ff380b41866ba2e0";
const gIPCpath = "/home/alex/_Work/Eth_AeroNet_t/geth.ipc";
const gPswd = "ira";
const key = `{"address":"a6f23407d139508fa38706140c56bf6487f87395","crypto":{"cipher":"aes-128-ctr","ciphertext":"5468728413080efcf6191dc3b5eaaf6ff34fd401277b9c08ac5aa93b6c3b3e44","cipherparams":{"iv":"efae50c39291e191193635ec715e96f4"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"a1e83898937dcd79e480d642aab80032a983a86e55c1ff8c64e91a1fc4fd9bdd"},"mac":"35456f2b746b847b4f1f5af2291f8f1518d6eed4e29ce950c1123da7ac3ee01d"},"id":"71412808-a83f-4677-ac6d-3d0f744d3ef0","version":3}`
const gPrivateKeyPath = "a_pr.key";
*/

var gPrivateKey *rsa.PrivateKey

var gConfig struct {
	//CryptoModulHash    string `json:"cryptoModulHash"`
	ContractHash   string `json:"contractHash"`
	IPCpath        string `json:"IPCpath"`
	Pswd           string `json:"pswd"`
	KeyPath        string `json:"keyPath"`
	PrivateKeyPath string `json:"privateKeyPath"`
	HttpPort       int    `json:"httpPort"`
	WebMode        int    `json:"webMode"`
	JsonMode       int    `json:"jsonMode"`
}

func init() {
	LoadConfig()
}

func main() {

	// how to connect to different servers
	//          (https://gist.github.com/bas-vk/299f4a686b66a22cf87302c561ee5866):
	//    geth --testnet --rpc
	// client, err := ethclient.Dial("http://localhost:8545")
	//    parity --testnet --port 31313 --jsonrpc-port 8546
	// client, err = ethclient.Dial("http://localhost:8546")

	// Create an IPC based RPC connection to a remote node
	// conn, err := rpc.NewIPCClient(gIPCpath)
	gClient, err := ethclient.Dial(gConfig.IPCpath)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Instantiate the contract, the address is taken from eth at the moment of contract initiation
	// kyc, err := NewLuxUni_KYC(common.HexToAddress(gContractHash), backends.NewRPCBackend(conn))
	pkiContract, err := NewLuxUni_Parse(common.HexToAddress(gConfig.ContractHash), gClient)
	if err != nil {
		log.Fatalf("Failed to instantiate a smart contract: %v", err)
	}

	// Logging into Ethereum as a user
	//key []byte;
	/* key, e := ioutil.ReadFile(gConfig.KeyPath)
	if e != nil {
		fmt.Printf("Key File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("Found Ethereum Key File \n")
	*/

	cert, e := ioutil.ReadFile("./cert/BlockchainTestIssuingCA3.cer")
	if e != nil {
		fmt.Printf("Certificate File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("Found Certificate Key File \n")

	callOpts := &bind.CallOpts{
		Pending: true,
	}

	//ret, err := pkiContract.ParseCert(callOpts, cert)
	retAddrCA, err := pkiContract.ParseAddrCA(callOpts, cert)
	if e != nil {
		fmt.Printf("Parse Cert CA Addr error: %v\n", e)
		os.Exit(1)
	}
	retAddrParent, err := pkiContract.ParseAddrParent(callOpts, cert)
	if e != nil {
		fmt.Printf("Parse Cert Parent Addr error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("addr pr: %v\n", retAddrParent.String())
	fmt.Printf("addr ca: %v\n", retAddrCA.String())

	/*key, err := ioutil.ReadFile(gConfig.KeyPath)
	if err != nil {
		log.Fatalf("Failed to read auth file: %v, err: %v", gConfig.KeyPath, err)
		return
	}

	auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	session := &LuxUni_ParseSession{
		Contract: pkiContract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:   auth.From,
			Signer: auth.Signer,
			//GasLimit: big.NewInt(2000000),
		},
	}
	session.TransactOpts = *auth
	session.TransactOpts.GasLimit = big.NewInt(200000000000)

	ret, e := session.ParseCert(cert)
	ret, e :=
	if e != nil {
		fmt.Printf("Parse Cert error: %v\n", e)
	}
	fmt.Printf(" Addr CA    : %v\n Addr parent: %v\n", ret._addrCA.String(), ret._addrParent.String())
	*/

}

/*
func ValidateCert(w http.ResponseWriter, r *http.Request) {
	var isCertFound bool = false

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
		return
	}

	//isCurl = r.MultipartForm.Value["Curl"]
	certHash := r.MultipartForm.Value["CertHash"][0]
	contrAddr := r.MultipartForm.Value["ContrAddr"][0]
	if common.IsHexAddress(contrAddr) == false {
		http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
			http.StatusInternalServerError)
	}

	currContract, err := NewLuxUni_PKI(common.HexToAddress(contrAddr), gClient)
	if err != nil {
		log.Fatalf("Failed to instantiate a smart contract: %v", err)
	}
	callOpts := &bind.CallOpts{
		Pending: true,
	}

	numCert, err := currContract.NumRegData(callOpts)
	if err != nil {
		http.Error(w, GeneralError{"NumRegData problems"}.Error(),
			http.StatusInternalServerError)
	}

	for i := int64(0); i < numCert.Int64(); i++ {
		bi := big.NewInt(i)
		//var certData CRegData;
		certData, err := currContract.RegData(callOpts, bi)
		if err != nil {
			http.Error(w, GeneralError{"RegData retrieval problems"}.Error(),
				http.StatusInternalServerError)
		}
		if string(certData.DataHash) == certHash {
			isCertFound = true
			break
		}
	}
	if isCertFound == true {
		err := CheckCa(common.HexToAddress(contrAddr))
		if err != nil {
			http.Error(w, GeneralError{"CA Cert not approved"}.Error(),
				http.StatusInternalServerError)
		}
	} else {
		http.Error(w, GeneralError{"Cert not found"}.Error(),
			http.StatusInternalServerError)
	}
}

func UpdateCertificate(file multipart.File, contrAddress string) error {
	return nil
}

func CheckCa(initAddr common.Address) error {
	//addr := common.HexToAddress(gConfig.ContractHash);
	addr := initAddr
	var maxIter int = 1000
	for i := 0; i < maxIter; i++ {
		addr = GetParent(addr)
		nilAddr := common.Address{}
		if addr == nilAddr {
			break
		}
		if (i == (maxIter - 1)) && (addr != nilAddr) {
			return GeneralError{"MaxIter limit is reached"}
		}
	}
	return nil
}

func GetParent(currAddr common.Address) common.Address {
	currContract, err := NewLuxUni_PKI(currAddr, gClient)
	if err != nil {
		log.Fatalf("Failed to instantiate a smart contract: %v", err)
	}
	callOpts := &bind.CallOpts{
		Pending: true,
	}

	//caCert, err := currContract.NumRegData(callOpts)
	caCert, err := currContract.CaCertificate(callOpts)
	parentAddr := caCert[gCaCertOffset : gCaCertOffset+20]
	return common.HexToAddress(string(parentAddr))
}
*/

func LoadConfig() error {
	file, err := os.Open(gConfigFile)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found configuration file: %s\n", gConfigFile)

	//jsonparser.Get(data, "person", "name", "fullName")
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&gConfig); err != nil {
		fmt.Printf("Parsing config file: %s\n", err.Error())
		os.Exit(1)
	}

	b, err := json.Marshal(gConfig)
	if err != nil {
		fmt.Printf("Cannot convert conf file into string: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded configuration file: %s\n", string(b))
	file.Close()
	return nil
}

type GeneralError struct {
	errMsg string
}

func (e GeneralError) Error() string {
	return e.errMsg
}
