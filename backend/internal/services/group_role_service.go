package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"

	"github.com/google/uuid"
)

func GetGroupRoleByName(groupRoleName string) (uuid.UUID, error) {
	var result string

	var query = database.DB.
		Model(&models.MsGroupRoles{}).
		Select("group_role_id").
		Where("group_role_name = ? and is_deleted = ?", groupRoleName, false).
		Limit(1).
		Scan(&result)

	if err := query.Error; err != nil {
		return uuid.Nil, err
	}

	groupRoleID, err := uuid.Parse(result)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing user status ID: %v", err)
	}

	return groupRoleID, nil
}

func GetAllGroupRoles() ([]schemas.GroupRoleResponse, error) {
	groupRoles := []schemas.GroupRoleResponse{}

	query := database.DB.
		Model(&models.MsGroupRoles{}).
		Select("group_role_id, group_role_name").
		Where("is_deleted = ?", false).
		Scan(&groupRoles)

	if err := query.Error; err != nil {
		return nil, err
	}

	return groupRoles, nil
}
