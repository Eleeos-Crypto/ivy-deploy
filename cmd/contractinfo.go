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
	contractId          string
	contractInfoOutFile string
)

func init() {
	rootCmd.AddCommand(contractInfoQueryCmd)
	contractInfoQueryCmd.PersistentFlags().StringVar(&contractId, "contract", "", "The contract ID to query")
	contractInfoQueryCmd.PersistentFlags().StringVar(&contractInfoOutFile, "out", "contractInfo.json", "The file to save the contract info to")
}

func WriteContractInfoFile(filename string, dtl hedera.ContractInfo) error {
	file, err := json.MarshalIndent(dtl, "", " ")

	err = ioutil.WriteFile(filename, file, 0644)

	return err
}

var contractInfoQueryCmd = &cobra.Command{
	Use:   "contract-info",
	Short: "Get info for a contract",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := hederapkg.SetupClient(viper.Get("Network").(string), viper.Get("Operator ID").(string), viper.Get("Operator Key").(string))
		if err != nil {
			notify.PrettyError("Could not setup client")
			os.Exit(1)
		}

		contract, err := hedera.ContractIDFromString(contractId)
		if err != nil {
			notify.PrettyError("Invalid Hedera Contract ID")
			os.Exit(1)
		}

		//Create the query
		contractInfo, err := hedera.NewContractInfoQuery().
			SetContractID(contract).
			SetQueryPayment(hedera.NewHbar(1)).
			SetMaxQueryPayment(hedera.NewHbar(5)).Execute(client)

		if err != nil {
			notify.PrettyError(fmt.Sprintf("Error retrieving contract info: %s", err.Error()))
			os.Exit(1)
		}

		WriteContractInfoFile(contractInfoOutFile, contractInfo)

		notify.PrettyInfo(fmt.Sprintf("Contract info written to: %s", contractInfoOutFile))

		defer func() {
			err = client.Close()
			if err != nil {
				println(err.Error(), ": error closing client")
				return
			}
		}()
	},
}
