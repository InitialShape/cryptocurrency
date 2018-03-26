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

func (b *Block) GetHash() ([]byte, error) {
	block, err := b.GetCBOR()
	if err != nil {
		log.Fatal("Error encoding to CBOR %s", err)
		return []byte{}, err
	}
	hasher := sha256.New()
	hasher.Write(block.Bytes())
	return hasher.Sum(nil), err
}

func (b *Block) GetBase58Hash() (string, error) {
	hash, err := b.GetHash()
	if err != nil {
		return "", err
	}
	return base58.Encode(hash), err
}
