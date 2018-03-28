package blockchain

import (
	"testing"
	"bytes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
	"log"
	"crypto/rand"
	cbor "github.com/whyrusleeping/cbor/go"
)

func TestMarshal(t *testing.T) {
	block := Block{0, []byte{}, nil, "", 1}
	marshalledBlock, err := block.GetCBOR()
	if err != nil {
		t.Error(err)
	}

	buf := new(bytes.Buffer)
	enc := cbor.NewEncoder(buf)
	err = enc.Encode(block)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, marshalledBlock, buf)
}

func TestBlockGetBase58Hash(t *testing.T) {
	block := Block{0, []byte{}, nil, "", 1}
	hash, err := block.GetBase58Hash()
	if err != nil {
		t.Error(err)
	}
	expected := "wdnX1DXTvrfkPuGHed3B4m8eTYeP7E8evegqkzwtn9T"
	assert.Equal(t, expected, hash)
}

func TestUnmarshal(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	block, err := GenerateGenesisBlock(publicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	enc := cbor.NewEncoder(buf)
	err = enc.Encode(block)
	if err != nil {
		log.Fatal("Error decoding ", err)
	}

	var newBlock Block
	dec := cbor.NewDecoder(bytes.NewReader(buf.Bytes()))
	err = dec.Decode(&newBlock)
	if err != nil {
		log.Fatal("Error decoding ", err)
	}

	assert.Equal(t, block, newBlock)
}
