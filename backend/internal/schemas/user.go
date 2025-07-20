package schemas

import (
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	GoogleID string `json:"google_id" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
}

type UserResponse struct {
	UserID                 uuid.UUID  `json:"user_id" binding:"required"`
	GoogleID               string     `json:"google_id" binding:"required"`
	Email                  string     `json:"email" binding:"required,email"`
	Name                   string     `json:"name" binding:"required"`
	UserApprovalStatusID   uuid.UUID  `json:"user_approval_status_id" binding:"required"`
	UserApprovalStatusName string     `json:"user_approval_status_name" binding:"required"`
	ApprovedAt             *time.Time `json:"approved_at"`
}
