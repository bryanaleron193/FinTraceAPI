package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	UserIn    uuid.UUID  `gorm:"type:uuid;not null" json:"user_in"`
	UserUp    *uuid.UUID `gorm:"type:uuid" json:"user_up"`
	DateIn    time.Time  `gorm:"type:timestamp with time zone;not null" json:"date_in"`
	DateUp    *time.Time `gorm:"type:timestamp with time zone" json:"date_up"`
	IsDeleted bool       `gorm:"type:boolean;not null" json:"is_deleted"`
}
