package main

import (
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
)

var gConfigFile string = "./config/pki-conf.json"

var gPrivateKey *rsa.PrivateKey

var gConfig struct {
	//CryptoModulHash    string `json:"cryptoModulHash"`
	ContractHash    string `json:"contractHash"`
	ContractWebHash string `json:"contractWebHash"`
	IPCpath         string `json:"IPCpath"`
	Pswd            string `json:"pswd"`
	KeyDir          string `json:"keyDir"`
	AccountAddr     string `json:"accountAddr"` // !! should start with 0x
	PrivateKeyPath  string `json:"privateKeyPath"`
	RestHttpPort    int    `json:"restHttpPort"`
	AppHttpPort     int    `json:"appHttpPort"`
	EthereumRpcPort int    `json:"ethereumRpcPort"`
	RestUrlServer   string `json:"restUrlServer"`  //for example -- localhost
	AppUrlServer    string `json:"appUrlServer"`   //for example -- localhost
	EthereumRpcUrl  string `json:"ethereumRpcUrl"` //for example -- localhost
	WebMode         int    `json:"webMode"`
	JsonMode        int    `json:"jsonMode"`
	FileWebPath     string `json:"fileWebPath"` // !! should start with 0x
}

type JsonValidateNode struct {
	ContrAddr  string
	ParentAddr string
	Hash       string
	RevokeDate string
	IsCertOK   string
	Warn       string
}
type JsonValidateResponse struct {
	Status      int
	IsCertOK    bool
	RevokeDate  time.Time
	Iter        int
	ProcMilisec int64
	CertPath    []JsonValidateNode
}

/*
   idFormStr -- id of the loading files in the html form, for example "UplFiles"
   return params
       []byte -- hash (not uint32)
       string -- fileName
       //string -- path to the temporary file on the disk
       []byte -- pointer to the file content
       *GeneralCodeError
   return Err code:
       1 -- no file is found
       2 -- too many files are uploaded
       3 -- others
*/
func UploadFile(w http.ResponseWriter, r *http.Request, idFormStr string,
	isCopyContent bool) ([]byte, string, []byte, GeneralCodeError) {
	m := r.MultipartForm
	files := m.File[idFormStr]

	if len(files) == 0 {
		return nil, "", nil, GeneralCodeError{"Upload File: no files are found)", 1}
	}

	if len(files) > 1 {
		return nil, "", nil, GeneralCodeError{fmt.Sprintf(
			"Upload File: number of aploaded files greater than 1 = ", len(files)), 2}
	}

	// for each fileheader, get a handle to the actual file
	// https://gobyexample.com/errors
	file, err := files[0].Open()
	defer file.Close()
	if err != nil {
		return nil, "", nil, GeneralCodeError{
			fmt.Sprintf("Upload File: parsing form: ", err.Error()), 3}
	}

	// https://stackoverflow.com/questions/36111777/golang-how-to-read-a-text-file
	// https://stackoverflow.com/questions/30182538/why-can-not-i-copy-a-slice-with-copy-in-golang
	hasher := sha3.NewKeccak256()
	//dst4hash, err := ioutil.ReadFile(dstFName);
	dst4hash, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", nil, GeneralCodeError{
			fmt.Sprintf("Error in open file for hash: ", err.Error()), 3}
	}
	var dst4return []byte = nil
	if isCopyContent == true {
		//dst4return = make([]byte,len(dst4hash))
		//copy(dst4hash, dst4return)
		dst4return = dst4hash
	}
	hasher.Write(dst4hash)
	hashSum := hasher.Sum(nil)
	fmt.Printf("Hash: %x\n", hashSum)

	return hashSum, files[0].Filename, dst4return, GeneralCodeError{"OK", 0}
}

