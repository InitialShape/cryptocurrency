package blockchain

import (
	"bytes"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	cbor "github.com/whyrusleeping/cbor/go"
	"golang.org/x/crypto/ed25519"
	"log"
	"testing"
)

func TestMarshal(t *testing.T) {
	// change this to a static nonce once mining algorithm is implemented
	block := Block{0, []byte{}, nil, []byte{}, 1, 1}
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
	// change this to a static nonce once mining algorithm is implemented
	block := Block{0, []byte{}, nil, []byte{}, 1, 1}
	hash, err := block.GetBase58Hash()
	if err != nil {
		t.Error(err)
	}
	expected := "8peBAxqs9Hsq9TTQHUebgR6iLAUeJcXx4CFw6pJ8aVR9"
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
