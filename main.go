package main

import (
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/InitialShape/blockchain/web"
	"fmt"
	"log"
	"os"
	"net/http"
)

func main() {
	var store blockchain.Store
	var peer blockchain.Peer

	store = blockchain.Store{}
	store.Open(os.Args[1], &peer)
	peer = blockchain.Peer{"localhost", os.Args[2], store}
	go peer.Start()

	_, err := store.StoreGenesisBlock(20)
	if err != nil {
		log.Fatal(err)
	}

	r := web.Handlers(store)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Args[3]), r))
}
