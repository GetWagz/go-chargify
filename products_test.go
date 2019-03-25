package chargify

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductFamilyCRUD(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	name := "Test Name"
	description := "This is a test product family that must be deleted manually"
	acctCode := fmt.Sprintf("test-acct-%d", rand.Int63())
	handle := fmt.Sprintf("test-handle-%d", rand.Int63())
	family, err := CreateProductFamily(name, description, handle, acctCode)
	require.Nil(t, err)
	require.NotZero(t, family.ID)
	assert.Equal(t, name, family.Name)
	assert.Equal(t, description, family.Description)

	found, err := GetProductFamily(family.ID)
	require.Nil(t, err)
	assert.Equal(t, name, found.Name)
	assert.Equal(t, description, found.Description)
}

func TestProductCRUD(t *testing.T) {
	randID := rand.Int63()
	family, err := CreateProductFamily("Test Name",
		"This is a test product family that must be deleted manually",
		fmt.Sprintf("test-acct-%d", randID),
		fmt.Sprintf("test-handle-%d", randID))
	require.Nil(t, err)

	product := Product{
		PriceInCents:  1000,
		Name:          fmt.Sprintf("Test Product-%d", randID),
		Handle:        fmt.Sprintf("test-product-handle-%d", randID),
		Description:   "Test product",
		IntervalUnit:  ProductIntervalDay,
		IntervalValue: 30,
	}
	err = CreateProduct(family.ID, &product)
	require.Nil(t, err)

	found, err := GetProductByID(product.ID)
	assert.Nil(t, err)
	require.NotNil(t, found)
	assert.Equal(t, product.ID, found.ID)
	assert.Equal(t, product.Name, found.Name)
	assert.Equal(t, product.Handle, found.Handle)
	assert.Equal(t, product.Description, found.Description)

	foundByHandle, err := GetProductByHandle(product.Handle)
	assert.Nil(t, err)
	require.NotNil(t, found)
	assert.Equal(t, product.ID, foundByHandle.ID)
	assert.Equal(t, product.Name, foundByHandle.Name)
	assert.Equal(t, product.Handle, foundByHandle.Handle)
	assert.Equal(t, product.Description, foundByHandle.Description)

	inFamily, err := GetProductsInFamily(family.ID)
	assert.Nil(t, err)
	assert.True(t, len(inFamily) > 0)

	err = ArchiveProduct(product.ID)
	require.Nil(t, err)
}
