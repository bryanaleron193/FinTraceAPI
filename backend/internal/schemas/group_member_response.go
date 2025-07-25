package schemas

import (
	"time"

	"github.com/google/uuid"
)

type GroupMemberResponse struct {
	GroupMemberId         uuid.UUID `json:"group_member_id"`
	MemberEmail           string    `json:"member_email"`
	MemberName            string    `json:"member_name"`
	GroupRoleId           uuid.UUID `json:"group_role_id"`
	GroupRoleName         string    `json:"group_role_name"`
	GroupMemberStatusId   uuid.UUID `json:"group_member_status_id"`
	GroupMemberStatusName string    `json:"group_member_status_name"`
	ApprovedAt            time.Time `json:"approved_at,omitempty"`
}
