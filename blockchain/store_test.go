package blockchain_test

import (
	"testing"
	"github.com/InitialShape/blockchain/miner"
	"github.com/stretchr/testify/assert"
	"github.com/InitialShape/blockchain/blockchain"
)

const LEVEL_DB = "db"

func TestPutBlock(t *testing.T) {
	store := blockchain.Store{LEVEL_DB}
	genesis, err := store.StoreGenesisBlock(5)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.GenerateBlock(2, 5, genesis, ch)
	newBlock := <-ch

	err = store.AddBlock(newBlock)
	if err != nil {
		t.Error(err)
	}
}

func TestPutBlockWithTooLowDifficulty(t *testing.T) {
	t.Skip()
	store := blockchain.Store{LEVEL_DB}
	genesis, err := store.StoreGenesisBlock(6)
	if err != nil {
		t.Error(err)
	}

	ch := make(chan blockchain.Block)
	go miner.GenerateBlock(2, 5, genesis, ch)
	newBlock := <-ch

	err = store.AddBlock(newBlock)
	assert.Error(t, err)
}
