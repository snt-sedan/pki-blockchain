testnet 
wallet/user "0x78714b8663f11a388b9b74c3edbc42d253d2b0da"

solc --bin nexus_test.sol
solc --abi nexus_test.sol > nexus_test.abi
abigen --abi nexus_test.abi --pkg main --type LuxUni_EET --out bind_eet.go

cd ./go/src/github.com/ethereum/go-ethereum/
godep go build ~/wrk/Dropbox/WORK/RD/Eethiq/Nexus_test/eet-web.go ~/wrk/Dropbox/WORK/RD/Eethiq/Nexus_test/bind_eet.go
godep go build ~/wrk/Dropbox/WORK/RD/Eethiq/Nexus_test/nxt-load.go ~/wrk/Dropbox/WORK/RD/Eethiq/Nexus_test/bind_eet.go

//GETH

primary	= eth.accounts[0];
web3.fromWei(eth.getBalance(primary), "ether");
personal.unlockAccount(primary);

// getting access to the contract from geth
//# inside parathesis is abi.json file for this contract
etf_abi = ''

etf_c = eth.contract(
[{"constant":false,"inputs":[{"name":"_cryptoID","type":"uint256"},{"name":"_bData","type":"bytes"}],"name":"encryptCallBack","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"deletedRegData","outputs":[{"name":"nodeSender","type":"address"},{"name":"deletionDate","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_regID","type":"uint256"}],"name":"deleteRegDatum","outputs":[{"name":"err","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"founder","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"LuxUni_KYC","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"numRegData","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"encryptRegData","outputs":[{"name":"nodeSender","type":"address"},{"name":"data","type":"bytes"},{"name":"encryptDate","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_dataHash","type":"bytes"},{"name":"_fileName","type":"string"},{"name":"_description","type":"string"},{"name":"_linkFile","type":"string"},{"name":"_encrypted","type":"uint256"},S{"name":"_cryptoModule","type":"address"}],"name":"newRegDatum","outputs":[{"name":"_regID","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"regData","outputs":[{"name":"nodeSender","type":"address"},{"name":"dataHash","type":"bytes"},{"name":"fileName","type":"string"},{"name":"description","type":"string"},{"name":"encrypted","type":"uint256"},{"name":"cryptoModule","type":"address"},{"name":"linkFile","type":"string"},{"name":"creationDate","type":"uint256"},{"name":"active","type":"bool"}],"payable":false,"type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_regID","type":"uint256"}],"name":"evDataEncrypted","type":"event"}]
)
etf = etf_c.at("0xf1918d06d7e66e60153d7109ff380b41866ba2e0")

