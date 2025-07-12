package controllers

import (
	"net/http"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary Register a new user
// @Description Create a new user by providing email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body schemas.UserSchemaIn true "User Data"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Bad request"
// @Router /sign-up [post]
func RegisterUser(c *gin.Context) {
	var input schemas.UserSchemaIn

	// Bind the JSON input to the UserSchemaIn DTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service to create the user, passing the entire input schema
	user, err := services.CreateUser(input)
	if err != nil {
		if err.Error() == "email already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

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
	var input schemas.UserLoginSchemaIn

	// Bind the JSON input to the UserLoginSchemaIn DTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Authenticate the user and generate JWT token using the service
	token, err := services.AuthenticateUser(input)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "incorrect password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
