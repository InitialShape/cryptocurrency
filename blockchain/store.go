package blockchain

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
	cbor "github.com/whyrusleeping/cbor/go"
	"errors"
	"bytes"
	"github.com/mr-tron/base58/base58"
	"github.com/davecgh/go-spew/spew"
)

type Store struct {
	DB string
}

func (s *Store) StoreGenesisBlock(difficulty int) ([]byte, error) {
	publicKey, _ := base58.Decode("6zjRZQyp47BjwArFoLpvzo8SHwwWeW571kJNiqWfSrFT")
	privateKey, _ := base58.Decode("35DxrJipeuCAakHNnnPkBjwxQffYWKM1632kUFv9vKGRNREFSyM6awhyrucxTNbo9h693nPKeWonJ9sFkw6Tou4d")

	block, err := GenerateGenesisBlock(publicKey, privateKey, difficulty)
	if err != nil {
		return []byte{}, err
	}

	db, err := leveldb.OpenFile(s.DB, nil)
	if err != nil {
		return []byte{}, err
	}
	cbor, err := block.GetCBOR()
	if err != nil {
		return []byte{}, err
	}

	err = db.Put(block.Hash, cbor.Bytes(), nil)
	if err != nil {
		return []byte{}, err
	}
	err = db.Put([]byte("root"), block.Hash, nil)
	if err != nil {
		return []byte{}, err
	}
	db.Close()

	base58Hash, err := block.GetBase58Hash()
	if err != nil {
		return []byte{}, err
	}
	fmt.Println(base58Hash)
	return block.Hash, err
}

func (s *Store) AddBlock(block Block) (error) {
	db, err := leveldb.OpenFile(s.DB, nil)
	defer db.Close()
	if err != nil {
		return err
	}

	spew.Dump(block)
	data, err := db.Get(block.PreviousBlock, nil)
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

		db.Put(block.Hash, cbor.Bytes(), nil)
		db.Put([]byte("root"), block.Hash, nil)

		fmt.Println("Put new block as root with hash and difficulty ", root.Difficulty)
		for _, n := range block.Hash {
			fmt.Printf("%08b", n)
		}
		return err
	} else {
		return errors.New("Difficulty too low")
	}

}
