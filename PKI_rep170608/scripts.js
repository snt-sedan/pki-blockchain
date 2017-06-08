testnet 
wallet/user "0x78714b8663f11a388b9b74c3edbc42d253d2b0da"



abigen --abi pki_abi.json --pkg main --type LuxUni_PKI --out bind_pki.go --bin pki_scont.bin

cd ./.go/src/github.com/ethereum/go-ethereum/
godep go build ~/DocsFS/Dropbox/WORK/RD/LuxBCh/PKI/pki-web.go ~/DocsFS/Dropbox/WORK/RD/LuxBCh/PKI/bind_pki.go

//GETH

primary	= eth.accounts[0];
web3.fromWei(eth.getBalance(primary), "ether");
personal.unlockAccount(primary);

// getting access to the contract from geth
//# inside parathesis is abi.json file for this contract
pki_abi = '[{
    constant: false,
    inputs: [{
        name: "_cryptoID",
        type: "uint256"
    }, {
        name: "_bData",
        type: "bytes"
    }],
    name: "encryptCallBack",
    outputs: [],
    payable: false,
    type: "function"
}, {
    constant: true,
    inputs: [{
        name: "",
        type: "uint256"
    }],
    name: "deletedRegData",
    outputs: [{
        name: "nodeSender",
        type: "address"
    }, {
        name: "deletionDate",
        type: "uint256"
    }],
    payable: false,
    type: "function"
}, {
    constant: false,
    inputs: [{
        name: "_regID",
        type: "uint256"
    }],
    name: "deleteRegDatum",
    outputs: [{
        name: "err",
        type: "uint256"
    }],
    payable: false,
    type: "function"
}, {
    constant: true,
    inputs: [],
    name: "founder",
    outputs: [{
        name: "",
        type: "address"
    }],
    payable: false,
    type: "function"
}, {
    constant: true,
    inputs: [{
        name: "",
        type: "address"
    }],
    name: "deletedCAData",
    outputs: [{
        name: "nodeSender",
        type: "address"
    }, {
        name: "deletionDate",
        type: "uint256"
    }],
    payable: false,
    type: "function"
}, {
    constant: false,
    inputs: [{
        name: "_CA",
        type: "address"
    }, {
        name: "_parent",
        type: "address"
    }, {
        name: "_dataHash",
        type: "bytes"
    }, {
        name: "_fileName",
        type: "string"
    }, {
        name: "_description",
        type: "string"
    }, {
        name: "_linkFile",
        type: "string"
    }],
    name: "newCADatum",
    outputs: [{
        name: "_error",
        type: "uint256"
    }],
    payable: false,
    type: "function"
}, {
    constant: true,
    inputs: [],
    name: "numRegData",
    outputs: [{
        name: "",
        type: "uint256"
    }],
    payable: false,
    type: "function"
}, {
    constant: true,
    inputs: [],
    name: "numCAData",
    outputs: [{
        name: "",
        type: "uint256"
    }],
    payable: false,
    type: "function"
}, {
    constant: true,
    inputs: [{
        name: "",
        type: "uint256"
    }],
    name: "encryptRegData",
    outputs: [{
        name: "nodeSender",
        type: "address"
    }, {
        name: "data",
        type: "bytes"
    }, {
        name: "encryptDate",
        type: "uint256"
    }],
    payable: false,
    type: "function"
}, {
    constant: true,
    inputs: [{
        name: "",
        type: "address"
    }],
    name: "CAData",
    outputs: [{
        name: "nodeSender",
        type: "address"
    }, {
        name: "CAParent",
        type: "address"
    }, {
        name: "dataHash",
        type: "bytes"
    }, {
        name: "fileName",
        type: "string"
    }, {
        name: "description",
        type: "string"
    }, {
        name: "linkFile",
        type: "string"
    }, {
        name: "creationDate",
        type: "uint256"
    }],
    payable: false,
    type: "function"
}, {
    constant: false,
    inputs: [{
        name: "_CA",
        type: "address"
    }],
    name: "deleteCADatum",
    outputs: [{
        name: "error",
        type: "uint256"
    }],
    payable: false,
    type: "function"
}, {
    constant: false,
    inputs: [{
        name: "_dataHash",
        type: "bytes"
    }, {
        name: "_fileName",
        type: "string"
    }, {
        name: "_description",
        type: "string"
    }, {
        name: "_linkFile",
        type: "string"
    }, {
        name: "_encrypted",
        type: "uint256"
    }, {
        name: "_cryptoModule",
        type: "address"
    }],
    name: "newRegDatum",
    outputs: [{
        name: "_regID",
        type: "uint256"
    }],
    payable: false,
    type: "function"
}, {
    constant: true,
    inputs: [{
        name: "",
        type: "uint256"
    }],
    name: "regData",
    outputs: [{
        name: "nodeSender",
        type: "address"
    }, {
        name: "dataHash",
        type: "bytes"
    }, {
        name: "fileName",
        type: "string"
    }, {
        name: "description",
        type: "string"
    }, {
        name: "encrypted",
        type: "uint256"
    }, {
        name: "cryptoModule",
        type: "address"
    }, {
        name: "linkFile",
        type: "string"
    }, {
        name: "creationDate",
        type: "uint256"
    }, {
        name: "active",
        type: "bool"
    }],
    payable: false,
    type: "function"
}, {
    inputs: [],
    payable: false,
    type: "constructor"
}, {
    anonymous: false,
    inputs: [{
        indexed: false,
        name: "_regID",
        type: "uint256"
    }],
    name: "evDataEncrypted",
    type: "event"
}]
'

