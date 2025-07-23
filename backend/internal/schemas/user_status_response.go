package schemas

import "github.com/google/uuid"

type UserStatusResponse struct {
	UserStatusID   uuid.UUID `json:"user_member_status_id" binding:"required"`
	UserStatusName string    `json:"user_member_status_name" binding:"required"`
}
