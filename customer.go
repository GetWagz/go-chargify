package chargify

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Customer is a single customer in the chargify account
type Customer struct {
	ID                         int64    `json:"id" mapstructure:"id"`                 //	The customer ID in Chargify
	FirstName                  string   `json:"first_name" mapstructure:"first_name"` //	The first name of the customer
	LastName                   string   `json:"last_name" mapstructure:"last_name"`   //	The last name of the customer
	Email                      string   `json:"email" mapstructure:"email"`           //	The email address of the customer
	CCEmailsRaw                string   `json:"cc_emails" mapstructure:"cc_emails"`   //	(Optional) A comma-separated list of emails that should be cc’d on all customer communications (i.e. “joe@example.com, sue@example.com”)
	CCEmails                   []string // the proccessed CC emails
	Organization               string   `json:"organization" mapstructure:"organization"`                                     //	The organization of the customer
	Reference                  string   `json:"reference" mapstructure:"reference"`                                           //	The unique identifier used within your own application for this customer
	CreatedAt                  string   `json:"created_at" mapstructure:"created_at"`                                         //	The timestamp in which the customer object was created in Chargify
	UpdatedAt                  string   `json:"updated_at" mapstructure:"updated_at"`                                         //	The timestamp in which the customer object was last edited
	Address                    string   `json:"address" mapstructure:"address"`                                               //	The customer’s shipping street address (i.e. “123 Main St.”)
	Address2                   string   `json:"address_2" mapstructure:"address_2"`                                           //	Second line of the customer’s shipping address i.e. “Apt. 100”
	City                       string   `json:"city" mapstructure:"city"`                                                     //	The customer’s shipping address city (i.e. “Boston”)
	State                      string   `json:"state" mapstructure:"state"`                                                   //	The customer’s shipping address state (i.e. “MA”)
	Zip                        string   `json:"zip" mapstructure:"zip"`                                                       //	The customer’s shipping address zip code (i.e. “12345”)
	Country                    string   `json:"country" mapstructure:"country"`                                               //	The customer shipping address country, perferably in  format (i.e. “US”)
	Phone                      string   `json:"phone" mapstructure:"phone"`                                                   //	The phone number of the customer
	Verified                   bool     `json:"verified" mapstructure:"verified"`                                             //	Is the customer verified to use ACH as a payment method. Available only on Authorize.Net gateway
	PortalCustomerCreatedAt    string   `json:"portal_customer_created_at" mapstructure:"portal_customer_created_at"`         //	The timestamp of when the Billing Portal entry was created at for the customer
	PortalInviteLastSentAt     string   `json:"portal_invite_last_sent_at" mapstructure:"portal_invite_last_sent_at"`         //	The timestamp of when the Billing Portal invite was last sent at
	PortalInviteLastAcceptedAt string   `json:"portal_invite_last_accepted_at" mapstructure:"portal_invite_last_accepted_at"` //	The timestamp of when the Billing Portal invite was last accepted
	TaxExempt                  bool     `json:"tax_exempt" mapstructure:"tax_exempt"`                                         //	(Optional) The tax exempt status for the customer. Acceptable values are true or 1 for true and false or 0 for false.
	VatNumber                  string   `json:"vat_number" mapstructure:"vat_number"`                                         //	(Optional) The VAT number, if applicable
}

// CreateCustomer creates a new customer on chargify
func CreateCustomer(input *Customer) (*Customer, error) {
	if input.FirstName == "" || input.LastName == "" || input.Email == "" {
		return nil, errors.New("first name, last name, and email are all required")
	}
	body := map[string]Customer{
		"customer": *input,
	}

	ret, err := makeCall(endpoints[endpointCustomerCreate], body, nil)
	if err != nil {
		return nil, err
	}
	// if successful, the customer should come back in a map[customer]Customer format
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	customer := &Customer{}
	err = mapstructure.Decode(apiBody["customer"], customer)
	return customer, err
}

// UpdateCustomer updates a customer in chargify
func UpdateCustomer(input *Customer) error {
	body := map[string]Customer{
		"customer": *input,
	}
	ret, err := makeCall(endpoints[endpointCustomerUpdate], body, &map[string]string{
		"id": fmt.Sprintf("%d", input.ID),
	})
	if err != nil {
		return err
	}
	if ret.HTTPCode != http.StatusOK {
		return errors.New("could not update that customer")
	}

	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["customer"], input)
	return err
}

// GetCustomerByID gets a customer by chargify id
func GetCustomerByID(id int) (*Customer, error) {
	ret, err := makeCall(endpoints[endpointCustomerGet], nil, &map[string]string{
		"id": fmt.Sprintf("%d", id),
	})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return nil, err
	}

	entry := ret.Body.(map[string]interface{})
	custRaw := entry["customer"]
	customer := Customer{}
	err = mapstructure.Decode(custRaw, &customer)
	if err == nil {
		return &customer, nil
	}
	return nil, err
}

