package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"

	"github.com/google/uuid"
)

func GetUserStatusByName(statusName string) (uuid.UUID, error) {
	var result string

	var query = database.DB.
		Model(&models.MsUserStatuses{}).
		Select("user_member_status_id").
		Where("user_member_status_name = ? and is_deleted = ?", statusName, false).
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

func GetUserStatusByID(statusID uuid.UUID) (string, error) {
	userApprovalStatusName := ""

	var query = database.DB.
		Model(&models.MsUserStatuses{}).
		Select("user_member_status_name").
		Where("user_member_status_id = ? and is_deleted = ?", statusID, false).
		Limit(1).
		Scan(&userApprovalStatusName)

	// Find the user in the database by email
	if err := query.Error; err != nil {
		return "", err
	}

	return userApprovalStatusName, nil
}

func GetAllUserStatuses() (*[]schemas.UserStatusResponse, error) {
	userApprovalStatuses := []schemas.UserStatusResponse{}

	query := database.DB.
		Model(&models.MsUserStatuses{}).
		Select("user_member_status_id, user_member_status_name").
		Where("is_deleted = ?", false).
		Scan(&userApprovalStatuses)

	if err := query.Error; err != nil {
		return nil, err
	}

	return &userApprovalStatuses, nil
}
