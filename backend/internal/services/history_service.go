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

func InsertHistoryGroup(user *models.TrGroups) error {
	historyGroup := models.HTrGroups{
		BaseModel: models.BaseModel{
			UserIn:    user.UserIn,
			UserUp:    user.UserUp,
			DateIn:    user.DateIn,
			DateUp:    user.DateUp,
			IsDeleted: user.IsDeleted,
		},
		HGroupID:  uuid.New(),
		GroupID:   user.GroupID,
		GroupName: user.GroupName,
	}

	if err := database.DB.Create(&historyGroup).Error; err != nil {
		return fmt.Errorf("error inserting history group: %v", err)
	}

	return nil
}

func InsertHistoryGroupMembers(members []models.TrGroupMembers) error {
	var historyGroups []models.HTrGroupMembers

	for _, user := range members {
		historyGroups = append(historyGroups, models.HTrGroupMembers{
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

	if err := database.DB.Create(&historyGroups).Error; err != nil {
		return fmt.Errorf("error inserting history group members: %v", err)
	}

	return nil
}
