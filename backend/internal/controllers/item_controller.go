package controllers

import (
	"net/http"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/services"
	"simple-gin-backend/internal/validators"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetItems godoc
// @Summary List all items
// @Description Get all items in the database
// @Tags Items
// @Accept json
// @Produce json
// @Success 200 {array} models.Item
// @Security BearerAuth
// @Router /items [get]
func GetItems(c *gin.Context) {
	// Retrieve userID from the context (set by the JWT middleware)
	userID := c.MustGet("user_id").(uint)
	items, err := services.GetAllItems(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateItem godoc
// @Summary Create a new item
// @Description Create a new item in the database
// @Tags Items
// @Accept json
// @Produce json
// @Param item body schemas.CreateItemSchemaIn true "Create Item Data"
// @Success 201 {object} models.Item
// @Security BearerAuth
// @Router /items [post]
func CreateItem(c *gin.Context) {
	var input schemas.CreateItemSchemaIn

	userID := c.MustGet("user_id").(uint)

	// Validation
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if _, err := services.AddItem(input, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "item created"})
}

// GetItem godoc
// @Summary Get an item by ID
// @Description Retrieve an item by its ID (protected route)
// @Tags Items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} models.Item
// @Failure 400 {string} string "Invalid item ID"
// @Failure 404 {string} string "Item not found"
// @Security BearerAuth
// @Router /items/{id} [get]
func GetItem(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := strconv.Atoi(idStr)
	userID := c.MustGet("user_id").(uint)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}
	// Validation
	// Validate ownership and item data
	if err := validators.ValidateItemOwnership(userID, uint(itemID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	item, err := services.GetItemByID(uint(itemID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateItem godoc
// @Summary Update an item
// @Description Update an existing item by its ID (protected route)
// @Tags Items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param item body schemas.UpdateItemSchemaIn true "Updated Item Data"
// @Success 200 {string} string "Item updated"
// @Failure 400 {string} string "Invalid item ID" or "Invalid input"
// @Failure 500 {string} string "Internal Server Error"
// @Security BearerAuth
// @Router /items/{id} [put]
func UpdateItem(c *gin.Context) {
	var input schemas.UpdateItemSchemaIn
	idStr := c.Param("id")
	itemID, err := strconv.Atoi(idStr)
	userID := c.MustGet("user_id").(uint)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	// Validation
	// Validate JSON
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	// Validate ownership
	if err := validators.ValidateItemOwnership(userID, uint(itemID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateItem(uint(itemID), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "item updated"})
}

// DeleteItem godoc
// @Summary Delete an item
// @Description Delete an item by its ID (protected route)
// @Tags Items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {string} string "Item deleted"
// @Failure 400 {string} string "Invalid item ID"
// @Failure 500 {string} string "Internal Server Error"
// @Security BearerAuth
// @Router /items/{id} [delete]
func DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := strconv.Atoi(idStr)
	userID := c.MustGet("user_id").(uint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	// Validation
	// Validate ownership
	if err := validators.ValidateItemOwnership(userID, uint(itemID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if err := services.DeleteItem(uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "item deleted"})
}
