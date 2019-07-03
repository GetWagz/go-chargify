module github.com/GetWagz/go-chargify

go 1.12

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-resty/resty v1.12.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/stretchr/testify v1.3.0
	golang.org/x/net v0.0.0-20190313082753-5c2c250b6a70 // indirect
	gopkg.in/resty.v1 v1.12.0 // indirect
)

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0
