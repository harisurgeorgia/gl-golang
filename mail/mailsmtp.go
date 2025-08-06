package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendMail sends an email via Gmail SMTP
func SendMail(to, subject, body string) error {
	user := os.Getenv("GMAIL_USER")
	pass := os.Getenv("GMAIL_PASS")
	if user == "" || pass == "" {
		return fmt.Errorf("GMAIL_USER or GMAIL_PASS not set")
	}

	// Gmail SMTP server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-version: 1.0;\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
			"\r\n"+
			"%s\r\n",
		user, to, subject, body,
	))

	// Authentication.
	auth := smtp.PlainAuth("", user, pass, smtpHost)

	// Send email.
	addr := smtpHost + ":" + smtpPort
	return smtp.SendMail(addr, auth, user, []string{to}, msg)
}