/*
  returns contrAddr, parentAddr, desc string of the client
*/
func ParseCert(binCert []byte) (common.Address, common.Address, string, error) {
	var contrAddr, parentAddr common.Address
	var retDesc string
	ca, err := x509.ParseCertificate(binCert)
	if err != nil {
		return common.Address{}, common.Address{}, "", err
	}

	for i := 0; i < len(ca.Subject.Names); i++ {
		retDesc += fmt.Sprint(ca.Subject.Names[i].Value) + " "
	}
	// iterate in the extension to get the information
	for _, element := range ca.Extensions {
		if element.Id.String() == "1.2.752.115.33.2" { // CA Address
			fmt.Printf("\tCaContractIdentifier: %+#+x\n", element.Value)
			val := element.Value[2:]
			if (len(val) != len(common.Address{}.Bytes())) {
				return common.Address{}, common.Address{}, "",
					GeneralError{"ParseCert: wrong length of CA addr"}
			}
			contrAddr = common.BytesToAddress(val)
		}
		if element.Id.String() == "1.2.752.115.33.1" { //Parent Address
			fmt.Printf("\tIssuerCaContractIdentifier: %+#+x\n", element.Value)
			val := element.Value[2:]
			if (len(val) != len(common.Address{}.Bytes())) {
				return common.Address{}, common.Address{}, "",
					GeneralError{"ParseCert: wrong length of CA addr"}
			}
			parentAddr = common.BytesToAddress(val)
		}
	}
	return contrAddr, parentAddr, retDesc, nil
}

/*
    params:
	  hash - hash to check in the currAddr,
	  		if hash==0, no user certificate check is conducted
             	and the info on the present CA is returned
      isGetCertData - if true the CA cert data is returned in the []byte variable
   returns:
	  indHashFound: - index of the hash available in this contract,
	  		or -1 if no hash is found
      revokeDate: date of revocation or zero if no revocation
      parentAddr,
      caHash of the present CA certificate
      []byte -- data of the CA certificate if the isGetCertData == true, otherwise nil
*/
func ConfirmHashCAData(currAddr common.Address, hash []byte,
	isGetCaCertData bool) (retIndHashFound int, retRevokeDate time.Time,
	retParentAddr common.Address, retCaHash []byte, retCertData []byte, err error) {

	client, err := ethclient.Dial(gConfig.IPCpath)
	if err != nil {
		return -1, time.Time{}, common.Address{}, nil, nil,
			GeneralError{fmt.Sprintf("Failed to connect to the Ethereum client: %v", err)}
	}

	return ConfirmHashCADataLight(client, currAddr, hash, isGetCaCertData)
}

