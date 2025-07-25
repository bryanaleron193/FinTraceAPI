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

	router.POST("/login", controllers.LoginUser)

	// Protected routes (Require authentication)
	api := router.Group("")
	api.Use(middleware.JWTAuthMiddleware())
	{
		api.GET("/users/get-all-statuses", controllers.GetAllUserStatuses)
		api.GET("/users/get-all-users", controllers.GetAllUsers)
		api.PUT("/users/update-status", controllers.UpdateUserStatus)

		api.GET("/groups/joined", controllers.GetAllGroupsByUserId)
		api.GET("/groups/not-joined", controllers.GetAllGroupsNotJoined)
		api.POST("/groups/create", controllers.CreateGroup)
		api.PUT("/groups/update", controllers.UpdateGroup)
		api.DELETE("/groups/delete", controllers.DisbandGroup)
		api.POST("/groups/request-join", controllers.RequestJoinGroup)
	}
}
