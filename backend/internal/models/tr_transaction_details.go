package models

import (
	"time"

	"github.com/google/uuid"
)

type TrTransactionDetails struct {
	BaseModel

	TransactionDetailID uuid.UUID  `gorm:"primaryKey;type:uuid" json:"transaction_detail_id"`
	TransactionHeaderID uuid.UUID  `gorm:"type:uuid" json:"transaction_header_id"`
	BorrowerID          uuid.UUID  `gorm:"type:uuid;not null" json:"borrower_id"`
	BorrowerName        string     `gorm:"type:text;not null" json:"borrower_name"`
	DetailName          string     `gorm:"type:text;not null" json:"detail_name"`
	Amount              float64    `gorm:"type:numeric(12,2);not null" json:"amount"`
	Quantity            int        `gorm:"type:int;not null;default:1;check:quantity > 0"`
	IsPaid              bool       `gorm:"type:boolean;not null" json:"is_paid"`
	PaidAt              *time.Time `gorm:"type:timestamp with time zone" json:"paid_at"`
}
