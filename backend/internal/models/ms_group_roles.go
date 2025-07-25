package models

import (
	"github.com/google/uuid"
)

type MsGroupRoles struct {
	BaseModel

	GroupRoleID   uuid.UUID `gorm:"primaryKey;type:uuid" json:"group_role_id"`
	GroupRoleName string    `gorm:"type:text;not null" json:"group_role_name"`
}
