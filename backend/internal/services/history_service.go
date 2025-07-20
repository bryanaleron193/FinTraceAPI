package services

import (
	"fmt"
	"simple-gin-backend/internal/database"
	"simple-gin-backend/internal/models"

	"github.com/google/uuid"
)

func InsertHistoryUser(user *models.MsUser) error {
	historyUser := models.HMsUser{
		BaseModel: models.BaseModel{
			UserIn:    user.UserIn,
			UserUp:    user.UserUp,
			DateIn:    user.DateIn,
			DateUp:    user.DateUp,
			IsDeleted: user.IsDeleted,
		},
		HUserID:              uuid.New(),
		UserID:               user.UserID,
		GoogleID:             user.GoogleID,
		Email:                user.Email,
		Name:                 user.Name,
		UserApprovalStatusID: user.UserApprovalStatusID,
		ApprovedAt:           user.ApprovedAt,
	}

	if err := database.DB.Create(&historyUser).Error; err != nil {
		return fmt.Errorf("error inserting history user: %v", err)
	}

	return nil
}
