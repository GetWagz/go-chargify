package chargify

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/GetWagz/go-chargify/internal"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

// Subscription represents a subscription
type Subscription struct {
	ID                            int64     `json:"id"`
	CancellationMessage           string    `json:"cancellation_message" mapstructure:"cancellation_message"`                                   // (Optional) Can be used when canceling a subscription (via the HTTP DELETE method) to make a note about the reason for cancellation.
	CancellationMethod            string    `json:"cancellation_method" mapstructure:"cancellation_method"`                                     // (Optional) Can be used when canceling a subscription (via the HTTP DELETE method) to make a note about how the subscription was canceled.
	ReasonCode                    string    `json:"reason_code" mapstructure:"reason_code"`                                                     // (Optional) Can be used when canceling a subscription (via the HTTP DELETE method) to indicate why a subscription was canceled.
	NextBillingAt                 string    `json:"next_billing_at" mapstructure:"next_billing_at"`                                             // (Optional) Set this attribute to a future date/time to sync imported subscriptions to your existing renewal schedule. See the notes on “Date/Time Format” at https://help.chargify.com/subscriptions/subscriptions-import.html. If you provide a next_billing_at timestamp that is in the future, no trial or initial charges will be applied when you create the subscription. In fact, no payment will be captured at all. The first payment will be captured, according to the prices defined by the product, near the time specified by next_billing_at. If you do not provide a value for next_billing_at, any trial and/or initial charges will be assessed and charged at the time of subscription creation. If the card cannot be successfully charged, the subscription will not be created. See further notes in the section on Importing Subscriptions.
	ExpiresAt                     string    `json:"expires_at" mapstructure:"expires_at"`                                                       // Timestamp giving the expiration date of this subscription (if any). You may manually change the expiration date at any point during a subscription period.
	ExpirationTracksChange        bool      `json:"expiration_tracks_next_billing_change" mapstructure:"expiration_tracks_next_billing_change"` // (Optional, default false) When set to true, and when next_billing_at is present, if the subscription expires, the expires_at will be shifted by the same amount of time as the difference between the old and new “next billing” dates.
	VATNumber                     string    `json:"vat_number" mapstructure:"vat_number"`                                                       // (Optional) Supplying the VAT number allows EU customer’s to opt-out of the Value Added Tax assuming the merchant address and customer billing address are not within the same EU country. It’s important to omit the country code from the VAT number upon entry. Otherwise, taxes will be assessed upon the purchase.
	CouponCode                    string    `json:"coupon_code" mapstructure:"coupon_code"`                                                     // (Optional) The coupon code of the coupon to apply ()
	PaymentCollectionMethod       string    `json:"payment_collection_method" mapstructure:"payment_collection_method"`                         // (Optional) The type of payment collection to be used in the subscription. May be automatic, or invoice.
	AgreementTerms                string    `json:"agreement_terms" mapstructure:"agreement_terms"`                                             // (Optional) The ACH authorization agreement terms. If enabled, an email will be sent to the customer with a copy of the terms.
	ACHFirstName                  string    `json:"authorizer_first_name" mapstructure:"authorizer_first_name"`                                 // (Optional) The first name of the person authorizing the ACH agreement.
	ACHLastName                   string    `json:"authorizer_last_name" mapstructure:"authorizer_last_name"`                                   // (Optional) The last name of the person authorizing the ACH agreement.
	ChangeDelayed                 bool      `json:"product_change_delayed" mapstructure:"product_change_delayed"`                               // (Optional, used only for https://reference.chargify.com/v1/subscriptions-product-changes-migrations-upgrades-downgrades/update-subscription-product-change When set to true, indicates that a changed value for product_handle should schedule the product change to the next subscription renewal.
	CalendarBilling               string    `json:"calendar_billing" mapstructure:"calendar_billing"`                                           // (Optional, see https://reference.chargify.com/v1/subscriptions/subscriptions-intro#https://help.chargify.com/subscriptions/billing-dates.html#calendar-billing for more details). Cannot be used when also specifying next_billing_at
	SnapDay                       int       `json:"snap_day" mapstructure:"snap_day"`                                                           // A value between 1 and 28, or end
	CalendarBillingFirstDayCharge string    `json:"calendar_billing_first_charge" mapstructure:"calendar_billing_first_charge"`                 // (Optional) One of “prorated” (the default – the prorated product price will be charged immediately), “immediate” (the full product price will be charged immediately), or “delayed” (the full product price will be charged with the first scheduled renewal).
	ReceivesInvoiceEmails         bool      `json:"receives_invoice_emails" mapstructure:"receives_invoice_emails"`                             // (Optional) Default: True - Whether or not this subscription is set to receive emails related to this subscription.
	Customer                      *Customer `json:"customer,omitempty" mapstructure:"customer"`
	Product                       *Product  `json:"product,omitempty" mapstructure:"product"`
	// some of these are only used on the return
	State string `json:"state,omitempty" mapstructure:"state"` // the state of the subscription
}

