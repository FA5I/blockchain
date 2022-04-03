package main

import (
	"encoding/json"
	"fmt"

	"github.com/FA5I/blockchain/database"
	"github.com/spf13/cobra"
)

const flagFrom = "from"
const flagTo = "to"
const flagValue = "value"
const flagData = "data"

func txCmd() *cobra.Command {
	var txCmd = &cobra.Command{
		Use:   "transaction",
		Short: "Interact with tx (add...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	txCmd.AddCommand(txAddCmd())

	return txCmd
}

func txAddCmd() *cobra.Command {
	var txAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a transaction to database.",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)
			data, _ := cmd.Flags().GetString(flagData)

			fromAcc := database.NewAccount(from)
			toAcc := database.NewAccount(to)

			tx := database.NewTransaction(fromAcc, toAcc, value, data)

			state, err := database.NewStateFromDisc()
			if err != nil {
				panic(err)
			}

			// add transaction the mempool
			err = state.Add(tx)
			if err != nil {
				panic(err)
			}

			fmt.Println("Persisting new TX to disk:")
			txJson, _ := json.MarshalIndent(tx, "", "\t")
			fmt.Println(string(txJson))
			fmt.Printf("New snapshot: %x\n", state.LatestSnapshot())

			// persist mempool to the database
			_, err = state.Persist()
			if err != nil {
				panic(err)
			}

			fmt.Println("TX successfully added to the ledger.")

		},
	}

	txAddCmd.Flags().String(flagFrom, "", "From what account to send tokens")
	txAddCmd.MarkFlagRequired(flagFrom)

	txAddCmd.Flags().String(flagTo, "", "To what account to send tokens")
	txAddCmd.MarkFlagRequired(flagTo)

	txAddCmd.Flags().Uint(flagValue, 0, "How many tokens to send")
	txAddCmd.MarkFlagRequired(flagValue)

	txAddCmd.Flags().String(flagData, "", "Possible values: 'reward'")

	return txAddCmd
}
