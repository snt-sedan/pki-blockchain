pragma solidity ^0.4.0;
import "./owned.sol";


contract LuxUni_PKI is owned {

    bytes private caCertificate;

    function getCaCertificate() constant returns(bytes) { return caCertificate; }

    uint private numRegData;

    function getNumRegData() constant returns(uint) { return numRegData; }

    RegDatum[] private regData;

    function getRegDataHash(uint _i) constant returns(bytes32) { return regData[_i].dataHash; }

    function getRegAlgoHashID(uint _i) constant returns(bytes) { return regData[_i].algoHashID; }

    function getRegCreationDate(uint _i) constant returns(uint) { return regData[_i].creationDate; }

    mapping(uint => DeletedRegDatum) private deletedRegData;

    function getDeletedRegNodeSender(uint _i) constant returns(address) { return deletedRegData[_i].nodeSender; }

    function getDeletedRegDate(uint _i) constant returns(uint) { return deletedRegData[_i].deletionDate; }

    struct RegDatum {
        bytes32 dataHash;
        bytes algoHashID;
        uint creationDate;
    }

    struct DeletedRegDatum {
        address nodeSender;
        uint deletionDate;
    }

    event evLuxUni_NewRegDatumReturn(uint _ind, uint _regID);

    function getOwner() constant returns(address) { return owner; }

    function setOwner(address _addr) onlyOwner { owner = _addr; }

    function populateCertificate(bytes _cert) onlyOwner {
        caCertificate = _cert;
    }

    /* nodeSender - the user account of subCA or an odinary user
       algoHashID - "" (default) for Keccak256, otherwise OID (ASN.1)
       // https://tools.ietf.org/html/rfc3279
       // http://www.alvestrand.no/objectid/1.3.14.3.2.26.html
       returns:
           ID of the new register record
           throws if executed by not an owner
    */
    function newRegDatum(bytes32 _dataHash, bytes _algoHashID) onlyOwner returns(uint _regID) {

        _regID = regData.length++;
        RegDatum reg = regData[_regID];
        reg.dataHash = _dataHash;
        reg.algoHashID = _algoHashID;
        reg.creationDate = now;
        numRegData = _regID + 1;

        evLuxUni_NewRegDatumReturn(uint(_dataHash), _regID);
    }

    /* !! our black list is based on the white list.
         If someone wants to put unknown certificate to the blacklist,
             he has to add it to the white list first and then immidiately put it into the black list.
         Is it OK or shall we make both lists separately? */
    function deleteRegDatum(uint _regID) onlyOwner returns(uint err) {

        if (_regID >= numRegData) {
            return 1;
        }
        if (deletedRegData[_regID].deletionDate != 0) {
            return 2;
        }
        deletedRegData[_regID] = DeletedRegDatum(msg.sender, now);
        return 0;
    }


    struct Asn1Item {
        uint iFirst; /* index of the first byte -- ixs */
        uint iFirstContent; /* index of the first content byte -- ixf */
        uint iLastContent; /* index of the last content byte -- ixl */
    }

    /* The following 4 functions are all you need to parse an ASN1 structure */

    /* # gets the first ASN1 structure in der */
    function asn1_node_root(bytes der) private constant returns(Asn1Item) {
        return ReadLength(der, 0);
    }

    /* # gets the next ASN1 structure following (ixs,ixf,ixl) */
    function asn1_node_next(bytes der, Asn1Item _asnCurrItem) private constant returns(Asn1Item) {
        return ReadLength(der, _asnCurrItem.iLastContent + 1);
    }

    /* # opens the container (ixs,ixf,ixl) and returns the first ASN1 inside */
    function asn1_node_first_child(bytes der, Asn1Item _asnCurrItem)
    private constant returns(Asn1Item) {
        if ((der[_asnCurrItem.iFirst] & 0x20) != 0x20) {
            /* raise ValueError('Error: can only open constructed types. '
                                +'Found type: 0x'+der[ixs].encode("hex")) */
            throw;
        }
        return ReadLength(der, _asnCurrItem.iFirstContent /*ixf*/ );
    }

    /* # is true if one ASN1 chunk is inside another chunk. (ixs,ixf,ixl), (jxs,jxf,jxl) */
    function asn1_node_is_child_of(Asn1Item _asnParent, Asn1Item _asnChild)
    private constant returns(bool) {
        /* return ( (ixf <= jxs ) and (jxl <= ixl) ) or \
           ( (jxf <= ixs ) and (ixl <= jxl) ) */
        /* return ( (_asnParent.iFirstContent <= _asnChild.iFirst ) || 
            (_asnChild.iLastContent <= _asnParent.iLastContent) ); */
        if (_asnChild.iLastContent > _asnParent.iLastContent) {
            return false;
        }
        if (_asnChild.iFirst < _asnParent.iFirstContent) {
            return false;
        }
        return true;
    }

    /* ##### END NAVIGATE */


    /* ##### ACCESS PRIMITIVES */

    /* # get content and verify type byte */
    /*function asn1_get_value_of_type(Asn1Item _asnParent,asn1_type) {
        asn1_type_table = {
        'BOOLEAN':           0x01,        'INTEGER':           0x02,
        'BIT STRING':        0x03,        'OCTET STRING':      0x04,
        'NULL':              0x05,        'OBJECT IDENTIFIER': 0x06,
        'SEQUENCE':          0x70,        'SET':               0x71,
        'PrintableString':   0x13,        'IA5String':         0x16,
        'UTCTime':           0x17,        'ENUMERATED':        0x0A,    
        'UTF8String':        0x0C,        'PrintableString':   0x13,
        }
        if asn1_type_table[asn1_type] != ord(der[ixs]):
                raise ValueError('Error: Expected type was: '+
                        hex(asn1_type_table[asn1_type])+
                        ' Found: 0x'+der[ixs].encode('hex'))
        return der[ixf:ixl+1]
    }*/

    function BytesSubArr(bytes _a, uint _iStart, uint _iEnd) private constant returns(bytes) {
        if (_iEnd > _a.length) {
            throw;
        }
        bytes memory _ret = new bytes(_iEnd + 1 - _iStart);
        for (uint _i = _iStart; _i < _iEnd + 1; _i++) { _ret[_i - _iStart] = _a[_i]; }
        return _ret;
    }

    /* # get value */
    function asn1_get_value(bytes der, Asn1Item _asnItem) private constant returns(bytes) {
        return BytesSubArr(der, _asnItem.iFirstContent, _asnItem.iLastContent); /* der[ixf:ixl+1] */
    }

    /* # get type+length+value */
    function asn1_get_all(bytes der, Asn1Item _asnItem) private constant returns(bytes) {
        return BytesSubArr(der, _asnItem.iFirst, _asnItem.iLastContent); /* der[ixs:ixl+1] */
    }

    /* # get tag */
    function asn1_get_tag(bytes der, Asn1Item _asnItem) private constant returns(byte) {
        //return BytesSubArr(der, _asnItem.iFirst, _asnItem.iFirstContent);   /* der[ixs:ixs+1] ??? */
        return byte(der[_asnItem.iFirst]); /* der[ixs:ixs+1] ??? */
    }

    /* ##### END ACCESS PRIMITIVES */


    /* ##### HELPER FUNCTIONS */
    /*
      # converter
    def bitstr_to_bytestr(bitstr):
        if bitstr[0] != '\x00':
                raise ValueError('Error: only 00 padded bitstr can be converted to bytestr!')
        return bitstr[1:]
    */

    function stringsEqual(string memory _a, string memory _b) internal constant returns(bool) {
        bytes memory a = bytes(_a);
        bytes memory b = bytes(_b);
        if (a.length != b.length)
            return false;
        // @todo unroll this loop
        for (uint i = 0; i < a.length; i++)
            if (a[i] != b[i])
                return false;
        return true;
    }

    function bytesEqual(bytes memory a, bytes memory b) internal constant returns(bool) {
        if (a.length != b.length)
            return false;
        // @todo unroll this loop
        for (uint i = 0; i < a.length; i++)
            if (a[i] != b[i])
                return false;
        return true;
    }

    function bytes32Equal(bytes32 a, bytes32 b) internal constant returns(bool) {
        for (uint i = 0; i < a.length; i++)
            if (a[i] != b[i])
                return false;
        return true;
    }

    function bytesToAddress(bytes b) public constant returns(address) {
        uint result = bytestr_to_uint(b);
        return address(result);
    }

    /* # converter */
    function bytestr_to_uint(bytes _ba) constant returns(uint) {
        /* # converts bytestring to integer */
        uint _ret = 0;
        for (uint _i = 0; _i < _ba.length; _i++) {
            _ret <<= 8;
            _ret = _ret | uint(_ba[_i]);
        }
        return _ret;
    }

    /* # ix points to the first byte of the asn1 structure
       # Returns first byte pointer, first content byte pointer and last.  */
    function ReadLength(bytes der, uint ix) private constant returns(Asn1Item _asnItem) {
        uint ix_first_content_byte;
        uint ix_last_content_byte;
        byte first = der[ix + 1];
        if ((der[ix + 1] & 0x80) == 0) {
            uint length = uint(first);
            ix_first_content_byte = ix + 2;
            ix_last_content_byte = ix_first_content_byte + uint(length) - 1;
        } else {
            uint lengthbytes = uint(first & 0x7F);
            length = bytestr_to_uint(BytesSubArr(der, ix + 2, ix + 1 + lengthbytes)); /* der[ix+2:ix+2+lengthbytes] */
            ix_first_content_byte = ix + 2 + lengthbytes;
            ix_last_content_byte = ix_first_content_byte + length - 1;
        }

        _asnItem.iFirst = ix;
        _asnItem.iFirstContent = ix_first_content_byte;
        _asnItem.iLastContent = ix_last_content_byte;

        return _asnItem;
    }
    /*  ##### END HELPER FUNCTIONS */

    function hexStrToBytes(string hex_str) public constant returns(bytes) {
        //Check hex string is valid
        if (bytes(hex_str)[0] != '0' ||
            bytes(hex_str)[1] != 'x' ||
            bytes(hex_str).length % 2 != 0 ||
            bytes(hex_str).length < 4) {
            throw;
        }

        bytes memory bytes_array = new bytes((bytes(hex_str).length - 2) / 2);

        for (uint i = 2; i < bytes(hex_str).length; i += 2) {
            uint tetrad1 = 16;
            uint tetrad2 = 16;

            //left digit
            if (uint(bytes(hex_str)[i]) >= 48 && uint(bytes(hex_str)[i]) <= 57)
                tetrad1 = uint(bytes(hex_str)[i]) - 48;

            //right digit
            if (uint(bytes(hex_str)[i + 1]) >= 48 && uint(bytes(hex_str)[i + 1]) <= 57)
                tetrad2 = uint(bytes(hex_str)[i + 1]) - 48;

            //left A->F
            if (uint(bytes(hex_str)[i]) >= 65 && uint(bytes(hex_str)[i]) <= 70)
                tetrad1 = uint(bytes(hex_str)[i]) - 65 + 10;

            //right A->F
            if (uint(bytes(hex_str)[i + 1]) >= 65 && uint(bytes(hex_str)[i + 1]) <= 70)
                tetrad2 = uint(bytes(hex_str)[i + 1]) - 65 + 10;

            //left a->f
            if (uint(bytes(hex_str)[i]) >= 97 && uint(bytes(hex_str)[i]) <= 102)
                tetrad1 = uint(bytes(hex_str)[i]) - 97 + 10;

            //right a->f
            if (uint(bytes(hex_str)[i + 1]) >= 97 && uint(bytes(hex_str)[i + 1]) <= 102)
                tetrad2 = uint(bytes(hex_str)[i + 1]) - 97 + 10;

            //Check all symbols are allowed
            if (tetrad1 == 16 || tetrad2 == 16)
                throw;

            bytes_array[i / 2 - 1] = byte(16 * tetrad1 + tetrad2);
        }

        return bytes_array;
    }

    /* ####################### END ASN1 DECODER ############################ */


    function ParseAddrCA(bytes _der) public constant
    returns(address _addrCA) {
        address _addrPar;
        uint _errCode;
        (_addrPar, _addrCA, _errCode) = ParseCert(_der);
        return _addrCA;
    }

    function ParseAddrParent(bytes _der) public constant
    returns(address _addrParent) {
        address _addrCA;
        uint _errCode;
        (_addrParent, _addrCA, _errCode) = ParseCert(_der);
        return _addrParent;
    }

    function ParseCert(bytes _der) public constant
    returns(address _addrParent, address _addrCA, uint errCode) {

            bytes memory _OID;
            Asn1Item memory i;
            Asn1Item memory j;
            Asn1Item memory k;
            Asn1Item memory l;

            i = asn1_node_root(_der); /* # Get root node */

            i = asn1_node_first_child(_der, i); /* # Get first node in Certificate */
            j = asn1_node_first_child(_der, i); /* # Get first node in TBSCertificate */

            uint16 ind;
            for (ind = 0; ind < 100; ind++) {
                if (asn1_node_is_child_of(i /*parent*/ , j) == false) {
                    break;
                }
                if (asn1_get_tag(_der, j) == 0xA3) {
                    break;
                }
                /* # Loop through to find extensions tag A3 */
                j = asn1_node_next(_der, j);
            }

            k = asn1_node_first_child(_der, j); /* # Get sequence in extensions */
            k = asn1_node_first_child(_der, k); /* # Get first node in extensions */

            for (ind = 0; ind < 100; ind++) { /* # Loop through extensions */
                if (asn1_node_is_child_of(j /*parent*/ , k) == false) {
                    break;
                }
                l = asn1_node_first_child(_der, k); /* # Get OID */
                _OID = asn1_get_all(_der, l); /*.encode("hex") */

                l = asn1_node_next(_der, l); /* # Get value */
                //if (string(strOID)=="06062a8570732100") {
                /* print "HashAlgo[" + strOID + "]=" + asn1_get_value(cert_der,l).encode("hex") */
                //    _retStatus.hashAlgo = asn1_get_value(_der, l);
                //}
                if (bytesEqual(_OID, hexStrToBytes("0x06062a8570732101")) == true) {
                    _addrParent = bytesToAddress(asn1_get_value(_der, l));
                }
                if (bytesEqual(_OID, hexStrToBytes("0x06062a8570732102")) == true) {
                    _addrCA = bytesToAddress(asn1_get_value(_der, l));
                }
                //if (string(strOID)=="06062a8570732103") {
                /* print "BlockchainName[" + strOID + "]=" + asn1_get_value(cert_der,l).encode("hex") */
                //    _retStatus.blockchainName = asn1_get_value(_der, l);
                //}
                k = asn1_node_next(_der, k);
            }
            return (_addrParent, _addrCA /*addrParent, addrCaContract*/ , 0);
        }
        /* #----------------------------- */

    /*       
        returns:
          newHash - hash of the cert for this CA - we should check for this hash in 
          result:
            0  - certificate OK
            1  - certificate not found
            2  - certificate revoked
            11 - error in parsing
            12 - CA addr in the certificate does not correspond to _addrCA
            13 - empty cert received for this CA
            14 - empty addrCA parsed in this CA cert
    */
    function CheckCertForCA(bytes32 _hash) constant returns(uint _result, address _addrParent, bytes32 _newHash) {
        uint i;

        address _crtAddrParent;
        address _crtAddrCA;
        uint _errCode;
        uint _revokeDate;

        if (caCertificate.length == 0) {
            return (13, 0, "");
        }

        (_crtAddrParent, _crtAddrCA, _errCode) = ParseCert(caCertificate);
        if (_errCode != 0) { return (11, 0, ""); }
        if (_crtAddrCA == 0) { return (14, 0, ""); }
        /*if( _crtAddrCA != _addrCA ) { // important check - to do
            return (12, 0, "");
        }*/

        _newHash = sha3(caCertificate);
        _revokeDate = 0;
        for (i = 0; i < numRegData; i++) {
            if (bytes32Equal(getRegDataHash(i), _hash) == true) {
                _revokeDate = getDeletedRegDate(i);
                if (_revokeDate == 0) {
                    return (0, _crtAddrParent, _newHash);
                }
            }
        }
        if (_revokeDate != 0) {
            return (2, _crtAddrParent, _newHash);
        }
        return (1, _crtAddrParent, _newHash);
    }


}


