package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/InitialShape/blockchain/blockchain"
	"golang.org/x/crypto/ed25519"
	"log"
	"net/http"
)

const server = "http://localhost:8000"

var transactionsUrl string

func init() {
	transactionsUrl = fmt.Sprintf("%s/mempool/transactions", server)
}

func main() {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	outputs := []blockchain.Output{blockchain.Output{publicKey, 10}}
	inputs := []blockchain.Input{blockchain.Input{[]byte{}, "", 0}}
	transaction := blockchain.Transaction{[]byte{}, inputs, outputs}
	hash, err := transaction.GetHash()
	if err != nil {
		log.Fatal(err)
	}
	transaction.Hash = hash
	transaction.Sign(privateKey, 0)

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, transactionsUrl,
		bytes.NewReader(transactionJSON))
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(transaction.GetBase58Hash())

	if res.StatusCode != 201 {
		log.Fatal(errors.New("Transaction wasn't created"))
	}
}
