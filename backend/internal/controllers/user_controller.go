package controllers

import (
	"net/http"
	"simple-gin-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	res, err := services.GetAllUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
