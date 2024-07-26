// email.go
package main

import (
	"fmt"
	"net/smtp"
)

func sendVerificationEmail(email, code string) error {
	from := "your_email@example.com"
	password := "your_email_password"

	to := []string{email}
	subject := "Verification Code"
	body := fmt.Sprintf("Your verification code is: %s", code)
	message := []byte("Subject: " + subject + "\r\n" +
		"To: " + email + "\r\n" +
		"From: " + from + "\r\n" +
		"\r\n" + body)

	auth := smtp.PlainAuth("", from, password, "smtp.example.com")

	err := smtp.SendMail("smtp.example.com:587", auth, from, to, message)
	if err != nil {
		return err
	}
	return nil
}
