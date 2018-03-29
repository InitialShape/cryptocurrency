package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	"github.com/InitialShape/blockchain/miner"
	"github.com/InitialShape/blockchain/blockchain"
	"strconv"
	"log"
	"github.com/davecgh/go-spew/spew"
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

	difficulty := block.Difficulty
	spew.Dump("Difficulty is ", difficulty)
	previousBlock := block.Hash

	ch := make(chan blockchain.Block)
	workers, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < workers; i++ {
		go miner.GenerateBlock(block.Height + 1, difficulty, previousBlock, ch)
	}
	newBlock := <-ch
	miner.SubmitBlock(newBlock)
}
