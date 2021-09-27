package customers

import (
	"fmt"

	goChargify "github.com/GetWagz/go-chargify"
	"github.com/GetWagz/go-chargify/example/cli/cmd/customers/shared"
	"github.com/GetWagz/go-chargify/example/cli/cmd/customers/subscriptions"
	"github.com/GetWagz/go-chargify/example/cli/internal"
	"github.com/GetWagz/go-chargify/example/cli/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var command = &cobra.Command{
	Use:               "customers",
	Short:             "query customers",
	Long:              `example: customers --subdomain=mapped-3`,
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
	if !(sort == "desc" || sort == "asc") {
		err = fmt.Errorf("--sort must be asc or desc")
		return err
	}
	return nil
}

var reference string

var page int
var sort string

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(command)

	// reference-id
	command.Flags().StringVar(&reference, "reference", internal.NotSetStringParam, "the customers referenceID.  This is an external foreign guid")
	// customer_id
	command.PersistentFlags().IntVar(&shared.CustomerID, "customer-id", internal.NotSetIntParam, "the customers id.")
	// page
	command.Flags().IntVar(&page, "page", 1, "the page number")
	// sort
	command.Flags().StringVar(&sort, "sort", "asc", "sort direction, asc or desc")

	subscriptions.Init(command)
}
func execute() error {
	apiKey := viper.Get("chargify-api-key")
	log.Debug().
		Str("api-key", apiKey.(string)).
		Int("id", shared.CustomerID).
		Str("reference", reference).
		Str("sort", sort).
		Int("page", page).
		Send()
	if len(reference) != 0 {
		found, err := goChargify.GetCustomerByReference(reference)
		if err != nil {
			return err
		}
		fmt.Println(utils.PrettyJSON(found))
		return nil
	}
	if shared.CustomerID != internal.NotSetIntParam {
		found, err := goChargify.GetCustomerByID(shared.CustomerID)
		if err != nil {
			return err
		}
		fmt.Println(utils.PrettyJSON(found))
		return nil
	}
	found, err := goChargify.GetCustomers(page, sort)
	if err != nil {
		return err
	}
	fmt.Println(utils.PrettyJSON(found))
	return nil
}
