package main

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	//"math/rand"
)

func init() {

}

/*
the only command lne parameter:
   -config with the path to the config file
*/
func main() {
	ptrConfigFile := flag.String("config", gConfigFile, "the path to config file, default: "+gConfigFile)
	gConfigFile = *ptrConfigFile
	err := LoadConfig()
	if err != nil {
		fmt.Printf("CONFIG ERROR: %v\n", err)
		os.Exit(1)
	}
	gConfig.WebMode = 1

	if gConfig.PrivateKeyPath != "" {
		gPrivateKey = new(rsa.PrivateKey)
		err = LoadPrivateKey(gConfig.PrivateKeyPath, gPrivateKey)
		if err != nil {
			gPrivateKey = nil
			log.Printf("Private key is not loaded, omitting : %v\n", gConfig.PrivateKeyPath)
		} else {
			log.Printf("Private key successfully loaded : %v\n", gConfig.PrivateKeyPath)
		}
	}

	// how to connect to different servers
	//          (https://gist.github.com/bas-vk/299f4a686b66a22cf87302c561ee5866):
	//    geth --testnet --rpc
	// client, err := ethclient.Dial("http://localhost:8545")
	//    parity --testnet --port 31313 --jsonrpc-port 8546
	// client, err = ethclient.Dial("http://localhost:8546")

	// http://stackoverflow.com/questions/15834278/serving-static-content-with-a-root-url-with-the-gorilla-toolkit
	// subrouter - http://stackoverflow.com/questions/18720526/how-does-pathprefix-work-in-gorilla-mux-library-for-go
	r := mux.NewRouter()
	//r.HandleFunc("/pki-test", PkiForm);
	r.HandleFunc("/enroll_user", rstEnrollUser)
	r.HandleFunc("/blacklist_user", rstBlacklistUser) // uses UserID as a paramenter
	r.HandleFunc("/blacklist_hash", rstBlacklistHash) // uses Cetificate Hash as a parameter
	//r.HandleFunc("/enroll_ca", EnrollCA);
	r.HandleFunc("/create_contract", rstCreateContract)
	r.HandleFunc("/populate_contract", rstPopulateContract)
	//r.HandleFunc("/validate_form", ValidateForm);
	r.HandleFunc("/validate_cert", rstValidateCert)
	r.HandleFunc("/download_cacert", rstDownloadCaCert)
	//r.HandleFunc("/generate_user_cert", GenerateUserCert);

	fs := http.FileServer(http.Dir(gConfig.FileWebPath))
	spref := http.StripPrefix("/public/", fs)
	r.PathPrefix("/public/").Handler(spref)
	http.Handle("/", r)

	//https://gist.github.com/denji/12b3a568f092ab951456 - SSL info
	//https://golanglibs.com/top?q=webrtc - webrtc server side for golang

	//var server = &http.Server{
	//    Addr : ":8071",
	//    Handler : r,
	//}

	log.Println("RESTful service is listening...")
	//http.ListenAndServeTLS(":8071", "server.pem", "server.key", r)
	http.ListenAndServe(":"+strconv.Itoa(gConfig.RestHttpPort), nil)
}

/*
/blacklist_user, all parameters are in POST
Puts certificate (either ordinary or CA) from the white list to the black list
	Parameters:
		ParentAddr: the address of the CA smart contract where the certificate's hash is stored
		UserAddr: the ID (address) of the user who has the privilage to modify the smart contract.
			The key of this user should be available in key storage
		Deletion: array of strings with HASHES of the items to be deleted in the user list.
			It is produced with checkbox HTML forms
	Returns:
		200 and "OK" in the html body in case of success
		Errors (details are in html body):
			480 : Hash in Deletion array has wrong length or is not correct
			481 : Hash in Deletion array is not available in the white list
			482 : Hash in Deletion array is already revoked
			484 : ParentAddr is incorrect
			485 : Deletion array is incorrect
			580 : Ethereum executionn error (out of gas and others)
			581 : Ethereum connection error
			500 : Other error
*/
func rstBlacklistHash(w http.ResponseWriter, r *http.Request) {
	internalBlackList(w, r, true)
}

/*
/blacklist_user, all parameters are in POST
Puts certificate (either ordinary or CA) from the white list to the black list
	Parameters:
		ParentAddr: the address of the CA smart contract where the certificate's hash is stored
		UserAddr: the ID (address) of the user who has the privilage to modify the smart contract.
			The key of this user should be available in key storage
		Deletion: array of strings with IDs of the items to be deleted in the user list.
			It is produced with checkbox HTML forms
	Returns:
		200 and "OK" in the html body in case of success
		Errors (details are in html body):
			484 : ParentAddr is incorrect
			485 : Deletion is incorrect
			580 : Ethereum executionn error (out of gas and others)
			581 : Ethereum connection error
			500 : Other error
*/
func rstBlacklistUser(w http.ResponseWriter, r *http.Request) {
	internalBlackList(w, r, false)
}

