package internal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"
)

//PrettyJSON from object
func PrettyJSON(obj interface{}) string {
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

// JSON from object
func JSON(obj interface{}) string {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

// MergeStringToStringMap, last one wins on conflict
func MergeStringToStringMap(ms ...map[string]string) map[string]string {
	res := map[string]string{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

func ToMapStringToString(m map[string]interface{}) map[string]string {
	var found = make(map[string]string)

	for key, value := range m {
		if value == nil {
			continue
		}
		switch t := value.(type) {
		case *uint:
			found[key] = strconv.FormatUint(uint64(*t), 10)
		case *uint8:
			found[key] = strconv.FormatUint(uint64(*t), 10)
		case *int:
			found[key] = strconv.FormatInt(int64(*t), 10)
		case *int16:
			found[key] = strconv.FormatInt(int64(*t), 10)
		case *int32:
			found[key] = strconv.FormatInt(int64(*t), 10)
		case *int64:
			found[key] = strconv.FormatInt(int64(*t), 10)
		case *float32:
			found[key] = strconv.FormatFloat(float64(*t), 'f', 32, 32)
		case *float64:
			found[key] = strconv.FormatFloat(float64(*t), 'f', 64, 64)
		case *bool:
			found[key] = strconv.FormatBool(*t)
		case *string:
			found[key] = *t
		}

	}
	return found
}
func JoinUrls(basePath string, paths ...string) (*url.URL, error) {
	u, err := url.Parse(basePath)
	if err != nil {
		return nil, fmt.Errorf("invalid url")
	}
	p2 := append([]string{u.Path}, paths...)
	result := path.Join(p2...)
	u.Path = result
	return u, nil
}
