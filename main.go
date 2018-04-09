package main

import (
	"github.com/InitialShape/blockchain/server"
	"github.com/InitialShape/blockchain/blockchain"
	"log"
	"net/http"
)

const DB = "db"

func main() {
	store := blockchain.Store{DB}
	_, err := store.StoreGenesisBlock(10)
	if err != nil {
		log.Fatal(err)
	}

	r := server.Handlers()
	log.Fatal(http.ListenAndServe(":8000", r))
}
