
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
    cryptoRand "crypto/rand"
    "crypto/rsa"
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
    "mime/multipart"
    "io"
)

var mTempl map[string]*template.Template
//var session *LuxUni_PKISession
//var gClient *ethclient.Client
//var gCaCertName = "_ca_cert.crt"

/*
CORE PARAMETERS ARE STORED IN gConfig 

const gCryptoModulHash = "0x3f2ed40488d0a9586013faa415718f3f64644fa1";
const gContractHash = "0xf1918d06d7e66e60153d7109ff380b41866ba2e0";
const gIPCpath = "/home/alex/_Work/Eth_AeroNet_t/geth.ipc";
const gPswd = "ira";
const key = `{"address":"a6f23407d139508fa38706140c56bf6487f87395","crypto":{"cipher":"aes-128-ctr","ciphertext":"5468728413080efcf6191dc3b5eaaf6ff34fd401277b9c08ac5aa93b6c3b3e44","cipherparams":{"iv":"efae50c39291e191193635ec715e96f4"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"a1e83898937dcd79e480d642aab80032a983a86e55c1ff8c64e91a1fc4fd9bdd"},"mac":"35456f2b746b847b4f1f5af2291f8f1518d6eed4e29ce950c1123da7ac3ee01d"},"id":"71412808-a83f-4677-ac6d-3d0f744d3ef0","version":3}`
const gPrivateKeyPath = "a_pr.key";
*/

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
/*
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


    //https://stackoverflow.com/questions/36111777/golang-how-to-read-a-text-file
    //https://stackoverflow.com/questions/30182538/why-can-not-i-copy-a-slice-with-copy-in-golang
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
*/

