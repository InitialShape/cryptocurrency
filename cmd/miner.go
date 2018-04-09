package main

import (
	"os"
	"github.com/InitialShape/blockchain/miner"
	"github.com/InitialShape/blockchain/blockchain"
	"strconv"
	"log"
)

func main() {
	ch := make(chan blockchain.Block)
	workers, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < workers; i++ {
		go miner.GenerateBlock(os.Args[1], ch)
	}
	newBlock := <-ch
	miner.SubmitBlock(newBlock)
}
