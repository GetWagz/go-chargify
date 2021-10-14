package chargify

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCouponCRD(t *testing.T) {
	customID := rand.Int63n(999999999)
	input := PercentageCoupon{
		Name:            fmt.Sprintf("Name-%d", customID),
		Code:            fmt.Sprintf("CODE%d", customID),
		Description:     fmt.Sprintf("test-description+%d", customID),
		Percentage:      25,
		Recurring:       "false",
		ProductFamilyID: fmt.Sprintf("%d", 1182341),
	}

	err := CreatePercentageCoupon(1182341, &input)
	require.Nil(t, err)

	found, err := GetCouponByCode(1182341, input.Code)

	err = ArchiveCoupon(1182341, found.ID)
	assert.Nil(t, err)

	inputFlat := FlatCoupon{
		Name:            fmt.Sprintf("Name-%d", customID),
		Code:            fmt.Sprintf("C0DE%d", customID),
		Description:     fmt.Sprintf("test-description+%d", customID),
		AmountInCents:   500,
		Recurring:       "false",
		ProductFamilyID: fmt.Sprintf("%d", 1182341),
	}

	err = CreateFlatCoupon(1182341, &inputFlat)
	require.Nil(t, err)

	found, err = GetCouponByCode(1182341, inputFlat.Code)

	err = ArchiveCoupon(1182341, found.ID)
	assert.Nil(t, err)
}
