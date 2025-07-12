package routes

import (
	"simple-gin-backend/internal/controllers"
	"simple-gin-backend/internal/middleware"

	"github.com/gin-gonic/gin"

	docs "simple-gin-backend/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes sets up the application routes
func RegisterRoutes(router *gin.Engine) {
	// Swagger routes
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Hello world routes
	router.GET("/", controllers.GetHelloWorld)
	router.POST("/send-test-email", controllers.PostSendEmail)

	router.POST("/sign-up", controllers.RegisterUser)
	router.POST("/login", controllers.LoginUser)

	// Protected routes (Require authentication)
	api := router.Group("")
	api.Use(middleware.JWTAuthMiddleware())
	{
		// CRUD Routes for items
		api.POST("/items", controllers.CreateItem)
		api.GET("/items/:id", controllers.GetItem)
		api.PUT("/items/:id", controllers.UpdateItem)
		api.DELETE("/items/:id", controllers.DeleteItem)
		api.GET("/items", controllers.GetItems)
	}
}
