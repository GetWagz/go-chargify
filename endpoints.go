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

	endpointCustomerCreate = "customer_create"
	endpointCustomerDelete = "customer_delete"
	endpointCustomersGet   = "customers_get"

	endpointPaymentProfileCreate = "payment_profile_create"
	endpointPaymentProfileDelete = "payment_profile_delete"
	endpointPaymentProfileUpdate = "payment_profile_update"

	endpointProductFamilyCreate = "product_family_create"
	endpointProductFamilyGet    = "product_family_get"

	endpointProductCreate       = "product_create"
	endpointProductUpdate       = "product_update"
	endpointProductArchive      = "product_archive"
	endpointProductGetByID      = "product_get_by_id"
	endpointProductGetByHandle  = "product_get_by_handle"
	endpointProductGetForFamily = "product_get_for_family"

	endpointSubscriptionCreate              = "subscription_create"
	endpointSubscriptionGet                 = "subscription_get"
	endpointSubscriptionGetMetaData         = "subscription_get_meta_data"
	endpointSubscriptionCancelImmediately   = "subscription_cancel_immediately"
	endpointSubscriptionCancelDelayed       = "subscription_cancel_delayed"
	endpointSubscriptionRemoveDelayedCancel = "subscription_remove_delayed_cancel"

	endpointSubscriptionMigrate   = "subscription_migrate"
	endpointSubscriptionUpdateNow = "subscription_update_now"

	endpointSubscriptionRefund = "subscription_refund"

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
}
