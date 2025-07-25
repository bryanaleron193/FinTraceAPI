package models

import (
	"github.com/google/uuid"
)

type TrGroups struct {
	BaseModel

	GroupID   uuid.UUID `gorm:"primaryKey;type:uuid" json:"group_id"`
	GroupName string    `gorm:"type:text;not null" json:"group_name"`
}