// GetCustomerByReference gets a customer by reference
func GetCustomerByReference(reference string) (*Customer, error) {
	ret, err := makeCall(endpoints[endpointCustomerByReferenceGet], nil, &map[string]string{
		"reference": reference,
	})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return nil, err
	}

	entry := ret.Body.(map[string]interface{})
	custRaw := entry["customer"]
	customer := Customer{}
	err = mapstructure.Decode(custRaw, &customer)
	if err == nil {
		return &customer, nil
	}
	return nil, err
}

// DeleteCustomerByID deletes a customer from chargify permanently
func DeleteCustomerByID(id int64) error {
	_, err := makeCall(endpoints[endpointCustomerDelete], nil, &map[string]string{
		"id": fmt.Sprintf("%d", id),
	})
	return err
}

// GetCustomers gets the customers for the site
func GetCustomers(page int, sortDir string) (found []Customer, err error) {
	sortDir = strings.ToLower(sortDir)
	if sortDir != "asc" && sortDir != "desc" {
		return found, errors.New("sortDir must be asc or desc")
	}
	if page < 1 {
		return found, errors.New("page must be 1 or higher, not 0 indexed")
	}

	ret, err := makeCall(endpoints[endpointCustomersGet], &map[string]string{
		"direction": sortDir,
		"page":      fmt.Sprintf("%d", page),
	}, nil)
	if err != nil || ret.HTTPCode != http.StatusOK {
		return
	}

	// so, Chargify violates OWASP best practices by returning these in an array
	temp := ret.Body.([]interface{})
	for i := range temp {
		entry := temp[i].(map[string]interface{})
		custRaw := entry["customer"]
		customer := Customer{}
		err = mapstructure.Decode(custRaw, &customer)
		if err == nil {
			found = append(found, customer)
		}
	}
	return found, nil

}

// GetCustomerSubscriptions
func GetCustomerSubscriptions(customerID int) (found []Subscription, err error) {
	ret, err := makeCall(endpoints[endpointCustomerSubscriptionsList], nil, &map[string]string{
		"customer_id": fmt.Sprintf("%d", customerID),
	})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return
	}
	temp := ret.Body.([]interface{})
	for i := range temp {
		entry := temp[i].(map[string]interface{})
		raw := entry["subscription"]
		entity := Subscription{}
		err = mapstructure.Decode(raw, &entity)
		if err == nil {
			found = append(found, entity)
		} else {
			return nil, err
		}
	}
	return
}

// SearchForCustomerByReference searches for a customer by it's reference value. It first performs the large search then
// looks for the substring in the returned values
func SearchForCustomerByReference(reference string) (Customer, error) {
	found := Customer{}

	customers, err := SearchForCustomersByReference(reference)
	if err != nil {
		return found, err
	}
	for i := range customers {
		if customers[i].Reference == reference {
			found = customers[i]
			break
		}
	}
	if found.ID == 0 {
		return found, errors.New("customer not found")
	}
	return found, nil
}

// SearchForCustomersByReference searches all of the customers for a specific reference
func SearchForCustomersByReference(reference string) ([]Customer, error) {
	found := []Customer{}
	var err error
	ret, err := makeCall(endpoints[endpointCustomersGet], map[string]string{
		"q": reference,
	}, nil)
	if err != nil || ret.HTTPCode != http.StatusOK {
		return found, err
	}

	// so, Chargify violates OWASP best practices by returning these in an array
	temp := ret.Body.([]interface{})
	for i := range temp {
		cust := Customer{}
		entry := temp[i].(map[string]interface{})
		custRaw := entry["customer"]
		err = mapstructure.Decode(custRaw, &cust)
		if err != nil {
			return []Customer{}, err
		}
		found = append(found, cust)
	}
	return found, err
}

// SearchForCustomersByEmail searches for customers with a specific email address; multiple can exist
func SearchForCustomersByEmail(email string) ([]Customer, error) {
	found := []Customer{}
	var err error
	ret, err := makeCall(endpoints[endpointCustomersGet], map[string]string{
		"q": email,
	}, nil)
	if err != nil || ret.HTTPCode != http.StatusOK {
		return found, err
	}

	// so, Chargify violates OWASP best practices by returning these in an array
	temp := ret.Body.([]interface{})
	for i := range temp {
		entry := temp[i].(map[string]interface{})
		custRaw := entry["customer"]
		customer := Customer{}
		err = mapstructure.Decode(custRaw, &customer)
		if err == nil {
			found = append(found, customer)
		}
	}
	return found, err
}

func createTestCustomer() (*Customer, *PaymentProfile, error) {
	customID := rand.Int63n(999999999)
	input := Customer{
		FirstName: fmt.Sprintf("First-%d", customID),
		LastName:  fmt.Sprintf("Last-%d", customID),
		Email:     fmt.Sprintf("test+%d@example.com", customID),
		Reference: fmt.Sprintf("test-lib-%d", customID),
	}

	customer, err := CreateCustomer(&input)
	if err != nil {
		return nil, nil, err
	}
	profile, err := SavePaymentProfileVault(customer.ID, VaultBogus, "12345")
	if err != nil {
		return nil, nil, err
	}
	return customer, profile, nil
}
