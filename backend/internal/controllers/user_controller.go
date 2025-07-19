package controllers

import (
	"net/http"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// LoginUser godoc
// @Summary Log in a user
// @Description Log in by providing email and password to receive a JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body schemas.UserLoginSchemaIn true "User Login Data"
// @Success 200 {string} string "JWT Token"
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
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
	}

	if res.Message != "" {
		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	}

	c.JSON(http.StatusOK, gin.H{"token": res.Token})
}
