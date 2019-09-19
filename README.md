# AsimovRPC

- [x] web3_clientVersion
- [x] web3_sha3
- [x] net_version
- [x] net_peerCount
- [x] net_listening
- [x] flow_protocolVersion
- [x] flow_syncing
- [x] flow_coinbase
- [x] flow_mining
- [x] flow_hashrate
- [x] flow_gasPrice
- [x] flow_accounts
- [x] flow_blockNumber
- [x] flow_getBalance
- [x] flow_getStorageAt
- [x] flow_getTransactionCount
- [x] flow_getBlockTransactionCountByHash
- [x] flow_getBlockTransactionCountByNumber
- [x] flow_getUncleCountByBlockHash
- [x] flow_getUncleCountByBlockNumber
- [x] flow_getCode
- [x] flow_sign
- [x] flow_sendTransaction
- [x] flow_sendRawTransaction
- [x] flow_call
- [x] flow_estimateGas
- [x] flow_getBlockByHash
- [x] flow_getBlockByNumber
- [x] flow_getTransactionByHash
- [x] flow_getTransactionByBlockHashAndIndex
- [x] flow_getTransactionByBlockNumberAndIndex
- [x] flow_getTransactionReceipt
- [x] flow_getCompilers (DEPRECATED)
- [x] flow_newFilter
- [x] flow_newBlockFilter
- [x] flow_newPendingTransactionFilter
- [x] flow_uninstallFilter
- [x] flow_getFilterChanges
- [x] flow_getFilterLogs
- [x] flow_getLogs

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/mistdex/mist-asimov-rpc"
)

func main() {
    client := asimovrpc.New("http://127.0.0.1:8545")

    version, err := client.Web3ClientVersion()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(version)

    // Send 1 Asim
    txid, err := client.AsimovSendTransaction(ethrpc.T{
        From:  "0x6247cf0412c6462da2a51d05139e2a3c6c630f0a",
        To:    "0xcfa202c4268749fbb5136f2b68f7402984ed444b",
        Value: asimovrpc.Asim1(),
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(txid)
}
```