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
	expected := "DzqgSYkaavhwvZvaoRVbRnkeiK5FbTohPmZ9WCazuLmc"
	assert.Equal(t, expected, hash)
}

func TestUnmarshal(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	block, err := GenerateGenesisBlock(publicKey, privateKey, 3)
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

func TestHashMatchesDifficulty(t *testing.T) {
	hash := []byte{0x1F, 0x00} // 00011111 00000000
	assert.True(t, HashMatchesDifficulty(hash, 3))

	hash = []byte{0x0F, 0x00} // 00001111 00000000
	assert.True(t, HashMatchesDifficulty(hash, 4))

	hash = []byte{0x2F, 0x00} // 00101111 00000000
	assert.False(t, HashMatchesDifficulty(hash, 4))

	hash = []byte{0x00, 0x7F} // 00000000 01111111
	assert.True(t, HashMatchesDifficulty(hash, 9))
}
