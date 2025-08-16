package models

import (
	"time"

	"github.com/google/uuid"
)

type HTrTransactionHeaders struct {
	BaseModel

	HTransactionHeaderID  uuid.UUID `gorm:"primaryKey;type:uuid" json:"h_transaction_header_id"`
	TransactionHeaderID   uuid.UUID `gorm:"type:uuid" json:"transaction_header_id"`
	TransactionCategoryID uuid.UUID `gorm:"type:uuid;not null" json:"transaction_category_id"`
	LenderID              uuid.UUID `gorm:"type:uuid" json:"lender_id"`
	LenderName            string    `gorm:"type:text" json:"lender_name"`
	TransactionName       string    `gorm:"type:text;not null" json:"transaction_name"`
	TransactionDate       time.Time `gorm:"type:timestamp with time zone;not null" json:"transaction_date"`
}
