
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
    //"math/rand"
    "hash/crc32"
    cryptoRand "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "encoding/gob"

    "github.com/gorilla/mux"

    //"github.com/ethereum/go-ethereum/rpc"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    //"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
    "github.com/ethereum/go-ethereum/ethclient"

    "encoding/json"
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
    KeyPath            string `json:"keyPath"`
    PrivateKeyPath     string `json:"privateKeyPath"`
    HttpPort           int `json:"httpPort"`
    WebMode            int `json:"webMode"`
    JsonMode           int `json:"jsonMode"`
}

type CRegData struct {
	NodeSender   common.Address
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
    Docs         []CDocPrez
    //AddrCAs      []string
}

type CDocPrez struct {
    Id           int
    Name         string
    ContrAddr    string
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
    Params   []CParam
    CallForm     string
    IsUploadFileForm bool
}
type CParam struct {
    Name string
    Value string
}

func init() {
    LoadConfig();
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
    var revokedParam CRevokedFormParam;
    var parentAddr common.Address = common.Address{};

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

    dels := r.MultipartForm.Value["Deletion"]
    if len(dels) > 0 {
        fmt.Printf("Debug: I am in deletion block")
        //dels := r.MultipartForm.Value["Deletion"]
        //dels := r.Form["Deletion"]
        for _, del := range dels {
            fmt.Printf("del=%v\n", del);
            delid, err := strconv.Atoi(del);
            if err != nil {
                log.Fatalf("Deletion conversion error: %v", err.Error())
            }
            revokedParam.RevokedIds = append(revokedParam.RevokedIds, delid);

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

            // Logging into Ethereum as a user
            //key []byte;
            key, e := ioutil.ReadFile( gConfig.KeyPath )
            if e != nil {
                fmt.Printf("Key File error: %v\n", e)
            }
            fmt.Printf("Found Ethereum Key File \n")

            auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
            if err != nil {
                fmt.Printf("Failed to create authorized transactor: %v", err)
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
                fmt.Printf("Deletion error: %v", nerr)
            }
        }
    }

    tmpl := mTempl["HashResult"]
    terr := tmpl.Execute(w, revokedParam)
    if terr != nil {
        http.Error(w, terr.Error(), http.StatusInternalServerError)
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
    var parentAddr common.Address = common.Address{}
    var contrAddr common.Address = common.Address{}

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("EnrollUser: No change data -- ", err.Error())
        http.Error(w, GeneralError{fmt.Sprintf(
            "EnrollUser: No change data -- ",err.Error())}.Error(),
            http.StatusInternalServerError)
        return
    }

    //isCurl = r.MultipartForm.Value["Curl"]

    hashSum, fileName, certCnt, cerr :=
            UploadFile(w, r, "UplFiles", common.Address{},true)
    if (cerr.errCode!=0){
        if(cerr.errCode == 1){
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
        if(common.IsHexAddress(strParentAddrArr[0])==false){
            http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        parentAddr = common.HexToAddress(strParentAddrArr[0])
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

    if (parentAddr == common.Address{}){
        http.Error(w, GeneralError{"Enroll: Parent address is not established"}.Error(),
            http.StatusInternalServerError)
        return
    }

    var desc string = ""
    var isEncrypt int64;
    descArr := r.MultipartForm.Value["Desc"]
    if (len(descArr) > 0) {
        desc = descArr[0];
    }
    fmt.Printf("DEBUG before newRegDatum: hash=%v, fname=%v, desc=%v, encrypt=%v\n",
        strconv.Itoa(int(hashSum)), fileName, desc, isEncrypt)

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

    // Logging into Ethereum as a user
    //key []byte;
    key, e := ioutil.ReadFile( gConfig.KeyPath )
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

    _, err = sess.NewRegDatum( []byte(strconv.Itoa(int(hashSum))), contrAddr,
        fileName, desc, "",
        big.NewInt(isEncrypt), common.Address{})
    if err != nil {
        log.Fatalf("Failed to add a record to blockchain: %v. ", err)
        fmt.Fprintln(w, err)
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

    numRD, err := pkiContr.NumRegData(callOpts)
    if err != nil {
        log.Fatalf("Failed to retrieve a total number of data records: %v", err)
    }
    fmt.Println("Number of data records (including those deleted)", numRD)

    for i := int64(0); i < numRD.Int64(); i++ {
        bi := big.NewInt(i)
        regDatum, err := pkiContr.RegData(callOpts, bi)
        if err != nil {
            log.Fatalf("Failed to get deleted data: %v", err)
        }

        delRegDatum, err := pkiContr.DeletedRegData(callOpts, bi)
        if err != nil {
            log.Fatalf("Deletion retrieval error: %v", err)
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
        var desc string
        var decrypt string

        /*if ( regDatum.Encrypted.Int64() != 0 ) {
            encryptRegDatum, err := session.EncryptRegData( bi)
            if err != nil {
                log.Fatalf("Encryption mapping retrieval error: %v", err)
            }
            if encryptRegDatum.EncryptDate.Int64() != 0 {
                desc = fmt.Sprintf("%v...", encryptRegDatum.Data[0:10]);
                bDecrypt, err := Decrypt(encryptRegDatum.Data, gPrivateKey)
                if err != nil {
                    decrypt = ""
                } else {
                    decrypt = string(bDecrypt)
                }
            }
        } else {*/
            desc = regDatum.Description;
        //}

        docPrez := CDocPrez{ int(i), regDatum.FileName, 
                       strContrAddr, desc, ""/* link */, decrypt,
                       string(regDatum.DataHash), crDate,
            crDate.String(), isCA}


        // scanning of recently revoked ids for formation of a string with results of revocation
        /*
        if (len(revokedParam.RevokedStr) == 0) {
            tmpStr := strconv.Itoa(len(revokedParam.RevokedIds)) + " key(s) revoked:"
            for iRev := 0; iRev < len(revokedParam.RevokedIds); iRev++ {
                if(revokedParam.RevokedIds[iRev] == docPrez.Id) {
                    tmpStr = tmpStr + " " + docPrez.Name;
                }
            }
            if (len(revokedParam.RevokedIds) > 0) {
                revokedParam.RevokedStr = tmpStr;
            }
        }
        */

        // formation of data for (presentation of) HTML forms
        if delRegDatum.DeletionDate.Int64() == 0 {
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

    /*
    // Create an IPC based RPC connection to a remote node
    // conn, err := rpc.NewIPCClient(gIPCpath)
    gClient, err = ethclient.Dial(gConfig.IPCpath)
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }

    // Instantiate the contract, the address is taken from eth at the moment of contract initiation
    // kyc, err := NewLuxUni_KYC(common.HexToAddress(gContractHash), backends.NewRPCBackend(conn))
    pkiContract, err := NewLuxUni_PKI(common.HexToAddress(gConfig.ContractHash), gClient)
    if err != nil {
        log.Fatalf("Failed to instantiate a smart contract: %v", err)
    }

    // Logging into Ethereum as a user
    //key []byte;
    key, e := ioutil.ReadFile( gConfig.KeyPath )
    if e != nil {
        fmt.Printf("Key File error: %v\n", e)
        os.Exit(1)
    }
    fmt.Printf("Found Ethereum Key File \n")

    auth, err := bind.NewTransactor(strings.NewReader(string(key)), gConfig.Pswd)
    if err != nil {
        log.Fatalf("Failed to create authorized transactor: %v", err)
    }

    session = &LuxUni_PKISession{
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
    session.TransactOpts = *auth;
    session.TransactOpts.GasLimit = big.NewInt(2000000)
    */

    // http://stackoverflow.com/questions/15834278/serving-static-content-with-a-root-url-with-the-gorilla-toolkit
    // subrouter - http://stackoverflow.com/questions/18720526/how-does-pathprefix-work-in-gorilla-mux-library-for-go
    r := mux.NewRouter();
    r.HandleFunc("/pki-test", PkiForm);
    r.HandleFunc("/enroll_user", EnrollUser);
    r.HandleFunc("/blacklist_user", BlacklistUser);
    //r.HandleFunc("/enroll_ca", EnrollCA);
    r.HandleFunc("/create_contract", CreateContract);
    r.HandleFunc("/populate_contract", PopulateContract);
    r.HandleFunc("/validate_cert", ValidateCert);

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
    contrAddrStr := r.MultipartForm.Value["ContrAddr"][0]
    if(common.IsHexAddress(contrAddrStr)==false){
         http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
             http.StatusInternalServerError)
        return
    }

    parentAddrStr := r.MultipartForm.Value["ParentAddr"][0]
    if(common.IsHexAddress(parentAddrStr)==false){
        http.Error(w, GeneralError{"Parent contract address is incorrect"}.Error(),
            http.StatusInternalServerError)
        return
    }

    certHash, _, fileCont, cerr := UploadFile(w,r,"UplFiles",
                common.HexToAddress(parentAddrStr), true)
    if(cerr.errCode!=0){
        fmt.Printf(fmt.Sprintf("Populate Uploadfile: %v\n", cerr.Error()))
        http.Error(w, cerr.Error(), http.StatusInternalServerError)
        return
    }

    client, err := ethclient.Dial(gConfig.IPCpath)
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }

    // Instantiate the contract, the address is taken from eth at the moment of contract initiation
    // kyc, err := NewLuxUni_KYC(common.HexToAddress(gContractHash), backends.NewRPCBackend(conn))
    pkiContract, err := NewLuxUni_PKI(common.HexToAddress(contrAddrStr), client)
    if err != nil {
        log.Fatalf("Failed to instantiate a smart contract: %v", err)
    }

    // Logging into Ethereum as a user
    //key []byte;
    key, e := ioutil.ReadFile( gConfig.KeyPath )
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

    _, err = sess.PopulateCertificate(fileCont)
    if err != nil {
        fmt.Printf(fmt.Sprintf("Failed to populate blockchain: %v.\n", err))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if (gConfig.WebMode!=0){
        var formParam CSimpleFormParam
        formParam.Params = append(formParam.Params,
            CParam{"Hash",strconv.Itoa(int(certHash))})
        formParam.Params = append(formParam.Params,
            CParam{"ContrAddr",contrAddrStr})
        formParam.Params = append(formParam.Params,
            CParam{"ParentAddr",parentAddrStr})
        formParam.Result = "Contract was successfully populated"
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
        }else{
            parentAddrStr = gConfig.ContractHash
        }
    }

    client, err := ethclient.Dial(gConfig.IPCpath)
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }

    key, e := ioutil.ReadFile( gConfig.KeyPath )
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
    contrAddr, _, /*contr*/_, err := DeployLuxUni_PKI(&trOpts, client, common.Address{})
    /*
    https://stackoverflow.com/questions/40096750/set-status-code-on-http-responsewriter
    */
    if err!=nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("500 - Something bad happened!"))
    }

    if (gConfig.WebMode!=0) {
        var formParam CSimpleFormParam
        formParam.IsUploadFileForm = true
        formParam.Params = append(formParam.Params,
            CParam{"ContrAddr",common.ToHex(contrAddr.Bytes())} )
        formParam.Params = append(formParam.Params,
            CParam{"ParentAddr",parentAddrStr})
        formParam.Result = "Contract was successfully created.\n Populate the contract"
        formParam.CallForm = "populate_contract"

        tmpl := mTempl["SimpleForm"]
        terr := tmpl.Execute(w, formParam)
        if terr != nil {
            http.Error(w, terr.Error(), http.StatusInternalServerError)
        }
    }
}

func ValidateCert(w http.ResponseWriter, r *http.Request){
    var isCertFound bool = false;

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
        return;
    }

    certHash, _, _, cerr := UploadFile(w,r,"UplFiles", common.Address{},true)
    if(cerr.errCode!=0){
        http.Error(w, cerr.Error(),
            http.StatusInternalServerError)
    }

    //isCurl = r.MultipartForm.Value["Curl"]
    //certHash := r.MultipartForm.Value["CertHash"][0]
    contrAddr := r.MultipartForm.Value["ContrAddr"][0]
    if(common.IsHexAddress(contrAddr)==false){
        http.Error(w, GeneralError{"Contract address is incorrect"}.Error(),
            http.StatusInternalServerError)
    }

    client, err := ethclient.Dial(gConfig.IPCpath)
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }

    currContract, err := NewLuxUni_PKI(common.HexToAddress(contrAddr), client)
    if err != nil {
        log.Fatalf("Failed to instantiate a smart contract: %v", err)
    }
    callOpts:= &bind.CallOpts{
        Pending: true,
    }

    numCert, err := currContract.NumRegData(callOpts)
    if(err!=nil){
        http.Error(w, GeneralError{"NumRegData problems"}.Error(),
            http.StatusInternalServerError)
    }

    for i:=int64(0); i<numCert.Int64(); i++ {
        bi := big.NewInt(i)
        //var certData CRegData;
        certData, err := currContract.RegData(callOpts, bi)
        if(err!=nil){
            http.Error(w, GeneralError{"RegData retrieval problems"}.Error(),
                http.StatusInternalServerError)
        }
        if( string(certData.DataHash) == strconv.Itoa(int(certHash)) ){
            isCertFound = true;
            break;
        }
    }
    if isCertFound == true {
        err := CheckCa(common.HexToAddress(contrAddr))
        if err!=nil {
            http.Error(w, GeneralError{"CA Cert not approved"}.Error(),
                http.StatusInternalServerError)
        }
    }else{
        http.Error(w, GeneralError{"Cert not found"}.Error(),
            http.StatusInternalServerError)
    }
}

