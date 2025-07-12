package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ApplicationEnvironment string
	IsDevelopmentEnv       bool

	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDb        string
	PostgresPort     string

	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
	EmailFromEmail string
	UseMailTLS     bool
}

// AppConfig holds the global configuration accessible by the entire application
var AppConfig Config

// LoadConfig loads environment variables into the AppConfig struct
func LoadConfig() {

	// Load configuration from environment variables
	AppConfig.ApplicationEnvironment = os.Getenv("APPLICATION_ENVIRONMENT")
	AppConfig.IsDevelopmentEnv = AppConfig.ApplicationEnvironment == "development"

	AppConfig.PostgresHost = os.Getenv("POSTGRES_HOST")
	AppConfig.PostgresUser = os.Getenv("POSTGRES_USER")
	AppConfig.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	AppConfig.PostgresDb = os.Getenv("POSTGRES_DB")
	AppConfig.PostgresPort = os.Getenv("POSTGRES_PORT")

	AppConfig.SMTPHost = os.Getenv("SMTP_HOST")
	AppConfig.SMTPUser = os.Getenv("SMTP_USER")
	AppConfig.SMTPPassword = os.Getenv("SMTP_PASSWORD")
	AppConfig.EmailFromEmail = os.Getenv("EMAILS_FROM_EMAIL")

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("Error parsing SMTP_PORT: %v", err)
	}
	AppConfig.SMTPPort = port

	// Convert MAIL_TLS to a boolean value
	tlsEnabled, err := strconv.ParseBool(os.Getenv("MAIL_TLS"))
	if err != nil {
		tlsEnabled = true // Default safe value
	}
	AppConfig.UseMailTLS = tlsEnabled
}
