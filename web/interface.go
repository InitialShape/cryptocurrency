package web

import (
	"bytes"
	"encoding/json"
	"github.com/InitialShape/cryptocurrency/blockchain"
	"github.com/gorilla/mux"
	"github.com/mr-tron/base58/base58"
	cbor "github.com/whyrusleeping/cbor/go"
	"log"
	"net/http"
)

var Store blockchain.Store

func Handlers(store blockchain.Store) *mux.Router {
	Store = store
	r := mux.NewRouter()
	r.HandleFunc("/blocks/{hash}", GetBlock).Methods("GET")
	r.HandleFunc("/blocks", PutBlock).Methods("PUT")
	r.HandleFunc("/mempool/transactions", PutTransaction).Methods("PUT")
	r.HandleFunc("/mempool/transactions", GetTransactions).Methods("GET")
	r.HandleFunc("/mempool/transactions/{hash}", GetTransaction).Methods("GET")
	r.HandleFunc("/root", GetRootBlock).Methods("GET")
	return r
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := Store.GetTransactions()
	if err != nil && err.Error() == "Bucket access error" {
		json.NewEncoder(w).Encode([]blockchain.Transaction{})
	} else if err == nil {
		json.NewEncoder(w).Encode(transactions)
	} else {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - couldn't get transactions"))
	}

}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hash, err := base58.Decode(params["hash"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Couldn't decode base58"))
	}

	transaction, err := Store.GetTransaction(hash, true)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - couldn't get transaction"))
	}
	json.NewEncoder(w).Encode(transaction)
}

func GetBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	hash, err := base58.Decode(params["hash"])
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Couldn't decode base58"))
	}

	data, err := Store.Get([]byte("blocks"), hash)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - couldn't get block"))
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
	hash, err := Store.Get([]byte("blocks"), []byte("root"))
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - couldn't get block"))
	}

	blockData, err := Store.Get([]byte("blocks"), hash)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - couldn't get block"))
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
	err = Store.AddBlock(block)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - block isn't valid"))
	}

	w.WriteHeader(http.StatusCreated)
}

func PutTransaction(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var transaction blockchain.Transaction
	err := dec.Decode(&transaction)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - couldn't decode transaction"))
	}
	err = Store.AddTransaction(transaction)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
	}

	w.WriteHeader(http.StatusCreated)
}
