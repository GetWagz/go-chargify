package chargify

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// PaymentProfile represents a payment profile. Note that many of the fields that are "numbers" are actually strings due to leading 0s.
// A lot of these are omitempty since certain fields are only used in certain calls and the API can get confused easily
type PaymentProfile struct {
	ID                    int64       `json:"id" mapstructure:"id"`                                             // the id after the profile is created
	SubscriptionID        int64       `json:"subscription" mapstructure:"subscription"`                         // the subscription id after the profile is created
	PaymentType           string      `json:"payment_type" mapstructure:"payment_type"`                         // (Optional) Default is credit_card. May be bank_account or credit_card or paypal_account.
	CustomerID            int64       `json:"customer_id" mapstructure:"customer_id"`                           // 	(Required when creating a new payment profile) The Chargify customer id.
	FirstName             string      `json:"first_name" mapstructure:"first_name"`                             // 	First name on card or bank account
	LastName              string      `json:"last_name" mapstructure:"last_name"`                               // 	Last name on card or bank account
	FullNumber            string      `json:"full_number" mapstructure:"full_number"`                           // 	(Required when payment_type is credit_card unless you provide the vault_token) The full credit card number (string representation, i.e. 5424000000000015)
	ExpirationMonth       string      `json:"expiration_month" mapstructure:"expiration_month"`                 // 	(Required when payment_type is credit_card unless you provide the vault_token) The 1- or 2-digit credit card expiration month, as an integer or string, i.e. 5
	ExpirationYear        string      `json:"expiration_year" mapstructure:"expiration_year"`                   // 	(Required when payment_type is credit_card unless you provide the vault_token) The 4-digit credit card expiration year, as an integer or string, i.e. 2012
	CVV                   string      `json:"cvv" mapstructure:"cvv"`                                           // 	(Optional, may be required by your gateway settings) The 3- or 4-digit Card Verification Value. This value is merely passed through to the payment gateway.
	BillingAddress        string      `json:"billing_address" mapstructure:"billing_address"`                   // 	(Optional, may be required by your product configuration or gateway settings) The credit card or bank account billing street address (i.e. 123 Main St.). This value is merely passed through to the payment gateway.
	BillingAddress2       string      `json:"billing_address_2" mapstructure:"billing_address_2"`               // 	(Optional) Second line of the customer’s billing address i.e. Apt. 100
	BillingCity           string      `json:"billing_city" mapstructure:"billing_city"`                         // 	(Optional, may be required by your product configuration or gateway settings) The credit card or bank account billing address city (i.e. Boston). This value is merely passed through to the payment gateway.
	BillingState          string      `json:"billing_state" mapstructure:"billing_state"`                       // 	(Optional, may be required by your product configuration or gateway settings) The credit card or bank account billing address state (i.e. MA). This value is merely passed through to the payment gateway.
	BillingZip            string      `json:"billing_zip" mapstructure:"billing_zip"`                           // 	(Optional, may be required by your product configuration or gateway settings) The credit card or bank account billing address zip code (i.e. “12345”). This value is merely passed through to the payment gateway.
	BillingCountry        string      `json:"billing_country" mapstructure:"billing_country"`                   // 	(Optional, may be required by your product configuration or gateway settings) The credit card or bank account billing address country, preferably in  format (i.e. “US”). This value is merely passed through to the payment gateway. Some gateways require country codes in a specific format. Please check your gateway’s documentation. If creating an ACH subscription, only US is supported at this time.
	BankName              string      `json:"bank_name" mapstructure:"bank_name"`                               // 	(Required when creating a subscription with ACH) The name of the bank where the customer’s account resides
	BankRouting           string      `json:"bank_routing_number" mapstructure:"bank_routing_number"`           // 	(Required when creating a subscription with ACH) The routing number of the bank
	BankAccount           string      `json:"bank_account_number" mapstructure:"bank_account_number"`           // 	Required when creating a subscription with ACH) The customer’s bank account number
	BankAccountType       string      `json:"bank_account_type" mapstructure:"bank_account_type"`               // 	When payment_type is bank_account, this defaults to checking and cannot be changed
	BankAccountHolderType string      `json:"bank_account_holder_type" mapstructure:"bank_account_holder_type"` // 	When payment_type is bank_account, may be personal (default) or business
	Verified              bool        `json:"verified,omitempty" mapstructure:"verified"`                       // 	When payment type is bank_account and current_vault is stripe_connect, may be set to true to indicate that the bank account has already been verified.
	PaypalEmail           string      `json:"paypal_email" mapstructure:"paypal_email"`                         //
	PaymentMethodNonce    string      `json:"payment_method_nonce" mapstructure:"payment_method_nonce"`         //
	VaultToken            string      `json:"vault_token" mapstructure:"vault_token"`                           // 	(Only allowed during the creation of a new payment profile.) If you have an existing vault_token from your gateway, you may associate it with this new payment profile.
	ChargifyToken         string      `json:"chargify_token" mapstructure:"chargify_token"`                     // 	(Optional) Token received after sending billing informations using . This token must be passed along with customer_id attribute (i.e. tok_9g6hw85pnpt6knmskpwp4ttt)
	CurrentVault          VaultMethod `json:"current_vault" mapstructure:"current_vault"`                       // 	(Required when you pass in a vault_token.) Will be one of the following: bogus (for testing), authorizenet, authorizenet_cim, beanstream, bpoint, braintree_blue, chargify, cybersource, elavon, eway, eway_rapid_std , firstdata, fusebox, litle, moneris, moneris_us, orbital, payment_express, paymill, pin, quickpay, square, stripe_connect, trust_commerce, wirecard. Provides a hint about where the credit card details represented by vault_token are stored, however transactions will always be sent to the gateway configured in the Site's settings.
	CardType              string      `json:"card_type" mapstructure:"card_type"`                               // 	Can be any of the following visa, master, discover, american_express, diners_club, jcb, switch, solo, dankort, maestro, forbrugsforeningen, laser
}

