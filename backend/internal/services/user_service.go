package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user in the database using the UserSchemaIn struct
func CreateUser(input schemas.UserSchemaIn) (*models.User, error) {
	// Check if the email is already registered
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	// Create a new User model instance and set fields
	newUser := models.User{
		Email:        input.Email,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		PasswordHash: string(hashedPassword),
	}

	// Save the user to the database
	if err := database.DB.Create(&newUser).Error; err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return &newUser, nil
}

// AuthenticateUser checks user credentials and returns a JWT if successful
func AuthenticateUser(input schemas.UserLoginSchemaIn) (string, error) {
	var user models.User

	// Find the user in the database by email
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", fmt.Errorf("incorrect password")
	}

	// Generate JWT token if authentication is successful
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return token, nil
}
