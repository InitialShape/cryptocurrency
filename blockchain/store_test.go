package blockchain_test

import (
	"bytes"
	"crypto/rand"
	"errors"
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/InitialShape/blockchain/miner"
	"github.com/mr-tron/base58/base58"
	"github.com/stretchr/testify/assert"
	cbor "github.com/whyrusleeping/cbor/go"
	"golang.org/x/crypto/ed25519"
	"testing"
)

const DB = "/tmp/db321"

var (
	store blockchain.Store
	peer  blockchain.Peer
)

func init() {
	store = blockchain.Store{}
	store.Open(DB, &peer)
	peer = blockchain.Peer{"localhost", "1234", store}
}

func TestPutBlock(t *testing.T) {
	genesis, err := store.StoreGenesisBlock(5)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, genesis.Hash, []blockchain.Transaction{}, ch)
	newBlock := <-ch

	err = store.AddBlock(newBlock)
	if err != nil {
		t.Error(err)
	}
}

func TestGetChain(t *testing.T) {
	genesis, err := store.StoreGenesisBlock(5)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.SearchBlock(1, 5, genesis.Hash, nil, ch)
	firstBlock := <-ch

	err = store.AddBlock(firstBlock)
	if err != nil {
		t.Error(err)
	}

	ch = make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, firstBlock.Hash, nil, ch)
	secondBlock := <-ch

	err = store.AddBlock(secondBlock)
	if err != nil {
		t.Error(err)
	}

	blocks, err := store.GetChain()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, []blockchain.Block{genesis, firstBlock, secondBlock}, blocks)
}

func TestEvaluateChains(t *testing.T) {
	genesis, err := store.StoreGenesisBlock(5)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.SearchBlock(1, 5, genesis.Hash, nil, ch)
	firstBlock := <-ch
	ch = make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, firstBlock.Hash, nil, ch)
	secondBlock := <-ch
	ch = make(chan blockchain.Block)
	go miner.SearchBlock(3, 5, secondBlock.Hash, nil, ch)
	thirdBlock := <-ch

	var firstChain = []blockchain.Block{firstBlock, secondBlock, thirdBlock}
	var secondChain = []blockchain.Block{firstBlock}

	var chains = [][]blockchain.Block{
		secondChain,
		firstChain,
	}

	chain := store.EvaluateChains(chains)
	assert.Equal(t, firstChain, chain)
}

func TestGetTransactionWithNothingInBucket(t *testing.T) {
	transaction, err := store.GetTransaction([]byte("transactions"), false)
	if assert.Error(t, err) {
		assert.Equal(t, err, errors.New("EOF"))
	}

	privateKey, _ := base58.Decode("35DxrJipeuCAakHNnnPkBjwxQffYWKM1632kUFv9vKGRNREFSyM6awhyrucxTNbo9h693nPKeWonJ9sFkw6Tou4d")
	publicKey, _ := base58.Decode("6zjRZQyp47BjwArFoLpvzo8SHwwWeW571kJNiqWfSrFT")
	outputs := []blockchain.Output{blockchain.Output{publicKey, 100}}
	inputs := []blockchain.Input{blockchain.Input{[]byte{}, []byte{}, 0}}
	transaction = blockchain.Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		t.Error(err)
	}
	transaction.Hash = hash
	transaction.Sign(privateKey, 0)
	cbor, err := transaction.GetCBOR()
	if err != nil {
		t.Error(err)
	}
	err = store.Put([]byte("transactions"), transaction.Hash, cbor)

	transaction, err = store.GetTransaction(transaction.Hash, false)

}

func TestPutBlockWithUnsignedTransferTransaction(t *testing.T) {
	genesis, err := store.StoreGenesisBlock(5)
	if err != nil {
		t.Error(err)
	}

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	outputs := []blockchain.Output{blockchain.Output{publicKey, 123}}
	inputs := []blockchain.Input{blockchain.Input{[]byte{},
		genesis.Transactions[0].Hash, 0}}
	transaction := blockchain.Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		t.Error(err)
	}
	transaction.Hash = hash

	coinbase, err := blockchain.GenerateCoinbase(publicKey, privateKey, 100)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, genesis.Hash,
		[]blockchain.Transaction{coinbase, transaction}, ch)
	newBlock := <-ch

	err = store.AddBlock(newBlock)
	assert.Error(t, err)
}

func TestPutBlockWithTransferTransaction(t *testing.T) {
	genesis, err := store.StoreGenesisBlock(5)
	if err != nil {
		t.Error(err)
	}

	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	privateKey, _ := base58.Decode("35DxrJipeuCAakHNnnPkBjwxQffYWKM1632kUFv9vKGRNREFSyM6awhyrucxTNbo9h693nPKeWonJ9sFkw6Tou4d")
	outputs := []blockchain.Output{blockchain.Output{publicKey, 123}}
	inputs := []blockchain.Input{blockchain.Input{[]byte{},
		genesis.Transactions[0].Hash, 0}}
	transaction := blockchain.Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		t.Error(err)
	}
	transaction.Hash = hash
	transaction.Sign(privateKey, 0)

	coinbase, err := blockchain.GenerateCoinbase(publicKey, privateKey, 100)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, genesis.Hash,
		[]blockchain.Transaction{coinbase, transaction}, ch)
	newBlock := <-ch

	err = store.AddBlock(newBlock)
	if err != nil {
		t.Error(err)
	}
}

