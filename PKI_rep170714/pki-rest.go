
package main

import (
    "time"
    "fmt"
    "io/ioutil"
    "os"
    "log"
    "net/http"
    "strings"
    "strconv"
    "math/big"
    "crypto/rsa"
    "github.com/gorilla/mux"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/ethclient"
    "encoding/json"
    //"math/rand"
    "bytes"
)


func init() {
    err := LoadConfig();
    if err != nil {
        fmt.Printf("CONFIG ERROR: %v\n", err)
        os.Exit(1)
    }
    gConfig.WebMode = 1

}



func rstBlacklistUser(w http.ResponseWriter, r *http.Request){
    var revokeResult string
    var parentAddr common.Address = common.Address{};
    var userAddr common.Address = common.Address{};

    fmt.Println("REST: inside blacklist")
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

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}


/*
post paremeters:
Hash or Files
ParentAddr
ContrAddr
CurrentUserAddr -- corresponds to userAddr of a user holding the parent contract
NewUserAddr -- corresponds to userAddr associated with new contract
*/
func rstEnrollUser(w http.ResponseWriter, r *http.Request){
    //var isCurl []string;
    var parentAddr common.Address = common.Address{} // this is addr of the contract which is going to hold the hash
    var contrAddr common.Address = common.Address{} // this is address of the new SubCA contract or zero if end user
    var curUserAddr common.Address = common.Address{} // !! this is the user_id of the owner of parent contr
    var newUserAddr common.Address = common.Address{} // !! this is the new owner of contrAddr contr.
    var isNoUpload bool = false
    var desc string

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
            UploadFile(w, r, "UplFiles",true)
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
            http.Error(w, fmt.Sprintf("Parent address as a parameter is incorrect: %v",
                strParentAddrArr[0]),
                http.StatusInternalServerError)
            return
        }
        parentAddr = common.HexToAddress(strParentAddrArr[0])
    }

    if isNoUpload==false {
        var caContrAddr, insertAddr common.Address
        caContrAddr, insertAddr, desc, err = ParseCert(dataCert)
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

    strUserAddrArr := r.MultipartForm.Value["CurrentUserAddr"]
    if (len(strUserAddrArr) > 0) {
        if(common.IsHexAddress(strUserAddrArr[0])==false){
            http.Error(w, GeneralError{"CurrentUser address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        curUserAddr = common.HexToAddress(strUserAddrArr[0])
    }

    strUserAddrArr = r.MultipartForm.Value["NewUserAddr"]
    if (len(strUserAddrArr) > 0) {
        if(common.IsHexAddress(strUserAddrArr[0])==false){
            http.Error(w, GeneralError{"NewUser address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        newUserAddr = common.HexToAddress(strUserAddrArr[0])
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

    fmt.Printf("DEBUG before newRegDatum: hash=%v, fname=%v, desc=%v, newUserAddr=%v\n",
        strconv.Itoa(int(hashSum)), fileName, desc, newUserAddr.String())

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
    if (curUserAddr == common.Address{}) {
        fmt.Printf("Attention! Enroll: user address is zero, default config account is used\n")
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
    sess.TransactOpts.GasLimit = big.NewInt(2000000)

    res, err := sess.NewRegDatum( []byte(strconv.Itoa(int(hashSum))), contrAddr,
            fileName, desc, "", newUserAddr)
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

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}


func main() {

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
    //r.HandleFunc("/pki-test", PkiForm);
    r.HandleFunc("/enroll_user", rstEnrollUser);
    r.HandleFunc("/blacklist_user", rstBlacklistUser);
    //r.HandleFunc("/enroll_ca", EnrollCA);
    r.HandleFunc("/create_contract", rstCreateContract);
    r.HandleFunc("/populate_contract", rstPopulateContract);
    //r.HandleFunc("/validate_form", ValidateForm);
    r.HandleFunc("/validate_cert", rstValidateCert);
    r.HandleFunc("/download_cacert", rstDownloadCaCert);
    //r.HandleFunc("/generate_user_cert", GenerateUserCert);

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

    log.Println("RESTful service is listening...")
    //http.ListenAndServeTLS(":8071", "server.pem", "server.key", r)
    http.ListenAndServe(":"+strconv.Itoa(gConfig.RestHttpPort), nil)
}


func rstPopulateContract(w http.ResponseWriter, r *http.Request) {

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

    /*
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
    */

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


    hashCert, _, dataCert, cerr := UploadFile(w,r,"UplFiles", true)
    if(cerr.errCode!=0){
        fmt.Printf(fmt.Sprintf("Populate Uploadfile: %v\n", cerr.Error()))
        http.Error(w, cerr.Error(), http.StatusInternalServerError)
        return
    }
    /*
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
    */

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

    w.WriteHeader(http.StatusOK)
    w.Write([]byte( strconv.Itoa(int(hashCert)) ))
}


func rstCreateContract(w http.ResponseWriter, r *http.Request){
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

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(contrAddr.String()))
}


/*
  https://stackoverflow.com/questions/35496233/go-how-to-i-make-download-service
  https://play.golang.org/p/UMKgI_NLwO
*/
func rstDownloadCaCert(w http.ResponseWriter, r *http.Request){

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


func rstValidateCert(w http.ResponseWriter, r *http.Request){

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
        return;
    }

    certHash, _, dataCert, cerr := UploadFile(w,r,"UplFiles", true)
    if(cerr.errCode!=0){
        http.Error(w, cerr.Error(),
            http.StatusInternalServerError)
        return
    }

    contrAddr, parentAddr, _, err:= ParseCert(dataCert)
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

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}




