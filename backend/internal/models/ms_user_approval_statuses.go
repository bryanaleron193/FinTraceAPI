package models

import (
	"github.com/google/uuid"
)

// User represents a user in the system
type MsUserApprovalStatus struct {
	BaseModel

	UserApprovalStatusID   uuid.UUID `gorm:"primaryKey;type:uuid" json:"user_approval_status_id"`
	UserApprovalStatusName string    `gorm:"type:text;not null" json:"user_approval_status_name"`
}