/*
	rstBlackListUser and rstBlackListHash are shell functions to this one
	Please refer to rstBlackListUser for parameter and return details

	isHashes == true:
		dels contains hashes
	isHashes == false:
		dels contains IDs of the whitelist array

*/
func internalBlackList(w http.ResponseWriter, r *http.Request, isHashes bool) {
	var revokeResult string
	var parentAddr common.Address = common.Address{}
	var userAddr common.Address = common.Address{}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("No data: Parsing blacklist multipart form: %v\n", err.Error())
		http.Error(w, GeneralError{fmt.Sprintf(
			"BlacklistUser: error in parsing -- ", err.Error())}.Error(),
			http.StatusInternalServerError)
		return
	}

	strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
	if len(strParentAddrArr) > 0 {
		if common.IsHexAddress(strParentAddrArr[0]) == false {
			http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
				484 /*http.StatusInternalServerError*/)
			return
		}
		parentAddr = common.HexToAddress(strParentAddrArr[0])
	}

	if (parentAddr == common.Address{}) {
		http.Error(w, GeneralError{"Delete: Parent address is not established"}.Error(),
			484 /*http.StatusInternalServerError*/)
		return
	}

	strUserAddrArr := r.MultipartForm.Value["UserAddr"]
	if len(strUserAddrArr) > 0 {
		if common.IsHexAddress(strUserAddrArr[0]) == false {
			http.Error(w, GeneralError{"User address is incorrect"}.Error(),
				http.StatusInternalServerError)
			return
		}
		userAddr = common.HexToAddress(strUserAddrArr[0])
	}

	dels := r.MultipartForm.Value["Deletion"]
	if len(dels) > 0 {
		//fmt.Printf("Debug: I am in deletion block")
		//dels := r.MultipartForm.Value["Deletion"]
		//dels := r.Form["Deletion"]
		for _, del := range dels {
			fmt.Printf("Rest Debug: del=%v\n", del)
			var delid int

			if isHashes == true {
				if del[0:2] == "0x" || del[0:2] == "0X" {
					del = del[2:]
				}
				delHash, err := hex.DecodeString(del)
				if err != nil {
					http.Error(w, fmt.Sprintf("Hash conversion error: %v, hash: %s", err.Error(), del),
						480 /*http.StatusInternalServerError*/)
					return
				}
				indHashFound, revokeDate, _ /*superParentAddr*/, _ /*retCaHash []byte*/, _ /*retCertData []byte*/, err :=
					ConfirmHashCAData(parentAddr, delHash /*isGetCaCertData=*/, false)
				if err != nil {
					http.Error(w, fmt.Sprintf("Confirm hash error: %v", err.Error()),
						480 /*http.StatusInternalServerError*/)
					return
				}
				if indHashFound == -1 {
					if (revokeDate == time.Time{}) {
						http.Error(w, fmt.Sprintf("Revocation: the hash is not found in the white list"),
							481 /*http.StatusInternalServerError*/)
						return
						/*} else {
						http.Error(w, fmt.Sprintf("Revocation: the hash is already revoked"),
							482 )//http.StatusInternalServerError
						return*/
					}
				}
				delid = indHashFound
			} else {
				delid, err = strconv.Atoi(del)
				if err != nil {
					http.Error(w, fmt.Sprintf("Deletion conversion error: %v", err.Error()),
						485 /*http.StatusInternalServerError*/)
					return
				}
			}

			//revokedParam.RevokedIds = append(revokedParam.RevokedIds, delid);
			revokeResult += del + " "

			client, err := ethclient.Dial(gConfig.IPCpath)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to connect to the Ethereum client: %v", err),
					581 /*http.StatusInternalServerError*/)
				return
			}

			// Instantiate the contract, the address is taken from eth at the moment of contract initiation
			// kyc, err := NewLuxUni_KYC(common.HexToAddress(gContractHash), backends.NewRPCBackend(conn))
			pkiContract, err := NewLuxUni_PKI(parentAddr, client)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to instantiate a smart contract: %v", err),
					581 /*http.StatusInternalServerError*/)
				return
			}

			// Logging into Ethereum as a user
			if (userAddr == common.Address{}) {
				fmt.Printf("Attention! Revoke: user address is zero, default config account is used\n")
				userAddr = common.HexToAddress(gConfig.AccountAddr)
			}
			keyFile, err := FindKeyFile(userAddr)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to find key file for account %v. %v ",
					userAddr.String(), err), 581 /*http.StatusInternalServerError*/)
				return
			}
			key, err := ioutil.ReadFile(gConfig.KeyDir + keyFile)
			if err != nil {
				http.Error(w, fmt.Sprintf("Key File error: %v\n", err),
					581 /*http.StatusInternalServerError*/)
				return
			}

			auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to create authorized transactor: %v", err),
					581 /*http.StatusInternalServerError*/)
				return
			}

			sess := &LuxUni_PKISession{
				Contract: pkiContract,
				CallOpts: bind.CallOpts{
					Pending: true,
				},
				TransactOpts: bind.TransactOpts{
					From:     auth.From,
					Signer:   auth.Signer,
					GasLimit: big.NewInt(2000000),
				},
			}
			/* sess.TransactOpts = *auth
			sess.TransactOpts.GasLimit = big.NewInt(2000000) // Rinkeby block gas limit 6124970 */

			_, nerr := sess.DeleteRegDatum(big.NewInt(int64(delid)))
			if nerr != nil {
				http.Error(w, fmt.Sprintf("Deletion error: %v", nerr),
					580 /*http.StatusInternalServerError*/)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

/*
/enroll_user, all parameters in POST
	Parameters:
		Hash or UplFiles (hash isgConfig.Pswd a hex string without a leading "0x")
		UplFiles : uploaded certificate
		ParentAddr: the address of the CA smart contract where the certificate's hash is stored
			This address of this contract should be called at user account CurrentUserAddr
		CurrentUserAddr: the ID (address) of the user who has the privilage to modify the parent smart contract.
			The key of this user should be available in key storage
	Returns:
		200 and "OK" in the html body in case of success
		Errors (details are in html body):
			480 : hash has the wrong length or hash is incorrect
			481 : hash is already enrolled
			482 : Certificate errors in case it was provided instead of hash
			484 : ParentAddr is incorrect
			485 : CurrentUserAddr is incorrect
			580 : Ethereum execution error (out of gas and others)
			581 : Ethereum connection error
			500 : Other error
    Not used paremeters which were deleited and used in web application to store the data for CA tree navigation
    	// ContrAddr - REMOVED AS WEB APP STORES IT ITSELF
    	// NewUserAddr -- corresponds to userAddr associated with new contract -- REMOVED AS WEB APP STORES IT ITSELF
*/
func rstEnrollUser(w http.ResponseWriter, r *http.Request) {
	var parentAddr common.Address = common.Address{} // this is addr of the contract which is going to hold the hash
	// REMOVED - var contrAddr common.Address = common.Address{}   // this is address of the new SubCA contract or zero if end user
	var curUserAddr common.Address = common.Address{} // !! this is the user_id of the owner of parent contr
	// REMOVED var newUserAddr common.Address = common.Address{} // !! this is the new owner of contrAddr contr.
	var isNoUpload bool = false

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("EnrollUser: No change data -- ", err.Error())
		http.Error(w, GeneralError{fmt.Sprintf(
			"EnrollUser: No change data -- ", err.Error())}.Error(),
			http.StatusInternalServerError)
		return
	}

	hashSum, _ /*fileName*/, dataCert, cerr :=
		UploadFile(w, r, "UplFiles", true)
	if cerr.errCode != 0 {
		if cerr.errCode == 1 {
			isNoUpload = true
			hashArr := r.MultipartForm.Value["Hash"]
			if len(hashArr) == 0 {
				http.Error(w, GeneralError{fmt.Sprintf(
					"EnrollUser: No hashes in request")}.Error(),
					480 /* http.StatusInternalServerError */)
				return
			}

			hashStr := hashArr[0]
			if (hashStr[0:2] == "0x") || (hashStr[0:2] == "0X") {
				hashStr = hashStr[2:]
			}
			hashInt := big.NewInt(0)
			hashInt, _ = hashInt.SetString(hashStr, 16) /* tmpInt, err := strconv.Atoi(hashArr[0]); */
			if hashInt == nil {
				http.Error(w, fmt.Sprintf("EnrollUser: Hash string %s is incorrect", hashArr[0]),
					480 /*http.StatusInternalServerError*/)
				return
			}
			hashSum = hashInt.Bytes()
		} else {
			http.Error(w, GeneralError{fmt.Sprintf(
				"EnrollUser UplFiles:", cerr.Error())}.Error(),
				482 /*http.StatusInternalServerError*/)
			return
		}
	}

	strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
	if len(strParentAddrArr) > 0 {
		if common.IsHexAddress(strParentAddrArr[0]) == false {
			http.Error(w, fmt.Sprintf("Parent address as a parameter is incorrect: %v",
				strParentAddrArr[0]),
				484 /*http.StatusInternalServerError*/)
			return
		}
		parentAddr = common.HexToAddress(strParentAddrArr[0])
	}

	if isNoUpload == false {
		var caContrAddr, insertAddr common.Address
		caContrAddr, insertAddr, _ /*desc*/, err = ParseCert(dataCert)
		if err != nil {
			http.Error(w, fmt.Sprintf("CERTIFICATE: Parsing error: %v", err),
				482 /*http.StatusInternalServerError*/)
			return
		}
		if (insertAddr == common.Address{}) {
			http.Error(w, "CERTIFICATE: No Parent Address is provided in the Cert",
				482 /*http.StatusInternalServerError*/)
			return
		}
		if (caContrAddr != common.Address{}) {
			http.Error(w, "CERTIFICATE: Non-CA certificates should not include non-zero CA contract address",
				482 /*http.StatusInternalServerError*/)
			return
		}
		if insertAddr != parentAddr {
			http.Error(w, "Address in the certificate does not correspond to the contract address of the Authority (CA)",
				482 /*http.StatusInternalServerError*/)
			return
		}
	}

	strUserAddrArr := r.MultipartForm.Value["CurrentUserAddr"]
	if len(strUserAddrArr) > 0 {
		if common.IsHexAddress(strUserAddrArr[0]) == false {
			http.Error(w, GeneralError{"CurrentUser address is incorrect"}.Error(),
				485 /*http.StatusInternalServerError*/)
			return
		}
		curUserAddr = common.HexToAddress(strUserAddrArr[0])
	}

	if (parentAddr == common.Address{}) {
		http.Error(w, GeneralError{"Enroll: Parent address is not established"}.Error(),
			484 /*http.StatusInternalServerError*/)
		return
	}

	indHashFound, _ /*revokeDate*/, _ /*superParentAddr*/, _ /*retCaHash []byte*/, _ /*retCertData []byte*/, err :=
		ConfirmHashCAData(parentAddr, hashSum, false /*isGetCaCertData*/)
	if err != nil {
		http.Error(w, fmt.Sprintf("Confirm hash error: %v", err.Error()),
			480 /*http.StatusInternalServerError*/)
		return
	}
	if indHashFound != -1 {
		http.Error(w, fmt.Sprintf("Enroll: certificate already exists, hash: %x", hashSum),
			481 /*http.StatusInternalServerError*/)
		return
	}

	client, err := ethclient.Dial(gConfig.IPCpath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Enroll: Failed to connect to the Ethereum client: %v", err),
			581 /*http.StatusInternalServerError*/)
		return
	}

	// Instantiate the contract, the address is taken from eth at the moment of contract initiation
	// kyc, err := NewLuxUni_KYC(common.HexToAddress(gContractHash), backends.NewRPCBackend(conn))
	pkiContract, err := NewLuxUni_PKI(parentAddr, client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Enroll: Failed to instantiate a smart contract: %v", err),
			581 /*http.StatusInternalServerError*/)
		return
	}

	callOpts := &bind.CallOpts{
		Pending: true,
	}
	initNumRegData, err := pkiContract.GetNumRegData(callOpts)
	if err != nil {
		http.Error(w, fmt.Sprintf("EnrollUser: Failed to get numRegData from blockchain: %v. ", err),
			580 /*http.StatusInternalServerError*/)
		return
	}

	// Logging into Ethereum as a user
	if (curUserAddr == common.Address{}) {
		fmt.Printf("Attention! Enroll: user address is zero, default config account is used\n")
		curUserAddr = common.HexToAddress(gConfig.AccountAddr)
	}
	keyFile, err := FindKeyFile(curUserAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find key file for account %v. %v ",
			curUserAddr.String(), err), 581 /*http.StatusInternalServerError*/)
		return
	}
	key, err := ioutil.ReadFile(gConfig.KeyDir + keyFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Enroll: Ethereum connect -- Key File error: %v\n", err),
			581 /*http.StatusInternalServerError*/)
		return
	}
	//fmt.Printf("DEBUG: Found Ethereum Key File \n")

	auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Enroll: Failed to create authorized transactor: %v", err),
			581 /*http.StatusInternalServerError*/)
		return
	}

	sess := &LuxUni_PKISession{
		Contract: pkiContract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: big.NewInt(2000000),
		},
	}

	var tmpHash [32]byte
	copy(tmpHash[:], hashSum)
	res, err := sess.NewRegDatum(tmpHash, []byte("")) /* contrAddr, fileName, desc, "", newUserAddr */
	if err != nil {
		http.Error(w, fmt.Sprintf("EnrollUser: Failed to add a record to blockchain: %v. ", err),
			580 /*http.StatusInternalServerError*/)
		return
	}

	finalNumRegData, err := pkiContract.GetNumRegData(callOpts)
	if err != nil {
		http.Error(w, fmt.Sprintf("EnrollUser: Failed to get numRegData from blockchain: %v. ", err),
			580 /*http.StatusInternalServerError*/)
		return
	}

	if finalNumRegData.Int64() != initNumRegData.Int64()+1 {
		http.Error(w, fmt.Sprintf("EnrollUser: Failed to add a record, wrong function return: %х",
			res.Data()), 580 /*http.StatusInternalServerError*/)
		return
	}

	/*!!!!!*/ // var result uint64 = uint64(finalNumRegData.Int64() - 1)
	result, err := GetEventReturn(tmpHash, parentAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("EnrollUser: Failed to retreive result: %v",
			err), 580) //*http.StatusInternalServerError
		return
	}
	/*!!!!!*/ //result = uint64(finalNumRegData.Int64() - 1)

	if result != uint64(finalNumRegData.Int64()-1) {
		http.Error(w, fmt.Sprintf("EnrollUser: Retreived result does not correspond to Number of RegData",
			err), 580) //*http.StatusInternalServerError
		return
	}

	// UplFile is id in the input "file" component of the form
	// http://stackoverflow.com/questions/33771167/handle-file-uploading-with-go
	// file, handler, err := r.FormFile("UplFile")
	//out, err := os.Create("/tmp/tst_"+handler.Filename);

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(int(result))))
}

