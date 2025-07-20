package controllers

import (
	"net/http"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllUserApprovalStatuses(c *gin.Context) {
	res, err := services.GetAllUserApprovalStatuses()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpdateUserApprovalStatus(c *gin.Context) {
	input := new(schemas.UserApprovalRequest)
	auditedUserID := c.MustGet("user_id").(uuid.UUID)

	// Validation
	// Validate JSON
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := services.UpdateUserApprovalStatus(auditedUserID, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user approval status successfully updated."})
}
