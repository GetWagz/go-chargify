package productFamilies

import (
	"fmt"

	goChargify "github.com/GetWagz/go-chargify"
	"github.com/GetWagz/go-chargify/example/cli/cmd/productFamilies/components"
	"github.com/GetWagz/go-chargify/example/cli/cmd/productFamilies/products"
	"github.com/GetWagz/go-chargify/example/cli/cmd/productFamilies/shared"
	"github.com/GetWagz/go-chargify/example/cli/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var command = &cobra.Command{
	Use:               "product-families",
	Short:             "query product families",
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
	productFamilyIDInterface := viper.Get("product-family-id")

	if productFamilyIDInterface == nil {
		err = fmt.Errorf("--product-family-id is not set or empty")
		return err
	}
	shared.ProductFamilyID = int64(productFamilyIDInterface.(int))

	return nil
}

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(command)

	// product-family-id
	command.PersistentFlags().Int64Var(&shared.ProductFamilyID, "product-family-id", 0, "the product family id, i.e. 1804402")
	viper.BindPFlag("product-family-id", command.PersistentFlags().Lookup("product-family-id"))

	components.Init(command)
	products.Init(command)

}
func execute() error {
	apiKey := viper.Get("chargify-api-key")
	log.Debug().Str("api-key", apiKey.(string)).Send()

	if shared.ProductFamilyID == 0 {
		result, err := goChargify.GetProductFamilies()
		if err != nil {
			return err
		}
		fmt.Println(utils.PrettyJSON(result))
	} else {
		result, err := goChargify.GetProductFamily(shared.ProductFamilyID)
		if err != nil {
			return err
		}
		fmt.Println(utils.PrettyJSON(result))
	}
	return nil
}
