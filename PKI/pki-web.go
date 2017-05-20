
package main

import (
    "html/template"
    "time"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "log"
    "net/http"
    "strings"
    "strconv"
    "math/big"
    "math/rand"
    "hash/crc32"
    cryptoRand "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "encoding/gob"

    "github.com/gorilla/mux"

    //"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
    //"github.com/ethereum/go-ethereum/rpc"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"

    "github.com/ethereum/go-ethereum/ethclient"

    "encoding/json"
)

var mTempl map[string]*template.Template
var session *LuxUni_PKISession

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
    CryptoModulHash    string `json:"cryptoModulHash"`
    ContractHash       string `json:"contractHash"`
    IPCpath            string `json:"IPCpath"`
    Pswd               string `json:"pswd"`
    KeyPath            string `json:"keyPath"`
    PrivateKeyPath     string `json:"privateKeyPath"`
    HttpPort           int `json:"httpPort"`
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
type CKYCFormParam struct {
    RandomNum    int
    RandomStr    string
    Docs         []CDocPrez
}

type CDocPrez struct {
    Id           int
    Name         string
    Desc         string
    Link         string
    Decryption   string
    Hash         string
    CreationDate  time.Time
    CreationStr   string
}

type CRevokedFormParam struct {
    RevokedIds    []int
    RevokedStr   string
    Docs         []CDocPrez
}


func init() {
    LoadConfig();

    if mTempl == nil {
        mTempl = make( map[string]*template.Template )
    }

    ///home/alex/DocsFS/Dropbox/WORK/RD/LuxBCh/PKI
    t, e := template.ParseFiles("./pki-form.html")
    if e != nil {
        fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "form", e.Error())
        mTempl = nil
    } else {
        mTempl["UploadFile"] = template.Must(t, e)
    }

    t, e = template.ParseFiles("./hashres-form.html")
    if e != nil {
        fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "form", e.Error())
        mTempl = nil
    } else {
        mTempl["HashResult"] = template.Must(t, e)
    }
}


