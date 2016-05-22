package email

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

// SendEmail sents an email using smtp credentials provided in the config
func SendEmail(email *Email) error {
	//config = ev.GetConfig()
	return nil
}
