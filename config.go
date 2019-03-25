package chargify

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type configDetails struct {
	environment string
	subdomain   string
	root        string
	apiKey      string
}

var config = &configDetails{}

func setup() (err error) {
	// setup the application; this is broken out from the init()
	// so that it may be called by unit tests to change env vars
	config.environment = strings.ToLower(envHelper("CHARGIFY_ENV", "production"))

	config.subdomain = strings.ToLower(envHelper("CHARGIFY_SUBDOMAIN", ""))
	if config.subdomain == "" {
		panic("subdomain not provided!")
	}
	config.apiKey = envHelper("CHARGIFY_API_KEY", "")
	if config.apiKey == "" {
		panic("apikey not provided!")
	}

	config.root = fmt.Sprintf("https://%s.chargify.com/", config.subdomain)
	return nil
}

func envHelper(variable, defaultValue string) string {
	found := os.Getenv(variable)
	if found == "" {
		found = defaultValue
	}
	return found
}

func init() {
	rand.Seed(time.Now().UnixNano())
	setup()
}
