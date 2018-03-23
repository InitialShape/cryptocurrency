package blockchain

import (
	"testing"
	"github.com/2tvenom/cbor"
	"bytes"
	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	block := Block{0, "abc", nil, ""}
	marshalledBlock, err := block.GetCBOR()
	if err != nil {
		t.Error(err)
	}

	var buffer bytes.Buffer
	encoder := cbor.NewEncoder(&buffer)
	ok, err := encoder.Marshal(block)
	if !ok {
		t.Error(err)
	}
	assert.Equal(t, marshalledBlock, buffer)
}

func TestGetHash(t *testing.T) {
	block := Block{0, "abc", nil, ""}
	hash, err := block.GetHash()
	if err != nil {
		t.Error(err)
	}
	expected := "HF58qWumjA9W9Kd7bBcg9G1HM74dZFhvdNgGnVj7AhPu"
	assert.Equal(t, expected, hash)
}
