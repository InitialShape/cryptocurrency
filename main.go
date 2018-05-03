package main

import (
	"fmt"
	"github.com/InitialShape/cryptocurrency/blockchain"
	"github.com/InitialShape/cryptocurrency/web"
	"github.com/InitialShape/cryptocurrency/utils"
	"log"
	"net/http"
	"os"
	"flag"
)

func main() {
	var store blockchain.Store
	var peer blockchain.Peer

	keys := flag.Bool("generate_keys", false,
					  "Generates keys for the wallet and miner")
	flag.Parse()
	if *keys {
			// key generation mode
			utils.GenerateWallet()
	} else {
		// normal operation mode
		store = blockchain.Store{}
		// TODO: For other os.Args use flag lib
		store.Open(os.Args[1], &peer)

		ip, err := utils.GetExternalIP()
		if err != nil {
			log.Fatal(err)
		}

		peer = blockchain.Peer{os.Args[2], ip, store}
		go peer.Start()

		_, err = store.StoreGenesisBlock(20)
		if err != nil {
			log.Fatal(err)
		}

		r := web.Handlers(store)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Args[3]), r))
	}

}
