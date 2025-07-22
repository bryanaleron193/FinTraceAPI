package schemas

import (
	"github.com/google/uuid"
)

type UserRequest struct {
	Email                string    `json:"email"`
	Name                 string    `json:"name"`
	UserApprovalStatusID uuid.UUID `json:"user_approval_status_id"`
	Limit                int       `json:"limit"`
	Page                 int       `json:"page"`
}
