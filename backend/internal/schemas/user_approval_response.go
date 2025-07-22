package schemas

import "github.com/google/uuid"

type UserApprovalResponse struct {
	UserApprovalStatusID   uuid.UUID `json:"user_approval_status_id" binding:"required"`
	UserApprovalStatusName string    `json:"user_approval_status_name" binding:"required"`
}
