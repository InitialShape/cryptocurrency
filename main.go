package main

import (
	"github.com/InitialShape/blockchain/server"
	"github.com/InitialShape/blockchain/blockchain"
	"log"
	"net/http"
)

const LEVEL_DB = "db"

func main() {
	store := blockchain.Store{LEVEL_DB}
	_, err := store.StoreGenesisBlock(10)
	if err != nil {
		log.Fatal(err)
	}

	r := server.Handlers()
	log.Fatal(http.ListenAndServe(":8000", r))
}
