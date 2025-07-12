package factories

import (
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/tests/testutils"

	"golang.org/x/crypto/bcrypt"
)

// CreateUserFactory generates a user and saves it to the test database
// Email: random, Password: "password"
func UserFactory() models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := models.User{
		FirstName:    testutils.GenerateRandomString(10),
		LastName:     testutils.GenerateRandomString(10),
		Email:        testutils.GenerateRandomString(10) + "@example.com",
		PasswordHash: string(hashedPassword),
	}
	testutils.TestDB.Create(&user)
	testutils.TestDB.First(&user, user.ID)
	return user
}
