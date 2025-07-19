package models

import (
	"time"

	"github.com/google/uuid"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	UserIn    uuid.UUID  `json:"user_in"`
	UserUp    *uuid.UUID `json:"user_up"`
	DateIn    time.Time  `json:"date_in"`
	DateUp    *time.Time `json:"date_up"`
	IsDeleted bool       `json:"is_deleted"`
}
