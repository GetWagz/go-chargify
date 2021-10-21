package products

import (
	"fmt"

	goChargify "github.com/GetWagz/go-chargify"
	"github.com/GetWagz/go-chargify/example/cli/cmd/productFamilies/shared"
	"github.com/GetWagz/go-chargify/example/cli/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var command = &cobra.Command{
	Use:               "products",
	Short:             "query products",
	PersistentPreRunE: persistentPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		return execute()
	},
}
var id string

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	err := utils.ParentPersistentPreRunE(cmd, args)
	if err != nil {
		return err
	}

	return nil
}

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(command)

	// component-id
	command.Flags().StringVar(&id, "id", "", "the product id, i.e. 1804402")

}
func execute() error {
	apiKey := viper.Get("chargify-api-key")

	log.Debug().
		Str("api-key", apiKey.(string)).
		Str("id", id).
		Interface("productFamilyID", shared.ProductFamilyID).
		Send()

	result, err := goChargify.GetProductFamilyProducts(shared.ProductFamilyID)
	if err != nil {
		return err
	}
	fmt.Println(utils.PrettyJSON(result))
	return nil
}