type SubscriptionComponent struct {
	ComponentID       int64  `json:"component_id" mapstructure:"component_id"`
	SubscriptionID    int64  `json:"subscription_id" mapstructure:"subscription_id"`
	AllocatedQuantity int64  `json:"allocated_quantity" mapstructure:"allocated_quantity"`
	PricingScheme     string `json:"pricing_scheme" mapstructure:"pricing_scheme"`
	Name              string `json:"name" mapstructure:"name"`
	Kind              string `json:"kind" mapstructure:"kind"`
	UnitName          string `json:"unit_name" mapstructure:"unit_name"`
	PricePointID      int64  `json:"price_point_id" mapstructure:"price_point_id"`
	PricePointHandle  string `json:"price_point_handle" mapstructure:"price_point_handle"`
	PricePointType    string `json:"price_point_type" mapstructure:"price_point_type"`
	PricePointName    string `json:"price_point_type" mapstructure:"price_point_type"`
	Enabled           bool   `json:"enabled" mapstructure:"enabled"`
	UnitBalance       int64  `json:"unit_balance" mapstructure:"unit_balance"`
	ID                int64  `json:"id" mapstructure:"id"`
	CreatedAt         string `json:"created_at" mapstructure:"created_at"`
	UpdatedAt         string `json:"updated_at" mapstructure:"updated_at"`
	ComponentHandle   string `json:"component_handle" mapstructure:"component_handle"`
	ArchivedAt        string `json:"archived_at" mapstructure:"archived_at"`
}

type ListSubscriptionEventsQueryParams struct {
	Page      *int    `json:"page,omitempty" mapstructure:"page,omitempty"`
	PerPage   *int    `json:"per_page,omitempty" mapstructure:"per_page,omitempty"`
	SinceID   *int    `json:"since_id,omitempty" mapstructure:"since_id,omitempty"`
	MaxID     *int    `json:"max_id,omitempty" mapstructure:"max_id,omitempty"`
	Direction *string `json:"direction,omitempty" mapstructure:"direction,omitempty"`
	Filter    *string `json:"filter,omitempty" mapstructure:"filter,omitempty"`
}

// ListSubscriptionsQueryParams are the query parameters for listing subscriptions; see:
// https://developers.chargify.com/docs/api-docs/51c68dd4dcb2b-list-subscriptions
type ListSubscriptionsQueryParams struct {
	Coupon              *int64  `json:"coupon"`
	DateField           *string `json:"date_field"`
	Direction           *string `json:"direction"`
	EndDate             *string `json:"end_date"`
	EndDateTime         *string `json:"end_datetime"`
	Page                *int    `json:"page"`
	PerPage             *int    `json:"per_page"`
	Product             *string `json:"product"`
	ProductPricePointID *int64  `json:"product_price_point_id"`
	Sort                *string `json:"sort"`
	StartDate           *string `json:"start_date"`
	StartDateTime       *string `json:"start_datetime"`
	State               *string `json:"state"`
}

func (input *ListSubscriptionsQueryParams) toMap() *map[string]string {
	m := map[string]string{}
	if input.Coupon != nil {
		m["coupon"] = fmt.Sprintf("%d", ToInt64(input.Coupon))
	}
	if input.DateField != nil {
		m["date_field"] = ToString(input.DateField)
	}
	if input.Direction != nil {
		m["direction"] = ToString(input.Direction)
	}
	if input.EndDate != nil {
		m["end_date"] = ToString(input.EndDate)
	}
	if input.EndDateTime != nil {
		m["end_datetime"] = ToString(input.EndDateTime)
	}
	if input.Page != nil {
		m["page"] = fmt.Sprintf("%d", ToInt(input.Page))
	}
	if input.PerPage != nil {
		m["per_page"] = fmt.Sprintf("%d", ToInt(input.PerPage))
	}
	if input.Product != nil {
		m["product"] = ToString(input.Product)
	}
	if input.ProductPricePointID != nil {
		m["product_price_point_id"] = fmt.Sprintf("%d", ToInt64(input.ProductPricePointID))
	}
	if input.Sort != nil {
		m["sort"] = ToString(input.Sort)
	}
	if input.StartDate != nil {
		m["start_date"] = ToString(input.StartDate)
	}
	if input.StartDateTime != nil {
		m["start_datetime"] = ToString(input.StartDateTime)
	}
	if input.State != nil {
		m["state"] = ToString(input.State)
	}
	return &m
}

