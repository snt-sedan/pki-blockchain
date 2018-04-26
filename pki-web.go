package main

import (
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	//"math/rand"
	"bytes"
	"io"
	"mime/multipart"
)

var mTempl map[string]*template.Template

type CRegData struct {
	EthAccCA     common.Address
	ContrAddr    common.Address
	DataHash     []byte
	FileName     string
	Description  string
	Encrypted    *big.Int
	CryptoModule common.Address
	CreationDate *big.Int
	Active       bool
}

// alternative aproach to passing the params to template - http://stackoverflow.com/questions/23802008/how-do-you-pass-multiple-objects-to-go-template
// general Go template description https://gohugo.io/templates/go-templates/
type CPKIFormParam struct {
	RandomNum       int
	RandomStr       string
	ParentAddr      string
	SuperParentAddr string
	SuperUserAddr   string
	IsRoot          bool
	Docs            []CDocPrez
	//AddrCAs      []string
}

type CDocPrez struct {
	Id           int
	Name         string
	ParentAddr   string
	ContrAddr    string
	UserAddr     string
	Desc         string
	Link         string
	Decryption   string
	Hash         string
	CreationDate time.Time
	CreationStr  string
	IsCA         bool
}

type CRevokedFormParam struct {
	RevokedIds []int
	RevokedStr string
	Docs       []CDocPrez
}

type CSimpleFormParam struct {
	Result           string
	Data             string
	Params           []CParam
	CallForm         string
	IsUploadFileForm bool
}

type CParam struct {
	Name  string
	Value string
}

func init() {

	if mTempl == nil {
		mTempl = make(map[string]*template.Template)
	}

	///home/alex/DocsFS/Dropbox/WORK/RD/LuxBCh/PKI
	t, e := template.ParseFiles("./html/form_main.html")
	if e != nil {
		fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "main", e.Error())
		mTempl = nil
	} else {
		mTempl["MainForm"] = template.Must(t, e)
	}

	t, e = template.ParseFiles("./html/form_hashres.html")
	if e != nil {
		fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "hashes", e.Error())
		mTempl = nil
	} else {
		mTempl["HashResult"] = template.Must(t, e)
	}

	t, e = template.ParseFiles("./html/form_simple.html")
	if e != nil {
		fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "simple", e.Error())
		mTempl = nil
	} else {
		mTempl["SimpleForm"] = template.Must(t, e)
	}

}

