package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type BlockHeader struct {
	Parent Hash
	Time   uint64
	Number uint64
}

type Block struct {
	Header          BlockHeader
	NewTransactions []Transaction
}

type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
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

func NewBlock(parentHash Hash, number uint64, time uint64, newTransactions []Transaction) Block {
	return Block{BlockHeader{parentHash, time, number}, newTransactions}
}
