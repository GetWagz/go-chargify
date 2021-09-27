package count

import (
	"fmt"

	goChargify "github.com/GetWagz/go-chargify"
	"github.com/GetWagz/go-chargify/example/cli/cmd/subscriptions/shared"
	"github.com/GetWagz/go-chargify/example/cli/internal"
	"github.com/GetWagz/go-chargify/example/cli/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var command = &cobra.Command{
	Use:               "count",
	Short:             "events count",
	PersistentPreRunE: utils.ParentPersistentPreRunE,
	PreRunE:           preRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		return execute()
	},
}

type optionalListEventsCountQueryParams struct {
	Page      int    `json:"page,omitempty" mapstructure:"page,omitempty"`
	PerPage   int    `json:"per_page,omitempty" mapstructure:"per_page,omitempty"`
	SinceID   int    `json:"since_id,omitempty" mapstructure:"since_id,omitempty"`
	MaxID     int    `json:"max_id,omitempty" mapstructure:"max_id,omitempty"`
	Direction string `json:"direction,omitempty" mapstructure:"direction,omitempty"`
	Filter    string `json:"filter,omitempty" mapstructure:"filter,omitempty"`
}

var optionalQueryParams = optionalListEventsCountQueryParams{}
var listEventQueryParams = goChargify.ListEventsCountQueryParams{}

func preRunE(cmd *cobra.Command, args []string) error {
	if optionalQueryParams.Page != internal.NotSetIntParam {
		listEventQueryParams.Page = &optionalQueryParams.Page
	}
	if optionalQueryParams.PerPage != internal.NotSetIntParam {
		listEventQueryParams.PerPage = &optionalQueryParams.PerPage
	}
	if optionalQueryParams.SinceID != internal.NotSetIntParam {
		listEventQueryParams.SinceID = &optionalQueryParams.SinceID
	}
	if optionalQueryParams.MaxID != internal.NotSetIntParam {
		listEventQueryParams.MaxID = &optionalQueryParams.MaxID
	}
	if optionalQueryParams.Direction != internal.NotSetStringParam {
		if !(optionalQueryParams.Direction == "asc" || optionalQueryParams.Direction == "desc") {
			return fmt.Errorf("--direction must be asc or desc")
		}
		listEventQueryParams.Direction = &optionalQueryParams.Direction
	}
	if optionalQueryParams.Filter != internal.NotSetStringParam {
		listEventQueryParams.Filter = &optionalQueryParams.Filter
	}

	return nil
}
func Init(rootCmd *cobra.Command) {
	// page
	command.Flags().IntVar(&optionalQueryParams.Page, "page", internal.NotSetIntParam, "page, i.e. 1 or above")
	// per-page
	command.Flags().IntVar(&optionalQueryParams.PerPage, "per-page", internal.NotSetIntParam, "per-page, i.e. 1 or above")
	// since_id
	command.Flags().IntVar(&optionalQueryParams.SinceID, "since_id", internal.NotSetIntParam, "since_id, i.e. 1 or above")
	// max_id
	command.Flags().IntVar(&optionalQueryParams.MaxID, "max_id", internal.NotSetIntParam, "max_id, i.e. 1 or above")
	// direction
	command.Flags().StringVar(&optionalQueryParams.Direction, "direction", internal.NotSetStringParam, "asc,desc")
	// filter
	command.Flags().StringVar(&optionalQueryParams.Filter, "filter", internal.NotSetStringParam, "filter")

	rootCmd.AddCommand(command)
}

func execute() error {
	apiKey := viper.Get("chargify-api-key")
	log.Debug().Str("api-key", apiKey.(string)).Send()
	log.Debug().Int("", shared.SubscriptionID).Interface("listEventQueryParams", listEventQueryParams).Send()
	found, err := goChargify.GetEventsCount(&listEventQueryParams)
	if err != nil {
		return err
	}
	fmt.Println(utils.PrettyJSON(found))
	return nil

}
