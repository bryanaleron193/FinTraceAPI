package schemas

import "github.com/google/uuid"

type UserStatusResponse struct {
	UserStatusID   uuid.UUID `json:"user_status_id" binding:"required"`
	UserStatusName string    `json:"user_status_name" binding:"required"`
}