func TestSpendTransactionTwice(t *testing.T) {
	genesis, err := store.StoreGenesisBlock(5)
	if err != nil {
		t.Error(err)
	}

	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	privateKey, _ := base58.Decode("35DxrJipeuCAakHNnnPkBjwxQffYWKM1632kUFv9vKGRNREFSyM6awhyrucxTNbo9h693nPKeWonJ9sFkw6Tou4d")
	outputs := []blockchain.Output{blockchain.Output{publicKey, 123}}
	inputs := []blockchain.Input{blockchain.Input{[]byte{},
		genesis.Transactions[0].Hash, 0}}
	transaction := blockchain.Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		t.Error(err)
	}
	transaction.Hash = hash
	transaction.Sign(privateKey, 0)

	coinbase, err := blockchain.GenerateCoinbase(publicKey, privateKey, 100)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, genesis.Hash,
		[]blockchain.Transaction{coinbase, transaction}, ch)
	newBlock := <-ch
	err = store.AddBlock(newBlock)
	if err != nil {
		t.Error(err)
	}

	publicKey, _, err = ed25519.GenerateKey(rand.Reader)
	privateKey, _ = base58.Decode("35DxrJipeuCAakHNnnPkBjwxQffYWKM1632kUFv9vKGRNREFSyM6awhyrucxTNbo9h693nPKeWonJ9sFkw6Tou4d")
	outputs = []blockchain.Output{blockchain.Output{publicKey, 123}}
	inputs = []blockchain.Input{blockchain.Input{[]byte{},
		genesis.Transactions[0].Hash, 0}}
	transaction = blockchain.Transaction{[]byte{}, inputs, outputs}
	hash, err = transaction.GetHash()
	if err != nil {
		t.Error(err)
	}
	transaction.Hash = hash
	transaction.Sign(privateKey, 0)

	ch = make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, genesis.Hash,
		[]blockchain.Transaction{coinbase, transaction}, ch)
	newBlock = <-ch

	err = store.AddBlock(newBlock)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("Output doesn't exist (anymore?)"), err)
	}
}

func TestPutCoinbaseTwiceInBlock(t *testing.T) {
	genesis, err := store.StoreGenesisBlock(5)
	if err != nil {
		t.Error(err)
	}

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	transaction, err := blockchain.GenerateCoinbase(publicKey, privateKey, 100)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, genesis.Hash,
		[]blockchain.Transaction{transaction, transaction}, ch)
	newBlock := <-ch

	err = store.AddBlock(newBlock)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("Transaction duplicate in block"), err)
	}
}

func TestPutBlockWithTooLowDifficulty(t *testing.T) {
	// TODO: This function fails every once in a while
	genesis, err := store.StoreGenesisBlock(6)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, genesis.Hash, []blockchain.Transaction{}, ch)
	newBlock := <-ch

	err = store.AddBlock(newBlock)
	assert.Error(t, err)
}

func TestPutAndGetData(t *testing.T) {
	expected := []byte("def")
	err := store.Put([]byte("123"), []byte("abc"), expected)
	if err != nil {
		t.Error(err)
	}

	data, err := store.Get([]byte("123"), []byte("abc"))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, data, expected)
}

func TestGetTransactions(t *testing.T) {
	// NOTE: This test is not ideal. It just tests if the new transaction is
	// contained in the result set. Instead, it should create a new database,
	// add multiple transactions and then check for equal.
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}

	outputs := []blockchain.Output{blockchain.Output{publicKey, 10}}
	inputs := []blockchain.Input{blockchain.Input{[]byte{}, []byte{}, 0}}
	transaction := blockchain.Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		t.Error(err)
	}
	transaction.Hash = hash
	store.AddTransaction(transaction)

	storeTransactions, err := store.GetTransactions()
	if err != nil {
		t.Error(err)
	}
	assert.Contains(t, storeTransactions, transaction)
}

func TestAddTransaction(t *testing.T) {
	publicKey, _ := base58.Decode("6zjRZQyp47BjwArFoLpvzo8SHwwWeW571kJNiqWfSrFT")
	outputs := []blockchain.Output{blockchain.Output{publicKey, 10}}
	inputs := []blockchain.Input{blockchain.Input{[]byte{}, []byte{}, 0}}
	transaction := blockchain.Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		t.Error(err)
	}
	transaction.Hash = hash
	store.AddTransaction(transaction)

	data, err := store.Get([]byte("mempool"), transaction.Hash)
	if err != nil {
		t.Error(err)
	}
	var storeTransaction blockchain.Transaction
	dec := cbor.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&storeTransaction)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, transaction, storeTransaction)

}
