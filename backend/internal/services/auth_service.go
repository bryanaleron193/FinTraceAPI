package services

import (
	"errors"
	"fmt"
	"log"
	"simple-gin-backend/internal/cache"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/utils"
	"time"

	"gorm.io/gorm"
)

func AuthenticateUser(input *schemas.AuthRequest) (*schemas.AuthResponse, error) {
	user, err := FindUserByGoogleID(input.GoogleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := CreateUser(input); err != nil {
				return nil, err
			}

			err := SendEmail([]string{input.Email}, "Account Created", "Your account has been successfully created. Please wait for approval.")

			if err != nil {
				return nil, fmt.Errorf("failed to send email notification: %v", err)
			}

			return &schemas.AuthResponse{Message: "Account successfully created. Please wait for approval."}, nil
		}

		return nil, err
	}

	if IsUserDataChanged(user, input) {
		if err := UpdateUserProfile(user, input); err != nil {
			return nil, err
		}
	}

	approvalStatus, err := GetUserStatusByID(user.UserStatusID)
	if err != nil {
		return nil, err
	}

	if approvalStatus == "Pending" {
		return &schemas.AuthResponse{Message: "Your account has not been approved yet. Please wait for approval."}, nil
	}

	if approvalStatus == "Rejected" {
		return &schemas.AuthResponse{Message: "Your account has been rejected."}, nil
	}

	if approvalStatus == "Removed" {
		return &schemas.AuthResponse{Message: "Your account has been removed."}, nil
	}

	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	err = cache.RedisClient.Set(cache.Ctx, user.UserID.String(), "exists", time.Hour).Err()
	if err != nil {
		log.Printf("Failed to set user in Redis: %v", err)
	}

	return &schemas.AuthResponse{Token: token}, nil
}
