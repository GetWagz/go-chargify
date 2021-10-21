package ingest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validJsonObject = []byte(`{
		"chargify": {
			"uniqueness_token": "1234"
		}
	}`)
	invalidJsonObject = []byte(`{
		"chargify-bad": {
			"uniqueness_token": "1234"
		}
	}`)
	invalidJsonObject2 = []byte(`{
		"chargify": {
			"uniqueness_token-bad": "1234"
		}
	}`)
	validJsonArray = []byte(`[{
		"chargify": {
			"uniqueness_token": "1234"
		}
	}]`)
	invalidJsonArray = []byte(`[{
		"chargify-bad": {
			"uniqueness_token": "1234"
		}
	}]`)
	invalidJsonArray2 = []byte(`[{
		"chargify": {
			"uniqueness_token-bad": "1234"
		}
	}]`)
)

func TestSingleJSONSchemaValidation(t *testing.T) {

	errs, err := jsonSchemaObject.ValidateBytes(context.TODO(), validJsonObject)
	assert.NoError(t, err)
	assert.True(t, len(errs) == 0)

	errs, err = jsonSchemaObject.ValidateBytes(context.TODO(), invalidJsonObject)
	assert.NoError(t, err)
	assert.True(t, len(errs) > 0)

	errs, err = jsonSchemaObject.ValidateBytes(context.TODO(), invalidJsonObject2)
	assert.NoError(t, err)
	assert.True(t, len(errs) > 0)

	errs, err = jsonSchemaObject.ValidateBytes(context.TODO(), validJsonArray)
	assert.NoError(t, err)
	assert.True(t, len(errs) > 0)

	errs, err = jsonSchemaArray.ValidateBytes(context.TODO(), invalidJsonObject2)
	assert.NoError(t, err)
	assert.True(t, len(errs) > 0)

	errs, err = jsonSchemaArray.ValidateBytes(context.TODO(), validJsonArray)
	assert.NoError(t, err)
	assert.True(t, len(errs) == 0)

	errs, err = jsonSchemaArray.ValidateBytes(context.TODO(), invalidJsonArray)
	assert.NoError(t, err)
	assert.True(t, len(errs) > 0)

	errs, err = jsonSchemaArray.ValidateBytes(context.TODO(), invalidJsonArray2)
	assert.NoError(t, err)
	assert.True(t, len(errs) > 0)
}
