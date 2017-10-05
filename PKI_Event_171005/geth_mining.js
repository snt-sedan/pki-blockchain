// geth js script.js
// --exec 'loadScript("/home/user/app/pki/geth_mining.js")'
// geth --identity "NexusTest" --rpc --rpcport "8084" --rpccorsdomain "*" --datadir "/home/user/ethpriv/1" --port "30303" --nodiscover --rpcapi "db,eth,net,web3" --networkid 3578 --nat "any" console
// loadScript("/home/alex/wrk/Dropbox/WORK/RD/codes/BlockChain/Private_Ethereum/geth_mining.js")
// loadScript("/home/user/app/pki/geth_mining.js")


var mining_threads = 1

function checkWork() {
    if (eth.getBlock("pending").transactions.length > 0) {
        if (eth.mining) return;
        console.log("== Pending transactions! Mining...");
        miner.start(mining_threads);
    } else {
        miner.stop();  // This param means nothing - CHANGED TO NO PARAMS!!
        console.log("== No transactions! Mining stopped.");
    }
}

eth.filter("latest", function(err, block) { checkWork(); });
eth.filter("pending", function(err, block) { checkWork(); });

checkWork();

