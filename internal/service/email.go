package service

import (
	"fmt"
	"net/smtp"
	"os"
)

// Simple SMTP-based email sender
func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("EMAIL_SENDER")       // your email
	password := os.Getenv("EMAIL_PASSWORD") // app password / SMTP password
	smtpHost := os.Getenv("SMTP_HOST")      // e.g., smtp.gmail.com
	smtpPort := os.Getenv("SMTP_PORT")      // e.g., 587

	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		fmt.Println("Failed to send email:", err)
	}
	return err
}
