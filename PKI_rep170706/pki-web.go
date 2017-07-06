
package main

import (
    "html/template"
    "time"
    "fmt"
    "io/ioutil"
    "os"
    "log"
    "net/http"
    "strings"
    "strconv"
    "math/big"
    "hash/crc32"
    cryptoRand "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "encoding/gob"
    "github.com/gorilla/mux"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/ethclient"
    "encoding/json"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/asn1"
    //"math/rand"
    "bytes"
)

var mTempl map[string]*template.Template
//var session *LuxUni_PKISession
//var gClient *ethclient.Client
var gCaCertOffset int = 64;
var gCaCertName = "_ca_cert.crt"

const gConfigFile = "./pki-conf.json";

/*
CORE PARAMETERS ARE STORED IN gConfig 

const gCryptoModulHash = "0x3f2ed40488d0a9586013faa415718f3f64644fa1";
const gContractHash = "0xf1918d06d7e66e60153d7109ff380b41866ba2e0";
const gIPCpath = "/home/alex/_Work/Eth_AeroNet_t/geth.ipc";
const gPswd = "ira";
const key = `{"address":"a6f23407d139508fa38706140c56bf6487f87395","crypto":{"cipher":"aes-128-ctr","ciphertext":"5468728413080efcf6191dc3b5eaaf6ff34fd401277b9c08ac5aa93b6c3b3e44","cipherparams":{"iv":"efae50c39291e191193635ec715e96f4"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"a1e83898937dcd79e480d642aab80032a983a86e55c1ff8c64e91a1fc4fd9bdd"},"mac":"35456f2b746b847b4f1f5af2291f8f1518d6eed4e29ce950c1123da7ac3ee01d"},"id":"71412808-a83f-4677-ac6d-3d0f744d3ef0","version":3}`
const gPrivateKeyPath = "a_pr.key";
*/

var gPrivateKey *rsa.PrivateKey;

var gConfig struct {
    //CryptoModulHash    string `json:"cryptoModulHash"`
    ContractHash       string `json:"contractHash"`
    IPCpath            string `json:"IPCpath"`
    Pswd               string `json:"pswd"`
    KeyDir             string `json:"keyDir"`
    AccountAddr        string `json:"accountAddr"`  // !! should start with 0x
    PrivateKeyPath     string `json:"privateKeyPath"`
    HttpPort           int `json:"httpPort"`
    WebMode            int `json:"webMode"`
    JsonMode           int `json:"jsonMode"`
}

type CRegData struct {
	NodeSender   common.Address
    ContrAddr    common.Address
	DataHash     []byte
	FileName     string
	Description  string
	Encrypted    *big.Int
	CryptoModule common.Address
	LinkFile     string
	CreationDate *big.Int
	Active       bool
}

// alternative aproach to passing the params to template - http://stackoverflow.com/questions/23802008/how-do-you-pass-multiple-objects-to-go-template
// general Go template description https://gohugo.io/templates/go-templates/
type CPKIFormParam struct {
    RandomNum    int
    RandomStr    string
    ParentAddr   string
    SuperParentAddr   string
    SuperUserAddr   string
    Docs         []CDocPrez
    //AddrCAs      []string
}

type CDocPrez struct {
    Id           int
    Name         string
    ParentAddr    string
    ContrAddr    string
    UserAddr     string
    Desc         string
    Link         string
    Decryption   string
    Hash         string
    CreationDate  time.Time
    CreationStr   string
    IsCA        bool
}

type CRevokedFormParam struct {
    RevokedIds    []int
    RevokedStr   string
    Docs         []CDocPrez
}

type CSimpleFormParam struct {
    Result       string
    Data         string
    Params   []CParam
    CallForm     string
    IsUploadFileForm bool
}

type CParam struct {
    Name string
    Value string
}

func init() {
    err := LoadConfig();
    if err != nil {
        fmt.Printf("CONFIG ERROR: %v\n", err)
        os.Exit(1)
    }
    gConfig.WebMode = 1

    if mTempl == nil {
        mTempl = make( map[string]*template.Template )
    }

    ///home/alex/DocsFS/Dropbox/WORK/RD/LuxBCh/PKI
    t, e := template.ParseFiles("./form_main.html")
    if e != nil {
        fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "main", e.Error())
        mTempl = nil
    } else {
        mTempl["MainForm"] = template.Must(t, e)
    }

    t, e = template.ParseFiles("./form_hashres.html")
    if e != nil {
        fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "hashes", e.Error())
        mTempl = nil
    } else {
        mTempl["HashResult"] = template.Must(t, e)
    }

    t, e = template.ParseFiles("./form_simple.html")
    if e != nil {
        fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "simple", e.Error())
        mTempl = nil
    } else {
        mTempl["SimpleForm"] = template.Must(t, e)
    }

}


/*
    idFormStr -- id of the loading files in the html form, for example "UplFiles"
    return params
        uint32 -- hash
        string -- fileName
        //string -- path to the temporary file on the disk
        []byte -- pointer to the file content
        *GeneralCodeError
    return Err code:
        1 -- no file is found
        2 -- too many files are uploaded
        3 -- others
*/
func UploadFile(w http.ResponseWriter, r *http.Request, idFormStr string, contrAddr common.Address,
            isCopyContent bool) (uint32, string, []byte, GeneralCodeError){
    m := r.MultipartForm
    files := m.File[idFormStr];

    if len(files)==0 {
        return 0, "", nil, GeneralCodeError{ "Upload File: no files are found)", 1}
    }

    if len(files)>1 {
        return 0, "", nil, GeneralCodeError{ fmt.Sprintf(
            "Upload File: number of aploaded files greater than 1 = ", len(files)), 2}
    }

    // for each fileheader, get a handle to the actual file
    // https://gobyexample.com/errors
    file, err := files[0].Open()
    defer file.Close()
    if err != nil {
        return 0, "", nil, GeneralCodeError{
            fmt.Sprintf("Upload File: parsing form: ", err.Error()),3}
    }

    /*
    //create destination file making sure the path is writeable.
    dstFName := "/tmp/tst" + "_" + files[0].Filename
    dst, err := os.Create(dstFName)
    defer dst.Close()
    if err != nil {
        return nil, nil, nil, &GeneralCodeError{
            fmt.Sprintf(
                "Upload File: unable to open the file for writing: ", err.Error()), 0}
    } else {
        fmt.Printf("Created file: %v\n", dstFName)
    }

    //copy the uploaded file to the destination file
    _, err = io.Copy(dst, file);
    if err != nil {
        return nil, nil, nil, nil, &GeneralCodeError{
            fmt.Sprintf("Upload File: unable to copy: ", err.Error()),0}
    }*/

    /*
    https://stackoverflow.com/questions/36111777/golang-how-to-read-a-text-file
    https://stackoverflow.com/questions/30182538/why-can-not-i-copy-a-slice-with-copy-in-golang
    */
    hasher := crc32.NewIEEE()
    //dst4hash, err := ioutil.ReadFile(dstFName);
    dst4hash, err := ioutil.ReadAll(file)
    if err != nil {
        return 0, "", nil, GeneralCodeError{
            fmt.Sprintf("Error in open file for hash: ", err.Error()),3}
    }
    if (contrAddr != common.Address{}){
        if cap(dst4hash) < (gCaCertOffset+len(contrAddr.Bytes())) {
            return 0, "", nil, GeneralCodeError{ fmt.Sprintf(
                "Upload File: certOffset is too large, cap(dst4Hash)= ", cap(dst4hash)), 3}
        }
        copy(dst4hash[gCaCertOffset:],contrAddr.Bytes())
    }
    var dst4return []byte = nil;
    if (isCopyContent == true){
        //dst4return = make([]byte,len(dst4hash))
        //copy(dst4hash, dst4return)
        dst4return = dst4hash
    }
    hasher.Write(dst4hash);
    hashSum := hasher.Sum32()
    fmt.Printf("Hash: %x\n", hashSum)

    return hashSum, files[0].Filename, dst4return, GeneralCodeError{"OK",0}
}


