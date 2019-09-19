package asimovrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
)

// AsimovError - ethereum error
type AsimovError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err AsimovError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type asimovResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *AsimovError    `json:"error"`
}

type asimovRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// AsimovRPC - Ethereum rpc client
type AsimovRPC struct {
	url    string
	client httpClient
	log    logger
	Debug  bool
}

// New create new rpc client with given url
func New(url string, options ...func(rpc *AsimovRPC)) *AsimovRPC {
	rpc := &AsimovRPC{
		url:    url,
		client: http.DefaultClient,
		log:    log.New(os.Stderr, "", log.LstdFlags),
	}
	for _, option := range options {
		option(rpc)
	}

	return rpc
}

// NewAsimovRPC create new rpc client with given url
func NewAsimovRPC(url string, options ...func(rpc *AsimovRPC)) *AsimovRPC {
	return New(url, options...)
}

func (rpc *AsimovRPC) call(method string, target interface{}, params ...interface{}) error {
	result, err := rpc.Call(method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}

// URL returns client url
func (rpc *AsimovRPC) URL() string {
	return rpc.url
}

// Call returns raw response of method call
func (rpc *AsimovRPC) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := asimovRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := rpc.client.Post(rpc.url, "application/json", bytes.NewBuffer(body))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if rpc.Debug {
		rpc.log.Println(fmt.Sprintf("%s\nRequest: %s\nResponse: %s\n", method, body, data))
	}

	resp := new(asimovResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}

// RawCall returns raw response of method call (Deprecated)
func (rpc *AsimovRPC) RawCall(method string, params ...interface{}) (json.RawMessage, error) {
	return rpc.Call(method, params...)
}

// Web3ClientVersion returns the current client version.
func (rpc *AsimovRPC) Web3ClientVersion() (string, error) {
	var clientVersion string

	err := rpc.call("web3_clientVersion", &clientVersion)
	return clientVersion, err
}

// Web3Sha3 returns Keccak-256 (not the standardized SHA3-256) of the given data.
func (rpc *AsimovRPC) Web3Sha3(data []byte) (string, error) {
	var hash string

	err := rpc.call("web3_sha3", &hash, fmt.Sprintf("0x%x", data))
	return hash, err
}

// NetVersion returns the current network protocol version.
func (rpc *AsimovRPC) NetVersion() (string, error) {
	var version string

	err := rpc.call("net_version", &version)
	return version, err
}

// NetListening returns true if client is actively listening for network connections.
func (rpc *AsimovRPC) NetListening() (bool, error) {
	var listening bool

	err := rpc.call("net_listening", &listening)
	return listening, err
}

// NetPeerCount returns number of peers currently connected to the client.
func (rpc *AsimovRPC) NetPeerCount() (int, error) {
	var response string
	if err := rpc.call("net_peerCount", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthProtocolVersion returns the current ethereum protocol version.
func (rpc *AsimovRPC) AsimovProtocolVersion() (string, error) {
	var protocolVersion string

	err := rpc.call("flow_protocolVersion", &protocolVersion)
	return protocolVersion, err
}

// EthSyncing returns an object with data about the sync status or false.
func (rpc *AsimovRPC) AsimovSyncing() (*Syncing, error) {
	result, err := rpc.RawCall("flow_syncing")
	if err != nil {
		return nil, err
	}
	syncing := new(Syncing)
	if bytes.Equal(result, []byte("false")) {
		return syncing, nil
	}
	err = json.Unmarshal(result, syncing)
	return syncing, err
}

// EthCoinbase returns the client coinbase address
func (rpc *AsimovRPC) AsimovCoinbase() (string, error) {
	var address string

	err := rpc.call("flow_coinbase", &address)
	return address, err
}

// EthMining returns true if client is actively mining new blocks.
func (rpc *AsimovRPC) AsimovMining() (bool, error) {
	var mining bool

	err := rpc.call("flow_mining", &mining)
	return mining, err
}

// EthHashrate returns the number of hashes per second that the node is mining with.
func (rpc *AsimovRPC) AsimovHashrate() (int, error) {
	var response string

	if err := rpc.call("flow_hashrate", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGasPrice returns the current price per gas in wei.
func (rpc *AsimovRPC) AsimovGasPrice() (big.Int, error) {
	var response string
	if err := rpc.call("flow_gasPrice", &response); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// EthAccounts returns a list of addresses owned by client.
func (rpc *AsimovRPC) AsimovAccounts() ([]string, error) {
	accounts := []string{}

	err := rpc.call("flow_accounts", &accounts)
	return accounts, err
}

// EthBlockNumber returns the number of most recent block.
func (rpc *AsimovRPC) AsimovBlockNumber() (int, error) {
	var response string
	if err := rpc.call("flow_blockNumber", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetBalance returns the balance of the account of given address in wei.
func (rpc *AsimovRPC) AsimovGetBalance(address, block string) (big.Int, error) {
	var response string
	if err := rpc.call("flow_getBalance", &response, address, block); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// EthGetStorageAt returns the value from a storage position at a given address.
func (rpc *AsimovRPC) AsimovGetStorageAt(data string, position int, tag string) (string, error) {
	var result string

	err := rpc.call("flow_getStorageAt", &result, data, IntToHex(position), tag)
	return result, err
}

// EthGetTransactionCount returns the number of transactions sent from an address.
func (rpc *AsimovRPC) AsimovGetTransactionCount(address, block string) (int, error) {
	var response string

	if err := rpc.call("flow_getTransactionCount", &response, address, block); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetBlockTransactionCountByHash returns the number of transactions in a block from a block matching the given block hash.
func (rpc *AsimovRPC) AsimovGetBlockTransactionCountByHash(hash string) (int, error) {
	var response string

	if err := rpc.call("flow_getBlockTransactionCountByHash", &response, hash); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetBlockTransactionCountByNumber returns the number of transactions in a block from a block matching the given block
func (rpc *AsimovRPC) AsimovGetBlockTransactionCountByNumber(number int) (int, error) {
	var response string

	if err := rpc.call("flow_getBlockTransactionCountByNumber", &response, IntToHex(number)); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetUncleCountByBlockHash returns the number of uncles in a block from a block matching the given block hash.
func (rpc *AsimovRPC) AsimovGetUncleCountByBlockHash(hash string) (int, error) {
	var response string

	if err := rpc.call("flow_getUncleCountByBlockHash", &response, hash); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetUncleCountByBlockNumber returns the number of uncles in a block from a block matching the given block number.
func (rpc *AsimovRPC) AsimovGetUncleCountByBlockNumber(number int) (int, error) {
	var response string

	if err := rpc.call("flow_getUncleCountByBlockNumber", &response, IntToHex(number)); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// EthGetCode returns code at a given address.
func (rpc *AsimovRPC) AsimovGetCode(address, block string) (string, error) {
	var code string

	err := rpc.call("flow_getCode", &code, address, block)
	return code, err
}

// EthSign signs data with a given address.
// Calculates an Ethereum specific signature with: sign(keccak256("\x19Ethereum Signed Message:\n" + len(message) + message)))
func (rpc *AsimovRPC) AsimovSign(address, data string) (string, error) {
	var signature string

	err := rpc.call("flow_sign", &signature, address, data)
	return signature, err
}

// EthSendTransaction creates new message call transaction or a contract creation, if the data field contains code.
func (rpc *AsimovRPC) AsimovSendTransaction(transaction T) (string, error) {
	var hash string

	err := rpc.call("flow_sendTransaction", &hash, transaction)
	return hash, err
}

// EthSendRawTransaction creates new message call transaction or a contract creation for signed transactions.
func (rpc *AsimovRPC) AsimovSendRawTransaction(data string) (string, error) {
	var hash string

	err := rpc.call("flow_sendRawTransaction", &hash, data)
	return hash, err
}

// EthCall executes a new message call immediately without creating a transaction on the block chain.
func (rpc *AsimovRPC) AsimovCall(transaction T, tag string) (string, error) {
	var data string

	err := rpc.call("flow_call", &data, transaction, tag)
	return data, err
}

// EthEstimateGas makes a call or transaction, which won't be added to the blockchain and returns the used gas, which can be used for estimating the used gas.
func (rpc *AsimovRPC) AsimovEstimateGas(transaction T) (int, error) {
	var response string

	err := rpc.call("flow_estimateGas", &response, transaction)
	if err != nil {
		return 0, err
	}

	return ParseInt(response)
}

func (rpc *AsimovRPC) getBlock(method string, withTransactions bool, params ...interface{}) (*Block, error) {
	result, err := rpc.RawCall(method, params...)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(result, []byte("null")) {
		return nil, nil
	}

	var response proxyBlock
	if withTransactions {
		response = new(proxyBlockWithTransactions)
	} else {
		response = new(proxyBlockWithoutTransactions)
	}

	err = json.Unmarshal(result, response)
	if err != nil {
		return nil, err
	}

	block := response.toBlock()
	return &block, nil
}

// EthGetBlockByHash returns information about a block by hash.
func (rpc *AsimovRPC) AsimovGetBlockByHash(hash string, withTransactions bool) (*Block, error) {
	return rpc.getBlock("flow_getBlockByHash", withTransactions, hash, withTransactions)
}

// EthGetBlockByNumber returns information about a block by block number.
func (rpc *AsimovRPC) AsimovGetBlockByNumber(number int, withTransactions bool) (*Block, error) {
	return rpc.getBlock("flow_getBlockByNumber", withTransactions, IntToHex(number), withTransactions)
}

func (rpc *AsimovRPC) getTransaction(method string, params ...interface{}) (*Transaction, error) {
	transaction := new(Transaction)

	err := rpc.call(method, transaction, params...)
	return transaction, err
}

// EthGetTransactionByHash returns the information about a transaction requested by transaction hash.
func (rpc *AsimovRPC) AsimovGetTransactionByHash(hash string) (*Transaction, error) {
	return rpc.getTransaction("flow_getTransactionByHash", hash)
}

// EthGetTransactionByBlockHashAndIndex returns information about a transaction by block hash and transaction index position.
func (rpc *AsimovRPC) AsimovGetTransactionByBlockHashAndIndex(blockHash string, transactionIndex int) (*Transaction, error) {
	return rpc.getTransaction("flow_getTransactionByBlockHashAndIndex", blockHash, IntToHex(transactionIndex))
}

// EthGetTransactionByBlockNumberAndIndex returns information about a transaction by block number and transaction index position.
func (rpc *AsimovRPC) AsimovGetTransactionByBlockNumberAndIndex(blockNumber, transactionIndex int) (*Transaction, error) {
	return rpc.getTransaction("flow_getTransactionByBlockNumberAndIndex", IntToHex(blockNumber), IntToHex(transactionIndex))
}

// EthGetTransactionReceipt returns the receipt of a transaction by transaction hash.
// Note That the receipt is not available for pending transactions.
func (rpc *AsimovRPC) AsimovGetTransactionReceipt(hash string) (*TransactionReceipt, error) {
	transactionReceipt := new(TransactionReceipt)

	err := rpc.call("flow_getTransactionReceipt", transactionReceipt, hash)
	if err != nil {
		return nil, err
	}

	return transactionReceipt, nil
}

// EthGetCompilers returns a list of available compilers in the client.
func (rpc *AsimovRPC) AsimovGetCompilers() ([]string, error) {
	compilers := []string{}

	err := rpc.call("flow_getCompilers", &compilers)
	return compilers, err
}

// EthNewFilter creates a new filter object.
func (rpc *AsimovRPC) AsimovNewFilter(params FilterParams) (string, error) {
	var filterID string
	err := rpc.call("flow_newFilter", &filterID, params)
	return filterID, err
}

// EthNewBlockFilter creates a filter in the node, to notify when a new block arrives.
// To check if the state has changed, call EthGetFilterChanges.
func (rpc *AsimovRPC) AsimovNewBlockFilter() (string, error) {
	var filterID string
	err := rpc.call("flow_newBlockFilter", &filterID)
	return filterID, err
}

// EthNewPendingTransactionFilter creates a filter in the node, to notify when new pending transactions arrive.
// To check if the state has changed, call EthGetFilterChanges.
func (rpc *AsimovRPC) AsimovNewPendingTransactionFilter() (string, error) {
	var filterID string
	err := rpc.call("flow_newPendingTransactionFilter", &filterID)
	return filterID, err
}

// EthUninstallFilter uninstalls a filter with given id.
func (rpc *AsimovRPC) AsimovUninstallFilter(filterID string) (bool, error) {
	var res bool
	err := rpc.call("flow_uninstallFilter", &res, filterID)
	return res, err
}

// EthGetFilterChanges polling method for a filter, which returns an array of logs which occurred since last poll.
func (rpc *AsimovRPC) AsimovGetFilterChanges(filterID string) ([]Log, error) {
	var logs = []Log{}
	err := rpc.call("flow_getFilterChanges", &logs, filterID)
	return logs, err
}

// EthGetFilterLogs returns an array of all logs matching filter with given id.
func (rpc *AsimovRPC) AsimovGetFilterLogs(filterID string) ([]Log, error) {
	var logs = []Log{}
	err := rpc.call("flow_getFilterLogs", &logs, filterID)
	return logs, err
}

// EthGetLogs returns an array of all logs matching a given filter object.
func (rpc *AsimovRPC) AsimovGetLogs(params FilterParams) ([]Log, error) {
	var logs = []Log{}
	err := rpc.call("flow_getLogs", &logs, params)
	return logs, err
}

// Asim1 returns 1 ethereum value (10^18 xin)
func (rpc *AsimovRPC) Asim1() *big.Int {
	return Asim1()
}

// Asim1 returns 1 ethereum value (10^18 xin)
func Asim1() *big.Int {
	return big.NewInt(1000000000000000000)
}
