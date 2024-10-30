package common

import (
	"log"
	"net/smtp"
	"todoGoApi/utils"
)

var smtp_auth smtp.Auth

func InitSMTP() {
	smtp_auth = smtp.PlainAuth("", utils.SMTP_EMAIL, utils.SMTP_PASSWORD, utils.SMTP_EMAIL_HOST)
	log.Println("SMTP Loaded")
}

func SendEmail(to string, subject string, msg string) error {
	emailTo := []string{to}
	emailMsg := []byte("To: "+to +"\r\n" + "Subject: "+subject+"\n"+ "\r\n" + msg)

	err := smtp.SendMail(utils.SMTP_EMAIL_HOST+":"+utils.SMTP_PORT, smtp_auth, utils.SMTP_EMAIL, emailTo, emailMsg)

	return err
}