pki_c = eth.contract(
[{"constant":false,"inputs":[{"name":"_cryptoID","type":"uint256"},{"name":"_bData","type":"bytes"}],"name":"encryptCallBack","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"deletedRegData","outputs":[{"name":"nodeSender","type":"address"},{"name":"deletionDate","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_regID","type":"uint256"}],"name":"deleteRegDatum","outputs":[{"name":"err","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"founder","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"LuxUni_KYC","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"numRegData","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"encryptRegData","outputs":[{"name":"nodeSender","type":"address"},{"name":"data","type":"bytes"},{"name":"encryptDate","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_dataHash","type":"bytes"},{"name":"_fileName","type":"string"},{"name":"_description","type":"string"},{"name":"_linkFile","type":"string"},{"name":"_encrypted","type":"uint256"},{"name":"_cryptoModule","type":"address"}],"name":"newRegDatum","outputs":[{"name":"_regID","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"regData","outputs":[{"name":"nodeSender","type":"address"},{"name":"dataHash","type":"bytes"},{"name":"fileName","type":"string"},{"name":"description","type":"string"},{"name":"encrypted","type":"uint256"},{"name":"cryptoModule","type":"address"},{"name":"linkFile","type":"string"},{"name":"creationDate","type":"uint256"},{"name":"active","type":"bool"}],"payable":false,"type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_regID","type":"uint256"}],"name":"evDataEncrypted","type":"event"}]
)
pki = pki_c.at("0xf1918d06d7e66e60153d7109ff380b41866ba2e0")

