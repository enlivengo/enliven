package email

import (
	"errors"
	"fmt"
	"net/smtp"

	"github.com/enlivengo/enliven/config"
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
		From: conf["email_from_default"],
	}
}

// Enabled returns whether or not we are configured for email
func (c Core) Enabled() bool {
	conf := config.GetConfig()
	if conf["email_smtp_host"] != "" {
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
func (e *Email) Send() error {
	conf := config.GetConfig()

	if e.From == "" {
		return errors.New("Enliven Core Email: Unable to send email without 'From' address.")
	}
	if len(e.To) == 0 {
		return errors.New("Enliven Core Email: At least one 'To' address must be specified.")
	}

	// The smtp server may have no auth mechanism
	if conf["email_smtp_auth"] == "none" {
		// Opening connection to smtp server
		smtpConn, err := smtp.Dial(conf["email_smtp_host"] + ":" + conf["email_smtp_port"])
		if err != nil {
			return err
		}

		// Adding from address
		if err := smtpConn.Mail(e.From); err != nil {
			return err
		}
		// Adding all recipients to the message
		for _, recipient := range e.To {
			if err := smtpConn.Rcpt(recipient); err != nil {
				return err
			}
		}

		// Writing out the message data
		messageWriter, err := smtpConn.Data()
		if err != nil {
			return err
		}
		if _, err = fmt.Fprintf(messageWriter, "From: "+e.From+"\nSubject: "+e.Subject+"\n\n"+e.Message+"\n"); err != nil {
			return err
		}
		if err = messageWriter.Close(); err != nil {
			return err
		}

		// Closing connection to the smtp server
		if err = smtpConn.Quit(); err != nil {
			return err
		}

		return nil
	}

	auth := smtp.PlainAuth(conf["email_smtp_identity"], conf["email_smtp_username"], conf["email_smtp_password"], conf["email_smtp_host"])
	message := []byte("From: " + e.From + "\nSubject: " + e.Subject + "\r\n\r\n" + e.Message + "\r\n")
	err := smtp.SendMail(conf["email_smtp_host"]+":"+conf["email_smtp_port"], auth, e.From, e.To, message)

	// If we failed with encryption error, and the setting for insecurity is allowed, we insecure send it (recommended only for testing)
	if err != nil && err.Error() == "unencrypted connection" && conf["email_allow_insecure"] == "1" {
		uAuth := unencryptedAuth{auth}
		err = smtp.SendMail(conf["email_smtp_host"]+":"+conf["email_smtp_port"], uAuth, e.From, e.To, message)
	}

	if err != nil {
		fmt.Println(err)
	}

	return err
}

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}
