package main

import (
	"bytes"
	"fmt"
	"net/smtp"
)

type Credentials struct {
	Email    string
	Password string
}

type Smtp struct {
	Host string
	Port string
}

type EmailOptions struct {
	Credentials Credentials
	Smtp        Smtp
	Subject     string
	Body        string
}

func sendEmail(opt EmailOptions, to []string) {
	// Sender data.
	from := opt.Credentials.Email
	password := opt.Credentials.Password

	// smtp server configuration.
	smtpHost := opt.Smtp.Host
	smtpPort := opt.Smtp.Port

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n%s", opt.Subject, mimeHeaders, opt.Body)))

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}
