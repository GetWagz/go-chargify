package chargify

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// BillingPortal represents a self-service management portal on the Chargify web site.
// If your customer has been invited to the Billing Portal, then they will receive a link to manage their subscription (the “Management URL”) automatically
// at the bottom of their statements, invoices, and receipts. This link changes periodically for security and is only valid for 65 days.
type BillingPortal struct {
	URL                string `json:"url" mapstructure:"url"`
	FetchCount         int64  `json:"fetch_count" mapstructure:"fetch_count"`
	CreatedAt          string `json:"created_at" mapstructure:"created_at"`
	NewLinkAvailableAt string `json:"new_link_available_at" mapstructure:"new_link_available_at"`
	ExpiresAt          string `json:"expires_at" mapstructure:"expires_at"`
}

// EnableBillingPortal enables billing portal management for the customer. Note that it will return an error
// if the portal is already enabled. Confusingly, the decision to send an invite is a query string parameter here
// rather than a HTTP body data object: https://reference.chargify.com/v1/billing-portal/enabling-billing-portal-for-customer
func EnableBillingPortal(customerID int64, sendInvitation bool) error {
	var err error
	if sendInvitation {
		_, err = makeCall(endpoints[endpointBillingPortalEnableAndInvite], nil, &map[string]string{
			"id": fmt.Sprintf("%d", customerID),
		}, nil)
	} else {
		_, err = makeCall(endpoints[endpointBillingPortalEnable], nil, &map[string]string{
			"id": fmt.Sprintf("%d", customerID),
		}, nil)
	}
	return err

}

// GetBillingPortal gets the billing portal information for the customer
func GetBillingPortal(customerID int64) (*BillingPortal, error) {
	ret, err := makeCall(endpoints[endpointBillingPortalGet], nil, &map[string]string{
		"id": fmt.Sprintf("%d", customerID),
	}, nil)
	if err != nil {
		return nil, err
	}
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	portal := &BillingPortal{}
	err = mapstructure.Decode(apiBody, portal)
	return portal, err
}
