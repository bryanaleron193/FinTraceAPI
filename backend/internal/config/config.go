package config

import (
	"os"
)

type Config struct {
	ApplicationEnvironment string
	IsDevelopmentEnv       bool

	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
	PostgresPort     string

	RedisAddr     string
	RedisPassword string
	RedisDB       string

	ResendAPIKey    string
	ResendFromName  string
	ResendFromEmail string

	AdminEmail string
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

	AppConfig.RedisAddr = os.Getenv("REDIS_ADDR")
	AppConfig.RedisPassword = os.Getenv("REDIS_PASSWORD")
	AppConfig.RedisDB = os.Getenv("REDIS_DB")

	AppConfig.ResendAPIKey = os.Getenv("RESEND_API_KEY")
	AppConfig.ResendFromName = os.Getenv("RESEND_FROM_NAME")
	AppConfig.ResendFromEmail = os.Getenv("RESEND_FROM_EMAIL")

	AppConfig.AdminEmail = os.Getenv("ADMIN_EMAIL")
}
