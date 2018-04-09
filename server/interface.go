package server

import (
	"bytes"
	"encoding/json"
	"github.com/InitialShape/blockchain/blockchain"
	"github.com/gorilla/mux"
	"github.com/mr-tron/base58/base58"
	cbor "github.com/whyrusleeping/cbor/go"
	"log"
	"net/http"
)

const DB = "/tmp/hello"

func GetBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hash, err := base58.Decode(params["hash"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Couldn't decode base58"))
	}

	store := blockchain.Store{DB}
	data, err := store.Get(hash)
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

func GetRootBlock(w http.ResponseWriter, r *http.Request) {
	store := blockchain.Store{DB}
	hash, err := store.Get([]byte("root"))
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("root wasn't found"))
	}

	blockData, err := store.Get(hash)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("root data wasn't found"))
	}

	var block blockchain.Block
	dec := cbor.NewDecoder(bytes.NewReader(blockData))
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
	store := blockchain.Store{DB}
	err = store.AddBlock(block)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - block's difficulty doesn't match"))
	}

	w.WriteHeader(http.StatusCreated)
}

func Handlers() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/blocks/{hash}", GetBlock).Methods("GET")
	r.HandleFunc("/blocks", PutBlock).Methods("PUT")
	r.HandleFunc("/root", GetRootBlock).Methods("GET")
	return r
}
