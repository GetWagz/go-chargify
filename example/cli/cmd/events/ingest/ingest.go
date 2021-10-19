package ingest

import (
	"context"
	"encoding/json"
	"fmt"

	goChargify "github.com/GetWagz/go-chargify"
	"github.com/GetWagz/go-chargify/example/cli/internal"
	"github.com/GetWagz/go-chargify/example/cli/internal/utils"
	"github.com/qri-io/jsonschema"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	jsonBodySchemaBulk = []byte(`{
		"$schema": "http://json-schema.org/draft-07/schema",
		"$id": "http://example.com/example.json",
		"type": "array",
		"title": "The root schema",
		"description": "The root schema comprises the entire JSON document.",
		"default": [],
		"examples": [
			[
				{
					"chargify": {
						"uniqueness_token": "1234"
					}
				}
			]
		],
		"additionalItems": true,
		"items": {
			"$id": "#/items",
			"anyOf": [
				{
					"$id": "#/items/anyOf/0",
					"type": "object",
					"title": "The first anyOf schema",
					"description": "An explanation about the purpose of this instance.",
					"default": {},
					"examples": [
						{
							"chargify": {
								"uniqueness_token": "1234"
							}
						}
					],
					"required": [
						"chargify"
					],
					"properties": {
						"chargify": {
							"$id": "#/items/anyOf/0/properties/chargify",
							"type": "object",
							"title": "The chargify schema",
							"description": "An explanation about the purpose of this instance.",
							"default": {},
							"examples": [
								{
									"uniqueness_token": "1234"
								}
							],
							"required": [
								"uniqueness_token"
							],
							"properties": {
								"uniqueness_token": {
									"$id": "#/items/anyOf/0/properties/chargify/properties/uniqueness_token",
									"type": "string",
									"title": "The uniqueness_token schema",
									"description": "An explanation about the purpose of this instance.",
									"default": "",
									"examples": [
										"1234"
									]
								}
							},
							"additionalProperties": true
						}
					},
					"additionalProperties": true
				}
			]
		}
	}`)
	jsonSchemaArray *jsonschema.Schema

	jsonBodySchema = []byte(`{
		"$schema": "http://json-schema.org/draft-07/schema",
		"$id": "http://example.com/example.json",
		"type": "object",
		"title": "The root schema",
		"description": "The root schema comprises the entire JSON document.",
		"default": {},
		"examples": [
			{
				"a": {},
				"chargify": {
					"uniqueness_token": "1234"
				}
			}
		],
		"required": [
			"chargify"
		],
		"properties": {
			"chargify": {
				"$id": "#/properties/chargify",
				"type": "object",
				"title": "The chargify schema",
				"description": "An explanation about the purpose of this instance.",
				"default": {},
				"examples": [
					{
						"uniqueness_token": "1234"
					}
				],
				"required": [
					"uniqueness_token"
				],
				"properties": {
					"uniqueness_token": {
						"$id": "#/properties/chargify/properties/uniqueness_token",
						"type": "string",
						"title": "The uniqueness_token schema",
						"description": "An explanation about the purpose of this instance.",
						"default": "",
						"examples": [
							"1234"
						]
					}
				},
				"additionalProperties": true
			}
		},
		"additionalProperties": true
	}`)
	jsonSchemaObject *jsonschema.Schema
)

func init() {
	rs := &jsonschema.Schema{}
	if err := json.Unmarshal(jsonBodySchema, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}
	jsonSchemaObject = rs

	rs = &jsonschema.Schema{}
	if err := json.Unmarshal(jsonBodySchemaBulk, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}
	jsonSchemaArray = rs
}

var command = &cobra.Command{
	Use:   "ingest",
	Short: "events ingest",
	Long: `
	example-single: events ingest --subdomain=mapped-3 --api-handle=ru_api_usage  --body='{\"customer\":{\"orgID\":\"2223\"},\"chargify\":{\"subscription_id\":\"46085700\",\"uniqueness_token\":\"1234\"}}'
	example-bulk:   events ingest --subdomain=mapped-3 --api-handle=ru_api_usage  --body='[{\"customer\":{\"orgID\":\"2223\"},\"chargify\":{\"subscription_id\":\"46085700\",\"uniqueness_token\":\"1234\"}}]'
	`,
	PersistentPreRunE: utils.ParentPersistentPreRunE,
	PreRunE:           preRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		return execute()
	},
}

type optionalEventsIngestQueryParams struct {
	StoreUID string `json:"store-uid,omitempty" mapstructure:"store-uid,omitempty"`
}

var body interface{}
var optionalQueryParams = optionalEventsIngestQueryParams{}
var queryParams = goChargify.EventsIngestQueryParams{}
var bulkIngest bool

func preRunE(cmd *cobra.Command, args []string) error {
	if optionalQueryParams.StoreUID != internal.NotSetStringParam {
		queryParams.StoreUID = &optionalQueryParams.StoreUID
	}

	var valid = []byte(jsonBody)

	// check single first
	var singleError error
	errs, singleError := jsonSchemaObject.ValidateBytes(context.TODO(), valid)
	if len(errs) > 0 {
		singleError = errs[0]
	}
	if singleError != nil {
		// maybe bulk error will work
		var bulkError error
		errs, bulkError := jsonSchemaArray.ValidateBytes(context.TODO(), valid)
		if len(errs) > 0 {
			bulkError = errs[0]
		}
		if bulkError != nil {
			log.Error().Err(bulkError).Msg("bulk Error")
			log.Error().Err(singleError).Msg("single Error")
			return singleError // this was the first error
		} else {
			jsonArrayMap := make([]map[string]interface{}, 0)

			err := json.Unmarshal(valid, &jsonArrayMap)
			if err != nil {
				log.Error().Err(err).Msg("failed to unmarshal json")
				return err
			}
			body = jsonArrayMap
			bulkIngest = true
		}
	} else {
		jsonMap := make(map[string](interface{}))
		err := json.Unmarshal(valid, &jsonMap)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal json")
			return err
		}
		body = jsonMap
	}

	return nil
}

var jsonBody string
var apiHandle string

func Init(rootCmd *cobra.Command) {
	// store-uid
	command.Flags().StringVar(&optionalQueryParams.StoreUID, "store-uid", internal.NotSetStringParam, "[optional] store-uid")

	// api-handle
	command.Flags().StringVar(&apiHandle, "api-handle", internal.NotSetStringParam, "[required] api-handle")
	command.MarkFlagRequired("api-handle")
	// body
	command.Flags().StringVar(&jsonBody, "body", internal.NotSetStringParam, "[required] body arbitrary json")
	command.MarkFlagRequired("body")

	rootCmd.AddCommand(command)
}

func execute() error {
	apiKey := viper.Get("chargify-api-key")
	log.Debug().Str("api-key", apiKey.(string)).Send()
	log.Debug().Interface("body", body).Interface("queryParams", queryParams).Send()

	pathParams := map[string]string{}
	pathParams["api_handle"] = apiHandle

	var err error

	if bulkIngest {
		err = goChargify.PostBulkEventsIngestion(body, &pathParams, &queryParams)

	} else {
		err = goChargify.PostEventsIngestion(body, &pathParams, &queryParams)
	}

	if err != nil {
		return err
	}

	fmt.Println("success...")
	return nil

}