// CreateSubscriptionForCustomer creates a new subscription. When creating a subscription, you must specify a product and a customer.
// The product should be specificed by productHandle and the customer should be specified with customerReference. The subscriptionOptions
// pointer is useful for specifying select additional options. Right now, only NextChargeAt is supported.
// The paymentProfileID is optional and is used to associate the subscription with a payment profile. If one is already setup,
// pass in 0.
func CreateSubscriptionForCustomer(customerReference, productHandle string, paymentProfileID int64, subscriptionOptions *Subscription) (*Subscription, error) {
	body := map[string]map[string]interface{}{
		"subscription": {
			"customer_reference": customerReference,
			"product_handle":     productHandle,
		},
	}
	if paymentProfileID != 0 {
		body["subscription"]["payment_profile_id"] = paymentProfileID
	}
	if subscriptionOptions != nil {
		if subscriptionOptions.NextBillingAt != "" {
			body["subscription"]["next_billing_at"] = subscriptionOptions.NextBillingAt
		}
		if subscriptionOptions.CouponCode != "" {
			body["subscription"]["coupon_code"] = subscriptionOptions.CouponCode
		}
	}

	ret, err := makeCall(endpoints[endpointSubscriptionCreate], body, nil)
	if err != nil {
		return nil, err
	}
	// if successful, the subscription should come back in a map["subscription"]Subscription format
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	subscription := &Subscription{}
	err = mapstructure.Decode(apiBody["subscription"], subscription)
	return subscription, err
}

// CancelSubscription cancels a subscription. You can choose to cancel now or delay it. If you choose to delay, you can provide a reason code and message
func CancelSubscription(subscriptionID int64, cancelImmediately bool, reasonCode string, cancellationMessage string) error {
	var err error
	if cancelImmediately {
		// it is a delete so no body
		_, err = makeCall(endpoints[endpointSubscriptionCancelImmediately], nil, &map[string]string{
			"subscriptionID": fmt.Sprintf("%d", subscriptionID),
		})
	} else {
		// if the reason info is set, then add it
		reason := map[string]string{}
		if reasonCode != "" && cancellationMessage != "" {
			reason = map[string]string{
				"cancellation_message": cancellationMessage,
				"reason_code":          reasonCode,
			}
		}
		_, err = makeCall(endpoints[endpointSubscriptionCancelDelayed], reason, &map[string]string{
			"subscriptionID": fmt.Sprintf("%d", subscriptionID),
		})
	}
	return err
}

// UpdateSubscription updates a subscription for a customer
func UpdateSubscription(subscriptionID int64, productHandle string) error {
	body := map[string]map[string]interface{}{
		"subscription": {
			"product_handle": productHandle,
		},
	}
	_, err := makeCall(endpoints[endpointSubscriptionUpdate], body, &map[string]string{
		"subscriptionID": fmt.Sprintf("%d", subscriptionID),
	})
	return err
}

// RemoveDelayedSubscriptionCancellation removes a delayed cancellation request, ensuring the subscription does not cancel
func RemoveDelayedSubscriptionCancellation(subscriptionID int64) error {
	_, err := makeCall(endpoints[endpointSubscriptionRemoveDelayedCancel], nil, &map[string]string{
		"subscriptionID": fmt.Sprintf("%d", subscriptionID),
	})
	return err
}

// MigrateSubscription migrates an existing subscription to a new subscription
func MigrateSubscription(targetProductHandle string, currentSubscriptionID int64, includeTrial bool, includeInitialCharge bool, includeCoupons bool, preservePeriod bool) error {
	body := map[string]map[string]interface{}{
		"migration": {
			"product_handle":         targetProductHandle,
			"include_trial":          includeTrial,
			"include_initial_charge": includeInitialCharge,
			"include_coupons":        includeCoupons,
			"preserve_period":        preservePeriod,
		},
	}

	_, err := makeCall(endpoints[endpointSubscriptionMigrate], body, &map[string]string{
		"subscriptionID": fmt.Sprintf("%d", currentSubscriptionID),
	})
	return err
}

// GetSubscription gets a subscription. The docs show it comes back as an array, but as of this implementation it comes back as a map
func GetSubscription(subscriptionID int64) (*Subscription, error) {
	ret, err := makeCall(endpoints[endpointSubscriptionGet], nil, &map[string]string{
		"subscriptionID": fmt.Sprintf("%d", subscriptionID),
	})
	if err != nil {
		return nil, err
	}
	// if successful, the subscription should come back in a map["subscription"]Subscription format
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	subscription := &Subscription{}
	err = mapstructure.Decode(apiBody["subscription"], subscription)
	return subscription, err
}

