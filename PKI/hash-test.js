pragma solidity ^0.4.8;

contract HashTest {

    function testSHA2(bytes _data) constant returns(bytes32) {
        return sha256(_data);
    }

    function testRIPEMD(bytes _data) constant returns(bytes20) {
        return ripemd160(_data);
    }
}


