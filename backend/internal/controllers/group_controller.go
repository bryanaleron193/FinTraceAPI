package controllers

import (
	"net/http"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllGroupsByUserId(c *gin.Context) {
	input := new(schemas.GroupRequest)

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is missing",
		})
		return
	}

	auditedUserID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is not a valid UUID",
		})
		return
	}

	res, pagination, err := services.GetAllGroupsByUserId(auditedUserID, input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.BaseResponse{
		Code:       http.StatusOK,
		Message:    "Success",
		Data:       res,
		Pagination: pagination,
	})
}

func GetAllGroupsNotJoined(c *gin.Context) {
	input := new(schemas.GroupRequest)

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is missing",
		})
		return
	}

	auditedUserID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is not a valid UUID",
		})
		return
	}

	res, pagination, err := services.GetAllGroupsNotJoined(auditedUserID, input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.BaseResponse{
		Code:       http.StatusOK,
		Message:    "Success",
		Data:       res,
		Pagination: pagination,
	})
}

func CreateGroup(c *gin.Context) {
	input := new(schemas.GroupRequest)

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is missing",
		})
		return
	}

	auditedUserID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is not a valid UUID",
		})
		return
	}

	if err := services.CreateGroup(auditedUserID, input); err != nil {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.BaseResponse{
		Code:    http.StatusOK,
		Message: "Group successfully created.",
	})
}

func UpdateGroup(c *gin.Context) {
	input := new(schemas.GroupRequest)

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is missing",
		})
		return
	}

	auditedUserID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is not a valid UUID",
		})
		return
	}

	if err := services.UpdateGroup(auditedUserID, input); err != nil {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.BaseResponse{
		Code:    http.StatusOK,
		Message: "Group successfully updated.",
	})
}

func DisbandGroup(c *gin.Context) {
	input := new(schemas.GroupRequest)

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is missing",
		})
		return
	}

	auditedUserID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is not a valid UUID",
		})
		return
	}

	if err := services.DisbandGroup(auditedUserID, input); err != nil {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.BaseResponse{
		Code:    http.StatusOK,
		Message: "Group successfully disbanded.",
	})
}

func RequestJoinGroup(c *gin.Context) {
	input := new(schemas.GroupRequest)

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is missing",
		})
		return
	}

	auditedUserID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: "User ID is not a valid UUID",
		})
		return
	}

	if err := services.RequestJoinGroup(auditedUserID, input); err != nil {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.BaseResponse{
		Code:    http.StatusOK,
		Message: "Request join group successfully created.",
	})
}
