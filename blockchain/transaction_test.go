package blockchain

import (
	"crypto/rand"
	"github.com/mr-tron/base58/base58"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
	"testing"
)

const PUBLIC_KEY = "mVHLEtFHLYQE7mwvkhkUp9uKqq5VDCMLvjYtePtMix5"

func TestCoinbaseTransaction(t *testing.T) {
	t.Skip() // signed and hashed transaction
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	transaction, err := GenerateCoinbase(publicKey, privateKey, 100)
	if err != nil {
		t.Error(err)
	}

	outputs := []Output{Output{publicKey, 100}}
	inputs := []Input{Input{[]byte{}, []byte{}, 0}}
	expected := Transaction{[]byte{}, inputs, outputs}

	assert.Equal(t, transaction, expected)
}

func TestTransactionGetBase58Hash(t *testing.T) {
	publicKey, err := base58.Decode(PUBLIC_KEY)
	if err != nil {
		t.Error(err)
	}
	outputs := []Output{Output{publicKey, 100}}
	inputs := []Input{Input{[]byte{}, []byte{}, 0}}
	transaction := Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetBase58Hash()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "4LFP1GFFMRaB3ay6ooisKzzEUbDqnFNpWFCrqKrxhhi1", hash)
}

func TestTransactionSign(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	outputs := []Output{Output{publicKey, 100}}
	inputs := []Input{Input{[]byte{}, []byte{}, 0}}
	transaction := Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		t.Error(err)
	}
	signature := ed25519.Sign(privateKey, hash)
	transaction.Sign(privateKey, 0)

	assert.Equal(t, transaction.Inputs[0].Signature, signature)
}

func TestTransactionVerifySignature(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	publicKey2, privateKey2, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	outputs := []Output{Output{publicKey, 100}}
	inputs := []Input{
		Input{[]byte{}, []byte{}, 0},
		Input{[]byte{}, []byte{}, 0},
	}
	transaction := Transaction{[]byte{}, inputs, outputs}
	transaction.Sign(privateKey, 0)
	transaction.Sign(privateKey2, 1)

	result, err := transaction.Verify(publicKey, 0)
	if err != nil {
		t.Error(err)
	}

	assert.True(t, result)
	result2, err := transaction.Verify(publicKey2, 1)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, result2)
}
