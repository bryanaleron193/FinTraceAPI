package schemas

import "github.com/google/uuid"

type GroupMemberStatusResponse struct {
	GroupMemberStatusID   uuid.UUID `json:"group_member_status_id" binding:"required"`
	GroupMemberStatusName string    `json:"group_member_status_name" binding:"required"`
}
