package main

import (
	"fmt"

	"github.com/FA5I/blockchain/database"
	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	balancesCmd.AddCommand(balancesListCmd())

	return balancesCmd
}

func balancesListCmd() *cobra.Command {
	var balancesListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all balances.",
		Run: func(cmd *cobra.Command, args []string) {
			state, err := database.NewStateFromDisc()
			if err != nil {
				panic(err)
			}

			// fmt.Printf("Account balances at %x\n", state.)
			fmt.Println("=================")
			for account, balance := range state.Balances {
				fmt.Printf("%s: %d\n", account, balance)
			}
		},
	}

	return balancesListCmd
}
