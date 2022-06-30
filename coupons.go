package chargify

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// PercentageCoupon is the structure of what we need to send to Chargify when creating a percentage-based coupon
type PercentageCoupon struct {
	Name            string `json:"name" mapstructure:"name"`                           //	The coupon name
	Code            string `json:"code" mapstructure:"code"`                           //	The coupon code
	Description     string `json:"description" mapstructure:"description"`             //   The (optionally required?) description for the coupon
	Percentage      int    `json:"percentage" mapstructure:"percentage"`               //	The percentage value of the coupon
	Recurring       string `json:"recurring" mapstructure:"recurring"`                 //	A string value for the boolean of whether or not this coupon is recurring
	ProductFamilyID string `json:"product_family_id" mapstructure:"product_family_id"` //	The id for the product family
}

// PercentageCouponReturn is here because Chargify send back different types than what it asks for
type PercentageCouponReturn struct {
	Name            string  `json:"name" mapstructure:"name"`                           //	The coupon name
	Code            string  `json:"code" mapstructure:"code"`                           //	The coupon code
	Description     string  `json:"description" mapstructure:"description"`             //  The (optionally required?) description for the coupon
	Percentage      string  `json:"percentage" mapstructure:"percentage"`               //	The percentage value of the coupon
	Recurring       bool    `json:"recurring" mapstructure:"recurring"`                 //	A boolean of whether or not this coupon is recurring
	ProductFamilyID float64 `json:"product_family_id" mapstructure:"product_family_id"` //	The id for the product family
}

// FlatCoupon is the structure of what we need to send to Chargify when creating a flat rate coupon
type FlatCoupon struct {
	Name            string `json:"name" mapstructure:"name"`                           //	The coupon name
	Code            string `json:"code" mapstructure:"code"`                           //	The coupon code
	Description     string `json:"description" mapstructure:"description"`             //   The (optionally required?) description for the coupon
	AmountInCents   int64  `json:"amount_in_cents" mapstructure:"amount_in_cents"`     //	The amount_in_cents value of the coupon
	Recurring       string `json:"recurring" mapstructure:"recurring"`                 //	A string value for the boolean of whether or not this coupon is recurring
	ProductFamilyID string `json:"product_family_id" mapstructure:"product_family_id"` //	The id for the product family
}

// FlatCouponReturn is here because Chargify send back different types than what it asks for
type FlatCouponReturn struct {
	Name            string  `json:"name" mapstructure:"name"`                           //	The coupon name
	Code            string  `json:"code" mapstructure:"code"`                           //	The coupon code
	Description     string  `json:"description" mapstructure:"description"`             //  The (optionally required?) description for the coupon
	AmountInCents   int64   `json:"amount_in_cents" mapstructure:"amount_in_cents"`     //	The amount_in_cents value of the coupon
	Recurring       bool    `json:"recurring" mapstructure:"recurring"`                 //	A boolean of whether or not this coupon is recurring
	ProductFamilyID float64 `json:"product_family_id" mapstructure:"product_family_id"` //	The id for the product family
}

// Coupon is a coupons structure
type Coupon struct {
	ID              int64  `json:"id"`
	Name            string `json:"name" mapstructure:"name"`                           //	The coupon name
	Code            string `json:"code" mapstructure:"code"`                           //	The coupon code
	Description     string `json:"description" mapstructure:"description"`             // The (optionally required?) description for the coupon
	Percentage      int    `json:"percentage" mapstructure:"percentage"`               //	The percentage value of the coupon
	AmountInCents   int64  `json:"amount_in_cents" mapstructure:"amount_in_cents"`     //	The amount_in_cents value of the coupon
	Recurring       string `json:"recurring" mapstructure:"recurring"`                 //	A string value for the boolean of whether or not this coupon is recurring
	ProductFamilyID string `json:"product_family_id" mapstructure:"product_family_id"` //	The id for the product family
}

// CouponReturn is what we give back on the request
type CouponReturn struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name" mapstructure:"name"`                           //	The coupon name
	Code            string  `json:"code" mapstructure:"code"`                           //	The coupon code
	Description     string  `json:"description" mapstructure:"description"`             // The (optionally required?) description for the coupon
	Percentage      string  `json:"percentage" mapstructure:"percentage"`               //	The percentage value of the coupon
	AmountInCents   int64   `json:"amount_in_cents" mapstructure:"amount_in_cents"`     //	The amount_in_cents value of the coupon
	Recurring       bool    `json:"recurring" mapstructure:"recurring"`                 //	A string value for the boolean of whether or not this coupon is recurring
	ProductFamilyID float64 `json:"product_family_id" mapstructure:"product_family_id"` //	The id for the product family
}

// ListCouponsQueryParams are the query params for coupon listing. Note that a lot of the fields were deprecated as they now prefer
// a query string array of filter indexes, but we expose the higher-level fields and consolidate it on the call
type ListCouponsQueryParams struct {
	CurrencyPrices *bool    `json:"currency_prices"`
	Page           *int     `json:"page"`
	PerPage        *int     `json:"per_page"`
	Codes          []string `json:"codes"`
	DateField      *string  `json:"date_field"`
	EndDate        *string  `json:"end_date"`
	EndDateTime    *string  `json:"end_datetime"`
	IDs            []string `json:"ids"`
	StartDate      *string  `json:"start_date"`
	StartDateTime  *string  `json:"start_date_time"`
}

