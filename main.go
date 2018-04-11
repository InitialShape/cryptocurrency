package main

import (
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/InitialShape/blockchain/web"
	"log"
	"net/http"
)

const DB = "db"

func main() {
	store := blockchain.Store{}
	store.Open(DB)

	_, err := store.StoreGenesisBlock(10)
	if err != nil {
		log.Fatal(err)
	}

	r := web.Handlers(store)
	log.Fatal(http.ListenAndServe(":8000", r))
}
