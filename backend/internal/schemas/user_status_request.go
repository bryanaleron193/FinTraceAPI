package schemas

import "github.com/google/uuid"

type UserStatusRequest struct {
	UserID       uuid.UUID `json:"user_id" binding:"required"`
	UserStatusID uuid.UUID `json:"user_member_status_id" binding:"required"`
}