// creating the contract
etf_source = 'pragma solidity ^0.4.8; contract Nexus_test { address owner; uint public numErrLog; string[] public errLog; struct Priority { string name; uint score; } /*uint constant numPriorities = 4;*/ Priority[4] public priorities; struct Addition { uint num; uint date; } struct Sender{ string name; uint creationDate; uint numTransactions; uint[] transactions; /* link to transactions array */ uint numPoints; Addition[] points; /* coins */ } uint public numSenders; address[] public senderList; mapping (address=>Sender) public senders; struct Receiver{ string name; uint creationDate; uint numTransactions; uint[] transactions; /* link to transactions array */ uint numPoints; Addition[] points; /* amount of priority points */ uint numPriorities; Addition[] priorities; /* links to priorities global table, TODO - add removal */ } uint public numReceivers; address[] public receiverList; mapping (address => Receiver) public receivers; struct Transaction{ address senderAddr; address receiverAddr; uint amount; uint priorityQty; uint priorityID; uint transDate; bytes certificate; } uint public numTransactions; Transaction[] public transactions; function Eethiq() { owner = msg.sender; priorities[0] = Priority("Poor",100); priorities[1] = Priority("Refugee",1000); priorities[2] = Priority("Children",3000); priorities[3] = Priority("Orphant",5000); } function isAdmin() constant returns (bool){ if (msg.sender == owner) { return true; } else { return false; } } /* returns -1 -- adress already exists */ function newReceiver(address _receiverAddr, string _name) returns (int _ret) { if(_receiverAddr == 0){ _receiverAddr = msg.sender; } if(receivers[_receiverAddr].creationDate != 0){ return -1; } receivers[_receiverAddr].name = _name; receivers[_receiverAddr].creationDate = now; receivers[_receiverAddr].numTransactions = 0; receivers[_receiverAddr].transactions = new uint[](0); receivers[_receiverAddr].numPoints = 0; /* receivers[_receiverAddr].points = new Addition[](0); */ uint _priorID = receivers[_receiverAddr].priorities.length++; receivers[_receiverAddr].priorities[_priorID].num = 0; receivers[_receiverAddr].priorities[_priorID].date = now; receivers[_receiverAddr].numPriorities = _priorID + 1; uint _receiverID = receiverList.length++; receiverList[_receiverID] = _receiverAddr; numReceivers = _receiverID + 1; _ret = int(_receiverID); return _ret; } /* returns -1 -- adress already exists */ function newSender(address _senderAddr, string _name) returns (int _ret) { if(_senderAddr == 0){ _senderAddr = msg.sender; } if(senders[_senderAddr].creationDate != 0){ return -1; } senders[_senderAddr].name = _name; senders[_senderAddr].creationDate = now; senders[_senderAddr].numTransactions = 0; senders[_senderAddr].transactions = new uint[](0); senders[_senderAddr].numPoints = 0; /* senders[_senderAddr].points = new Addition[](0); */ uint _senderID = senderList.length++; senderList[_senderID] = _senderAddr; numSenders = _senderID + 1; _ret = int(_senderID); return _ret; } /* returns 0 -- OK -1 -- adress not found */ function addPoints(address _addr, bool isReceiver, uint _amount) returns (int error) { uint _pointID; if(_addr == 0){ _addr = msg.sender; } if(isReceiver == false) { if(senders[_addr].creationDate == 0){ return -1; } _pointID = senders[_addr].points.length++; Addition coin = senders[_addr].points[_pointID]; coin.num = _amount; coin.date = now; senders[_addr].numPoints = _pointID + 1; } else { if(receivers[_addr].creationDate == 0){ throw; } _pointID = receivers[_addr].points.length++; Addition pnt = receivers[_addr].points[_pointID]; pnt.num = _amount; pnt.date = now; receivers[_addr].numPoints = _pointID + 1; } return 0; } /* returns -1 -- adress not found (-10 - Array.length) -- if numArray is not the same length */ function getTotDonations(address _addr, bool isReceiver, bool isFullScan) constant returns (int _total) { _total = 0; uint i; if(_addr == 0) { _addr = msg.sender; } if(isReceiver == false) { if(senders[_addr].creationDate == 0) { return -1; } if(senders[_addr].transactions.length != senders[_addr].numTransactions){ return (-10 - int(senders[_addr].transactions.length)); } if (isFullScan == false){ for(i=0; i<senders[_addr].transactions.length; i++) { _total = _total + int(transactions[senders[_addr].transactions[i]].amount); } }else{ for(i=0; i<transactions.length; i++) { if (transactions[i].senderAddr == _addr) { _total = _total + int(transactions[i].amount); } } } } else { if(receivers[_addr].creationDate == 0) { return -1; } if(receivers[_addr].transactions.length != receivers[_addr].numTransactions){ return (-10 - int(receivers[_addr].transactions.length)); } if (isFullScan == false){ for(i=0; i<receivers[_addr].transactions.length; i++) { _total = _total + int(transactions[receivers[_addr].transactions[i]].amount); } }else{ for(i=0; i<transactions.length; i++) { if (transactions[i].receiverAddr == _addr) { _total = _total + int(transactions[i].amount); } } } } return _total; } /* returns -1 -- address not found (-10 - Array.length) -- if numArray is not the same length */ function getTotPoints(address _addr, bool isReceiver) constant returns (int _total) { _total = 0; uint i; if(_addr == 0) { _addr = msg.sender; } if(isReceiver == false) { if(senders[_addr].creationDate == 0) { return -1; } if(senders[_addr].points.length != senders[_addr].numPoints){ return (-10 - int(senders[_addr].points.length)); } for(i=0; i<senders[_addr].points.length; i++) { _total = _total + int(senders[_addr].points[i].num); } } else { if(receivers[_addr].creationDate == 0) { return -1; } if(receivers[_addr].points.length != receivers[_addr].numPoints){ return (-10 - int(receivers[_addr].points.length)); } for(i=0; i<receivers[_addr].points.length; i++) { _total = _total + int(receivers[_addr].points[i].num); } } return _total; } /* -1 -- error in TotPoints -2 -- error in TotDonations -3 -- balance is negative */ function getBalance(address _addr, bool isReceiver, bool isTotalScan) constant returns (int _balance) { _balance = getTotPoints(_addr, isReceiver); if(_balance < 0){ return -1; } int _donat = getTotDonations(_addr, isReceiver, isTotalScan); if(_donat < 0){ return -2; } if (_balance < _donat) { return -3; } _balance = _balance - _donat; return _balance; } /* priorityID = -1 : we just take the first priority by default returns -1 : if sender balance <0 -2 : if receiver balance <0 -3 : if sender balance <amount -4 : if receiver balance <amount */ function makeDonation(address _senderAddr, address _receiverAddr, uint _amount, int _priorityID, bytes _certificate) returns (int _ret) { if(_senderAddr == 0){ _senderAddr = msg.sender; } int _senderBalance = getBalance(_senderAddr, false, false); if (_senderBalance<0){ return -1; } int _receiverBalance = getBalance(_receiverAddr, true, false); if (_receiverBalance<0){ return -2; } if(_senderBalance < int(_amount)) { return -3; } if(_receiverBalance < int(_amount)) { return -4; } if(_priorityID == -1) { _priorityID = int(receivers[_receiverAddr].priorities[0].num); } uint _transID = transactions.length++; Transaction _tr = transactions[_transID]; _tr.senderAddr = _senderAddr; _tr.receiverAddr = _receiverAddr; _tr.amount = _amount; _tr.priorityQty = 0; _tr.priorityID = uint(_priorityID); _tr.transDate = now; if (_certificate.length > 0) { _tr.certificate = _certificate; } numTransactions = _transID + 1; uint _senderTransID = senders[_senderAddr].transactions.length++; senders[_senderAddr].transactions[_senderTransID] = _transID; senders[_senderAddr].numTransactions = _senderTransID + 1; uint _receiverTransID = receivers[_receiverAddr].transactions.length++; receivers[_receiverAddr].transactions[_receiverTransID] = _transID; receivers[_receiverAddr].numTransactions = _receiverTransID + 1; _ret = int(_transID); _ret = int(_transID); return _ret; } }'

