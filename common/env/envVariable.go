package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	AUTH_TOKEN_LENGTH int
	TOKEN_KEY_LENGTH int
	DB_URL string
	SMTP_EMAIL_HOST string
	SMTP_EMAIL string
	SMTP_PASSWORD string
	SMTP_PORT string
	SMTP_TLS string
	CLOUDINARY_URL string
	SQS_URL string
)

func LoadEnv() error {
	var err error
	err = godotenv.Load("../.env")
	if err != nil {
        return err
    }
	AUTH_TOKEN_LENGTH, err = strconv.Atoi(os.Getenv("AUTH_TOKEN_LENGTH"))
	if err != nil {
        return err
    }
	TOKEN_KEY_LENGTH, err = strconv.Atoi(os.Getenv("TOKEN_KEY_LENGTH"))
	DB_URL = os.Getenv("DB_URL")
	SMTP_EMAIL_HOST = os.Getenv("SMTP_EMAIL_HOST")
	SMTP_EMAIL = os.Getenv("SMTP_EMAIL")
	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	SMTP_PORT = os.Getenv("SMTP_PORT")
	SMTP_TLS = os.Getenv("SMTP_TLS")
	CLOUDINARY_URL = os.Getenv("CLOUDINARY_URL")
	SQS_URL = os.Getenv("SQS_URL")

	if err != nil {
        return err
    }
	log.Println("ENV loaded")
	return nil
}
