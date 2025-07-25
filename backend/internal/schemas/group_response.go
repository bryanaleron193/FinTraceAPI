package schemas

import (
	"time"

	"github.com/google/uuid"
)

type GroupResponse struct {
	GroupID               uuid.UUID `json:"group_id"`
	GroupName             string    `json:"group_name"`
	GroupRoleId           uuid.UUID `json:"group_role_id"`
	GroupRoleName         string    `json:"group_role_name"`
	GroupMemberStatusId   uuid.UUID `json:"group_member_status_id"`
	GroupMemberStatusName string    `json:"group_member_status_name"`
	ApprovedAt            time.Time `json:"approved_at,omitempty"`
}
