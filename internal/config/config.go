package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	RecaptchaSecret string
	BrevoAPIKey     string
	BrevoSender     string
	DbHost          string
	DbUser          string
	DbPassword      string
	DbName          string
	DbPort          string
	GoEnv           string
	Port            string
}

func NewConfig() *EnvVars {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return &EnvVars{
		// recaptcha
		RecaptchaSecret: os.Getenv("RECAPTCHA_SECRET"),
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
	}
}
