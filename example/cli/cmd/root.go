/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	goChargify "github.com/GetWagz/go-chargify"
	"github.com/GetWagz/go-chargify/example/cli/cmd/customers"
	"github.com/GetWagz/go-chargify/example/cli/cmd/events"
	"github.com/GetWagz/go-chargify/example/cli/cmd/productFamilies"
	"github.com/GetWagz/go-chargify/example/cli/cmd/products"
	"github.com/GetWagz/go-chargify/example/cli/cmd/subscriptions"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: rootPersistentPreRunE,
}

func rootPersistentPreRunE(cmd *cobra.Command, args []string) error {
	var err error
	chargifyApiKey := viper.GetString("chargify-api-key")
	if len(chargifyApiKey) == 0 {
		err = fmt.Errorf("--chargify-api-key missing and env:CHARGIFY_API_KEY not present")
	}
	subdomain := viper.GetString("subdomain")
	if len(subdomain) == 0 {
		err = fmt.Errorf("--subdomain missing and env:CHARGIFY_SUBDOMAIN not present")
	}
	goChargify.SetCredentials(subdomain, chargifyApiKey)
	return err
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var apiKey string
var subdomain string

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// chargify-api-key
	rootCmd.PersistentFlags().StringVar(&apiKey, "chargify-api-key", "", "the chargify api key, aka personal access token")
	viper.BindPFlag("chargify-api-key", rootCmd.PersistentFlags().Lookup("chargify-api-key"))
	viper.BindEnv("chargify-api-key", "CHARGIFY_API_KEY")

	// subdomain
	rootCmd.PersistentFlags().StringVar(&subdomain, "subdomain", "", "the chargify subdomain")
	viper.BindPFlag("subdomain", rootCmd.PersistentFlags().Lookup("subdomain"))
	viper.BindEnv("subdomain", "CHARGIFY_SUBDOMAIN")

	productFamilies.Init(rootCmd)
	products.Init(rootCmd)
	customers.Init(rootCmd)
	subscriptions.Init(rootCmd)
	events.Init(rootCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
