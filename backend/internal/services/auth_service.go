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

func AuthenticateUser(input *schemas.UserRequest) (*schemas.AuthResponse, error) {
	user, err := FindUserByGoogleID(input.GoogleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := CreateUser(input); err != nil {
				return nil, err
			}

			return &schemas.AuthResponse{Message: "Account successfully created. Please wait for approval."}, nil
		}

		return nil, err
	}

	if IsUserDataChanged(user, input) {
		if err := UpdateUser(user, input); err != nil {
			return nil, err
		}
	}

	approvalStatus, err := GetUserApprovalStatusByID(user.UserApprovalStatusID)
	if err != nil {
		return nil, err
	}

	if approvalStatus == "Waiting For Approval" {
		return &schemas.AuthResponse{Message: "Your account has not been approved yet. Please wait for approval."}, nil
	}

	if approvalStatus == "Rejected" {
		return &schemas.AuthResponse{Message: "Your account has been rejected."}, nil
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
