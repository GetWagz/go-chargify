package chargify

import "net/http"

type endpoint struct {
	method     string
	uri        string
	pathParams []string
}

const (
	endpointBillingPortalEnable          = "billing_portal_enable"
	endpointBillingPortalEnableAndInvite = "billing_portal_enable_and_invite"
	endpointBillingPortalGet             = "billing_portal_get"

	endpointCouponCreate    = "coupon_create"
	endpointCouponGetByCode = "find_coupon"
	endpointCouponArchive   = "coupon_archive"

	endpointCustomerCreate            = "customer_create"
	endpointCustomerDelete            = "customer_delete"
	endpointCustomerUpdate            = "customer_update"
	endpointCustomersGet              = "customers_get"
	endpointCustomerGet               = "customer_get"
	endpointCustomerByReferenceGet    = "customer_get_by_reference"
	endpointCustomerSubscriptionsList = "customer_subscriptions_list"

	endpointEvents             = "events"
	endpointEventsCount        = "events_count"
	endpointEventIngestion     = "events_ingestion"
	endpointBulkEventIngestion = "events_bulk_ingestion"

	endpointPaymentProfileCreate = "payment_profile_create"
	endpointPaymentProfileDelete = "payment_profile_delete"
	endpointPaymentProfileUpdate = "payment_profile_update"

	endpointProductFamilyCreate               = "product_family_create"
	endpointProductFamilyGet                  = "product_family_get"
	endpointProductFamilyProductsGet          = "product_family_products_get"
	endpointProductFamiliesGet                = "product_families_get"
	endpointProductFamilyComponentsGet        = "product_family_components_get"
	endpointProductFamilyComponentByIdGet     = "product_family_component_by_id_get"
	endpointProductFamilyComponentByHandleGet = "product_family_component_by_handle_get"
	endpointProductCreate                     = "product_create"
	endpointProductUpdate                     = "product_update"
	endpointProductArchive                    = "product_archive"
	endpointProductGetByID                    = "product_get_by_id"
	endpointProductGetByHandle                = "product_get_by_handle"
	endpointProductGetForFamily               = "product_get_for_family"

	endpointSubscriptionCreate              = "subscription_create"
	endpointSubscriptionGet                 = "subscription_get"
	endpointSubscriptionUpdate              = "subscription_update"
	endpointSubscriptionGetMetaData         = "subscription_get_meta_data"
	endpointSubscriptionCancelImmediately   = "subscription_cancel_immediately"
	endpointSubscriptionCancelDelayed       = "subscription_cancel_delayed"
	endpointSubscriptionRemoveDelayedCancel = "subscription_remove_delayed_cancel"

	endpointSubscriptionMigrate   = "subscription_migrate"
	endpointSubscriptionUpdateNow = "subscription_update_now"

	endpointSubscriptionRefund           = "subscription_refund"
	endpointSubscriptionEvents           = "subscription_events"
	endpointSubscriptionComponentsUsages = "subscription_components_usages"

	endpointGetInvoices   = "invoices_get"
	endpointGetInvoice    = "invoice_get"
	endpointRefundInvoice = "invoice_refund"
)

