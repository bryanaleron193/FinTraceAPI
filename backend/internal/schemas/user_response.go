package schemas

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	UserID         uuid.UUID  `json:"user_id" binding:"required"`
	GoogleID       string     `json:"google_id" binding:"required"`
	Email          string     `json:"email" binding:"required,email"`
	Name           string     `json:"name" binding:"required"`
	UserStatusID   uuid.UUID  `json:"user_member_status_id" binding:"required"`
	UserStatusName string     `json:"user_member_status_name" binding:"required"`
	ApprovedAt     *time.Time `json:"approved_at"`
}
