pragma solidity ^0.4.0;


contract LuxUni_PKI {

    bytes public caCertificate;
    address public caAddr;

    uint public numRegData;
    RegDatum[] public regData;
    mapping (uint => DeletedRegDatum) public deletedRegData;


    /* mapping (uint => RevokedRegDatum) revokedRegData; -- presently DeletedRegData instead */

    struct RegDatum {
        address nodeSender;
        bytes dataHash;
        string fileName;
        string description;
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

    function populateCertificate(bytes _cert){
        caCertificate = _cert;
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
        reg.creationDate = now;
        reg.active = true;
        numRegData = _regID+1;

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

