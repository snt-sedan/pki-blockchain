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
