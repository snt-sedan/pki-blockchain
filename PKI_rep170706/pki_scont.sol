pragma solidity ^0.4.0;

contract owned {
    function owned() { owner = msg.sender; }
    address internal owner;

    // This contract only defines a modifier but does not use
    // it - it will be used in derived main contracts.
    // The function body of the main contract is inserted where the special symbol
    // "_;" in the definition of a modifier appears.
    // If the owner calls this function, the function is executed
    // and otherwise, an exception is thrown.
    modifier onlyOwner {
        require(msg.sender == owner);
        _;
    }
}

contract LuxUni_PKI is owned {

    bytes private caCertificate;
    function getCaCertificate() constant returns(bytes) { return caCertificate; }

    uint private numRegData;
    function getNumRegData() constant returns(uint) { return numRegData; }

    RegDatum[] private regData;
    function getRegNodeSender(uint _i) constant returns( address ) { return regData[_i].nodeSender; }
    function getRegDataHash(uint _i) constant returns( bytes ) { return regData[_i].dataHash; }
    function getRegContrAddr(uint _i) constant returns( address ) { return regData[_i].contrAddr; }
    function getRegFileName(uint _i) constant returns( string ) { return regData[_i].fileName; }
    function getRegDescription(uint _i) constant returns( string ) { return regData[_i].description; }
    function getRegLinkFile(uint _i) constant returns( string ) { return regData[_i].linkFile; }
    function getRegCreationDate(uint _i) constant returns( uint ) { return regData[_i].creationDate; }

    mapping (uint => DeletedRegDatum) private deletedRegData;
    function getDeletedRegNodeSender(uint _i) constant returns(address) { return deletedRegData[_i].nodeSender; }
    function getDeletedRegDate(uint _i) constant returns(uint) { return deletedRegData[_i].deletionDate; }

    struct RegDatum {
        address nodeSender;
        bytes dataHash;
        address contrAddr;
        string fileName;
        string description;
        string linkFile;
        uint creationDate;
    }

    struct DeletedRegDatum {
        address nodeSender;
        uint deletionDate;
    }

    /*struct RegIDStruct {
        uint regID;
        uint creationDate;
    }*/

    function getOwner() constant returns(address) { return owner; }
    function setOwner(address _addr) { owner = _addr; }

    function populateCertificate(bytes _cert)  onlyOwner {
        caCertificate = _cert;
    }

    /* nodeSender - the account of subCA or an odinary user
       returns:
           ID of the new register record
           throws if executed by not an owner
    */
    function newRegDatum(bytes _dataHash, address _contrAddr, string _fileName, string _description,
            string _linkFile, address _nodeSender) onlyOwner returns(uint _regID) {

        _regID = regData.length++;
        RegDatum reg = regData[_regID];
        reg.nodeSender = _nodeSender;
        reg.dataHash = _dataHash;
        reg.contrAddr = _contrAddr;
        reg.fileName = _fileName;
        reg.description = _description;
        reg.contrAddr = _contrAddr;
        reg.linkFile = _linkFile;
        reg.creationDate = now;
        numRegData = _regID+1;

    }

    /* !! our black list is based on the white list.
         If someone wants to put unknown certificate to the blacklist,
             he has to add it to the white list first and then immidiately put it into the black list.
         Is it OK or shall we make both lists separately? */
    function deleteRegDatum(uint _regID) onlyOwner returns (uint err) {

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

