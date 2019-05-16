package chargify

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionCRUD(t *testing.T) {
	customer, _, err := createTestCustomer()
	require.Nil(t, err)
	// defer DeleteCustomerByID(customer.ID)

	_, product, err := createTestProductAndFamily()
	require.Nil(t, err)
	// defer ArchiveProduct(product.ID)
	// note for the programmer: there is no way to archive a product family at this point
	profileID := int64(0)

	subscription, err := CreateSubscriptionForCustomer(customer.Reference, product.Handle, profileID, nil)
	require.Nil(t, err)
	assert.NotZero(t, subscription.ID)

	err = CancelSubscription(subscription.ID, false, "MY_REASON", "Testing")
	assert.Nil(t, err)

	err = RemoveDelayedSubscriptionCancellation(subscription.ID)
	assert.Nil(t, err)

	err = CancelSubscription(subscription.ID, true, "MY_REASON", "Testing")
	assert.Nil(t, err)
}
