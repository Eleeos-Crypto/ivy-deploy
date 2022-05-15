package cmd

import (
	"fmt"
	"os"

	"github.com/eleeos-crypto/ivy-deploy/internal/notify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{

		Use:   "ivy",
		Short: "Ivy deploy makes it easy to quickly deploy, test, and interact with the Hedera Hashgraph",
		Long:  `Ivy deploy makes it easy to quickly deploy, test, and interact with the Hedera Hashgraph`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Ivy config file")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		notify.PrettyInfo(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	} else {
		notify.PrettyError("Could not read config file")
		os.Exit(1)
	}
}
