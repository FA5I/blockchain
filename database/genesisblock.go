package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const layoutISO = "2006-01-02"

type GenesisBlock struct {
	GenesisTime time.Time       `json:"genesis_time"`
	ChainID     string          `json:"chain_id"`
	Balances    map[Account]int `json:"balances"`
}

func loadGenesis(path string) (GenesisBlock, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return GenesisBlock{}, err
	}

	var loadedGenesis GenesisBlock

	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return GenesisBlock{}, nil
	}

	fmt.Println("Genesis file loaded.")
	return loadedGenesis, nil
}

func writeGenesisToDisk(path string) error {
	time, err := time.Parse(layoutISO, "2019-03-18")
	if err != nil {
		return err
	}

	genesis := GenesisBlock{time, "gochain", map[Account]int{"alice": 10, "bob": 10, "charlie": 10}}
	genesisJson, _ := json.Marshal(genesis)
	ioutil.WriteFile(path, []byte(genesisJson), 0644)

	fmt.Println("Genesis file written to disk.")
	return nil
}

func writeEmptyBlocksDbToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
}