etf_Compiled = eth.compile.solidity(etf_source);

etf_Contract = eth.contract(etf_Compiled['<stdin>:Eethiq'].info.abiDefinition);
etf_Contract = eth.contract(etf_Compiled['<stdin>:Nexus_test'].info.abiDefinition);
//pki_Contract = eth.contract(pki_abi);

etf = etf_Contract.new(
    {
      from:web3.eth.accounts[0], 
      /* data:etf_Compiled['<stdin>:Eethiq'].code, */
      data:etf_Compiled['<stdin>:Nexus_test'].code,
      gas: 20000000
    }, function(e, contract){
      if(!e) {
        if(!contract.address) {
          console.log("Contract transaction send: TransactionHash: " + contract.transactionHash + " waiting to be mined...");

        } else {
          console.log("Contract mined! Address: " + contract.address);
          console.log(contract);
        }
      } else {
        console.log(e)
      }
    })

//Contract mined! Address: 0x9ad32a8274e7c2d93d4cfac84c32c7f64e77c786
//Second blockchain! Nexus_test addr: 0x31c3e3abc9d35d12e140482a2ff7533e40e8befd
//Second blockchain! Nexus_test addr2: 0xcdaf78680655afaf824bc6bcd07bad042dafc272
//Second blockchain! Address: 0xcc2281a517096d2d2962507c77233ffa2ebeb412
// TEST DEBUG 0xa340053b89ee859d524556dd6acdb93708fdc957
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




