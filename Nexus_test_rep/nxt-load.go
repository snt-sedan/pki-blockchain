
package main

import (
    "time"
    "fmt"
    "io/ioutil"
    "os"
    "log"
    "strings"
    //"strconv"
    "math/big"
    //"math/rand"
    cryptoRand "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "encoding/gob"

    //"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
    //"github.com/ethereum/go-ethereum/rpc"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"

    "github.com/ethereum/go-ethereum/ethclient"

    "encoding/json"
    "encoding/hex"
)

var gSession1 *LuxUni_EETSession
var gSession2 *LuxUni_EETSession

const gConfigFile = "./load-conf.json";

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
    IPCpath            string `json:"IPCpath"`
    IterNum            int `json:"iterNum"`
    AddrSender       string `json:"addrSender"`
    AddrReceiver       string `json:"addrReceiver"`
    CertPath       string `json:"certPath"`
    DecesecDelay       int `json:"decesecDelay"`
    ContractHash       string `json:"contractHash"`

    Pswd1               string `json:"pswd1"`
    KeyPath1            string `json:"keyPath1"`

    Pswd2               string `json:"pswd2"`
    KeyPath2            string `json:"keyPath2"`
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

}

func AddPoints(session *LuxUni_EETSession, isReceiver bool, amount int, sAddr string ) (error) {
    if(gConfig.DecesecDelay>0){
        time.Sleep(time.Second / 10 * time.Duration(gConfig.DecesecDelay))
    }

    var addr common.Address;
    byteAddr, err := hex.DecodeString(sAddr)
    if err!= nil || len(byteAddr)>len(addr) {
        return GeneralError{
                        fmt.Sprintf("Error in hex decoding: len(addr)=%v, str=%v",
                            len(sAddr), sAddr)}
    }
    addr.SetBytes(byteAddr)
    fmt.Printf("Addpoints addr=%x; Init string=%v; Amount=%v \n",
                    addr.Bytes(), sAddr, amount);

    ret, nerr := session.AddPoints(addr, isReceiver,
                    big.NewInt(int64(amount)))
    fmt.Printf("Result of the AddPoints : %v", ret.Value().Int64())
    if nerr != nil {
        return GeneralError{
            fmt.Sprintf("Add points for sender error: %v", nerr)}
    }
    return nil
}

func Donate(session *LuxUni_EETSession, amount int, sAddr string ) (error) {
    if(gConfig.DecesecDelay>0){
        time.Sleep(time.Second / 10 * time.Duration(gConfig.DecesecDelay))
    }

    var addr common.Address;
    byteAddr, err := hex.DecodeString(sAddr)
    if err!= nil || len(byteAddr)>len(addr) {
        return GeneralError{
            fmt.Sprintf("Error in hex decoding: len(addr)=%v, str=%v",
                len(sAddr), sAddr) }
    }
    addr.SetBytes(byteAddr)
    //fmt.Printf("Donate addr=%x; Init string=%v \n", addr.Bytes(), sAddr);

    var certf []byte;
    if len(gConfig.CertPath)>0 {
        certf, err = ioutil.ReadFile( gConfig.CertPath )
        if err != nil {
            fmt.Printf("Cert file load error: %v\n", err)
        }
    }

    var zeroAddr common.Address;
    //zeroAddr.SetString("0");
    ret, nerr := session.MakeDonation(zeroAddr, addr,
            big.NewInt(int64(amount)), big.NewInt(int64(0)), certf )
    if nerr != nil {
        /*
        fmt.Printf("Donat error %v, Second Attempt", nerr)
        if(gConfig.DecesecDelay>0){
            time.Sleep(time.Second / 5 * time.Duration(gConfig.DecesecDelay))
        }
        _, nerr := session.MakeDonation(zeroAddr, addr,
            big.NewInt(int64(amount)), big.NewInt(int64(0)), certf )*/
        if nerr != nil {
            return GeneralError{
               fmt.Sprintf("Make Donation error: %v", nerr)}
        }
    }
    if ret.Value().Int64()<0 {
        return GeneralError{
            fmt.Sprintf("Make Donation return error: %v", ret.Value().Int64() )}
    }
    return nil
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
    client, err := ethclient.Dial(gConfig.IPCpath)
    if err != nil {
        log.Fatalf("Failed to connect to the Ethereum client1: %v", err)
    }

    // Instantiate the contract, the address is taken from eth at the moment of contract initiation
    // kyc, err := NewLuxUni_KYC(common.HexToAddress(gContractHash), backends.NewRPCBackend(conn))
    eetContract, err := NewLuxUni_EET(common.HexToAddress(gConfig.ContractHash), client)
    if err != nil {
        log.Fatalf("Failed to instantiate a smart contract: %v", err)
    }

    // Logging into Ethereum as a user
    //key []byte;
    key1, e := ioutil.ReadFile( gConfig.KeyPath1 )
    if e != nil {
        fmt.Printf("Key File error: %v\n", e)
        os.Exit(1)
    }
    fmt.Printf("Found Ethereum Key File \n")
    /*key2, e := ioutil.ReadFile( gConfig.KeyPath2 )
    if e != nil {
        fmt.Printf("Key File error: %v\n", e)
        os.Exit(1)
    }
    fmt.Printf("Found Ethereum Key File \n")*/

    /* this is the admin session - only the owner of the contact should add coins/points*/
    auth1, err := bind.NewTransactor(strings.NewReader(string(key1)), gConfig.Pswd1)
    if err != nil {
        log.Fatalf("Failed to create authorized transactor: %v", err)
    }
    /* this is the sender session - Donate is done based on it, user SHOULD BE REGISTERED in senders*/
    /*
    auth2, err := bind.NewTransactor(strings.NewReader(string(key2)), gConfig.Pswd2)
    if err != nil {
        log.Fatalf("Failed to create authorized transactor: %v", err)
    }*/

    gSession1 = &LuxUni_EETSession{
        Contract: eetContract,
        CallOpts: bind.CallOpts{
            Pending: true,
        },
        TransactOpts: bind.TransactOpts{
            From:     auth1.From,
            Signer:   auth1.Signer,
            GasLimit: big.NewInt(2000000),
        },
    }
    gSession1.TransactOpts = *auth1;
    gSession1.TransactOpts.GasLimit = big.NewInt(2000000)

    /*gSession2 = &LuxUni_EETSession{
        Contract: eetContract,
        CallOpts: bind.CallOpts{
            Pending: true,
        },
        TransactOpts: bind.TransactOpts{
            From:     auth2.From,
            Signer:   auth2.Signer,
            GasLimit: big.NewInt(2000000),
        },
    }
    gSession2.TransactOpts = *auth2;
    gSession2.TransactOpts.GasLimit = big.NewInt(2000000)*/

    for i:=0; i<gConfig.IterNum; i++ {
        payment := 10
        transNum := 500
        err := AddPoints(gSession1, false, (transNum * payment), gConfig.AddrSender)
        if err != nil {
            log.Fatalf("Error in add coins to sender: %v", err)
        }
        err = AddPoints(gSession1, true, (transNum * payment), gConfig.AddrReceiver)
        if err != nil {
            log.Fatalf("Error in add limits to receiver: %v", err)
        }
        for j:=0; j<transNum; j++ {
            err = Donate(gSession1,payment,gConfig.AddrReceiver)
            if err != nil {
                log.Fatalf("Error in donation: %v", err)
            }
            fmt.Print(".")
        }
        fmt.Print("!")
    }
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




