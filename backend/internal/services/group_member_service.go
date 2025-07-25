package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"
	"simple-gin-backend/internal/schemas"
	"strings"
	"time"

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

func GetAllMembers(input *schemas.GroupMemberRequest) ([]schemas.GroupMemberResponse, *schemas.Pagination, error) {
	groupMembers := []schemas.GroupMemberResponse{}
	pagination := new(schemas.Pagination)
	pagination.SetPagination(input.Limit, input.Page)

	input.MemberName = strings.TrimSpace(input.MemberName)

	// Base query for filtering
	baseQuery := database.DB.
		Table("tr_group_members AS a").
		Joins("JOIN ms_users AS b ON a.user_id = b.user_id AND b.is_deleted = FALSE").
		Joins("JOIN ms_group_roles AS c ON a.group_role_id = c.group_role_id AND c.is_deleted = FALSE").
		Joins("JOIN ms_group_member_statuses AS d ON a.group_member_status_id = d.group_member_status_id AND d.is_deleted = FALSE").
		Where("a.group_id = ? AND a.is_deleted = ?", input.GroupID, false)

	if input.MemberName != "" {
		baseQuery = baseQuery.Where("b.name ILIKE ?", input.MemberName+"%")
	}

	if input.GroupRoleID != uuid.Nil {
		baseQuery = baseQuery.Where("a.group_role_id = ?", input.GroupRoleID)
	}

	if input.GroupMemberStatusID != uuid.Nil {
		baseQuery = baseQuery.Where("a.group_member_status_id = ?", input.GroupMemberStatusID)
	}

	// Count total matching records
	var total int64
	countQuery := baseQuery.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	pagination.SetTotalPages(total)

	// Fetch paginated data
	if err := baseQuery.
		Select(`
			a.group_member_id,
			b.user_id,
			b.email AS member_email,
			b.name AS member_name,
			c.group_role_id,
			c.group_role_name,
			d.group_member_status_id,
			d.group_member_status_name,
			a.approved_at
		`).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Scan(&groupMembers).Error; err != nil {
		return nil, nil, err
	}

	return groupMembers, pagination, nil
}

func FindMemberByGroupMemberId(groupMemberID uuid.UUID) (*models.TrGroupMembers, error) {
	groupMember := new(models.TrGroupMembers)

	var query = database.DB.
		Model(&models.TrGroupMembers{}).
		Where("group_member_id = ? AND is_deleted = ?", groupMemberID, false).
		Limit(1).
		Find(groupMember)

	if err := query.Error; err != nil {
		return nil, err
	}

	return groupMember, nil
}

func UpdateGroupRole(auditedUserID uuid.UUID, input *schemas.GroupMemberRequest) error {
	existingGroupMember, err := FindMemberByGroupMemberId(input.GroupMemberID)
	if err != nil {
		return fmt.Errorf("error finding group member: %v", err)
	}

	if err := InsertHistoryGroupMember(existingGroupMember); err != nil {
		return err
	}

	query := database.DB.Model(&models.TrGroupMembers{}).
		Where("group_member_id = ? AND is_deleted = ?", input.GroupMemberID, false).
		Updates(map[string]interface{}{
			"user_up":       auditedUserID,
			"date_up":       time.Now(),
			"group_role_id": input.GroupRoleID,
		})

	if err := query.Error; err != nil {
		return fmt.Errorf("error updating group role: %v", err)
	}

	return nil
}

func UpdateGroupMemberStatus(auditedUserID uuid.UUID, input *schemas.GroupMemberRequest) error {
	existingGroupMember, err := FindMemberByGroupMemberId(input.GroupMemberID)
	if err != nil {
		return fmt.Errorf("error finding group member: %v", err)
	}

	if err := InsertHistoryGroupMember(existingGroupMember); err != nil {
		return err
	}

	query := database.DB.Model(&models.TrGroupMembers{}).
		Where("group_member_id = ? AND is_deleted = ?", input.GroupMemberID, false).
		Updates(map[string]interface{}{
			"user_up":                auditedUserID,
			"date_up":                time.Now(),
			"group_member_status_id": input.GroupMemberStatusID,
		})

	if err := query.Error; err != nil {
		return fmt.Errorf("error updating group member status: %v", err)
	}

	return nil
}
