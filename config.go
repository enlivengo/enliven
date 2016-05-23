package enliven

import "github.com/hickeroar/enliven/config"

// DefaultEnlivenConfig represents the default config for enliven
var DefaultEnlivenConfig = config.Config{
	"email_smtp_identity":  "",
	"email_smtp_username":  "",
	"email_smtp_password":  "",
	"email_smtp_host":      "",
	"email_smtp_port":      "25",
	"email_from_default":   "",
	"email_smtp_auth":      "plain", // "plain" and "none" are supported
	"email_allow_insecure": "0",

	"server_address": ":8000",

	"site_name": "Enliven",
	"site_url":  "http://localhost:8000",
}
