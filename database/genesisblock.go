package database

import "time"

type GenesisBlock struct {
	GenesisTime time.Time `json:"genesis_time"`
	ChainID     string    `json:"chain_id"`
	Balances map[Account]int `json:"balances"`
}