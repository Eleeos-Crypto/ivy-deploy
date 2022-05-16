package cmd

import (
	"fmt"
	"os"

	hederapkg "github.com/eleeos-crypto/ivy-deploy/internal/hedera"
	"github.com/eleeos-crypto/ivy-deploy/internal/notify"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(balanceQueryCmd)
	balanceQueryCmd.PersistentFlags().StringVar(&account, "account", "", "The account to get the balance of")
}

var balanceQueryCmd = &cobra.Command{
	Use:   "balance",
	Short: "Get the balance of an account",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := hederapkg.SetupClient(viper.Get("Network").(string), viper.Get("Operator ID").(string), viper.Get("Operator Key").(string))
		if err != nil {
			notify.PrettyError("Could not setup client")
			os.Exit(1)
		}

		account, err := hedera.AccountIDFromString(account)
		if err != nil {
			notify.PrettyError("Invalid Hedera Account ID")
			os.Exit(1)
		}

		//Create the query
		tokenBalance, err := hedera.NewAccountBalanceQuery().
			SetAccountID(account).
			Execute(client)

		if err != nil {
			notify.PrettyError(fmt.Sprintf("Error retrieving account balance: %s", err.Error()))
			os.Exit(1)
		}

		notify.PrettyInfo(fmt.Sprintf("Account Balance Tinybars: %d", tokenBalance.Hbars.AsTinybar()))

		defer func() {
			err = client.Close()
			if err != nil {
				println(err.Error(), ": error closing client")
				return
			}
		}()
	},
}
