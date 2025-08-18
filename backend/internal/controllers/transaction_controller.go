package controllers

import (
	"net/http"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllTransactionCategories(c *gin.Context) {
	res, err := services.GetAllTransactionCategories()

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

func GetAllBorrowTransactions(c *gin.Context) {
	input := new(schemas.BorrowTransactionRequest)

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

	res, pagination, err := services.GetAllBorrowTransactions(auditedUserID, input)

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

func GetAllLendTransactions(c *gin.Context) {
	input := new(schemas.LendTransactionRequest)

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

	res, pagination, err := services.GetAllLendTransactions(auditedUserID, input)

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

func GetTransactionDetailById(c *gin.Context) {
	input := new(schemas.TransactionDetailRequest)

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	res, err := services.GetTransactionDetailById(input)

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

func GetBorrowersTotalByTransaction(c *gin.Context) {
	input := new(schemas.TransactionDetailRequest)

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	res, err := services.GetBorrowersTotalByTransaction(input)

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

func CreateTransaction(c *gin.Context) {
	input := new(schemas.CreateTransactionRequest)

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

	err := services.CreateTransaction(auditedUserID, input)

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
	})
}

func UpdateTransaction(c *gin.Context) {
	input := new(schemas.UpdateTransactionRequest)

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

	err := services.UpdateTransaction(auditedUserID, input)

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
	})
}

func DeleteTransaction(c *gin.Context) {
	input := new(schemas.DeleteTransactionRequest)

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

	err := services.DeleteTransaction(auditedUserID, input)

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
	})
}
