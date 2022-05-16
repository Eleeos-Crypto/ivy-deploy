package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	hederapkg "github.com/eleeos-crypto/ivy-deploy/internal/hedera"
	"github.com/eleeos-crypto/ivy-deploy/internal/notify"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	numAccounts     int
	accountsOutFile string
)

func init() {
	rootCmd.AddCommand(createAccountsCmd)
	createAccountsCmd.PersistentFlags().IntVar(&numAccounts, "count", 1, "The number of accounts to create")
	createAccountsCmd.PersistentFlags().StringVar(&accountsOutFile, "out", "accounts.json", "The file to save the accounts to")
}

type HederaAccount struct {
	AccountID  string `json:"accountId"`
	PrivateKey string `json:"privKey"`
}

func WriteAccountsFile(filename string, dtl []HederaAccount) error {
	file, err := json.MarshalIndent(dtl, "", " ")

	err = ioutil.WriteFile(filename, file, 0644)

	return err
}

var createAccountsCmd = &cobra.Command{
	Use:   "create-accounts",
	Short: "Create a set of accounts",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := hederapkg.SetupClient(viper.Get("Network").(string), viper.Get("Operator ID").(string), viper.Get("Operator Key").(string))
		if err != nil {
			notify.PrettyError("Could not setup client")
			os.Exit(1)
		}

		notify.PrettyInfo(fmt.Sprintf("Creating %d accounts", numAccounts))

		var createdAccounts []HederaAccount

		for i := 0; i < numAccounts; i++ {
			// Generate new key to use with new account
			newKey, err := hedera.GeneratePrivateKey()
			if err != nil {
				notify.PrettyError("Error creating private key")
				os.Exit(1)
			}

			response, err := hedera.NewAccountCreateTransaction().
				SetKey(newKey.PublicKey()).
				SetReceiverSignatureRequired(false).
				SetMaxAutomaticTokenAssociations(1).
				SetTransactionMemo("Test Account").
				Execute(client)

			if err != nil {
				notify.PrettyError(fmt.Sprintf("Error creating account: %s", err.Error()))
				os.Exit(1)
			}

			// Get receipt to see if transaction succeeded, and has the account ID
			transactionReceipt, err := response.GetReceipt(client)
			if err != nil {
				notify.PrettyError(fmt.Sprintf("Error retrieving receipt"))
				os.Exit(1)
			}

			createdAccounts = append(createdAccounts, HederaAccount{
				AccountID:  transactionReceipt.AccountID.String(),
				PrivateKey: newKey.String(),
			})
		}

		WriteAccountsFile(accountsOutFile, createdAccounts)

		notify.PrettyInfo(fmt.Sprintf("Accounts created"))

		defer func() {
			err = client.Close()
			if err != nil {
				println(err.Error(), ": error closing client")
				return
			}
		}()
	},
}
