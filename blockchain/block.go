package blockchain

import (
	"bytes"
	"github.com/2tvenom/cbor"
	"log"
	"crypto/sha256"
	"github.com/mr-tron/base58/base58"
)

type Block struct {
	Height        int
	Hash          string
	Transactions  []Transaction
	PreviousBlock string
}

func (b *Block) GetCBOR() (bytes.Buffer, error) {
	var buffer bytes.Buffer
	encoder := cbor.NewEncoder(&buffer)
	ok, err := encoder.Marshal(b)

	if !ok {
		log.Fatal("Error decoding %s", err)
		return bytes.Buffer{}, err
	}

	return buffer, err
}

func (b *Block) GetHash() (string, error) {
	block, err := b.GetCBOR()
	if err != nil {
		return "", err
	}
	hasher := sha256.New()
	hasher.Write(block.Bytes())
	return base58.Encode(hasher.Sum(nil)), err
}
