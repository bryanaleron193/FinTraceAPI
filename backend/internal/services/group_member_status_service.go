package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"

	"github.com/google/uuid"
)

func GetGroupMemberStatusByName(groupMemberStatusName string) (uuid.UUID, error) {
	var result string

	var query = database.DB.
		Model(&models.MsGroupMemberStatuses{}).
		Select("group_member_status_id").
		Where("group_member_status_name = ? and is_deleted = ?", groupMemberStatusName, false).
		Limit(1).
		Scan(&result)

	if err := query.Error; err != nil {
		return uuid.Nil, err
	}

	groupMemberStatusID, err := uuid.Parse(result)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing group member status ID: %v", err)
	}

	return groupMemberStatusID, nil
}

func GetAllGroupMemberStatuses() ([]schemas.GroupMemberStatusResponse, error) {
	groupMemberStatuses := []schemas.GroupMemberStatusResponse{}

	query := database.DB.
		Model(&models.MsGroupMemberStatuses{}).
		Select("group_member_status_id, group_member_status_name").
		Where("is_deleted = ?", false).
		Scan(&groupMemberStatuses)

	if err := query.Error; err != nil {
		return nil, err
	}

	return groupMemberStatuses, nil
}
