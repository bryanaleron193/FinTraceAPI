package schemas

import "github.com/google/uuid"

type GroupMemberRequest struct {
	UserID              uuid.UUID `json:"user_id"`
	GroupID             uuid.UUID `json:"group_id"`
	GroupMemberStatusID uuid.UUID `json:"group_member_status_id"`
	GroupRoleID         uuid.UUID `json:"group_role_id"`
}
