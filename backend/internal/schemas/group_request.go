package schemas

import "github.com/google/uuid"

type GroupRequest struct {
	GroupID   uuid.UUID `json:"group_id"`
	GroupName string    `json:"group_name"`
	Limit     int       `json:"limit"`
	Page      int       `json:"page"`
}
