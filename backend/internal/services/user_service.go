package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"
	"simple-gin-backend/internal/utils"
	"time"

	"github.com/google/uuid"
)

func FindUser(googleID string) (*models.MsUser, error) {
	user := &models.MsUser{}

	var query = database.DB.
		Model(&models.MsUser{}).
		Where("google_id = ? and is_deleted = ?", googleID, false).
		First(user)

	// Find the user in the database by email
	if err := query.Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func CreateUser(input *schemas.UserRequest) error {
	userID := uuid.New()

	userApprovalStatusID, err := GetUserApprovalStatusByName("Waiting For Approval")
	if err != nil {
		return err
	}

	newUser := models.MsUser{
		BaseModel: models.BaseModel{
			UserIn:    userID,
			UserUp:    nil,
			DateIn:    time.Now(),
			DateUp:    nil,
			IsDeleted: false,
		},
		UserID:               userID,
		GoogleID:             input.GoogleID,
		Email:                input.Email,
		Name:                 input.Name,
		UserApprovalStatusID: userApprovalStatusID,
		ApprovedAt:           nil,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		return fmt.Errorf("error creating account: %v", err)
	}

	return nil
}

func IsUserDataChanged(user *models.MsUser, input *schemas.UserRequest) bool {
	return user.Email != input.Email || user.Name != input.Name
}

func InsertHistoryUser(user *models.MsUser) error {
	historyUser := models.HMsUser{
		BaseModel: models.BaseModel{
			UserIn:    user.UserIn,
			UserUp:    user.UserUp,
			DateIn:    user.DateIn,
			DateUp:    user.DateUp,
			IsDeleted: user.IsDeleted,
		},
		HUserID:              uuid.New(),
		UserID:               user.UserID,
		GoogleID:             user.GoogleID,
		Email:                user.Email,
		Name:                 user.Name,
		UserApprovalStatusID: user.UserApprovalStatusID,
		ApprovedAt:           user.ApprovedAt,
	}

	if err := database.DB.Create(&historyUser).Error; err != nil {
		return fmt.Errorf("error inserting history user: %v", err)
	}

	return nil
}

func UpdateUser(user *models.MsUser, input *schemas.UserRequest) error {
	if err := InsertHistoryUser(user); err != nil {
		return err
	}

	query := database.DB.Model(&models.MsUser{}).
		Where("user_id = ?", user.UserID).
		Updates(map[string]interface{}{
			"user_up": user.UserID,
			"date_up": time.Now(),
			"email":   input.Email,
			"name":    input.Name,
		})

	if err := query.Error; err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	return nil
}

func AuthenticateUser(input *schemas.UserRequest) (*schemas.AuthResponse, error) {
	user, err := FindUser(input.GoogleID)
	if err != nil {
		if err := CreateUser(input); err != nil {
			return nil, err
		}

		return &schemas.AuthResponse{Message: "Account successfully created. Please wait for approval."}, nil
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

	return &schemas.AuthResponse{Token: token}, nil
}
