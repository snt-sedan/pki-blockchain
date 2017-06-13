testnet 
wallet/user "0x78714b8663f11a388b9b74c3edbc42d253d2b0da"


solc --bin pki_scont.sol > pki_scont.bin
solc --abi pki_scont.sol > pki_abi.json
abigen --abi pki_abi.json --pkg main --type LuxUni_PKI --out bind_pki.go --bin pki_scont.bin

cd ~/go/src/github.com/ethereum/go-ethereum/
godep go build ~/wrk/Dropbox/WORK/RD/LuxBCh/PKI/pki-web.go ~/wrk/Dropbox/WORK/RD/LuxBCh/PKI/bind_pki.go

//GETH

primary	= eth.accounts[0];
web3.fromWei(eth.getBalance(primary), "ether");
personal.unlockAccount(primary);

// getting access to the contract from geth
//# inside parathesis is abi.json file for this contract
pki_abi = ''

pki_c = eth.contract(
[{"constant":false,"inputs":[{"name":"_cryptoID","type":"uint256"},{"name":"_bData","type":"bytes"}],"name":"encryptCallBack","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"deletedRegData","outputs":[{"name":"nodeSender","type":"address"},{"name":"deletionDate","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_regID","type":"uint256"}],"name":"deleteRegDatum","outputs":[{"name":"err","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"founder","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"LuxUni_KYC","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"numRegData","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"encryptRegData","outputs":[{"name":"nodeSender","type":"address"},{"name":"data","type":"bytes"},{"name":"encryptDate","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_dataHash","type":"bytes"},{"name":"_fileName","type":"string"},{"name":"_description","type":"string"},{"name":"_linkFile","type":"string"},{"name":"_encrypted","type":"uint256"},{"name":"_cryptoModule","type":"address"}],"name":"newRegDatum","outputs":[{"name":"_regID","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"regData","outputs":[{"name":"nodeSender","type":"address"},{"name":"dataHash","type":"bytes"},{"name":"fileName","type":"string"},{"name":"description","type":"string"},{"name":"encrypted","type":"uint256"},{"name":"cryptoModule","type":"address"},{"name":"linkFile","type":"string"},{"name":"creationDate","type":"uint256"},{"name":"active","type":"bool"}],"payable":false,"type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_regID","type":"uint256"}],"name":"evDataEncrypted","type":"event"}]
)
pki = pki_c.at("0xf1918d06d7e66e60153d7109ff380b41866ba2e0")

// creating the contract
pki_source = 'pragma solidity ^0.4.0;   contract LuxUni_PKI {      bytes public caCertificate;     address public caAddr;      uint public numRegData;     RegDatum[] public regData;     mapping (uint => DeletedRegDatum) public deletedRegData;       /* mapping (uint => RevokedRegDatum) revokedRegData; -- presently DeletedRegData instead */      struct RegDatum {         address nodeSender;         bytes dataHash;         address contrAddr;         string fileName;         string description;         string linkFile;         uint creationDate;         bool active;         Confirm[] confirms;         mapping (address => bool) confirmed;     }      struct DeletedRegDatum {         address nodeSender;         uint deletionDate;     }      struct RegIDStruct {         uint regID;         uint creationDate;     }      struct Confirm {         int position;         address Confirmer;     }      function LuxUni_PKI(address _addr) {         if (_addr == 0) {             _addr = msg.sender;         }         caAddr = _addr;     }      function populateCertificate(bytes _cert){         caCertificate = _cert;     }      function newRegDatum(bytes _dataHash, address _contrAddr, string _fileName, string _description, string _linkFile, uint _encrypted, address _cryptoModule) returns (uint _regID) {         if (msg.sender != caAddr) {             throw;         }          _regID = regData.length++;         RegDatum reg = regData[_regID];         reg.nodeSender = msg.sender;         reg.dataHash = _dataHash;         reg.contrAddr = _contrAddr;         reg.fileName = _fileName;         reg.description = _description;         reg.contrAddr = _contrAddr;         reg.linkFile = _linkFile;         reg.creationDate = now;         reg.active = true;         numRegData = _regID+1;      }      /* !! our black list is based on the white list.          If someone wants to put unknown certificate to the blacklist, he has to add it to the white list first and then immidiately put it into the black list.          Is it OK or shall we make both lists separately? */     function deleteRegDatum(uint _regID) returns (uint err) {         if (msg.sender != caAddr) {             throw;         }          if (_regID >= numRegData) {            return 1;         }         if (deletedRegData[_regID].deletionDate != 0) {            return 2;         }         deletedRegData[_regID] = DeletedRegDatum(msg.sender, now);         return 0;     }       /*     function strConcat(string _a, string _b, string _c, string _d, string _e) internal returns (string){         bytes memory _ba = bytes(_a);         bytes memory _bb = bytes(_b);         bytes memory _bc = bytes(_c);         bytes memory _bd = bytes(_d);         bytes memory _be = bytes(_e);         string memory abcde = new string(_ba.length + _bb.length + _bc.length + _bd.length + _be.length);         bytes memory babcde = bytes(abcde);         uint k = 0;         for (uint i = 0; i < _ba.length; i++) babcde[k++] = _ba[i];         for (i = 0; i < _bb.length; i++) babcde[k++] = _bb[i];         for (i = 0; i < _bc.length; i++) babcde[k++] = _bc[i];         for (i = 0; i < _bd.length; i++) babcde[k++] = _bd[i];         for (i = 0; i < _be.length; i++) babcde[k++] = _be[i];         return string(babcde);     }     */ }  '

pki_Compiled = eth.compile.solidity(pki_source);

pki_Contract = eth.contract(pki_Compiled['<stdin>:LuxUni_PKI'].info.abiDefinition);
//pki_Contract = eth.contract(pki_abi);

pki = pki_Contract.new(
    {
      from:web3.eth.accounts[0], 
      data:pki_Compiled['<stdin>:LuxUni_PKI'].code, 
      gas: 3000000
    }, function(e, contract){
      if(!e) {
        if(!contract.address) {
          console.log("Contract transaction send: TransactionHash: " + contract.transactionHash + " waiting to be mined...");

        } else {
          console.log("Contract mined! Address: " + contract.address);
          console.log(contract);
        }
      }      
    })

// contract addr: 0xbf6f4ded6fa3724ddd4f28351ede893575803f98
// OLD -- Contract mined! Address: 0xf1918d06d7e66e60153d7109ff380b41866ba2e0
// RADULA CONTRACT 0x799e30c8873658039c5ff18eeb21dc2eec84f310
// TESTNET CONTRACT 0x5a6d8e8db0130ba82396f8a56f461d607705ac11 // BALANCE after contr mining
//4.98060668
//4.97589898537300933

pki.newRegDatum.sendTransaction(
    "","1234","", "File7.pdf", "This is the 7th trial from geth",
    {from:eth.coinbase, gas:2000000},
    function (e, result){
       if(!e){
          console.log("newRegDatum id: " + result.args.regID);
       }
    } 
)



