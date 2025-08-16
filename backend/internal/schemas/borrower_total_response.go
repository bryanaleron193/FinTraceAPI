package schemas

import "github.com/google/uuid"

type BorrowerTotalResponse struct {
	BorrowerID   uuid.UUID `gorm:"column:borrower_id"`
	BorrowerName string    `gorm:"column:borrower_name"`
	TotalOwed    float64   `gorm:"column:total_owed"`
}
