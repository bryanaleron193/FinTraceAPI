package schemas

import (
	"github.com/google/uuid"
)

type UserRequest struct {
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	UserStatusID uuid.UUID `json:"user_member_status_id"`
	Limit        int       `json:"limit"`
	Page         int       `json:"page"`
}
