package chargify

import (
	"errors"
	"fmt"

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

// Coupon is what we give back on the request
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

// CreateCoupon creates a new percent based coupon
func CreatePercentageCoupon(productFamilyID int64, input *PercentageCoupon) error {
	if input.Name == "" || input.Code == "" || input.Recurring == "" {
		return errors.New("name, code, and recurring are required")
	}
	if input.Percentage <= 0 {
		return errors.New("a value greater than 0 must be included for percentage")
	}

	input.ProductFamilyID = fmt.Sprintf("%d", productFamilyID)
	body := map[string]PercentageCoupon{
		"coupon": *input,
	}

	ret, err := makeCall(endpoints[endpointCouponCreate], body, &map[string]string{
		"familyID": fmt.Sprintf("%d", productFamilyID),
	})
	if err != nil {
		return err
	}

	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return errors.New("could not understand server response")
	}
	handleRet := PercentageCouponReturn{}
	err = mapstructure.Decode(apiBody["coupon"], &handleRet)
	return err
}

// CreateFlatCoupon creates a new flat rate coupon
func CreateFlatCoupon(productFamilyID int64, input *FlatCoupon) error {
	if input.Name == "" || input.Code == "" || input.Recurring == "" {
		return errors.New("name, code, and recurring are required")
	}
	if input.AmountInCents <= 0 {
		return errors.New("a value greater than 0 must be included for amount_in_cents")
	}

	input.ProductFamilyID = fmt.Sprintf("%d", productFamilyID)
	body := map[string]FlatCoupon{
		"coupon": *input,
	}

	ret, err := makeCall(endpoints[endpointCouponCreate], body, &map[string]string{
		"familyID": fmt.Sprintf("%d", productFamilyID),
	})
	if err != nil {
		return err
	}

	apiBody, bodyOK := ret.Body.(map[string]interface{})
	if !bodyOK {
		return errors.New("could not understand server response")
	}
	handleRet := FlatCouponReturn{}
	err = mapstructure.Decode(apiBody["coupon"], &handleRet)
	return err
}

// GetCouponByCode gets a coupon by its code
func GetCouponByCode(productFamilyID int64, code string) (*Coupon, error) {
	coupon := &Coupon{}
	ret, err := makeCall(endpoints[endpointCouponGetByCode], nil, &map[string]string{
		"familyID": fmt.Sprintf("%d", productFamilyID),
		"code":     fmt.Sprintf("%s", code),
	})
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
