package chargify

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentProfileCreation(t *testing.T) {
	customID := rand.Int63n(999999999)
	input := Customer{
		FirstName: fmt.Sprintf("First-%d", customID),
		LastName:  fmt.Sprintf("Last-%d", customID),
		Email:     fmt.Sprintf("test+%d@example.com", customID),
		Reference: fmt.Sprintf("test-lib-%d", customID),
	}

	customer, err := CreateCustomer(&input)
	require.Nil(t, err)
	defer DeleteCustomerByID(customer.ID)

	profile, err := SavePaymentProfileVault(customer.ID, VaultBogus, "12345")
	assert.Nil(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, VaultBogus, profile.CurrentVault)
	assert.Equal(t, "12345", profile.VaultToken)
	assert.NotZero(t, profile.ID)

	DeletePaymentProfile(profile.SubscriptionID, profile.ID)
}
