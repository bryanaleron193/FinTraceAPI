package schemas

import "github.com/google/uuid"

type TransactionCategoryResponse struct {
	TransactionCategoryID   uuid.UUID `json:"transaction_category_id" binding:"required"`
	TransactionCategoryName string    `json:"transaction_category_name" binding:"required"`
}