func BlacklistUser(w http.ResponseWriter, r *http.Request){
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

    var b bytes.Buffer
    newWri := multipart.NewWriter(&b)

    strParentAddrArr := r.MultipartForm.Value["ParentAddr"]
    if (len(strParentAddrArr) > 0) {
        if(common.IsHexAddress(strParentAddrArr[0])==false){
            http.Error(w, GeneralError{"Parent address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        parentAddr = common.HexToAddress(strParentAddrArr[0])
    }
    if (parentAddr == common.Address{}){
        http.Error(w, "Delete: Parent address is not found",
            http.StatusInternalServerError)
        return
    }
    fw, err := newWri.CreateFormField("ParentAddr");
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
    if (len(strUserAddrArr) > 0) {
        if(common.IsHexAddress(strUserAddrArr[0])==false){
            http.Error(w, GeneralError{"User address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        userAddr = common.HexToAddress(strUserAddrArr[0])

        fw, err := newWri.CreateFormField("UserAddr");
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
        fw, err := newWri.CreateFormField("Deletion");
        if err != nil {
            http.Error(w, "Delete: Deletion field is not created in request to REST",
                http.StatusInternalServerError)
            return
        }
        for _, del := range dels {
            fmt.Printf("del=%v\n", del);
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
            gConfig.RestUrlServer+":"+strconv.Itoa(gConfig.RestHttpPort)+"/blacklist_user", &b)
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
        http.Error(w, fmt.Sprintf("Delete: REST's bad status: %v, %v", res.Status, strResult),
            http.StatusInternalServerError)
        return
    }

    if (gConfig.WebMode!=0){
        var formParam CSimpleFormParam
        //formParam.Param = ""
        formParam.Result = "CA "+parentAddr.String()+": certificate was successfully revoked"
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
Hash or Files
ParentAddr
ContrAddr
CurrentUserAddr -- corresponds to userAddr of a user holding the parent contract
NewUserAddr -- corresponds to userAddr associated with new contract
*/
func EnrollUser(w http.ResponseWriter, r *http.Request){
    //var isCurl []string;
    var parentAddr common.Address = common.Address{} // this is addr of the contract which is going to hold the hash
    var contrAddr common.Address = common.Address{} // this is address of the new SubCA contract or zero if end user
    var curUserAddr common.Address = common.Address{} // !! this is the user_id of the owner of parent contr
    var newUserAddr common.Address = common.Address{} // !! this is the new owner of contrAddr contr.
    var isNoUpload bool = false

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("EnrollUser: No change data -- ", err.Error())
        http.Error(w, GeneralError{fmt.Sprintf(
            "EnrollUser: No change data -- ",err.Error())}.Error(),
            http.StatusInternalServerError)
        return
    }

    var b bytes.Buffer
    newWri := multipart.NewWriter(&b)

    //isCurl = r.MultipartForm.Value["Curl"]
    hashSum, fileName, dataCert, cerr :=
            UploadFile(w, r, "UplFiles", true)
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

            fw, err := newWri.CreateFormField("Hash");
            if err != nil {
                http.Error(w, "EnrollUser: Hash is not created in request to REST",
                    http.StatusInternalServerError)
                return
            }
            if _, err = fw.Write([]byte( strconv.Itoa(int(hashSum)) )); err != nil {
                http.Error(w, "EnrollUser: Hash is not written in request to REST",
                    http.StatusInternalServerError)
                return
            }

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

        fw, err := newWri.CreateFormField("ParentAddr");
        if err != nil {
            http.Error(w, "EnrollUser: ParentAddr is not created in request to REST",
                http.StatusInternalServerError)
            return
        }
        if _, err = fw.Write([]byte( parentAddr.String() )); err != nil {
            http.Error(w, "EnrollUser: ParentAddr is not written in request to REST",
                http.StatusInternalServerError)
            return
        }
    }

    if isNoUpload==false {
        fw, err := newWri.CreateFormFile("UplFiles", fileName)
        if err != nil {
            http.Error(w, fmt.Sprintf("EnrollUser: cannot create form file: %v", err),
                http.StatusInternalServerError)
            return
        }
        if _, err = io.Copy(fw, bytes.NewReader(dataCert) ); err != nil {
            http.Error(w, fmt.Sprintf("EnrollUser: cannot read/copy cert data: %v", err),
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

        fw, err := newWri.CreateFormField("CurrentUserAddr");
        if err != nil {
            http.Error(w, "EnrollUser: CurrentUserAddr is not created in request to REST",
                http.StatusInternalServerError)
            return
        }
        if _, err = fw.Write([]byte( curUserAddr.String() )); err != nil {
            http.Error(w, "EnrollUser: CurrentUserAddr is not written in request to REST",
                http.StatusInternalServerError)
            return
        }
    }

    strUserAddrArr = r.MultipartForm.Value["NewUserAddr"]
    if (len(strUserAddrArr) > 0) {
        if(common.IsHexAddress(strUserAddrArr[0])==false){
            http.Error(w, GeneralError{"NewUser address is incorrect"}.Error(),
                http.StatusInternalServerError)
            return
        }
        newUserAddr = common.HexToAddress(strUserAddrArr[0])

        fw, err := newWri.CreateFormField("NewUserAddr");
        if err != nil {
            http.Error(w, "EnrollUser: NewUserAddr is not created in request to REST",
                http.StatusInternalServerError)
            return
        }
        if _, err = fw.Write([]byte( newUserAddr.String() )); err != nil {
            http.Error(w, "EnrollUser: NewUserAddr is not written in request to REST",
                http.StatusInternalServerError)
            return
        }
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

        fw, err := newWri.CreateFormField("ContrAddr");
        if err != nil {
            http.Error(w, "EnrollUser: ContrAddr is not created in request to REST",
                http.StatusInternalServerError)
            return
        }
        if _, err = fw.Write([]byte( contrAddr.String() )); err != nil {
            http.Error(w, "EnrollUser: ContrAddr is not written in request to REST",
                http.StatusInternalServerError)
            return
        }
    }

    if (parentAddr == common.Address{}){
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
    if res.StatusCode != http.StatusOK {
        restResult, _ := ioutil.ReadAll(res.Body)
        http.Error(w, fmt.Sprintf("EnrollUser: REST's bad status: %v, %v", res.Status, string(restResult)),
            http.StatusInternalServerError)
        return
    }


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
        //http.Error(w, GeneralError{"Parent address is nil: "}.Error(),
        //    http.StatusInternalServerError)
        //return
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
    http.ListenAndServe(":"+strconv.Itoa(gConfig.AppHttpPort), nil)
}



func PopulateContract(w http.ResponseWriter, r *http.Request) {

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
        return;
    }

    var b bytes.Buffer
    newWri := multipart.NewWriter(&b)

    strContrAddr, err := CopyRequestField(r, newWri,
        "ContrAddr", "ContrAddr")
    if err!=nil {
        http.Error(w, fmt.Sprintf("Populate contract: %v", err),
            http.StatusInternalServerError)
        return
    }
    if (common.IsHexAddress(strContrAddr)==false) {
        http.Error(w, fmt.Sprintf("Contract address is incorrect: %v", strContrAddr),
            http.StatusInternalServerError)
        return
    }
    contrAddr := common.HexToAddress(strContrAddr)

    strParentAddr, err := CopyRequestField(r, newWri,
        "ParentAddr", "ParentAddr")
    if err!=nil {
        http.Error(w, fmt.Sprintf("Populate contract: %v", err),
            http.StatusInternalServerError)
        return
    }
    if (common.IsHexAddress(strContrAddr)==false) {
        http.Error(w, fmt.Sprintf("Contract address is incorrect: %v", strContrAddr),
            http.StatusInternalServerError)
        return
    }
    parentAddr := common.HexToAddress(strParentAddr)

    strNewUserAddr, err := CopyRequestField(r, newWri,
        "NewUserAddr", "NewUserAddr")
    if err!=nil {
        http.Error(w, fmt.Sprintf("Populate contract: %v", err),
            http.StatusInternalServerError)
        return
    }
    if (common.IsHexAddress(strContrAddr)==false) {
        http.Error(w, fmt.Sprintf("Contract address is incorrect: %v", strContrAddr),
            http.StatusInternalServerError)
        return
    }
    newUserAddr := common.HexToAddress(strNewUserAddr)

    strCurUserAddr, err := CopyRequestField(r, newWri,
        "CurrentUserAddr", "CurrentUserAddr")
    if err!=nil {
        http.Error(w, fmt.Sprintf("Populate contract: %v", err),
            http.StatusInternalServerError)
        return
    }
    if (common.IsHexAddress(strContrAddr)==false) {
        http.Error(w, fmt.Sprintf("Populate: Current user address is incorrect: %v", strContrAddr),
            http.StatusInternalServerError)
        return
    }
    curUserAddr := common.HexToAddress(strCurUserAddr)


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

    fw, err := newWri.CreateFormFile("UplFiles", "CA_cert.crt")
    if err != nil {
        http.Error(w, fmt.Sprintf("EnrollUser: cannot create form file: %v", err),
            http.StatusInternalServerError)
        return
    }
    if _, err = io.Copy(fw, bytes.NewReader(dataCert) ); err != nil {
        http.Error(w, fmt.Sprintf("EnrollUser: cannot read/copy cert data: %v", err),
            http.StatusInternalServerError)
        return
    }

    /*
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
	*/

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
    detailStr := "Contract Addr = "+contrAddr.String()+"\n"
    detailStr += "Parent Addr = "+common.ToHex(parentAddr.Bytes())+"\n"
    detailStr += "New User = "+common.ToHex(newUserAddr.Bytes())+"\n"
    detailStr += "Current User = "+common.ToHex(curUserAddr.Bytes())+"\n"
    if (gConfig.WebMode!=0){
        var formParam CSimpleFormParam
        formParam.Params = append(formParam.Params,
            CParam{"Hash",strconv.Itoa(int(hashCert))})
        formParam.Params = append(formParam.Params,
            CParam{"ContrAddr",contrAddr.String()})
        formParam.Params = append(formParam.Params,
            CParam{"ParentAddr",common.ToHex(parentAddr.Bytes())})
        formParam.Params = append(formParam.Params,
            CParam{"CurrentUserAddr",common.ToHex(curUserAddr.Bytes())})
        formParam.Params = append(formParam.Params,
            CParam{"NewUserAddr",common.ToHex(newUserAddr.Bytes())})
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

    certHash, _, dataCert, cerr := UploadFile(w,r,"UplFiles",true)
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


func GenerateCert(contrAddr common.Address, parentAddr common.Address,
          isCA bool, strName string) ([]byte, error){

    // see Certificate structure at
    // http://golang.org/pkg/crypto/x509/#Certificate
    template := &x509.Certificate {
        IsCA : isCA,
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
   processes and copies field from incoming request to outgoing request to REST service
 */
func CopyRequestField(r *http.Request, newWri *multipart.Writer,
        fieldNameIn string, fieldNameOut string) (fieldVal string, err error) {
    strFieldArr := r.MultipartForm.Value[fieldNameIn]
    if (len(strFieldArr) == 0) {
        return "", GeneralError{fmt.Sprintf("The number of items for %v is %v (>1)",
            fieldNameIn, len(strFieldArr) )}
    }
    if (len(strFieldArr) > 1) {
        return "", GeneralError{fmt.Sprintf("The number of items for %v is %v (>1)",
            fieldNameIn, len(strFieldArr) )}
    }

    fieldVal = strFieldArr[0]

    fw, err := newWri.CreateFormField(fieldNameOut);
    if err != nil {
        return "", GeneralError{ fmt.Sprintf("%v is not created in request to REST",fieldNameOut)}
    }
    if _, err = fw.Write([]byte( fieldVal )); err != nil {
        return "", GeneralError{ fmt.Sprintf("%v is not written in request to REST",fieldNameOut)}
    }
    return fieldVal, nil
}
