package blockchain

import (
	"bytes"
	"crypto/sha256"
	"github.com/mr-tron/base58/base58"
	cbor "github.com/whyrusleeping/cbor/go"
	"golang.org/x/crypto/ed25519"
	"log"
)

type Input struct {
	Signature       []byte `json:"signature"`
	TransactionHash string `json:"transaction_hash"`
	OutputID        int    `json:"output_id"`
}

type Output struct {
	PublicKey ed25519.PublicKey `json:"public_key"`
	Amount    int               `json:"amount"`
}

type Transaction struct {
	Hash    []byte   `json:"hash"`
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

func GenerateCoinbase(publicKey ed25519.PublicKey, amount int) Transaction {
	outputs := []Output{Output{publicKey, amount}}
	inputs := []Input{Input{[]byte{}, "", 0}}
	transaction := Transaction{[]byte{}, inputs, outputs}
	return transaction
}

func (t *Transaction) GetCBOR() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	enc := cbor.NewEncoder(buf)
	err := enc.Encode(t)

	if err != nil {
		log.Fatal("Error decoding ", err)
		return new(bytes.Buffer), err
	}

	return buf, err
}

func (t *Transaction) GetHash() ([]byte, error) {
	hash := t.Hash
	t.Hash = []byte{}
	transaction, err := t.GetCBOR()
	if err != nil {
		return []byte{}, err
	}
	hasher := sha256.New()
	hasher.Write(transaction.Bytes())
	t.Hash = hash
	return hasher.Sum(nil), err
}

func (t *Transaction) GetBase58Hash() (string, error) {
	hash, err := t.GetHash()
	if err != nil {
		return "", err
	}
	return base58.Encode(hash), err
}

func (t *Transaction) Sign(privateKey ed25519.PrivateKey, index int) error {
	hash, err := t.GetHash()
	signature := ed25519.Sign(privateKey, hash)
	t.Inputs[index].Signature = signature
	return err
}
