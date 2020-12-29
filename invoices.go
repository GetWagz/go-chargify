package chargify

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Invoice is a relationship invoice on Chargify. Note that not all fields are currently implemented, as there
// are over 145 properties on an invoice
type Invoice struct {
	UID               string
	SiteID            int64
	CustomerID        int64
	SubscriptionID    int64
	Number            string // oof
	SequenceNumber    int64
	IssueDate         string
	DueDate           string
	PaidDate          string
	Status            string
	Currency          string
	ProductName       string
	ProductFamilyName string
	Customer          Customer
	Payments          []Payment
	Refunds           []Refund
}

// Refund is a single refund issued against an invoice
type Refund struct {
	TransactionID  int64
	PaymentID      int64
	Memo           string
	OriginalAmount string
	AppliedAmount  string
}

// InvoiceQueryParams are a collection of implemented query params to pass in to the invoice
// get call
type InvoiceQueryParams struct {
	StartDate      *string `json:"start_date,omit_empty" mapstructure:"start_date"`
	EndDate        *string `json:"end_date,omit_empty" mapstructure:"end_date"`
	Status         *string `json:"status,omit_empty" mapstructure:"status"`
	SubscriptionID *int64  `json:"subscription_id,omit_empty" mapstructure:"subscription_id"`
	Page           *int64  `json:"page,omit_empty" mapstructure:"page"`
	PerPage        *int64  `json:"per_page,omit_empty" mapstructure:"per_page"`
	Direction      *string `json:"direction,omit_empty" mapstructure:"direction"`
}

// GetInvoices searched for invoices based upon passed-in params
func GetInvoices(queryParams *InvoiceQueryParams) ([]Invoice, error) {
	invoices := []Invoice{}

	ret, err := makeCall(endpoints[endpointGetInvoices], queryParams, &map[string]string{})
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
