package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/stretchr/testify/assert"
	"github.com/mr-tron/base58/base58"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/InitialShape/blockchain/miner"
)

var (
	server    *httptest.Server
	blocksUrl string
	rootUrl   string
)

func init() {
	server = httptest.NewServer(Handlers())

	blocksUrl = fmt.Sprintf("%s/blocks", server.URL)
	rootUrl = fmt.Sprintf("%s/root", server.URL)
}

func TestPutBlock(t *testing.T) {
	store := blockchain.Store{LEVEL_DB}
	err := store.StoreGenesisBlock(3)
	genesis, err := base58.Decode("6gMgy5V3nyQyue8wWXqo3buiZcXzwR3qNgv5SexeGZLG")

	ch := make(chan blockchain.Block)
	go miner.GenerateBlock(2, 2, genesis, ch)
	newBlock := <-ch

	newBlockJSON, err := json.Marshal(newBlock)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(newBlockJSON))

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
		t.Errorf("Expected status code 201 but got ", res.StatusCode)
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
