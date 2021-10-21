package components

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
	Use:               "components",
	Short:             "query components",
	PersistentPreRunE: persistentPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		return execute()
	},
}
var componentID int64
var handle string

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
	command.Flags().Int64Var(&componentID, "component-id", 0, "the component id, i.e. 1804402")
	// handle
	command.Flags().StringVar(&handle, "handle", "", "the component handle, i.e. api-usage-credits")

}
func execute() error {

	apiKey := viper.Get("chargify-api-key")

	log.Debug().
		Str("api-key", apiKey.(string)).
		Int64("productFamilyID", shared.ProductFamilyID).
		Int64("componentID", componentID).
		Str("handle", handle).
		Send()

	if handle == "" && componentID == 0 {
		result, err := goChargify.GetProductFamilyComponents(shared.ProductFamilyID)
		if err != nil {
			return err
		}
		fmt.Println(utils.PrettyJSON(result))
		return nil
	}
	if componentID > 0 {
		result, err := goChargify.GetProductFamilyComponentById(shared.ProductFamilyID, componentID)
		if err != nil {
			return err
		}
		fmt.Println(utils.PrettyJSON(result))
		return nil
	}
	if len(handle) > 0 {
		result, err := goChargify.GetProductFamilyComponentByHandle(shared.ProductFamilyID, handle)
		if err != nil {
			return err
		}
		fmt.Println(utils.PrettyJSON(result))
		return nil
	}
	return nil

}