/*
/populate_contract, all parameters in POST
	Pupulation of the CA smart contract:
		a. putting a certificate into the contract referencing its parent, and
		b. setting ownership of the smartcontract to the user
	Params:
		UplFiles : uploaded certificate
		NewUserAddr - owner is set to this address at the end of the proc. If empty, then new owner is not set
			At the end of the population procedure only the NewUserAddr can modify the smart contract in the future
		CurrentUserAddr: - the user addr to connect to Ethereum. If empty, then set to root user addr
		ContrAddr: the address of the CA smart contract which should be populated
			This address of this contract should be called at user account CurrentUserAddr
	Returns:
		200 and hash string WITHOUT heading "0x" in the html body in case of success
		Errors (details are in html body):
			482 : Certificate errors
			483 : NewUserAddr is incorrect
			484 : ContrAddr is incorrect
			485 : CurrentUserAddr is incorrect
			580 : Ethereum execution error (out of gas and others)
			581 : Ethereum connection error
			500 : Other error
*/
func rstPopulateContract(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, GeneralError{"Rest populate: No change data: Parsing multipart form: %v"}.Error(),
			http.StatusInternalServerError)
		return
	}

	if len(r.MultipartForm.Value["ContrAddr"]) == 0 {
		http.Error(w, GeneralError{"Rest populate: No contrAddr is provided"}.Error(),
			484 /*http.StatusInternalServerError*/)
		return
	}
	contrAddrStr := r.MultipartForm.Value["ContrAddr"][0]
	if common.IsHexAddress(contrAddrStr) == false {
		http.Error(w, GeneralError{"Rest populate: Contract address is incorrect"}.Error(),
			484 /*http.StatusInternalServerError*/)
		return
	}
	contrAddr := common.HexToAddress(contrAddrStr)

	newUserAddr := common.Address{}
	if len(r.MultipartForm.Value["NewUserAddr"]) != 0 {
		userAddrStr := r.MultipartForm.Value["NewUserAddr"][0]
		if common.IsHexAddress(userAddrStr) == false {
			http.Error(w, GeneralError{"Rest populate: New User address is incorrect"}.Error(),
				483 /*http.StatusInternalServerError*/)
			return
		}
		newUserAddr = common.HexToAddress(userAddrStr)
	}

	curUserAddr := common.Address{}
	if len(r.MultipartForm.Value["CurrentUserAddr"]) != 0 {
		userAddrStr := r.MultipartForm.Value["CurrentUserAddr"][0]
		if common.IsHexAddress(userAddrStr) == false {
			http.Error(w, GeneralError{"Current User address is incorrect"}.Error(),
				http.StatusInternalServerError)
			return
		}
		curUserAddr = common.HexToAddress(userAddrStr)
	}

	hashCert, _, dataCert, cerr := UploadFile(w, r, "UplFiles", true)
	if cerr.errCode != 0 {
		fmt.Printf(fmt.Sprintf("Rest Populate: Uploadfile: %v\n", cerr.Error()))
		http.Error(w, cerr.Error(), 482 /*http.StatusInternalServerError*/)
		return
	}

	client, err := ethclient.Dial(gConfig.IPCpath)
	if err != nil {
		http.Error(w, err.Error(), 581 /*http.StatusInternalServerError*/)
		return
	}

	// Instantiate the contract, the address is taken from eth at the moment of contract initiation
	pkiContract, err := NewLuxUni_PKI(contrAddr, client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Populate: Failed to instantiate a smart contract: %v", err),
			581 /*http.StatusInternalServerError*/)
		return
	}

	// Logging into Ethereum as a user
	if (curUserAddr == common.Address{}) {
		fmt.Printf("Attention! Populate contract: user address is zero, default config account is used\n")
		curUserAddr = common.HexToAddress(gConfig.AccountAddr)
	}
	keyFile, err := FindKeyFile(curUserAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Populate: Failed to find key file for account %v. %v ",
			curUserAddr.String(), err), 581 /*http.StatusInternalServerError*/)
		return
	}
	key, err := ioutil.ReadFile(gConfig.KeyDir + keyFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Populatre: Key File %v error: %v\n",
			gConfig.KeyDir+keyFile, err), 581 /*http.StatusInternalServerError*/)
		return
	}
	fmt.Printf("Found Ethereum Key File \n")

	auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
		http.Error(w, fmt.Sprintf("Rest Populatre: Failed to create authorized transactor: %v", err),
			581 /*http.StatusInternalServerError*/)
		return
	}

	sess := &LuxUni_PKISession{
		Contract: pkiContract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: big.NewInt(2000000),
		},
	}

	_, err = sess.PopulateCertificate(dataCert)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Populate: Failed to populate blockchain: %v", err),
			580 /*http.StatusInternalServerError*/)
		return
	}
	if (newUserAddr != common.Address{}) {
		_, err := sess.SetOwner(newUserAddr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Rest Populate: Failed to update owner addr: %v", err),
				580 /*http.StatusInternalServerError*/)
			return
		}
		newOwner, err := sess.GetOwner()
		if err != nil {
			http.Error(w, fmt.Sprintf("Rest Populate: Failed to check new owner addr: %v", err),
				580 /*http.StatusInternalServerError*/)
			return
		}
		if newOwner != newUserAddr {
			http.Error(w, fmt.Sprintf("OwnerAddr (%v) does not equal to newUserAddr (%v) despite SetOwner - probably lack of permissions",
				newOwner.String(), newUserAddr.String()), 580 /*http.StatusInternalServerError*/)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hex.EncodeToString(hashCert)))
}

