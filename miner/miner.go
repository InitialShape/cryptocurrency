package main

import (
	"encoding/json"
	"encoding/binary"
	"crypto/rand"
	"fmt"
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	for i := 0; i < 4; i++ {
		go generateBlock(block, difficulty, previousBlock, ch)
	}
	newBlock := <-ch
	spew.Dump(newBlock)
}

func generateBlock(block blockchain.Block, difficulty int, previousBlock []byte, ch chan<- blockchain.Block) {
	for {
		var nonce int32
		binary.Read(rand.Reader, binary.LittleEndian, &nonce)

		newBlock := blockchain.Block{block.Height + 1, []byte{}, []blockchain.Transaction{},
			previousBlock, difficulty, nonce}
		hash, err := newBlock.GetHash()
		if err != nil {
			log.Fatal(err)
		}
		if blockchain.HashMatchesDifficulty(hash, difficulty) {
			for _, n := range hash {
				fmt.Printf("%b", n)
			}
			newBlock.Hash = hash
			ch <- newBlock
			break
		} else {
			nonce += 1
		}
	}
}

func SubmitBlock(block blockchain.Block) {

}