func main() {

	if mTempl == nil {
		return
	}

	ptrConfigFile := flag.String("config", gConfigFile, "the path to config file, default: "+gConfigFile)
	gConfigFile = *ptrConfigFile
	err := LoadConfig()
	if err != nil {
		fmt.Printf("CONFIG ERROR: %v\n", err)
		os.Exit(1)
	}
	gConfig.WebMode = 1

	gPrivateKey = new(rsa.PrivateKey)
	err = LoadPrivateKey(gConfig.PrivateKeyPath, gPrivateKey)
	if err != nil {
		gPrivateKey = nil
		log.Printf("Private key is not loaded, omitting : %v\n", gConfig.PrivateKeyPath)
	} else {
		log.Printf("Private key successfully loaded : %v\n", gConfig.PrivateKeyPath)
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
	r.HandleFunc("/pki-test", PkiForm)
	r.HandleFunc("/enroll_user", EnrollUser)
	r.HandleFunc("/blacklist_user", BlacklistUser)
	//r.HandleFunc("/enroll_ca", EnrollCA);
	r.HandleFunc("/create_contract", CreateContract)
	r.HandleFunc("/populate_contract", PopulateContract)
	r.HandleFunc("/validate_form", ValidateForm)
	r.HandleFunc("/validate_cert", ValidateCert)
	r.HandleFunc("/download_cacert", DownloadCaCert)
	r.HandleFunc("/generate_user_cert", GenerateUserCert)

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

	log.Println("Listening...")
	//http.ListenAndServeTLS(":8071", "server.pem", "server.key", r)
	http.ListenAndServe(":"+strconv.Itoa(gConfig.AppHttpPort), nil)
}

/*

 */
func BlacklistUser(w http.ResponseWriter, r *http.Request) {
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

	var b bytes.Buffer
	newWri := multipart.NewWriter(&b)

	strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
	if len(strParentAddrArr) > 0 {
		if common.IsHexAddress(strParentAddrArr[0]) == false {
			http.Error(w, GeneralError{"Parent address is incorrect"}.Error(),
				http.StatusInternalServerError)
			return
		}
		parentAddr = common.HexToAddress(strParentAddrArr[0])
	}
	if (parentAddr == common.Address{}) {
		http.Error(w, "Delete: Parent address is not found",
			http.StatusInternalServerError)
		return
	}
	fw, err := newWri.CreateFormField("ParentAddr")
	if err != nil {
		http.Error(w, "Delete: Parent address is not created in request to REST",
			http.StatusInternalServerError)
		return
	}
	if _, err = fw.Write([]byte(parentAddr.String())); err != nil {
		http.Error(w, "Delete: Parent address is not written in request to REST",
			http.StatusInternalServerError)
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

		fw, err := newWri.CreateFormField("UserAddr")
		if err != nil {
			http.Error(w, "Delete: User address is not created in request to REST",
				http.StatusInternalServerError)
			return
		}
		if _, err = fw.Write([]byte(userAddr.String())); err != nil {
			http.Error(w, "Delete: User address is not written in request to REST",
				http.StatusInternalServerError)
			return
		}
	}

	dels := r.MultipartForm.Value["Deletion"]
	if len(dels) > 0 {
		if len(dels) > 1 {
			http.Error(w, "Delete: Hey, I dont know how to prepare mutiple DELs in request to REST",
				http.StatusInternalServerError)
			return
		}
		fw, err := newWri.CreateFormField("Deletion")
		if err != nil {
			http.Error(w, "Delete: Deletion field is not created in request to REST",
				http.StatusInternalServerError)
			return
		}
		for _, del := range dels {
			fmt.Printf("del=%v\n", del)
			if _, err = fw.Write([]byte(del)); err != nil {
				http.Error(w, "Delete: Deletion field is not written in request to REST",
					http.StatusInternalServerError)
				return
			}
		}
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	newWri.Close()
	// Now that you have a form, you can submit it to your handler.
	newReq, err := http.NewRequest("POST",
		gConfig.RestUrlServer+":"+strconv.Itoa(gConfig.RestHttpPort)+"/blacklist_hash", &b) // /blacklist_user to send IDs
	if err != nil {
		http.Error(w, fmt.Sprintf("Delete: error in creation of REST request - %v", err),
			http.StatusInternalServerError)
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	newReq.Header.Set("Content-Type", newWri.FormDataContentType())

	// Submit the request
	// https://stackoverflow.com/questions/20205796/golang-post-data-using-the-content-type-multipart-form-data
	client := &http.Client{}
	res, err := client.Do(newReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Delete: error in REST request - %v", err),
			http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		strResult, _ := ioutil.ReadAll(res.Body)
		http.Error(w, fmt.Sprintf("Delete: REST's bad status: %v, %s", res.Status, strResult),
			http.StatusInternalServerError)
		return
	}

	if gConfig.WebMode != 0 {
		var formParam CSimpleFormParam
		//formParam.Param = ""
		formParam.Params = append(formParam.Params,
			CParam{"ParentAddr", parentAddr.String()})
		formParam.Result = "CA " + parentAddr.String() + ": certificate was successfully revoked"
		formParam.CallForm = "pki-test"

		tmpl := mTempl["SimpleForm"]
		terr := tmpl.Execute(w, formParam)
		if terr != nil {
			http.Error(w, terr.Error(), http.StatusInternalServerError)
		}
	}
}

/*
post paremeters:
Hash or Files (hash is a hex string WITHOUT leading "0x")
ParentAddr
ContrAddr
CurrentUserAddr -- corresponds to userAddr of a user holding the parent contract
NewUserAddr -- corresponds to userAddr associated with new contract
*/
func EnrollUser(w http.ResponseWriter, r *http.Request) {
	//var isCurl []string;
	var parentAddr common.Address = common.Address{}  // this is addr of the contract which is going to hold the hash
	var contrAddr common.Address = common.Address{}   // this is address of the new SubCA contract or zero if end user
	var curUserAddr common.Address = common.Address{} // !! this is the user_id of the owner of parent contr
	var newUserAddr common.Address = common.Address{} // !! this is the new owner of contrAddr contr.
	var isNoUpload bool = false

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("EnrollUser: No change data -- %v\n", err.Error())
		http.Error(w, GeneralError{fmt.Sprintf(
			"EnrollUser: No change data -- %v", err.Error())}.Error(),
			http.StatusInternalServerError)
		return
	}

	var b bytes.Buffer
	newWri := multipart.NewWriter(&b)

	//isCurl = r.MultipartForm.Value["Curl"]
	hashSum, fileName, dataCert, cerr :=
		UploadFile(w, r, "UplFiles", true)
	if cerr.errCode != 0 {
		if cerr.errCode == 1 {
			isNoUpload = true
			hashArr := r.MultipartForm.Value["Hash"]
			if len(hashArr) == 0 {
				http.Error(w, GeneralError{fmt.Sprintf(
					"EnrollUser: No hashes in request")}.Error(),
					http.StatusInternalServerError)
				return
			}

			hashStr := hashArr[0]
			if (hashStr[0:2] == "0x") || (hashStr[0:2] == "0X") {
				hashStr = hashStr[2:]
			}
			hashInt := big.NewInt(0)
			hashInt, _ = hashInt.SetString(hashStr, 16) /* tmpInt, err := strconv.Atoi(hashArr[0]); */
			if hashInt == nil {
				http.Error(w, fmt.Sprintf("EnrollUser: Hash string %v is incorrect", hashArr[0]),
					http.StatusInternalServerError)
				return
			}
			hashSum = hashInt.Bytes()

			fw, err := newWri.CreateFormField("Hash")
			if err != nil {
				http.Error(w, "EnrollUser: Hash is not created in request to REST",
					http.StatusInternalServerError)
				return
			}
			//fmt.Printf("Debug Enroll User initiation: hash to send = %s", hex.EncodeToString(hashSum))
			if _, err = fw.Write([]byte(hex.EncodeToString(hashSum))); err != nil {
				http.Error(w, "EnrollUser: Hash is not written in request to REST",
					http.StatusInternalServerError)
				return
			}

		} else {
			http.Error(w, GeneralError{fmt.Sprintf(
				"EnrollUser UplFiles: %v", cerr.Error())}.Error(),
				http.StatusInternalServerError)
			return
		}
	}

	strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
	if len(strParentAddrArr) > 0 {
		if common.IsHexAddress(strParentAddrArr[0]) == false {
			http.Error(w, "Parent address as a parameter is incorrect",
				http.StatusInternalServerError)
			return
		}
		parentAddr = common.HexToAddress(strParentAddrArr[0])

		fw, err := newWri.CreateFormField("ParentAddr")
		if err != nil {
			http.Error(w, "EnrollUser: ParentAddr is not created in request to REST",
				http.StatusInternalServerError)
			return
		}
		if _, err = fw.Write([]byte(parentAddr.String())); err != nil {
			http.Error(w, "EnrollUser: ParentAddr is not written in request to REST",
				http.StatusInternalServerError)
			return
		}
	}

	if isNoUpload == false {
		fw, err := newWri.CreateFormFile("UplFiles", fileName)
		if err != nil {
			http.Error(w, fmt.Sprintf("EnrollUser: cannot create form file: %v", err),
				http.StatusInternalServerError)
			return
		}
		if _, err = io.Copy(fw, bytes.NewReader(dataCert)); err != nil {
			http.Error(w, fmt.Sprintf("EnrollUser: cannot read/copy cert data: %v", err),
				http.StatusInternalServerError)
			return
		}
	}

	strUserAddrArr := r.MultipartForm.Value["CurrentUserAddr"]
	if len(strUserAddrArr) > 0 {
		if common.IsHexAddress(strUserAddrArr[0]) == false {
			http.Error(w, GeneralError{"CurrentUser address is incorrect"}.Error(),
				http.StatusInternalServerError)
			return
		}
		curUserAddr = common.HexToAddress(strUserAddrArr[0])

		fw, err := newWri.CreateFormField("CurrentUserAddr")
		if err != nil {
			http.Error(w, "EnrollUser: CurrentUserAddr is not created in request to REST",
				http.StatusInternalServerError)
			return
		}
		if _, err = fw.Write([]byte(curUserAddr.String())); err != nil {
			http.Error(w, "EnrollUser: CurrentUserAddr is not written in request to REST",
				http.StatusInternalServerError)
			return
		}
	}

	strUserAddrArr = r.MultipartForm.Value["NewUserAddr"]
	if len(strUserAddrArr) > 0 {
		if common.IsHexAddress(strUserAddrArr[0]) == false {
			http.Error(w, GeneralError{"NewUser address is incorrect"}.Error(),
				http.StatusInternalServerError)
			return
		}
		newUserAddr = common.HexToAddress(strUserAddrArr[0])

		fw, err := newWri.CreateFormField("NewUserAddr")
		if err != nil {
			http.Error(w, "EnrollUser: NewUserAddr is not created in request to REST",
				http.StatusInternalServerError)
			return
		}
		if _, err = fw.Write([]byte(newUserAddr.String())); err != nil {
			http.Error(w, "EnrollUser: NewUserAddr is not written in request to REST",
				http.StatusInternalServerError)
			return
		}
	}

	//if fileName==gCaCertName {
	strContrAddrArr := r.MultipartForm.Value["ContrAddr"]
	if len(strContrAddrArr) > 0 {
		if common.IsHexAddress(strContrAddrArr[0]) == false {
			http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
				http.StatusInternalServerError)
			return
		}
		contrAddr = common.HexToAddress(strContrAddrArr[0])

		fw, err := newWri.CreateFormField("ContrAddr")
		if err != nil {
			http.Error(w, "EnrollUser: ContrAddr is not created in request to REST",
				http.StatusInternalServerError)
			return
		}
		if _, err = fw.Write([]byte(contrAddr.String())); err != nil {
			http.Error(w, "EnrollUser: ContrAddr is not written in request to REST",
				http.StatusInternalServerError)
			return
		}
	}

	if (parentAddr == common.Address{}) {
		http.Error(w, GeneralError{"Enroll: Parent address is not established"}.Error(),
			http.StatusInternalServerError)
		return
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	newWri.Close()
	// Now that you have a form, you can submit it to your handler.
	newReq, err := http.NewRequest("POST",
		gConfig.RestUrlServer+":"+strconv.Itoa(gConfig.RestHttpPort)+"/enroll_user", &b)
	if err != nil {
		http.Error(w, fmt.Sprintf("EnrollUser: error in creation of REST request - %v", err),
			http.StatusInternalServerError)
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	newReq.Header.Set("Content-Type", newWri.FormDataContentType())

	// Submit the request
	// https://stackoverflow.com/questions/20205796/golang-post-data-using-the-content-type-multipart-form-data
	client := &http.Client{}
	res, err := client.Do(newReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("EnrollUser: error in REST request - %v", err),
			http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	restResult, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("EnrollUser: REST's bad status: %v, %v", res.Status, string(restResult)),
			http.StatusInternalServerError)
		return
	}

	//fmt.Printf("DEBUG Received Rest result: %v", string(restResult))
	var jsonRestRes struct {
		ArrayInd int `json:"arrayInd"`
	}
	jsonRestRes.ArrayInd, err = strconv.Atoi(string(restResult))
	/*jsonParser := json.NewDecoder(bufio.NewReader(res.Body))
	err = jsonParser.Decode(&jsonRestRes)*/
	if err != nil {
		http.Error(w, fmt.Sprintf("EnrollUser: error in parsing REST's result JSON: %v, err=%v", string(restResult), err),
			http.StatusInternalServerError)
		return
	}
	//fmt.Printf("DEBUG Rest result ArrayInd = %v, restResult = %v, restResult Length = %v\n", jsonRestRes.ArrayInd, restResult, len(restResult))

	ethClient, err := ethclient.Dial(gConfig.IPCpath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to the Ethereum client: %v", err), http.StatusInternalServerError)
		return
	}

	pkiWebContract, err := NewLuxUni_PKI_web(common.HexToAddress(gConfig.ContractWebHash), ethClient)
	if err != nil {
		log.Fatalf("Failed to instantiate a smart contract: %v", err)
	}

	keyFile, err := FindKeyFile(common.HexToAddress(gConfig.AccountAddr))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find key file for account %v. %v ",
			curUserAddr.String(), err), http.StatusInternalServerError)
		return
	}
	key, err := ioutil.ReadFile(gConfig.KeyDir + keyFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read key file for the Ethereum client: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Found Ethereum Key File \n")

	auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	sess := &LuxUni_PKI_webSession{
		Contract: pkiWebContract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: big.NewInt(2000000),
		},
	}
	/*sess.TransactOpts = *auth
	sess.TransactOpts.GasLimit = big.NewInt(2000000) // Rinkeby block gas limit 6124970 */

	_, err = sess.NewRegDatum(parentAddr, big.NewInt(int64(jsonRestRes.ArrayInd)),
		newUserAddr, contrAddr, fileName, "") /* contrAddr, fileName, desc, "", newUserAddr */
	if err != nil {
		http.Error(w, fmt.Sprintf("EnrollUser: Failed to add a record to blockchain: %v. ", err),
			http.StatusInternalServerError)
		return
	}

	hashResult := "Hash info is successfully loaded to blockchain."
	if gConfig.WebMode != 0 {
		var formParam CSimpleFormParam
		//formParam.Param = ""
		formParam.Params = append(formParam.Params,
			CParam{"ParentAddr", parentAddr.String()})
		formParam.Result = hashResult
		formParam.CallForm = "pki-test"

		tmpl := mTempl["SimpleForm"]
		terr := tmpl.Execute(w, formParam)
		if terr != nil {
			http.Error(w, terr.Error(), http.StatusInternalServerError)
		}
	}

}

