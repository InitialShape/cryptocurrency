package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/mr-tron/base58/base58"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"log"
	"github.com/InitialShape/blockchain/miner"
)

var (
	server    *httptest.Server
	blocksUrl string
	transactionsUrl string
	rootUrl   string
	store blockchain.Store
)

func init() {
	store = blockchain.Store{}
	err := store.Open(DB)
	if err != nil {
		log.Fatal(err)
	}
	server = httptest.NewServer(Handlers(store))

	blocksUrl = fmt.Sprintf("%s/blocks", server.URL)
	transactionsUrl = fmt.Sprintf("%s/transactions", server.URL)
	rootUrl = fmt.Sprintf("%s/root", server.URL)
}

func TestPutTransaction(t *testing.T) {
	publicKey, _ := base58.Decode("6zjRZQyp47BjwArFoLpvzo8SHwwWeW571kJNiqWfSrFT")
	outputs := []blockchain.Output{blockchain.Output{publicKey, 10}}
	inputs := []blockchain.Input{blockchain.Input{[]byte{}, "", 0}}
	transaction := blockchain.Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		t.Error(err)
	}
	transaction.Hash = hash

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, transactionsUrl,
								bytes.NewReader(transactionJSON))
	if err != nil {
		t.Error(err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 201 {
		t.Errorf("Expected status code 201 but got %d", res.StatusCode)
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

	newBlockJSON, err := json.Marshal(newBlock)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, blocksUrl, bytes.NewReader(newBlockJSON))
	if err != nil {
		t.Error(err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Expected status code 201 but got %d", res.StatusCode)
	}

	req, err = http.NewRequest(http.MethodGet, rootUrl, nil)
	if err != nil {
		t.Error(err)
	}
	res, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(body))
	var root blockchain.Block
	err = json.Unmarshal(body, &root)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, newBlock, root)
}