func BlacklistUser(w http.ResponseWriter, r *http.Request){
    var revokeResult string
    var parentAddr common.Address = common.Address{};
    var userAddr common.Address = common.Address{};

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("No data: Parsing blacklist multipart form: %v\n", err.Error())
        http.Error(w, GeneralError{fmt.Sprintf(
            "BlacklistUser: error in parsing -- ",err.Error())}.Error(),
            http.StatusInternalServerError)
        return;
    }

    strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
    if (len(strParentAddrArr) > 0) {
        if(common.IsHexAddress(strParentAddrArr[0])==false){
            http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        parentAddr = common.HexToAddress(strParentAddrArr[0])
    }

    if (parentAddr == common.Address{}){
        http.Error(w, GeneralError{"Delete: Parent address is not established"}.Error(),
            http.StatusInternalServerError)
        return
    }

    strUserAddrArr := r.MultipartForm.Value["UserAddr"]
    if (len(strUserAddrArr) > 0) {
        if(common.IsHexAddress(strUserAddrArr[0])==false){
            http.Error(w, GeneralError{"User address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        userAddr = common.HexToAddress(strUserAddrArr[0])
    }

    dels := r.MultipartForm.Value["Deletion"]
    if len(dels) > 0 {
        fmt.Printf("Debug: I am in deletion block")
        //dels := r.MultipartForm.Value["Deletion"]
        //dels := r.Form["Deletion"]
        for _, del := range dels {
            fmt.Printf("del=%v\n", del);
            delid, err := strconv.Atoi(del);
            if err != nil {
                http.Error(w, fmt.Sprintf("Deletion conversion error: %v", err.Error()),
                    http.StatusInternalServerError)
                return
            }
            //revokedParam.RevokedIds = append(revokedParam.RevokedIds, delid);
            revokeResult += del+" ";

            client, err := ethclient.Dial(gConfig.IPCpath)
            if err != nil {
                http.Error(w, fmt.Sprintf("Failed to connect to the Ethereum client: %v", err),
                    http.StatusInternalServerError)
                return
            }

            // Instantiate the contract, the address is taken from eth at the moment of contract initiation
            // kyc, err := NewLuxUni_KYC(common.HexToAddress(gContractHash), backends.NewRPCBackend(conn))
            pkiContract, err := NewLuxUni_PKI(parentAddr, client)
            if err != nil {
                http.Error(w, fmt.Sprintf("Failed to instantiate a smart contract: %v", err),
                    http.StatusInternalServerError)
                return
            }

            // Logging into Ethereum as a user
            if (userAddr == common.Address{}) {
                fmt.Printf("Attention! Revoke: user address is zero, default config account is used\n")
                userAddr = common.HexToAddress(gConfig.AccountAddr)
            }
            keyFile, err := FindKeyFile(userAddr)
            if (err!=nil) {
                http.Error(w, fmt.Sprintf("Failed to find key file for account %v. %v ",
                    userAddr.String(), err), http.StatusInternalServerError)
                return
            }
            key, err := ioutil.ReadFile( gConfig.KeyDir+keyFile )
            if err != nil {
                http.Error(w, fmt.Sprintf("Key File error: %v\n", err),
                    http.StatusInternalServerError)
                return
            }

            auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
            if err != nil {
                http.Error(w, fmt.Sprintf("Failed to create authorized transactor: %v", err),
                    http.StatusInternalServerError)
                return
            }

            sess := &LuxUni_PKISession{
                Contract: pkiContract,
                CallOpts: bind.CallOpts{
                    Pending: true,
                },
                TransactOpts: bind.TransactOpts{},
            }
            sess.TransactOpts = *auth;
            sess.TransactOpts.GasLimit = big.NewInt(2000000)

            _, nerr := sess.DeleteRegDatum(big.NewInt(int64(delid)))
            if nerr != nil {
                http.Error(w, fmt.Sprintf("Deletion error: %v", nerr),
                    http.StatusInternalServerError)
                return
            }
        }
    }

    if (gConfig.WebMode!=0){
        var formParam CSimpleFormParam
        //formParam.Param = ""
        if len(revokeResult)>0 {
            formParam.Result = "CA "+parentAddr.String()+": certificate " + revokeResult +
                "was successfully revoked"
        } else {
            formParam.Result = "No revocation was conducted"
        }
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
ParentAddr
ContrAddr
Hash or Files
 */
func EnrollUser(w http.ResponseWriter, r *http.Request){
    //var isCurl []string;
    var parentAddr common.Address = common.Address{} // this is addr of the contract which is going to hold the hash
    var contrAddr common.Address = common.Address{} // this is address of the new SubCA contract or zero if end user
    var insertAddr common.Address = common.Address{} // this is address in the certificate corresponding to the parent
    var userAddr common.Address = common.Address{}
    var isNoUpload bool = false

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("EnrollUser: No change data -- ", err.Error())
        http.Error(w, GeneralError{fmt.Sprintf(
            "EnrollUser: No change data -- ",err.Error())}.Error(),
            http.StatusInternalServerError)
        return
    }

    //isCurl = r.MultipartForm.Value["Curl"]
    hashSum, fileName, dataCert, cerr :=
            UploadFile(w, r, "UplFiles", common.Address{},true)
    if (cerr.errCode!=0){
        if(cerr.errCode == 1){
            isNoUpload = true
            hashArr := r.MultipartForm.Value["Hash"]
            if (len(hashArr) == 0) {
                http.Error(w, GeneralError{fmt.Sprintf(
                    "EnrollUser: No hashes in request")}.Error(),
                    http.StatusInternalServerError)
                return
            }
            tmpInt, err := strconv.Atoi(hashArr[0]);
            if(err!=nil){
                http.Error(w, GeneralError{fmt.Sprintf(
                    "EnrollUser: Hash is incorrect",err.Error())}.Error(),
                    http.StatusInternalServerError)
                return
            }
            hashSum = uint32(tmpInt)
        }else{
            http.Error(w, GeneralError{fmt.Sprintf(
                "EnrollUser UplFiles:",cerr.Error())}.Error(),
                http.StatusInternalServerError)
            return
        }
    }

    strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
    if (len(strParentAddrArr) > 0) {
        if(common.IsHexAddress(strParentAddrArr[0])==false) {
            http.Error(w, "Parent address as a parameter is incorrect",
                http.StatusInternalServerError)
            return
        }
        parentAddr = common.HexToAddress(strParentAddrArr[0])
    }

    if isNoUpload==false {
        var caContrAddr common.Address;
        caContrAddr, insertAddr, err = ParseCert(dataCert)
        if (err != nil) {
            http.Error(w, fmt.Sprintf("CERTIFICATE: Parsing error: %v", err),
                http.StatusInternalServerError)
            return
        }
        if (insertAddr == common.Address{}) {
            http.Error(w, "CERTIFICATE: No Parent Address is provided in the Cert",
                http.StatusInternalServerError)
            return
        }
        if (caContrAddr != common.Address{}) {
            http.Error(w, "CERTIFICATE: Non-CA certificates should not include parent address data",
                http.StatusInternalServerError)
            return
        }
        if insertAddr != parentAddr {
            http.Error(w, "Address in the certificate does not correspond to the contract address of the Authority (CA)",
                http.StatusInternalServerError)
            return
        }
    }

    strUserAddrArr := r.MultipartForm.Value["UserAddr"]
    if (len(strUserAddrArr) > 0) {
        if(common.IsHexAddress(strUserAddrArr[0])==false){
            http.Error(w, GeneralError{"User address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        userAddr = common.HexToAddress(strUserAddrArr[0])
    }

    //if fileName==gCaCertName {
    strContrAddrArr := r.MultipartForm.Value["ContrAddr"]
    if( len(strContrAddrArr)>0 ){
        if(common.IsHexAddress(strContrAddrArr[0])==false){
            http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        contrAddr = common.HexToAddress(strContrAddrArr[0])
    }

    /*
    if (contrAddr!=common.Address{} && certCnt!= nil && parentAddr==common.Address{}) {
        parentAddr := common.Address{}
        if len(certCnt)< (gCaCertOffset+len( parentAddr.Bytes() )){
            http.Error(w, GeneralError{fmt.Sprintf(
                "EnrollUser: Certificate is too small")}.Error(),
                http.StatusInternalServerError)
            return
        }
        // TO DO: ADD a check if the chunk below is an address with common.isAddress (isHex)
        parentAddr.SetBytes( certCnt[gCaCertOffset : gCaCertOffset+len( contrAddr.Bytes() )] )
    }
    */

    if (parentAddr == common.Address{}){
        http.Error(w, GeneralError{"Enroll: Parent address is not established"}.Error(),
            http.StatusInternalServerError)
        return
    }

    var desc string = ""
    descArr := r.MultipartForm.Value["Desc"]
    if (len(descArr) > 0) {
        desc = descArr[0];
    }
    fmt.Printf("DEBUG before newRegDatum: hash=%v, fname=%v, desc=%v, userAddr=%v\n",
        strconv.Itoa(int(hashSum)), fileName, desc, userAddr.String())

    client, err := ethclient.Dial(gConfig.IPCpath)
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }

    // Instantiate the contract, the address is taken from eth at the moment of contract initiation
    // kyc, err := NewLuxUni_KYC(common.HexToAddress(gContractHash), backends.NewRPCBackend(conn))
    pkiContract, err := NewLuxUni_PKI(parentAddr, client)
    if err != nil {
        log.Fatalf("Failed to instantiate a smart contract: %v", err)
    }

    callOpts:= &bind.CallOpts{
        Pending: true,
    }
    initNumRegData,err := pkiContract.GetNumRegData(callOpts)
    if err != nil {
        http.Error(w, fmt.Sprintf("EnrollUser: Failed to get numRegData from blockchain: %v. ", err),
            http.StatusInternalServerError)
        return
    }

    // Logging into Ethereum as a user
    if (userAddr == common.Address{}) {
        fmt.Printf("Attention! Enroll: user address is zero, default config account is used\n")
        userAddr = common.HexToAddress(gConfig.AccountAddr)
    }
    keyFile, err := FindKeyFile(userAddr)
    if (err!=nil) {
        http.Error(w, fmt.Sprintf("Failed to find key file for account %v. %v ",
            userAddr.String(), err), http.StatusInternalServerError)
        return
    }
    key, e := ioutil.ReadFile( gConfig.KeyDir+keyFile )
    if e != nil {
        fmt.Printf("Key File error: %v\n", e)
        os.Exit(1)
    }
    fmt.Printf("Found Ethereum Key File \n")

    auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
    if err != nil {
        log.Fatalf("Failed to create authorized transactor: %v", err)
    }

    sess := &LuxUni_PKISession{
        Contract: pkiContract,
        CallOpts: bind.CallOpts{
            Pending: true,
        },
        TransactOpts: bind.TransactOpts{},
    }
    sess.TransactOpts = *auth;
    sess.TransactOpts.GasLimit = big.NewInt(2000000)

    res, err := sess.NewRegDatum( []byte(strconv.Itoa(int(hashSum))), contrAddr,
            fileName, desc, "", userAddr)
    if err != nil {
        http.Error(w, fmt.Sprintf("EnrollUser: Failed to add a record to blockchain: %v. ", err),
            http.StatusInternalServerError)
        return
    }

    finalNumRegData,err := pkiContract.GetNumRegData(callOpts)
    if err != nil {
        http.Error(w, fmt.Sprintf("EnrollUser: Failed to get numRegData from blockchain: %v. ", err),
            http.StatusInternalServerError)
        return
    }

    if( finalNumRegData.Int64() != initNumRegData.Int64()+1 ){
        http.Error(w, fmt.Sprintf("EnrollUser: Failed to add a record, wrong function return: %v",
            res.Value().Int64()), http.StatusInternalServerError)
        return
    }
    //hashResult = strconv.Itoa(i+1) + " file(s) processed:"

    // UplFile is id in the input "file" component of the form
    // http://stackoverflow.com/questions/33771167/handle-file-uploading-with-go
    // file, handler, err := r.FormFile("UplFile")
    //out, err := os.Create("/tmp/tst_"+handler.Filename);

    hashResult := "Hash info is successfully loaded to blockchain.";

    if (gConfig.WebMode!=0){
        var formParam CSimpleFormParam
        //formParam.Param = ""
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

    var formParam CPKIFormParam;
    var revokedParam CRevokedFormParam;
    var isRevokeListRequest bool;
    var isCurl []string;
    //var addrCAs []common.Address
    var parentAddr common.Address
    var superParentAddr common.Address
    var superUserAddr common.Address
    //var isChangeData bool = false;

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("PKI Main form: error in parsing multipart form: %v\n", err.Error())
        parentAddr = common.HexToAddress(gConfig.ContractHash)
    } else {
        strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
        if( len(strParentAddrArr)>0 ){
            if(common.IsHexAddress(strParentAddrArr[0])==false){
                http.Error(w, GeneralError{fmt.Sprintf(
                    "Parent address is incorrect: ",strParentAddrArr[0])}.Error(),
                    http.StatusInternalServerError)
                return
            }
            parentAddr = common.HexToAddress(strParentAddrArr[0])
        }

        strSuperParentAddrArr := r.MultipartForm.Value["SuperParentAddr"]
        if( len(strSuperParentAddrArr)>0 ){
            if(common.IsHexAddress(strSuperParentAddrArr[0])==false){
                http.Error(w, GeneralError{fmt.Sprintf(
                    "SuperParent address is incorrect: ",strSuperParentAddrArr[0])}.Error(),
                    http.StatusInternalServerError)
                return
            }
            superParentAddr = common.HexToAddress(strSuperParentAddrArr[0])
        }

        strSuperUserAddrArr := r.MultipartForm.Value["UserAddr"]
        if( len(strSuperUserAddrArr)>0 ){
            if(common.IsHexAddress(strSuperUserAddrArr[0])==false){
                http.Error(w, GeneralError{fmt.Sprintf(
                    "User address is incorrect: ",strSuperUserAddrArr[0])}.Error(),
                    http.StatusInternalServerError)
                return
            }
            superUserAddr = common.HexToAddress(strSuperUserAddrArr[0])
        }

        strRevokeListRequest := r.MultipartForm.Value["RevokeListButton"]
        if len(strRevokeListRequest)>0 {
            isRevokeListRequest = true;
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
    strParentAddr := r.FormValue("ParentAddr");
    if(len(strParentAddr)>0){
        if(common.IsHexAddress(strParentAddr)==true){
            parentAddr = common.HexToAddress(strParentAddr)
        }else{
            http.Error(w, GeneralError{fmt.Sprintf(
                "Parent address is incorrect: ",strParentAddr)}.Error(),
                http.StatusInternalServerError)
            return
        }
    }

    if(parentAddr == common.Address{}){
        /*http.Error(w, GeneralError{"Parent address is nil: "}.Error(),
            http.StatusInternalServerError)
        return*/
        parentAddr = common.HexToAddress(gConfig.ContractHash);
    }
    fmt.Printf("Contract address in Main Form: %v\n", common.ToHex(parentAddr.Bytes()))

    if(superUserAddr == common.Address{}){
        //fmt.Printf("Attention! pki-form: user address is zero, default config account is used\n")
        superUserAddr = common.HexToAddress(gConfig.AccountAddr);
    }
    fmt.Printf("Contract address in Main Form: %v\n", common.ToHex(parentAddr.Bytes()))

    client, err := ethclient.Dial(gConfig.IPCpath)
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }

    pkiContr, err := NewLuxUni_PKI(parentAddr, client)
    if err != nil {
        http.Error(w, GeneralError{fmt.Sprintf(
            "Failed to instantiate a smart contract: ", err)}.Error(),
            http.StatusInternalServerError)
        return
    }
    callOpts:= &bind.CallOpts{
        Pending: true,
    }

    numRD, err := pkiContr.GetNumRegData(callOpts)
    if err != nil {
        log.Fatalf("Failed to retrieve a total number of data records: %v", err)
    }
    fmt.Println("Number of data records (including those deleted)", numRD)

    for i := int64(0); i < numRD.Int64(); i++ {
        bi := big.NewInt(i)
        var regDatum CRegData
        regDatum.ContrAddr, err = pkiContr.GetRegContrAddr(callOpts, bi)
        if err != nil {
            http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
            return
        }
        regDatum.NodeSender, err = pkiContr.GetRegNodeSender(callOpts, bi)
        if err != nil {
            http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
            return
        }
        regDatum.DataHash, err = pkiContr.GetRegDataHash(callOpts, bi)
        if err != nil {
            http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
            return
        }
        regDatum.CreationDate, err = pkiContr.GetRegCreationDate(callOpts, bi)
        if err != nil {
            http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
            return
        }
        regDatum.FileName, err = pkiContr.GetRegFileName(callOpts, bi)
        if err != nil {
            http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
            return
        }
        regDatum.LinkFile, err = pkiContr.GetRegLinkFile(callOpts, bi)
        if err != nil {
            http.Error(w, fmt.Sprintf("Data is not retrieved: %v", err), http.StatusInternalServerError)
            return
        }
        regDatum.Description, err = pkiContr.GetRegDescription(callOpts, bi)
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

        var strContrAddr string;
        var isCA bool;
        if (regDatum.ContrAddr != common.Address{}) {
            //addrCAs = append(addrCAs, regDatum.ContrAddr)
            //formParam.AddrCAs = append(formParam.AddrCAs, common.ToHex(regDatum.ContrAddr.Bytes()))
            strContrAddr = common.ToHex(regDatum.ContrAddr.Bytes())
            isCA = true
        }
        crDate := time.Unix(regDatum.CreationDate.Int64(), 0);

        docPrez := CDocPrez{ int(i), regDatum.FileName, parentAddr.String(),
                       strContrAddr, regDatum.NodeSender.String(),
                       regDatum.Description, ""/* link */,
                       "" /*decrypt*/, string(regDatum.DataHash), crDate,
                       crDate.String(), isCA}

        // formation of data for (presentation of) HTML forms
        if delRegDate.Int64() == 0 {
            formParam.Docs = append(formParam.Docs, docPrez )
            //fmt.Printf("Data received: %v, %v \n", regDatum.Description, crDate.String()) 
        } else {
            revokedParam.Docs = append(revokedParam.Docs, docPrez )
        }
    }

    // CURL does neet any web return - just "OK" of no errors
    if (isCurl != nil) {
        jsonResp := `{"verbose":"OK","result":0}`
        fmt.Fprintln(w, jsonResp)
        return;
    }

    if ( (len(revokedParam.RevokedIds) > 0) || ( isRevokeListRequest != false ) ) {
        tmpl := mTempl["HashResult"]
        terr := tmpl.Execute(w, revokedParam)
        if terr != nil {
            http.Error(w, terr.Error(), http.StatusInternalServerError)
        }
    }

    tmpl := mTempl["MainForm"]
    formParam.ParentAddr=common.ToHex(parentAddr.Bytes())
    formParam.SuperParentAddr=common.ToHex(superParentAddr.Bytes())
    formParam.SuperUserAddr=common.ToHex(superUserAddr.Bytes())
    terr := tmpl.Execute(w, formParam)
    if terr != nil {
        http.Error(w, terr.Error(), http.StatusInternalServerError)
    }
}


func main() {
    
    if mTempl == nil { return; }
    
    gPrivateKey = new(rsa.PrivateKey);
    err := LoadPrivateKey(gConfig.PrivateKeyPath, gPrivateKey)
    if err != nil {
        gPrivateKey = nil;
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
    r := mux.NewRouter();
    r.HandleFunc("/pki-test", PkiForm);
    r.HandleFunc("/enroll_user", EnrollUser);
    r.HandleFunc("/blacklist_user", BlacklistUser);
    //r.HandleFunc("/enroll_ca", EnrollCA);
    r.HandleFunc("/create_contract", CreateContract);
    r.HandleFunc("/populate_contract", PopulateContract);
    r.HandleFunc("/validate_form", ValidateForm);
    r.HandleFunc("/validate_cert", ValidateCert);
    r.HandleFunc("/download_cacert", DownloadCaCert);
    r.HandleFunc("/generate_user_cert", GenerateUserCert);

    fs := http.FileServer(http.Dir("/home/alex/DocsFS/Dropbox/WORK/RD/LuxBCh/PKI/public"));
    spref := http.StripPrefix("/public/", fs);
    r.PathPrefix("/public/").Handler(spref);
    http.Handle("/", r);

    //https://gist.github.com/denji/12b3a568f092ab951456 - SSL info
    //https://golanglibs.com/top?q=webrtc - webrtc server side for golang

    //var server = &http.Server{
    //    Addr : ":8071",
    //    Handler : r,
    //}

    log.Println("Listening...")
    //http.ListenAndServeTLS(":8071", "server.pem", "server.key", r)
    http.ListenAndServe(":"+strconv.Itoa(gConfig.HttpPort), nil)
}

func PopulateContract(w http.ResponseWriter, r *http.Request) {

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
        return;
    }

    //isCurl = r.MultipartForm.Value["Curl"]
    if len(r.MultipartForm.Value["ContrAddr"])==0 {
        http.Error(w, GeneralError{"No contrAddr is provided"}.Error(),
            http.StatusInternalServerError)
        return
    }
    contrAddrStr := r.MultipartForm.Value["ContrAddr"][0]
    if(common.IsHexAddress(contrAddrStr)==false){
        http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
            http.StatusInternalServerError)
        return
    }
    contrAddr := common.HexToAddress(contrAddrStr)

    parentAddr := common.Address{}
    if len(r.MultipartForm.Value["ParentAddr"])!=0 {
        parentAddrStr := r.MultipartForm.Value["ParentAddr"][0]
        if (common.IsHexAddress(parentAddrStr) == false) {
            http.Error(w, GeneralError{"Parent contract address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        parentAddr = common.HexToAddress(parentAddrStr)
    }

    newUserAddr := common.Address{}
    if len(r.MultipartForm.Value["NewUserAddr"])!=0 {
        userAddrStr := r.MultipartForm.Value["NewUserAddr"][0]
        if (common.IsHexAddress(userAddrStr) == false) {
            http.Error(w, GeneralError{"New User address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        newUserAddr = common.HexToAddress(userAddrStr)
    } else {
        http.Error(w, "New User address is not available in params", http.StatusInternalServerError)
        return
    }

    curUserAddr := common.Address{}
    if len(r.MultipartForm.Value["CurrentUserAddr"])!=0 {
        userAddrStr := r.MultipartForm.Value["CurrentUserAddr"][0]
        if (common.IsHexAddress(userAddrStr) == false) {
            http.Error(w, GeneralError{"Current User address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        curUserAddr = common.HexToAddress(userAddrStr)
    } else {
        http.Error(w, "Current User address is not available in params", http.StatusInternalServerError)
        return
    }


    /*certHash, _, fileCont, cerr := UploadFile(w,r,"UplFiles",
                common.HexToAddress(parentAddrStr), true)
    if(cerr.errCode!=0){
        fmt.Printf(fmt.Sprintf("Populate Uploadfile: %v\n", cerr.Error()))
        http.Error(w, cerr.Error(), http.StatusInternalServerError)
        return
    }*/
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

    client, err := ethclient.Dial(gConfig.IPCpath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Instantiate the contract, the address is taken from eth at the moment of contract initiation
    // kyc, err := NewLuxUni_KYC(common.HexToAddress(gContractHash), backends.NewRPCBackend(conn))
    pkiContract, err := NewLuxUni_PKI(contrAddr, client)
    if err != nil {
        http.Error(w, GeneralError{
            fmt.Sprintf("Failed to instantiate a smart contract: %v", err)}.Error(),
            http.StatusInternalServerError)
        return
    }

    // Logging into Ethereum as a user
    if (curUserAddr == common.Address{}) {
        fmt.Printf("Attention! Populate contract: user address is zero, default config account is used\n")
        curUserAddr = common.HexToAddress(gConfig.AccountAddr)
    }
    keyFile, err := FindKeyFile(curUserAddr)
    if (err!=nil) {
        http.Error(w, fmt.Sprintf("Failed to find key file for account %v. %v ",
            curUserAddr.String(), err), http.StatusInternalServerError)
        return
    }
    key, e := ioutil.ReadFile( gConfig.KeyDir+keyFile )
    if e != nil {
        fmt.Printf("Key File error: %v\n", e)
        os.Exit(1)
    }
    fmt.Printf("Found Ethereum Key File \n")

    auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
    if err != nil {
        log.Fatalf("Failed to create authorized transactor: %v", err)
    }

    sess := &LuxUni_PKISession{
        Contract: pkiContract,
        CallOpts: bind.CallOpts{
            Pending: true,
        },
        TransactOpts: bind.TransactOpts{},
    }
    sess.TransactOpts = *auth;
    sess.TransactOpts.GasLimit = big.NewInt(50000000)

    _, err = sess.PopulateCertificate(dataCert)
    if err != nil {
        fmt.Printf(fmt.Sprintf("Failed to populate blockchain: %v.\n", err))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if(newUserAddr != common.Address{}){
        _, err := sess.SetOwner(newUserAddr)
        if err != nil {
            fmt.Printf(fmt.Sprintf("Failed to update owner addr: %v.\n", err))
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    } else {
        http.Error(w, "New User addr is null", http.StatusInternalServerError)
        return
    }

    resStr := "Contract was successfully populated\n"
    detailStr := "Contract Addr = "+contrAddrStr+"\n"
    detailStr += "Parent Addr = "+common.ToHex(parentAddr.Bytes())+"\n"
    detailStr += "New User = "+common.ToHex(newUserAddr.Bytes())+"\n"
    detailStr += "Current User = "+common.ToHex(curUserAddr.Bytes())+"\n"
    if (gConfig.WebMode!=0){
        var formParam CSimpleFormParam
        formParam.Params = append(formParam.Params,
            CParam{"Hash",strconv.Itoa(int(hashCert))})
        formParam.Params = append(formParam.Params,
            CParam{"ContrAddr",contrAddrStr})
        formParam.Params = append(formParam.Params,
            CParam{"ParentAddr",common.ToHex(parentAddr.Bytes())})
        formParam.Params = append(formParam.Params,
            CParam{"UserAddr",common.ToHex(newUserAddr.Bytes())})
        formParam.Result = resStr
        formParam.Data = detailStr
        formParam.CallForm = "enroll_user"

        tmpl := mTempl["SimpleForm"]
        terr := tmpl.Execute(w, formParam)
        if terr != nil {
            http.Error(w, terr.Error(), http.StatusInternalServerError)
        }
    }
}

func CreateContract(w http.ResponseWriter, r *http.Request){
    /*
    https://vincentserpoul.github.io/post/binding-ethereum-golang/
    https://ethereum.stackexchange.com/questions/7499/how-are-addresses-created-if-deploying-a-new-bound-contract
    */
    var parentAddrStr string;
    var curUserAddrStr string;   // !!! presently current user not used - addr=contr.GetOwner used instead
    var newUserAddrStr string;

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("CreateContract: No data in multipart form: %v\n", err.Error())
        parentAddrStr = gConfig.ContractHash
    } else {
        strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
        if len(strParentAddrArr)>0 {
            parentAddrStr = strParentAddrArr[0];
            if(common.IsHexAddress(parentAddrStr)==false){
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
    if len(strUserAddrArr)>0 {
        curUserAddrStr = strUserAddrArr[0];
        if(common.IsHexAddress(curUserAddrStr)==false) {
            fmt.Println("Create Contract: Current user address is incorrect")
            http.Error(w, GeneralError{"Current user address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
    }

    strUserAddrArr = r.MultipartForm.Value["NewUserAddr"]
    if len(strUserAddrArr)>0 {
        newUserAddrStr = strUserAddrArr[0];
        if(common.IsHexAddress(newUserAddrStr)==false) {
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
    callOpts:= &bind.CallOpts{
        Pending: true,
    }
    execUserAddr,err := pkiContr.GetOwner(callOpts)
    if err != nil {
        http.Error(w, GeneralError{fmt.Sprintf(
            "CreateCont - failed to get owner addr: ", err)}.Error(),
            http.StatusInternalServerError)
        return
    }
    if execUserAddr!=common.HexToAddress(curUserAddrStr) {
        http.Error(w, "contract.GetOwner does not correspond to the Current User param",
            http.StatusInternalServerError)
        return
    }

    //keyFile := gConfig.KeyFile
    keyFile, err := FindKeyFile(execUserAddr)
    if (err!=nil) {
        http.Error(w, fmt.Sprintf("CreateContract -- FindKeyFile: %v. ", err),
            http.StatusInternalServerError)
        return
    }
    key, e := ioutil.ReadFile( gConfig.KeyDir+keyFile )
    if e != nil {
        fmt.Printf("Key File error: %v\n", e)
        os.Exit(1)
    }
    fmt.Printf("Found Ethereum Key File \n")

    auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
    if err != nil {
        log.Fatalf("Failed to create authorized transactor: %v", err)
    }
    var trOpts bind.TransactOpts = *auth;
    trOpts.GasLimit = big.NewInt(50000000)
    contrAddr, _, /*contr*/_, err := DeployLuxUni_PKI(&trOpts, client)
    /*
    https://stackoverflow.com/questions/40096750/set-status-code-on-http-responsewriter
    */
    if err!=nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("500 - Something bad happened!"))
    }

    if (gConfig.WebMode!=0) {
        var formParam CSimpleFormParam
        formParam.IsUploadFileForm = false
        formParam.Params = append(formParam.Params,
            CParam{"ContrAddr",common.ToHex(contrAddr.Bytes())} )
        formParam.Params = append(formParam.Params,
            CParam{"ParentAddr",parentAddrStr})
        formParam.Params = append(formParam.Params,
            CParam{"CurrentUserAddr",curUserAddrStr})
        formParam.Params = append(formParam.Params,
            CParam{"NewUserAddr",newUserAddrStr})
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
        return;
    }

    //isCurl = r.MultipartForm.Value["Curl"]
    //certHash := r.MultipartForm.Value["CertHash"][0]
    if len(r.MultipartForm.Value["ContrAddr"])==0 {
        http.Error(w, GeneralError{"No contrAddr is provided"}.Error(),
            http.StatusInternalServerError)
        return
    }
    strContrAddr := r.MultipartForm.Value["ContrAddr"][0]
    if(common.IsHexAddress(strContrAddr)==false){
        http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
            http.StatusInternalServerError)
        return
    }
    contrAddr:=common.HexToAddress(strContrAddr)

    if (gConfig.WebMode!=0) {
        var formParam CSimpleFormParam
        formParam.IsUploadFileForm = true
        formParam.Params = append(formParam.Params,
            CParam{"ContrAddr",common.ToHex(contrAddr.Bytes())} )
        //formParam.Params = append(formParam.Params,
        //    CParam{"ParentAddr",parentAddrStr})
        formParam.Result = "Please upload a certificate to check for CA "+
            common.ToHex(contrAddr.Bytes())
        formParam.CallForm = "validate_cert"

        tmpl := mTempl["SimpleForm"]
        terr := tmpl.Execute(w, formParam)
        if terr != nil {
            http.Error(w, terr.Error(), http.StatusInternalServerError)
        }
    }
}

func GenerateUserCert(w http.ResponseWriter, r *http.Request){

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
        return;
    }

    //isCurl = r.MultipartForm.Value["Curl"]
    //certHash := r.MultipartForm.Value["CertHash"][0]
    if len(r.MultipartForm.Value["InsertAddr"])==0 {
        http.Error(w, GeneralError{"No insertAddr is provided"}.Error(),
            http.StatusInternalServerError)
        return
    }
    strInsertAddr := r.MultipartForm.Value["InsertAddr"][0]
    if(common.IsHexAddress(strInsertAddr)==false){
        http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
            http.StatusInternalServerError)
        return
    }
    insertAddr:=common.HexToAddress(strInsertAddr)

    if len(r.MultipartForm.Value["Name"])==0 {
        http.Error(w, GeneralError{"No user name is provided"}.Error(),
            http.StatusInternalServerError)
        return
    }
    strName := r.MultipartForm.Value["Name"][0]

    dataCert, err := GenerateCert(common.Address{}, insertAddr,false, strName)
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
func DownloadCaCert(w http.ResponseWriter, r *http.Request){

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
        return;
    }

    if len(r.MultipartForm.Value["ContrAddr"])==0 {
        http.Error(w, GeneralError{"No contrAddr is provided"}.Error(),
            http.StatusInternalServerError)
        return
    }
    strContrAddr := r.MultipartForm.Value["ContrAddr"][0]
    if(common.IsHexAddress(strContrAddr)==false){
        http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
            http.StatusInternalServerError)
        return
    }
    contrAddr:=common.HexToAddress(strContrAddr)

    /*isCertOK*/ _, /*revokDate*/_, /*parentAddr*/_, /*retCaHash*/_, certData, err :=
        ConfirmHashCAData(contrAddr, 0, true)

    w.Header().Set("Content-Disposition", "attachment; filename=ca.crt.out")
    w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
    http.ServeContent(w, r, "ca.crt.out", time.Now(), bytes.NewReader(certData))
}


func ValidateCert(w http.ResponseWriter, r *http.Request){

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
        return;
    }

    certHash, _, dataCert, cerr := UploadFile(w,r,"UplFiles", common.Address{},true)
    if(cerr.errCode!=0){
        http.Error(w, cerr.Error(),
            http.StatusInternalServerError)
        return
    }

    contrAddr, parentAddr, err:= ParseCert(dataCert)
    if(err!=nil){
        http.Error(w, fmt.Sprintf("CERTIFICATE: Parsing error: %v",err),
            http.StatusInternalServerError)
        return
    }
    if (parentAddr==common.Address{}) {
        http.Error(w, "CERTIFICATE: no parent address in the certificate is provided",
            http.StatusInternalServerError)
        return
    }
    if (contrAddr!=common.Address{}) {
        http.Error(w, "CERTIFICATE: Non-CA certificates should not include own contract address data",
            http.StatusInternalServerError)
        return
    }

    //isCurl = r.MultipartForm.Value["Curl"]
    //certHash := r.MultipartForm.Value["CertHash"][0]
    /* // ContrAddr SHOULD BE SPECIFIED IN THE CONTRACT ITSELF
    if len(r.MultipartForm.Value["ContrAddr"])==0 {
        http.Error(w, GeneralError{"No contrAddr is provided"}.Error(),
            http.StatusInternalServerError)
        return
    }
    strContrAddr := r.MultipartForm.Value["ContrAddr"][0]
    if(common.IsHexAddress(strContrAddr)==false){
        http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
            http.StatusInternalServerError)
        return
    }
    contrAddr:=common.HexToAddress(strContrAddr)
    */

    isCertOK, revokeDate, jsonCertTree, iter, err := CheckCertTree(parentAddr, certHash)
    if(err != nil){
        http.Error(w, err.Error(),
            http.StatusInternalServerError)
        return
    }
    //https://stackoverflow.com/questions/19038598/how-can-i-pretty-print-json-using-go
    var prettyJSON bytes.Buffer
    err = json.Indent(&prettyJSON, []byte(jsonCertTree), "", "\t")
    if err != nil {
        http.Error(w, fmt.Sprintf("ValidateCert: Json parse: %v",err.Error()),
            http.StatusInternalServerError)
        return
    }

    revokeDateStr:="No"
    if (revokeDate!=time.Time{}) {
        revokeDateStr="Yes, a certificate was revoked at "+ strconv.Itoa(iter) +
            " level above on "+revokeDate.String()
    }
    retDataStr := "Is the certificate valid? " + strconv.FormatBool(isCertOK) + "\n"
    retDataStr += "Was there revocation (date)? " + revokeDateStr + "\n"
    retDataStr += "Level of the initial certificate (check level iterations)? " +
        strconv.Itoa(iter) + "\n"
    retDataStr += "\nJson path of the certificate check: \n" + string(prettyJSON.Bytes())

    if (gConfig.WebMode!=0) {
        var formParam CSimpleFormParam
        formParam.IsUploadFileForm = false
        //formParam.Params = append(formParam.Params,
        //    CParam{"ContrAddr",common.ToHex(contrAddr.Bytes())} )
        //formParam.Params = append(formParam.Params,
        //    CParam{"ParentAddr",parentAddrStr})
        formParam.Result = "Results of certificate check (certHash="+
            strconv.Itoa(int(certHash))+") for CA "+contrAddr.String()
        formParam.Data = retDataStr
        formParam.CallForm = "pki-test"

        tmpl := mTempl["SimpleForm"]
        terr := tmpl.Execute(w, formParam)
        if terr != nil {
            http.Error(w, terr.Error(), http.StatusInternalServerError)
        }
    }
}

/*
  returns Json string with the path of certificates to the root
          int with the number of iteractions
 */
func CheckCertTree(parentAddr common.Address, userHash uint32) (retIsCertOK bool,
         retRevokeDate time.Time, retJson string, retIter int, err error) {
    //addr := common.HexToAddress(gConfig.ContractHash);
    type JsonNode struct {
        ContrAddr string
        ParentAddr string
        Hash string
        RevokeDate string
        IsCertOK string
        Warn string
    }
    var jsonPath []JsonNode;
    var maxIter int = 1000;

    iterHash := userHash;
    for retIter=0; retIter < maxIter; retIter++ {
        var jsonNode JsonNode
        jsonNode.ContrAddr=parentAddr.String();
        jsonNode.Hash=strconv.Itoa(int(iterHash));

        retIsCertOK, retRevokeDate, parentAddr, iterHash, _, err =
            ConfirmHashCAData(parentAddr,iterHash,false)
        if err != nil {
            return false, time.Time{}, "", retIter, err
        }

        jsonNode.ParentAddr=parentAddr.String();
        jsonNode.IsCertOK=strconv.FormatBool(retIsCertOK)
        jsonNode.RevokeDate=retRevokeDate.String()
        jsonPath = append(jsonPath, jsonNode)
        if (retIsCertOK == false) {
            break
        }
        if (parentAddr == common.Address{}) {
            break
        }
        if (retIter>=(maxIter-1)) && (parentAddr != common.Address{}) {
            return false, time.Time{}, "", retIter,GeneralError{"MaxIter limit is reached"}
        }
    }
    //bJson, err := json.Marshal(jsonPath)
    bJson, err := json.Marshal(jsonPath)
    if err != nil {
        return false, time.Time{}, "", retIter, GeneralError{fmt.Sprintf("Json Marshal:",err)}
    }
    retJson = string(bJson)
    return retIsCertOK, retRevokeDate, retJson, retIter, nil
}

/*
    params:
      hash - hash to check in the currAddr, if hash==0, no user certificate check is conducted
             and the info on the present CA is returned
      isGetCertData - if true the CA cert data is returned in the []byte variable
   returns:
      isHashFound: - is the hash available in this contract:
      revokeDate: date of revocation or zero if no revocation
      parentAddr,
      caHash of the present CA certificate
      []byte -- data of the CA certificate if the isGetCertData == true, otherwise nil
 */
func ConfirmHashCAData(currAddr common.Address, hash uint32,
        isGetCaCertData bool) (retIsHashFound bool, retRevokeDate time.Time,
             retParentAddr common.Address, retCaHash uint32, retCertData []byte, err error) {

    retIsHashFound = false;
    retRevokeDate = time.Time{};
    retCertData = nil;

    client, err := ethclient.Dial(gConfig.IPCpath)
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }

    currContract, err := NewLuxUni_PKI(currAddr, client)
    if err != nil {
        log.Fatalf("Failed to instantiate a smart contract: %v", err)
    }
    callOpts:= &bind.CallOpts{
        Pending: true,
    }

    //caCert, err := currContract.NumRegData(callOpts)
    caCert, err := currContract.GetCaCertificate(callOpts)
    if err!=nil {
        return false, time.Time{}, common.Address{}, 0, nil, err
    }
    if isGetCaCertData == true {
        retCertData = caCert
    }
    //parentAddr := caCert[gCaCertOffset:gCaCertOffset+len(common.Address{}.Bytes())]
    contrAddr, retParentAddr, err := ParseCert(caCert);
    if err!= nil {
        return false, time.Time{}, common.Address{}, 0, nil, err
    }
    if contrAddr != currAddr {
        return false, time.Time{}, common.Address{}, 0, nil,
            GeneralError{"GetParent: contrAddr does not correspond to the CA smart contract"}
    }
    retCaHash, err = CalcHash(caCert)
    if err!= nil {
        return false, time.Time{}, common.Address{}, 0, nil, err
    }

    if hash==0 {
        return false, time.Time{}, retParentAddr, retCaHash, retCertData, nil;
    }

    numRD, err := currContract.GetNumRegData(callOpts)
    if err!=nil {
        return false, time.Time{}, common.Address{}, 0, nil, err
    }

    for i := int64(0); i < numRD.Int64(); i++ {
        bi := big.NewInt(i)
        regDataHash, err := currContract.GetRegDataHash(callOpts, bi)
        if err != nil {
            return false, time.Time{}, common.Address{}, 0, nil,
                    GeneralError{fmt.Sprintf("Failed to get hash data: ", err)}
        }
        delRegDate, err := currContract.GetDeletedRegDate(callOpts, bi)
        if err != nil {
            return false, time.Time{}, common.Address{}, 0, nil,
                    GeneralError{fmt.Sprintf("Failed to get deleted data: ", err)}
        }
        blchHash, err := strconv.Atoi(string(regDataHash))
        if err != nil {
            return false, time.Time{}, common.Address{}, 0, nil,
                    GeneralError{fmt.Sprintf("Hash string is not parsable: ", err)}
        }
        if(uint32(blchHash) == hash){
            if delRegDate.Int64() != 0 {
                retRevokeDate = time.Unix(delRegDate.Int64(), 0)
            }else{
                retIsHashFound = true
            }
            break
        }
    }

    return retIsHashFound, retRevokeDate, retParentAddr, retCaHash, retCertData, err
}


func GenerateCert(contrAddr common.Address, parentAddr common.Address,
          isCA bool, strName string) ([]byte, error){

    // see Certificate structure at
    // http://golang.org/pkg/crypto/x509/#Certificate
    template := &x509.Certificate {
        IsCA : true,
        BasicConstraintsValid : true,
        SubjectKeyId : []byte{1,2,3},
        SerialNumber : big.NewInt(1234),
        Subject : pkix.Name{
            Country : []string{"Earth"},
            Organization: []string{strName},
        },
        NotBefore : time.Now(),
        NotAfter : time.Now().AddDate(5,5,5),
        // see http://golang.org/pkg/crypto/x509/#KeyUsage
        ExtKeyUsage : []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
        KeyUsage : x509.KeyUsageDigitalSignature|x509.KeyUsageCertSign,
    }
    // https://stackoverflow.com/questions/26441547/go-how-do-i-add-an-extension-subjectaltname-to-a-x509-certificate
    template.ExtraExtensions = append(template.ExtraExtensions, // blockchain name
        GenerateCertExtension(asn1.ObjectIdentifier{1,2,752,115,33,3},
            common.Hex2Bytes("0x130d457468657265756d2d54657374"), false ) )
    if(isCA == true) {
        template.ExtraExtensions = append(template.ExtraExtensions,
            GenerateCertExtension(asn1.ObjectIdentifier{1,2,752,115,33,2},
                contrAddr.Bytes(), true ) )
    }
    template.ExtraExtensions = append(template.ExtraExtensions,
        GenerateCertExtension(asn1.ObjectIdentifier{1,2,752,115,33,1},
            parentAddr.Bytes(), true ) )
    template.ExtraExtensions = append(template.ExtraExtensions, // hash algo
        GenerateCertExtension(asn1.ObjectIdentifier{1,2,752,115,33,0},
            common.Hex2Bytes("0x0609608648016503040201"), false ) )

    // generate private key
    privatekey, err := rsa.GenerateKey(cryptoRand.Reader, 2048)

    if err != nil {
        return nil, err
    }

    publickey := &privatekey.PublicKey

    // create a self-signed certificate. template = parent
    var parent = template
    cert, err := x509.CreateCertificate(cryptoRand.Reader, template, parent, publickey,privatekey)

    if err != nil {
        return nil, err
    }

    return cert, nil

}

func GenerateCertExtension(ID asn1.ObjectIdentifier, val []byte,
    isAddPrefix bool) (pkix.Extension) {
    var extVal []byte
    if isAddPrefix==true {
        extVal = []byte{4,20}
    }
    extVal = append(extVal, val...)
    ext := pkix.Extension{}
    ext.Critical = false
    ext.Id = ID
    ext.Value = extVal;
    //template.ExtraExtensions = append(template.ExtraExtensions, ext)
    return ext
}

/*
  returns contrAddr, parentAddr
 */
func ParseCert(binCert []byte) (common.Address, common.Address, error) {
    var contrAddr, parentAddr common.Address
    ca, err := x509.ParseCertificate(binCert)
    if err!=nil {
        return common.Address{}, common.Address{}, err
    }
    // iterate in the extension to get the information
    for _, element := range ca.Extensions {
        if element.Id.String() == "1.2.752.115.33.2" { // CA Address
            fmt.Printf("\tCaContractIdentifier: %+#+x\n", element.Value)
            val:=element.Value[2:]
            if( len(val) != len(common.Address{}.Bytes()) ) {
                return common.Address{}, common.Address{},
                    GeneralError{"ParseCert: wrong length of CA addr"}
            }
            contrAddr = common.BytesToAddress(val)
        }
        if element.Id.String() == "1.2.752.115.33.1" { //Parent Address
            fmt.Printf("\tIssuerCaContractIdentifier: %+#+x\n", element.Value)
            val:=element.Value[2:]
            if( len(val) != len(common.Address{}.Bytes()) ) {
                return common.Address{}, common.Address{},
                    GeneralError{"ParseCert: wrong length of CA addr"}
            }
            parentAddr = common.BytesToAddress(val)
        }
    }
    return contrAddr, parentAddr, nil
}

func CalcHash(data []byte) (uint32, error){
    hasher := crc32.NewIEEE()
    hasher.Write(data)
    return hasher.Sum32(), nil
}

func FindKeyFile(addr common.Address) (keyFileName string, err error) {
    files, err := ioutil.ReadDir(gConfig.KeyDir)
    if(err!=nil){
        return "", err
    }
    for _, f := range files {
        if( strings.Contains(f.Name(), addr.String()[2:]) == true ) {
            return f.Name(), nil
        }
    }
    return "", GeneralError{"Ethereum Key File not found for this address"}
}

func LoadConfig() (error) {
    file, err := os.Open( gConfigFile )
    if err != nil {
        return GeneralError{fmt.Sprintf("File error: %v\n", err)}
    }
    fmt.Printf("Found configuration file: %s\n", gConfigFile)

    //jsonparser.Get(data, "person", "name", "fullName")
    jsonParser := json.NewDecoder(file)
    if err = jsonParser.Decode( &gConfig ); err != nil {
        return GeneralError{fmt.Sprintf("Parsing config file: %s\n", err.Error())}
    }
    
    b, err := json.Marshal(gConfig)
    if err != nil {
        return GeneralError{fmt.Sprintf("Cannot convert conf file into string: %s", err)}
    }
    fmt.Printf("Loaded configuration file: %s\n", string(b) )
    file.Close()

    if common.IsHexAddress(gConfig.AccountAddr) == false {
        return GeneralError{"Config: account address is not correct"}
    }
    if common.IsHexAddress(gConfig.ContractHash) == false {
        return GeneralError{"Config: contract hash is not correct"}
    }

    return nil;
}

func LoadPrivateKey( path string, ptrPrivateKey *rsa.PrivateKey ) (error) {
    privatekeyfile, err := os.Open( path )
    if err != nil {
        return err;
    }

    decoder := gob.NewDecoder(privatekeyfile)
    err = decoder.Decode(ptrPrivateKey)
    if err != nil {
        return err;
    }
    privatekeyfile.Close()
    return nil;
}

func Decrypt(msg []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
    if privateKey == nil {
        err := GeneralError{"Private keys is not loaded"}
        return nil, err;
    }

    label := []byte("")  
    hash := sha256.New()

    return rsa.DecryptOAEP(hash, cryptoRand.Reader, privateKey, msg, label)
}

type GeneralError struct {
    errMsg string;
}
func (e GeneralError) Error() string {
    return e.errMsg
}

type GeneralCodeError struct {
    errMsg string;
    errCode uint64;
}
func (e *GeneralCodeError) Error() string {
    return e.errMsg
}


