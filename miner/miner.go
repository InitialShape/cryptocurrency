package miner

import (
	"encoding/json"
	"fmt"
	"github.com/InitialShape/blockchain/blockchain"
	"log"
	"net/http"
	"os"
	"math/rand"
	"bytes"
	"errors"
)


func GenerateBlock(height int, difficulty int, previousBlock []byte,
				   ch chan<- blockchain.Block) {
	newBlock := blockchain.Block{height, []byte{}, nil, previousBlock,
	difficulty, 0}

	for {
		// Use 256 bits
		newBlock.Nonce = rand.Int31()

		hash, err := newBlock.GetHash()
		if err != nil {
			log.Fatal(err)
		}
		if blockchain.HashMatchesDifficulty(hash, difficulty) {
			newBlock.Hash = hash
			ch <- newBlock
			break
		}
	}
}

func SubmitBlock(block blockchain.Block) {
	fmt.Println(block.Nonce)

	blockJSON, err := json.Marshal(block)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	blocksUrl := fmt.Sprintf("%s/blocks", os.Args[1])
	fmt.Println(blocksUrl)
	req, err := http.NewRequest(http.MethodPut, blocksUrl,
								bytes.NewReader(blockJSON))
	fmt.Println("Sending new found block")
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 201 {
		errors.New(fmt.Sprintf("Expected status code 201 but got %s",
							   res.StatusCode))
	}
}
