package models

// User represents a user in the system
type User struct {
	BaseModel
	FirstName    string `json:"first_name" example:"John"`
	LastName     string `json:"last_name" example:"Doe"`
	Email        string `gorm:"unique;not null" json:"email" example:"john.doe@example.com"`
	PasswordHash string `gorm:"not null" json:"-"` // Password will be hashed
}
