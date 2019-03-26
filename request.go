package chargify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-resty/resty"
)

// APIReturn represents the return of the API calls
type APIReturn struct {
	StatusCode string      `json:"statusCode"`
	HTTPCode   int         `json:"httpCode"`
	Body       interface{} `json:"body"`
}

func makeCall(end endpoint, body interface{}, pathParams *map[string]string) (ret APIReturn, err error) {
	if config.subdomain == "" || config.apiKey == "" {
		return ret, errors.New("configuration is invalid for chargify")
	}
	endpointURI := end.uri
	if pathParams != nil {
		for k, v := range *pathParams {
			endpointURI = strings.Replace(endpointURI, "{"+k+"}", v, -1)
		}
	}
	url := fmt.Sprintf("%s%s", config.root, endpointURI)
	var response *resty.Response

	httpRequest := resty.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(config.apiKey, "x")

	if end.method == http.MethodGet {
		if body != nil {
			params, paramsOK := body.(map[string]string)
			if !paramsOK {
				return ret, errors.New("get calls must send in a map[string]string body")
			}
			httpRequest.SetQueryParams(params)
		}
		response, err = httpRequest.Get(url)
	} else if end.method == http.MethodPost {
		response, err = httpRequest.SetBody(body).
			Post(url)
	} else if end.method == http.MethodPut {
		response, err = httpRequest.SetBody(body).
			Put(url)
	} else if end.method == http.MethodDelete {
		response, err = httpRequest.Delete(url)
	}

	if err != nil {
		return
	}

	ret.HTTPCode = response.StatusCode()
	switch ret.HTTPCode {
	case http.StatusUnprocessableEntity:
		json.Unmarshal(response.Body(), &ret.Body)
		err = apiErrorToError(ret.Body)
	case http.StatusForbidden, http.StatusUnauthorized:
		err = errors.New("permission denied")
	case http.StatusNotFound:
		err = errors.New("not found")
	case http.StatusInternalServerError:
		err = errors.New("chargify server error")
	case http.StatusOK, http.StatusCreated:
		err = json.Unmarshal(response.Body(), &ret.Body)
		if err != nil {
			fmt.Printf("\nError: %+v\n", err)
			err = errors.New("could not unmarshal the JSON response; check the API")
			return
		}
	case http.StatusNoContent:
	default:
		fmt.Printf("Found unexpected: %d\n", ret.HTTPCode)
		fmt.Printf("\n%+v\n", ret.Body)
	}

	return
}

func convertStructToMap(i interface{}) (result map[string]string) {
	result = map[string]string{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		tag := typ.Field(i).Tag.Get("json")
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		// if it is blank, we don't send it
		if v != "" {
			result[tag] = v
		}
	}
	return
}

func apiErrorToError(input interface{}) error {
	// the body is likely a map of errors to []string
	errsI, errsOK := input.(map[string]interface{})
	if errsOK {
		errs := errsI["errors"].([]interface{})
		errorStrings := make([]string, len(errs))
		for i := range errs {
			errorStrings = append(errorStrings, errs[i].(string))
		}
		e := strings.Join(errorStrings, " ")
		return errors.New(e)
	}
	return errors.New("errors not found or not valid")
}

// ConvertJSONFloatToInt converts a float64 to an int64 from the JSON field interface
func ConvertJSONFloatToInt(input interface{}) (int64, error) {
	i, ok := input.(float64)
	if !ok {
		return 0, errors.New("could not convert")
	}
	return int64(i), nil
}
