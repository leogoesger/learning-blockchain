package main

import (
	"fmt"
	"log"

	"github.com/leogoesger/learning-blockchain/database"
	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	balancesCmd := &cobra.Command{
		Use:   "balances",
		Short: "",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	balancesCmd.AddCommand(balancesListCmd)

	return balancesCmd
}

var balancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := database.NewStateFromDisk()
		if err != nil {
			log.Fatal("init new state from disk")
		}
		defer state.Close()

		for account, balance := range state.Balances {
			fmt.Printf("%s: %d\n", account, balance)
		}
	},
}