// creating the contract
pki_source = 'pragma solidity ^0.4.0; contract LUCrypProxy { function encryptCallBack(uint _cryptoID, bytes _bData); } contract LUCrypModule { address public founder; uint public numRegData; RegDatum[] public regData; mapping (uint => DeletedRegDatum) public deletedRegData; string[] public errLog; uint public errLogNum; event evSendData(uint ind, bytes _bData); struct RegDatum { address nodeSender; bytes bData; string descript; uint creationDate; } struct DeletedRegDatum { address nodeSender; uint deletionDate; } function LUCrypModule(); function encryptRequest(bytes _data) returns (uint _regID); function encryptResponse(uint _regID, bytes _data) returns (uint err); function deleteRegDatum(uint _regID) internal returns (uint err); function logError(string _sError) internal returns (uint _logID); } contract LuxUni_PKI is LUCrypProxy { address public founder; uint public numRegData; RegDatum[] public regData; mapping (uint => DeletedRegDatum) public deletedRegData; mapping (uint => EncryptRegDatum) public encryptRegData; mapping (uint => RegIDStruct) internal crypto2regID; uint public numCAData; mapping (address => CADatum) public CAData; mapping (address => DeletedCADatum) public deletedCAData; string[] newCAErrorLog; /* mapping (uint => RevokedRegDatum) revokedRegData; -- presently DeletedRegData instead */ event evDataEncrypted(uint _regID); /* "encrypted" is binary flags indicating what fields are encrypted and which fields should be looked at encryptRegData */ struct RegDatum { address nodeSender; bytes dataHash; string fileName; string description; uint encrypted; address cryptoModule; string linkFile; uint creationDate; bool active; Confirm[] confirms; mapping (address => bool) confirmed; } struct DeletedRegDatum { address nodeSender; uint deletionDate; } struct CADatum { /* address CAAddress; the address of the CA which has been given this certificate */ address nodeSender; /* the one who approved this CA certificate (root or parent) */ address CAParent; /* the address of the parent in this array */ bytes dataHash; string fileName; string description; string linkFile; uint creationDate; mapping (address => bool) children; /* is it used? bool - if the certificate is revoked*/ } struct DeletedCADatum { address nodeSender; uint deletionDate; } /* data is an encryption of a several string fields (separated with char "0") in accordance with uint encrypted flags */ struct EncryptRegDatum { address nodeSender; bytes data; uint encryptDate; } struct RegIDStruct { uint regID; uint creationDate; } struct Confirm { int position; address Confirmer; } function LuxUni_PKI() { founder = msg.sender; newCADatum(founder, 0, "", "", "", ""); } function newRegDatum(bytes _dataHash, string _fileName, string _description, string _linkFile, uint _encrypted, address _cryptoModule) returns (uint _regID) { _regID = regData.length++; RegDatum reg = regData[_regID]; reg.nodeSender = msg.sender; reg.dataHash = _dataHash; reg.fileName = _fileName; reg.linkFile = _linkFile; reg.encrypted = _encrypted; reg.cryptoModule = _cryptoModule; reg.creationDate = now; reg.active = true; numRegData = _regID+1; if( _encrypted == 0 ) { reg.description = _description; } else { LUCrypModule theModule = LUCrypModule(_cryptoModule); uint _cryptoID = theModule.encryptRequest( bytes(_description) ); crypto2regID[ _cryptoID ] = RegIDStruct( _regID, now ); } } function deleteRegDatum(uint _regID) returns (uint err) { if (_regID >= numRegData) { return 1; } if (deletedRegData[_regID].deletionDate != 0) { return 2; } deletedRegData[_regID] = DeletedRegDatum(msg.sender, now); return 0; } function encryptCallBack(uint _cryptoID, bytes _bData) { uint _regID; if (crypto2regID[_cryptoID].creationDate == 0) { return; } _regID = crypto2regID[_cryptoID].regID; if (_regID >= numRegData) { throw; } if (encryptRegData[_regID].encryptDate != 0) { return; } encryptRegData[_regID] = EncryptRegDatum(msg.sender, _bData, now); evDataEncrypted(_regID); } /* returns uint _error = 0 in case of success. _error = 1 - the msg.sender not found in the CAData. Error desc can be found in newCADatumLog */ function newCADatum(address _CA, address _parent, bytes _dataHash, string _fileName, string _description, string _linkFile) returns (uint _error) { if (CAData[msg.sender].creationDate == 0) { return 1; } if (_parent == 0) { _parent = msg.sender; } CAData[_CA] = CADatum( msg.sender, /* reg.nodeSender = msg.sender */ _parent, /* CAparent = _parent */ _dataHash, /* reg.dataHash = _dataHash; */ _fileName, /* reg.fileName = _fileName; */ _description, /* reg.description = _description; */ _linkFile, /* reg.linkFile = _linkFile; */ now ); } function deleteCADatum(address _CA) returns (uint error) { if (deletedCAData[_CA].deletionDate != 0) { return 2; } if (CAData[_CA].CAParent != msg.sender || founder!=msg.sender) { return 1; } deletedCAData[_CA] = DeletedCADatum(msg.sender, now); return 0; } /* function strConcat(string _a, string _b, string _c, string _d, string _e) internal returns (string){ bytes memory _ba = bytes(_a); bytes memory _bb = bytes(_b); bytes memory _bc = bytes(_c); bytes memory _bd = bytes(_d); bytes memory _be = bytes(_e); string memory abcde = new string(_ba.length + _bb.length + _bc.length + _bd.length + _be.length); bytes memory babcde = bytes(abcde); uint k = 0; for (uint i = 0; i < _ba.length; i++) babcde[k++] = _ba[i]; for (i = 0; i < _bb.length; i++) babcde[k++] = _bb[i]; for (i = 0; i < _bc.length; i++) babcde[k++] = _bc[i]; for (i = 0; i < _bd.length; i++) babcde[k++] = _bd[i]; for (i = 0; i < _be.length; i++) babcde[k++] = _be[i]; return string(babcde); } */ } '

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

//Contract mined! Address: 0xf1918d06d7e66e60153d7109ff380b41866ba2e0
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



