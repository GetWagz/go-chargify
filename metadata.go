package chargify

// MetaDataEntry represents a single key/value meta data entry
type MetaDataEntry struct {
	Value      string `json:"value" mapstructure:"value"`
	ResourceID int64  `json:"resource_id" mapstructure:"resource_id"`
	Name       string `json:"name" mapstructure:"name"`
}

// MetaData represents a pageable return of a metadata request
type MetaData struct {
	TotalCount  int64           `json:"total_count" mapstructure:"total_count"`
	CurrentPage int64           `json:"current_page" mapstructure:"current_page"`
	TotalPages  int64           `json:"total_pages" mapstructure:"total_pages"`
	PerPage     int64           `json:"per_page" mapstructure:"per_page"`
	MetaData    []MetaDataEntry `json:"metadata" mapstructure:"metadata"`
}
