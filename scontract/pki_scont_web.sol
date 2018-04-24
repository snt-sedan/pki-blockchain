pragma solidity ^0.4.0;
import "./pki_scont.sol";


contract LuxUni_PKI_web is owned {

    /* uint256 -- 32 bytes
       address -- 20 bytes
       uint256 in mapping below : 
            first 20 bytes -- address of parent smart contract holding the hashes in regData arrays
            second 12 bytes -- uint96: the index of RegData array of the corresponding 
    */
    function EncodeMapID(address _addrContr, uint96 _index) constant returns(uint256 _res) {
        _res = _index;
        _res <<= 8 * 20; /* 20 bytes -- the length of address type in bytes */
        _res = _res | uint(_addrContr);
        return _res;
    }

    /* function DecodeMapID (uint256) constant returns((address _addrContr, uint96 _index)) {} */

    mapping(uint256 => RegDatum) private regData;

    function getRegEthAccCA(address _addrContr, uint96 _index) constant returns(address) {
        return regData[EncodeMapID(_addrContr, _index)].ethAccCA;
    }

    function getRegContrAddr(address _addrContr, uint96 _index) constant returns(address) {
        return regData[EncodeMapID(_addrContr, _index)].contrAddr;
    }

    function getRegFileName(address _addrContr, uint96 _index) constant returns(string) {
        return regData[EncodeMapID(_addrContr, _index)].fileName;
    }

    function getRegDescription(address _addrContr, uint96 _index) constant returns(string) {
        return regData[EncodeMapID(_addrContr, _index)].description;
    }

    function getRegCreationDate(address _addrContr, uint96 _index) constant returns(uint) {
        return regData[EncodeMapID(_addrContr, _index)].creationDate;
    }

    struct RegDatum {
        address ethAccCA; /* ethAccCA - the user account of subCA which controls the corresponding smart contract */
        address contrAddr; /* address of the smart contract of CA - allows to browse CA tree from the root to leafs*/
        string fileName; /* name of the certificate file loaded - for information purposes*/
        string description; /* for information purposes */
        uint creationDate;
    }

    function newRegDatum(address _parentAddr, uint _arrInd, address _ethAccCA, address _contrAddr,
        string _fileName, string _description) onlyOwner returns(uint err) {

        uint _regID = EncodeMapID(_parentAddr, uint96(_arrInd));
        // assert should be used
        if (regData[_regID].creationDate != 0) {
            return 1;
        }

        RegDatum memory reg;
        reg.ethAccCA = _ethAccCA;
        reg.contrAddr = _contrAddr;
        reg.fileName = _fileName;
        reg.description = _description;
        reg.creationDate = now;

        regData[_regID] = reg;
        return 0;
    }

}

