package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
	"encoding/json"
	"github.com/mr-tron/base58/base58"
	"log"
	"bytes"
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/davecgh/go-spew/spew"
	cbor "github.com/whyrusleeping/cbor/go"
)

const LEVEL_DB = "db"

func GetBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	db, err := leveldb.OpenFile(LEVEL_DB, nil)
	defer db.Close()
	if err != nil {
		log.Fatal("database not found", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - leveldb couldn't be opened"))
	}

	hash, err := base58.Decode(params["hash"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Couldn't decode base58"))
	}

	data, err := db.Get(hash, nil)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("hash wasn't found"))
	}

	var block blockchain.Block
	dec := cbor.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&block)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - couldn't unmarshal binary"))
	}

	json.NewEncoder(w).Encode(block)
}


func PutBlock(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var block blockchain.Block
	err := dec.Decode(&block)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - couldn't decode block"))
	}
	spew.Dump(block)
}

func Handlers() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/blocks/{hash}", GetBlock).Methods("GET")
	r.HandleFunc("/blocks", PutBlock).Methods("PUT")
	return r
}
