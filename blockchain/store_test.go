package blockchain_test

import (
	"testing"
	"github.com/InitialShape/blockchain/miner"
	"github.com/stretchr/testify/assert"
	"github.com/InitialShape/blockchain/blockchain"
	"log"
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
	go miner.SearchBlock(2, 5, genesis, ch)
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
	go miner.SearchBlock(2, 5, genesis, ch)
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
