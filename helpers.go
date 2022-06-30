package chargify

// these helpers are designed to make query params a little more manageable
// for inclsion and is modeled off of libraries such as Stripes. These each to have a FromX
// and a ToX funcs

// FromString converts a value to a pointer
func FromString(input string) *string {
	return &input
}

// ToString takes a pointer and gives the value
func ToString(input *string) string {
	return *input
}

// FromInt64 converts a value to a pointer
func FromInt64(input int64) *int64 {
	return &input
}

// ToInt64 takes a pointer and gives the value
func ToInt64(input *int64) int64 {
	return *input
}

// FromInt converts a value to a pointer
func FromInt(input int) *int {
	return &input
}

// ToInt takes a pointer and gives the value
func ToInt(input *int) int {
	return *input
}

// FromFloat64 converts a value to a pointer
func FromFloat64(input float64) *float64 {
	return &input
}

// ToFloat64 takes a pointer and gives the value
func ToFloat64(input *float64) float64 {
	return *input
}

// FromBool converts a value to a pointer
func FromBool(input bool) *bool {
	return &input
}

// ToBool takes a pointer and gives the value
func ToBool(input *bool) bool {
	return *input
}
