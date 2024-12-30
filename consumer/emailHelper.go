package main

import (
	"log"
	"net/smtp"
	"common/env"
)

var smtp_auth smtp.Auth

func InitSMTP() {
	smtp_auth = smtp.PlainAuth("", env.SMTP_EMAIL, env.SMTP_PASSWORD, env.SMTP_EMAIL_HOST)
	log.Println("SMTP Loaded")
}

func SendEmail(to string, subject string, msg string) error {
	emailTo := []string{to}
	emailMsg := []byte("To: "+to +"\r\n" + "Subject: "+subject+"\n"+ "\r\n" + msg)

	err := smtp.SendMail(env.SMTP_EMAIL_HOST+":"+env.SMTP_PORT, smtp_auth, env.SMTP_EMAIL, emailTo, emailMsg)

	return err
}
