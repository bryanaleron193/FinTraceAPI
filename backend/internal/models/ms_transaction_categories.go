package models

import (
	"github.com/google/uuid"
)

type MsTransactionCategories struct {
	BaseModel

	TransactionCategoryID   uuid.UUID `gorm:"primaryKey;type:uuid" json:"transaction_category_id"`
	TransactionCategoryName string    `gorm:"type:text;not null" json:"transaction_category_name"`
}
