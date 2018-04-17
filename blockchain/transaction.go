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
	TransactionHash []byte `json:"transaction_hash"`
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

func GenerateCoinbase(publicKey ed25519.PublicKey,
	privateKey ed25519.PrivateKey, amount int) (Transaction,
	error) {
	outputs := []Output{Output{publicKey, amount}}
	inputs := []Input{Input{[]byte{}, []byte{}, 0}}
	transaction := Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		return transaction, err
	}
	transaction.Hash = hash
	transaction.Sign(privateKey, 0)
	return transaction, err
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
	if err != nil {
		log.Fatal("Error signing transaction: ", err)
	}
	signature := ed25519.Sign(privateKey, hash)
	t.Inputs[index].Signature = signature
	return err
}

func (t *Transaction) Verify(publicKey ed25519.PublicKey, index int) (bool, error) {
	hash, err := t.GetHash()
	if err != nil {
		log.Fatal("Error validating transaction: ", err)
		return false, err
	}
	return ed25519.Verify(publicKey, hash, t.Inputs[index].Signature), err
}
