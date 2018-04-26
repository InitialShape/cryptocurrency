package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	cbor "github.com/whyrusleeping/cbor/go"
	"github.com/InitialShape/blockchain/utils"
	"log"
	"time"
)

type Store struct {
	DB   *bolt.DB
	Peer *Peer
}

func (s *Store) Open(location string, peer *Peer) error {
	db, err := bolt.Open(location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	s.DB = db
	s.Peer = peer
	return err
}

func (s *Store) Put(bucket []byte, key []byte, value []byte) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		err = b.Put(key, value)
		return err
	})

	return err
}

func (s *Store) Get(bucket []byte, key []byte) ([]byte, error) {
	var data []byte
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b != nil {
			v := b.Get(key)
			data = append(data, v...)
			if v == nil {
				return errors.New("EOF")
			}
			return nil
		} else {
			return errors.New("Bucket access error")
		}
	})

	return data, err
}

func (s *Store) Delete(bucket []byte, key []byte) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b != nil {
			err := b.Delete(key)
			return err
		} else {
			return errors.New("Bucket access error")
		}
	})

	return err
}

func (s *Store) StoreGenesisBlock(difficulty int) (Block, error) {
	publicKey, privateKey, err := utils.GetWallet()

	block, err := GenerateGenesisBlock(publicKey, privateKey, difficulty)
	if err != nil {
		return Block{}, err
	}

	s.storeBlock(block)
	base58Hash, err := block.GetBase58Hash()
	if err != nil {
		return Block{}, err
	}
	fmt.Println(base58Hash)
	return block, err
}

func (s *Store) AddTransaction(transaction Transaction) error {
	cbor, err := transaction.GetCBOR()
	if err != nil {
		return err
	}

	transactionMempool, err := s.GetTransaction(transaction.Hash, true)
	// TODO: Ideally we'd check for an empty struct here
	if len(transactionMempool.Hash) == 0 {
		transactionBlocks, _ := s.GetTransaction(transaction.Hash, false)
		if len(transactionBlocks.Hash) == 0 {
			go s.Peer.GossipTransaction(transaction)
		}
	}

	err = s.Put([]byte("mempool"), transaction.Hash, cbor)
	return err
}

func (s *Store) AddPeer(peer string) error {
	err := s.Put([]byte("peers"), []byte(peer), []byte(peer))
	return err
}

func (s *Store) GetTransaction(hash []byte, mempool bool) (Transaction, error) {
	var data []byte
	var err error
	if mempool {
		data, err = s.Get([]byte("mempool"), hash)
	} else {
		data, err = s.Get([]byte("transactions"), hash)
	}
	if err != nil {
		return Transaction{}, err
	}

	var transaction Transaction
	dec := cbor.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&transaction)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, err
}

func (s *Store) GetTransactions() ([]Transaction, error) {
	var transactions []Transaction

	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("mempool"))

		if b != nil {
			b.ForEach(func(k, v []byte) error {
				var transaction Transaction
				dec := cbor.NewDecoder(bytes.NewReader(v))
				err := dec.Decode(&transaction)
				if err != nil {
					return err
				}
				transactions = append(transactions, transaction)

				return nil
			})
		} else {
			return errors.New("Bucket access error")
		}
		return nil
	})

	return transactions, err
}

func (s *Store) GetPeers() ([]string, error) {
	var peers []string
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("peers"))

		if b != nil {
			b.ForEach(func(k, v []byte) error {
				peers = append(peers, string(v))

				return nil
			})
			return nil
		} else {
			return errors.New("Bucket access error")
		}
	})
	return peers, err
}

func (s *Store) DeletePeer(peer string) error {
	err := s.Delete([]byte("peers"), []byte(peer))
	return err
}

func (s *Store) GetChain() ([]Block, error) {
	var blocks []Block
	previousBlock, err := s.Get([]byte("blocks"), []byte("root"))
	if err != nil {
		return blocks, err
	}

	for {
		data, err := s.Get([]byte("blocks"), previousBlock)
		if err != nil {
			return blocks, err
		}

		var root Block
		dec := cbor.NewDecoder(bytes.NewReader(data))
		err = dec.Decode(&root)
		if err != nil {
			return blocks, err
		}

		// preprend
		blocks = append([]Block{root}, blocks...)

		if len(root.PreviousBlock) > 0 {
			previousBlock = root.PreviousBlock
		} else {
			break
		}
	}
	return blocks, nil
}

