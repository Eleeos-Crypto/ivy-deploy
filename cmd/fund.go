package cmd

import (
	"github.com/spf13/cobra"
)

var (
	account string
	hbar    int64
	tinybar int64
)

func init() {
	rootCmd.AddCommand(fundCmd)
	rootCmd.PersistentFlags().StringVar(&account, "account", "", "The account to fund in the form <shard>.<realm>.<account #>")
	rootCmd.PersistentFlags().Int64Var(&hbar, "hbar", 0, "The amount in HBAR to send")
	rootCmd.PersistentFlags().Int64Var(&tinybar, "tinybar", 0, "The amount in tiny bar to send")
}

var fundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Fund an account using the provided operator ID/key",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
