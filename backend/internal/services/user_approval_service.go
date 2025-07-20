package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"
	"time"

	"github.com/google/uuid"
)

func GetUserApprovalStatusByName(statusName string) (uuid.UUID, error) {
	var result string

	var query = database.DB.
		Model(&models.MsUserApprovalStatus{}).
		Select("user_approval_status_id").
		Where("user_approval_status_name = ? and is_deleted = ?", statusName, false).
		Limit(1).
		Scan(&result)

	// Find the user in the database by email
	if err := query.Error; err != nil {
		return uuid.Nil, err
	}

	userApprovalStatusID, err := uuid.Parse(result)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing user approval status ID: %v", err)
	}

	return userApprovalStatusID, nil
}

func GetUserApprovalStatusByID(statusID uuid.UUID) (string, error) {
	userApprovalStatusName := ""

	var query = database.DB.
		Model(&models.MsUserApprovalStatus{}).
		Select("user_approval_status_name").
		Where("user_approval_status_id = ? and is_deleted = ?", statusID, false).
		Limit(1).
		Scan(&userApprovalStatusName)

	// Find the user in the database by email
	if err := query.Error; err != nil {
		return "", err
	}

	return userApprovalStatusName, nil
}

func GetAllUserApprovalStatuses() (*[]schemas.UserApprovalResponse, error) {
	userApprovalStatuses := new([]schemas.UserApprovalResponse)

	var query = database.DB.
		Model(&models.MsUserApprovalStatus{}).
		Select("user_approval_status_id, user_approval_status_name").
		Where("is_deleted = ?", false).
		Scan(userApprovalStatuses)

	// Find the user in the database by email
	if err := query.Error; err != nil {
		return nil, err
	}

	return userApprovalStatuses, nil
}

func UpdateUserApprovalStatus(auditedUserID uuid.UUID, input *schemas.UserApprovalRequest) error {
	existingUser, err := FindUserByUserID(input.UserID)
	if err != nil {
		return fmt.Errorf("error finding user: %v", err)
	}

	err = InsertHistoryUser(existingUser)
	if err != nil {
		return fmt.Errorf("error inserting history user: %v", err)
	}

	userApprovalStatusName, err := GetUserApprovalStatusByID(input.UserApprovalStatusID)
	if err != nil {
		return fmt.Errorf("error getting user approval status by ID: %v", err)
	}

	timeNow := time.Now()
	var approvedAt *time.Time
	if userApprovalStatusName == "Approved" {
		approvedAt = &timeNow
	} else {
		approvedAt = nil
	}

	query := database.DB.Model(&models.MsUser{}).
		Where("user_id = ?", input.UserID).
		Updates(map[string]interface{}{
			"user_up":                 auditedUserID,
			"date_up":                 timeNow,
			"user_approval_status_id": input.UserApprovalStatusID,
			"approved_at":             approvedAt,
		})

	if err := query.Error; err != nil {
		return fmt.Errorf("error updating user approval status: %v", err)
	}

	return nil
}
