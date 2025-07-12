package main

import (
	"simple-gin-backend/internal/config"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/routes"

	"github.com/gin-gonic/gin"
)

// @title Simple Gin Backend API
// @version 1.0
// @description This is a simple backend using Gin and GORM.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load the configuration from environment variables
	config.LoadConfig()

	// Initialize the database
	database.InitDB()
	// Auto-migrate the database schema
	database.AutoMigrate()

	r := gin.Default()

	// Set up routes
	routes.RegisterRoutes(r)

	// Run the server
	r.Run(":8080")
}
