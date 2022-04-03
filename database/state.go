package database

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Snapshot [32]byte

type State struct {
	Balances  map[Account]int `json:"balances"`
	txMempool []Transaction
	dbFile    *os.File
	snapshot  Snapshot
}

// persist transactions in the mempool to the database
func (s *State) Persist() (Snapshot, error) {
	mempool := make([]Transaction, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {

		txjsn, err := json.Marshal(mempool[i])
		if err != nil {
			return Snapshot{}, err
		}

		if _, err := s.dbFile.Write(append(txjsn, '\n')); err != nil {
			return Snapshot{}, err
		}

		s.txMempool = s.txMempool[1:]
	}

	return s.snapshot, nil

}

// apply a transaction to the balances
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

func (s *State) Add(tx Transaction) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

/*
Create snapshot of blockchain state
*/
func (s *State) DoSnapshot() error {
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		return err
	}

	s.snapshot = sha256.Sum256(data)

	return nil
}

func (s *State) LatestSnapshot() Snapshot {
	return s.snapshot
}

/*
Creates the latest state by building up all historical transactions
from the Genesis block until current.
*/
func NewStateFromDisc() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// read in the genesis block and seed the initial account balances
	genesisPath := filepath.Join(cwd, "database", "genesis.json")

	genesisFile, err := os.Open(genesisPath)
	if err != nil {
		panic(err)
	}
	// defer genesisFile.Close()

	var genesisBlock GenesisBlock
	jsonParser := json.NewDecoder(genesisFile)
	jsonParser.Decode(&genesisBlock)

	balances := make(map[Account]int)
	for account, balance := range genesisBlock.Balances {
		balances[account] = balance
	}

	// Next, read in all the new transactions and update balances sequentially
	txdbPath := filepath.Join(cwd, "database", "transactions.db")
	transactionsFile, err := os.OpenFile(txdbPath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	// defer transactionsFile.Close()

	scanner := bufio.NewScanner(transactionsFile)
	state := &State{balances, make([]Transaction, 0), transactionsFile, Snapshot{}}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Transaction
		json.Unmarshal(scanner.Bytes(), &tx)

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	err = state.DoSnapshot()
	if err != nil {
		return nil, err
	}

	return state, nil
}