/*#
  # https://blog.saush.com/2015/03/18/html-forms-and-go/
*/
func PkiForm(w http.ResponseWriter, r *http.Request) {

	var formParam CPKIFormParam
	var revokedParam CRevokedFormParam
	var isRevokeListRequest bool
	var isCurl []string
	//var addrCAs []common.Address
	var parentAddr common.Address
	var superParentAddr common.Address
	var superUserAddr common.Address
	var warningMessage string
	//var isChangeData bool = false;

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("PKI Main form: error in parsing multipart form: %v\n", err.Error())
		parentAddr = common.HexToAddress(gConfig.ContractHash)
	} else {
		strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
		if len(strParentAddrArr) > 0 {
			if common.IsHexAddress(strParentAddrArr[0]) == false {
				http.Error(w, GeneralError{fmt.Sprintf(
					"Parent address is incorrect: %s", strParentAddrArr[0])}.Error(),
					http.StatusInternalServerError)
				return
			}
			parentAddr = common.HexToAddress(strParentAddrArr[0])
		}

		strSuperParentAddrArr := r.MultipartForm.Value["SuperParentAddr"]
		if len(strSuperParentAddrArr) > 0 {
			if common.IsHexAddress(strSuperParentAddrArr[0]) == false {
				http.Error(w, GeneralError{fmt.Sprintf(
					"SuperParent address is incorrect: %s", strSuperParentAddrArr[0])}.Error(),
					http.StatusInternalServerError)
				return
			}
			superParentAddr = common.HexToAddress(strSuperParentAddrArr[0])
		}

		strSuperUserAddrArr := r.MultipartForm.Value["UserAddr"]
		if len(strSuperUserAddrArr) > 0 {
			if common.IsHexAddress(strSuperUserAddrArr[0]) == false {
				http.Error(w, GeneralError{fmt.Sprintf(
					"User address is incorrect: %v", strSuperUserAddrArr[0])}.Error(),
					http.StatusInternalServerError)
				return
			}
			superUserAddr = common.HexToAddress(strSuperUserAddrArr[0])
		}

		strRevokeListRequest := r.MultipartForm.Value["RevokeListButton"]
		if len(strRevokeListRequest) > 0 {
			isRevokeListRequest = true
		}
	}

	/*
	   randomBase := 1000000
	   random := rand.New(rand.NewSource(time.Now().UnixNano()));
	   formParam.RandomNum = random.Intn(9*randomBase) + randomBase
	   formParam.RandomStr = strconv.Itoa(formParam.RandomNum)
	*/

	fmt.Println("Debug: passing in GET mode")

	err = r.ParseForm()
	if err != nil {
		fmt.Printf("Error in parsing get form: %v\n", err.Error())
	}
	strParentAddr := r.FormValue("ParentAddr")
	if len(strParentAddr) > 0 {
		if common.IsHexAddress(strParentAddr) == true {
			parentAddr = common.HexToAddress(strParentAddr)
		} else {
			http.Error(w, GeneralError{fmt.Sprintf(
				"Parent address is incorrect: %v", strParentAddr)}.Error(),
				http.StatusInternalServerError)
			return
		}
	}

	if (parentAddr == common.Address{}) {
		//http.Error(w, GeneralError{"Parent address is nil: "}.Error(),
		//    http.StatusInternalServerError)
		//return
		parentAddr = common.HexToAddress(gConfig.ContractHash)
		fmt.Printf("Contract address set to default: %v\n", common.ToHex(parentAddr.Bytes()))
	}

	if (superUserAddr == common.Address{}) {
		//fmt.Printf("Attention! pki-form: user address is zero, default config account is used\n")
		superUserAddr = common.HexToAddress(gConfig.AccountAddr)
		fmt.Printf("User address in the Main Form: %v\n", superUserAddr.String())
	}

	client, err := ethclient.Dial(gConfig.IPCpath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to the Ethereum client: %v", err),
			http.StatusInternalServerError)
		return
	}

	pkiContr, err := NewLuxUni_PKI(parentAddr, client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to instantiate a smart contract: ", err),
			http.StatusInternalServerError)
		return
	}
	pkiWebContr, err := NewLuxUni_PKI_web(common.HexToAddress(gConfig.ContractWebHash), client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to instantiate a smart contract: ", err),
			http.StatusInternalServerError)
		return
	}
	callOpts := &bind.CallOpts{
		Pending: true,
	}

	ownerAddr, err := pkiContr.GetOwner(callOpts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to GetOwner: %v", err),
			http.StatusInternalServerError)
		return
	}
	if ownerAddr != superUserAddr {
		superUserAddr = ownerAddr
		fmt.Printf("User address in the Main Form changed from default to %v\n", superUserAddr.String())
	}

	if (superParentAddr == common.Address{}) {
		certData, err := pkiContr.GetCaCertificate(callOpts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to Get Certificate: %v", err),
				http.StatusInternalServerError)
			return
		}
		if len(certData) > 100 {
			certContrAddr, certParentAddr /*desc*/, _, err := ParseCert(certData)
			if err != nil {
				warningMessage = warningMessage + fmt.Sprintf(
					" WARNING: The certificate of this CA smart contract is incorrect! Parse error: %v",
					err)
			}
			if certContrAddr != parentAddr {
				warningMessage = warningMessage +
					" WARNING: The address in the request does not correspond to the address in the certificate!"
			} else {
				superParentAddr = certParentAddr
			}
		} else {
			warningMessage = warningMessage +
				"WARNING: The certificate is empty! Please populate the certificate for this contract."
		}
	}

	numRD, err := pkiContr.GetNumRegData(callOpts)
	if err != nil {
		log.Fatalf("Failed to retrieve a total number of data records: %v", err)
	}
	fmt.Println("Number of data records (including those deleted)", numRD)

	for i := int64(0); i < numRD.Int64(); i++ {
		bi := big.NewInt(i)
		var regDatum CRegData
		regDatum.ContrAddr, err = pkiWebContr.GetRegContrAddr(callOpts, parentAddr, bi)
		if err != nil {
			http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
			return
		}
		regDatum.EthAccCA, err = pkiWebContr.GetRegEthAccCA(callOpts, parentAddr, bi)
		if err != nil {
			http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
			return
		}
		tmpHash, err := pkiContr.GetRegDataHash(callOpts, bi)
		if err != nil {
			http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
			return
		}
		regDatum.DataHash = tmpHash[:]

		regDatum.CreationDate, err = pkiContr.GetRegCreationDate(callOpts, bi)
		if err != nil {
			http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
			return
		}
		regDatum.FileName, err = pkiWebContr.GetRegFileName(callOpts, parentAddr, bi)
		if err != nil {
			http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
			return
		}
		regDatum.Description, err = pkiWebContr.GetRegDescription(callOpts, parentAddr, bi)
		if err != nil {
			http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
			return
		}

		delRegDate, err := pkiContr.GetDeletedRegDate(callOpts, bi)
		if err != nil {
			http.Error(w, fmt.Sprintf("Deleted data is not retrieved: %v", err),
				http.StatusInternalServerError)
			return
		}

		var strContrAddr string
		var isCA bool
		if (regDatum.ContrAddr != common.Address{}) {
			//addrCAs = append(addrCAs, regDatum.ContrAddr)
			//formParam.AddrCAs = append(formParam.AddrCAs, common.ToHex(regDatum.ContrAddr.Bytes()))
			strContrAddr = common.ToHex(regDatum.ContrAddr.Bytes())
			isCA = true
		}
		crDate := time.Unix(regDatum.CreationDate.Int64(), 0)

		hashStr := hex.EncodeToString(regDatum.DataHash)
		docPrez := CDocPrez{int(i), regDatum.FileName, parentAddr.String(),
			strContrAddr, regDatum.EthAccCA.String(),
			regDatum.Description, "", /* link */
			"" /*decrypt*/, "0x" + hashStr /*[:10] + "..."*/, crDate,
			crDate.String(), isCA}

		// formation of data for (presentation of) HTML forms
		if delRegDate.Int64() == 0 {
			formParam.Docs = append(formParam.Docs, docPrez)
			//fmt.Printf("Data received: %v, %v \n", regDatum.Description, crDate.String())
		} else {
			revokedParam.Docs = append(revokedParam.Docs, docPrez)
		}
	}

	// CURL does neet any web return - just "OK" of no errors
	if isCurl != nil {
		jsonResp := `{"verbose":"OK","result":0}`
		fmt.Fprintln(w, jsonResp)
		return
	}

	if (len(revokedParam.RevokedIds) > 0) || (isRevokeListRequest != false) {
		tmpl := mTempl["HashResult"]
		terr := tmpl.Execute(w, revokedParam)
		if terr != nil {
			http.Error(w, terr.Error(), http.StatusInternalServerError)
		}
	}

	tmpl := mTempl["MainForm"]
	formParam.ParentAddr = common.ToHex(parentAddr.Bytes())
	formParam.SuperParentAddr = common.ToHex(superParentAddr.Bytes())
	if (superParentAddr == common.Address{}) {
		formParam.IsRoot = true
	}
	formParam.SuperUserAddr = common.ToHex(superUserAddr.Bytes())
	terr := tmpl.Execute(w, formParam)
	if terr != nil {
		http.Error(w, terr.Error(), http.StatusInternalServerError)
	}
}

