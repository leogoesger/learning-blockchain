package main

import (
	"fmt"
	"log"

	"github.com/leogoesger/learning-blockchain/database"
	"github.com/spf13/cobra"
)

func txCmd() *cobra.Command {
	txCmd := cobra.Command{
		Use:   "tx",
		Short: "",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	txCmd.AddCommand(txAddCmd())

	return &txCmd
}

const (
	flagFrom  = "from"
	flagTo    = "to"
	flagValue = "value"
	flagData  = "data"
)

func txAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds new TX to db",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)
			data, _ := cmd.Flags().GetString(flagData)

			tx := database.NewTx(database.NewAccount(from), database.NewAccount(to), value, data)

			state, err := database.NewStateFromDisk()
			if err != nil {
				log.Fatal(err)
			}
			defer state.Close()

			if err := state.AddTx(tx); err != nil {
				log.Fatal(err)
			}

			if _, err := state.Persist(); err != nil {
				log.Fatal(err)
			}

			fmt.Println("TX successfully persisted to the ledger!")
		},
	}

	cmd.Flags().String(flagFrom, "", "From what account to send tokens")
	cmd.MarkFlagRequired(flagFrom)

	cmd.Flags().String(flagTo, "", "To what account to send tokens")
	cmd.MarkFlagRequired(flagTo)

	cmd.Flags().Uint(flagValue, 0, "How many tokens to send")
	cmd.MarkFlagRequired(flagValue)

	cmd.Flags().String(flagData, "", "Possible values: 'reward'")

	return cmd
}