/*
/create_contract, all params as POST
	Creation of the "empty" CA smart contract:
		a. CA certificate should be added to smart contract through population procedure
		b. the right to execute the smart contract should be changed to the CA account with population procedure as well
	Params:
		ParentAddr: the address of the CA smart contract which is used for creation (it has the bin code)
			This address of this contract should be called at user account CurrentUserAddr
		NewUserAddr - owner is set to this address at the end of the proc. If empty, then new owner is not set
			At the end of the population procedure only the NewUserAddr can modify the smart contract in the future
		CurrentUserAddr: - the user addr to connect to Ethereum. If empty, then set to root user addr
	Returns:
		200 and the smart contract address WITH heading "0x" in the html body in case of success
		Errors (details are in html body):
			480 : Current user does not have rights to execute the creation of the CA certificate
			483 : NewUserAddr is incorrect
			484 : ParentAddr is incorrect
			485 : CurrentUserAddr is incorrect
			580 : Ethereum execution error (out of gas and others)
			581 : Ethereum connection error
			500 : Other error
*/
func rstCreateContract(w http.ResponseWriter, r *http.Request) {
	/*
	   https://vincentserpoul.github.io/post/binding-ethereum-golang/
	   https://ethereum.stackexchange.com/questions/7499/how-are-addresses-created-if-deploying-a-new-bound-contract
	*/
	var parentAddrStr string
	var curUserAddrStr string // !!! presently current user not used - addr=contr.GetOwner used instead
	var newUserAddrStr string

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("CreateContract: No data in multipart form: %v\n", err.Error())
		parentAddrStr = gConfig.ContractHash
	} else {
		strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
		if len(strParentAddrArr) > 0 {
			parentAddrStr = strParentAddrArr[0]
			if common.IsHexAddress(parentAddrStr) == false {
				fmt.Println("Create Contract: Parent address is incorrect")
				http.Error(w, GeneralError{"Rest Create contract: Parent address is incorrect"}.Error(),
					484 /*http.StatusInternalServerError*/)
			}
		} else {
			parentAddrStr = gConfig.ContractHash
		}
	}

	// !!! presently current user not used - addr=contr.GetOwner used instead
	strUserAddrArr := r.MultipartForm.Value["CurrentUserAddr"]
	if len(strUserAddrArr) > 0 {
		curUserAddrStr = strUserAddrArr[0]
		if common.IsHexAddress(curUserAddrStr) == false {
			fmt.Println("Create Contract: Current user address is incorrect")
			http.Error(w, GeneralError{"Rest Create contract: Current user address is incorrect"}.Error(),
				485 /*http.StatusInternalServerError*/)
			return
		}
	}

	strUserAddrArr = r.MultipartForm.Value["NewUserAddr"]
	if len(strUserAddrArr) > 0 {
		newUserAddrStr = strUserAddrArr[0]
		if common.IsHexAddress(newUserAddrStr) == false {
			fmt.Println("Create Contract: New user address is incorrect")
			http.Error(w, "Rest Create contreact: New user address is incorrect", 483 /*http.StatusInternalServerError*/)
			return
		}
	} else {
		http.Error(w, "New user address is not available", 483 /*http.StatusInternalServerError*/)
		return
	}

	client, err := ethclient.Dial(gConfig.IPCpath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Create contract: Failed to connect to the Ethereum client: %v", err),
			581 /*http.StatusInternalServerError*/)
		return
	}

	pkiContr, err := NewLuxUni_PKI(common.HexToAddress(parentAddrStr), client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Create contract: Failed to instantiate a smart contract: %v", err),
			581 /*http.StatusInternalServerError*/)
		return
	}
	callOpts := &bind.CallOpts{
		Pending: true,
	}
	execUserAddr, err := pkiContr.GetOwner(callOpts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Create contr - failed to get owner addr: ", err),
			581 /*http.StatusInternalServerError*/)
		return
	}
	if execUserAddr != common.HexToAddress(curUserAddrStr) {
		http.Error(w, "Rest Create contract: GetOwner does not correspond to the Current User param",
			480 /*http.StatusInternalServerError*/)
		return
	}

	keyFile, err := FindKeyFile(execUserAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Create contract: FindKeyFile: %v. ", err),
			581 /*http.StatusInternalServerError*/)
		return
	}
	key, err := ioutil.ReadFile(gConfig.KeyDir + keyFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Create contract: Key File error: %v", err),
			581 /*http.StatusInternalServerError*/)
		return
	}

	auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Create contract: Failed to create authorized transactor: %v", err),
			581 /*http.StatusInternalServerError*/)
		return
	}
	var trOpts bind.TransactOpts = *auth
	trOpts.GasLimit = big.NewInt(4000000) // 6124970 - block gas limit in Rinkeby
	contrAddr, _ /*contr*/, _, err := DeployLuxUni_PKI(&trOpts, client)
	/*
	   https://stackoverflow.com/questions/40096750/set-status-code-on-http-responsewriter
	*/
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Create contract: CreateContract -- Etherreum error in contract creation: %v", err),
			580 /*http.StatusInternalServerError*/)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(contrAddr.String()))
}

