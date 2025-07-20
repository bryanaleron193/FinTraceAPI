package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"
	"strings"
	"time"

	"github.com/google/uuid"
)

func FindUserByGoogleID(googleID string) (*models.MsUser, error) {
	user := &models.MsUser{}

	var query = database.DB.
		Model(&models.MsUser{}).
		Where("google_id = ?", googleID).
		First(user)

	// Find the user in the database by email
	if err := query.Error; err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByUserID(userID uuid.UUID) (*models.MsUser, error) {
	user := &models.MsUser{}

	var query = database.DB.
		Model(&models.MsUser{}).
		Where("user_id = ?", userID).
		First(user)

	// Find the user in the database by email
	if err := query.Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetAllUsers() (*[]schemas.UserResponse, error) {
	users := new([]schemas.UserResponse)

	query := database.DB.
		Table("ms_users AS a").
		Select(`
			a.user_id,
			a.email,
			a.name,
			a.user_approval_status_id,
			b.user_approval_status_name,
			a.approved_at
		`).
		Joins("LEFT JOIN ms_user_approval_statuses AS b ON a.user_approval_status_id = b.user_approval_status_id").
		Scan(users)

	if err := query.Error; err != nil {
		return nil, err
	}

	return users, nil
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
	return strings.TrimSpace(user.Email) != strings.TrimSpace(input.Email) ||
		strings.TrimSpace(user.Name) != strings.TrimSpace(input.Name)
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
