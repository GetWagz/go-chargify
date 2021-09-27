package events

import (
	"fmt"

	goChargify "github.com/GetWagz/go-chargify"
	"github.com/GetWagz/go-chargify/example/cli/cmd/events/count"
	"github.com/GetWagz/go-chargify/example/cli/cmd/events/ingest"
	"github.com/GetWagz/go-chargify/example/cli/internal"
	"github.com/GetWagz/go-chargify/example/cli/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var command = &cobra.Command{
	Use:               "events",
	Short:             "events streaming",
	PersistentPreRunE: utils.ParentPersistentPreRunE,
	PreRunE:           preRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		return execute()
	},
}

type optionalListEventsQueryParams struct {
	Page          int    `json:"page,omitempty" mapstructure:"page,omitempty"`
	PerPage       int    `json:"per_page,omitempty" mapstructure:"per_page,omitempty"`
	SinceID       int    `json:"since_id,omitempty" mapstructure:"since_id,omitempty"`
	MaxID         int    `json:"max_id,omitempty" mapstructure:"max_id,omitempty"`
	Direction     string `json:"direction,omitempty" mapstructure:"direction,omitempty"`
	Filter        string `json:"filter,omitempty" mapstructure:"filter,omitempty"`
	DateField     string `json:"date_field,omitempty" mapstructure:"date_field,omitempty"`
	StartDate     string `json:"start_date,omitempty" mapstructure:"start_date,omitempty"`
	EndDate       string `json:"end_date,omitempty" mapstructure:"end_date,omitempty"`
	StartDatetime string `json:"start_datetime,omitempty" mapstructure:"start_datetime,omitempty"`
	EndDatetime   string `json:"end_datetime,omitempty" mapstructure:"end_datetime,omitempty"`
}

var optionalQueryParams optionalListEventsQueryParams = optionalListEventsQueryParams{}
var listEventQueryParams = goChargify.ListEventsQueryParams{}

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
	if optionalQueryParams.DateField != internal.NotSetStringParam {
		listEventQueryParams.DateField = &optionalQueryParams.DateField
	}
	if optionalQueryParams.StartDate != internal.NotSetStringParam {
		listEventQueryParams.StartDate = &optionalQueryParams.StartDate
	}
	if optionalQueryParams.EndDate != internal.NotSetStringParam {
		listEventQueryParams.EndDate = &optionalQueryParams.EndDate
	}
	if optionalQueryParams.StartDatetime != internal.NotSetStringParam {
		listEventQueryParams.StartDatetime = &optionalQueryParams.StartDatetime
	}
	if optionalQueryParams.EndDatetime != internal.NotSetStringParam {
		listEventQueryParams.EndDatetime = &optionalQueryParams.EndDatetime
	}
	return nil
}
func Init(rootCmd *cobra.Command) {
	// page
	command.Flags().IntVar(&optionalQueryParams.Page, "page", internal.NotSetIntParam, "[optional] page, i.e. 1 or above")
	// per-page
	command.Flags().IntVar(&optionalQueryParams.PerPage, "per-page", internal.NotSetIntParam, "[optional] per-page, i.e. 1 or above")
	// since_id
	command.Flags().IntVar(&optionalQueryParams.SinceID, "since_id", internal.NotSetIntParam, "[optional] since_id, i.e. 1 or above")
	// max_id
	command.Flags().IntVar(&optionalQueryParams.MaxID, "max_id", internal.NotSetIntParam, "[optional] max_id, i.e. 1 or above")
	// direction
	command.Flags().StringVar(&optionalQueryParams.Direction, "direction", internal.NotSetStringParam, "[optional] asc,desc")
	// filter
	command.Flags().StringVar(&optionalQueryParams.Filter, "filter", internal.NotSetStringParam, "[optional] filter")
	// date_field
	command.Flags().StringVar(&optionalQueryParams.DateField, "date_field", internal.NotSetStringParam, "[optional] date_field")
	// start_date
	command.Flags().StringVar(&optionalQueryParams.StartDate, "start_date", internal.NotSetStringParam, "[optional] start_date")
	// end_date
	command.Flags().StringVar(&optionalQueryParams.EndDate, "end_date", internal.NotSetStringParam, "[optional] end_date")
	// start_datetime
	command.Flags().StringVar(&optionalQueryParams.StartDatetime, "start_datetime", internal.NotSetStringParam, "[optional] start_datetime")
	// end_datetime
	command.Flags().StringVar(&optionalQueryParams.EndDatetime, "end_datetime", internal.NotSetStringParam, "[optional] end_datetime")

	count.Init(command)
	ingest.Init(command)
	rootCmd.AddCommand(command)
}

//https://help.chargify.com/events/getting-data-in-guide.html
/*
curl https://events.chargify.com/mapped-3/events/ru_api_usage \
-u API_KEY:x \
-H 'Content-Type: application/json' \
-d '{
  "chargify": {
    "subscription_id": 123456
  },
  "clicks": 100
}'
*/

func execute() error {
	apiKey := viper.Get("chargify-api-key")
	log.Debug().Str("api-key", apiKey.(string)).Send()
	log.Debug().Interface("listEventQueryParams", listEventQueryParams).Send()
	found, err := goChargify.ListEvents(&listEventQueryParams)
	if err != nil {
		return err
	}
	fmt.Println(utils.PrettyJSON(found))
	return nil
}
