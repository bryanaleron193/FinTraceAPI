package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"

	"github.com/google/uuid"
)

func InsertHistoryUser(user *models.MsUsers) error {
	historyUser := models.HMsUsers{
		BaseModel: models.BaseModel{
			UserIn:    user.UserIn,
			UserUp:    user.UserUp,
			DateIn:    user.DateIn,
			DateUp:    user.DateUp,
			IsDeleted: user.IsDeleted,
		},
		HUserID:      uuid.New(),
		UserID:       user.UserID,
		GoogleID:     user.GoogleID,
		Email:        user.Email,
		Name:         user.Name,
		UserStatusID: user.UserStatusID,
		ApprovedAt:   user.ApprovedAt,
	}

	if err := database.DB.Create(&historyUser).Error; err != nil {
		return fmt.Errorf("error inserting history user: %v", err)
	}

	return nil
}

func InsertHistoryGroup(group *models.TrGroups) error {
	historyGroup := models.HTrGroups{
		BaseModel: models.BaseModel{
			UserIn:    group.UserIn,
			UserUp:    group.UserUp,
			DateIn:    group.DateIn,
			DateUp:    group.DateUp,
			IsDeleted: group.IsDeleted,
		},
		HGroupID:  uuid.New(),
		GroupID:   group.GroupID,
		GroupName: group.GroupName,
	}

	if err := database.DB.Create(&historyGroup).Error; err != nil {
		return fmt.Errorf("error inserting history group: %v", err)
	}

	return nil
}

func InsertHistoryGroupMember(member *models.TrGroupMembers) error {
	historyGroupMember := models.HTrGroupMembers{
		BaseModel: models.BaseModel{
			UserIn:    member.UserIn,
			UserUp:    member.UserUp,
			DateIn:    member.DateIn,
			DateUp:    member.DateUp,
			IsDeleted: member.IsDeleted,
		},
		HGroupMemberID:      uuid.New(),
		GroupMemberID:       member.GroupMemberID,
		GroupID:             member.GroupID,
		UserID:              member.UserID,
		GroupRoleID:         member.GroupRoleID,
		GroupMemberStatusID: member.GroupMemberStatusID,
		ApprovedAt:          member.ApprovedAt,
	}

	if err := database.DB.Create(&historyGroupMember).Error; err != nil {
		return fmt.Errorf("error inserting history group member: %v", err)
	}

	return nil
}

func InsertHistoryGroupMembers(members []models.TrGroupMembers) error {
	var historyGroupMembers []models.HTrGroupMembers

	for _, user := range members {
		historyGroupMembers = append(historyGroupMembers, models.HTrGroupMembers{
			BaseModel: models.BaseModel{
				UserIn:    user.UserIn,
				UserUp:    user.UserUp,
				DateIn:    user.DateIn,
				DateUp:    user.DateUp,
				IsDeleted: user.IsDeleted,
			},
			HGroupMemberID:      uuid.New(),
			GroupMemberID:       user.GroupMemberID,
			GroupID:             user.GroupID,
			UserID:              user.UserID,
			GroupRoleID:         user.GroupRoleID,
			GroupMemberStatusID: user.GroupMemberStatusID,
			ApprovedAt:          user.ApprovedAt,
		})
	}

	if err := database.DB.Create(&historyGroupMembers).Error; err != nil {
		return fmt.Errorf("error inserting history group members: %v", err)
	}

	return nil
}
