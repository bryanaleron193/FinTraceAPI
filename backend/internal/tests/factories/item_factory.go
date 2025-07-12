package factories

import (
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/tests/testutils"
)

// ItemFactory generates an item and saves it to the test database
func ItemFactory(userID uint) models.Item {
	item := models.Item{
		Name:   testutils.GenerateRandomString(10),
		UserID: userID,
	}
	testutils.TestDB.Create(&item)
	testutils.TestDB.First(&item, item.ID)
	return item
}