func PopulateContract(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
		return
	}

	var b bytes.Buffer
	newWri := multipart.NewWriter(&b)

	strContrAddr, err := CopyRequestField(r, newWri,
		"ContrAddr", "ContrAddr")
	if err != nil {
		http.Error(w, fmt.Sprintf("Populate contract: %v", err),
			http.StatusInternalServerError)
		return
	}
	if common.IsHexAddress(strContrAddr) == false {
		http.Error(w, fmt.Sprintf("Contract address is incorrect: %v", strContrAddr),
			http.StatusInternalServerError)
		return
	}
	contrAddr := common.HexToAddress(strContrAddr)

	parentAddr := common.Address{}
	strParentAddr, err := CopyRequestField(r, newWri,
		"ParentAddr", "ParentAddr")
	if err == nil {
		if common.IsHexAddress(strContrAddr) == false {
			http.Error(w, fmt.Sprintf("Contract address is incorrect: %v", strContrAddr),
				http.StatusInternalServerError)
			return
		}
		parentAddr = common.HexToAddress(strParentAddr)
	}

	newUserAddr := common.Address{}
	strNewUserAddr, err := CopyRequestField(r, newWri,
		"NewUserAddr", "NewUserAddr")
	if err == nil {
		if common.IsHexAddress(strContrAddr) == false {
			http.Error(w, fmt.Sprintf("Contract address is incorrect: %v", strContrAddr),
				http.StatusInternalServerError)
			return
		}
		newUserAddr = common.HexToAddress(strNewUserAddr)
	}

	curUserAddr := common.Address{}
	strCurUserAddr, err := CopyRequestField(r, newWri,
		"CurrentUserAddr", "CurrentUserAddr")
	if err == nil {
		if common.IsHexAddress(strContrAddr) == false {
			http.Error(w, fmt.Sprintf("Populate: Current user address is incorrect: %v", strContrAddr),
				http.StatusInternalServerError)
			return
		}
		curUserAddr = common.HexToAddress(strCurUserAddr)
	}

	if (curUserAddr == common.Address{}) && (newUserAddr == common.Address{}) {
		if (contrAddr != common.HexToAddress(gConfig.ContractHash)) || (parentAddr != common.Address{}) {
			http.Error(w, fmt.Sprintf(
				"Root Populating: the Curr Contr should be default in config (now is %v) and Parrent Contr should be null (now is %v)",
				contrAddr.String(), parentAddr.String()), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Populating root certificate\n")
	} else {
		if (curUserAddr == common.Address{}) || (newUserAddr == common.Address{}) {
			http.Error(w, fmt.Sprintf(
				"Populate error: curUserAddr (%v) or newUserAddr (%v) is null",
				curUserAddr.String(), newUserAddr.String()), http.StatusInternalServerError)
			return
		}
	}

	dataCert, err := GenerateCert(contrAddr, parentAddr, true, "Mother Nature CA")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashCert, err := CalcHash(dataCert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fw, err := newWri.CreateFormFile("UplFiles", "CA_cert.crt")
	if err != nil {
		http.Error(w, fmt.Sprintf("EnrollUser: cannot create form file: %v", err),
			http.StatusInternalServerError)
		return
	}
	if _, err = io.Copy(fw, bytes.NewReader(dataCert)); err != nil {
		http.Error(w, fmt.Sprintf("EnrollUser: cannot read/copy cert data: %v", err),
			http.StatusInternalServerError)
		return
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	newWri.Close()
	// Now that you have a form, you can submit it to your handler.
	newReq, err := http.NewRequest("POST",
		gConfig.RestUrlServer+":"+strconv.Itoa(gConfig.RestHttpPort)+"/populate_contract", &b)
	if err != nil {
		http.Error(w, fmt.Sprintf("Populate: error in creation of REST request - %v", err),
			http.StatusInternalServerError)
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	newReq.Header.Set("Content-Type", newWri.FormDataContentType())

	// Submit the request
	// https://stackoverflow.com/questions/20205796/golang-post-data-using-the-content-type-multipart-form-data
	client := &http.Client{}
	res, err := client.Do(newReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Populate: error in REST request - %v", err),
			http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		restResult, _ := ioutil.ReadAll(res.Body)
		http.Error(w, fmt.Sprintf("Populate: REST's bad status: %v, %v", res.Status, string(restResult)),
			http.StatusInternalServerError)
		return
	}

	resStr := "Contract was successfully populated\n"
	detailStr := "Contract Addr = " + contrAddr.String() + "\n"
	detailStr += "Parent Addr = " + common.ToHex(parentAddr.Bytes()) + "\n"
	detailStr += "New User = " + common.ToHex(newUserAddr.Bytes()) + "\n"
	detailStr += "Current User = " + common.ToHex(curUserAddr.Bytes()) + "\n"
	if gConfig.WebMode != 0 {

		var formParam CSimpleFormParam
		formParam.Result = resStr
		formParam.Data = detailStr

		if (curUserAddr == common.Address{}) && (newUserAddr == common.Address{}) { /* Root populating */
			formParam.CallForm = "pki-test"
		} else { /* if it is not a root populating, we should add a record with the CA certificate to the blockchain */
			formParam.Params = append(formParam.Params,
				CParam{"Hash", hex.EncodeToString(hashCert)})
			formParam.Params = append(formParam.Params,
				CParam{"ContrAddr", contrAddr.String()})
			formParam.Params = append(formParam.Params,
				CParam{"ParentAddr", common.ToHex(parentAddr.Bytes())})
			formParam.Params = append(formParam.Params,
				CParam{"CurrentUserAddr", common.ToHex(curUserAddr.Bytes())})
			formParam.Params = append(formParam.Params,
				CParam{"NewUserAddr", common.ToHex(newUserAddr.Bytes())})
			formParam.CallForm = "enroll_user"
		}

		tmpl := mTempl["SimpleForm"]
		terr := tmpl.Execute(w, formParam)
		if terr != nil {
			http.Error(w, terr.Error(), http.StatusInternalServerError)
		}
	}
}

func CreateContract(w http.ResponseWriter, r *http.Request) {
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
				http.Error(w, GeneralError{"Parent address is incorrect"}.Error(),
					http.StatusInternalServerError)
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
			http.Error(w, GeneralError{"Current user address is incorrect"}.Error(),
				http.StatusInternalServerError)
			return
		}
	}

	strUserAddrArr = r.MultipartForm.Value["NewUserAddr"]
	if len(strUserAddrArr) > 0 {
		newUserAddrStr = strUserAddrArr[0]
		if common.IsHexAddress(newUserAddrStr) == false {
			fmt.Println("Create Contract: New user address is incorrect")
			http.Error(w, "New user address is incorrect", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "New user address is not available", http.StatusInternalServerError)
		return
	}

	client, err := ethclient.Dial(gConfig.IPCpath)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	pkiContr, err := NewLuxUni_PKI(common.HexToAddress(parentAddrStr), client)
	if err != nil {
		http.Error(w, GeneralError{fmt.Sprintf(
			"Failed to instantiate a smart contract: ", err)}.Error(),
			http.StatusInternalServerError)
		return
	}
	callOpts := &bind.CallOpts{
		Pending: true,
	}
	execUserAddr, err := pkiContr.GetOwner(callOpts)
	if err != nil {
		http.Error(w, GeneralError{fmt.Sprintf(
			"CreateCont - failed to get owner addr: ", err)}.Error(),
			http.StatusInternalServerError)
		return
	}
	if execUserAddr != common.HexToAddress(curUserAddrStr) {
		http.Error(w, "contract.GetOwner does not correspond to the Current User param",
			http.StatusInternalServerError)
		return
	}

	//keyFile := gConfig.KeyFile
	keyFile, err := FindKeyFile(execUserAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("CreateContract -- FindKeyFile: %v. ", err),
			http.StatusInternalServerError)
		return
	}
	key, e := ioutil.ReadFile(gConfig.KeyDir + keyFile)
	if e != nil {
		fmt.Printf("Key File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("Found Ethereum Key File \n")

	auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	var trOpts bind.TransactOpts = *auth
	trOpts.GasLimit = big.NewInt(4000000) // 6124970 - block gas limit in Rinkeby
	contrAddr, _ /*contr*/, _, err := DeployLuxUni_PKI(&trOpts, client)
	/*
	   https://stackoverflow.com/questions/40096750/set-status-code-on-http-responsewriter
	*/
	if err != nil {
		http.Error(w, fmt.Sprintf("CreateContract -- Etherreum error in contract creation: %v", err),
			http.StatusInternalServerError)
		return
	}

	if gConfig.WebMode != 0 {
		var formParam CSimpleFormParam
		formParam.IsUploadFileForm = false
		formParam.Params = append(formParam.Params,
			CParam{"ContrAddr", common.ToHex(contrAddr.Bytes())})
		formParam.Params = append(formParam.Params,
			CParam{"ParentAddr", parentAddrStr})
		formParam.Params = append(formParam.Params,
			CParam{"CurrentUserAddr", curUserAddrStr})
		formParam.Params = append(formParam.Params,
			CParam{"NewUserAddr", newUserAddrStr})
		formParam.Result = "Contract was successfully created.\n Populate the contract"
		formParam.CallForm = "populate_contract"

		tmpl := mTempl["SimpleForm"]
		terr := tmpl.Execute(w, formParam)
		if terr != nil {
			http.Error(w, terr.Error(), http.StatusInternalServerError)
		}
	}
}

func ValidateForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
		return
	}

	arrStrParentAddr := r.MultipartForm.Value["ParentAddr"]
	if len(arrStrParentAddr) == 0 {
		http.Error(w, GeneralError{"No parentAddr is provided"}.Error(),
			http.StatusInternalServerError)
		return
	}
	strParentAddr := arrStrParentAddr[0]
	if common.IsHexAddress(strParentAddr) == false {
		http.Error(w, GeneralError{"Parent address is incorrect"}.Error(),
			http.StatusInternalServerError)
		return
	}
	parentAddr := common.HexToAddress(strParentAddr)

	if gConfig.WebMode != 0 {
		var formParam CSimpleFormParam
		formParam.IsUploadFileForm = true
		formParam.Params = append(formParam.Params,
			CParam{"ParentAddr", common.ToHex(parentAddr.Bytes())})
		formParam.Result = "Please upload a certificate to check for CA " +
			common.ToHex(parentAddr.Bytes())
		formParam.CallForm = "validate_cert"

		tmpl := mTempl["SimpleForm"]
		terr := tmpl.Execute(w, formParam)
		if terr != nil {
			http.Error(w, terr.Error(), http.StatusInternalServerError)
		}
	}
}

func ValidateCert(w http.ResponseWriter, r *http.Request) {

	var parentAddr common.Address = common.Address{}
	var isNoUpload bool = false

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
		return
	}

	var b bytes.Buffer
	newWri := multipart.NewWriter(&b)

	certHashStr := ""
	certHash, uplFileName, dataCert, cerr := UploadFile(w, r, "UplFiles", true)
	if cerr.errCode != 0 {
		if cerr.errCode == 1 {
			isNoUpload = true
			hashArr := r.MultipartForm.Value["Hash"]
			if len(hashArr) == 0 {
				http.Error(w, fmt.Sprintf("Rest Validate: No hashes in request"),
					480 /* http.StatusInternalServerError */)
				return
			}

			fw, err := newWri.CreateFormField("Hash")
			if err != nil {
				http.Error(w, "VaidateCert: Hash is not created in request to REST",
					http.StatusInternalServerError)
				return
			}
			if _, err = fw.Write([]byte(hashArr[0])); err != nil {
				http.Error(w, "ValidateCert: Hash is not written in request to REST",
					http.StatusInternalServerError)
				return
			}
			if len(certHash) == 0 {
				certHashStr = hashArr[0]
			}
		} else {
			http.Error(w, cerr.Error(), http.StatusInternalServerError)
			return
		}
	}
	if len(certHash) > 0 {
		certHashStr = fmt.Sprintf("%x", certHash)
	}

	if isNoUpload == true {
		arrStrParentAddr := r.MultipartForm.Value["ParentAddr"]
		if len(arrStrParentAddr) != 0 {
			strParentAddr := arrStrParentAddr[0]
			if common.IsHexAddress(strParentAddr) == false {
				http.Error(w, GeneralError{"Parent address is incorrect"}.Error(),
					http.StatusInternalServerError)
				return
			}
			parentAddr = common.HexToAddress(strParentAddr)

			fw, err := newWri.CreateFormField("ParentAddr")
			if err != nil {
				http.Error(w, "VaidateCert: ParentAddr is not created in request to REST",
					http.StatusInternalServerError)
				return
			}
			if _, err = fw.Write([]byte(parentAddr.String())); err != nil {
				http.Error(w, "ValidateCert: ParentAddr is not written in request to REST",
					http.StatusInternalServerError)
				return
			}
		}
	}

	if isNoUpload == false {
		fw, err := newWri.CreateFormFile("UplFiles", uplFileName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Validate: cannot create form file: %v", err),
				http.StatusInternalServerError)
			return
		}
		if _, err = io.Copy(fw, bytes.NewReader(dataCert)); err != nil {
			http.Error(w, fmt.Sprintf("validate: cannot read/copy cert data: %v", err),
				http.StatusInternalServerError)
			return
		}
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	newWri.Close()
	// Now that you have a form, you can submit it to your handler.
	newReq, err := http.NewRequest("POST",
		gConfig.RestUrlServer+":"+strconv.Itoa(gConfig.RestHttpPort)+"/validate_cert", &b)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validate: error in creation of REST request - %v", err),
			http.StatusInternalServerError)
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	newReq.Header.Set("Content-Type", newWri.FormDataContentType())

	// Submit the request
	// https://stackoverflow.com/questions/20205796/golang-post-data-using-the-content-type-multipart-form-data
	client := &http.Client{}
	res, err := client.Do(newReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validate: error in REST request - %v", err),
			http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		restResult, _ := ioutil.ReadAll(res.Body)
		http.Error(w, fmt.Sprintf("Validate: REST's bad status: %v, %v", res.Status, string(restResult)),
			http.StatusInternalServerError)
		return
	}

	jsonResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validate: cannot read answer the body of REST answer: %v", err),
			http.StatusInternalServerError)
		return
	}

	var response JsonValidateResponse
	err = json.Unmarshal(jsonResponse, &response)
	if err != nil {
		http.Error(w, fmt.Sprintf("ValidateCert: Json parse: %v", err.Error()),
			http.StatusInternalServerError)
		return
	}
	//https://stackoverflow.com/questions/19038598/how-can-i-pretty-print-json-using-go
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(jsonResponse), "", "\t")
	if err != nil {
		http.Error(w, fmt.Sprintf("ValidateCert: Json indent: %v", err.Error()),
			http.StatusInternalServerError)
		return
	}

	revokeDateStr := "No"
	if (response.RevokeDate.String() != time.Time{}.String()) {
		revokeDateStr = "Yes, a certificate was revoked at " + strconv.Itoa(response.Iter) +
			" level above on " + response.RevokeDate.String() //+ " zero: " + zeroTime.String()
	}
	retDataStr := "Is the certificate valid? " + strconv.FormatBool(response.IsCertOK) + "\n"
	retDataStr += "Was there revocation (date)? " + revokeDateStr + "\n"
	retDataStr += "Level of the initial certificate (check level iterations)? " +
		strconv.Itoa(response.Iter) + "\n"
	retDataStr += "\nJson path of the certificate check: \n" + string(prettyJSON.Bytes())

	if gConfig.WebMode != 0 {
		var formParam CSimpleFormParam
		formParam.IsUploadFileForm = false
		formParam.Params = append(formParam.Params,
			CParam{"ParentAddr", parentAddr.String()})
		formParam.Result = "Results of certificate check (certHash=" +
			certHashStr + ") for CA " // + contrAddr.String()
		formParam.Data = retDataStr
		formParam.CallForm = "pki-test"

		tmpl := mTempl["SimpleForm"]
		terr := tmpl.Execute(w, formParam)
		if terr != nil {
			http.Error(w, terr.Error(), http.StatusInternalServerError)
		}
	}
}

