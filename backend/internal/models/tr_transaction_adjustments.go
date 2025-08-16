package models

import (
	"github.com/google/uuid"
)

type TrTransactionAdjustments struct {
	BaseModel

	TransactionAdjustmentID uuid.UUID `gorm:"primaryKey;type:uuid" json:"transaction_adjustment_id"`
	TransactionHeaderID     uuid.UUID `gorm:"type:uuid" json:"transaction_header_id"`
	AdjustmentName          string    `gorm:"type:text;not null" json:"adjustment_name"`
	AdjustmentAmount        float64   `gorm:"type:numeric(12,2);not null" json:"adjustment_amount"`
}
