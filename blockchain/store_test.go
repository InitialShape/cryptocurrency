package blockchain_test

import (
	"bytes"
	"crypto/rand"
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/InitialShape/blockchain/miner"
	"github.com/mr-tron/base58/base58"
	"github.com/stretchr/testify/assert"
	cbor "github.com/whyrusleeping/cbor/go"
	"golang.org/x/crypto/ed25519"
	"log"
	"testing"
)

const DB = "/tmp/db"

var store blockchain.Store

func init() {
	store = blockchain.Store{}
	err := store.Open(DB)
	if err != nil {
		log.Fatal(err)
	}
}

func TestPutBlock(t *testing.T) {
	genesis, err := store.StoreGenesisBlock(5)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, genesis, []blockchain.Transaction{}, ch)
	newBlock := <-ch

	err = store.AddBlock(newBlock)
	if err != nil {
		t.Error(err)
	}
}

func TestPutBlockWithTooLowDifficulty(t *testing.T) {
	t.Skip()
	genesis, err := store.StoreGenesisBlock(6)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.SearchBlock(2, 5, genesis, []blockchain.Transaction{}, ch)
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
	inputs := []blockchain.Input{blockchain.Input{[]byte{}, "", 0}}
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
	inputs := []blockchain.Input{blockchain.Input{[]byte{}, "", 0}}
	transaction := blockchain.Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		t.Error(err)
	}
	transaction.Hash = hash
	store.AddTransaction(transaction)

	data, err := store.Get([]byte("transactions"), transaction.Hash)
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