/*
/download_cacert
	Extracting (download) of certificate from CA smart contract
	Params:
		ContrAddr: the address of the CA smart contract
	Returns:
		200 and the smart contract address WITH heading "0x" in the html body in case of success
		Errors (details are in html body):
			484 : ContrAddr is incorrect
			580 : Ethereum execution error (out of gas and others)
			581 : Ethereum connection error
			500 : Other error
	  https://stackoverflow.com/questions/35496233/go-how-to-i-make-download-service
	  https://play.golang.org/p/UMKgI_NLwO
*/
func rstDownloadCaCert(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
		return
	}

	if len(r.MultipartForm.Value["ContrAddr"]) == 0 {
		http.Error(w, GeneralError{"No contrAddr is provided"}.Error(),
			484 /*http.StatusInternalServerError*/)
		return
	}
	strContrAddr := r.MultipartForm.Value["ContrAddr"][0]
	if common.IsHexAddress(strContrAddr) == false {
		http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
			484 /*http.StatusInternalServerError*/)
		return
	}
	contrAddr := common.HexToAddress(strContrAddr)

	/*isCertOK*/ _ /*revokDate*/, _ /*parentAddr*/, _ /*retCaHash*/, _, certData, err :=
		ConfirmHashCAData(contrAddr, nil, true)

	w.Header().Set("Content-Disposition", "attachment; filename=ca.crt.out")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	http.ServeContent(w, r, "ca.crt.out", time.Now(), bytes.NewReader(certData))
}