func GenerateUserCert(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
		return
	}

	//isCurl = r.MultipartForm.Value["Curl"]
	//certHash := r.MultipartForm.Value["CertHash"][0]
	if len(r.MultipartForm.Value["InsertAddr"]) == 0 {
		http.Error(w, GeneralError{"No insertAddr is provided"}.Error(),
			http.StatusInternalServerError)
		return
	}
	strInsertAddr := r.MultipartForm.Value["InsertAddr"][0]
	if common.IsHexAddress(strInsertAddr) == false {
		http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
			http.StatusInternalServerError)
		return
	}
	insertAddr := common.HexToAddress(strInsertAddr)

	if len(r.MultipartForm.Value["Name"]) == 0 {
		http.Error(w, GeneralError{"No user name is provided"}.Error(),
			http.StatusInternalServerError)
		return
	}
	strName := GetCorrectFileName(r.MultipartForm.Value["Name"][0])

	dataCert, err := GenerateCert(common.Address{}, insertAddr, false, strName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+strName+".crt.out")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	http.ServeContent(w, r, strName+".crt.out", time.Now(), bytes.NewReader(dataCert))
}

/*
  https://stackoverflow.com/questions/35496233/go-how-to-i-make-download-service
  https://play.golang.org/p/UMKgI_NLwO
*/
func DownloadCaCert(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
		return
	}

	if len(r.MultipartForm.Value["ContrAddr"]) == 0 {
		http.Error(w, GeneralError{"No contrAddr is provided"}.Error(),
			http.StatusInternalServerError)
		return
	}
	strContrAddr := r.MultipartForm.Value["ContrAddr"][0]
	if common.IsHexAddress(strContrAddr) == false {
		http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
			http.StatusInternalServerError)
		return
	}
	contrAddr := common.HexToAddress(strContrAddr)

	_ /*isCertOK*/, _ /*revokDate*/, _ /*parentAddr*/, _ /*retCaHash*/, certData, err :=
		ConfirmHashCAData(contrAddr, nil, true)

	w.Header().Set("Content-Disposition", "attachment; filename=ca.crt.out")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	http.ServeContent(w, r, "ca.crt.out", time.Now(), bytes.NewReader(certData))
}

