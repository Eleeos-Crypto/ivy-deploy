package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	hederapkg "github.com/eleeos-crypto/ivy-deploy/internal/hedera"
	hscsdeploy "github.com/eleeos-crypto/ivy-deploy/internal/hedera/hscs/deploy"
	"github.com/eleeos-crypto/ivy-deploy/internal/notify"
	"github.com/eleeos-crypto/ivy-deploy/internal/truffle"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	truffleFile     string
	contractOutFile string
	gasAmount       uint64
	maxFee          int64
)

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.PersistentFlags().StringVar(&truffleFile, "truffle", "", "The truffle file to deploy from")
	deployCmd.PersistentFlags().StringVar(&contractOutFile, "out", "", "The file to save the smart contract info in")
	deployCmd.PersistentFlags().Uint64Var(&gasAmount, "gas", 1000000, "The amount of gas to use in the transaction")
	deployCmd.PersistentFlags().Int64Var(&maxFee, "maxFee", hedera.NewHbar(5).AsTinybar(), "The max amount to pay for the fee")
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a smart contract onto the hedera network",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := hederapkg.SetupClient(viper.Get("Network").(string), viper.Get("Operator ID").(string), viper.Get("Operator Key").(string))
		if err != nil {
			notify.PrettyError("Could not setup client")
		}

		// If the truffle file is provided, use that
		if truffleFile != "" {
			err = DeployFromTruffle(client)
			if err != nil {
				notify.PrettyError(fmt.Sprintf("Error deploying smart contract to Hedera from truffle file: %s", err.Error()))
				os.Exit(1)
			}
		}

		defer func() {
			err = client.Close()
			if err != nil {
				println(err.Error(), ": error closing client")
				return
			}
		}()
	},
}

func WriteDeploymentToFile(filename string, dtl hscsdeploy.DeploymentDetail) error {
	file, err := json.MarshalIndent(dtl, "", " ")

	err = ioutil.WriteFile(filename, file, 0644)

	return err
}

func DeployFromTruffle(client *hedera.Client) error {
	notify.PrettyWarning("Deploying from truffle file")
	trfl, err := truffle.LoadTruffleFile(truffleFile)
	if err != nil {
		notify.PrettyError("Error reading truffle file")
		return err
	}
	deployDetail, err := hscsdeploy.Deploy(client, maxFee, gasAmount, []byte(trfl.ByteCode))
	if err != nil {
		notify.PrettyError("Error deploying truffle file bytecode")
		return err
	}

	err = WriteDeploymentToFile(contractOutFile, *deployDetail)
	if err != nil {
		notify.PrettyError("Error writing deployment details")
		return err
	}

	return nil
}
