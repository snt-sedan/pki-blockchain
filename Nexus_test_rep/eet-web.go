
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
    "encoding/hex"
)

var mTempl map[string]*template.Template
var session *LuxUni_EETSession

const gConfigFile = "./eet-conf.json";

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
    DebugAdmin         int `json:"debugAdmin"`
}

var gSession struct{
    IsAdmin     bool
}


// alternative aproach to passing the params to template - http://stackoverflow.com/questions/23802008/how-do-you-pass-multiple-objects-to-go-template
// general Go template description https://gohugo.io/templates/go-templates/
type CDonateFormParam struct {
    IsAdmin      bool
    Desc         string
    Receivers    []CReceiverPrez
}

type CReceiverPrez struct {
    Index        int
    Addr         string
    Name         string
    Priorities   string
    Balance      float64
    TotDonation  float64
    BalanceStr   string
    TotDonationStr  string
    CreationDate time.Time
    CreationStr  string
}

type CSenderFormParam struct {
    Desc         string
    Senders      []CReceiverPrez
}


func init() {
    LoadConfig();

    if mTempl == nil {
        mTempl = make( map[string]*template.Template )
    }

    ///home/alex/DocsFS/Dropbox/WORK/RD/LuxBCh/PKI
    t, e := template.ParseFiles("./eet-form.html")
    if e != nil {
        fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "eet-form", e.Error())
        mTempl = nil
    } else {
        mTempl["MainForm"] = template.Must(t, e)
    }

    t, e = template.ParseFiles("./eetadm-form.html")
    if e != nil {
        fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "eetadm-form", e.Error())
        mTempl = nil
    } else {
        mTempl["AdminForm"] = template.Must(t, e)
    }

    t, e = template.ParseFiles("./eet-err.html")
    if e != nil {
        fmt.Printf("IMPORTANT: Error in parsing of template %v, %v\n", "eet-err", e.Error())
        mTempl = nil
    } else {
        mTempl["ErrForm"] = template.Must(t, e)
    }
}

