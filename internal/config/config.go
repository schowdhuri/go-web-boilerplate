package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	BrevoAPIKey     string
	BrevoSender     string
	DbHost          string
	DbUser          string
	DbPassword      string
	DbName          string
	DbPort          string
	GoEnv           string
	Port            string
	PublicUrl       string
	RecaptchaSecret string
}

func NewConfig() *EnvVars {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return &EnvVars{

		// brevo (SendInBlue)
		BrevoAPIKey: os.Getenv("BREVO_API_KEY"),
		BrevoSender: os.Getenv("BREVO_SENDER"),
		// database
		DbHost:     os.Getenv("DB_HOST"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		DbPort:     os.Getenv("DB_PORT"),
		GoEnv:      os.Getenv("GO_ENV"),
		Port:       os.Getenv("PORT"),
		PublicUrl:  os.Getenv("PUBLIC_URL"),
		// recaptcha
		RecaptchaSecret: os.Getenv("RECAPTCHA_SECRET"),
	}
}
