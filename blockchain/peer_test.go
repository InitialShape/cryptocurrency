package blockchain_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/InitialShape/cryptocurrency/blockchain"
	"net"
	"encoding/json"
	"io/ioutil"
)

func TestPingPong(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		t.Error(err)
	}
	msg := "PING localhost:12345"

	conn.Write([]byte(msg))
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "PONG", string(resp))
}

func TestPingPongWithoutPeer(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		t.Error(err)
	}
	msg := "PING"

	conn.Write([]byte(msg))
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "PONG", string(resp))
}

func TestGettingPeers(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		t.Error(err)
	}
	msg := "PEERS"

	conn.Write([]byte(msg))
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "localhost:12345", string(resp))
}

func TestGetChainFromPeer(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		t.Error(err)
	}
	msg := "CHAIN"

	conn.Write([]byte(msg))
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Error(err)
	}
	var chain []blockchain.Block
	err = json.Unmarshal(resp, &chain)
	if err != nil {
		t.Error(err)
	}

	expected, err := store.GetChain()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expected, chain)
}