func GenerateCert(contrAddr common.Address, parentAddr common.Address,
	isCA bool, strName string) ([]byte, error) {

	// see Certificate structure at
	// http://golang.org/pkg/crypto/x509/#Certificate
	template := &x509.Certificate{
		IsCA: isCA,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1, 2, 3},
		SerialNumber:          big.NewInt(1234),
		Subject: pkix.Name{
			Country:      []string{"Earth"},
			Organization: []string{strName},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(5, 5, 5),
		// see http://golang.org/pkg/crypto/x509/#KeyUsage
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	// https://stackoverflow.com/questions/26441547/go-how-do-i-add-an-extension-subjectaltname-to-a-x509-certificate
	template.ExtraExtensions = append(template.ExtraExtensions, // blockchain name
		GenerateCertExtension(asn1.ObjectIdentifier{1, 2, 752, 115, 33, 3},
			common.Hex2Bytes("0x130d457468657265756d2d54657374"), false))
	if isCA == true {
		template.ExtraExtensions = append(template.ExtraExtensions,
			GenerateCertExtension(asn1.ObjectIdentifier{1, 2, 752, 115, 33, 2},
				contrAddr.Bytes(), true))
	}
	template.ExtraExtensions = append(template.ExtraExtensions,
		GenerateCertExtension(asn1.ObjectIdentifier{1, 2, 752, 115, 33, 1},
			parentAddr.Bytes(), true))
	template.ExtraExtensions = append(template.ExtraExtensions, // hash algo
		GenerateCertExtension(asn1.ObjectIdentifier{1, 2, 752, 115, 33, 0},
			common.Hex2Bytes("0x0609608648016503040201"), false))

	// generate private key
	privatekey, err := rsa.GenerateKey(cryptoRand.Reader, 2048)

	if err != nil {
		return nil, err
	}

	publickey := &privatekey.PublicKey

	// create a self-signed certificate. template = parent
	var parent = template
	cert, err := x509.CreateCertificate(cryptoRand.Reader, template, parent, publickey, privatekey)

	if err != nil {
		return nil, err
	}

	return cert, nil
}

func GenerateCertExtension(ID asn1.ObjectIdentifier, val []byte,
	isAddPrefix bool) pkix.Extension {
	var extVal []byte
	if isAddPrefix == true {
		extVal = []byte{4, 20}
	}
	extVal = append(extVal, val...)
	ext := pkix.Extension{}
	ext.Critical = false
	ext.Id = ID
	ext.Value = extVal
	return ext
}

/*
   processes and copies field from incoming request to outgoing request to REST service
*/
func CopyRequestField(r *http.Request, newWri *multipart.Writer,
	fieldNameIn string, fieldNameOut string) (fieldVal string, err error) {
	strFieldArr := r.MultipartForm.Value[fieldNameIn]
	if len(strFieldArr) == 0 {
		return "", GeneralError{fmt.Sprintf("Request fields are incomplete: The number of items for %v is %v (should be = 1)",
			fieldNameIn, len(strFieldArr))}
	}
	if len(strFieldArr) > 1 {
		return "", GeneralError{fmt.Sprintf("Request fields are excessive: The number of items for %v is %v (>1)",
			fieldNameIn, len(strFieldArr))}
	}

	fieldVal = strFieldArr[0]

	fw, err := newWri.CreateFormField(fieldNameOut)
	if err != nil {
		return "", GeneralError{fmt.Sprintf("%v is not created in request to REST", fieldNameOut)}
	}
	if _, err = fw.Write([]byte(fieldVal)); err != nil {
		return "", GeneralError{fmt.Sprintf("%v is not written in request to REST", fieldNameOut)}
	}
	return fieldVal, nil
}

/*
   processes and copies field from incoming request to outgoing request to REST service
*/
func GetCorrectFileName(originalName string) string {
	return strings.Replace(originalName, " ", "_", -1)
}