func ChangeData(r *http.Request) (error){
    isDonateRequest := r.MultipartForm.Value["DonateButton"]
    isAddPointsRequest := r.MultipartForm.Value["AddPointsButton"]
    if len(isDonateRequest)!=0 || len(isAddPointsRequest)!=0 {
        numRece, err := session.NumReceivers()
        if err != nil {
            log.Fatalf("Failed to retrieve a total number of receivers: %v", err)
        }
        fmt.Println("Number of receivers (including those deleted) %v", numRece.Int64())

        if numRece.Int64() > 0 {
            fmt.Printf("Debug: I am in donation block")
            for i:=int64(0); i<numRece.Int64(); i++ {
                amountStr := r.MultipartForm.Value["Amount_"+strconv.Itoa(int(i))]
                addrStr := r.MultipartForm.Value["Addr_"+strconv.Itoa(int(i))]

                if(len(amountStr)>0 && len(amountStr[0])>0) {
                    amount, err := strconv.Atoi(amountStr[0]);
                    if err != nil {
                        return GeneralError{
                            fmt.Sprintf("Amount conversion error: %v", err.Error())};
                    }
                    if len(addrStr)== 0 {
                        return GeneralError{
                            fmt.Sprintf("AddStr nil - concept error.")}
                    }
                    var addr common.Address;
                    byteAddr, err := hex.DecodeString(addrStr[0])
                    if err!= nil || len(byteAddr)>len(addr) {
                        return GeneralError{
                            fmt.Sprintf("Error in hex decoding: len(addr)=%v, str=%v",
                            len(addrStr[0]), addrStr[0]) }
                    }
                    addr.SetBytes(byteAddr)
                    fmt.Printf("addr=%x; Init string=%v \n", addr.Bytes(), addrStr[0]);

                    if len(isDonateRequest)!=0 {
                        var zeroAddr common.Address;
                        //zeroAddr.SetString("0");
                        _, nerr := session.MakeDonation(zeroAddr, addr,
                            big.NewInt(int64(amount)), big.NewInt(int64(0)), nil )
                        if nerr != nil {
                            return GeneralError{
                                fmt.Sprintf("Make Donation error: %v", nerr)}
                        }
                    }
                    if len(isAddPointsRequest)!=0 {
                        _, nerr := session.AddPoints(addr,true,
                            big.NewInt(int64(amount)))
                        if nerr != nil {
                            return GeneralError{
                                fmt.Sprintf("Add points for receiver error: %v", nerr)}
                        }
                    }
                }
            }

        }

    }

    isAddCoinsSndRequest := r.MultipartForm.Value["AddCoinSndButton"]
    if len(isAddCoinsSndRequest)!=0 {
        numSender, err := session.NumSenders()
        if err != nil {
            return GeneralError{
                fmt.Sprintf("Failed to retrieve a total number of senders: %v", err)}
        }
        fmt.Println("Number of senders (including those deleted) %v", numSender.Int64())

        if numSender.Int64() > 0 {
            fmt.Printf("Debug: I am in add coins block")
            for i:=int64(0); i<numSender.Int64(); i++ {
                amountStr := r.MultipartForm.Value["SndAmount_"+strconv.Itoa(int(i))]
                addrStr := r.MultipartForm.Value["SndAddr_"+strconv.Itoa(int(i))]

                if(len(amountStr)>0 && len(amountStr[0])>0) {
                    amount, err := strconv.Atoi(amountStr[0]);
                    if err != nil {
                        return GeneralError{
                            fmt.Sprintf("Amount conversion error: %v", err)}
                    }
                    if len(addrStr)== 0 {
                        return GeneralError{
                            fmt.Sprintf("SndAddStr nil - concept error")}
                    }
                    var addr common.Address;
                    byteAddr, err := hex.DecodeString(addrStr[0])
                    if err!= nil || len(byteAddr)>len(addr) {
                        return GeneralError{
                            fmt.Sprintf("Error in hex decoding: len(addr)=%v, str=%v",
                            len(addrStr[0]), addrStr[0])}
                    }
                    addr.SetBytes(byteAddr)
                    fmt.Printf("addr=%x; Init string=%v; Amount=%v \n",
                        addr.Bytes(), addrStr[0], amount);

                    ret, nerr := session.AddPoints(addr,false,
                        big.NewInt(int64(amount)))
                    fmt.Printf("Result of the AddPoints : %v", ret.Value().Int64())
                    if nerr != nil {
                        return GeneralError{
                            fmt.Sprintf("Add points for sender error: %v", nerr)}
                    }
                }
            }

        }

    }

    isNewPersonRequest := r.MultipartForm.Value["NewPersonButton"]
    if len(isNewPersonRequest)!=0 {
        var sName, sType string
        var isReceiver bool
        var addr common.Address

        fmt.Printf("I am in the add person block \n");
        nameStr := r.MultipartForm.Value["NewName"]
        if len(nameStr)>0 {
            sName = nameStr[0]
            fmt.Printf("name=%v; ", sName);
        }
        addrStr := r.MultipartForm.Value["NewAddress"]
        if len(addrStr)>0 {
            byteAddr, err := hex.DecodeString(addrStr[0])
            if err!= nil || len(byteAddr)>len(addr) {
                return GeneralError{
                    fmt.Sprintf("Error in hex decoding: len(addr)=%v, str=%v",
                    len(addrStr[0]), addrStr[0])}
            }
            addr.SetBytes(byteAddr)
            fmt.Printf("addr=%x; Init string=%v \n", addr.Bytes(), addrStr[0]);
        }

        typeStr := r.MultipartForm.Value["NewTypePerson"]
        if len(typeStr)>0 {
            sType = typeStr[0]
            fmt.Printf("Person's type: %v; ", sType)
        }
        if(sType == "Receiver") {
            isReceiver = true
        }

        if(isReceiver == true) {
            ret, err := session.NewReceiver(addr, sName)
            if err != nil {
                return GeneralError{
                    fmt.Sprintf("Failed to add new person: %v, isReceiver: %v",
                    err, isReceiver)}
            }
            fmt.Printf("NewReceiver resulted in", ret.Value().Int64())
        } else {
            ret, err := session.NewSender(addr, sName)
            if err != nil {
                log.Fatalf("Failed to add new person: %v, isReceiver: %v",
                    err, isReceiver)
            }
            fmt.Printf("NewSender resulted in", ret.Value().Int64())
        }
    }
    return nil;
}