var endpoints = map[string]endpoint{
	// customers
	endpointCustomerCreate: {
		method:     http.MethodPost,
		uri:        "customers",
		pathParams: []string{},
	},
	endpointCustomerDelete: {
		method: http.MethodDelete,
		uri:    "customers/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	endpointCustomerUpdate: {
		method: http.MethodPut,
		uri:    "customers/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	endpointCustomerGet: {
		method: http.MethodGet,
		uri:    "customers/{id}.json",
		pathParams: []string{
			"{id}",
		},
	},
	endpointCustomerByReferenceGet: {
		method: http.MethodGet,
		uri:    "customers/lookup.json?reference={reference}",
		pathParams: []string{
			"{reference}",
		},
	},
	endpointCustomerSubscriptionsList: {
		method: http.MethodGet,
		uri:    "customers/{customer_id}/subscriptions.json",
		pathParams: []string{
			"{customer_id}",
		},
	},

	endpointEvents: {
		method:     http.MethodGet,
		uri:        "events.json",
		pathParams: []string{},
	},
	endpointEventsCount: {
		method:     http.MethodGet,
		uri:        "events/count.json",
		pathParams: []string{},
	},
	endpointEventIngestion: {
		method: http.MethodPost,
		uri:    "events/{api_handle}.json",
		pathParams: []string{
			"{api_handle}",
		},
	},
	endpointBulkEventIngestion: {
		method: http.MethodPost,
		uri:    "events/{api_handle}/bulk.json",
		pathParams: []string{
			"{api_handle}",
		},
	},
	endpointCustomersGet: {
		method:     http.MethodGet,
		uri:        "customers",
		pathParams: []string{},
	},
	// billing portals
	endpointBillingPortalEnable: {
		method: http.MethodPost,
		uri:    "portal/customers/{id}/enable",
		pathParams: []string{
			"{id}",
		},
	},
	endpointBillingPortalEnableAndInvite: {
		method: http.MethodPost,
		uri:    "portal/customers/{id}/enable?invite=1",
		pathParams: []string{
			"{id}",
		},
	},
	endpointBillingPortalGet: {
		method: http.MethodGet,
		uri:    "portal/customers/{id}/management_link",
		pathParams: []string{
			"{id}",
		},
	},
	// payment profiles
	endpointPaymentProfileCreate: {
		method:     http.MethodPost,
		uri:        "payment_profiles",
		pathParams: []string{},
	},
	endpointPaymentProfileUpdate: {
		method: http.MethodPut,
		uri:    "payment_profiles/{paymentProfileID}",
		pathParams: []string{
			"{paymentProfileID}",
		},
	},
	endpointPaymentProfileDelete: {
		method: http.MethodDelete,
		uri:    "/subscriptions/{subscriptionID}/payment_profiles/{paymentProfileID}",
		pathParams: []string{
			"{subscriptionID}",
			"{paymentProfileID}",
		},
	},
	// product families
	endpointProductFamilyCreate: {
		method:     http.MethodPost,
		uri:        "product_families",
		pathParams: []string{},
	},
	endpointProductFamilyGet: {
		method: http.MethodGet,
		uri:    "product_families/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	// products
	endpointProductFamilyProductsGet: {
		method: http.MethodGet,
		uri:    "product_families/{id}/products.json",
		pathParams: []string{
			"{id}",
		},
	},
	endpointProductFamilyComponentsGet: {
		method: http.MethodGet,
		uri:    "product_families/{product_family_id}/components.json",
		pathParams: []string{
			"{product_family_id}",
		},
	},
	endpointProductFamilyComponentByIdGet: {
		method: http.MethodGet,
		uri:    "product_families/{product_family_id}/components/{component_id}.json",
		pathParams: []string{
			"{product_family_id}",
			"{component_id}",
		},
	},
	endpointProductFamilyComponentByHandleGet: {
		method: http.MethodGet,
		uri:    "product_families/{product_family_id}/components/handle:{component_handle}.json",
		pathParams: []string{
			"{product_family_id}",
			"{component_handle}",
		},
	},
	endpointProductFamiliesGet: {
		method:     http.MethodGet,
		uri:        "product_families.json",
		pathParams: []string{},
	},

	endpointProductCreate: {
		method: http.MethodPost,
		uri:    "product_families/{familyID}/products",
		pathParams: []string{
			"{familyID}",
		},
	},
	endpointProductUpdate: {
		method: http.MethodPut,
		uri:    "products/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	endpointProductArchive: {
		method: http.MethodDelete,
		uri:    "products/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	endpointProductGetByID: {
		method: http.MethodGet,
		uri:    "products/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	endpointProductGetByHandle: {
		method: http.MethodGet,
		uri:    "products/handle/{handle}",
		pathParams: []string{
			"{handle}",
		},
	},
	endpointProductGetForFamily: {
		method: http.MethodGet,
		uri:    "product_families/{familyID}/products",
		pathParams: []string{
			"{familyID}",
		},
	},
	// subscriptions
	endpointSubscriptionCreate: {
		method:     http.MethodPost,
		uri:        "subscriptions",
		pathParams: []string{},
	},
	endpointSubscriptionGet: {
		method: http.MethodGet,
		uri:    "subscriptions/{subscriptionID}",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionUpdate: {
		method: http.MethodPut,
		uri:    "subscriptions/{subscriptionID}",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionGetMetaData: {
		method: http.MethodGet,
		uri:    "subscriptions/{subscriptionID}/metadata",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionCancelImmediately: {
		method: http.MethodDelete,
		uri:    "subscriptions/{subscriptionID}",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionCancelDelayed: {
		method: http.MethodPost,
		uri:    "subscriptions/{subscriptionID}/delayed_cancel",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionRemoveDelayedCancel: {
		method: http.MethodDelete,
		uri:    "subscriptions/{subscriptionID}/delayed_cancel",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionMigrate: {
		method: http.MethodPost,
		uri:    "subscriptions/{subscriptionID}/migrations",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionRefund: {
		method: http.MethodPost,
		uri:    "subscriptions/{subscriptionID}/refunds",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionEvents: {
		method: http.MethodGet,
		uri:    "subscriptions/{subscriptionID}/events.json",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionComponentsUsages: {
		method: http.MethodPost,
		uri:    "subscriptions/{subscriptionID}/components/{componentID}}/usages.json",
		pathParams: []string{
			"{subscriptionID}",
			"{componentID}",
		},
	},

	// invoices
	endpointGetInvoices: {
		method:     http.MethodGet,
		uri:        "invoices",
		pathParams: []string{},
	},
	endpointGetInvoice: {
		method: http.MethodPost,
		uri:    "invoices/{invoiceID}",
		pathParams: []string{
			"{invoiceID}",
		},
	},
	endpointRefundInvoice: {
		method: http.MethodPost,
		uri:    "invoices/{invoiceID}/refunds",
		pathParams: []string{
			"{invoiceID}",
		},
	},
	// coupons
	endpointCouponCreate: {
		method: http.MethodPost,
		uri:    "product_families/{familyID}/coupons",
		pathParams: []string{
			"{familyID}",
		},
	},
	endpointCouponGetByCode: {
		method: http.MethodGet,
		uri:    "coupons/find?code={code}&product_family_id={familyID}",
		pathParams: []string{
			"{code}",
			"{familyID}",
		},
	},
	endpointCouponArchive: {
		method: http.MethodDelete,
		uri:    "product_families/{familyID}/coupons/{couponID}.json",
		pathParams: []string{
			"{familyID}",
			"{couponID}",
		},
	},
}