func CheckCa(initAddr common.Address) (error) {
    //addr := common.HexToAddress(gConfig.ContractHash);
    addr := initAddr;
    var maxIter int = 1000;
    for i:=0; i < maxIter; i++ {
        addr = GetParent(addr)
        nilAddr := common.Address{}
        if addr == nilAddr{
            break
        }
        if (i==(maxIter-1)) && (addr != nilAddr) {
            return GeneralError{"MaxIter limit is reached"}
        }
    }
    return nil;
}

func GetParent(currAddr common.Address) (common.Address) {

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
    caCert, err := currContract.CaCertificate(callOpts)
    parentAddr := caCert[gCaCertOffset:gCaCertOffset+len(common.Address{}.Bytes())]
    return common.HexToAddress(string(parentAddr));
}

func LoadConfig() (error) {
    file, err := os.Open( gConfigFile )
    if err != nil {
        fmt.Printf("File error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Found configuration file: %s\n", gConfigFile)

    //jsonparser.Get(data, "person", "name", "fullName")
    jsonParser := json.NewDecoder(file)
    if err = jsonParser.Decode( &gConfig ); err != nil {
        fmt.Printf("Parsing config file: %s\n", err.Error())
        os.Exit(1)
    }
    
    b, err := json.Marshal(gConfig)
    if err != nil {
        fmt.Printf("Cannot convert conf file into string: %s", err)
        os.Exit(1);
    }    
    fmt.Printf("Loaded configuration file: %s\n", string(b) )
    file.Close()
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


