package chargify

// MetaDataEntry represents a single key/value meta data entry
type MetaDataEntry struct {
	Value      string `json:"value"`
	ResourceID int64  `json:"resource_id"`
	Name       string `json:"name"`
}

// MetaData represents a pageable return of a metadata request
type MetaData struct {
	TotalCount  int64           `json:"total_count"`
	CurrentPage int64           `json:"current_page"`
	TotalPages  int64           `json:"total_pages"`
	PerPage     int64           `json:"per_page"`
	MetaData    []MetaDataEntry `json:"metadata"`
}
