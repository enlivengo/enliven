package enliven

import "github.com/hickeroar/enliven/config"

// DefaultEnlivenConfig represents the default config for enliven
var DefaultEnlivenConfig = config.Config{
	"email_smtp_identity":        "",
	"email_smtp_username":        "",
	"email_smtp_password":        "",
	"email_smtp_host":            "",
	"email_default_from_address": "",

	"server_address": ":8000",
}
