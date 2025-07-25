package models

import (
	"github.com/google/uuid"
)

type HTrGroups struct {
	BaseModel

	HGroupID  uuid.UUID `gorm:"primaryKey;type:uuid" json:"h_group_id"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null" json:"group_id"`
	GroupName string    `gorm:"type:text;not null" json:"group_name"`
}
