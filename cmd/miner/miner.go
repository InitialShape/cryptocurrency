package main

import (
	"github.com/InitialShape/cryptocurrency/blockchain"
	"github.com/InitialShape/cryptocurrency/miner"
	"log"
	"os"
	"strconv"
)

func main() {
	for {
		mine()
	}
}

func mine() {
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
