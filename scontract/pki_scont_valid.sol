pragma solidity ^0.4.0;
import "./pki_scont.sol";


contract LuxUni_validateCert {

    /*
        _addrRootCA - the address of the root CA, we check if the parentAddr == 0 exactly at _addrRootCA
        returns:
        the last byte contains the error code
            0  - OK, where n is the number of iteraction
            1  - certificate not found
            2  - certificate revoked
            11 - error in parsing
            12 - CA addr in the certificate does not correspond to _addrCA
            13 - empty cert received for this CA
            14 - empty addrCA parsed in this CA cert
            15 - parent addr is null, but CA addr does not correspond to Root addr  
            16 - too many iterations: limit (100?) is exceeded
        other bytes contain the number of iteractions (path to the root ot to the level at which the error occured)
    */
    function CheckCert(bytes32 _newHash, address _addrCA,
        address _addrRoot) constant returns(uint _result) {
        uint i;
        uint _resCheck;
        address _addrParent;
        address _addrNewCA = _addrCA;

        LuxUni_PKI _contrCA = LuxUni_PKI(_addrCA);

        for (i = 0; i < 2000; i++) {
            (_resCheck, _addrParent, _newHash) = _contrCA.CheckCertForCA(_newHash);
            if (_resCheck != 0) {
                return EncodeReturn(i, _resCheck);
            }
            if (_addrParent == 0) {
                if (_addrNewCA != _addrRoot) {
                    return EncodeReturn(i, 15);
                } else {
                    return EncodeReturn(i, 0);
                }
            }
            _addrNewCA = _addrParent;
            _contrCA = LuxUni_PKI(_addrNewCA);
        }
        return 16;
    }

    function EncodeReturn(uint _iter, uint _err) constant returns(uint _res) {
        _res = _iter;
        _res <<= 8; /* 8 bit = 1 byte -- the size of error code*/
        _res = _res | uint(_err);
        return _res;
    }

    function DecodeReturnIter(uint _res) constant returns(uint _iter) {
        _iter = _res;
        _iter >>= 8; /* 8 bit = 1 byte -- the size of error code*/
        return _iter;
    }

    function DecodeReturnErr(uint _res) constant returns(uint _err) {
        _err = _res & 0xff; /* we look at the last byte for error code */
        return _err;
    }

    /* 
    Ethereum uses KECCAK-256. It should be noted that it does not follow 
    the FIPS-202 based standard of Keccak, which was finalized in August 2015.

    Hashing the string "testing":
    Ethereum SHA3 function in Solidity = 5f16f4c7f149ac4f9510d9cf8cf384038ad348b3bcdc01915f95de12df9d1b02
    Keccak-256 (Original Padding) = 5f16f4c7f149ac4f9510d9cf8cf384038ad348b3bcdc01915f95de12df9d1b02
    SHA3-256 (NIST Standard) = 7f5979fb78f082e8b1c676635db8795c4ac6faba03525fb708cb5fd68fd40c5e

    More info:
    https://github.com/ethereum/EIPs/issues/59
    http://ethereum.stackexchange.com/questions/550/which-cryptographic-hash-function-does-ethereum-use    
    */

}
