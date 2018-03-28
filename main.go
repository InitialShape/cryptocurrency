package main

import (
	"github.com/InitialShape/blockchain/server"
	"github.com/InitialShape/blockchain/blockchain"
	"golang.org/x/crypto/ed25519"
	"github.com/syndtr/goleveldb/leveldb"
	"crypto/rand"
	"log"
	"fmt"
	"net/http"
)

const LEVEL_DB = "db"

func main() {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	block, err := blockchain.GenerateGenesisBlock(publicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := leveldb.OpenFile(LEVEL_DB, nil)
	if err != nil {
		log.Fatal(err)
	}
	cbor, err := block.GetCBOR()
	if err != nil {
		log.Fatal(err)
	}

	db.Put(block.Hash, cbor.Bytes(), nil)
	db.Put([]byte("root"), block.Hash, nil)
	db.Close()

	base58Hash, err := block.GetBase58Hash()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(base58Hash)

	r := server.Handlers()
	log.Fatal(http.ListenAndServe(":8000", r))
}