/*
  the same like ConfirmHashCAData, but without connecting to Ethereum IPC
*/
func ConfirmHashCADataLight(client *ethclient.Client, currAddr common.Address, hash []byte,
	isGetCaCertData bool) (retIndHashFound int, retRevokeDate time.Time,
	retParentAddr common.Address, retCaHash []byte, retCertData []byte, err error) {

	retIndHashFound = -1
	retRevokeDate = time.Time{}
	retCertData = nil

	currContract, err := NewLuxUni_PKI(currAddr, client)
	if err != nil {
		return -1, time.Time{}, common.Address{}, nil, nil,
			GeneralError{fmt.Sprintf(
				"The address to smart contract is incorrect. Failed to instantiate a smart contract: %v", err)}
	}
	callOpts := &bind.CallOpts{
		Pending: true,
	}

	caCert, err := currContract.GetCaCertificate(callOpts)
	if err != nil {
		return -1, time.Time{}, common.Address{}, nil, nil, err
	}
	if (caCert == nil) && (len(caCert) == 0) {
		return -1, time.Time{}, common.Address{}, nil, nil,
			GeneralError{fmt.Sprintf(
				"The certificate of the smart contract is empty")}
	}
	if isGetCaCertData == true {
		retCertData = caCert
	}
	contrAddr, retParentAddr, _, err := ParseCert(caCert)
	if err != nil {
		return -1, time.Time{}, common.Address{}, nil, nil, err
	}
	if contrAddr != currAddr {
		return -1, time.Time{}, common.Address{}, nil, nil,
			GeneralError{"ConfirmCA: contrAddr does not correspond to the CA smart contract"}
	}
	retCaHash, err = CalcHash(caCert)
	if err != nil {
		return -1, time.Time{}, common.Address{}, nil, nil, err
	}

	if hash == nil {
		return -1, time.Time{}, retParentAddr, retCaHash, retCertData, nil
	}

	numRD, err := currContract.GetNumRegData(callOpts)
	if err != nil {
		return -1, time.Time{}, common.Address{}, nil, nil, err
	}

	for i := int64(0); i < numRD.Int64(); i++ {
		bi := big.NewInt(i)
		regDataHash, err := currContract.GetRegDataHash(callOpts, bi)
		if err != nil {
			return -1, time.Time{}, common.Address{}, nil, nil,
				GeneralError{fmt.Sprintf("Failed to get hash data: ", err)}
		}
		delRegDate, err := currContract.GetDeletedRegDate(callOpts, bi)
		if err != nil {
			return -1, time.Time{}, common.Address{}, nil, nil,
				GeneralError{fmt.Sprintf("Failed to get deleted data: ", err)}
		}
		if bytes.Equal(hash, regDataHash[:]) == true {
			if delRegDate.Int64() != 0 {
				retRevokeDate = time.Unix(delRegDate.Int64(), 0)
			} else {
				retRevokeDate = time.Time{}
				retIndHashFound = int(i)
				break
			}
		}
	}
	if (retRevokeDate != time.Time{}) {
		retIndHashFound = -1
	}

	return retIndHashFound, retRevokeDate, retParentAddr, retCaHash, retCertData, err
}

func CalcHash(data []byte) ([]byte, error) {
	hasher := sha3.NewKeccak256()
	hasher.Write(data)
	return hasher.Sum(nil), nil
}

func FindKeyFile(addr common.Address) (keyFileName string, err error) {
	files, err := ioutil.ReadDir(gConfig.KeyDir)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Name()), strings.ToLower(addr.String()[2:])) == true {
			return f.Name(), nil
		}
	}
	return "", GeneralError{"Ethereum Key File not found for this address"}
}

func LoadConfig() error {
	file, err := os.Open(gConfigFile)
	if err != nil {
		return GeneralError{fmt.Sprintf("File error: %v\n", err)}
	}
	fmt.Printf("Found configuration file: %s\n", gConfigFile)

	//jsonparser.Get(data, "person", "name", "fullName")
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&gConfig); err != nil {
		return GeneralError{fmt.Sprintf("Parsing config file: %s\n", err.Error())}
	}

	b, err := json.Marshal(gConfig)
	if err != nil {
		return GeneralError{fmt.Sprintf("Cannot convert conf file into string: %s", err)}
	}
	fmt.Printf("Loaded configuration file: %s\n", string(b))
	file.Close()

	if common.IsHexAddress(gConfig.AccountAddr) == false {
		return GeneralError{"Config: account address is not correct"}
	}
	if common.IsHexAddress(gConfig.ContractHash) == false {
		return GeneralError{"Config: contract hash is not correct"}
	}

	return nil
}

func LoadPrivateKey(path string, ptrPrivateKey *rsa.PrivateKey) error {
	privatekeyfile, err := os.Open(path)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(privatekeyfile)
	err = decoder.Decode(ptrPrivateKey)
	if err != nil {
		return err
	}
	privatekeyfile.Close()
	return nil
}

func Decrypt(msg []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		err := GeneralError{"Private keys is not loaded"}
		return nil, err
	}

	label := []byte("")
	hash := sha256.New()

	return rsa.DecryptOAEP(hash, cryptoRand.Reader, privateKey, msg, label)
}

type GeneralError struct {
	errMsg string
}

func (e GeneralError) Error() string {
	return e.errMsg
}

type GeneralCodeError struct {
	errMsg  string
	errCode uint64
}

func (e *GeneralCodeError) Error() string {
	return e.errMsg
}
