package controllers

import (
	"net/http"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	input := new(schemas.UserRequest)

	if err := c.BindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, schemas.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	res, pagination, err := services.GetAllUsers(input)

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