// VaultMethod represents one of the payment vaults for use with tokenization. This is generally the recommended way to handle payment methods.
type VaultMethod string

var (
	// VaultBogus represents a bogus vault used for testing
	VaultBogus VaultMethod = "bogus"
	// VaultAuthorize represents the AuthorizeNet vault
	VaultAuthorize VaultMethod = "authorizenet"
	// VaultAuthorizeCIM represents the AuthorizeCIM vault
	VaultAuthorizeCIM VaultMethod = "authorizenet_cim"
	// VaultBeanStream represents the BeanStream vault
	VaultBeanStream VaultMethod = "beanstream"
	// VaultBPoint represents the BPoint vault
	VaultBPoint VaultMethod = "bpoint"
	// VaultBraintree represents the BrainTree vault
	VaultBraintree VaultMethod = "braintree_blue"
	// VaultChargify represents the Chargify vault
	VaultChargify VaultMethod = "chargify"
	// VaultCyberSource represents the CyberSource vault
	VaultCyberSource VaultMethod = "cybersource"
	// VaultElavon represents the Elavon vault
	VaultElavon VaultMethod = "elavon"
	// VaultEWay represents the EWay vault
	VaultEWay VaultMethod = "eway"
	// VaultEWayRapid represents the EWayRapid vault
	VaultEWayRapid VaultMethod = "eway_rapid_std"
	// VaultFirstData represents the FirstData vault
	VaultFirstData VaultMethod = "firstdata"
	// VaultFuseBox represents the FuseBox vault
	VaultFuseBox VaultMethod = "fusebox"
	// VaultLittle represents the Little vault
	VaultLittle VaultMethod = "litle"
	// VaultMoneris represents the Moneris vault
	VaultMoneris VaultMethod = "moneris"
	// VaultMonerisUS represents the Moneris US vault
	VaultMonerisUS VaultMethod = "moneris_us"
	// VaultOrbital represents the Orbital vault
	VaultOrbital VaultMethod = "orbital"
	// VaultPaymentExpress represents the PaymentExpress vault
	VaultPaymentExpress VaultMethod = "payment_express"
	// VaultPaymill represents the PayMill vault
	VaultPaymill VaultMethod = "paymill"
	// VaultPIN represents the PIN vault
	VaultPIN VaultMethod = "pin"
	// VaultQuickPay represents the QuickPay vault
	VaultQuickPay VaultMethod = "quickpay"
	// VaultSquare represents the Square vault
	VaultSquare VaultMethod = "square"
	// VaultStripe represents the Stripe vault
	VaultStripe VaultMethod = "stripe_connect"
	// VaultTrustCommerce represents the TrustCommerce vault
	VaultTrustCommerce VaultMethod = "trust_commerce"
	// VaultWireCard represents the WireCard vault
	VaultWireCard VaultMethod = "wirecard"
)

