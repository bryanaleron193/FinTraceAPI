package services

import (
	"fmt"
	"log"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindUserByGoogleID(googleID string) (*models.MsUsers, error) {
	user := &models.MsUsers{}

	var query = database.DB.
		Model(&models.MsUsers{}).
		Where("google_id = ?", googleID).
		First(user)

	if err := query.Error; err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByUserID(userID uuid.UUID) (*models.MsUsers, error) {
	user := &models.MsUsers{}

	var query = database.DB.
		Model(&models.MsUsers{}).
		Where("user_id = ?", userID).
		First(user)

	if err := query.Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetAllUsers(input *schemas.UserRequest) ([]schemas.UserResponse, *schemas.Pagination, error) {
	users := []schemas.UserResponse{}
	pagination := new(schemas.Pagination)
	pagination.SetPagination(input.Limit, input.Page)

	// Base query for filtering
	baseQuery := database.DB.
		Table("ms_users AS a").
		Joins("LEFT JOIN ms_user_statuses AS b ON a.user_status_id = b.user_status_id")

	// Apply filters
	if input.UserStatusID != uuid.Nil {
		baseQuery = baseQuery.Where("a.user_status_id = ?", input.UserStatusID)
	}
	if len(input.Name) == 0 && strings.TrimSpace(input.Name) != "" {
		baseQuery = baseQuery.Where("a.name ILIKE ?", input.Name+"%")
	}
	if len(input.Email) == 0 && strings.TrimSpace(input.Email) != "" {
		baseQuery = baseQuery.Where("a.email ILIKE ?", input.Email+"%")
	}

	// Count total matching records
	var total int64
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, nil, err
	}

	pagination.SetTotalPages(total)

	// Fetch paginated data
	if err := baseQuery.
		Select(`
			a.user_id,
			a.email,
			a.name,
			a.user_status_id,
			b.user_status_name,
			a.approved_at
		`).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Scan(&users).Error; err != nil {
		return nil, nil, err
	}

	return users, pagination, nil
}

func CreateUser(input *schemas.AuthRequest) error {
	userID := uuid.New()

	userApprovalStatusID, err := GetUserStatusByName("Pending")
	if err != nil {
		return err
	}

	newUser := models.MsUsers{
		BaseModel: models.BaseModel{
			UserIn:    userID,
			UserUp:    nil,
			DateIn:    time.Now(),
			DateUp:    nil,
			IsDeleted: false,
		},
		UserID:       userID,
		GoogleID:     input.GoogleID,
		Email:        input.Email,
		Name:         input.Name,
		UserStatusID: userApprovalStatusID,
		ApprovedAt:   nil,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		return fmt.Errorf("error creating account: %v", err)
	}

	return nil
}

func IsUserDataChanged(user *models.MsUsers, input *schemas.AuthRequest) bool {
	return strings.TrimSpace(user.Email) != strings.TrimSpace(input.Email) ||
		strings.TrimSpace(user.Name) != strings.TrimSpace(input.Name)
}

func UpdateUserProfile(user *models.MsUsers, input *schemas.AuthRequest) error {
	if err := InsertHistoryUser(user); err != nil {
		return err
	}

	query := database.DB.Model(&models.MsUsers{}).
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

func UpdateUserStatus(auditedUserID uuid.UUID, input *schemas.UserStatusRequest) error {
	existingUser, err := FindUserByUserID(input.UserID)
	if err != nil {
		return fmt.Errorf("error finding user: %v", err)
	}

	err = InsertHistoryUser(existingUser)
	if err != nil {
		return fmt.Errorf("error inserting history user: %v", err)
	}

	userApprovalStatusName, err := GetUserStatusByID(input.UserStatusID)
	if err != nil {
		return fmt.Errorf("error getting user status by ID: %v", err)
	}

	timeNow := time.Now()
	var approvedAt *time.Time
	if userApprovalStatusName == "Approved" {
		approvedAt = &timeNow
	} else {
		approvedAt = nil
	}

	query := database.DB.Model(&models.MsUsers{}).
		Where("user_id = ?", input.UserID).
		Updates(map[string]interface{}{
			"user_up":        auditedUserID,
			"date_up":        timeNow,
			"user_status_id": input.UserStatusID,
			"approved_at":    approvedAt,
		})

	if err := query.Error; err != nil {
		return fmt.Errorf("error updating user status: %v", err)
	}

	err = SendEmail([]string{existingUser.Email}, "Account Status Updated", fmt.Sprintf("Your account status has been updated to: %s", userApprovalStatusName))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email notification")
	}

	return nil
}
