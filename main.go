package main

import (
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/InitialShape/blockchain/web"
	"fmt"
	"github.com/InitialShape/blockchain/p2p"
	"log"
	"os"
	"net/http"
)

func main() {
	store := blockchain.Store{}
	store.Open(os.Args[1])

	_, err := store.StoreGenesisBlock(10)
	if err != nil {
		log.Fatal(err)
	}

	p := p2p.Peer{"localhost", os.Args[2], store}
	go p.Start()

	r := web.Handlers(store)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Args[3]), r))
}
