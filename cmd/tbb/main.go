package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	tbbCmd := &cobra.Command{
		Use:   "tbb",
		Short: "Learning blockchain",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balancesCmd())
	tbbCmd.AddCommand(txCmd())

	if err := tbbCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
