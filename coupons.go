package chargify

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Coupon struct {
	ID              int64  `json:"id"`
	Name            string `json:"name" mapstructure:"name"`                           //	The coupon name
	Code            string `json:"code" mapstructure:"code"`                           //	The coupon code
	Percentage      int    `json:"percentage" mapstructure:"percentage"`               //	The percentage value of the coupon
	AmountInCents   int64  `json:"amount_in_cents" mapstructure:"amount_in_cents"`     //	The amount_in_cents value of the coupon
	Recurring       string `json:"recurring" mapstructure:"recurring"`                 //	A string value for the boolean of whether or not this coupon is recurring
	ProductFamilyID string `json:"product_family_id" mapstructure:"product_family_id"` //	The id for the product family
}

// CreatePercentCoupon creates a new percent based coupon
func CreatePercentCoupon(productFamilyID int64, input *Coupon) error {
	if input.Name == "" || input.Code == "" || input.Recurring == "" {
		return errors.New("name, code, and percentage are required")
	}
	if input.Percentage <= 0 && input.AmountInCents <= 0 {
		return errors.New("a value greater than 0 must be included for either percentage or amount_in_cents")
	}
	body := map[string]Coupon{
		"coupon": *input,
	}

	ret, err := makeCall(endpoints[endpointCouponCreate], body, &map[string]string{
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
	err = mapstructure.Decode(apiBody["coupon"], input)
	return err
}

// GetProductByHandle gets a product by its handle
func GetCouponByCode(productFamilyID int64, code string) (*Coupon, error) {
	coupon := &Coupon{}
	ret, err := makeCall(endpoints[endpointProductGetByHandle], nil, &map[string]string{
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
func ArchiveCoupon(productFamilyID, couponId int64) error {
	_, err := makeCall(endpoints[endpointArchiveCoupon], nil, &map[string]string{
		"familyID": fmt.Sprintf("%d", productFamilyID),
		"couponId": fmt.Sprintf("%d", couponId),
	})
	return err
}
