package chargify

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
)

// ProductInterval represents an interval used for various calculations in a product
type ProductInterval string

var (
	// ProductIntervalMonth represents an interval of month
	ProductIntervalMonth ProductInterval = "month"
	// ProductIntervalDay represents an interval of day
	ProductIntervalDay ProductInterval = "day"
)

// Component represents a single component
type Component struct {
	ID                        int64   `json:"id"`
	Name                      string  `json:"name" mapstructure:"name"`
	Handle                    string  `json:"handle" mapstructure:"handle"`
	Description               string  `json:"description" mapstructure:"description"`
	PricingScheme             string  `json:"pricing_scheme" mapstructure:"pricing_scheme"`
	UnitName                  string  `json:"unit_name" mapstructure:"unit_name"`
	UnitPrice                 *string `json:"unit_price" mapstructure:"unit_price"`
	ProductFamilyID           int64   `json:"product_family_id" mapstructure:"product_family_id"`
	ProductFamilyName         string  `json:"product_family_name" mapstructure:"product_family_name"`
	Kind                      string  `json:"kind" mapstructure:"kind"`
	Archived                  bool    `json:"archived" mapstructure:"archived"`
	Taxable                   bool    `json:"taxable" mapstructure:"taxable"`
	DefaultPricePointID       int64   `json:"default_price_point_id" mapstructure:"default_price_point_id"`
	PricePointCount           int64   `json:"price_point_count" mapstructure:"price_point_count"`
	PricePointsUrl            string  `json:"price_points_url" mapstructure:"price_points_url"`
	TaxCode                   int64   `json:"tax_code" mapstructure:"tax_code"`
	Recurring                 bool    `json:"recurring" mapstructure:"recurring"`
	UpgradeCharge             *string `json:"upgrade_charge" mapstructure:"upgrade_charge"`
	DowngradeCredit           *string `json:"downgrade_credit" mapstructure:"downgrade_credit"`
	DefaultPricePointName     string  `json:"default_price_point_name" mapstructure:"default_price_point_name"`
	HideDateRangeOnInvoice    bool    `json:"hide_date_range_on_invoice" mapstructure:"hide_date_range_on_invoice"`
	Prices                    []Price `json:"prices" mapstructure:"prices"`
	OveragePrices             []Price `json:"overage_prices" mapstructure:"overage_prices"`
	CreatedAt                 string  `json:"created_at" mapstructure:"created_at"`
	UpdatedAt                 string  `json:"updated_at" mapstructure:"updated_at"`
	AllowFractionalQuantities bool    `json:"allow_fractional_quantities" mapstructure:"allow_fractional_quantities"`
}

type Price struct {
	ID                 int64  `json:"id"`
	ComponentID        int64  `json:"component_id" mapstructure:"component_id"`
	StartingQuantity   int64  `json:"starting_quantity" mapstructure:"starting_quantity"`
	EndingQuantity     int64  `json:"ending_quantity" mapstructure:"ending_quantity"`
	UnitPrice          string `json:"unit_price" mapstructure:"unit_price"`
	PricePointID       int64  `json:"price_point_id" mapstructure:"price_point_id"`
	FormattedUnitPrice string `json:"formatted_unit_price" mapstructure:"formatted_unit_price"`
}

