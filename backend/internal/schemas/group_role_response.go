package schemas

import "github.com/google/uuid"

type GroupRoleResponse struct {
	GroupRoleID   uuid.UUID `json:"group_role_id" binding:"required"`
	GroupRoleName string    `json:"group_role_name" binding:"required"`
}
