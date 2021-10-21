package chargify

import (
	"net/http"

	"github.com/GetWagz/go-chargify/internal"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type AllocationDetail struct {
	AllocationID  int `json:"allocation_id" mapstructure:"allocation_id"`
	ChargeID      int `json:"charge_id" mapstructure:"charge_id"`
	UsageQuantity int `json:"usage_quantity" mapstructure:"usage_quantity"`
}
type EventSpecificData struct {
	PreviousUnitBalance        string             `json:"previous_unit_balance" mapstructure:"previous_unit_balance"`
	PreviousOverageUnitBalance string             `json:"previous_overage_unit_balance" mapstructure:"previous_overage_unit_balance"`
	NewUnitBalance             int                `json:"new_unit_balance" mapstructure:"new_unit_balance"`
	NewOverageUnitBalance      int                `json:"new_overage_unit_balance" mapstructure:"new_overage_unit_balance"`
	UsageQuantity              int                `json:"usage_quantity" mapstructure:"usage_quantity"`
	OverageUsageQuantity       int                `json:"overage_usage_quantity" mapstructure:"overage_usage_quantity"`
	ComponentID                int                `json:"component_id" mapstructure:"component_id"`
	ComponentHandle            string             `json:"component_handle" mapstructure:"component_handle"`
	Memo                       string             `json:"memo" mapstructure:"memo"`
	AllocationDetails          []AllocationDetail `json:"allocation_details" mapstructure:"allocation_details"`
}
type Event struct {
	ID                int64             `json:"id" mapstructure:"id"` //	The customer ID in Chargify
	Key               string            `json:"key" mapstructure:"key"`
	Message           string            `json:"message" mapstructure:"message"`
	Subscription_id   int               `json:"subscription_id" mapstructure:"subscription_id"`
	CustomerID        int               `json:"customer_id" mapstructure:"customer_id"`
	CreatedAt         string            `json:"created_at" mapstructure:"created_at"`
	EventSpecificData EventSpecificData `json:"event_specific_data" mapstructure:"event_specific_data"`
}

type ListEventsQueryParams struct {
	Page          *int    `json:"page,omitempty" mapstructure:"page,omitempty"`
	PerPage       *int    `json:"per_page,omitempty" mapstructure:"per_page,omitempty"`
	SinceID       *int    `json:"since_id,omitempty" mapstructure:"since_id,omitempty"`
	MaxID         *int    `json:"max_id,omitempty" mapstructure:"max_id,omitempty"`
	Direction     *string `json:"direction,omitempty" mapstructure:"direction,omitempty"`
	Filter        *string `json:"filter,omitempty" mapstructure:"filter,omitempty"`
	DateField     *string `json:"date_field,omitempty" mapstructure:"date_field,omitempty"`
	StartDate     *string `json:"start_date,omitempty" mapstructure:"start_date,omitempty"`
	EndDate       *string `json:"end_date,omitempty" mapstructure:"end_date,omitempty"`
	StartDatetime *string `json:"start_datetime,omitempty" mapstructure:"start_datetime,omitempty"`
	EndDatetime   *string `json:"end_datetime,omitempty" mapstructure:"end_datetime,omitempty"`
}
type ListEventsCountQueryParams struct {
	Page      *int    `json:"page,omitempty" mapstructure:"page,omitempty"`
	PerPage   *int    `json:"per_page,omitempty" mapstructure:"per_page,omitempty"`
	SinceID   *int    `json:"since_id,omitempty" mapstructure:"since_id,omitempty"`
	MaxID     *int    `json:"max_id,omitempty" mapstructure:"max_id,omitempty"`
	Direction *string `json:"direction,omitempty" mapstructure:"direction,omitempty"`
	Filter    *string `json:"filter,omitempty" mapstructure:"filter,omitempty"`
}

type EventsIngestQueryParams struct {
	StoreUID *string `json:"store_uid,omitempty" mapstructure:"store_uid,omitempty"`
}
type Count struct {
	Count int `json:"count,omitempty" mapstructure:"count,omitempty"`
}

// GetCustomerByID gets a customer by chargify id
func ListEvents(queryParams *ListEventsQueryParams) (found []Event, err error) {
	structs.DefaultTagName = "mapstructure"
	m := structs.Map(queryParams)
	body := internal.ToMapStringToString(m)
	ret, err := makeCall(endpoints[endpointEvents], body, &map[string]string{})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return nil, err
	}

	temp := ret.Body.([]interface{})
	for i := range temp {
		entry := temp[i].(map[string]interface{})
		raw := entry["event"]
		entity := Event{}
		err = mapstructure.Decode(raw, &entity)
		if err == nil {
			found = append(found, entity)
		}
	}
	return found, nil

}

// GetEventsCount ...
func GetEventsCount(queryParams *ListEventsCountQueryParams) (response *Count, err error) {
	structs.DefaultTagName = "mapstructure"
	m := structs.Map(queryParams)
	body := internal.ToMapStringToString(m)
	ret, err := makeCall(endpoints[endpointEventsCount], body, &map[string]string{})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return nil, err
	}

	raw := ret.Body
	response = &Count{}
	err = mapstructure.Decode(raw, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// PostEventsInjestion ...
func PostEventsIngestion(body interface{}, pathParams *map[string]string, queryParams *EventsIngestQueryParams) error {
	var qP *map[string]string
	if queryParams != nil {
		structs.DefaultTagName = "mapstructure"
		m := structs.Map(queryParams)
		m2 := internal.ToMapStringToString(m)
		qP = &m2
	}
	ret, err := makeEventsCall(endpoints[endpointEventIngestion], body, pathParams, qP)
	if err != nil || ret.HTTPCode != http.StatusOK {
		return err
	}
	return err
}

// PostBulkEventsIngestion ...
func PostBulkEventsIngestion(body interface{}, pathParams *map[string]string, queryParams *EventsIngestQueryParams) error {
	var qP *map[string]string
	if queryParams != nil {
		structs.DefaultTagName = "mapstructure"
		m := structs.Map(queryParams)
		m2 := internal.ToMapStringToString(m)
		qP = &m2
	}
	ret, err := makeEventsCall(endpoints[endpointBulkEventIngestion], body, pathParams, qP)
	if err != nil || ret.HTTPCode != http.StatusOK {
		return err
	}
	return err
}