// Product represents a single product
type Product struct {
	ID                      int64           `json:"id"`
	PriceInCents            int             `json:"price_in_cents" mapstructure:"price_in_cents"`                       //	The product price, in integer cents
	Name                    string          `json:"name" mapstructure:"name"`                                           //	The product name
	Handle                  string          `json:"handle" mapstructure:"handle"`                                       //	The product API handle
	Description             string          `json:"description" mapstructure:"description"`                             //	The product description
	ProductFamily           *ProductFamily  `json:"product_family" mapstructure:"product_family"`                       //	Nested attributes pertaining to the product family to which this product belongs
	IntervalUnit            ProductInterval `json:"interval_unit" mapstructure:"interval_unit"`                         // A string representing the interval unit for this product, either month or day
	IntervalValue           int             `json:"interval,omitempty" mapstructure:"interval"`                         // The numerical interval. i.e. an interval of ‘30’ coupled with an interval_unit of day would mean this product would renew every 30 days
	InitialChargeInCents    int             `json:"initial_charge_in_cents" mapstructure:"initial_charge_in_cents"`     // The up front charge you have specified.
	TrialPriceInCents       *int            `json:"trial_price_in_cents,omitempty" mapstructure:"trial_price_in_cents"` // The price of the trial period for a subscription to this product, in integer cents.
	TrialIntervalValue      *int            `json:"trial_interval,omitempty" mapstructure:"trial_interval"`             // A numerical interval for the length of the trial period of a subscription to this product. See the description of interval for a description of how this value is coupled with an interval unit to calculate the full interval
	TrialIntervalUnit       ProductInterval `json:"trial_interval_unit" mapstructure:"trial_interval_unit"`             // A string representing the trial interval unit for this product, either month or day
	ExpirationIntervalValue *int            `json:"expiration_interval,omitempty" mapstructure:"expiration_interval"`   // A numerical interval for the length a subscription to this product will run before it expires. See the description of interval for a description of how this value is coupled with an interval unit to calculate the full interval
	ExpirationIntervalUnit  ProductInterval `json:"expiration_interval_unit" mapstructure:"expiration_interval_unit"`   // A string representing the trial interval unit for this product, either month or day
	VersionNumber           float64         `json:"version_number" mapstructure:"version_number"`                       // The version of the product
	UpdateReturnURL         string          `json:"update_return_url" mapstructure:"update_return_url"`                 // The url to which a customer will be returned after a successful account update
	UpdateReturnParams      string          `json:"update_return_params" mapstructure:"update_return_params"`           // The parameters will append to the url after a successful account update
	RequireCreditCard       bool            `json:"require_credit_card" mapstructure:"require_credit_card"`             // Boolean
	RequestCreditCard       bool            `json:"request_credit_card" mapstructure:"request_credit_card"`             // Boolean
	CreatedAt               string          `json:"created_at" mapstructure:"created_at"`                               // Timestamp indicating when this product was created
	UpdatedAt               string          `json:"updated_at" mapstructure:"updated_at"`                               // Timestamp indicating when this product was last updated
	Archived                string          `json:"archived_at" mapstructure:"archived_at"`                             // Timestamp indicating when this product was archived
	SignupPages             *[]SignupPage   `json:"public_signup_pages" mapstructure:"public_signup_pages"`             // An array of signup pages
	AutoCreateSignupPage    bool            `json:"auto_create_signup_page" mapstructure:"auto_create_signup_page"`     // Whether or not to create a signup page
	TaxCode                 string          `json:"tax_code" mapstructure:"tax_code"`                                   // A string representing the tax code related to the product type. This is especially important when using the Avalara service to tax based on locale. This attribute has a max length of 10 characters.

}

// SignupPage represents a product's signup page, if needed
type SignupPage struct {
	ID           int64  `json:"id"`                                         // The id of the signup page (public_signup_pages only)
	URL          string `json:"url" mapstructure:"url"`                     // The url where the signup page can be viewed (public_signup_pages only)
	ReturnURL    string `json:"return_url" mapstructure:"return_url"`       // The url to which a customer will be returned after a successful signup (public_signup_pages only)
	ReturnParams string `json:"return_params" mapstructure:"return_params"` // The params to be appended to the return_url (public_signup_pages only)
}

// ProductFamily represents a product family
type ProductFamily struct {
	ID             int64  `json:"id"`
	Name           string `json:"name" mapstructure:"name"`                       //	The product family name
	Handle         string `json:"handle" mapstructure:"handle"`                   //	The product family API handle
	AccountingCode string `json:"accounting_code" mapstructure:"accounting_code"` // The product family accounting code (has no bearing in Chargify, may be used within your app)
	Description    string `json:"description" mapstructure:"description"`         // The product family description
	CreatedAt      string `json:"created_at" mapstructure:"created_at"`
	UpdatedAt      string `json:"updated_at" mapstructure:"updated_at"`
}

