package chargify

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomerCreation(t *testing.T) {
	customID := rand.Int63n(999999999)
	input := Customer{
		FirstName: fmt.Sprintf("First-%d", customID),
		LastName:  fmt.Sprintf("Last-%d", customID),
		Email:     fmt.Sprintf("test+%d@example.com", customID),
		Reference: fmt.Sprintf("test-lib-%d", customID),
	}

	customer, err := CreateCustomer(&input)
	require.Nil(t, err)
	assert.Equal(t, input.Email, customer.Email)
	assert.Equal(t, input.Reference, customer.Reference)

	found, err := GetCustomers(1, "asc")
	assert.Nil(t, err)
	require.NotNil(t, found)

	updatedInput := Customer{
		ID:        customer.ID,
		FirstName: fmt.Sprintf("First-%d", customID),
		LastName:  fmt.Sprintf("Last-%d", customID),
		Email:     fmt.Sprintf("test+%d@example.com", customID),
		Reference: fmt.Sprintf("test-lib-%d", customID),
		Address:   "123 Main Street",
		City:      "Portsmouth",
		State:     "NH",
		Zip:       "03801",
	}

	err = UpdateCustomer(&updatedInput)
	assert.Nil(t, err)

	assert.Nil(t, err)
	assert.NotZero(t, len(found))

	foundByEmail, err := SearchForCustomersByEmail(input.Email)
	assert.Nil(t, err)
	assert.NotZero(t, len(foundByEmail))

	foundByReference, err := SearchForCustomerByReference(input.Reference)
	assert.Nil(t, err)
	assert.Equal(t, input.Email, foundByReference.Email)

	err = DeleteCustomerByID(customer.ID)
	assert.Nil(t, err)
}
