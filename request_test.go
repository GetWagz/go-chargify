package chargify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {
	cust := Customer{
		FirstName: "Kevin",
		LastName:  "Eaton",
		Email:     "kevin@wagz.com",
	}

	result := convertStructToMap(&cust)
	assert.Equal(t, "Kevin", result["first_name"])
	assert.Equal(t, "Eaton", result["last_name"])
	_, foundAddressOK := result["address"]
	assert.False(t, foundAddressOK)
}
