package controllers

import (
	"net/http"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func LoginUser(c *gin.Context) {
	input := new(schemas.AuthRequest)

	// Bind the JSON input to the UserLoginSchemaIn DTO
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	// Authenticate the user and generate JWT token using the service
	res, err := services.AuthenticateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	if res.Message != "" {
		c.JSON(http.StatusOK, schemas.BaseResponse{
			Code:    http.StatusOK,
			Message: "Success",
			Data:    res,
		})
		return
	}

	c.JSON(http.StatusOK, schemas.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}
