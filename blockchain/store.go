package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/mr-tron/base58/base58"
	cbor "github.com/whyrusleeping/cbor/go"
	"time"
)

type Store struct {
	DB *bolt.DB
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
		err := b.Delete(key)
		return err
	})

	return err
}

func (s *Store) StoreGenesisBlock(difficulty int) ([]byte, error) {
	publicKey, _ := base58.Decode("6zjRZQyp47BjwArFoLpvzo8SHwwWeW571kJNiqWfSrFT")
	privateKey, _ := base58.Decode("35DxrJipeuCAakHNnnPkBjwxQffYWKM1632kUFv9vKGRNREFSyM6awhyrucxTNbo9h693nPKeWonJ9sFkw6Tou4d")

	block, err := GenerateGenesisBlock(publicKey, privateKey, difficulty)
	if err != nil {
		return []byte{}, err
	}

	cbor, err := block.GetCBOR()
	if err != nil {
		return []byte{}, err
	}

	err = s.Put([]byte("blocks"), block.Hash, cbor.Bytes())
	if err != nil {
		return []byte{}, err
	}
	err = s.Put([]byte("blocks"), []byte("root"), block.Hash)
	if err != nil {
		return []byte{}, err
	}

	base58Hash, err := block.GetBase58Hash()
	if err != nil {
		return []byte{}, err
	}
	fmt.Println(base58Hash)
	return block.Hash, err
}

func (s *Store) AddTransaction(transaction Transaction) error {
	cbor, err := transaction.GetCBOR()
	if err != nil {
		return err
	}

	newTransaction, err := s.GetTransaction(transaction.Hash)
	if newTransaction.Hash == nil{
		go s.Peer.GossipTransaction(transaction)
	}

	err = s.Put([]byte("mempool"), transaction.Hash, cbor.Bytes())
	return err
}

func (s *Store) AddPeer(peer string) error {
	err := s.Put([]byte("peers"), []byte(peer), []byte(peer))
	return err
}

func (s *Store) GetTransaction(hash []byte) (Transaction, error) {
	data, err := s.Get([]byte("mempool"), hash)
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

		b.ForEach(func(k, v []byte) error {
			peers = append(peers, string(v))

			return nil
		})
		return nil
	})
	return peers, err
}

func (s *Store) DeletePeer(peer string) error {
	err := s.Delete([]byte("peers"), []byte(peer))
	return err
}

func (s *Store) AddBlock(block Block) error {
	data, err := s.Get([]byte("blocks"), block.PreviousBlock)
	if err != nil {
		return err
	}

	var root Block
	dec := cbor.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&root)
	if err != nil {
		return err
	}

	// Check if transactions are valid here and delete from mempool
	if HashMatchesDifficulty(block.Hash, root.Difficulty) {
		cbor, err := block.GetCBOR()
		if err != nil {
			return err
		}

		err = s.Put([]byte("blocks"), block.Hash, cbor.Bytes())
		if err != nil {
			return err
		}
		err = s.Put([]byte("blocks"), []byte("root"), block.Hash)
		if err != nil {
			return err
		}

		fmt.Println("Put new block as root with hash and difficulty ", root.Difficulty)
		for _, n := range block.Hash {
			fmt.Printf("%08b", n)
		}
		return err
	} else {
		return errors.New("Difficulty too low")
	}

}
