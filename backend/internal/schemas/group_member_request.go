package schemas

import "github.com/google/uuid"

type GroupMemberRequest struct {
	GroupID             uuid.UUID `json:"group_id"`
	GroupMemberID       uuid.UUID `json:"group_member_id"`
	MemberName          string    `json:"member_name"`
	GroupRoleID         uuid.UUID `json:"group_role_id"`
	GroupMemberStatusID uuid.UUID `json:"group_member_status_id"`
	Limit               int       `json:"limit"`
	Page                int       `json:"page"`
}