/*#
  # https://blog.saush.com/2015/03/18/html-forms-and-go/
*/
func kycForm(w http.ResponseWriter, r *http.Request) {

    var formParam CKYCFormParam;
    var revokedParam CRevokedFormParam;
    var hashResult string ="";
    var changeData bool = true;
    var isRevokeListRequest []string;
    var isCurl []string;

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
        changeData = false;
    }

    if (changeData == true) {
        isCurl = r.MultipartForm.Value["Curl"]
        isRevokeListRequest = r.MultipartForm.Value["RevokeListButton"]
        if (isRevokeListRequest != nil) {
            changeData = false
        }
    }

    if (changeData == true) {

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

                _, nerr := session.DeleteRegDatum(big.NewInt(int64(delid)))
                if nerr != nil {
                    log.Fatalf("Deletion error: %v", nerr)
                }
            }
        }

        //get a ref to the parsed multipart form
        m := r.MultipartForm
        files := m.File["UplFiles"];

        for i, _ := range files {
            //for each fileheader, get a handle to the actual file
            file, err := files[i].Open()
            defer file.Close()
            if err != nil {
                fmt.Printf("Parsing form: %v\n", err.Error())
                return
            }

            //create destination file making sure the path is writeable.
            dstFName := "/tmp/tst" + "_" + files[i].Filename
            dst, err := os.Create(dstFName)
            defer dst.Close()
            if err != nil {
                fmt.Printf("Unable to open the file for writing: %v\n", err.Error())
                return
            } else {
                fmt.Printf("Created file: %v\n", dstFName)
            }

            //copy the uploaded file to the destination file
            _, err = io.Copy(dst, file); 
            if err != nil {
                fmt.Printf("Unable to copy: %v\n", err.Error())
                return
            }

            hasher := crc32.NewIEEE()
            //_, err = io.Copy(hasher, dst);
            dst4hash, err := ioutil.ReadFile(dstFName);
            if err != nil {
                fmt.Println("Error in open file for hash: %v\n", err.Error())
            }
            hasher.Write(dst4hash);
            hashSum := hasher.Sum32()
            fmt.Printf("Hash: %x\n", hashSum)

            var desc string = ""
            var isEncrypt int64;
            descArr := r.MultipartForm.Value["Desc"]
            if (len(descArr) > 0) {
                desc = descArr[0];
            }            
            encryptArr := r.MultipartForm.Value["Encryption"]
            if (len(encryptArr) > 0) {
                if encryptArr[0] != "" {
                    isEncrypt = 1;
                }
            }            
            fmt.Printf("DEBUG before newRegDatum: hash=%v, fname=%v, desc=%v, encrypt=%v\n", 
                    strconv.Itoa(int(hashSum)), files[i].Filename, desc, isEncrypt)

            _, err = session.NewRegDatum( []byte(strconv.Itoa(int(hashSum))), 
                        files[i].Filename, desc, "", 
                        big.NewInt(isEncrypt), common.HexToAddress(gConfig.CryptoModulHash))
            if err != nil {
               log.Fatalf("Failed to add a record to blockchain: %v. ", err)
               fmt.Fprintln(w, err)
            }
            hashResult = strconv.Itoa(i+1) + " file(s) processed:"
        }
        // UplFile is id in the input "file" component of the form
        // http://stackoverflow.com/questions/33771167/handle-file-uploading-with-go
        // file, handler, err := r.FormFile("UplFile")       
        //out, err := os.Create("/tmp/tst_"+handler.Filename);
        
        hashResult = "Hash info is successfully loaded to blockchain." + hashResult;
        /*
        tmpl := mTempl["HashResult"]
        terr := tmpl.Execute(w, hashResult)
        if terr != nil {
            http.Error(w, terr.Error(), http.StatusInternalServerError)
        }
        */

    } /// end of changeData condition

    randomBase := 1000000
    random := rand.New(rand.NewSource(time.Now().UnixNano()));
    formParam.RandomNum = random.Intn(9*randomBase) + randomBase
    formParam.RandomStr = strconv.Itoa(formParam.RandomNum)

    fmt.Println("Debug: passing in GET mode")

    numRD, err := session.NumRegData()
    if err != nil {
        log.Fatalf("Failed to retrieve a total number of data records: %v", err)
    }
    fmt.Println("Number of data records (including those deleted)", numRD)

    for i := int64(0); i < numRD.Int64(); i++ {
        bi := big.NewInt(i)
        regDatum, err := session.RegData(bi)
        if err != nil {
            log.Fatalf("Failed to get deleted data: %v", err)
        }

        delRegDatum, err := session.DeletedRegData( bi)
        if err != nil {
            log.Fatalf("Deletion retrieval error: %v", err)
        }

        crDate := time.Unix(regDatum.CreationDate.Int64(), 0);
        var desc string
        var decrypt string

        if ( regDatum.Encrypted.Int64() != 0 ) { 
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
        } else {
            desc = regDatum.Description;
        }

        docPrez := CDocPrez{ int(i), regDatum.FileName, 
                       desc, ""/* link */, decrypt,
                       string(regDatum.DataHash), crDate, crDate.String()}

        // scanning of recently revoked ids for formation of a string with results of revocation
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

    if ( (len(revokedParam.RevokedIds) > 0) || ( isRevokeListRequest != nil ) ) {
        tmpl := mTempl["HashResult"]
        terr := tmpl.Execute(w, revokedParam)
        if terr != nil {
            http.Error(w, terr.Error(), http.StatusInternalServerError)
        }
    }

    tmpl := mTempl["UploadFile"]
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


    // Create an IPC based RPC connection to a remote node
    // conn, err := rpc.NewIPCClient(gIPCpath)
    client, err := ethclient.Dial(gConfig.IPCpath)
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client: %v", err)
    }

    // Instantiate the contract, the address is taken from eth at the moment of contract initiation
    // kyc, err := NewLuxUni_KYC(common.HexToAddress(gContractHash), backends.NewRPCBackend(conn))
    pkiContract, err := NewLuxUni_PKI(common.HexToAddress(gConfig.ContractHash), client)
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
    
    // http://stackoverflow.com/questions/15834278/serving-static-content-with-a-root-url-with-the-gorilla-toolkit
    // subrouter - http://stackoverflow.com/questions/18720526/how-does-pathprefix-work-in-gorilla-mux-library-for-go
    r := mux.NewRouter();
    r.HandleFunc("/pki-test", kycForm);
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




