package blockchain

import (
	"golang.org/x/crypto/ed25519"
	"github.com/2tvenom/cbor"
	"bytes"
	"log"
	"crypto/sha256"
	"github.com/mr-tron/base58/base58"
)

type Input struct {
	Signature []byte
	TransactionHash string
	OutputID int
}

type Output struct {
	PublicKey ed25519.PublicKey
	Amount int
}

type Transaction struct {
	Hash string
	Inputs []Input
	Outputs []Output
}

func Genesis(publicKey ed25519.PublicKey, amount int)  (Transaction) {
	outputs := []Output{Output{publicKey, amount}}
	inputs := []Input{Input{[]byte{}, "", 0}}
	transaction := Transaction{"", inputs, outputs}
	return transaction
}

func (t *Transaction) GetCBOR() (bytes.Buffer, error) {
	var buffer bytes.Buffer
	encoder := cbor.NewEncoder(&buffer)
	ok, err := encoder.Marshal(t)

	if !ok {
		log.Fatal("Error decoding %s", err)
		return bytes.Buffer{}, err
	}

	return buffer, err
}

func (t *Transaction) GetHash() ([]byte, error) {
	transaction, err := t.GetCBOR()
	if err != nil {
		return []byte{}, err
	}
	hasher := sha256.New()
	hasher.Write(transaction.Bytes())
	return hasher.Sum(nil), err
}

func (t *Transaction) GetBase58Hash() (string, error) {
	hash, err := t.GetHash()
	if err != nil {
		return "", err
	}
	return base58.Encode(hash), err
}

func (t *Transaction) Sign(privateKey ed25519.PrivateKey, index int) (error) {
	hash, err := t.GetHash()
	signature := ed25519.Sign(privateKey, hash)
	t.Inputs[index].Signature = signature
	return err
}
