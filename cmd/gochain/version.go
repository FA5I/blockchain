package main

import (
	"fmt"

	"github.com/spf13/cobra"
)




const Major = "0"
const Minor = "1"
const Fix = "0"
const Verbal = "Transaction Add && Balances List"


var GitCommit string

var versionCmd = &cobra.Command {
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Version: %s.%s.%s-alpha %s %s", Major, Minor, Fix, shortGitCommit(GitCommit), Verbal))
	},
}

func shortGitCommit(fullGitCommit string) string {
	shortCommit := ""
	if len(fullGitCommit) >= 6 {
		shortCommit = fullGitCommit[0:6]
	}

	return shortCommit
}