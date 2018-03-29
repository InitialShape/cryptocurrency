package main

import (
	"os"
	"github.com/davecgh/go-spew/spew"
	"log"
	"io/ioutil"
	"net/http"
	"fmt"
	"github.com/InitialShape/blockchain/blockchain"
	"encoding/json"
)

func main() {
	rootUrl := fmt.Sprintf("%s/root", os.Args[1])
	res, err := http.Get(rootUrl)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var block blockchain.Block
	err = json.Unmarshal(body, &block)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump("Found block: ", block)

	difficulty := block.Difficulty
	spew.Dump("Difficulty is ", difficulty)
	previousBlock := block.Hash

	nonce := 0
	for {
		newBlock := blockchain.Block{block.Height+1, []byte{}, []blockchain.Transaction{},
						  previousBlock, difficulty, nonce}
		hash, err := newBlock.GetHash()
		if err != nil {
			log.Fatal(err)
		}
		if blockchain.HashMatchesDifficulty(hash, difficulty) {
			for _, n := range(hash) {
				fmt.Printf("%b", n)
			}
			break
		} else {
			fmt.Println("Need to work more")
			nonce += 1
		}
	}
}
