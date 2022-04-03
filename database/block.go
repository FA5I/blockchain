package database

import (
	"crypto/sha256"
	"encoding/json"
)

type Hash [32]byte

type BlockHeader struct {
	Parent Hash
	Time   uint64
}

type Block struct {
	Header          BlockHeader
	NewTransactions []Transaction
}

/*
Now, only hash the latest block
*/
func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}
