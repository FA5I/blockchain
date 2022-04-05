package main

import (
	"fmt"
	"os"
	"time"

	"github.com/FA5I/blockchain/database"
)

func main() {
	cwd, _ := os.Getwd()

	state, err := database.NewStateFromDisc(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	block0 := database.NewBlock(
		database.Hash{},
		0,
		uint64(time.Now().Unix()),
		[]database.Transaction{
			database.NewTransaction("alice", "bob", 3, ""),
			database.NewTransaction("alice", "alice", 7, "reward"),
		},
	)

	state.AddBlock(block0)
	block0hash, _ := state.Persist()

	block1 := database.NewBlock(
		block0hash,
		1,
		uint64(time.Now().Unix()),
		[]database.Transaction{
			database.NewTransaction("bob", "charlie", 1, ""),
			database.NewTransaction("alice", "bob", 3, ""),
			database.NewTransaction("alice", "charlie", 3, ""),
		},
	)

	state.AddBlock(block1)
	state.Persist()
}
