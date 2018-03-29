package blockchain

import (
	"github.com/syndtr/goleveldb/leveldb"
	"fmt"
	cbor "github.com/whyrusleeping/cbor/go"
	"errors"
	"bytes"
)

type Store struct {
	DB string
}

func (s *Store) AddBlock(block Block) (error) {
	db, err := leveldb.OpenFile(s.DB, nil)
	defer db.Close()
	if err != nil {
		return err
	}

	data, err := db.Get([]byte("root"), nil)
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
			fmt.Printf("%b", n)
		}
		return err
	} else {
		return errors.New("Difficulty too low")
	}

}
