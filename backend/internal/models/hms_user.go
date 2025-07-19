package models

import (
	"time"

	"github.com/google/uuid"
)

type HMsUser struct {
	BaseModel

	HUserID              uuid.UUID  `gorm:"primaryKey;type:uuid" json:"h_user_id"`
	UserID               uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	GoogleID             string     `gorm:"type:text;not null" json:"google_id"`
	Email                string     `gorm:"type:text;not null" json:"email"`
	Name                 string     `gorm:"type:text;not null" json:"name"`
	UserApprovalStatusID uuid.UUID  `gorm:"type:uuid;not null" json:"user_approval_status_id"`
	ApprovedAt           *time.Time `gorm:"type:timestamp with time zone" json:"approved_at"`
}
