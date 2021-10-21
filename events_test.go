package chargify

import (
	"fmt"
	"testing"

	"github.com/GetWagz/go-chargify/internal"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
)

func TestListEventParams(t *testing.T) {
	page := 1
	direction := "asc"
	queryParams := ListEventsQueryParams{
		Page:      &page,
		Direction: &direction,
	}
	structs.DefaultTagName = "mapstructure"
	m := structs.Map(queryParams)
	_, ok := m["page"]
	assert.True(t, ok)

	f := internal.ToMapStringToString(m)
	fmt.Println(internal.PrettyJSON(f))
	_, ok = m["per_page"]
	assert.False(t, ok)

}
