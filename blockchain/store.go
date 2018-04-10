package blockchain

import (
	"fmt"
	cbor "github.com/whyrusleeping/cbor/go"
	"errors"
	"bytes"
	"github.com/mr-tron/base58/base58"
	"github.com/boltdb/bolt"
	"time"
)

type Store struct {
	DB *bolt.DB
}

func (s *Store) Open(location string) error {
	db, err := bolt.Open(location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	s.DB = db
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
		v := b.Get(key)
		data = append(data, v...)
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

func (s *Store) AddBlock(block Block) (error) {
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