func (s *Store) VerifyTransaction(transaction Transaction, index int) (bool, error) {
	// TODO: Cannot verify if dependent transaction is in block
	if index == 0 && len(transaction.Inputs[0].TransactionHash) == 0 &&
		transaction.Inputs[0].OutputID == 0 {
		return true, nil
	}

	_, err := s.GetTransaction(transaction.Hash, false)
	if err != nil && (err.Error() == "Bucket access error" || err.Error() == "EOF") {
		// Check if transactions are valid here
		for index, input := range transaction.Inputs {
			inputTransaction, err := s.GetTransaction(input.TransactionHash, false)

			if err != nil && (err.Error() == "Bucket access error" ||
				err.Error() == "EOF") {
				return false, errors.New("Input transaction doesn't exist")
			} else {
				pointer := fmt.Sprintf("-%d", index)
				pointerBytes := append(input.TransactionHash, pointer...)
				_, err := s.Get([]byte("utxo"), pointerBytes)
				if err != nil && (err.Error() == "Bucket access error" ||
					err.Error() == "EOF") {
					// output unspendable as doesn't exist
					return false, errors.New("Output doesn't exist (anymore?)")
				} else {
					inputPublicKey := inputTransaction.Outputs[index].PublicKey
					return transaction.Verify(inputPublicKey, index)
				}
			}
		}
		return false, errors.New("Shouldn't reach statement")
	} else if err == nil {
		fmt.Println("Transaction with hash exists already", transaction.Hash, index)
		return false, errors.New("Transaction exists already")
	} else {
		// should not be reached
		log.Fatal(err)
		return false, err
	}
}

func (s *Store) storeBlock(block Block) error {
	cbor, err := block.GetCBOR()
	if err != nil {
		return err
	}

	err = s.Put([]byte("blocks"), block.Hash, cbor)
	if err != nil {
		log.Fatal("Error storing block", err)
		return err
	}
	err = s.Put([]byte("blocks"), []byte("root"), block.Hash)
	if err != nil {
		log.Fatal("Error overwriting root", err)
		return err
	}
	log.Println("Block added successfully")

	// store transactions in transactions bucket
	for _, transaction := range block.Transactions {
		transactionCbor, err := transaction.GetCBOR()
		if err != nil {
			log.Fatal("Couldn't transform transaction to cbor: ", err)
			return err
		}
		err = s.Put([]byte("transactions"), transaction.Hash,
			transactionCbor)
		if err != nil {
			log.Fatal("Error storing transaction", err)
			return err
		}

		for index, output := range transaction.Outputs {
			outputCbor, err := output.GetCBOR()
			if err != nil {
				log.Fatal("Error encoding output to cbor", err)
			}

			pointer := fmt.Sprintf("-%d", index)
			pointerBytes := append(transaction.Hash, pointer...)
			err = s.Put([]byte("utxo"), pointerBytes, outputCbor)
		}
	}
	return err
}

func (s *Store) EvaluateChains(chains [][]Block) []Block {
	difficulties := make([]int, len(chains))
	for index, chain := range chains {
		for _, block := range chain {
			// take actual block difficulty in bits here
			difficulties[index] += block.Difficulty
		}
	}

	var mostWork int
	for index, difficulty := range difficulties {
		if index == 0 {
			mostWork = 0
		} else if difficulty > difficulties[mostWork] {
			mostWork = index
		}
	}

	return chains[mostWork]
}

func (s *Store) AddBlock(block Block) error {
	data, err := s.Get([]byte("blocks"), block.PreviousBlock)
	if err != nil {
		log.Println("Chain ran out of sync, getting blocks from peers")
		chains, err := s.Peer.Download()
		if err != nil {
			log.Println("Error downloading peer chain", err)
		}

		chain := s.EvaluateChains(chains)
		chain = append(chain, block)

		for _, block := range chain[1:] {
			s.AddBlock(block)
		}
	}

	var root Block
	dec := cbor.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&root)
	if err != nil {
		return err
	}

	if HashMatchesDifficulty(block.Hash, root.Difficulty) {

		// check for duplicates in block
		visited := make(map[string]bool)
		for _, transaction := range block.Transactions {
			if visited[string(transaction.Hash)] {
				return errors.New("Transaction duplicate in block")
			} else {
				visited[string(transaction.Hash)] = true
			}
		}

		// verify transactions' integrity
		for index, transaction := range block.Transactions {
			_, err := s.VerifyTransaction(transaction, index)
			if err != nil {
				return err
			}
		}

		_, err := s.Get([]byte("blocks"), block.Hash)
		if err != nil {
			go s.Peer.GossipBlock(block)
		}

		err = s.storeBlock(block)
		if err != nil {
			log.Fatal("Error storing block", err)
			return err
		}

		// delete all transactions from mempool
		for _, transaction := range block.Transactions {
			s.Delete([]byte("mempool"), transaction.Hash)

			for index, input := range transaction.Inputs {
				pointer := fmt.Sprintf("-%d", index)
				pointerBytes := append(input.TransactionHash, pointer...)
				s.Delete([]byte("utxo"), pointerBytes)
			}
		}

		return err
	} else {
		return errors.New("Difficulty too low")
	}

}
