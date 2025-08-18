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

		api.GET("/group-members/get-all-roles", controllers.GetAllGroupRoles)
		api.GET("/group-members/get-all-member-statuses", controllers.GetAllGroupMemberStatuses)
		api.GET("/group-members/get-all-members", controllers.GetAllMembers)
		api.PUT("/group-members/update-role", controllers.UpdateGroupRole)
		api.PUT("/group-members/update-member-status", controllers.UpdateGroupMemberStatus)

		api.GET("/transaction-categories/get-list", controllers.GetAllTransactionCategories)

		api.GET("/transactions/get-borrow-list", controllers.GetAllBorrowTransactions)
		api.GET("/transactions/get-lend-list", controllers.GetAllBorrowTransactions)
		api.GET("/transactions/get-detail-by-id", controllers.GetTransactionDetailById)
		api.GET("/transactions/get-borrower-total-by-transaction", controllers.GetBorrowersTotalByTransaction)
		api.POST("/transactions/create", controllers.CreateTransaction)
		api.PUT("/transactions/update", controllers.UpdateTransaction)
		api.DELETE("/transactions/delete", controllers.DeleteTransaction)
	}
}
