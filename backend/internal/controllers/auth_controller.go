package controllers

import (
	"net/http"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func LoginUser(c *gin.Context) {
	input := new(schemas.UserRequest)

	// Bind the JSON input to the UserLoginSchemaIn DTO
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Authenticate the user and generate JWT token using the service
	res, err := services.AuthenticateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.Message != "" {
		c.JSON(http.StatusOK, gin.H{"message": res.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": res.Token})
}
