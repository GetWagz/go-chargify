package chargify

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Subscription represents a subscription
type Subscription struct {
	ID                            int64     `json:"id"`
	CancellationMessage           string    `json:"cancellation_message,omitEmpty"`                  // (Optional) Can be used when canceling a subscription (via the HTTP DELETE method) to make a note about the reason for cancellation.
	CancellationMethod            string    `json:"cancellation_method,omitEmpty"`                   // (Optional) Can be used when canceling a subscription (via the HTTP DELETE method) to make a note about how the subscription was canceled.
	ReasonCode                    string    `json:"reason_code,omitEmpty"`                           // (Optional) Can be used when canceling a subscription (via the HTTP DELETE method) to indicate why a subscription was canceled.
	NextBillingAt                 string    `json:"next_billing_at,omitEmpty"`                       // (Optional) Set this attribute to a future date/time to sync imported subscriptions to your existing renewal schedule. See the notes on “Date/Time Format” at https://help.chargify.com/subscriptions/subscriptions-import.html. If you provide a next_billing_at timestamp that is in the future, no trial or initial charges will be applied when you create the subscription. In fact, no payment will be captured at all. The first payment will be captured, according to the prices defined by the product, near the time specified by next_billing_at. If you do not provide a value for next_billing_at, any trial and/or initial charges will be assessed and charged at the time of subscription creation. If the card cannot be successfully charged, the subscription will not be created. See further notes in the section on Importing Subscriptions.
	ExpiresAt                     string    `json:"expires_at,omitEmpty"`                            // Timestamp giving the expiration date of this subscription (if any). You may manually change the expiration date at any point during a subscription period.
	ExpirationTracksChange        bool      `json:"expiration_tracks_next_billing_change,omitEmpty"` // (Optional, default false) When set to true, and when next_billing_at is present, if the subscription expires, the expires_at will be shifted by the same amount of time as the difference between the old and new “next billing” dates.
	VATNumber                     string    `json:"vat_number,omitEmpty"`                            // (Optional) Supplying the VAT number allows EU customer’s to opt-out of the Value Added Tax assuming the merchant address and customer billing address are not within the same EU country. It’s important to omit the country code from the VAT number upon entry. Otherwise, taxes will be assessed upon the purchase.
	CouponCode                    string    `json:"coupon_code,omitEmpty"`                           // (Optional) The coupon code of the coupon to apply ()
	PaymentCollectionMethod       string    `json:"payment_collection_method,omitEmpty"`             // (Optional) The type of payment collection to be used in the subscription. May be automatic, or invoice.
	AgreementTerms                string    `json:"agreement_terms,omitEmpty"`                       // (Optional) The ACH authorization agreement terms. If enabled, an email will be sent to the customer with a copy of the terms.
	ACHFirstName                  string    `json:"authorizer_first_name,omitEmpty"`                 // (Optional) The first name of the person authorizing the ACH agreement.
	ACHLastName                   string    `json:"authorizer_last_name,omitEmpty"`                  // (Optional) The last name of the person authorizing the ACH agreement.
	ChangeDelayed                 bool      `json:"product_change_delayed,omitEmpty"`                // (Optional, used only for https://reference.chargify.com/v1/subscriptions-product-changes-migrations-upgrades-downgrades/update-subscription-product-change When set to true, indicates that a changed value for product_handle should schedule the product change to the next subscription renewal.
	CalendarBilling               string    `json:"calendar_billing,omitEmpty"`                      // (Optional, see https://reference.chargify.com/v1/subscriptions/subscriptions-intro#https://help.chargify.com/subscriptions/billing-dates.html#calendar-billing for more details). Cannot be used when also specifying next_billing_at
	SnapDay                       int       `json:"snap_day,omitEmpty"`                              // A value between 1 and 28, or end
	CalendarBillingFirstDayCharge string    `json:"calendar_billing_first_charge,omitEmpty"`         // (Optional) One of “prorated” (the default – the prorated product price will be charged immediately), “immediate” (the full product price will be charged immediately), or “delayed” (the full product price will be charged with the first scheduled renewal).
	ReceivesInvoiceEmails         bool      `json:"receives_invoice_emails,omitEmpty"`               // (Optional) Default: True - Whether or not this subscription is set to receive emails related to this subscription.
	Customer                      *Customer `json:"customer,omitempty"`
	Product                       *Product  `json:"product,omitempty"`
}

// CreateSubscriptionForCustomer creates a new subscription. When creating a subscription, you must specify a product and a customer.
// The product should be specificed by productHandle and the customer should be specified with customerReference
func CreateSubscriptionForCustomer(customerReference, productHandle string) (*Subscription, error) {
	body := map[string]map[string]string{
		"subscription": map[string]string{
			"customer_reference": customerReference,
			"product_handle":     productHandle,
		},
	}

	ret, err := makeCall(endpoints[endpointSubscriptionCreate], body, nil)
	if err != nil {
		return nil, err
	}
	// if successful, the customer should come back in a map[customer]Customer format
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
		if reasonCode != "" && cancellationMessage != "" {
			reason := map[string]string{
				"cancellation_message": cancellationMessage,
				"reason_code":          reasonCode,
			}
			_, err = makeCall(endpoints[endpointSubscriptionCancelDelayed], reason, &map[string]string{
				"subscriptionID": fmt.Sprintf("%d", subscriptionID),
			})
		}
	}
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
		"migration": map[string]interface{}{
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
