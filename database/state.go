package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type State struct {
	Balances        map[Account]int `json:"balances"`
	txMempool       []Transaction
	dbFile          *os.File
	LatestBlockHash Hash
}

/*
Creates the latest state by building up all historical transactions
from the Genesis block until current.
*/
func NewStateFromDisc(dataDir string) (*State, error) {
	err := initDataDirIfNotExists(dataDir)
	if err != nil {
		return nil, err
	}

	gen, err := loadGenesis(getGenesisJsonFilePath(dataDir))
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]int)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	blockFile, err := os.OpenFile(getBlocksDbFilePath(dataDir), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(blockFile)
	state := &State{balances, make([]Transaction, 0), blockFile, Hash{}}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Bytes()
		var blockfs BlockFS
		err := json.Unmarshal(blockFsJson, &blockfs)
		if err != nil {
			return nil, err
		}

		err = state.applyBlock(blockfs.Value)
		if err != nil {
			return nil, err
		}

		state.LatestBlockHash = blockfs.Key
	}

	return state, nil
}

// persist transactions in the mempool to the database
func (s *State) Persist() (Hash, error) {
	// create a new block with new transactions
	block := NewBlock(s.LatestBlockHash, uint64(time.Now().Unix()), s.txMempool)

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, nil
	}

	blockfs := BlockFS{blockHash, block}

	blockfsJson, err := json.Marshal(blockfs)
	if err != nil {
		return Hash{}, nil
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n\n", blockfsJson)

	_, err = s.dbFile.Write(append(blockfsJson, '\n'))
	if err != nil {
		return Hash{}, nil
	}

	s.LatestBlockHash = blockHash

	s.txMempool = []Transaction{}

	return s.LatestBlockHash, nil

}

/*
Apply effect of a transaction to balances
*/
func (s *State) apply(tx Transaction) error {
	if tx.isReward() {
		s.Balances[tx.To] += int(tx.Value)
	}

	if s.Balances[tx.From] < int(tx.Value) {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= int(tx.Value)
	s.Balances[tx.To] += int(tx.Value)

	return nil
}

/*
Add a transaction to the mempool
*/
func (s *State) AddTx(tx Transaction) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)
	return nil
}

/*
Add all transactions in a block to the mempool
*/
func (s *State) AddBlock(b Block) error {
	for _, tx := range b.NewTransactions {
		if err := s.AddTx(tx); err != nil {
			return err
		}
	}
	return nil
}

/*
Apply the effect of transactions in a block to the balances
*/
func (s *State) applyBlock(b Block) error {
	for _, tx := range b.NewTransactions {
		if err := s.apply(tx); err != nil {
			return err
		}
	}

	fmt.Println("Block applied.")
	return nil
}
