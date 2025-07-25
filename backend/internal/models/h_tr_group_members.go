package models

import (
	"time"

	"github.com/google/uuid"
)

type HTrGroupMembers struct {
	BaseModel

	HGroupMemberID      uuid.UUID  `gorm:"primaryKey;type:uuid" json:"h_group_member_id"`
	GroupMemberID       uuid.UUID  `gorm:"type:uuid;not null" json:"group_member_id"`
	GroupID             uuid.UUID  `gorm:"type:uuid;not null" json:"group_id"`
	UserID              uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	GroupRoleID         uuid.UUID  `gorm:"type:uuid;not null" json:"group_role_id"`
	GroupMemberStatusID uuid.UUID  `gorm:"type:uuid;not null" json:"group_member_status_id"`
	ApprovedAt          *time.Time `gorm:"type:timestamp with time zone" json:"approved_at"`
}