// Payment represents a single payment on an invoice, for example
type Payment struct {
	TransactionTime string        `json:"transaction_time" mapstructure:"transaction_time"`
	Memo            string        `json:"memo" mapstructure:"memo"`
	OriginalAmount  string        `json:"original_amount" mapstructure:"original_amount"`
	AppliedAmount   string        `json:"applied_amount" mapstructure:"applied_amount"`
	TransactionID   int64         `json:"transaction_id" mapstructure:"transaction_id"`
	Prepayment      bool          `json:"prepayment" mapstructure:"prepayment"`
	PaymentMethod   PaymentMethod `json:"payment_method" mapstructure:"payment_method"`
}

// PaymentMethod represents a payment method, found on a payment struct
type PaymentMethod struct {
	Details          string `json:"details" mapstructure:"details"`
	Kind             string `json:"kind" mapstructure:"kind"`
	Memo             string `json:"memo" mapstructure:"memo"`
	PaymentType      string `json:"payment_type" mapstructure:"payment_type"`
	CardBrand        string `json:"card_brand" mapstructure:"card_brand"`
	CardExpiration   string `json:"card_expiration" mapstructure:"card_expiration"`
	LastFour         string `json:"last_four" mapstructure:"last_four"`
	MaskedCardNumber string `json:"masked_card_number" mapstructure:"masked_card_number"`
}

// SavePaymentProfileForCustomer saves a new payment profile. Note that this is a raw save; for ease of use it may be better to use one of the other SavePaymentProfile* methods
func SavePaymentProfileForCustomer(customerID int64, input *PaymentProfile) error {
	body := map[string]PaymentProfile{
		"payment_profile": *input,
	}

	ret, err := makeCall(endpoints[endpointPaymentProfileCreate], body, nil)
	if err != nil {
		return err
	}
	if ret.HTTPCode != http.StatusCreated {
		return errors.New("could not create that profile")
	}
	// if successful, the profile should come back in a map[payment_profile]PaymentProfile format
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["payment_profile"], input)
	return err
}

// SavePaymentProfileVault saves a payment profile using a vault
func SavePaymentProfileVault(customerID int64, vault VaultMethod, vaultToken string) (*PaymentProfile, error) {
	// TODO: make sure everything is valid
	profile := &PaymentProfile{
		CustomerID:   customerID,
		VaultToken:   vaultToken,
		CurrentVault: vault,
	}
	return profile, SavePaymentProfileForCustomer(customerID, profile)
}

// SavePaymentProfileACH saves a payment profile using ACH
func SavePaymentProfileACH(customerID int64, bankName, bankRoutingNumber, bankAccountNumber, bankAccountType, bankAccountHolderType string) (*PaymentProfile, error) {
	// TODO: make sure everything is valid
	profile := &PaymentProfile{
		CustomerID:            customerID,
		BankName:              bankName,
		BankRouting:           bankRoutingNumber,
		BankAccount:           bankAccountNumber,
		BankAccountType:       bankAccountType,
		BankAccountHolderType: bankAccountHolderType,
	}
	return profile, SavePaymentProfileForCustomer(customerID, profile)
}

// DeletePaymentProfile deletes a payment profile
func DeletePaymentProfile(subscriptionID int64, profileID int64) error {

	ret, err := makeCall(endpoints[endpointPaymentProfileDelete], nil, &map[string]string{
		"subscriptionID": fmt.Sprintf("%v", subscriptionID),
		"profileID":      fmt.Sprintf("%v", profileID),
	})
	if err != nil {
		return err
	}
	if ret.HTTPCode != http.StatusNoContent {
		return errors.New("could not delete that profile")
	}
	return nil
}

// UpdatePaymentProfile updates a payment profile
func UpdatePaymentProfile(input *PaymentProfile) error {
	body := map[string]PaymentProfile{
		"payment_profile": *input,
	}

	ret, err := makeCall(endpoints[endpointPaymentProfileUpdate], body, &map[string]string{
		"paymentProfileID": fmt.Sprintf("%d", input.ID),
	})
	if err != nil {
		return err
	}
	if ret.HTTPCode != http.StatusOK {
		return errors.New("could not update that profile")
	}
	// if successful, the profile should come back in a map[payment_profile]PaymentProfile format
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["payment_profile"], input)
	return err
}
