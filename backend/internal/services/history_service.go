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

func InsertHistoryTransactionHeader(header *models.TrTransactionHeaders) error {
	historyTransactionHeader := models.HTrTransactionHeaders{
		BaseModel: models.BaseModel{
			UserIn:    header.UserIn,
			UserUp:    header.UserUp,
			DateIn:    header.DateIn,
			DateUp:    header.DateUp,
			IsDeleted: header.IsDeleted,
		},
		HTransactionHeaderID:  uuid.New(),
		TransactionHeaderID:   header.TransactionHeaderID,
		TransactionCategoryID: header.TransactionCategoryID,
		LenderID:              header.LenderID,
		LenderName:            header.LenderName,
		TransactionName:       header.TransactionName,
		TransactionDate:       header.TransactionDate,
		IsRoundDown:           header.IsRoundDown,
	}

	if err := database.DB.Create(&historyTransactionHeader).Error; err != nil {
		return fmt.Errorf("error inserting history transaction header: %v", err)
	}

	return nil
}

func InsertHistoryTransactionDetails(details []models.TrTransactionDetails) error {
	var historyTransactionDetails []models.HTrTransactionDetails

	for _, detail := range details {
		historyTransactionDetails = append(historyTransactionDetails, models.HTrTransactionDetails{
			BaseModel: models.BaseModel{
				UserIn:    detail.UserIn,
				UserUp:    detail.UserUp,
				DateIn:    detail.DateIn,
				DateUp:    detail.DateUp,
				IsDeleted: detail.IsDeleted,
			},
			HTransactionDetailID: uuid.New(),
			TransactionDetailID:  detail.TransactionDetailID,
			TransactionHeaderID:  detail.TransactionHeaderID,
			BorrowerID:           detail.BorrowerID,
			BorrowerName:         detail.BorrowerName,
			DetailName:           detail.DetailName,
			Amount:               detail.Amount,
			Quantity:             detail.Quantity,
			IsPaid:               detail.IsPaid,
			PaidAt:               detail.PaidAt,
		})
	}

	if err := database.DB.Create(&historyTransactionDetails).Error; err != nil {
		return fmt.Errorf("error inserting history transaction details: %v", err)
	}

	return nil
}

func InsertHistoryTransactionAdjustments(adjustments []models.TrTransactionAdjustments) error {
	var historyTransactionAdjustments []models.HTrTransactionAdjustments

	for _, adjustment := range adjustments {
		historyTransactionAdjustments = append(historyTransactionAdjustments, models.HTrTransactionAdjustments{
			BaseModel: models.BaseModel{
				UserIn:    adjustment.UserIn,
				UserUp:    adjustment.UserUp,
				DateIn:    adjustment.DateIn,
				DateUp:    adjustment.DateUp,
				IsDeleted: adjustment.IsDeleted,
			},
			HTransactionAdjustmentID: uuid.New(),
			TransactionAdjustmentID:  adjustment.TransactionAdjustmentID,
			TransactionHeaderID:      adjustment.TransactionHeaderID,
			AdjustmentName:           adjustment.AdjustmentName,
			AdjustmentAmount:         adjustment.AdjustmentAmount,
		})
	}

	if err := database.DB.Create(&historyTransactionAdjustments).Error; err != nil {
		return fmt.Errorf("error inserting history transaction adjustments: %v", err)
	}

	return nil
}
