package services

import (
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindAllMembersByGroupId(groupID uuid.UUID) ([]models.TrGroupMembers, error) {
	groupMembers := []models.TrGroupMembers{}

	var query = database.DB.
		Model(&models.TrGroupMembers{}).
		Where("group_id = ? AND is_deleted = ?", groupID, false).
		Find(&groupMembers)

	if err := query.Error; err != nil {
		return nil, err
	}

	return groupMembers, nil
}

func GetAllMembers(input *schemas.GroupMemberRequest) ([]schemas.GroupResponse, *schemas.Pagination, error) {
	groups := []schemas.GroupResponse{}
	pagination := new(schemas.Pagination)
	pagination.SetPagination(input.Limit, input.Page)

	// Base query for filtering
	baseQuery := database.DB.
		Table("tr_groups AS a").
		Joins("JOIN tr_group_members AS b ON a.group_id = b.group_id AND b.is_deleted = FALSE").
		Joins("JOIN ms_group_roles AS c ON b.group_role_id = c.group_role_id AND c.is_deleted = FALSE").
		Joins("JOIN ms_group_member_statuses AS d ON b.group_member_status_id = d.group_member_status_id AND d.is_deleted = FALSE").
		Where("b.user_id = ? AND a.is_deleted = ?", auditedUserID, false)

	if len(input.GroupName) == 0 && strings.TrimSpace(input.GroupName) != "" {
		baseQuery = baseQuery.Where("a.group_name ILIKE ?", input.GroupName+"%")
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
			a.group_id,
			a.group_name,
			c.group_role_id,
			c.group_role_name,
			d.group_member_status_id,
			d.group_member_status_name,
			b.approved_at
		`).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Scan(&groups).Error; err != nil {
		return nil, nil, err
	}

	return groups, pagination, nil
}
