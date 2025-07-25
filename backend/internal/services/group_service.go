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

func FindGroupByGroupId(groupID uuid.UUID) (*models.TrGroups, error) {
	group := &models.TrGroups{}

	var query = database.DB.
		Model(&models.TrGroups{}).
		Where("group_id = ? AND is_deleted = ?", groupID, false).
		First(group)

	if err := query.Error; err != nil {
		return nil, err
	}

	return group, nil
}

func GetAllGroupsByUserId(auditedUserID uuid.UUID, input *schemas.GroupRequest) ([]schemas.GroupResponse, *schemas.Pagination, error) {
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

func GetAllGroupsNotJoined(auditedUserID uuid.UUID, input *schemas.GroupRequest) ([]schemas.GroupResponse, *schemas.Pagination, error) {
	groups := []schemas.GroupResponse{}
	pagination := new(schemas.Pagination)
	pagination.SetPagination(input.Limit, input.Page)

	// Base query for filtering
	baseQuery := database.DB.
		Table("tr_groups AS a").
		Where("a.is_deleted = FALSE").
		Where("NOT EXISTS ("+
			"SELECT 1 FROM tr_group_members AS b "+
			"WHERE a.group_id = b.group_id AND b.user_id = ? AND b.is_deleted = FALSE"+
			")", auditedUserID)

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
			a.group_name
		`).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Scan(&groups).Error; err != nil {
		return nil, nil, err
	}

	return groups, pagination, nil
}

func CreateGroup(auditedUserID uuid.UUID, input *schemas.GroupRequest) error {
	timeNow := time.Now()
	groupID := uuid.New()

	groupRoleID, err := GetGroupRoleByName("Admin")
	if err != nil {
		return err
	}

	groupMemberStatusID, err := GetGroupMemberStatusByName("Approved")
	if err != nil {
		return err
	}

	newGroup := models.TrGroups{
		BaseModel: models.BaseModel{
			UserIn:    auditedUserID,
			UserUp:    nil,
			DateIn:    timeNow,
			DateUp:    nil,
			IsDeleted: false,
		},
		GroupID:   groupID,
		GroupName: input.GroupName,
	}

	if err := database.DB.Create(&newGroup).Error; err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}

	newGroupMember := models.TrGroupMembers{
		BaseModel: models.BaseModel{
			UserIn:    auditedUserID,
			UserUp:    nil,
			DateIn:    timeNow,
			DateUp:    nil,
			IsDeleted: false,
		},
		GroupMemberID:       uuid.New(),
		GroupID:             groupID,
		UserID:              auditedUserID,
		GroupRoleID:         groupRoleID,
		GroupMemberStatusID: groupMemberStatusID,
		ApprovedAt:          &timeNow,
	}

	if err := database.DB.Create(&newGroupMember).Error; err != nil {
		return fmt.Errorf("error inserting group member: %v", err)
	}

	return nil
}

func UpdateGroup(auditedUserID uuid.UUID, input *schemas.GroupRequest) error {
	existingGroup, err := FindGroupByGroupId(input.GroupID)
	if err != nil {
		return fmt.Errorf("error finding group: %v", err)
	}

	if err := InsertHistoryGroup(existingGroup); err != nil {
		return err
	}

	query := database.DB.Model(&models.TrGroups{}).
		Where("group_id = ? AND is_deleted = ?", input.GroupID, false).
		Updates(map[string]interface{}{
			"user_up":    auditedUserID,
			"date_up":    time.Now(),
			"group_name": input.GroupName,
		})

	if err := query.Error; err != nil {
		return fmt.Errorf("error updating group status: %v", err)
	}

	return nil
}

func DisbandGroup(auditedUserID uuid.UUID, input *schemas.GroupRequest) error {
	existingGroup, err := FindGroupByGroupId(input.GroupID)
	if err != nil {
		return fmt.Errorf("error finding group: %v", err)
	}

	if err := InsertHistoryGroup(existingGroup); err != nil {
		return err
	}

	existingAllMembers, err := FindAllMembersByGroupId(input.GroupID)
	if err != nil {
		return fmt.Errorf("error finding all members: %v", err)
	}

	if err := InsertHistoryGroupMembers(existingAllMembers); err != nil {
		return err
	}

	timeNow := time.Now()

	query := database.DB.Model(&models.TrGroups{}).
		Where("group_id = ? AND is_deleted = ?", input.GroupID, false).
		Updates(map[string]interface{}{
			"user_up":    auditedUserID,
			"date_up":    timeNow,
			"is_deleted": true,
		})

	if err := query.Error; err != nil {
		return fmt.Errorf("error deleting group: %v", err)
	}

	query = database.DB.Model(&models.TrGroupMembers{}).
		Where("group_id = ? AND is_deleted = ?", input.GroupID, false).
		Updates(map[string]interface{}{
			"user_up":    auditedUserID,
			"date_up":    timeNow,
			"is_deleted": true,
		})

	if err := query.Error; err != nil {
		return fmt.Errorf("error deleting group members: %v", err)
	}

	return nil
}

func RequestJoinGroup(auditedUserID uuid.UUID, input *schemas.GroupRequest) error {
	groupRoleID, err := GetGroupRoleByName("Member")
	if err != nil {
		return err
	}

	groupMemberStatusID, err := GetGroupMemberStatusByName("Pending")
	if err != nil {
		return err
	}

	requestJoinMember := models.TrGroupMembers{
		BaseModel: models.BaseModel{
			UserIn:    auditedUserID,
			UserUp:    nil,
			DateIn:    time.Now(),
			DateUp:    nil,
			IsDeleted: false,
		},
		GroupMemberID:       uuid.New(),
		GroupID:             input.GroupID,
		UserID:              auditedUserID,
		GroupRoleID:         groupRoleID,
		GroupMemberStatusID: groupMemberStatusID,
		ApprovedAt:          nil,
	}

	if err := database.DB.Create(&requestJoinMember).Error; err != nil {
		return fmt.Errorf("error requesting join group: %v", err)
	}

	return nil
}
