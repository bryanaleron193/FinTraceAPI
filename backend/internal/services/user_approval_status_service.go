package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"

	"github.com/google/uuid"
)

func GetUserApprovalStatusByName(statusName string) (uuid.UUID, error) {
	userApprovalStatusID := uuid.Nil

	var query = database.DB.
		Model(&models.MsUserApprovalStatus{}).
		Select("user_approval_status_id").
		Where("user_approval_status_name = ? and is_deleted = ?", statusName, false).
		Limit(1).
		Scan(&userApprovalStatusID)

	// Find the user in the database by email
	if err := query.Error; err != nil {
		return uuid.Nil, fmt.Errorf("user approval status not found")
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
		return "", fmt.Errorf("user approval status not found")
	}

	return userApprovalStatusName, nil
}