func (input *ListCouponsQueryParams) toMap() *map[string]string {
	m := map[string]string{}
	if input.CurrencyPrices != nil {
		m["currency_prices"] = fmt.Sprintf("%v", ToBool(input.CurrencyPrices))
	}
	if input.Page != nil {
		m["page"] = fmt.Sprintf("%d", ToInt(input.Page))
	}
	if input.PerPage != nil {
		m["per_page"] = fmt.Sprintf("%d", ToInt(input.PerPage))
	}
	if len(input.Codes) > 0 {
		m["filter[codes]"] = strings.Join(input.Codes, ",")
	}
	if len(input.IDs) > 0 {
		m["filter[ids]"] = strings.Join(input.IDs, ",")
	}
	if input.DateField != nil {
		m["filter[date_field]"] = ToString(input.DateField)
	}
	if input.EndDate != nil {
		m["filter[end_date]"] = ToString(input.EndDate)
	}
	if input.EndDateTime != nil {
		m["filter[end_datetime]"] = ToString(input.EndDateTime)
	}
	if input.StartDate != nil {
		m["filter[start_date]"] = ToString(input.StartDate)
	}
	if input.StartDateTime != nil {
		m["filter[start_date_time]"] = ToString(input.StartDateTime)
	}
	return &m
}

// CreateCoupon creates a new percent based coupon
func CreatePercentageCoupon(productFamilyID int64, input *PercentageCoupon) (*PercentageCouponReturn, error) {
	handleRet := PercentageCouponReturn{}
	if input.Name == "" || input.Code == "" || input.Recurring == "" {
		return &handleRet, errors.New("name, code, and recurring are required")
	}
	if input.Percentage <= 0 {
		return &handleRet, errors.New("a value greater than 0 must be included for percentage")
	}

	input.ProductFamilyID = fmt.Sprintf("%d", productFamilyID)
	body := map[string]PercentageCoupon{
		"coupon": *input,
	}

	ret, err := makeCall(endpoints[endpointCouponCreate], body, &map[string]string{
		"familyID": fmt.Sprintf("%d", productFamilyID),
	})
	if err != nil {
		return &handleRet, err
	}

	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return &handleRet, errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["coupon"], &handleRet)
	return &handleRet, err
}

// CreateFlatCoupon creates a new flat rate coupon
func CreateFlatCoupon(productFamilyID int64, input *FlatCoupon) (*FlatCouponReturn, error) {
	handleRet := FlatCouponReturn{}
	if input.Name == "" || input.Code == "" || input.Recurring == "" {
		return &handleRet, errors.New("name, code, and recurring are required")
	}
	if input.AmountInCents <= 0 {
		return &handleRet, errors.New("a value greater than 0 must be included for amount_in_cents")
	}

	input.ProductFamilyID = fmt.Sprintf("%d", productFamilyID)
	body := map[string]FlatCoupon{
		"coupon": *input,
	}

	ret, err := makeCall(endpoints[endpointCouponCreate], body, &map[string]string{
		"familyID": fmt.Sprintf("%d", productFamilyID),
	})
	if err != nil {
		return &handleRet, err
	}

	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return &handleRet, errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["coupon"], &handleRet)
	return &handleRet, err
}

// GetCouponByCode gets a coupon by its code
func GetCouponByCode(productFamilyID int64, code string) (*CouponReturn, error) {
	coupon := &CouponReturn{}
	ret, err := makeCall(endpoints[endpointCouponGetByCode], map[string]string{
		"familyID": fmt.Sprintf("%d", productFamilyID),
		"code":     code,
	}, nil)
	if err != nil {
		return nil, err
	}
	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	err = mapstructure.Decode(apiBody["coupon"], coupon)
	return coupon, err
}

// ArchiveCoupon archives a coupon on use or expiration
func ArchiveCoupon(productFamilyID, couponID int64) error {
	_, err := makeCall(endpoints[endpointCouponArchive], nil, &map[string]string{
		"familyID": fmt.Sprintf("%d", productFamilyID),
		"couponID": fmt.Sprintf("%d", couponID),
	})
	return err
}

// ListCoupons lists out the coupons based upon the result of the passed in query params
func ListCoupons(params *ListCouponsQueryParams) ([]CouponReturn, error) {
	if params == nil {
		params = &ListCouponsQueryParams{}
	}
	queryParams := params.toMap()
	options := &makeCallOptions{
		End:         endpoints[endpointCouponsList],
		QueryParams: queryParams,
	}

	data := []CouponReturn{}

	ret, err := makeAPICall(options)
	if err != nil {
		return data, err
	}
	apiBody, bodyOK := ret.Body.([]interface{})
	if !bodyOK {
		return nil, errors.New("could not understand server response")
	}
	// the result is an array of objects that have a coupon key, similar to:
	// 	[
	//   {
	//     "coupon": {...}
	//   }
	// ]
	// so we have a [] of interfaces; loop and convert
	for i := range apiBody {
		if raw, rawOK := apiBody[i].(map[string]interface{}); rawOK {
			coupon := CouponReturn{}
			err = mapstructure.Decode(raw["coupon"], &coupon)
			if err == nil {
				data = append(data, coupon)
			}
		}

	}
	return data, nil
}
