package products

import (
	"fmt"

	"github.com/GetWagz/go-chargify/example/cli/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var command = &cobra.Command{
	Use:               "products",
	Short:             "query products",
	PersistentPreRunE: utils.ParentPersistentPreRunE,
	Run: func(cmd *cobra.Command, args []string) {
		execute()
	},
}
var productID string

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(command)

	// number-of-consumers
	command.PersistentFlags().StringVar(&productID, "id", "", "the product id, i.e. 1804402")
	command.MarkPersistentFlagRequired("id")

}
func execute() {
	apiKey := viper.Get("chargify-api-key")
	log.Debug().Str("api-key", apiKey.(string)).Send()

	fmt.Printf("products/%v", productID)

}
