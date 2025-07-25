package models

import (
	"github.com/google/uuid"
)

type MsGroupMemberStatuses struct {
	BaseModel

	GroupMemberStatusID   uuid.UUID `gorm:"primaryKey;type:uuid" json:"group_member_status_id"`
	GroupMemberStatusName string    `gorm:"type:text;not null" json:"group_member_status_name"`
}