// CreateProductFamily creates a new product family
func CreateProductFamily(name, description, handle string, accountingCode string) (*ProductFamily, error) {
	family := &ProductFamily{
		Name:           name,
		Description:    description,
		Handle:         handle,
		AccountingCode: accountingCode,
	}
	if family.Name == "" || family.Description == "" || family.Handle == "" {
		return nil, errors.New("name, handle, and description are all required")
	}
	body := map[string]ProductFamily{
		"product_family": *family,
	}

	ret, err := makeCall(endpoints[endpointProductFamilyCreate], body, nil)
	if err != nil {
		return nil, err
	}
	// if successful, the product family should come back in a map[product_family]ProductFamily format
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["product_family"], family)
	return family, err
}

// GetProductFamily gets a product family
func GetProductFamilies() ([]ProductFamily, error) {
	found := []ProductFamily{}

	ret, err := makeCall(endpoints[endpointProductFamiliesGet], nil, &map[string]string{})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return found, err
	}

	temp := ret.Body.([]interface{})
	for i := range temp {
		entity := ProductFamily{}
		entry := temp[i].(map[string]interface{})
		entityRaw := entry["product_family"]
		err = mapstructure.Decode(entityRaw, &entity)
		if err != nil {
			return []ProductFamily{}, err
		}
		found = append(found, entity)
	}

	return found, err
}

// GetProductFamilyProducts gets products in a family
func GetProductFamilyComponents(id int64) ([]Component, error) {
	found := []Component{}

	ret, err := makeCall(endpoints[endpointProductFamilyComponentsGet], nil, &map[string]string{
		"product_family_id": fmt.Sprintf("%d", id),
	})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return found, err
	}

	temp := ret.Body.([]interface{})
	for i := range temp {
		entity := Component{}
		entry := temp[i].(map[string]interface{})
		entityRaw := entry["component"]
		err = mapstructure.Decode(entityRaw, &entity)
		if err != nil {
			return []Component{}, err
		}
		found = append(found, entity)
	}

	return found, err
}

// GetProductFamilyComponentByHandle gets components in a family
func GetProductFamilyComponentByHandle(familyID int64, handle string) (*Component, error) {

	ret, err := makeCall(endpoints[endpointProductFamilyComponentByHandleGet], nil, &map[string]string{
		"product_family_id": fmt.Sprintf("%d", familyID),
		"component_handle":  handle,
	})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return nil, err
	}

	entity := Component{}
	entry := ret.Body.(map[string]interface{})
	entityRaw := entry["component"]
	err = mapstructure.Decode(entityRaw, &entity)
	if err != nil {
		return nil, err
	}

	return &entity, err
}

// GetProductFamilyProducts gets products in a family
func GetProductFamilyComponentById(familyID int64, componentID int64) (*Component, error) {

	ret, err := makeCall(endpoints[endpointProductFamilyComponentByIdGet], nil, &map[string]string{
		"product_family_id": fmt.Sprintf("%d", familyID),
		"component_id":      fmt.Sprintf("%d", componentID),
	})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return nil, err
	}

	entity := Component{}
	entry := ret.Body.(map[string]interface{})
	entityRaw := entry["component"]
	err = mapstructure.Decode(entityRaw, &entity)
	if err != nil {
		return nil, err
	}

	return &entity, err
}

// GetProductFamilyProducts gets products in a family
func GetProductFamilyProducts(id int64) ([]Product, error) {
	found := []Product{}

	ret, err := makeCall(endpoints[endpointProductFamilyProductsGet], nil, &map[string]string{
		"id": fmt.Sprintf("%d", id),
	})
	if err != nil || ret.HTTPCode != http.StatusOK {
		return found, err
	}

	temp := ret.Body.([]interface{})
	for i := range temp {
		entity := Product{}
		entry := temp[i].(map[string]interface{})
		entityRaw := entry["product"]
		err = mapstructure.Decode(entityRaw, &entity)
		if err != nil {
			return []Product{}, err
		}
		found = append(found, entity)
	}

	return found, err
}

// GetProductFamily gets a product family
func GetProductFamily(productFamilyID int64) (*ProductFamily, error) {
	family := &ProductFamily{}
	ret, err := makeCall(endpoints[endpointProductFamilyGet], nil, &map[string]string{
		"id": fmt.Sprintf("%d", productFamilyID),
	})
	if err != nil {
		return nil, err
	}
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["product_family"], family)
	return family, err
}

