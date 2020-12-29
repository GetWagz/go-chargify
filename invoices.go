package chargify

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Invoice is a relationship invoice on Chargify. Note that not all fields are currently implemented, as there
// are over 145 properties on an invoice
type Invoice struct {
	UID               string    `json:"uid,omitempty" mapstructure:"uid"`
	SiteID            int64     `json:"site_id,omitempty" mapstructure:"site_id"`
	CustomerID        int64     `json:"customer_id,omitempty" mapstructure:"customer_id"`
	SubscriptionID    int64     `json:"subscription_id,omitempty" mapstructure:"subscription_id"`
	Number            string    `json:"number,omitempty" mapstructure:"number"`
	SequenceNumber    int64     `json:"sequence_number,omitempty" mapstructure:"sequence_number"`
	IssueDate         string    `json:"issue_date,omitempty" mapstructure:"issue_date"`
	DueDate           string    `json:"due_date,omitempty" mapstructure:"due_date"`
	PaidDate          string    `json:"paid_date,omitempty" mapstructure:"paid_date"`
	Status            string    `json:"status,omitempty" mapstructure:"status"`
	Currency          string    `json:"currency,omitempty" mapstructure:"currency"`
	ProductName       string    `json:"product_name,omitempty" mapstructure:"product_name"`
	ProductFamilyName string    `json:"product_family_name,omitempty" mapstructure:"product_family_name"`
	Customer          Customer  `json:"customer,omitempty" mapstructure:"customer"`
	Payments          []Payment `json:"payments,omitempty" mapstructure:"payments"`
	Refunds           []Refund  `json:"refunds,omitempty" mapstructure:"refunds"`
}

// Refund is a single refund issued against an invoice
type Refund struct {
	TransactionID  int64  `json:"transaction_id,omitempty" mapstructure:"transaction_id"`
	PaymentID      int64  `json:"payment_id,omitempty" mapstructure:"payment_id"`
	Memo           string `json:"memo,omitempty" mapstructure:"memo"`
	OriginalAmount string `json:"original_amount,omitempty" mapstructure:"original_amount"`
	AppliedAmount  string `json:"applied_amount,omitempty" mapstructure:"applied_amount"`
}

// InvoiceQueryParams are a collection of implemented query params to pass in to the invoice
// get call
type InvoiceQueryParams struct {
	StartDate      string `json:"start_date,omitempty" mapstructure:"start_date"`
	EndDate        string `json:"end_date,omitempty" mapstructure:"end_date"`
	Status         string `json:"status,omitempty" mapstructure:"status"`
	SubscriptionID int64  `json:"subscription_id,omitempty" mapstructure:"subscription_id"`
	Page           int64  `json:"page,omitempty" mapstructure:"page"`
	PerPage        int64  `json:"per_page,omitempty" mapstructure:"per_page"`
	Direction      string `json:"direction,omitempty" mapstructure:"direction"`
}

// GetInvoices searched for invoices based upon passed-in params
func GetInvoices(queryParams *InvoiceQueryParams) ([]Invoice, error) {
	invoices := []Invoice{}

	// massage the params into map[string]string
	params := map[string]string{}
	if queryParams != nil && queryParams.StartDate != "" {
		params["start_date"] = queryParams.StartDate
	}
	if queryParams != nil && queryParams.EndDate != "" {
		params["end_date"] = queryParams.EndDate
	}
	if queryParams != nil && queryParams.Status != "" {
		params["status"] = queryParams.Status
	}
	if queryParams != nil && queryParams.SubscriptionID != 0 {
		params["subscription_id"] = fmt.Sprintf("%d", queryParams.SubscriptionID)
	}
	if queryParams != nil && queryParams.Page != -1 {
		params["page"] = fmt.Sprintf("%d", queryParams.Page)
	}
	if queryParams != nil && queryParams.PerPage != 0 {
		params["per_page"] = fmt.Sprintf("%d", queryParams.PerPage)
	}
	if queryParams != nil && queryParams.Direction != "" {
		params["direction"] = queryParams.Direction
	}

	ret, err := makeCall(endpoints[endpointGetInvoices], params, &map[string]string{})
	if err != nil {
		return invoices, err
	}

	// chargify violates OWASP best practices and returns an array, so
	apiBody, bodyOK := ret.Body.([]map[string]interface{})
	if !bodyOK {
		return invoices, errors.New("could not understand server response")
	}

	for i := range apiBody {
		invoice := Invoice{}
		err = mapstructure.Decode(apiBody[i], &invoice)
		if err != nil {
			return invoices, fmt.Errorf("could not decode element %d: %+v", i, err)
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

// GetInvoiceByID gets a single relationship invoice
func GetInvoiceByID(invoiceID int64) (*Invoice, error) {
	invoice := &Invoice{}

	ret, err := makeCall(endpoints[endpointGetInvoice], nil, &map[string]string{
		"invoiceID": fmt.Sprintf("%d", invoiceID),
	})
	if err != nil {
		return invoice, err
	}

	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return invoice, errors.New("could not understand server response")
	}

	err = mapstructure.Decode(apiBody, &invoice)
	if err != nil {
		return invoice, fmt.Errorf("could not decode return value: %+v", err)
	}

	return invoice, nil
}

// RefundInvoice refunds a single invoice. Note that the amount is a string, which expects a decimal. This is unusual and will catch you
// off guard if you are not carefule. So, for example, pass in "10.50" for ten dollars and fifty cents. Also note that the required fields are
// amount, memo, and paymentID
func RefundInvoice(invoiceID, amount, memo string, paymentID int64, external, applyCredit, voidInvoice bool) (*Invoice, error) {
	invoice := &Invoice{}

	params := map[string]map[string]string{
		"refund": {
			"amount":       amount,
			"memo":         memo,
			"payment_id":   fmt.Sprintf("%d", paymentID),
			"external":     fmt.Sprintf("%v", external),
			"apply_credit": fmt.Sprintf("%v", applyCredit),
			"void_invoice": fmt.Sprintf("%v", voidInvoice),
		},
	}

	ret, err := makeCall(endpoints[endpointRefundInvoice], params, &map[string]string{
		"invoiceID": invoiceID,
	})
	if err != nil {
		return invoice, err
	}
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody, invoice)
	return invoice, err

}
