package server

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	block, err := blockchain.GenerateGenesisBlock(publicKey, privateKey)
	if err != nil {
		t.Error(err)
	}
	b, err := json.Marshal(block)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, blocksUrl, bytes.NewReader(b))
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

	var root blockchain.Block
	err = json.Unmarshal(body, &root)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, block, root)
}
