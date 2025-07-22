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
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}

func UpdateUserApprovalStatus(c *gin.Context) {
	input := new(schemas.UserApprovalRequest)
	auditedUserID := c.MustGet("user_id").(uuid.UUID)

	// Validation
	// Validate JSON
	if err := c.BindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	if err := services.UpdateUserApprovalStatus(auditedUserID, input); err != nil {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.BaseResponse{
		Code:    http.StatusOK,
		Message: "user approval status successfully updated.",
	})
}