// CreateProduct creates a new product and places the result in the input
func CreateProduct(productFamilyID int64, input *Product) error {
	if input.Name == "" || input.Handle == "" || input.Description == "" {
		return errors.New("name, handle, and description are required")
	}
	if input.PriceInCents <= 0 {
		return errors.New("price in cents must be greater than 0")
	}
	if input.IntervalUnit == "" || input.IntervalValue == 0 {
		return errors.New("interval and interval value must be provided")
	}
	body := map[string]Product{
		"product": *input,
	}

	ret, err := makeCall(endpoints[endpointProductCreate], body, &map[string]string{
		"familyID": fmt.Sprintf("%d", productFamilyID),
	})
	if err != nil {
		return err
	}
	// if successful, the product family should come back in a map[product_family]ProductFamily format
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["product"], input)
	return err
}

// GetProductByID gets a single product by id
func GetProductByID(productID int64) (*Product, error) {
	product := &Product{}
	ret, err := makeCall(endpoints[endpointProductGetByID], nil, &map[string]string{
		"id": fmt.Sprintf("%d", productID),
	})
	if err != nil {
		return nil, err
	}
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["product"], product)
	return product, err
}

// GetProductsInFamily gets all of the products in a family
func GetProductsInFamily(productFamilyID int64) ([]Product, error) {
	products := []Product{}
	ret, err := makeCall(endpoints[endpointProductGetForFamily], nil, &map[string]string{
		"familyID": fmt.Sprintf("%d", productFamilyID),
	})
	if err != nil {
		return nil, err
	}

	// so, Chargify violates OWASP best practices by returning these in an array
	temp := ret.Body.([]interface{})
	for i := range temp {
		entry := temp[i].(map[string]interface{})
		raw := entry["product"]
		product := Product{}
		err = mapstructure.Decode(raw, &product)
		if err == nil {
			products = append(products, product)
		}
	}
	return products, nil
}

// GetProductByHandle gets a product by its handle
func GetProductByHandle(handle string) (*Product, error) {
	product := &Product{}
	ret, err := makeCall(endpoints[endpointProductGetByHandle], nil, &map[string]string{
		"handle": handle,
	})
	if err != nil {
		return nil, err
	}
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["product"], product)
	return product, err
}

// UpdateProduct updates a product
func UpdateProduct(productID int64, input *Product) error {
	body := map[string]Product{
		"product": *input,
	}

	_, err := makeCall(endpoints[endpointProductUpdate], body, &map[string]string{
		"productID": fmt.Sprintf("%d", productID),
	})
	return err
}

// ArchiveProduct archives a product
func ArchiveProduct(productID int64) error {
	_, err := makeCall(endpoints[endpointProductArchive], nil, &map[string]string{
		"id": fmt.Sprintf("%d", productID),
	})
	return err
}

func createTestProductAndFamily() (*ProductFamily, *Product, error) {
	rand.Seed(time.Now().UnixNano())
	randID := rand.Int63()
	name := fmt.Sprintf("test-product-family-name-%d", randID)
	description := fmt.Sprintf("test-product-family-desc-%d", randID)
	handle := fmt.Sprintf("test-product-family-handle-%d", randID)
	acctCode := fmt.Sprintf("test-product-family-acct-%d", randID)
	family, err := CreateProductFamily(name, description, handle, acctCode)
	if err != nil {
		return nil, nil, err
	}

	trialPrice := 0
	trialIntervalValue := 90

	product := &Product{
		PriceInCents:       1000,
		Name:               fmt.Sprintf("Test Product-%d", randID),
		Handle:             fmt.Sprintf("test-product-handle-%d", randID),
		Description:        "Test product",
		IntervalUnit:       ProductIntervalDay,
		IntervalValue:      30,
		TrialPriceInCents:  &trialPrice,
		TrialIntervalValue: &trialIntervalValue,
		TrialIntervalUnit:  "day",
	}
	err = CreateProduct(family.ID, product)
	return family, product, err
}
