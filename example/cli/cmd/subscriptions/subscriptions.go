package subscriptions

import (
	"fmt"

	"github.com/GetWagz/go-chargify/example/cli/cmd/subscriptions/events"
	"github.com/GetWagz/go-chargify/example/cli/cmd/subscriptions/shared"
	"github.com/GetWagz/go-chargify/example/cli/internal"
	"github.com/GetWagz/go-chargify/example/cli/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var command = &cobra.Command{
	Use:               "subscriptions",
	Short:             "query subscriptions",
	PersistentPreRunE: persistentPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		return execute()
	},
}

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	err := utils.ParentPersistentPreRunE(cmd, args)
	if err != nil {
		return err
	}

	return nil
}
func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(command)

	// number-of-consumers
	command.PersistentFlags().IntVar(&shared.SubscriptionID, "subscription-id", internal.NotSetIntParam, "the subscription id, i.e. 1804402")
	command.MarkPersistentFlagRequired("subscription-id")
	viper.BindPFlag("subscription-id", command.PersistentFlags().Lookup("subscription-id"))

	events.Init(command)

}
func execute() error {
	apiKey := viper.Get("chargify-api-key")
	log.Debug().Str("api-key", apiKey.(string)).Send()

	fmt.Printf("products/%v", shared.SubscriptionID)
	return nil

}
