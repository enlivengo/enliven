package email

import (
	"errors"
	"net/smtp"

	"github.com/hickeroar/enliven/config"
)

// Core is the core functionality for sending emails.
type Core struct{}

// New creates a new self-contained email that can be sent to a user.
func (c Core) New() Email {
	if !c.Enabled() {
		panic("Email functionality has not been configured.")
	}
	conf := config.GetConfig()
	return Email{
		From: conf["email_default_from_address"],
	}
}

// Enabled returns whether or not we are configured for email
func (c Core) Enabled() bool {
	conf := config.GetConfig()
	if conf["smtp_host"] != "" {
		return true
	}
	return false
}

// Email represents an email that someone wants to send
type Email struct {
	To      []string
	From    string
	Subject string
	Message string
}

// AddRecipient appends an email address to the To slice
func (e *Email) AddRecipient(address string) {
	e.To = append(e.To, address)
}

// Send senss an email using smtp credentials provided in the config
func (e *Email) Send(email *Email) error {
	conf := config.GetConfig()

	if e.From == "" {
		return errors.New("Enliven Core Email: Unable to send email without 'From' address.")
	}
	if len(e.To) == 0 {
		return errors.New("Enliven Core Email: At least one 'To' address must be specified.")
	}

	auth := smtp.PlainAuth(conf["email_smtp_identity"], conf["email_smtp_username"], conf["email_smtp_password"], conf["email_smtp_host"])

	message := []byte("Subject: " + e.Subject + "\n\n" + e.Message + "\n")

	return smtp.SendMail(conf["email_smtp_host"]+":25", auth, e.From, e.To, message)
}
