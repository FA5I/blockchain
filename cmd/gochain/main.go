package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}

func main() {
	var cmd = &cobra.Command{
		Use:   "gochain",
		Short: "blockchain in Go",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	cmd.AddCommand(versionCmd)
	cmd.AddCommand(balancesCmd())
	cmd.AddCommand(txCmd())

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
