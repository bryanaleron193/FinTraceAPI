package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id" example:"1"`
	CreatedAt time.Time      `json:"created_at" example:"2023-09-06T14:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-09-06T14:00:00Z"`
	DeletedAt gorm.DeletedAt `swaggertype:"primitive,string" json:"deleted_at,omitempty" example:"2023-09-06T14:00:00Z"`
}
