package validators

import (
	"errors"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/services"
)

// ValidateItemOwnership checks if the item exists and belongs to the user
func ValidateItemOwnership(userID uint, itemID uint) error {
	var item models.Item
	item, err := services.GetItemByID(itemID)
	if err != nil {
		return errors.New("item not found")
	}

	// Check if the item belongs to the user
	if item.UserID != userID {
		return errors.New("you are not authorized to update this item")
	}

	// Return the item if validation passes
	return nil
}
