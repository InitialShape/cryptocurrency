package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/mr-tron/base58/base58"
	cbor "github.com/whyrusleeping/cbor/go"
	"golang.org/x/crypto/ed25519"
	"log"
	"strings"
)

type Block struct {
	Height        int           `json:"height"`
	Hash          []byte        `json:"hash"`
	Transactions  []Transaction `json:"transactions"`
	PreviousBlock []byte        `json:"previous_block"`
	Difficulty    int           `json:"difficulty"`
	Nonce         int32         `json:"nonce"`
}

const COINBASE_AMOUNT = 25

func GenerateGenesisBlock(publicKey ed25519.PublicKey,
	privateKey ed25519.PrivateKey,
	difficulty int) (Block, error) {
	coinbase, err := GenerateCoinbase(publicKey, privateKey, COINBASE_AMOUNT)
	if err != nil {
		return Block{}, err
	}

	// change this to a static nonce once mining algorithm is implemented
	block := Block{0, []byte{}, []Transaction{coinbase}, []byte{}, difficulty,
		1}
	hash, err := block.GetHash()
	if err != nil {
		return Block{}, err
	}

	block.Hash = hash
	return block, err
}

func (b *Block) GetCBOR() (*bytes.Buffer, error) {
	// change this to type []byte
	buf := new(bytes.Buffer)
	enc := cbor.NewEncoder(buf)
	err := enc.Encode(b)

	if err != nil {
		log.Fatal("Error decoding ", err)
		return new(bytes.Buffer), err
	}

	return buf, err
}

func (b *Block) GetHash() ([]byte, error) {
	hash := b.Hash
	b.Hash = []byte{}
	block, err := b.GetCBOR()
	if err != nil {
		log.Fatal("Error encoding to CBOR", err)
		return []byte{}, err
	}
	hasher := sha256.New()
	hasher.Write(block.Bytes())
	b.Hash = hash
	return hasher.Sum(nil), err
}

func (b *Block) GetBase58Hash() (string, error) {
	hash, err := b.GetHash()
	if err != nil {
		return "", err
	}
	return base58.Encode(hash), err
}

func HashMatchesDifficulty(hash []byte, difficulty int) bool {
	var hashBinary bytes.Buffer
	var prefix bytes.Buffer
	for _, n := range hash {
		hashBinary.WriteString(fmt.Sprintf("%08b", n))
	}

	for i := 0; i < difficulty; i++ {
		prefix.WriteString("0")
	}

	prefixString := prefix.String()
	hashString := hashBinary.String()

	return strings.HasPrefix(hashString, prefixString)
}
