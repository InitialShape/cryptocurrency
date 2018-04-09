package blockchain

import (
	"fmt"
	cbor "github.com/whyrusleeping/cbor/go"
	"errors"
	"bytes"
	"github.com/mr-tron/base58/base58"
	"github.com/dgraph-io/badger"
)

type Store struct {
	DB *badger.DB
}

func (s *Store) Open(location string) error {
	opts := badger.DefaultOptions
	opts.Dir = location
	opts.ValueDir = location
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	s.DB = db
	return err
}

func (s *Store) Put(key []byte, value []byte) error {
	err := s.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})

	return err
}

func (s *Store) Get(key []byte) ([]byte, error) {
	var data []byte
	err := s.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		data = append(data, val...)
		return nil
	})

	return data, err
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

	err = s.Put(block.Hash, cbor.Bytes())
	if err != nil {
		return []byte{}, err
	}
	err = s.Put([]byte("root"), block.Hash)
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

func (s *Store) AddBlock(block Block) (error) {
	data, err := s.Get(block.PreviousBlock)
	if err != nil {
		return err
	}

	var root Block
	dec := cbor.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&root)

	if HashMatchesDifficulty(block.Hash, root.Difficulty) {
		cbor, err := block.GetCBOR()
		if err != nil {
			return err
		}

		err = s.Put(block.Hash, cbor.Bytes())
		if err != nil {
			return err
		}
		err = s.Put([]byte("root"), block.Hash)
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
