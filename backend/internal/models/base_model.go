package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	UserIn    uuid.UUID  `json:"user_in"`
	UserUp    *uuid.UUID `json:"user_up"`
	DateIn    time.Time  `json:"date_in"`
	DateUp    *time.Time `json:"date_up"`
	IsDeleted bool       `json:"is_deleted"`
}
