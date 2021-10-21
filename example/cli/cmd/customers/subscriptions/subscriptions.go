package subscriptions

import (
	"fmt"

	goChargify "github.com/GetWagz/go-chargify"
	"github.com/GetWagz/go-chargify/example/cli/cmd/customers/shared"
	"github.com/GetWagz/go-chargify/example/cli/internal"
	"github.com/GetWagz/go-chargify/example/cli/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var command = &cobra.Command{
	Use:               "subscriptions",
	Short:             "query customer subscriptions",
	Long:              `example:  customers subscriptions --subdomain=mapped-3 --customer-id=47169214`,
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

	if shared.CustomerID == internal.NotSetIntParam {
		err = fmt.Errorf("--customer-id has not been set")
	}
	return err
}

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(command)

}
func execute() error {
	apiKey := viper.Get("chargify-api-key")
	log.Debug().
		Str("api-key", apiKey.(string)).
		Int("customer-id", shared.CustomerID).
		Send()

	found, err := goChargify.GetCustomerSubscriptions(shared.CustomerID)
	if err != nil {
		return err
	}
	fmt.Println(utils.PrettyJSON(found))
	return nil
}
