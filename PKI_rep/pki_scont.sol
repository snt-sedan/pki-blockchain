pragma solidity ^0.4.0;

/*
contract LUCrypProxy {
    function encryptCallBack(uint _cryptoID, bytes _bData);
}


contract LUCrypModule {
    address public founder;
    uint public numRegData;
    RegDatum[] public regData;
    mapping (uint => DeletedRegDatum) public deletedRegData;

    string[] public errLog;
    uint public errLogNum;

    event evSendData(uint ind, bytes _bData);
        
    struct RegDatum {
        address nodeSender;
        bytes bData;
        string descript;
        uint creationDate;
    }

    struct DeletedRegDatum {
        address nodeSender;
        uint deletionDate;
    }

    function LUCrypModule();
    function encryptRequest(bytes _data) returns (uint _regID);
    function encryptResponse(uint _regID, bytes _data) returns (uint err);
    function deleteRegDatum(uint _regID) internal returns (uint err);
    function logError(string _sError) internal returns (uint _logID);
}
*/

contract LuxUni_PKI /*is LUCrypProxy*/ {

    bytes ownCertificate;

    address public caAddr;

    uint public numRegData;
    RegDatum[] public regData;
    mapping (uint => DeletedRegDatum) public deletedRegData;
    mapping (uint => EncryptRegDatum) public encryptRegData;
    mapping (uint => RegIDStruct) internal crypto2regID;


    /* mapping (uint => RevokedRegDatum) revokedRegData; -- presently DeletedRegData instead */

    event evDataEncrypted(uint _regID);
    
    /* "encrypted" is binary flags indicating what fields are encrypted and which fields should be looked at encryptRegData */
    struct RegDatum {
        address nodeSender;
        bytes dataHash;
        string fileName;
        string description;
        /*uint encrypted;*/
        address cryptoModule;
        string linkFile;
        uint creationDate;
        bool active;
        Confirm[] confirms;
        mapping (address => bool) confirmed;
    }

    struct DeletedRegDatum {
        address nodeSender;
        uint deletionDate;
    }

    /* data is an encryption of a several string fields (separated with char "0") in accordance with uint encrypted flags */
    struct EncryptRegDatum {
        address nodeSender;
        bytes data;
        uint encryptDate;
    }

    struct RegIDStruct {
        uint regID;
        uint creationDate;
    }

    struct Confirm {
        int position;
        address Confirmer;
    }

    function LuxUni_PKI(address _addr) {
        if (_addr == 0) {
            _addr = msg.sender;
        }
        caAddr = _addr;
    }

    function newRegDatum(bytes _dataHash, string _fileName, string _description, string _linkFile, uint _encrypted, address _cryptoModule) returns (uint _regID) {
        if (msg.sender != caAddr) {
            throw;
        }

        _regID = regData.length++;
        RegDatum reg = regData[_regID];
        reg.nodeSender = msg.sender;
        reg.dataHash = _dataHash;
        reg.fileName = _fileName;
        reg.linkFile = _linkFile;
        reg.encrypted = _encrypted;
        reg.cryptoModule = _cryptoModule;
        reg.creationDate = now;
        reg.active = true;
        numRegData = _regID+1;

        if( _encrypted == 0 ) {
            reg.description = _description;
        } else {
            LUCrypModule theModule = LUCrypModule(_cryptoModule);
            uint _cryptoID = theModule.encryptRequest( bytes(_description) );
            crypto2regID[ _cryptoID ] = RegIDStruct( _regID, now );
        }
    }

    /* !! our black list is based on the white list.
         If someone wants to put unknown certificate to the blacklist, he has to add it to the white list first and then immidiately put it into the black list.
         Is it OK or shall we make both lists separately? */
    function deleteRegDatum(uint _regID) returns (uint err) {
        if (msg.sender != caAddr) {
            throw;
        }

        if (_regID >= numRegData) {
           return 1;
        }
        if (deletedRegData[_regID].deletionDate != 0) {
           return 2;
        }
        deletedRegData[_regID] = DeletedRegDatum(msg.sender, now);
        return 0;
    }

    function encryptCallBack(uint _cryptoID, bytes _bData) {
        uint _regID;
        if (crypto2regID[_cryptoID].creationDate == 0) {
           return;
        }
        _regID = crypto2regID[_cryptoID].regID;

        if (_regID >= numRegData) {
           throw;
        }
        if (encryptRegData[_regID].encryptDate != 0) {
           return;
        }
        encryptRegData[_regID] = EncryptRegDatum(msg.sender, _bData, now);
        evDataEncrypted(_regID);
    }


    /*
    function strConcat(string _a, string _b, string _c, string _d, string _e) internal returns (string){
        bytes memory _ba = bytes(_a);
        bytes memory _bb = bytes(_b);
        bytes memory _bc = bytes(_c);
        bytes memory _bd = bytes(_d);
        bytes memory _be = bytes(_e);
        string memory abcde = new string(_ba.length + _bb.length + _bc.length + _bd.length + _be.length);
        bytes memory babcde = bytes(abcde);
        uint k = 0;
        for (uint i = 0; i < _ba.length; i++) babcde[k++] = _ba[i];
        for (i = 0; i < _bb.length; i++) babcde[k++] = _bb[i];
        for (i = 0; i < _bc.length; i++) babcde[k++] = _bc[i];
        for (i = 0; i < _bd.length; i++) babcde[k++] = _bd[i];
        for (i = 0; i < _be.length; i++) babcde[k++] = _be[i];
        return string(babcde);
    }
    */
}