/*
/validate_cert, all params as POST
	Parameters:
		Hash or UplFiles (hash is a hex string without a leading "0x")
		UplFiles : uploaded certificate
		ParentAddr: the address of the CA smart contract where the certificate's hash is stored
			If certificate is uploaded through UplFiles, ParentAddr may not be specified
	Returns:
		200 and JSON with the validation results in the html body in case of success
		Errors (details are in html body):
			480 : hash has wrong length or hash is incorrect
			482 : Certificate errors in case it was provided instead of hash
			484 : ParentAddr is incorrect
			580 : Ethereum execution error (out of gas and others)
			581 : Ethereum connection error
			500 : Other error
*/
func rstValidateCert(w http.ResponseWriter, r *http.Request) {

	var parentAddr common.Address
	var isNoUpload bool = false

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Velidat: No change data: Parsing multipart form: %v", err),
			http.StatusInternalServerError)
		return
	}

	certHash, _ /*fileName*/, dataCert, cerr := UploadFile(w, r, "UplFiles", true)
	if cerr.errCode != 0 {
		if cerr.errCode == 1 {
			isNoUpload = true
			hashArr := r.MultipartForm.Value["Hash"]
			if len(hashArr) == 0 {
				http.Error(w, fmt.Sprintf("Rest Validate: No hashes in request"),
					480 /* http.StatusInternalServerError */)
				return
			}

			hashStr := hashArr[0]
			if (hashStr[0:2] == "0x") || (hashStr[0:2] == "0X") {
				hashStr = hashStr[2:]
			}
			hashInt := big.NewInt(0)
			hashInt, _ = hashInt.SetString(hashStr, 16) /* tmpInt, err := strconv.Atoi(hashArr[0]); */
			if hashInt == nil {
				http.Error(w, fmt.Sprintf("Rest Validate: Hash string %s is incorrect", hashArr[0]),
					480 /*http.StatusInternalServerError*/)
				return
			}
			certHash = hashInt.Bytes()
		} else {
			http.Error(w, fmt.Sprintf("Rest Validate UplFiles: %v", cerr.Error()),
				482 /*http.StatusInternalServerError*/)
			return
		}
	}

	strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
	if len(strParentAddrArr) > 0 {
		if common.IsHexAddress(strParentAddrArr[0]) == false {
			http.Error(w, fmt.Sprintf("Rest Validate: Parent address as a parameter is incorrect: %v",
				strParentAddrArr[0]),
				484 /*http.StatusInternalServerError*/)
			return
		}
		parentAddr = common.HexToAddress(strParentAddrArr[0])
	}

	if isNoUpload == false {
		var insertAddr common.Address
		_ /*caContrAddr*/, insertAddr, _ /*desc*/, err = ParseCert(dataCert)
		if err != nil {
			http.Error(w, fmt.Sprintf("Rest CERTIFICATE: Parsing error: %v", err),
				482 /*http.StatusInternalServerError*/)
			return
		}
		if (insertAddr == common.Address{}) {
			http.Error(w, "Rest CERTIFICATE: No Parent Address is provided in the Cert",
				482 /*http.StatusInternalServerError*/)
			return
		}
		if (parentAddr != common.Address{}) {
			if insertAddr != parentAddr {
				http.Error(w, "Rest Validate: Address in the certificate does not correspond to the contract address of the Authority (CA)",
					482 /*http.StatusInternalServerError*/)
				return
			}
		} else {
			parentAddr = insertAddr
		}
	}
	if (parentAddr == common.Address{}) {
		http.Error(w, "Rest Validate: Parent Address is not defined",
			484 /*http.StatusInternalServerError*/)
		return
	}

	start := time.Now()
	isCertOK, revokeDate, _ /*certPath*/, iter, err := CheckCertTree(parentAddr, certHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	finish := time.Now()
	elapsed := finish.Sub(start)

	var jsonResponse JsonValidateResponse
	jsonResponse.ProcMilisec = elapsed.Nanoseconds() / 1000
	jsonResponse.RevokeDate = revokeDate
	jsonResponse.IsCertOK = isCertOK
	jsonResponse.Iter = iter
	//jsonResponse.CertPath = certPath
	jsonResponse.Status = 0
	bJson, err := json.Marshal(jsonResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rest Validate: Json Marshal:", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bJson)
}

/*
  returns Json string with the path of certificates to the root
          int with the number of iteractions
*/
func CheckCertTree(parentAddr common.Address, userHash []byte) (retIsCertOK bool,
	retRevokeDate time.Time, retCertPath []JsonValidateNode, retIter int, err error) {

	var maxIter int = 200000

	client, err := ethclient.Dial(gConfig.IPCpath)
	if err != nil {
		return false, time.Time{}, nil, retIter,
			GeneralError{fmt.Sprintf("Failed to connect to the Ethereum client: %v", err)}
	}

	iterHash := userHash
	for retIter = 0; retIter < maxIter; retIter++ {
		var indCertOK int
		var jsonNode JsonValidateNode
		jsonNode.ContrAddr = parentAddr.String()
		jsonNode.Hash = fmt.Sprintf("%x", iterHash)

		indCertOK, retRevokeDate, parentAddr, iterHash, _, err =
			ConfirmHashCADataLight(client, parentAddr, iterHash, false)
		if err != nil {
			return false, time.Time{}, nil, retIter, err
		}
		if indCertOK > -1 {
			retIsCertOK = true
		} else {
			retIsCertOK = false
		}

		jsonNode.ParentAddr = parentAddr.String()
		jsonNode.IsCertOK = strconv.FormatBool(retIsCertOK)
		jsonNode.RevokeDate = retRevokeDate.String()
		retCertPath = append(retCertPath, jsonNode)
		if retIsCertOK == false {
			break
		}
		if (parentAddr == common.Address{}) {
			break
		}
		if (retIter >= (maxIter - 1)) && (parentAddr != common.Address{}) {
			return false, time.Time{}, nil, retIter, GeneralError{"MaxIter limit is reached"}
		}
	}
	return retIsCertOK, retRevokeDate, retCertPath, retIter, nil
}

/*
	Getting results through Events
		the code of event evLuxUni_NewRegDatumReturn(uint256,uint256)
		keccak - web3.sha3("evLuxUni_NewRegDatumReturn(uint256,uint256)") --
				"0x75d1c4f2937517b1233bf95d8f3b4c1d077820b9bc4c5bc28adcd886a3ba7ab6"
				 -- this is the topics in creation of the filter for event logs
*/
func GetEventReturn(dataHash [32]byte, contrAddr common.Address) (result uint64, err error) {

	var evCallHash common.Hash
	bytesCallHash, _ := hex.DecodeString("75d1c4f2937517b1233bf95d8f3b4c1d077820b9bc4c5bc28adcd886a3ba7ab6")
	evCallHash.SetBytes(bytesCallHash)
	fmt.Printf("DEBUG: callHash: %s\n", evCallHash.Hex())

	query := ethereum.FilterQuery{
		FromBlock: nil,
		ToBlock:   nil,
		Topics:    [][]common.Hash{{evCallHash}}, //[][]common.Hash
		Addresses: []common.Address{contrAddr}}
	var logs = make(chan types.Log) //, 2)

	client, err := ethclient.Dial(gConfig.IPCpath)
	if err != nil {
		return 0, GeneralError{fmt.Sprintf("Failed to connect to the Ethereum client: %v", err)}
	}
	s, err := client.SubscribeFilterLogs(context.TODO(), query, logs)
	if err != nil {
		return 0, GeneralError{fmt.Sprintf("Failed to establish Ethereum event filter: %v", err)}
	}

	errChan := s.Err()
	for {
		select {
		case err := <-errChan:
			return 0, GeneralError{fmt.Sprintf("Event Logs subscription error: %v", err)}
		case l := <-logs:
			//fmt.Printf("DEBUG Event Data: %x\n", l.Data)
			return ProcEventInteger(l.Data, dataHash)
		}
	}

}

func ProcEventInteger(evData []byte, dataHash [32]byte) (result uint64, err error) {

	if len(evData) != 32*2 {
		return 0, GeneralError{fmt.Sprintf("Length of the data string is not valid: %v, dataString: %x, dataHash: %x",
			len(evData), evData, dataHash)}
	}

	eventHash := evData[:32]

	if bytes.Equal(eventHash, dataHash[:]) == true {
		eventReturn := binary.BigEndian.Uint64(evData[(32*2 - 64/8) : 32*2])
		return eventReturn, nil
	}
	return 0, GeneralError{fmt.Sprintf("Hash %x is not found in data string %x",
		dataHash, evData)}
}

func CallRPC(query string) ([]byte, error) {
	body := strings.NewReader(query)
	req, err := http.NewRequest("POST", gConfig.EthereumRpcUrl+":"+strconv.Itoa(gConfig.EthereumRpcPort), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
