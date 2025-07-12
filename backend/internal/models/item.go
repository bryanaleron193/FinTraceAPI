package models

// Item represents a model in the database
type Item struct {
	BaseModel
	Name   string `json:"name" example:"Sample Item"`
	UserID uint   `gorm:"not null" json:"user_id"`
	User   User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
