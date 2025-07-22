package schemas

import "github.com/google/uuid"

type UserApprovalRequest struct {
	UserID               uuid.UUID `json:"user_id" binding:"required"`
	UserApprovalStatusID uuid.UUID `json:"user_approval_status_id" binding:"required"`
}