/*#
  # https://blog.saush.com/2015/03/18/html-forms-and-go/
*/
func donateForm(w http.ResponseWriter, r *http.Request) {

    var formParam CDonateFormParam;
    formParam.IsAdmin = gSession.IsAdmin
    var formSenderParam CSenderFormParam;
    var changeData bool = true;
    var isSenderListRequest []string;

    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        fmt.Printf("No change data: Parsing multipart form: %v\n", err.Error())
        changeData = false;
    }

    if (changeData == true) {
        isSenderListRequest = r.MultipartForm.Value["SenderListButton"]
        if (isSenderListRequest != nil) {
            changeData = false
        }
    }

    if (changeData == true) {
        //get a ref to the parsed multipart form
        err := ChangeData(r);
        if err!=nil {
            tmpl := mTempl["ErrForm"]
            terr := tmpl.Execute(w, fmt.Sprintf("%v",err))
            if terr != nil {
                http.Error(w, terr.Error(), http.StatusInternalServerError)
            }
            return
        }
    } /// end of changeData condition

    fmt.Println("Debug: passing in GET mode")

    numRece, err := session.NumReceivers()
    if err != nil {
        log.Fatalf("Failed to retrieve a total number of receivers: %v", err)
    }
    fmt.Println("Number of receivers (including those deleted) %v", numRece.Int64())

    for i := int64(0); i < numRece.Int64(); i++ {
        bi := big.NewInt(i)
        receiverAddr, err := session.ReceiverList(bi)
        if err != nil {
            log.Fatalf("Failed to get Receivers List data: %v", err)
        }

        receiverDatum, err := session.Receivers( receiverAddr )
        if err != nil {
            log.Fatalf("Failed to get Receivers data: %v", err)
        }

        crDate := time.Unix(receiverDatum.CreationDate.Int64(), 0);

        receiverTotPoints, err := session.GetTotPoints( receiverAddr, /* isReceiver */ true )
        if err != nil {
            log.Fatalf("Failed to get Receivers TotPoints %v", err)
        }

        receiverTotDonats, err := session.GetTotDonations( receiverAddr,
                /* isReceiver */ true , true)
        if err != nil {
            log.Fatalf("Failed to get Receivers TotDonations: %v", err)
        }

        receiverBalance := receiverTotPoints.Int64() - receiverTotDonats.Int64()

/*
    Index        int
    Addr         string
    Name         string
    Priorities   string
    Balance      float64
    TotDonation  float64
    CreationDate time.Time
    CreationStr  string
*/
        if(receiverBalance > 0 || gSession.IsAdmin == true){
            recePrez := CReceiverPrez{ int(bi.Int64()),
                fmt.Sprintf("%x", receiverAddr.Bytes()),
                receiverDatum.Name, ""/* priorities */,
                float64(receiverBalance),
                float64(receiverTotDonats.Int64()),
                fmt.Sprintf("%v", int64(receiverBalance)),
                fmt.Sprintf("%v", receiverTotDonats.Int64()),
                crDate, crDate.String()}
            formParam.Receivers = append(formParam.Receivers, recePrez )
        }

    }

    if isSenderListRequest != nil {
        numSend, err := session.NumSenders()
        if err != nil {
            log.Fatalf("Failed to retrieve a total number of senders: %v", err)
        }
        fmt.Println("Number of senders (including those deleted) %v", numRece.Int64())

        for i := int64(0); i < numSend.Int64(); i++ {
            bi := big.NewInt(i)
            sendAddr, err := session.SenderList(bi)
            if err != nil {
                log.Fatalf("Failed to get Senders List data: %v", err)
            }else{
                fmt.Printf("Sender %v, address %x \n", strconv.Itoa(int(i)), sendAddr.Bytes())
            }

            sendDatum, err := session.Senders(sendAddr)
            if err != nil {
                log.Fatalf("Failed to get Sender data: %v", err)
            }

            crDate := time.Unix(sendDatum.CreationDate.Int64(), 0);
            senderTotPoints, err := session.GetTotPoints(sendAddr, /* isReceiver */ false)
            if err != nil {
                log.Fatalf("Failed to get Senders tot points: %v", err)
            }
            if(senderTotPoints.Int64()<0){
                log.Fatalf("Sender totPoints error: %v", senderTotPoints.Int64())
            }

            senderTotDonats, err := session.GetTotDonations(sendAddr,
                    /* isReceiver */ false, true)
            if err != nil {
                log.Fatalf("Failed to get Senders Tot donations: %v", err)
            }
            if(senderTotDonats.Int64()<0){
                log.Fatalf("Sender totDonats error: %v", senderTotDonats.Int64())
            }

            senderBalance := senderTotPoints.Int64() - senderTotDonats.Int64()

            if (gSession.IsAdmin == true) {
                sendPrez := CReceiverPrez{int(bi.Int64()),
                    fmt.Sprintf("%x", sendAddr.Bytes()),
                    sendDatum.Name, "" /* priorities */ ,
                    float64(senderBalance),
                    float64(senderTotDonats.Int64()),
                    fmt.Sprintf("%v", int64(senderBalance)),
                    fmt.Sprintf("%v", senderTotDonats.Int64()),
                    crDate, crDate.String()}
                formSenderParam.Senders = append(formSenderParam.Senders, sendPrez)
            }
        }
    }

    if isSenderListRequest != nil {
        tmpl := mTempl["AdminForm"]
        terr := tmpl.Execute(w, formSenderParam)
        if terr != nil {
            http.Error(w, terr.Error(), http.StatusInternalServerError)
        }
    }

    tmpl := mTempl["MainForm"]
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
    eetContract, err := NewLuxUni_EET(common.HexToAddress(gConfig.ContractHash), client)
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

    session = &LuxUni_EETSession{
        Contract: eetContract,
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

    gSession.IsAdmin, err = session.IsAdmin()
    if err != nil {
        log.Fatalf("Failed in the first function call: %v", err)
    }
    gSession.IsAdmin = true;
    // DEBUG MODE BELOW
    /*if gConfig.DebugAdmin != 1 {
        gSession.IsAdmin = false;
    }*/
    fmt.Printf("Is it an admin session(true/false): %v \n", gSession.IsAdmin)

    // http://stackoverflow.com/questions/15834278/serving-static-content-with-a-root-url-with-the-gorilla-toolkit
    // subrouter - http://stackoverflow.com/questions/18720526/how-does-pathprefix-work-in-gorilla-mux-library-for-go
    r := mux.NewRouter();
    r.HandleFunc("/eet-test", donateForm);
    fs := http.FileServer(http.Dir("/home/alex/wrk/Dropbox/WORK/RD/Eethiq/Nexus_test/public"));
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




