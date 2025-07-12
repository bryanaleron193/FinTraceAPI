package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"
)

// GetAllItems retrieves items from the database, optionally filtering by user ID
func GetAllItems(userId *uint) ([]models.Item, error) {
	var items []models.Item

	// If userId is provided (not nil), filter by userId
	if userId != nil {
		result := database.DB.Where("user_id = ?", *userId).Find(&items)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		// If no userId is provided, return all items
		result := database.DB.Find(&items)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return items, nil
}

// Get item by ID
func GetItemByID(id uint) (models.Item, error) {
	var item models.Item
	if err := database.DB.First(&item, id).Error; err != nil {
		return item, err
	}
	return item, nil
}

// Add a new item
func AddItem(input schemas.CreateItemSchemaIn, userID uint) (*models.Item, error) {
	// Create a new User model instance and set fields
	newItem := models.Item{
		Name:   input.Name,
		UserID: userID,
	}
	// Save the item to the database
	if err := database.DB.Create(&newItem).Error; err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}
	return &newItem, nil
}

// Update an item by a User
func UpdateItem(id uint, input schemas.UpdateItemSchemaIn) error {
	var item models.Item
	if err := database.DB.First(&item, id).Error; err != nil {
		return err
	}
	item.Name = input.Name
	return database.DB.Save(&item).Error
}

// Delete an item
func DeleteItem(id uint) error {
	return database.DB.Delete(&models.Item{}, id).Error
}
