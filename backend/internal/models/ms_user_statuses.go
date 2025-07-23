package models

import (
	"github.com/google/uuid"
)

type MsUserStatuses struct {
	BaseModel

	UserStatusID   uuid.UUID `gorm:"primaryKey;type:uuid" json:"user_status_id"`
	UserStatusName string    `gorm:"type:text;not null" json:"user_status_name"`
}
