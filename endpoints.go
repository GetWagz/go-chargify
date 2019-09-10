package chargify

import "net/http"

type endpoint struct {
	method     string
	uri        string
	pathParams []string
}

const (
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
)

var endpoints = map[string]endpoint{
	endpointCustomerCreate: endpoint{
		method:     http.MethodPost,
		uri:        "customers",
		pathParams: []string{},
	},
	endpointCustomerDelete: endpoint{
		method: http.MethodDelete,
		uri:    "customers/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	endpointCustomersGet: endpoint{
		method:     http.MethodGet,
		uri:        "customers",
		pathParams: []string{},
	},
	// payment profiles
	endpointPaymentProfileCreate: endpoint{
		method:     http.MethodPost,
		uri:        "payment_profiles",
		pathParams: []string{},
	},
	endpointPaymentProfileUpdate: endpoint{
		method: http.MethodPut,
		uri:    "payment_profiles/{paymentProfileID}",
		pathParams: []string{
			"{paymentProfileID}",
		},
	},
	endpointPaymentProfileDelete: endpoint{
		method: http.MethodDelete,
		uri:    "/subscriptions/{subscriptionID}/payment_profiles/{paymentProfileID}",
		pathParams: []string{
			"{subscriptionID}",
			"{paymentProfileID}",
		},
	},
	// product families
	endpointProductFamilyCreate: endpoint{
		method:     http.MethodPost,
		uri:        "product_families",
		pathParams: []string{},
	},
	endpointProductFamilyGet: endpoint{
		method: http.MethodGet,
		uri:    "product_families/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	// products
	endpointProductCreate: endpoint{
		method: http.MethodPost,
		uri:    "product_families/{familyID}/products",
		pathParams: []string{
			"{familyID}",
		},
	},
	endpointProductUpdate: endpoint{
		method: http.MethodPut,
		uri:    "products/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	endpointProductArchive: endpoint{
		method: http.MethodDelete,
		uri:    "products/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	endpointProductGetByID: endpoint{
		method: http.MethodGet,
		uri:    "products/{id}",
		pathParams: []string{
			"{id}",
		},
	},
	endpointProductGetByHandle: endpoint{
		method: http.MethodGet,
		uri:    "products/handle/{handle}",
		pathParams: []string{
			"{handle}",
		},
	},
	endpointProductGetForFamily: endpoint{
		method: http.MethodGet,
		uri:    "product_families/{familyID}/products",
		pathParams: []string{
			"{familyID}",
		},
	},
	// subscriptions
	endpointSubscriptionCreate: endpoint{
		method:     http.MethodPost,
		uri:        "subscriptions",
		pathParams: []string{},
	},
	endpointSubscriptionGet: endpoint{
		method: http.MethodGet,
		uri:    "subscriptions/{subscriptionID}",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionGetMetaData: endpoint{
		method: http.MethodGet,
		uri:    "subscriptions/{subscriptionID}/metadata",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionCancelImmediately: endpoint{
		method: http.MethodDelete,
		uri:    "subscriptions/{subscriptionID}",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionCancelDelayed: endpoint{
		method: http.MethodPost,
		uri:    "subscriptions/{subscriptionID}/delayed_cancel",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionRemoveDelayedCancel: endpoint{
		method: http.MethodDelete,
		uri:    "subscriptions/{subscriptionID}/delayed_cancel",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
	endpointSubscriptionMigrate: endpoint{
		method: http.MethodPost,
		uri:    "subscriptions/{subscriptionID}/migrations",
		pathParams: []string{
			"{subscriptionID}",
		},
	},
}
