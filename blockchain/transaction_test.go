package blockchain

import (
	"testing"
	"golang.org/x/crypto/ed25519"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
)

func TestGenesis(t *testing.T) {
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	transaction := Genesis(publicKey, 100)

	outputs := []Output{Output{publicKey, 100}}
	inputs := []Input{Input{[]byte{}, "", 0}}
	expected := Transaction{"", inputs, outputs}

	assert.Equal(t, transaction, expected)
}

func TestTransactionGetBase58Hash(t *testing.T) {
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	outputs := []Output{Output{publicKey, 100}}
	inputs := []Input{Input{[]byte{}, "", 0}}
	transaction := Transaction{"", inputs, outputs}
	hash, err := transaction.GetBase58Hash()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, hash, "JcDex1mjA9nUawCutP2oaTrzLby9L6ezZscLL1xdmY5")
}

func TestTransactionSign(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	outputs := []Output{Output{publicKey, 100}}
	inputs := []Input{Input{[]byte{}, "", 0}}
	transaction := Transaction{"", inputs, outputs}
	hash, err := transaction.GetHash()
	signature := ed25519.Sign(privateKey, hash)
	transaction.Sign(privateKey, 0)

	assert.Equal(t, transaction.Inputs[0].Signature, signature)
}