// GetProductFamilyProducts gets products in a family
func GetSubscriptionComponents(subscriptionID int64) ([]SubscriptionComponent, error) {
	found := []SubscriptionComponent{}

	ret, err := makeCall(endpoints[endpointProductFamilyComponentsGet], nil, &map[string]string{
		"subscriptionID": fmt.Sprintf("%d", subscriptionID),
	})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return found, err
	}

	temp := ret.Body.([]interface{})
	for i := range temp {
		entity := SubscriptionComponent{}
		entry := temp[i].(map[string]interface{})
		entityRaw := entry["component"]
		err = mapstructure.Decode(entityRaw, &entity)
		if err != nil {
			return []SubscriptionComponent{}, err
		}
		found = append(found, entity)
	}

	return found, err
}

// GetSubscriptionMetaData gets the subscription metadata
func GetSubscriptionMetaData(subscriptionID int64) (*MetaData, error) {
	ret, err := makeCall(endpoints[endpointSubscriptionGetMetaData], nil, &map[string]string{
		"subscriptionID": fmt.Sprintf("%d", subscriptionID),
	})
	if err != nil {
		return nil, err
	}
	// if successful, the subscription should come back in a map
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	data := &MetaData{}
	err = mapstructure.Decode(apiBody, data)
	return data, err
}

// RefundSubscriptionPayment refunds a specific payment for a subscription. This is supposedly deprecated to support relationship
// invoicing
func RefundSubscriptionPayment(subscriptionID string, paymentID string, amount string, memo string) (*Refund, error) {
	body := map[string]map[string]string{
		"refund": {
			"payment_id": paymentID,
			"amount":     amount,
			"memo":       memo,
		},
	}

	ret, err := makeCall(endpoints[endpointSubscriptionRefund], body, &map[string]string{
		"subscriptionID": subscriptionID,
	})
	if err != nil {
		return nil, err
	}
	// if successful, the subscription should come back in a map["refund"] format
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	refund := &Refund{}
	err = mapstructure.Decode(apiBody["refund"], refund)
	return refund, err
}

// GetCustomerByID gets a customer by chargify id
func ListSubscriptionEvents(subscriptionID int, queryParams *ListSubscriptionEventsQueryParams) (found []Event, err error) {
	structs.DefaultTagName = "mapstructure"
	m := structs.Map(queryParams)
	body := internal.ToMapStringToString(m)
	ret, err := makeCall(endpoints[endpointSubscriptionEvents], body, &map[string]string{
		"subscriptionID": fmt.Sprintf("%d", subscriptionID),
	})
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

// PurgeSubscription purges a subscription from an account IN TEST MODE. This WILL NOT WORK on production environments.
func PurgeSubscription(subscriptionID int64, customerID int64, cascadeCustomer bool, cascadePayment bool) error {
	cascade := []string{}
	if cascadeCustomer {
		cascade = append(cascade, "customer")
	}
	if cascadePayment {
		cascade = append(cascade, "payment_profile")
	}

	// note that we will eventually want to move the other calls to this cleaner approach at
	// some point in the future
	options := &makeCallOptions{
		End: endpoints[endpointSubscriptionPurge],
		MultiQueryParams: &map[string][]string{
			"ack":     {fmt.Sprintf("%d", customerID)},
			"cascade": cascade,
		},
		PathParams: &map[string]string{
			"subscriptionID": fmt.Sprintf("%d", subscriptionID),
		},
	}

	ret, err := makeAPICall(options)
	if err != nil {
		return err
	}

	// if successful, we can just return out from here
	if ret.HTTPCode != http.StatusOK {
		return fmt.Errorf("received a %d", ret.HTTPCode)
	}
	return nil
}

// ListSubscriptions lists out the subscriptions based upon the result of the passed in query params
func ListSubscriptions(params *ListSubscriptionsQueryParams) ([]Subscription, error) {
	if params == nil {
		params = &ListSubscriptionsQueryParams{}
	}
	queryParams := params.toMap()
	options := &makeCallOptions{
		End:         endpoints[endpointSubscriptionsList],
		QueryParams: queryParams,
	}

	data := []Subscription{}

	ret, err := makeAPICall(options)
	if err != nil {
		return data, err
	}
	apiBody, bodyOK := ret.Body.([]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	// the result is an array of objects that have a subscription key, similar to:
	// 	[
	//   {
	//     "subscription": {...}
	//   }
	// ]
	// so we have a [] of interfaces; loop and convert
	for i := range apiBody {
		if raw, rawOK := apiBody[i].(map[string]interface{}); rawOK {
			sub := Subscription{}
			err = mapstructure.Decode(raw["subscription"], &sub)
			if err == nil {
				data = append(data, sub)
			}
		}

	}
	return data, nil
}
