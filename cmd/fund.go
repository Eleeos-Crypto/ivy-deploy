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

var (
	account string
	memo    string
	hbar    int64
	tinybar int64
)

func init() {
	rootCmd.AddCommand(transferCmd)
	transferCmd.PersistentFlags().StringVar(&account, "account", "", "The account to fund in the form <shard>.<realm>.<account #>")
	transferCmd.PersistentFlags().StringVar(&memo, "memo", "", "Memo for the transaction")
	transferCmd.PersistentFlags().Int64Var(&hbar, "hbar", 0, "The amount in HBAR to send")
	transferCmd.PersistentFlags().Int64Var(&tinybar, "tinybar", 0, "The amount in tiny bar to send")
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Fund an account using the provided operator ID/key",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := hederapkg.SetupClient(viper.Get("Network").(string), viper.Get("Operator ID").(string), viper.Get("Operator Key").(string))
		if err != nil {
			notify.PrettyError("Could not setup client")
			os.Exit(1)
		}

		// Get the account ID we are sending to
		toAccountID, err := hedera.AccountIDFromString(account)
		if err != nil {
			notify.PrettyError("Invalid account ID")
			os.Exit(1)
		}

		amount := hedera.NewHbar(0)

		//If HBAR is set use that, otherwise use tinybar
		if hbar > 0 {
			amount = hedera.HbarFrom(float64(hbar), hedera.HbarUnits.Hbar)
		} else {
			amount = hedera.HbarFromTinybar(tinybar)
		}

		if amount.AsTinybar() < 1 {
			notify.PrettyError("Transfer amount must be nonzero. Please provide hbar or tinybar amount")
			os.Exit(1)
		}

		notify.PrettyInfo(fmt.Sprintf("Transfering: %s (%d)", amount.String(), amount.AsTinybar()))

		// Send HBARS to the derived address
		transactionResponse, err := hedera.NewTransferTransaction().
			// Hbar has to be negated to denote we are taking out from that account
			AddHbarTransfer(client.GetOperatorAccountID(), hedera.NewHbar(float64(hbar*-1))).
			// If the amount of these 2 transfers is not the same, the transaction will throw an error
			AddHbarTransfer(toAccountID, hedera.NewHbar(float64(hbar))).
			SetTransactionMemo(memo).
			Execute(client)

		if err != nil {
			notify.PrettyError(fmt.Sprintf("Error executing transfer: %s", err.Error()))
			os.Exit(1)
		}

		// Retrieve the receipt to make sure the transaction went through
		_, err = transactionResponse.GetReceipt(client)

		if err != nil {
			notify.PrettyError(fmt.Sprintf("Error retreiving receipt: %s", err.Error()))
			os.Exit(1)
		}

		notify.PrettyInfo(fmt.Sprintf("Transfer completed"))

		defer func() {
			err = client.Close()
			if err != nil {
				println(err.Error(), ": error closing client")
				return
			}
		}()
	},
}
